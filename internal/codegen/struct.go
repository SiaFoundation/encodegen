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
	var resetCode = ""
	if s.options.PoolObjects {
		resetCode, err = s.generateReset(structInfo.Fields())
		if err != nil {
			return "", err
		}
	}
	var data = struct {
		Receiver      string
		Alias         string
		InitEmbedded  string
		EncodingCases string
		DecodingCases string
		Reset         string
		FieldCount    int
	}{
		Receiver:      s.Alias + " *" + s.Name,
		DecodingCases: strings.Join(decodingCases, "\n"),
		EncodingCases: strings.Join(encodingCases, "\n"),
		FieldCount:    len(decodingCases),
		InitEmbedded:  initEmbedded,
		Reset:         resetCode,
		Alias:         s.Alias,
	}
	return expandBlockTemplate(encodingStructType, data)
}

func (s *Struct) generateReset(fields []*toolbox.FieldInfo) (string, error) {
	fieldReset, err := s.generateFieldReset(fields)
	if err != nil {
		return "", nil
	}
	return expandBlockTemplate(resetStruct, struct {
		Reset    string
		Receiver string
	}{
		Reset:    strings.Join(fieldReset, "\n"),
		Receiver: s.Alias + " *" + s.Name,
	})
}

func (s *Struct) generateFieldReset(fields []*toolbox.FieldInfo) ([]string, error) {
	fieldReset := []string{}
	for i := range fields {
		var templateKey = -1
		fieldTypeInfo := s.Type(normalizeTypeName(fields[i].TypeName))
		field, err := NewField(s, fields[i], fieldTypeInfo)
		if err != nil {
			return nil, err
		}
		if field.IsPointer || field.IsSlice || (fieldTypeInfo != nil && fieldTypeInfo.IsSlice) {
			templateKey = resetFieldValue
		} else {
			if isPrimitiveString(field.Type) || isPrimitiveArrayString(field.Type) {
				templateKey = resetFieldValue
			}
		}
		if templateKey != -1 {
			code, err := expandFieldTemplate(templateKey, field)
			if err != nil {
				return nil, err
			}
			fieldReset = append(fieldReset, code)
		}
	}
	return fieldReset, nil
}

func (s *Struct) generateFieldDecoding(fields []*toolbox.FieldInfo) (string, []string, error) {

	fieldCases := []string{}
	var initCode = ""
	for i := range fields {
		if isSkipable(s.options, fields[i]) {
			continue
		}
		var templateKey = -1
		fieldTypeInfo := s.Type(normalizeTypeName(fields[i].TypeName))
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
				if f, _, ok := s.typedFieldDecode(field, field.ComponentType); ok {
					templateKey = decodeStructSlice
					if err = f(field); err != nil {
						return "", nil, err
					}
				} else {
					templateKey = decodeStructSlice
					if err = s.generateObjectArray(field); err != nil {
						return "", nil, err
					}
				}
			} else if _, k, ok := s.typedFieldDecode(field, field.Type); ok {
				templateKey = k
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
		if isSkipable(s.options, fields[i]) {
			continue
		}
		var templateKey = -1
		fieldTypeInfo := s.Type(normalizeTypeName(fields[i].TypeName))
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
			} else if _, k, ok := s.typedFieldEncode(field, field.Type); ok {
				templateKey = k
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

func (s *Struct) typedFieldEncode(field *Field, typeName string) (func(*Field) error, int, bool) {
	return nil, 0, false
}

func (s *Struct) typedFieldDecode(field *Field, typeName string) (func(*Field) error, int, bool) {
	return nil, 0, false
}

func NewStruct(info *toolbox.TypeInfo, generator *Generator) *Struct {
	return &Struct{
		TypeInfo:  info,
		Generator: generator,
		Alias:     extractReceiverAlias(info.Name),
	}
}
