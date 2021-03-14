package codegen

import (
	"fmt"
	"go.sia.tech/encodegen/internal/toolbox"
	"strings"
)

//Field represents a field.
type Field struct {
	Accessor           string
	Alias              string //object in function (un)marshaler definition
	Derived            string // if the type is a type alias the original type goes here
	Type               string
	RawType            string
	ComponentType      string
	RawComponentType   string
	IsPointerComponent bool

	DecodingMethod string
	EncodingMethod string

	PrimitiveWriteCast string

	IsPointer bool
	IsSlice   bool
	Iterator string

	AnonymousChildFields []*toolbox.FieldInfo
}

//NewField returns a new field
func NewField(owner *Struct, field *toolbox.FieldInfo, fieldType *toolbox.TypeInfo) (*Field, error) {
	// fmt.Printf("\nOwner: {%+v}\nField: {%+v}\nFieldType: {%+v}\n", owner, field, fieldType)
	// fmt.Printf("Owner Alias: %s, Field Name: %s\n", owner.Alias, field.Name)

	result := &Field{
		RawType:              field.TypeName,
		IsPointer:            field.IsPointer,
		Type:                 field.TypeName,
		Accessor:             owner.Alias + "." + field.Name,
		ComponentType:        field.ComponentType,
		IsPointerComponent:   field.IsPointerComponent,
		IsSlice:              field.IsSlice,
		Alias:                owner.Alias,
		AnonymousChildFields: field.AnonymousChildFields,
	}

	if fieldType != nil {
		// alias
		if fieldType.Derived != "" {
			result.Derived = fieldType.Derived
		}
	}

	if field.IsPointer {
		if strings.HasPrefix(result.RawType, "**") {
			return nil, fmt.Errorf("Only single pointers are supported (error found in %+v)", field)
		}
		// toolbox library does not properly label pointers to slices but we dont support these anyways
		if strings.HasPrefix(result.RawType, "*[]") {
			return nil, fmt.Errorf("Pointers to slices (%+v) are not supported", field)
		}
	}

	componentType := field.ComponentType
	if componentType == "" && (fieldType == nil || fieldType.Derived != "") {
		componentType = result.Type
		result.ComponentType = strings.TrimPrefix(strings.TrimPrefix(result.Type, "[]"), "*")
	}

	if isPrimitiveString(componentType) {
		primitive := supportedPrimitives[componentType]
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
