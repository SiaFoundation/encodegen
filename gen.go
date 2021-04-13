package main

import (
	"fmt"
	"go/build"
	"go/format"
	"go/token"
	"go/types"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

var siaMarshaler, siaUnmarshaler *types.Interface

func makeInterface(name string, param types.Type) *types.Interface {
	params := types.NewTuple(types.NewVar(token.NoPos, nil, "", param))
	results := types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type()))
	meth := types.NewFunc(token.NoPos, nil, name, types.NewSignature(nil, params, results, false))
	return types.NewInterfaceType([]*types.Func{meth}, nil).Complete()
}

type generator struct {
	pkg     *packages.Package
	typs    map[string]string
	imports []string
}

func (g *generator) addImport(pkg string) {
	pkg = `"` + pkg + `"`
	for _, p := range g.imports {
		if p == pkg {
			return
		}
	}
	g.imports = append(g.imports, pkg)
}

func (g *generator) addImportType(t types.Type) {
	if named, ok := t.(*types.Named); ok {
		if pkg := named.Obj().Pkg(); pkg != g.pkg.Types {
			g.addImport(pkg.Path())
		}
	}
}

func (g *generator) willGenerate(t types.Type) bool {
	if named, ok := t.(*types.Named); ok && named.Obj().Pkg() == g.pkg.Types {
		return g.typs[types.TypeString(t, types.RelativeTo(g.pkg.Types))] != ""
	}
	return false
}

func Generate(dir string, typs ...string) (string, error) {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedTypes | packages.NeedImports}

	// load source package; also load "io", to construct interface types
	pkgs, err := packages.Load(cfg, dir, "io")
	if err != nil {
		return "", err
	}
	// ordering of pkgs is unspecified, so locate packages manually
	srcPkg, ioPkg := pkgs[0], pkgs[1]
	if ioPkg.PkgPath != "io" {
		srcPkg, ioPkg = ioPkg, srcPkg
	}

	// construct interface types
	siaMarshaler = makeInterface("MarshalSia", ioPkg.Types.Scope().Lookup("Writer").Type())
	siaUnmarshaler = makeInterface("UnmarshalSia", ioPkg.Types.Scope().Lookup("Reader").Type())

	g := &generator{
		pkg:  srcPkg,
		typs: make(map[string]string),
	}
	g.addImport("io") // for io.Reader/io.Writer in method signatures
	g.addImport("gitlab.com/NebulousLabs/encoding")

	// check that all types are legal
	for _, typ := range typs {
		if err := g.checkType(typ); err != nil {
			return "", err
		}
	}

	// initialize g.typs with all of the types we're going to generate methods
	// for. This allows us to emit MarshalSia/UnmarshalSia calls for other types
	// in the same codegen batch.
	for _, typ := range typs {
		g.typs[typ] = "<PLACEHOLDER>"
	}

	// generate marshal/unmarshal methods for each specified type
	for _, typ := range typs {
		g.genMethods(g.pkg.Types.Scope().Lookup(typ).Type())
	}

	// output
	var methods []string
	for _, code := range g.typs {
		methods = append(methods, code)
	}
	raw := fmt.Sprintf(`// Code generated by encodegen. DO NOT EDIT.
package %s
import (%s)

%s
`, g.pkg.Name, strings.Join(g.imports, "\n"), strings.Join(methods, "\n"))

	// fmt.Printf("UNFORMATTED:\n%s\n", string(raw))

	formatted, err := format.Source([]byte(raw))
	if err != nil {
		panic(err) // should never happen
	}
	return string(formatted), nil
}

func (g *generator) checkType(typName string) error {
	var check func(t types.Type, ctx string) error
	check = func(t types.Type, ctx string) error {
		switch t := t.Underlying().(type) {
		case *types.Basic:
			if t.Info()&types.IsInteger != 0 || t.Kind() == types.Bool || t.Kind() == types.String {
				return nil
			}
		case *types.Array:
			return check(t.Elem(), ctx+"[0]")
		case *types.Slice:
			return check(t.Elem(), ctx+"[0]")
		case *types.Struct:
			for i := 0; i < t.NumFields(); i++ {
				field := t.Field(i)
				if err := check(field.Type(), ctx+"."+field.Name()); err != nil {
					return err
				}
			}
			return nil
		case *types.Pointer:
			return check(t.Elem(), "*"+ctx)
		}

		if ctx != "" {
			return fmt.Errorf("unsupported type %s at (%s)%s", t, typName, ctx)
		}
		return fmt.Errorf("unsupported type %s", t)
	}

	obj := g.pkg.Types.Scope().Lookup(typName)
	if obj == nil {
		return fmt.Errorf("no declaration found for type %q", typName)
	}
	return check(obj.Type(), "")
}

func (g *generator) genMethods(t types.Type) error {
	var enc, dec string
	switch t := t.Underlying().(type) {
	case *types.Basic:
		enc = g.genEncodeBody(t.Name()+"(x)", t)
		dec = g.genDecodeBody(t.Name()+"(x)", t)
	case *types.Array:
		enc = g.genEncodeBody("x", t)
		dec = g.genDecodeBody("x", t)
	case *types.Slice:
		enc = g.genEncodeBody("x", t)
		dec = g.genDecodeBody("x", t)
	case *types.Struct:
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			enc += g.genEncodeBody("x."+field.Name(), field.Type())
			dec += g.genDecodeBody("x."+field.Name(), field.Type())
		}
	default:
		// checkType should catch unhandled types, making this a developer error
		panic(fmt.Sprintf("unhandled type %T", t))
	}
	typName := types.TypeString(t, types.RelativeTo(g.pkg.Types))
	g.typs[typName] = fmt.Sprintf(`
// MarshalSia implements encoding.SiaMarshaler.
func (x %s) MarshalSia(w io.Writer) error {
	e := encoding.NewEncoder(w)
	%s
	return e.Err()
}

// UnmarshalSia implements encoding.SiaUnmarshaler.
func (x *%s) UnmarshalSia(r io.Reader) error {
	d := encoding.NewDecoder(r, encoding.DefaultAllocLimit)
	%s
	return d.Err()
}
`, typName, strings.TrimSpace(enc), typName, strings.TrimSpace(dec))
	return nil
}

func (g *generator) genEncodeBody(ident string, tOriginal types.Type) string {
	switch t := tOriginal.Underlying().(type) {
	case *types.Basic:
		if t.Info()&types.IsInteger != 0 {
			return fmt.Sprintf("e.WriteUint64(%s)\n", cast(ident, t, types.Typ[types.Uint64]))
		} else if t.Kind() == types.Bool {
			return fmt.Sprintf("e.WriteBool(%s)\n", cast(ident, t, types.Typ[types.Bool]))
		} else if t.Kind() == types.String {
			return fmt.Sprintf("e.WritePrefixedBytes([]byte(%s))\n", ident)
		}
	case *types.Array:
		// check for [...]byte
		if basic, ok := t.Elem().(*types.Basic); ok && basic.Kind() == types.Byte {
			return fmt.Sprintf("e.Write(%s[:])\n", ident)
		}
		// NOTE: it's fine to always use "v" as the loop variable, even in
		// nested loops; the inner v will shadow the outer v, but inner
		// loops never need to reference the variables of outer loops.
		body := g.genEncodeBody("v", t.Elem())
		return fmt.Sprintf("for _, v := range &%s { %s }\n", ident, body)
	case *types.Slice:
		// check for []byte
		if basic, ok := t.Elem().(*types.Basic); ok && basic.Kind() == types.Byte {
			return fmt.Sprintf("e.WritePrefixedBytes(%s)\n", ident)
		}
		prefix := fmt.Sprintf("e.WriteInt(len(%s))\n", ident)
		body := g.genEncodeBody("v", t.Elem())
		return prefix + fmt.Sprintf("for _, v := range %s { %s }\n", ident, body)
	case *types.Struct:
		// If the type has a MarshalSia method defined (or if they *will* have such
		// a method defined when we're done), use it.
		if types.Implements(tOriginal, siaMarshaler) || g.willGenerate(tOriginal) {
			return fmt.Sprintf("%s.MarshalSia(e)\n", ident)
		} else {
			var body string
			for i := 0; i < t.NumFields(); i++ {
				field := t.Field(i)
				body += g.genEncodeBody(ident+"."+field.Name(), field.Type())
			}
			return body
		}
	case *types.Pointer:
		body := g.genEncodeBody(fmt.Sprintf("(*%s)", ident), t.Elem())
		return fmt.Sprintf("e.WriteBool(%s != nil); if %s != nil { %s }\n", ident, ident, body)
	}
	panic("unreachable")
}

func (g *generator) genDecodeBody(ident string, tOriginal types.Type) string {
	// If the type is defined in a separate package, import it
	// We only need to do this for slices and pointers because they
	// are the only cases where we actually reference the type by name.
	// For decoding pointers, we make a call to new(typeName) so it is necessary,
	// and in slices we need to call make().
	// Adding imports for other kinds of types results in unused import errors
	if arr, ok := tOriginal.(*types.Slice); ok {
		g.addImportType(arr.Elem())
	} else if ptr, ok := tOriginal.(*types.Pointer); ok {
		g.addImportType(ptr.Elem())
	}

	switch t := tOriginal.Underlying().(type) {
	case *types.Basic:
		if t.Info()&types.IsInteger != 0 {
			return fmt.Sprintf("%s = %s\n", ident, cast("d.NextUint64()", types.Typ[types.Uint64], t))
		} else if t.Kind() == types.Bool {
			return fmt.Sprintf("%s = %s(d.NextBool())\n", ident, t)
		} else if t.Kind() == types.String {
			return fmt.Sprintf("%s = %s(d.ReadPrefixedBytes())\n", ident, t)
		}
	case *types.Array:
		// check for [...]byte
		if basic, ok := t.Elem().(*types.Basic); ok && basic.Kind() == types.Byte {
			return fmt.Sprintf("d.Read(%s[:])\n", ident)
		}
		// NOTE: we can use the same variable shadowing trick as genEncodeBody,
		// but we have to use a pointer, so things are slightly uglier.
		body := g.genDecodeBody("(*v)", t.Elem())
		return fmt.Sprintf("for i := range %s { v := &%s[i]; %s }\n", ident, ident, body)
	case *types.Slice:
		// check for []byte
		if basic, ok := t.Elem().(*types.Basic); ok && basic.Kind() == types.Byte {
			return fmt.Sprintf("%s = %s(d.ReadPrefixedBytes())\n", ident, t)
		}

		/*
			by default imported types will be expressed using their full paths
			by using (*types.Package).Name as the qualifier we condense strings like go.sia.tech/encodegen/internal/test_imported.Imported to test_imported.Imported
			read more: https://github.com/golang/example/blob/master/gotypes/go-types.md#formatting-support
		*/
		typeString := types.TypeString(t, (*types.Package).Name)

		prefix := fmt.Sprintf("%s = make(%s, d.NextPrefix(%d))\n", ident, typeString, sizeof(t.Elem()))
		body := g.genDecodeBody("(*v)", t.Elem())
		return prefix + fmt.Sprintf("for i := range %s {v := &%s[i]; %s}\n", ident, ident, body)
	case *types.Struct:
		// If the type has an UnmarshalSia method defined (or if they *will* have such
		// a method defined when we're done), use it.
		if types.Implements(types.NewPointer(tOriginal), siaUnmarshaler) || g.willGenerate(tOriginal) {
			return fmt.Sprintf("%s.UnmarshalSia(d)\n", ident)
		} else {
			var body string
			for i := 0; i < t.NumFields(); i++ {
				field := t.Field(i)
				body += g.genDecodeBody(ident+"."+field.Name(), field.Type())
			}
			return body
		}
	case *types.Pointer:
		body := g.genDecodeBody(fmt.Sprintf("(*%s)", ident), t.Elem())
		return fmt.Sprintf("if d.NextBool() { %s = new(%s); %s }\n", ident, types.TypeString(t.Elem(), (*types.Package).Name), body)
	}
	panic("unreachable")
}

func sizeof(t types.Type) int {
	switch t := t.Underlying().(type) {
	case *types.Basic:
		if t.Info()&types.IsInteger != 0 {
			return 8
		} else if t.Kind() == types.Bool {
			return 1
		} else if t.Kind() == types.String {
			return 8
		}
	case *types.Array:
		return int(t.Len()) * sizeof(t.Elem())
	case *types.Slice:
		return 8
	case *types.Struct:
		var total int
		for i := 0; i < t.NumFields(); i++ {
			total += sizeof(t.Field(i).Type())
		}
		return total
	case *types.Pointer:
		// sizeof the true/false for whether its null or not
		return 1
	}
	panic("unreachable")
}

func cast(ident string, from types.Type, to types.Type) string {
	/*
		I *think* types.AssignableTo might be preferable here, but we should
		check that it doesn't skip any casts that are actually necessary
	*/
	if types.Identical(from, to) {
		return ident
	}
	return fmt.Sprintf("%s(%s)", types.TypeString(to, (*types.Package).Name), ident)
}

func gopath() string {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return gopath
	} else {
		return build.Default.GOPATH
	}
}
