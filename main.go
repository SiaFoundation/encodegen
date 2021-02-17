package main

import (
	// "gitlab.com/NebulousLabs/encoding"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	src := `
package main

type TestType1 struct {
	// A ***uint64
	// B int64
	// C string
	// E int
	// F byte
	// G *byte
	H []*byte
	I []TestType2
}
type TestType2 struct {
	I string
	// J *TestType3
	K byte
	X [][]TestType3
}

type TestType3 struct {
	L uint64
	// Y *TestType4
}

type TestType4 struct {
	K string
}

`
	// Create the AST by parsing src.
	fs := token.NewFileSet() // positions are relative to fset
	parsed, err := parser.ParseFile(fs, "test.go", src, 0)
	if err != nil {
		log.Fatal(err)
	}
	// ast.Print(fs, parsed)

	structs := make(map[string]*ast.StructType)

	// gather all types
	for _, decl := range parsed.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			structs[typeSpec.Name.Name] = structType
		}
	}

	var g generator

	for structName, structType := range structs {
		g.Printf("func (r *%s) unmarshalBuffer(b *objBuffer) {\n", structName)
		writeFieldAssignmentStatements(&g, structs, []string{}, structType)
		g.Printf("}\n")
	}

	// log.Printf("UNFORMATTED:\n%s", string(g.buf.Bytes()))
	log.Printf("FORMATTED:\n%s", string(g.format()))
}
