package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"log"
	"strconv"
	"strings"
)

// https://github.com/lukechampine/few/blob/master/generate.go

type generator struct {
	buf bytes.Buffer
}

// clear the buffer
func (g *generator) reset() {
	g.buf.Reset()
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

func nextRangeIdentifier(currentIdentifier string) string {
	// this function allows the generated the code to iterate over slices of structs that have slices within them without having iteration identifiers conflict (i.e., there'd be multiple "range i := r.Fields"s)

	idSplit := strings.Split(currentIdentifier, "i")
	if len(idSplit) != 2 {
		return "'"
	}
	if idSplit[1] != "" {
		num, err := strconv.Atoi(idSplit[1])
		if err != nil {
			return ""
		}
		return fmt.Sprintf("i%d", num+1)
	} else {
		return "i1"
	}
}

// I debated whether to use these or not but I felt that it basically just complicates the code for no real gain
// const bufferIdentifier string = "b"
// const structIdentifier string = "r"

func (g *generator) writeUnmarshalBufferStatements(structs map[string]*ast.StructType, prefixs []string, currentIdentifier string, currentStructType *ast.StructType) {
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

			// generate the stuff before the equals sign in the assignment, i.e. "r.A" or "r.A.C" by joining the prefixs and the current field name
			assignmentName = fmt.Sprintf("%s.%s", strings.Join(prefixs, "."), structField.Names[0].Name)

			if field.Primitive {
				// generate the stuff after the equals sign.  we first cast the type (note that this may generate redundancies like uint64(b.readUint64()) but the compiler is ok with this) then call the appropriate read function
				assignmentValue = fmt.Sprintf("%s(b.%s())\n", field.Name, field.PrimitiveFunctions.ReadFunction)

				g.Printf("%s = ", assignmentName)
				if field.ArrayCount > 0 {
					// initialize the slice with make
					g.Printf("make(%s%s, int(b.readUint64()))\n", strings.Repeat("[]", 1), field.Name)
					identifier := nextRangeIdentifier(currentIdentifier)

					g.Printf("for %s := range %s {\n", identifier, assignmentName)
					// if its a pointer do the nil check first
					for i := 0; i < field.StarCount; i++ {
						g.Printf("if b.readBool() {\n")
					}
					// write the actual assignment
					g.Printf("%s[%s] = %s\n", assignmentName, identifier, assignmentValue)
					// close the nil checks
					for i := 0; i < field.StarCount; i++ {
						g.Printf("}\n")
					}
					// close the range
					g.Printf("}\n")
				} else {
					g.Printf("%s", assignmentValue)
				}
			} else {
				// array of structs
				if field.ArrayCount > 0 {
					g.Printf("%s = make(%s%s, int(b.readUint64()))\n", assignmentName, strings.Repeat("[]", 1), field.Name)
					identifier := nextRangeIdentifier(currentIdentifier)
					g.Printf("for %s := range %s {\n", identifier, assignmentName)
					prefixs = append(prefixs, fmt.Sprintf("%s[%s]", structField.Names[0].Name, identifier))
					for i := 0; i < field.StarCount; i++ {
						g.Printf("if b.readBool() {\n")
					}
					g.writeUnmarshalBufferStatements(structs, prefixs, identifier, structs[field.Name])
					for i := 0; i < field.StarCount; i++ {
						g.Printf("}\n")
					}
					prefixs = prefixs[:len(prefixs)-1]
					g.Printf("}\n")
				} else {
					// single struct
					prefixs = append(prefixs, structField.Names[0].Name)
					g.writeUnmarshalBufferStatements(structs, prefixs, currentIdentifier, structs[field.Name])
					prefixs = prefixs[:len(prefixs)-1]
				}
			}

			// close if statements for nil pointer checking
			if field.ArrayCount == 0 {
				for i := 0; i < field.StarCount; i++ {
					g.Printf("}\n")
				}
			}

		}
	}
}

func (g *generator) writeMarshalBufferStatements(structs map[string]*ast.StructType, prefixs []string, currentIdentifier string, currentStructType *ast.StructType) {
	if g == nil || currentStructType == nil {
		return
	}
	for _, structField := range currentStructType.Fields.List {
		if len(structField.Names) > 0 {

			field := getFieldType(structField.Type, 0, 0)
			if field == nil {
				continue
			}

			fieldName, assignmentValue := "", ""

			// nil pointer checks for written pointers
			if field.ArrayCount == 0 {
				for i := 0; i < field.StarCount; i++ {
					g.Printf("if b.readBool() {\n")
				}
			}

			// generate the stuff before the equals sign in the assignment, i.e. "r.A" or "r.A.C" by joining the prefixs and the current field name
			fieldName = fmt.Sprintf("%s.%s", strings.Join(prefixs, "."), structField.Names[0].Name)

			if field.Primitive {
				// generate the stuff after the equals sign.  we first cast the type (note that this may generate redundancies like uint64(b.readUint64()) but the compiler is ok with this) then call the appropriate read function
				assignmentValue = fmt.Sprintf("b.%s(%s(%s))\n", field.PrimitiveFunctions.WriteFunction, field.PrimitiveFunctions.WriteCast, fieldName)

				if field.ArrayCount > 0 {
					// initialize the slice with make
					g.Printf("b.writePrefix(len(%s))\n", fieldName)
					identifier := nextRangeIdentifier(currentIdentifier)

					g.Printf("for %s := range %s {\n", identifier, fieldName)
					// if its a pointer do the nil check first
					for i := 0; i < field.StarCount; i++ {
						g.Printf("if b.readBool() {\n")
					}
					g.Printf("b.%s(%s(%s[%s]))\n", field.PrimitiveFunctions.WriteFunction, field.PrimitiveFunctions.WriteCast, fieldName, identifier)

					// close the nil checks
					for i := 0; i < field.StarCount; i++ {
						g.Printf("}\n")
					}
					// close the range
					g.Printf("}\n")
				} else {
					g.Printf("%s", assignmentValue)
				}
			} else {
				// array of structs
				if field.ArrayCount > 0 {
					g.Printf("%s = make(%s%s, int(b.readUint64()))\n", fieldName, strings.Repeat("[]", 1), field.Name)
					identifier := nextRangeIdentifier(currentIdentifier)
					g.Printf("for %s := range %s {\n", identifier, fieldName)
					prefixs = append(prefixs, fmt.Sprintf("%s[%s]", structField.Names[0].Name, identifier))
					for i := 0; i < field.StarCount; i++ {
						g.Printf("if b.readBool() {\n")
					}
					g.writeMarshalBufferStatements(structs, prefixs, identifier, structs[field.Name])
					for i := 0; i < field.StarCount; i++ {
						g.Printf("}\n")
					}
					prefixs = prefixs[:len(prefixs)-1]
					g.Printf("}\n")
				} else {
					// single struct
					prefixs = append(prefixs, structField.Names[0].Name)
					g.writeMarshalBufferStatements(structs, prefixs, currentIdentifier, structs[field.Name])
					prefixs = prefixs[:len(prefixs)-1]
				}
			}

			// close if statements for nil pointer checking
			if field.ArrayCount == 0 {
				for i := 0; i < field.StarCount; i++ {
					g.Printf("}\n")
				}
			}

		}
	}
}
