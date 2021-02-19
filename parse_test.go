package main

import (
	"go/ast"
	"reflect"
	"testing"
)

func TestIsPrimitiveString(t *testing.T) {
	if !isPrimitiveString("uint8") {
		t.Errorf("isPrimitiveString failed to recognize uint8 as supported")
	}
	if isPrimitiveString("wasdf99") {
		t.Errorf("isPrimitiveString recognized invalid type as supported")
	}
}

func TestGetFieldType(t *testing.T) {
	// struct (non primitive type)
	ident := &ast.Ident{Name: "UnknownStruct"}
	expectedOutput := &FieldType{
		Name:         "UnknownStruct",
		FunctionName: "",
		Primitive:    false,
		StarCount:    0,
		ArrayCount:   0,
	}

	fieldType := getFieldType(ident, 0, 0)
	if !reflect.DeepEqual(expectedOutput, fieldType) {
		t.Errorf("Output (%+v) did not equal expected output (%+v)", fieldType, expectedOutput)
	}

	// uint8 (primitive type)
	ident = &ast.Ident{Name: "uint8"}
	expectedOutput = &FieldType{
		Name:         "uint8",
		FunctionName: supportedPrimitives["uint8"],
		Primitive:    true,
		StarCount:    0,
		ArrayCount:   0,
	}

	fieldType = getFieldType(ident, 0, 0)
	if !reflect.DeepEqual(expectedOutput, fieldType) {
		t.Errorf("Output (%+v) did not equal expected output (%+v)", fieldType, expectedOutput)
	}

	// one dimensional array of uint8s
	arrayIdent := &ast.ArrayType{Elt: &ast.Ident{
		Name: "uint8",
	}}
	expectedOutput = &FieldType{
		Name:         "uint8",
		FunctionName: supportedPrimitives["uint8"],
		Primitive:    true,
		StarCount:    0,
		ArrayCount:   1,
	}

	fieldType = getFieldType(arrayIdent, 0, 0)
	if !reflect.DeepEqual(expectedOutput, fieldType) {
		t.Errorf("Output (%+v) did not equal expected output (%+v)", fieldType, expectedOutput)
	}

	// one dimensional array of UnknownStructs
	arrayIdent = &ast.ArrayType{Elt: &ast.Ident{
		Name: "UnknownStruct",
	}}
	expectedOutput = &FieldType{
		Name:         "UnknownStruct",
		FunctionName: "",
		Primitive:    false,
		StarCount:    0,
		ArrayCount:   1,
	}

	fieldType = getFieldType(arrayIdent, 0, 0)
	if !reflect.DeepEqual(expectedOutput, fieldType) {
		t.Errorf("Output (%+v) did not equal expected output (%+v)", fieldType, expectedOutput)
	}

	// pointer
	pointerIdent := &ast.StarExpr{X: &ast.Ident{
		Name: "UnknownStruct",
	}}

	expectedOutput = &FieldType{
		Name:         "UnknownStruct",
		FunctionName: "",
		Primitive:    false,
		StarCount:    1,
		ArrayCount:   0,
	}

	fieldType = getFieldType(pointerIdent, 0, 0)
	if !reflect.DeepEqual(expectedOutput, fieldType) {
		t.Errorf("Output (%+v) did not equal expected output (%+v)", fieldType, expectedOutput)
	}

	// double pointer
	pointerIdent = &ast.StarExpr{X: pointerIdent}

	expectedOutput = &FieldType{
		Name:         "UnknownStruct",
		FunctionName: "",
		Primitive:    false,
		StarCount:    2,
		ArrayCount:   0,
	}

	fieldType = getFieldType(pointerIdent, 0, 0)
	if !reflect.DeepEqual(expectedOutput, fieldType) {
		t.Errorf("Output (%+v) did not equal expected output (%+v)", fieldType, expectedOutput)
	}

	// array of double pointers
	arrayIdent = &ast.ArrayType{Elt: pointerIdent}

	expectedOutput = &FieldType{
		Name:         "UnknownStruct",
		FunctionName: "",
		Primitive:    false,
		StarCount:    2,
		ArrayCount:   1,
	}

	fieldType = getFieldType(arrayIdent, 0, 0)
	if !reflect.DeepEqual(expectedOutput, fieldType) {
		t.Errorf("Output (%+v) did not equal expected output (%+v)", fieldType, expectedOutput)
	}

}
