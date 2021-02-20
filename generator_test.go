package main

import (
	"go/ast"
	"testing"
)

func TestNextRangeIdentifier(t *testing.T) {
	output := nextRangeIdentifier("i")
	expectedOutput := "i1"
	if output != expectedOutput {
		t.Errorf("Identifier output (%+v) did not equal expected output (%+v)", output, expectedOutput)
	}

	output = nextRangeIdentifier("i1")
	expectedOutput = "i2"
	if output != expectedOutput {
		t.Errorf("Identifier output (%+v) did not equal expected output (%+v)", output, expectedOutput)
	}

	output = nextRangeIdentifier("i99")
	expectedOutput = "i100"
	if output != expectedOutput {
		t.Errorf("Identifier output (%+v) did not equal expected output (%+v)", output, expectedOutput)
	}

}

func TestGenerator(t *testing.T) {
	var g generator
	const helloWorldCode = `
import "log"
func main() {
	log.Print("hello, world")
}
`
	g.Printf("%s", helloWorldCode)
	output := g.buf.Bytes()
	if string(output) != helloWorldCode {
		t.Errorf("Generated output (%+v) did not equal expected output (%+v).  Were the bytes properly written to the generator buffer?", output, helloWorldCode)
	}
}

// the writeUnmarshalBufferStatements code is being changed a lot so this test is not particularly meaningful as of now
// func TestWriteUnmarshalBufferStatements(t *testing.T) {

// 	// single primitive variable struct
// 	structs := make(map[string]*ast.StructType)
// 	structs["TestType1"] = &ast.StructType{Fields: &ast.FieldList{
// 		List: []*ast.Field{
// 			{Names: []*ast.Ident{{Name: "A"}}, Type: &ast.Ident{Name: "uint64"}},
// 		},
// 	}}

// 	var g generator

// 	for structName, structType := range structs {
// 		g.Printf("func (r *%s) unmarshalBuffer(b *objBuffer) {\n", structName)
// 		g.writeUnmarshalBufferStatements(structs, []string{"r"}, "i", structType)
// 		g.Printf("}\n")
// 	}
// 	expectedOutput := `func (r *TestType1) unmarshalBuffer(b *objBuffer) {
// 	r.A = uint64(b.readUint64())
// }
// `

// 	output := string(g.format())

// 	if expectedOutput != output {
// 		t.Errorf("Generated output (%+v) did not equal expected output (%+v).", output, expectedOutput)
// 	}

// 	// reset the buffer
// 	g.reset()

// 	// two primitive variable struct
// 	structs = make(map[string]*ast.StructType)
// 	structs["TestType1"] = &ast.StructType{Fields: &ast.FieldList{
// 		List: []*ast.Field{
// 			{Names: []*ast.Ident{{Name: "A"}}, Type: &ast.Ident{Name: "int8"}},
// 			{Names: []*ast.Ident{{Name: "B"}}, Type: &ast.Ident{Name: "string"}},
// 		},
// 	}}

// 	for structName, structType := range structs {
// 		g.Printf("func (r *%s) unmarshalBuffer(b *objBuffer) {\n", structName)
// 		g.writeUnmarshalBufferStatements(structs, []string{"r"}, "i", structType)
// 		g.Printf("}\n")
// 	}
// 	expectedOutput = `func (r *TestType1) unmarshalBuffer(b *objBuffer) {
// 	r.A = int8(b.readUint64())
// 	r.B = string(b.readPrefixedBytes())
// }
// `
// 	output = string(g.format())

// 	if expectedOutput != output {
// 		t.Errorf("Generated output (%+v) did not equal expected output (%+v).", output, expectedOutput)
// 	}

// }
