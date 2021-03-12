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
)

var fieldTemplate = map[int]string{
	decodeBaseType: `

{{if .IsPointer}}
	if b.ReadBool() {
		if ({{.Accessor}} == nil) {
			{{.Accessor}} = new({{.Type}})
		}
		*{{.Accessor}} = {{.Type}}(b.{{.DecodingMethod}}())
	}
{{else}}
	{{.Accessor}} = {{.Type}}(b.{{.DecodingMethod}}())
{{end}}
`,
	encodeBaseType: `
{{if .IsPointer}}
	if {{.Accessor}} != nil {
		b.WriteBool(true)
{{end}}

	b.{{.EncodingMethod}}({{.PrimitiveWriteCast}}({{if .IsPointer}}*{{end}}{{.Accessor}}))

{{if .IsPointer}}
	} else {
		b.WriteBool(false)
	}
{{end}}
`,

	decodeBaseTypeSlice: `

length = int(b.ReadUint64())
if length > 0 {
	{{.Accessor}} = make({{.Type}}, length)

	for i := range {{.Accessor}} {

		{{if .IsPointerComponent}}
		if b.ReadBool() {
			if ({{.Accessor}}[i] == nil) {
				{{.Accessor}}[i] = new({{.ComponentType}})		
			}
			*{{.Accessor}}[i] = {{noPointer .ComponentType}}(b.{{.DecodingMethod}}())
		}
		{{else}}
			{{.Accessor}}[i] = {{.ComponentType}}(b.{{.DecodingMethod}}())
		{{end}}
	}
}

`,
	encodeBaseTypeSlice: `

b.WriteUint64(uint64(len({{if .IsPointer}}*{{end}}{{.Accessor}})))

for i := range {{if .IsPointer}}*{{end}}{{.Accessor}} {
	{{if .IsPointerComponent}}
	if {{.Accessor}}[i] != nil {
		b.WriteBool(true)

		b.{{.EncodingMethod}}({{.PrimitiveWriteCast}}(*{{.Accessor}}[i]))

	{{else}}
		b.{{.EncodingMethod}}({{.PrimitiveWriteCast}}({{.Accessor}}[i]))
	{{end}}

	{{if .IsPointerComponent}}
	} else {
		b.WriteBool(false)
	}
	{{end}}
}

`,
	decodeStruct: `

{{if .IsPointer}}
	if b.ReadBool() {
		if ({{.Accessor}} == nil) {
			{{.Accessor}} = new({{.Type}})
		}
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

length = int(b.ReadUint64())
if length > 0 {
	{{.Accessor}} = make({{.RawType}}, length)

	for i := range {{.Accessor}} {

		{{if .IsPointerComponent}}
		if b.ReadBool() {
			if ({{.Accessor}}[i] == nil) {
				{{.Accessor}}[i] = new({{.ComponentType}})
			}
			{{noPointer .Accessor}}[i].UnmarshalBuffer(b)
		}
		{{else}}
			(*{{.ComponentType}})(&{{.Accessor}}[i]).UnmarshalBuffer(b)
		{{end}}
	}
}


`,

	encodeStructSlice: `

b.WriteUint64(uint64(len({{if .IsPointer}}*{{end}}{{.Accessor}})))

for i := range {{if .IsPointer}}*{{end}}{{.Accessor}} {
	{{if .IsPointerComponent}}
	if {{.Accessor}}[i] != nil {
		b.WriteBool(true)
	{{end}}

	{{noPointer .Accessor}}[i].MarshalBuffer(b)

	{{if .IsPointerComponent}}
	} else {
		b.WriteBool(false)
	}
	{{end}}
}

`, decodeAliasBaseType: `

*{{.Accessor}} = {{.Name}}({{.Derived}}(b.{{.DecodingMethod}}()))

`, decodeAliasBaseTypeSlice: `

length = int(b.ReadUint64())
if length > 0 {
	temp := make([]{{.ComponentType}}, length)

	for i := range temp {

		{{if .IsPointerComponent}}
		if b.ReadBool() {
			if temp[i] == nil {
				temp[i] = new({{noPointer .ComponentType}})		
			}
			*temp[i] = {{noPointer .ComponentType}}(b.{{.DecodingMethod}}())
		}
		{{else}}
			temp[i] = {{.ComponentType}}(b.{{.DecodingMethod}}())
		{{end}}
	}

	*{{.Accessor}} = {{.Name}}(temp)
}
`, encodeAliasBaseTypeSlice: `

b.WriteUint64(uint64(len(*{{.Accessor}})))

temp := []{{.ComponentType}}(*{{.Accessor}})

for i := range temp {
	{{if .IsPointerComponent}}
	if temp[i] != nil {
		b.WriteBool(true)

		b.{{.EncodingMethod}}({{.PrimitiveWriteCast}}(*temp[i]))

	{{else}}
		b.{{.EncodingMethod}}({{.PrimitiveWriteCast}}(temp[i]))
	{{end}}

	{{if .IsPointerComponent}}
	} else {
		b.WriteBool(false)
	}
	{{end}}
}
`,
	encodeAliasBaseType: `

b.{{.EncodingMethod}}({{.PrimitiveWriteCast}}({{.Derived}}(*{{.Accessor}})))

`,
	decodeAliasStruct: `

(*{{.Derived}})({{.Accessor}}).UnmarshalBuffer(b)

`,
	encodeAliasStruct: `

(*{{.Derived}})({{.Accessor}}).MarshalBuffer(b)

`,
	decodeAliasStructSlice: `

length = int(b.ReadUint64())
if length > 0 {
	temp := make([]{{.ComponentType}}, length)

	for i := range temp {

		{{if .IsPointerComponent}}
		if b.ReadBool() {
			temp[i] = new({{noPointer .ComponentType}})		
			({{.ComponentType}})(temp[i]).UnmarshalBuffer(b)

		}
		{{else}}
			(*{{.ComponentType}})(&temp[i]).UnmarshalBuffer(b)
		{{end}}
	}

	*{{.Accessor}} = {{.Name}}(temp)
}

`, encodeAliasStructSlice: `

b.WriteUint64(uint64(len(*{{.Accessor}})))

temp := []{{.ComponentType}}(*{{.Accessor}})

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
`,
}

const (
	fileCode = iota
	encodingStructType
	baseTypeSlice
	structTypeSlice
	embeddedStructInit
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
	{{end}}

	{{.DecodingCases}}	
}
	return b.Err()
}`,
	embeddedStructInit: `if {{.Accessor}} == nil { 
		{{.Accessor}} = {{.Init}}
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
	return writer.String(), err
}

func expandFieldTemplate(key int, data interface{}) (string, error) {
	return expandTemplate("fieldTemplate", fieldTemplate, key, data)
}

func expandBlockTemplate(key int, data interface{}) (string, error) {
	return expandTemplate("blockTemplate", blockTemplate, key, data)
}
