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
		decodingCases, encodingCases = s.generateAliasCases(structInfo)
		// }
	}

	var data = struct {
		Receiver string
		Alias         string
		InitEmbedded  string
		EncodingCases string
		DecodingCases string
		FieldCount    int
	}{
		Receiver: s.Alias + " *" + s.Name,
		DecodingCases: strings.Join(decodingCases, "\n"),
		EncodingCases: strings.Join(encodingCases, "\n"),
		FieldCount:    len(decodingCases),
		InitEmbedded:  initEmbedded,
		Alias:         s.Alias,
	}
	return expandBlockTemplate(encodingStructType, data)
}

/*
The template system works on struct fields not structs and but these cases are very simple (we do not need any templating logic like if statements etc) so I decided to make them separate rather than add a bunch more code.
If we have too many cases like these I will make it generic but currently I don't think it is a problem.
*/
func (s *Struct) generateAliasCases(structInfo *toolbox.TypeInfo) ([]string, []string) {
	if !isPrimitiveString(structInfo.Derived) && !isPrimitiveArrayString(structInfo.Derived) {
		return []string{fmt.Sprintf("(*%s)(%s).UnmarshalBuffer(b)", structInfo.Derived, s.Alias)}, []string{fmt.Sprintf("(*%s)(%s).MarshalBuffer(b)", structInfo.Derived, s.Alias)}
	} else {
		primitiveType := supportedPrimitives[structInfo.Derived]
		return []string{
			fmt.Sprintf("*%s = %s(%s(b.%s()))", s.Alias, structInfo.Name, structInfo.Derived, primitiveType.ReadFunction)}, []string{fmt.Sprintf("b.%s(%s(%s(*%s)))", primitiveType.WriteFunction, primitiveType.WriteCast, structInfo.Derived, s.Alias)}
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
