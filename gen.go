package main

import (
	"fmt"
	"go.sia.tech/encodegen/internal/test_imported"
	"go/format"
	"go/types"
	"golang.org/x/tools/go/packages"
	"strings"
)

type generator struct {
	pkg     *packages.Package
	typs    map[string]string
	imports []string
}

func (g *generator) addImport(pkg string) {
	for _, p := range g.imports {
		if p == pkg {
			return
		}
	}
	g.imports = append(g.imports, `"`+pkg+`"`)
}

func Generate(pkgName string, typs ...string) (string, error) {
	// load source
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, pkgName)
	if err != nil {
		return "", err
	}

	g := &generator{
		pkg:  pkgs[0],
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

	// generate marshal/unmarshal methods for each specified type
	for _, typ := range typs {
		g.genMethods(typ)
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
`, g.pkg.Name, strings.Join(g.imports, ";"), strings.Join(methods, "\n"))

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

func (g *generator) genMethods(typName string) error {
	if _, ok := g.typs[typName]; ok {
		return nil // already generated
	}
	var enc, dec string
	switch t := g.pkg.Types.Scope().Lookup(typName).Type().Underlying().(type) {
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

func (g *generator) genEncodeBody(ident string, t types.Type) string {
	switch t := t.Underlying().(type) {
	case *types.Basic:
		if t.Info()&types.IsInteger != 0 {
			return fmt.Sprintf("e.WriteUint64(uint64(%s))\n", ident)
		} else if t.Kind() == types.Bool {
			return fmt.Sprintf("e.WriteBool(bool(%s))\n", ident)
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
			return fmt.Sprintf("e.WritePrefixedBytes([]byte(%s))\n", ident)
		}
		prefix := fmt.Sprintf("e.WriteInt(len(%s))\n", ident)
		body := g.genEncodeBody("v", t.Elem())
		return prefix + fmt.Sprintf("for _, v := range %s { %s }\n", ident, body)
	case *types.Struct:
		var body string
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			body += g.genEncodeBody(ident+"."+field.Name(), field.Type())
		}
		return body
	}
	panic("unreachable")
}

func (g *generator) genDecodeBody(ident string, t types.Type) string {
	switch t := t.Underlying().(type) {
	case *types.Basic:
		if t.Info()&types.IsInteger != 0 {
			return fmt.Sprintf("%s = %s(d.NextUint64())\n", ident, t)
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

		// if we have an imported type
		if named, ok := t.Elem().(*types.Named); ok {
			g.addImport(named.Obj().Pkg().Path())
		}

		prefix := fmt.Sprintf("%s = make(%s, d.NextPrefix(%d))\n", ident, typeString, sizeof(t.Elem()))
		body := g.genDecodeBody("(*v)", t.Elem())
		return prefix + fmt.Sprintf("for i := range %s {v := &%s[i]; %s}\n", ident, ident, body)
	case *types.Struct:
		var body string
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			body += g.genDecodeBody(ident+"."+field.Name(), field.Type())
		}
		return body
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
	}
	panic("unreachable")
}

// for testing; delete later

type Foo struct {
	A test_imported.Imported
	B test_imported.Hash
	C []test_imported.Imported
	D [][]test_imported.Imported
	E [40]test_imported.Imported
	X int
	Y uint64
	Z [][][]uint64
}

type FooAlias Foo

type Slice struct {
	b []byte
}

type Hash [32]byte

type Array struct {
	Bar int
	Str string
	LOL [10][]struct {
		Inner [3]uint32
	}
}
