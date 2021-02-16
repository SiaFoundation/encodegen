package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"log"
	"strings"
)

// https://github.com/lukechampine/few/blob/master/generate.go

type generator struct {
	buf bytes.Buffer
}

func (g *generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Fatalln("invalid Go generated:", err)
	}
	return src
}

func (g *generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

type FieldType struct {
	Name         string
	FunctionName string
	Primitive    bool
	StarCount    int
}

func writeFieldAssignmentStatements(g *generator, structs map[string]*ast.StructType, prefixs []string, currentStructName string) {
	for _, structField := range structs[currentStructName].Fields.List {
		if len(structField.Names) > 0 {
			field := getFieldType(structField.Type, 0)
			if field == nil {
				continue
			}

			// nil pointer checks for written pointers
			for i := 0; i < field.StarCount; i++ {
				g.Printf("if b.readBool() {\n")
			}

			if field.Primitive {
				if len(prefixs) == 0 {
					g.Printf("r.%s = ", structField.Names[0].Name)
				} else {
					g.Printf("r.%s.%s = ", strings.Join(prefixs, "."), structField.Names[0].Name)
				}
				g.Printf("%s(b.%s())\n", field.Name, field.FunctionName)
			} else {
				prefixs = append(prefixs, structField.Names[0].Name)
				writeFieldAssignmentStatements(g, structs, prefixs, field.Name)
				prefixs = prefixs[:len(prefixs)-1]
			}
			// close if statements for nil pointer checking
			for i := 0; i < field.StarCount; i++ {
				g.Printf("}\n")
			}
		}
	}
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

func getFieldType(exp ast.Expr, existingStarCount int) *FieldType {
	switch v := exp.(type) {
	case *ast.Ident:
		if isPrimitive(v) {
			log.Print(&FieldType{
				Name:         v.Name,
				FunctionName: supportedPrimitives[v.Name],
				Primitive:    true,
				StarCount:    existingStarCount,
			})
			return &FieldType{
				Name:         v.Name,
				FunctionName: supportedPrimitives[v.Name],
				Primitive:    true,
				StarCount:    existingStarCount,
			}
		} else {
			return &FieldType{
				Name:         v.Name,
				FunctionName: "read" + strings.Title(v.Name),
				Primitive:    false,
				StarCount:    existingStarCount,
			}
		}
		break
	// case *ast.ArrayType:
	// 	return getArrayType(v)
	case *ast.StarExpr:
		existingStarCount++
		fieldType := getFieldType(v.X, existingStarCount)
		return fieldType
		break
	default:
		log.Printf("WARNING: This struct contains unsupported data (%+v)", exp)
		return nil
		break
	}
	return nil
}
