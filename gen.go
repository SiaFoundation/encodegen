package main

import (
	"fmt"
	"go/format"
	"go/token"
	"go/types"
	"path"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	typesPackage       = "go.sia.tech/core/types"
	currencyDefinition = typesPackage + ".Currency"
)

var encoderTo, decoderFrom *types.Interface

func makeInterface(name string, param types.Type) *types.Interface {
	params := types.NewTuple(types.NewVar(token.NoPos, nil, "", param))
	results := types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type()))
	meth := types.NewFunc(token.NoPos, nil, name, types.NewSignature(nil, params, results, false))
	return types.NewInterfaceType([]*types.Func{meth}, nil).Complete()
}

type gentype struct {
	code string
}

type generator struct {
	pkg     *packages.Package
	typs    map[string]gentype
	imports map[string]string
}

func (g *generator) importQualifier(pkg *types.Package) string {
	name := pkg.Name()
	for g.imports[name] != "" && g.imports[name] != pkg.Path() {
		name += "_"
	}
	return name
}

func (g *generator) typeString(t types.Type) string {
	return types.TypeString(t, func(other *types.Package) string {
		if g.pkg.Types == other {
			return "" // same package; unqualified
		}

		// external package; add import and qualify with package name
		qual := g.importQualifier(other)
		g.imports[qual] = other.Path()
		return qual
	})
}

func (g *generator) cast(ident string, from types.Type, to types.Type) string {
	// I *think* types.AssignableTo might be preferable here, but we should
	// check that it doesn't skip any casts that are actually necessary
	if types.Identical(from, to) {
		return ident
	}
	return fmt.Sprintf("%s(%s)", g.typeString(to), ident)
}

func (g *generator) willGenerate(t types.Type) bool {
	// TODO: could this just be return g.typs[g.typeString(t)] != ""
	if named, ok := t.(*types.Named); ok && named.Obj().Pkg() == g.pkg.Types {
		return g.typs[g.typeString(t)].code != ""
	}
	return false
}

func Generate(dir string, typs ...string) ([]byte, error) {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports}

	// load source package; also load "go.sia.tech/core/types", to construct interface types
	pkgs, err := packages.Load(cfg, dir, typesPackage)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			return nil, err
		}
	}
	// ordering of pkgs is unspecified, so locate packages manually
	var srcPkg, typesPkg *packages.Package
	if len(pkgs) == 2 {
		srcPkg, typesPkg = pkgs[0], pkgs[1]
		if typesPkg.PkgPath != typesPackage {
			srcPkg, typesPkg = typesPkg, srcPkg
		}
	} else if len(pkgs) == 1 {
		// we are in go.sia.tech/core/types
		srcPkg, typesPkg = pkgs[0], pkgs[0]
	}

	// construct interface types
	encoderTo = typesPkg.Types.Scope().Lookup("EncoderTo").Type().Underlying().(*types.Interface)
	decoderFrom = typesPkg.Types.Scope().Lookup("DecoderFrom").Type().Underlying().(*types.Interface)

	g := &generator{
		pkg:  srcPkg,
		typs: make(map[string]gentype),
		imports: map[string]string{
			"types": typesPackage,
		},
	}

	// initialize g.typs with all of the types we're going to generate methods
	// for. This allows us to emit EncodeTo/DecodeFrom calls for other types
	// in the same codegen batch.
	for _, typ := range typs {
		// get unmarshaler allocation limit expression (if present)
		typSplit := strings.Split(typ, ":")
		g.typs[typSplit[0]] = gentype{
			code: "<PLACEHOLDER>",
		}

	}
	// check that all types are legal
	for typ := range g.typs {
		if err := g.checkType(typ); err != nil {
			return nil, fmt.Errorf("cannot generate methods for type %v: %w", typ, err)
		}
	}

	// generate marshal/unmarshal methods for each specified type
	for typ := range g.typs {
		g.genMethods(g.pkg.Types.Scope().Lookup(typ).Type())
	}

	// output
	var methods []string
	for _, typ := range g.typs {
		methods = append(methods, typ.code)
	}

	var importString string
	for qual, fullpath := range g.imports {
		// omit qualifier if possible
		if qual == path.Base(fullpath) {
			qual = ""
		}
		importString += fmt.Sprintf("%s %q\n", qual, fullpath)
	}

	raw := fmt.Sprintf(`// Code generated by encodegen. DO NOT EDIT.
package %s
import (%s)

%s
`, g.pkg.Name, importString, strings.Join(methods, "\n"))

	formatted, err := format.Source([]byte(raw))
	if err != nil {
		panic(err) // should never happen
	}
	return formatted, nil
}

func (g *generator) checkType(typName string) error {
	var check func(t types.Type, ctx string) error
	check = func(t types.Type, ctx string) error {
		// If the type already implements both the marshal and unmarshal interface
		// we can skip checking since the generater will use them.
		if (types.Implements(t, encoderTo) || types.Implements(types.NewPointer(t), encoderTo)) && types.Implements(types.NewPointer(t), decoderFrom) {
			return nil
		}

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
				if !field.Exported() {
					if ctx != "" {
						return fmt.Errorf("unexported field %s at (%s)%s", field.Name(), typName, ctx)
					}
					return fmt.Errorf("unexported field %s", field.Name())
				}
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
	typName := g.typeString(t)
	g.typs[typName] = gentype{
		code: fmt.Sprintf(`
// EncodeTo implements types.EncoderTo.
func (x %s) EncodeTo(e *types.Encoder) {
	%s
}

// DecodeFrom implements types.DecoderFrom.
func (x *%s) DecodeFrom(d *types.Decoder) {
	%s
}
`, typName, strings.TrimSpace(enc), typName, strings.TrimSpace(dec)),
	}
	return nil
}

func (g *generator) genEncodeBody(ident string, t types.Type) string {
	// If the type has a EncodeTo method defined (or if they *will* have such
	// a method defined when we're done), use it.
	if types.Implements(t, encoderTo) || types.Implements(types.NewPointer(t), encoderTo) || g.willGenerate(t) {
		// If t is a pointer type, don't duplicate the nil-check here; instead,
		// fallthrough to the logic below, which will end up calling EncodeTo
		// on t.Elem().
		if _, isPointer := t.Underlying().(*types.Pointer); !isPointer {
			return fmt.Sprintf("%s.EncodeTo(e)\n", ident)
		}
	} else if t.String() == currencyDefinition {
		return fmt.Sprintf("types.V1Currency(%s).EncodeTo(e)\n", ident)
	}

	switch t := t.Underlying().(type) {
	case *types.Basic:
		if t.Info()&types.IsInteger != 0 {
			return fmt.Sprintf("e.WriteUint64(%s)\n", g.cast(ident, t, types.Typ[types.Uint64]))
		} else if t.Kind() == types.Bool {
			return fmt.Sprintf("e.WriteBool(%s)\n", g.cast(ident, t, types.Typ[types.Bool]))
		} else if t.Kind() == types.String {
			return fmt.Sprintf("e.WriteBytes([]byte(%s))\n", ident)
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
			return fmt.Sprintf("e.WriteBytes(%s)\n", ident)
		}
		prefix := fmt.Sprintf("e.WritePrefix(len(%s))\n", ident)
		body := g.genEncodeBody("v", t.Elem())
		return prefix + fmt.Sprintf("for _, v := range %s { %s }\n", ident, body)
	case *types.Struct:
		var body string
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			body += g.genEncodeBody(ident+"."+field.Name(), field.Type())
		}
		return body
	case *types.Pointer:
		body := g.genEncodeBody(fmt.Sprintf("(*%s)", ident), t.Elem())
		return fmt.Sprintf("e.WriteBool(%s != nil); if %s != nil { %s }\n", ident, ident, body)
	}
	panic("unreachable")
}

func (g *generator) genDecodeBody(ident string, t types.Type) string {
	// If the type has an DecodeFrom method defined (or if they *will* have such
	// a method defined when we're done), use it.
	if types.Implements(types.NewPointer(t), decoderFrom) || g.willGenerate(t) {
		// If t is a pointer type, don't duplicate the nil-check here; instead,
		// fallthrough to the logic below, which will end up calling
		// DecodeFrom on t.Elem().
		if _, isPointer := t.Underlying().(*types.Pointer); !isPointer {
			return fmt.Sprintf("%s.DecodeFrom(d)\n", ident)
		}
	} else if t.String() == currencyDefinition {
		return fmt.Sprintf("(*types.V1Currency)(&%s).DecodeFrom(d)\n", ident)
	}

	switch t := t.Underlying().(type) {
	case *types.Basic:
		if t.Info()&types.IsInteger != 0 {
			return fmt.Sprintf("%s = %s\n", ident, g.cast("d.ReadUint64()", types.Typ[types.Uint64], t))
		} else if t.Kind() == types.Bool {
			return fmt.Sprintf("%s = %s(d.ReadBool())\n", ident, t)
		} else if t.Kind() == types.String {
			return fmt.Sprintf("%s = %s(d.ReadBytes())\n", ident, t)
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
			return fmt.Sprintf("%s = %s(d.ReadBytes())\n", ident, t)
		}
		prefix := fmt.Sprintf("%s = make(%s, d.ReadPrefix())\n", ident, g.typeString(t))
		body := g.genDecodeBody("(*v)", t.Elem())
		return prefix + fmt.Sprintf("for i := range %s {v := &%s[i]; %s}\n", ident, ident, body)
	case *types.Struct:
		var body string
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			body += g.genDecodeBody(ident+"."+field.Name(), field.Type())
		}
		return body
	case *types.Pointer:
		body := g.genDecodeBody(fmt.Sprintf("(*%s)", ident), t.Elem())
		return fmt.Sprintf("if d.ReadBool() { %s = new(%s); %s }\n", ident, g.typeString(t.Elem()), body)
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
