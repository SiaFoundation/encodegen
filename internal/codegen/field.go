package codegen

import (
	"fmt"
	"go.sia.tech/encodegen/internal/toolbox"
	"strings"
)

type Field struct {
	Name                 string // original struct field name
	Alias                string // object in function (un)marshaler definition
	Accessor             string // if a struct field is called ABC then this will be something like m.ABC
	Derived              string // if the type is a type alias the original type goes here
	Type                 string
	ComponentType        string // []byte -> byte, []*byte -> byte, []SubMessage -> SubMessage, []*SubMessage -> SubMessage
	RawComponentType     string // []*byte -> *byte, []byte -> byte
	IsPointerComponent   bool   // []byte -> false, []SubMessage -> false, []*byte -> true, []*SubMessage -> true
	PrimitiveFunctions   PrimitiveFunctions
	IsPointer            bool   // byte -> false, *byte -> true, []*byte -> false
	IsSlice              bool   // really should be IsArray, byte -> false, []byte -> true, [40]byte -> true
	IsFixed              bool   // []byte -> false, [40]byte -> true
	FixedSize            int    // [40]byte -> 40
	Iterator             string // used to prevent the same iterator from being used in loops in anonymous structs containing arrays
	ReuseMemory          bool   // copy of the per struct setting
	AnonymousChildFields []*toolbox.FieldInfo
}

//NewField returns a new field
func NewField(owner *Struct, field *toolbox.FieldInfo, fieldType *toolbox.TypeInfo) (*Field, error) {
	result := &Field{
		Name:                 field.Name,
		IsPointer:            field.IsPointer,
		Type:                 field.TypeName,
		Accessor:             owner.Alias + "." + field.Name,
		ComponentType:        field.ComponentType,
		IsPointerComponent:   field.IsPointerComponent,
		IsSlice:              field.IsSlice,
		Alias:                owner.Alias,
		AnonymousChildFields: field.AnonymousChildFields,
		IsFixed:              field.IsFixed,
		FixedSize:            field.FixedSize,
	}

	if fieldType != nil {
		// alias
		if fieldType.Derived != "" {
			result.Derived = fieldType.Derived
		}
	}

	if field.IsPointer {
		if strings.Contains(result.Type, "**") {
			return nil, fmt.Errorf("Only single pointers are supported (error found in %+v)", field)
		}
		// toolbox library does not properly label pointers to slices but we dont support these anyways
		if strings.Contains(result.Type, "*[]") {
			return nil, fmt.Errorf("Pointers to slices are not supported (error found in %+v)", field)
		}
	}

	componentType := field.ComponentType
	if componentType == "" && (fieldType == nil || fieldType.Derived != "") {
		componentType = result.Type
		result.ComponentType = strings.TrimPrefix(strings.TrimPrefix(result.Type, "[]"), "*")
	}

	if isPrimitiveString(componentType) {
		result.PrimitiveFunctions = supportedPrimitives[componentType]
	}

	if result.IsPointerComponent {
		result.RawComponentType = "*" + result.ComponentType
	} else {
		result.RawComponentType = result.ComponentType
	}
	return result, nil
}
