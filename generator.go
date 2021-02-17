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
	ArrayCount   int
}

func writeFieldAssignmentStatements(g *generator, structs map[string]*ast.StructType, prefixs []string, currentStructType *ast.StructType) {
	if g == nil || currentStructType == nil {
		return
	}
	for _, structField := range currentStructType.Fields.List {
		if len(structField.Names) > 0 {
			field := getFieldType(structField.Type, 0, 0)
			if field == nil {
				continue
			}
			assignmentName, assignmentValue := "", ""

			// nil pointer checks for written pointers
			if field.ArrayCount == 0 {
				for i := 0; i < field.StarCount; i++ {
					g.Printf("if b.readBool() {\n")
				}
			}

			if len(prefixs) == 0 {
				assignmentName = fmt.Sprintf("r.%s", structField.Names[0].Name)
			} else {
				assignmentName = fmt.Sprintf("r.%s.%s", strings.Join(prefixs, "."), structField.Names[0].Name)
			}

			if field.Primitive {
				assignmentValue = fmt.Sprintf("%s(b.%s())\n", field.Name, field.FunctionName)

				g.Printf("%s = ", assignmentName)
				if field.ArrayCount > 0 {
					g.Printf("make(%s%s, int(b.readUint64()))\n", strings.Repeat("[]", 1), field.Name)
					g.Printf("for i := range %s {\n", assignmentName)
					for i := 0; i < field.StarCount; i++ {
						g.Printf("if b.readBool() {\n")
					}
					g.Printf("%s[i] = %s\n", assignmentName, assignmentValue)
					for i := 0; i < field.StarCount; i++ {
						g.Printf("}\n")
					}
					g.Printf("}\n")
				} else {
					g.Printf("%s", assignmentValue)
				}
			} else {
				if field.ArrayCount > 0 {
					g.Printf("%s = make(%s%s, int(b.readUint64()))\n", assignmentName, strings.Repeat("[]", 1), field.Name)
					g.Printf("for i := range %s {\n", assignmentName)
					prefixs = append(prefixs, structField.Names[0].Name+"[i]")
					for i := 0; i < field.StarCount; i++ {
						g.Printf("if b.readBool() {\n")
					}
					writeFieldAssignmentStatements(g, structs, prefixs, structs[field.Name])
					for i := 0; i < field.StarCount; i++ {
						g.Printf("}\n")
					}
					prefixs = prefixs[:len(prefixs)-1]
					g.Printf("}\n")
				} else {
					prefixs = append(prefixs, structField.Names[0].Name)
					writeFieldAssignmentStatements(g, structs, prefixs, structs[field.Name])
					prefixs = prefixs[:len(prefixs)-1]
				}
			}

			if field.ArrayCount == 0 {

				// close if statements for nil pointer checking
				for i := 0; i < field.StarCount; i++ {
					g.Printf("}\n")
				}
			}
		}
	}
}
