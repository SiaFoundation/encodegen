package main

import (
	"go/ast"
	"log"
)

const ReadBoolFunction = "readBool"
const ReadIntFunction = "readUint64"
const ReadStringFunction = "readPrefixedBytes"
const ReadByteFunction = "readByte"

const WriteBoolFunction = "writeBool"
const WriteIntFunction = "writeUint64"
const WriteStringFunction = "writePrefixedBytes"
const WriteByteFunction = "writeByte"

type PrimitiveFunctions struct {
	ReadFunction  string
	WriteFunction string
	WriteCast     string
}

type FieldType struct {
	Name               string
	PrimitiveFunctions PrimitiveFunctions
	Primitive          bool
	StarCount          int
	ArrayCount         int
}

var IntPrimitiveFunctions = PrimitiveFunctions{
	ReadFunction:  ReadIntFunction,
	WriteFunction: WriteIntFunction,
	WriteCast:     "uint64",
}

var supportedPrimitives = map[string]PrimitiveFunctions{
	"bool":   {ReadFunction: ReadBoolFunction, WriteFunction: WriteBoolFunction, WriteCast: "bool"},
	"string": {ReadFunction: ReadStringFunction, WriteFunction: WriteStringFunction, WriteCast: "[]byte"},
	"byte":   {ReadFunction: ReadByteFunction, WriteFunction: WriteByteFunction, WriteCast: "byte"},
	"int":    IntPrimitiveFunctions,
	"int8":   IntPrimitiveFunctions,
	"int16":  IntPrimitiveFunctions,
	"int32":  IntPrimitiveFunctions,
	"int64":  IntPrimitiveFunctions,
	"uint":   IntPrimitiveFunctions,
	"uint8":  IntPrimitiveFunctions,
	"uint16": IntPrimitiveFunctions,
	"uint32": IntPrimitiveFunctions,
	"uint64": IntPrimitiveFunctions,
}

func isPrimitive(ty *ast.Ident) bool {
	return isPrimitiveString(ty.Name)
}

func isPrimitiveString(t string) bool {
	_, ok := supportedPrimitives[t]
	return ok
}

func getFieldType(exp ast.Expr, existingStarCount int, existingArrayCount int) *FieldType {
	switch v := exp.(type) {
	case *ast.Ident:
		if isPrimitive(v) {
			return &FieldType{
				Name:               v.Name,
				PrimitiveFunctions: supportedPrimitives[v.Name],
				Primitive:          true,
				StarCount:          existingStarCount,
				ArrayCount:         existingArrayCount,
			}
		} else {
			return &FieldType{
				Name:       v.Name,
				Primitive:  false,
				StarCount:  existingStarCount,
				ArrayCount: existingArrayCount,
			}
		}
		break
	case *ast.ArrayType:
		// count is for multidimensional arrays, currently not implemented however
		existingArrayCount++
		fieldType := getFieldType(v.Elt, existingStarCount, existingArrayCount)
		return fieldType
	case *ast.StarExpr:
		// the star count allows us to keep track of how many levels of dereferencing there is. i.e. ***int should have a starcount of 3
		existingStarCount++
		fieldType := getFieldType(v.X, existingStarCount, existingArrayCount)
		return fieldType
		break
	default:
		log.Printf("WARNING: This struct contains unsupported data (%+v)", exp)
		return nil
		break
	}
	return nil
}
