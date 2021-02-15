package main

import (
	// "gitlab.com/NebulousLabs/encoding"
	"fmt"
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
}

var a = 5;
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
		var statements []jen.Code
		unmarshalFunc := jen.Func().Params(jen.Id("r").Op("*").Id(structName)).Id("unmarshalBuffer").Params(jen.Id("b").Op("*").Id("objBuffer"))
		for _, structField := range structType.Fields.List {
			log.Printf("%+v", *structField)
			if len(structField.Names) > 0 {
				fieldNameString, aa := getFieldType(structField.Type, map[string]string{})
				if err != nil {
					continue
				}
				log.Print(aa)

				fieldAssignment := jen.Id("r").Dot(structField.Names[0].Name).Op("=").Id("b").Dot("read" + strings.Title(fieldNameString)).Call()
				statements = append(statements, fieldAssignment)
			}
		}
		codeOutput.Add(unmarshalFunc.Block(statements...))
	}
	log.Printf("%#v", codeOutput)

}

//Returns a string representation of the given expression if it was recognized.
//Refer to the implementation to see the different string representations.
func getFieldType(exp ast.Expr, aliases map[string]string) (string, []string) {
	switch v := exp.(type) {
	case *ast.Ident:
		return getIdent(v, aliases)
	case *ast.ArrayType:
		return getArrayType(v, aliases)
	case *ast.StarExpr:
		return getStarExp(v, aliases)
	}
	return "", []string{}
}

func getIdent(v *ast.Ident, aliases map[string]string) (string, []string) {

	if isPrimitive(v) {
		return v.Name, []string{}
	}
	t := fmt.Sprintf("%s", v.Name)
	return t, []string{t}
}

func getArrayType(v *ast.ArrayType, aliases map[string]string) (string, []string) {
	t, fundamentalTypes := getFieldType(v.Elt, aliases)
	return fmt.Sprintf("%ss", t), fundamentalTypes
}

func getInterfaceType(v *ast.InterfaceType, aliases map[string]string) (string, []string) {

	methods := make([]string, 0)
	for _, field := range v.Methods.List {
		methodName := ""
		if field.Names != nil {
			methodName = field.Names[0].Name
		}
		t, _ := getFieldType(field.Type, aliases)
		methods = append(methods, methodName+" "+t)
	}
	return fmt.Sprintf("{%s}", strings.Join(methods, "; ")), []string{}
}

func getStarExp(v *ast.StarExpr, aliases map[string]string) (string, []string) {
	t, f := getFieldType(v.X, aliases)
	// return fmt.Sprintf("*%s", t), f
	return fmt.Sprintf("%s", t), f
}

var globalPrimitives = map[string]struct{}{
	"bool":        {},
	"string":      {},
	"int":         {},
	"int8":        {},
	"int16":       {},
	"int32":       {},
	"int64":       {},
	"uint":        {},
	"uint8":       {},
	"uint16":      {},
	"uint32":      {},
	"uint64":      {},
	"uintptr":     {},
	"byte":        {},
	"rune":        {},
	"float32":     {},
	"float64":     {},
	"complex64":   {},
	"complex128":  {},
	"error":       {},
	"*bool":       {},
	"*string":     {},
	"*int":        {},
	"*int8":       {},
	"*int16":      {},
	"*int32":      {},
	"*int64":      {},
	"*uint":       {},
	"*uint8":      {},
	"*uint16":     {},
	"*uint32":     {},
	"*uint64":     {},
	"*uintptr":    {},
	"*byte":       {},
	"*rune":       {},
	"*float32":    {},
	"*float64":    {},
	"*complex64":  {},
	"*complex128": {},
	"*error":      {},
}

func isPrimitive(ty *ast.Ident) bool {
	return isPrimitiveString(ty.Name)
}

func isPrimitiveString(t string) bool {
	_, ok := globalPrimitives[t]
	return ok
}

const packageConstant = "{packageName}"

func replacePackageConstant(field, packageName string) string {
	if packageName != "" {
		packageName = fmt.Sprintf("%s.", packageName)
	}
	return strings.Replace(field, packageConstant, packageName, 1)
}
