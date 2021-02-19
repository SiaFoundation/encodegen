package main

import (
	"go/ast"
	"log"
)

const ReadBoolFunction = "readBool"
const ReadIntFunction = "readUint64"
const ReadStringFunction = "readPrefixedBytes"
const ReadByteFunction = "readByte"

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
			return &FieldType{
				Name:         v.Name,
				FunctionName: supportedPrimitives[v.Name],
				Primitive:    true,
				StarCount:    existingStarCount,
				ArrayCount:   existingArrayCount,
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
