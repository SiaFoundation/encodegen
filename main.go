package main

import (
	// "gitlab.com/NebulousLabs/encoding"
	"github.com/dave/jennifer/jen"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func main() {
	src := `
package main

type TestType1 struct {
	A uint64
	B *int64
	C uint32
	D int32
	E bool
	F string
	G []byte
	H TestType2
}
type TestType2 struct {
	I string
	J TestType3
}

type TestType3 struct {
	K uint64
	L TestType4
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

	codeOutput := jen.NewFile("main")

	for structName, structType := range structs {
		unmarshalFunc := jen.Func().Params(jen.Id("r").Op("*").Id(structName)).Id("unmarshalBuffer").Params(jen.Id("b").Op("*").Id("objBuffer"))
		statements := getFieldAssignmentStatements(structs, []string{}, structName, structType)

		codeOutput.Add(unmarshalFunc.Block(statements...))
	}
	log.Printf("%#v", codeOutput)

}

type FieldType struct {
	Name         string
	FunctionName string
	Primitive    bool
	Pointer      bool
}

func getFieldAssignmentStatements(structs map[string]*ast.StructType, prefixs []string, currentStructName string, currentStructType *ast.StructType) []jen.Code {
	var statements []jen.Code
	var fieldAssignment *jen.Statement
	for _, structField := range currentStructType.Fields.List {
		if len(structField.Names) > 0 {
			field := getFieldType(structField.Type)
			if field == nil {
				continue
			}

			if len(prefixs) == 0 {
				fieldAssignment = jen.Id("r").Dot(structField.Names[0].Name).Op("=")
			} else {
				fieldAssignment = jen.Id("r")
				for _, prefix := range prefixs {
					fieldAssignment = fieldAssignment.Dot(prefix)
				}
				fieldAssignment = fieldAssignment.Dot(structField.Names[0].Name).Op("=")
			}

			if field.Primitive {
				fieldAssignment = fieldAssignment.Op(field.Name).Parens(jen.Id("b").Dot(field.FunctionName).Call())
				statements = append(statements, fieldAssignment)
			} else {
				prefixs = append(prefixs, structField.Names[0].Name)
				statements = append(statements, getFieldAssignmentStatements(structs, prefixs, field.Name, structs[field.Name])...)
			}

		}
	}
	return statements
}

//Returns a string representation of the given expression if it was recognized.
//Refer to the implementation to see the different string representations.
func getFieldType(exp ast.Expr) *FieldType {
	switch v := exp.(type) {
	case *ast.Ident:
		if isPrimitive(v) {
			return &FieldType{
				Name:         v.Name,
				FunctionName: supportedPrimitives[v.Name],
				Primitive:    true,
				Pointer:      false,
			}
			log.Print("primitive")
		} else {
			return &FieldType{
				Name:         v.Name,
				FunctionName: "read" + strings.Title(v.Name),
				Primitive:    false,
				Pointer:      false,
			}
		}
		break
		// case *ast.ArrayType:
		// 	return getArrayType(v)
		// case *ast.StarExpr:
		// 	return getStarExp(v)
	default:
		log.Print("WARNING: This struct contains unsupported types")
	}
	return nil
}

const ReadBoolFunction = "readBool"
const ReadIntFunction = "readUint64"
const ReadStringFunction = "readPrefixedBytes"
const ReadByteFunction = "readByte"

const ReadBoolPointerFunction = "readBool"
const ReadIntPointerFunction = "readUint64"
const ReadStringPointerFunction = "readPrefixedBytes"
const ReadBytePointerFunction = "readByte"

// TEMPORARY
const UnsupportedTypeFunction = "panic"

var supportedPrimitives = map[string]string{
	"bool":   ReadBoolFunction,
	"string": ReadStringFunction,
	"int":    ReadIntFunction,
	"int8":   ReadIntFunction,
	"int16":  ReadIntFunction,
	"int32":  ReadIntFunction,
	"int64":  ReadIntFunction,
	"uint":   ReadIntFunction,
	"uint8":  ReadIntFunction,
	"uint16": ReadIntFunction,
	"uint32": ReadIntFunction,
	"uint64": ReadIntFunction,
	"byte":   ReadByteFunction,
}

func isPrimitive(ty *ast.Ident) bool {
	return isPrimitiveString(ty.Name)
}

func isPrimitiveString(t string) bool {
	_, ok := supportedPrimitives[t]
	return ok
}
