package codegen

import (
	"fmt"
	"strings"

	"github.com/viant/toolbox"
)

type Struct struct {
	*toolbox.TypeInfo
	referenced *toolbox.TypeInfo
	*Generator
	Alias string
	Init  string
	Body  string
}

//Generate generates decoderCode + structRelease + encoderCode
func (s *Struct) Generate() (string, error) {
	return s.generateEncoding(s.TypeInfo)
}

func (s *Struct) generateEncoding(structInfo *toolbox.TypeInfo) (string, error) {
	var initEmbedded, decodingCases, err = s.generateFieldDecoding(structInfo.Fields())
	if err != nil {
		return "", err
	}

	encodingCases, err := s.generateFieldEncoding(structInfo.Fields())
	if err != nil {
		return "", err
	}

	if structInfo.IsDerived {
		// if !isPrimitiveString(structInfo.Derived) && !isPrimitiveArrayString(structInfo.Derived) {
		decodingCases, encodingCases, err = s.generateAliasCases(structInfo)
		if err != nil {
			return "", err
		}
	}

	var data = struct {
		Receiver string
		Alias         string
		InitEmbedded  string
		EncodingCases string
		DecodingCases string
		FieldCount    int
		HasSlice bool
	}{
		Receiver: s.Alias + " *" + s.Name,
		DecodingCases: strings.Join(decodingCases, "\n"),
		EncodingCases: strings.Join(encodingCases, "\n"),
		FieldCount:    len(decodingCases),
		InitEmbedded:  initEmbedded,
		Alias:         s.Alias,
		HasSlice: hasSlice(structInfo.Fields()),
	}
	return expandBlockTemplate(encodingStructType, data)
}

func hasSlice(fields []*toolbox.FieldInfo) bool {
	for _, field := range fields {
		if field.IsSlice {
			return true
		}
	}
	return false
}

func (s *Struct) generateAliasCases(structInfo *toolbox.TypeInfo) ([]string, []string, error) {
	var newStructInfo = struct {
		Accessor string
		Derived string
		Name string
		ReadFunction string
		WriteFunction string
		WriteCast string
	}{
		Accessor: s.Alias,
		Derived: structInfo.Derived,
		Name: structInfo.Name,
		ReadFunction: supportedPrimitives[structInfo.Derived].ReadFunction,
		WriteFunction: supportedPrimitives[structInfo.Derived].WriteFunction,
		WriteCast: supportedPrimitives[structInfo.Derived].WriteCast,
	}

	if !isPrimitiveString(structInfo.Derived) && !isPrimitiveArrayString(structInfo.Derived) {
		decode, err := expandFieldTemplate(decodeAliasStruct, newStructInfo)
		if err != nil {
			return nil, nil, err
		}	
		encode, err := expandFieldTemplate(encodeAliasStruct, newStructInfo)
		if err != nil {
			return nil, nil, err
		}

		return []string{decode}, []string{encode}, nil
	} else {
		decode, err := expandFieldTemplate(decodeAliasBaseType, newStructInfo)
		if err != nil {
			return nil, nil, err
		}	
		encode, err := expandFieldTemplate(encodeAliasBaseType, newStructInfo)
		if err != nil {
			return nil, nil, err
		}

		return []string{decode}, []string{encode}, nil
		// return []string{
		// 	fmt.Sprintf("*%s = %s(%s(b.%s()))", s.Alias, structInfo.Name, structInfo.Derived, primitiveType.ReadFunction)}, []string{fmt.Sprintf("b.%s(%s(%s(*%s)))", primitiveType.WriteFunction, primitiveType.WriteCast, structInfo.Derived, s.Alias)}, nil
	}
}

func (s *Struct) generateFieldDecoding(fields []*toolbox.FieldInfo) (string, []string, error) {

	fieldCases := []string{}
	var initCode = ""
	for i := range fields {
		var templateKey = -1
		fieldTypeInfo := s.Type(fields[i].TypeName)
		field, err := NewField(s, fields[i], fieldTypeInfo)
		if err != nil {
			return "", nil, err
		}
		if fieldTypeInfo != nil {
			if err = s.generateStructCode(fieldTypeInfo.Name); err != nil {
				return "", nil, err
			}
		}

		if field.IsAnonymous {
			if fieldTypeInfo != nil {
				if field.IsPointer {
					init, err := expandBlockTemplate(embeddedStructInit, field)
					if err != nil {
						return "", nil, err
					}
					initCode += init
				}
				init, embeddedCases, err := s.generateFieldDecoding(fieldTypeInfo.Fields())
				if err != nil {
					return "", nil, err
				}
				initCode += init
				fieldCases = append(fieldCases, embeddedCases...)
			}
			continue
		}

	main:
		switch {
		case isPrimitiveString(field.Type):
			templateKey = decodeBaseType
		case isPrimitiveArrayString(field.Type):
			templateKey = decodeBaseTypeSlice
			s.generatePrimitiveArray(field)
		default:

			if fieldTypeInfo != nil {
				if !(field.IsSlice || fieldTypeInfo.IsSlice) {

					templateKey = decodeStruct
					break main
				}

				if isPrimitiveString(fieldTypeInfo.ComponentType) {
					s.generatePrimitiveArray(field)
					templateKey = decodeBaseTypeSlice
					break main

				}

				if err = s.generateStructCode(field.ComponentType); err != nil {
					return "", nil, err
				}

				templateKey = decodeStructSlice
				if err = s.generateObjectArray(field); err != nil {
					return "", nil, err
				}

				break main
			} else if field.IsSlice {
					templateKey = decodeStructSlice
					if err = s.generateObjectArray(field); err != nil {
						return "", nil, err
					}
			} else {
				// templateKey = decodeUnknown
				templateKey = decodeStruct
				// return "", nil, fmt.Errorf("Unknown type %s for field %s", field.Type, field.Name)
			}
		}
		if templateKey != -1 {
			decodingCase, err := expandFieldTemplate(templateKey, field)
			if err != nil {
				return "", nil, err
			}
			fieldCases = append(fieldCases, decodingCase)
		}

	}
	return initCode, fieldCases, nil
}

func (s *Struct) generateEmbeddedFieldEncoding(field *Field, fieldTypeInfo *toolbox.TypeInfo) ([]string, error) {
	var result = []string{}
	if fieldTypeInfo != nil {
		embeddedCases, err := s.generateFieldEncoding(fieldTypeInfo.Fields())
		if err != nil {
			return nil, err
		}
		if field.IsPointer {
			result = append(result, fmt.Sprintf("    if %v != nil {", field.Accessor))
			for _, code := range embeddedCases {
				result = append(result, "    "+code)
			}
			result = append(result, "    }")
		} else {
			result = append(result, embeddedCases...)
		}
	}
	return result, nil
}

func (s *Struct) generateFieldEncoding(fields []*toolbox.FieldInfo) ([]string, error) {
	fieldCases := []string{}
	for i := range fields {
		var templateKey = -1
		fieldTypeInfo := s.Type(fields[i].TypeName)
		field, err := NewField(s, fields[i], fieldTypeInfo)
		if err != nil {
			return nil, err
		}
		if field.IsAnonymous {
			embedded, err := s.generateEmbeddedFieldEncoding(field, fieldTypeInfo)
			if err != nil {
				return nil, err
			}
			fieldCases = append(fieldCases, embedded...)
			continue
		}
	main:
		switch {
		case isPrimitiveString(field.Type):
			templateKey = encodeBaseType
		case isPrimitiveArrayString(field.Type):
			templateKey = encodeBaseTypeSlice
			s.generatePrimitiveArray(field)
		default:
			if fieldTypeInfo != nil {
				if !(field.IsSlice || fieldTypeInfo.IsSlice) {
					templateKey = encodeStruct
					break main
				}
				switch {
				case isPrimitiveString(fieldTypeInfo.ComponentType):
					templateKey = decodeBaseTypeSlice
					break main
				}
				templateKey = encodeStructSlice
				break main
			} else if field.IsSlice {
				templateKey = encodeStructSlice
			} else {
				// templateKey = encodeUnknown
				templateKey = encodeStruct
				// return nil, fmt.Errorf("Unknown type %s for field %s", field.Type, field.Name)
			}
		}
		if templateKey != -1 {
			decodingCase, err := expandFieldTemplate(templateKey, field)
			if err != nil {
				return nil, err
			}
			fieldCases = append(fieldCases, decodingCase)
		}

	}
	return fieldCases, nil
}

func NewStruct(info *toolbox.TypeInfo, generator *Generator) *Struct {
	return &Struct{
		TypeInfo:  info,
		Generator: generator,
		Alias:     extractReceiverAlias(info.Name),
	}
}
