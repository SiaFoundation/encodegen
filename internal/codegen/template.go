package codegen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	decodeBaseType = iota
	encodeBaseType
	decodeBaseTypeSlice
	encodeBaseTypeSlice
	decodeRawType
	encodeRawType
	decodeStruct
	encodeStruct
	decodeStructSlice
	encodeStructSlice
	decodeAliasBaseType
	encodeAliasBaseType
	decodeAliasStruct
	encodeAliasStruct
	decodeAliasBaseTypeSlice
	encodeAliasBaseTypeSlice
	decodeAliasStructSlice
	encodeAliasStructSlice
	decodeAnonymousStructPointer
	encodeAnonymousStructPointer
	decodeAnonymousStructSlice
	encodeAnonymousStructSlice
)

var fieldTemplate = map[int]string{
	decodeBaseType: `
{{if .IsPointer}}
	if b.ReadBool() {
		{{if .ReuseMemory}}
		if {{.Accessor}} == nil {
		{{end}}
			{{.Accessor}} = new({{.Type}})
		{{if .ReuseMemory}}
		}
		{{end}}
		*{{.Accessor}} = {{.Type}}({{.PrimitiveFunction.ReadCast}}(b.{{.PrimitiveFunction.ReadFunction}}()))
	}
{{else}}
	{{.Accessor}} = {{.Type}}(b.{{.PrimitiveFunction.ReadFunction}}())
{{end}}
`,
	encodeBaseType: `
{{if .IsPointer}}
	if {{.Accessor}} != nil {
		b.WriteBool(true)
{{end}}
	b.{{.PrimitiveFunction.WriteFunction}}({{.PrimitiveFunction.WriteCast}}({{if .IsPointer}}*{{end}}{{.Accessor}}))
{{if .IsPointer}}
	} else {
		b.WriteBool(false)
	}
{{end}}
`,
	decodeBaseTypeSlice: `
{{if not .IsFixed}}
length = int(b.ReadUint64())
if length > 0 {
	{{if .ReuseMemory}}
	if len({{.Accessor}}) < length {
	{{end}}
	{{.Accessor}} = make({{.Type}}, length)
	{{if .ReuseMemory}}
	}
	{{end}}
	{{.Accessor}} = {{.Accessor}}[:length]
	{{end}}
	{{if and (eq .ComponentType "byte") (.IsSlice) (eq .IsPointerComponent false) (eq .IsPointer false)}}
	b.Read({{.Accessor}}{{if .IsFixed}}[:]{{end}})
	{{else}}
	for {{.Iterator}} := range {{.Accessor}} {
		{{if and .ReuseMemory (not .IsFixed)}}
		if {{.Iterator}} == length {
			break
		}
		{{end}}

		{{if .IsPointerComponent}}
		if b.ReadBool() {
			{{if .ReuseMemory}}
			if {{.Accessor}}[{{.Iterator}}] == nil {
			{{end}}
			{{.Accessor}}[{{.Iterator}}] = new({{.ComponentType}})		
			{{if .ReuseMemory}}
			}
			{{end}}
			*{{.Accessor}}[{{.Iterator}}] = {{noPointer .ComponentType}}({{.PrimitiveFunction.ReadCast}}(b.{{.PrimitiveFunction.ReadFunction}}()))
		}
		{{else}}
			{{.Accessor}}[{{.Iterator}}] = {{.ComponentType}}({{.PrimitiveFunction.ReadCast}}(b.{{.PrimitiveFunction.ReadFunction}}()))
		{{end}}
	}
	{{end}}
{{if not .IsFixed}}
}
{{end}}
`,
	encodeBaseTypeSlice: `
{{if not .IsFixed}}
b.WriteUint64(uint64(len({{.Accessor}})))
{{end}}
{{if and (eq .ComponentType "byte") (.IsSlice) (eq .IsPointerComponent false) (eq .IsPointer false)}}
b.Write({{.Accessor}}{{if .IsFixed}}[:]{{end}})
{{else}}
for {{.Iterator}} := range {{.Accessor}} {
	{{if .IsPointerComponent}}
	if {{.Accessor}}[{{.Iterator}}] != nil {
		b.WriteBool(true)
		b.{{.PrimitiveFunction.WriteFunction}}({{.PrimitiveFunction.WriteCast}}(*{{.Accessor}}[{{.Iterator}}]))
	{{else}}
		b.{{.PrimitiveFunction.WriteFunction}}({{.PrimitiveFunction.WriteCast}}({{.Accessor}}[{{.Iterator}}]))
	{{end}}
	{{if .IsPointerComponent}}
	} else {
		b.WriteBool(false)
	}
	{{end}}
}
{{end}}
`,
	decodeStruct: `
{{if .IsPointer}}
	if b.ReadBool() {
		{{if .ReuseMemory}}
		if ({{.Accessor}} == nil) {
		{{end}}
			{{.Accessor}} = new({{.Type}})
		{{if .ReuseMemory}}
		}
		{{end}}
		(*{{.Type}})({{.Accessor}}).UnmarshalBuffer(b)	
	}
{{else}}
	(*{{.Type}})(&{{.Accessor}}).UnmarshalBuffer(b)
{{end}}
`,
	encodeStruct: `
{{if .IsPointer}}
	if {{.Accessor}} != nil {
		b.WriteBool(true)
		(*{{.Type}})({{.Accessor}}).MarshalBuffer(b)
	} else {
		b.WriteBool(false)
	}
{{else}}
	(*{{.Type}})(&{{.Accessor}}).MarshalBuffer(b)
{{end}}
`,
	decodeStructSlice: `
{{if not .IsFixed}}
length = int(b.ReadUint64())
if length > 0 {
	{{if .ReuseMemory}}
	if len({{.Accessor}}) < length {
	{{end}}
	{{.Accessor}} = make({{.RawType}}, length)
	{{if .ReuseMemory}}
	}
	{{.Accessor}} = {{.Accessor}}[:length]
	{{end}}
	{{end}}
	for {{.Iterator}} := range {{.Accessor}} {
		{{if and .ReuseMemory (not .IsFixed)}}
		if {{.Iterator}} == length {
			break
		}
		{{end}}
		{{if .IsPointerComponent}}
		if b.ReadBool() {
			{{if .ReuseMemory}}
			if {{.Accessor}}[{{.Iterator}}] == nil {
			{{end}}
			{{.Accessor}}[{{.Iterator}}] = new({{.ComponentType}})
			{{if .ReuseMemory}}
			}
			{{end}}
			{{noPointer .Accessor}}[{{.Iterator}}].UnmarshalBuffer(b)
		}
		{{else}}
			(*{{.ComponentType}})(&{{.Accessor}}[{{.Iterator}}]).UnmarshalBuffer(b)
		{{end}}
	}
{{if not .IsFixed}}
}
{{end}}
`,
	encodeStructSlice: `
{{if not .IsFixed}}
b.WriteUint64(uint64(len({{.Accessor}})))
{{end}}
for {{.Iterator}} := range {{.Accessor}} {
	{{if .IsPointerComponent}}
	if {{.Accessor}}[{{.Iterator}}] != nil {
		b.WriteBool(true)
	{{end}}
	{{noPointer .Accessor}}[{{.Iterator}}].MarshalBuffer(b)
	{{if .IsPointerComponent}}
	} else {
		b.WriteBool(false)
	}
	{{end}}
}
`, decodeAliasBaseType: `
*{{.Accessor}} = {{.Name}}({{.Derived}}({{.PrimitiveFunction.ReadCast}}(b.{{.PrimitiveFunction.ReadFunction}}())))
`, decodeAliasBaseTypeSlice: `
{{if not .IsFixed}}
length = int(b.ReadUint64())
if length > 0 {
	{{if .ReuseMemory}}
	if len(*{{.Accessor}}) < length {
	{{end}}
	*{{.Accessor}} = make([]{{.ComponentType}}, length)
	{{if .ReuseMemory}}
	}
	{{end}}
	(*{{.Accessor}}) = (*{{.Accessor}})[:length]
{{end}}
	{{if and (eq .ComponentType "byte") (eq .IsPointerComponent false)}}

	{{if .IsFixed}}
	temp := [{{.FixedSize}}]{{.ComponentType}}(*{{.Accessor}})
	b.Read(temp[:])
	*{{.Accessor}} = temp

	{{else}}
	b.Read(*{{.Accessor}}{{if .IsFixed}}[:]{{end}})
	{{end}}

	{{else}}
	for i := range *{{.Accessor}} {
		{{if and .ReuseMemory (not .IsFixed)}}
		if i == length {
			continue
		}
		{{end}}

		{{if .IsPointerComponent}}
		if b.ReadBool() {
			{{if .ReuseMemory}}
			if (*{{.Accessor}})[i] == nil {
			{{end}}
			(*{{.Accessor}})[i] = new({{noPointer .ComponentType}})
			{{if .ReuseMemory}}
			}
			{{end}}
			*(*{{.Accessor}})[i] = {{noPointer .ComponentType}}({{.PrimitiveFunction.ReadCast}}(b.{{.PrimitiveFunction.ReadFunction}}()))
		}
		{{else}}
			(*{{.Accessor}})[i] = {{.ComponentType}}({{.PrimitiveFunction.ReadCast}}(b.{{.PrimitiveFunction.ReadFunction}}()))
		{{end}}
	}
	{{end}}
{{if not .IsFixed}}
}
{{end}}
`, encodeAliasBaseTypeSlice: `
{{if not .IsFixed}}
b.WriteUint64(uint64(len(*{{.Accessor}})))
{{end}}
{{if and (eq .ComponentType "byte") (eq .IsPointerComponent false)}}

{{if .IsFixed}}
temp := [{{.FixedSize}}]{{.ComponentType}}(*{{.Accessor}})
b.Write([]byte(temp{{if .IsFixed}}[:]{{end}}))
{{else}}
b.Write([]{{.ComponentType}}(*{{.Accessor}}))
{{end}}

{{else}}
temp := [{{if .IsFixed}}{{.FixedSize}}{{end}}]{{.ComponentType}}(*{{.Accessor}})
for i := range temp {
	{{if .IsPointerComponent}}
	if temp[i] != nil {
		b.WriteBool(true)
		b.{{.PrimitiveFunction.WriteFunction}}({{.PrimitiveFunction.WriteCast}}(*temp[i]))
	{{else}}
		b.{{.PrimitiveFunction.WriteFunction}}({{.PrimitiveFunction.WriteCast}}(temp[i]))
	{{end}}
	{{if .IsPointerComponent}}
	} else {
		b.WriteBool(false)
	}
	{{end}}
}
{{end}}
`,
	encodeAliasBaseType: `
b.{{.PrimitiveFunction.WriteFunction}}({{.PrimitiveFunction.WriteCast}}({{.Derived}}(*{{.Accessor}})))
`,
	decodeAliasStruct: `
(*{{.Derived}})({{.Accessor}}).UnmarshalBuffer(b)
`,
	encodeAliasStruct: `
(*{{.Derived}})({{.Accessor}}).MarshalBuffer(b)
`,
	decodeAliasStructSlice: `
{{if not .IsFixed}}
length = int(b.ReadUint64())
if length > 0 {
	{{if .ReuseMemory}}
	if len(*{{.Accessor}}) < length {
	{{end}}
	*{{.Accessor}} = make([]{{.ComponentType}}, length)
	{{if .ReuseMemory}}
	}
	{{end}}
	(*{{.Accessor}}) = (*{{.Accessor}})[:length]
{{end}}
	for i := range *{{.Accessor}} {
		{{if and .ReuseMemory (not .IsFixed)}}
		if i == length {
			continue
		}
		{{end}}

		{{if .IsPointerComponent}}
		if b.ReadBool() {
			{{if .ReuseMemory}}
			if (*{{.Accessor}})[i] == nil {
			{{end}}
			(*{{.Accessor}})[i] = new({{noPointer .ComponentType}})		
			{{if .ReuseMemory}}
			}
			{{end}}
			({{.ComponentType}})((*{{.Accessor}})[i]).UnmarshalBuffer(b)
		}
		{{else}}
			(*{{.ComponentType}})(&(*{{.Accessor}})[i]).UnmarshalBuffer(b)
		{{end}}
	}
{{if not .IsFixed}}
}
{{end}}
`, encodeAliasStructSlice: `
{{if not .IsFixed}}
b.WriteUint64(uint64(len(*{{.Accessor}})))
{{end}}
temp := [{{if .IsFixed}}{{.FixedSize}}{{end}}]{{.ComponentType}}(*{{.Accessor}})
for i := range temp {
	{{if .IsPointerComponent}}
	if temp[i] != nil {
		b.WriteBool(true)
		({{.ComponentType}})(temp[i]).MarshalBuffer(b)
	{{else}}
		(*{{.ComponentType}})(&temp[i]).MarshalBuffer(b)
	{{end}}
	{{if .IsPointerComponent}}
	} else {
		b.WriteBool(false)
	}
	{{end}}
}
`, decodeAnonymousStructPointer: `
if b.ReadBool() {
	{{if .ReuseMemory}}
	if {{.Accessor}} == nil {
	{{end}}
		{{.Accessor}} = new({{noPointer .Type}})
	{{if .ReuseMemory}}
	}
	{{end}}
	{{.Cases}}
}
`, encodeAnonymousStructPointer: `
if {{.Accessor}} != nil {
	b.WriteBool(true)
	{{.Cases}}
} else {
	b.WriteBool(false)
}
`, decodeAnonymousStructSlice: `
{{if not .IsFixed}}
length = int(b.ReadUint64())
if length > 0 {
	{{if .ReuseMemory}}
	if len({{.Accessor}}) < length {
	{{end}}
	{{.Accessor}} = make({{.Type}}, length) 
	{{if .ReuseMemory}}
	}
	{{end}}
	{{.Accessor}} = {{.Accessor}}[:length]
{{end}}
	for {{.Iterator}} := range {{.Accessor}} {
		{{if and .ReuseMemory (not .IsFixed)}}
		if {{.Iterator}} == length {
			break
		}
		{{end}}
		{{if .IsPointerComponent}}
		if b.ReadBool() {
			{{if .ReuseMemory}}
			if {{.Accessor}}[{{.Iterator}}] == nil {
			{{end}}
			{{.Accessor}}[{{.Iterator}}] = new({{.ComponentType}})
			{{if .ReuseMemory}}
			}
			{{end}}

		{{end}}
		{{.Cases}}
		{{if .IsPointerComponent}}
		}
		{{end}}
	}
{{if not .IsFixed}}
}
{{end}}
`, encodeAnonymousStructSlice: `
{{if not .IsFixed}}
b.WriteUint64(uint64(len({{.Accessor}})))
{{end}}
for {{.Iterator}} := range {{.Accessor}} {
	{{if .IsPointerComponent}}
	if {{.Accessor}}[{{.Iterator}}] != nil {
		b.WriteBool(true)
	{{end}}
	{{.Cases}}
	{{if .IsPointerComponent}}
	} else {
		b.WriteBool(false)
	}
	{{end}}
}
`,
}

const (
	fileCode = iota
	encodingStructType
	baseTypeSlice
	structTypeSlice
	typeSlice
)

var blockTemplate = map[int]string{
	fileCode: `// Code generated by encodegen. DO NOT EDIT.
package {{base .Pkg}}
import (
	{{.Imports}}
)
{{.Code}}
`,
	encodingStructType: `// MarshalBuffer implements MarshalerBuffer
func ({{.Receiver}}) MarshalBuffer(b *encodegen.ObjBuffer) {
if {{.Alias}} != nil {
	{{.EncodingCases}}
}
}
// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func ({{.Receiver}}) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
if {{.Alias}} != nil {
	{{if .HasSlice}}
	var length int = 0
	_ = length
	{{end}}

	{{.DecodingCases}}	
}
	return b.Err()
}`,
}

func noPointer(s string) string {
	return strings.TrimPrefix(s, "*")
}

func expandTemplate(namespace string, dictionary map[int]string, key int, data interface{}) (string, error) {
	var id = fmt.Sprintf("%v_%v", namespace, key)
	textTemplate, ok := dictionary[key]
	if !ok {
		return "", fmt.Errorf("failed to lookup template for %v.%v", namespace, key)
	}

	// add iter function to allow us to conveniently repeat code n times
	temlate, err := template.New(id).Funcs(template.FuncMap{"noPointer": noPointer, "base": filepath.Base}).Parse(textTemplate)
	if err != nil {
		return "", fmt.Errorf("fiailed to parse template %v %v, due to %v", namespace, key, err)
	}

	writer := new(bytes.Buffer)
	err = temlate.Execute(writer, data)
	// fmt.Printf("Call template with key, %d, data: %+v\n", key, data)
	if err != nil {
		panic(err)
	}
	return writer.String(), err
}

func expandFieldTemplate(key int, data interface{}) (string, error) {
	return expandTemplate("fieldTemplate", fieldTemplate, key, data)
}

func expandBlockTemplate(key int, data interface{}) (string, error) {
	return expandTemplate("blockTemplate", blockTemplate, key, data)
}
