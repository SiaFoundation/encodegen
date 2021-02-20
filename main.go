package main

import (
	// "gitlab.com/NebulousLabs/encoding"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify file(s) to generate code from.")
	}

	for _, path := range os.Args[1:] {
		// I know Go devs are trying to scrap ioutil but this was only done in 1.16 which came out just a few days ago which means os.ReadFile for the majority of Go users is not defined so they must use ioutil.ReadFile
		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal("Reading file: ", err)
		}

		// Create the AST by parsing src.
		fs := token.NewFileSet() // positions are relative to fset
		parsed, err := parser.ParseFile(fs, path, fileContent, 0)
		if err != nil {
			log.Fatal("Parsing file: ", path, ": ", err)
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

		// write the output to the generator
		for structName, structType := range structs {
			g.Printf("func (r *%s) unmarshalBuffer(b *objBuffer) {\n", structName)
			g.writeUnmarshalBufferStatements(structs, []string{"r"}, "i", structType)
			g.Printf("}\n")

			g.Printf("func (r *%s) marshalBuffer(b *objBuffer) {\n", structName)
			g.writeMarshalBufferStatements(structs, []string{"r"}, "i", structType)
			g.Printf("}\n")
		}

		// log.Printf("UNFORMATTED:\n%s", string(g.buf.Bytes()))
		log.Printf("FORMATTED %s:\n%s", path, string(g.format()))
	}
}
