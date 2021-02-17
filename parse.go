package main

import (
	"go/ast"
	"log"
	"strings"
)

const ReadBoolFunction = "readBool"
const ReadIntFunction = "readUint64"
const ReadStringFunction = "readPrefixedBytes"
const ReadByteFunction = "readByte"

const ReadBoolPointerFunction = "readBool"
const ReadIntPointerFunction = "readUint64"
const ReadStringPointerFunction = "readPrefixedBytes"
const ReadBytePointerFunction = "readByte"

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

func getFieldType(exp ast.Expr, existingStarCount int, existingArrayCount int) *FieldType {
	switch v := exp.(type) {
	case *ast.Ident:
		if isPrimitive(v) {
			log.Print(&FieldType{
				Name:         v.Name,
				FunctionName: supportedPrimitives[v.Name],
				Primitive:    true,
				StarCount:    existingStarCount,
				ArrayCount:   existingArrayCount,
			})
			return &FieldType{
				Name:         v.Name,
				FunctionName: supportedPrimitives[v.Name],
				Primitive:    true,
				StarCount:    existingStarCount,
				ArrayCount:   existingArrayCount,
			}
		} else {
			return &FieldType{
				Name:         v.Name,
				FunctionName: "read" + strings.Title(v.Name),
				Primitive:    false,
				StarCount:    existingStarCount,
				ArrayCount:   existingArrayCount,
			}
		}
		break
	case *ast.ArrayType:
		existingArrayCount++
		fieldType := getFieldType(v.Elt, existingStarCount, existingArrayCount)
		return fieldType
	case *ast.StarExpr:
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
