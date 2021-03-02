package codegen

import (
	"fmt"
	"github.com/viant/toolbox"
	"strings"
)

//Field represents a field.
type Field struct {
	Accessor           string
	Alias              string //object in function (un)marshaler definition
	AliasDerived string // if the type is a type alias the original type goes here
	Type               string
	RawType            string
	ComponentType      string
	RawComponentType   string
	IsPointerComponent bool

	DecodingMethod     string
	EncodingMethod     string

	PrimitiveWriteCast string

	IsAnonymous     bool
	IsPointer       bool
	IsSlice         bool
}

//NewField returns a new field
func NewField(owner *Struct, field *toolbox.FieldInfo, fieldType *toolbox.TypeInfo) (*Field, error) {
	var result = &Field{
		IsAnonymous:        field.IsAnonymous,
		RawType:            field.TypeName,
		IsPointer:          field.IsPointer,
		Type:               field.TypeName,
		Accessor:           owner.Alias + "." + field.Name,
		ComponentType:      field.ComponentType,
		IsPointerComponent: field.IsPointerComponent,
		IsSlice:            field.IsSlice,
		Alias:              owner.Alias,
	}

	if fieldType != nil && fieldType.Derived != "" {
		result.AliasDerived = fieldType.Derived
	}

	if field.IsPointer {
		if strings.Contains(result.RawType, "**") {
			return nil, fmt.Errorf("Only single pointers are supported (error found in %+v)", field)
		}
		// toolbox library does not properly label pointers to slices but we dont support these anyways
		if strings.Contains(result.RawType, "[]") {
			return nil, fmt.Errorf("Pointers to slices (%+v) are not supported", field)
			// field.IsSlice = true
			// result.IsSlice = field.IsSlice

			// arraySplit := strings.Split(result.RawType, "[]")
			// for _, split := range arraySplit[1:] {
			// 	if strings.Contains(split, "*") {
			// 		field.IsPointerComponent = true
			// 		result.IsPointerComponent = field.IsPointerComponent
			// 	}
			// }
		}
	}

	if field.IsPointer && field.IsSlice {
		field.ComponentType = strings.Replace(result.Type, "*[]", "", -1)
		result.ComponentType = field.ComponentType
	}

	encodingMethod := field.ComponentType
	if encodingMethod == "" {
		encodingMethod = result.Type
	}

	if isPrimitiveString(encodingMethod) {
		primitive := supportedPrimitives[encodingMethod]
		result.EncodingMethod = primitive.WriteFunction
		result.DecodingMethod = primitive.ReadFunction
		result.PrimitiveWriteCast = primitive.WriteCast
	}

	if result.IsPointerComponent {
		result.RawComponentType = "*" + result.ComponentType
	} else {
		result.RawComponentType = result.ComponentType
	}

	return result, nil
}
