package codegen

import (
	"fmt"
	"go.sia.tech/encodegen/internal/toolbox"
	"strings"
)

type Struct struct {
	*toolbox.TypeInfo
	referenced *toolbox.TypeInfo
	*Generator
	Alias string
	Init  string
	Body  string
}

//Generate generates decoderCode + encoderCode
func (s *Struct) Generate() (string, error) {
	return s.generateEncoding(s.TypeInfo)
}

func (s *Struct) generateEncoding(structInfo *toolbox.TypeInfo) (string, error) {
	hasSlice := fieldsHaveSlice(structInfo.Fields())

	decodingCases, err := s.generateFieldDecoding(structInfo.Fields(), "")
	if err != nil {
		return "", err
	}

	encodingCases, err := s.generateFieldEncoding(structInfo.Fields(), "")
	if err != nil {
		return "", err
	}

	if structInfo.IsDerived {
		decodingCases, encodingCases, err = s.generateAliasCases(structInfo)
		if err != nil {
			return "", err
		}
	} else if structInfo.ComponentType != "" {
		decodingCases, encodingCases, err = s.generateAliasCases(structInfo)
		if err != nil {
			return "", err
		}
		hasSlice = true
	}

	var data = struct {
		Receiver      string
		Alias         string
		EncodingCases string
		DecodingCases string
		FieldCount    int
		HasSlice      bool
	}{
		Receiver:      s.Alias + " *" + s.Name,
		DecodingCases: strings.Join(decodingCases, "\n"),
		EncodingCases: strings.Join(encodingCases, "\n"),
		FieldCount:    len(decodingCases),
		Alias:         s.Alias,
		HasSlice:      hasSlice,
	}
	return expandBlockTemplate(encodingStructType, data)
}

func fieldsHaveSlice(fields []*toolbox.FieldInfo) bool {
	for _, field := range fields {
		if field.IsSlice {
			return true
		}
		if len(field.AnonymousChildFields) > 0 {
			hasSlice := fieldsHaveSlice(field.AnonymousChildFields)
			if hasSlice {
				return hasSlice
			}
		}
	}
	return false
}

func (s *Struct) generateAliasCases(structInfo *toolbox.TypeInfo) ([]string, []string, error) {
	var err error
	var decodeKey int
	var encodeKey int
	var newStructInfo = struct {
		Accessor           string
		Derived            string
		Name               string
		DecodingMethod     string
		EncodingMethod     string
		PrimitiveWriteCast string
		ComponentType      string
		IsPointerComponent bool
	}{
		Accessor:           s.Alias,
		Derived:            structInfo.Derived,
		Name:               structInfo.Name,
		ComponentType:      structInfo.ComponentType,
		IsPointerComponent: structInfo.IsPointerComponentType,
	}
	if structInfo.IsPointerComponentType {
		newStructInfo.ComponentType = "*" + structInfo.ComponentType
	}

	if (isPrimitiveString(structInfo.Derived) || isPrimitiveArrayString(structInfo.Derived)) || (isPrimitiveString(structInfo.ComponentType) || isPrimitiveArrayString(structInfo.ComponentType)) {
		if structInfo.IsSlice {
			newStructInfo.DecodingMethod = supportedPrimitives[structInfo.ComponentType].ReadFunction
			newStructInfo.EncodingMethod = supportedPrimitives[structInfo.ComponentType].WriteFunction
			newStructInfo.PrimitiveWriteCast = supportedPrimitives[structInfo.ComponentType].WriteCast
			decodeKey = decodeAliasBaseTypeSlice
			encodeKey = encodeAliasBaseTypeSlice
		} else {
			newStructInfo.DecodingMethod = supportedPrimitives[structInfo.Derived].ReadFunction
			newStructInfo.EncodingMethod = supportedPrimitives[structInfo.Derived].WriteFunction
			newStructInfo.PrimitiveWriteCast = supportedPrimitives[structInfo.Derived].WriteCast

			decodeKey = decodeAliasBaseType
			encodeKey = encodeAliasBaseType
		}
	} else {
		if structInfo.IsSlice {
			decodeKey = decodeAliasStructSlice
			encodeKey = encodeAliasStructSlice
		} else {
			decodeKey = decodeAliasStruct
			encodeKey = encodeAliasStruct
		}
	}

	decode, err := expandFieldTemplate(decodeKey, newStructInfo)
	if err != nil {
		return nil, nil, err
	}
	encode, err := expandFieldTemplate(encodeKey, newStructInfo)
	if err != nil {
		return nil, nil, err
	}
	return []string{decode}, []string{encode}, nil
}

func (s *Struct) generateFieldDecoding(fields []*toolbox.FieldInfo, anonymousPrefix string) ([]string, error) {
	fieldCases := []string{}
	for i := range fields {
		templateKey := -1

		// dont modify the original
		fieldCopy := *fields[i]

		fieldTypeInfo := s.Type(fieldCopy.TypeName)
		if len(anonymousPrefix) > 0 {
			fieldCopy.Name = anonymousPrefix + "." + fieldCopy.Name
		}

		field, err := NewField(s, &fieldCopy, fieldTypeInfo)
		if err != nil {
			return nil, err
		}

		if fieldTypeInfo != nil {
			err = s.generateStructCode(fieldTypeInfo.Name)
			if err != nil {
				return nil, err
			}
		}

		if len(field.AnonymousChildFields) > 0 {
			oldPrefix := anonymousPrefix
			anonymousPrefix = fieldCopy.Name
			anonymousCases, err := s.generateFieldDecoding(field.AnonymousChildFields, anonymousPrefix)
			if err != nil {
				return nil, err
			}
			anonymousPrefix = oldPrefix
			if field.IsPointer {
				fieldCases = append(fieldCases, fmt.Sprintf(`
				if b.ReadBool() {
					if %s == nil {
						%s = new(%s)
					}
				`, field.Accessor, field.Accessor, noPointer(field.Type)))
			}
			fieldCases = append(fieldCases, anonymousCases...)
			if field.IsPointer {
				fieldCases = append(fieldCases, "}")
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
					// s.generatePrimitiveArray(field)
					// templateKey = decodeBaseTypeSlice
					templateKey = decodeStruct
					break main

				}

				if err = s.generateStructCode(field.ComponentType); err != nil {
					return nil, err
				}

				if field.ComponentType != "" {
					// templateKey = decodeStruct
					templateKey = decodeStructSlice
					if err = s.generateObjectArray(field); err != nil {
						return nil, err
					}
				} else {
					templateKey = decodeStruct
				}

				break main
			} else if field.IsSlice {
				templateKey = decodeStructSlice
				if err = s.generateObjectArray(field); err != nil {
					return nil, err
				}
			} else {
				// templateKey = decodeStruct
				continue
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

func (s *Struct) generateFieldEncoding(fields []*toolbox.FieldInfo, anonymousPrefix string) ([]string, error) {
	fieldCases := []string{}
	for i := range fields {
		templateKey := -1

		// dont modify the original
		fieldCopy := *fields[i]

		fieldTypeInfo := s.Type(fieldCopy.TypeName)
		if len(anonymousPrefix) > 0 {
			fieldCopy.Name = anonymousPrefix + "." + fieldCopy.Name
		}
		field, err := NewField(s, &fieldCopy, fieldTypeInfo)
		if err != nil {
			return nil, err
		}

		// if we have an anonymous struct
		if len(field.AnonymousChildFields) > 0 {
			oldPrefix := anonymousPrefix
			anonymousPrefix = fieldCopy.Name
			anonymousCases, err := s.generateFieldEncoding(fieldCopy.AnonymousChildFields, anonymousPrefix)
			if err != nil {
				return nil, err
			}
			anonymousPrefix = oldPrefix

			if field.IsPointer {
				fieldCases = append(fieldCases, fmt.Sprintf(`
				if %s != nil {
					b.WriteBool(true)
				`, field.Accessor))
			}
			fieldCases = append(fieldCases, anonymousCases...)
			if field.IsPointer {
				fieldCases = append(fieldCases, `
				} else {
					b.WriteBool(false)
				}`)
			}
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
				if isPrimitiveString(fieldTypeInfo.ComponentType) {
					templateKey = encodeStruct
					break main
				}
				if field.ComponentType != "" {
					templateKey = encodeStructSlice
				} else {
					templateKey = encodeStruct
				}
				break main
			} else if field.IsSlice {
				templateKey = encodeStructSlice
			} else {
				// anonymous struct
				// templateKey = encodeStruct
				continue
			}
		}
		if templateKey != -1 {
			encodingCase, err := expandFieldTemplate(templateKey, field)
			if err != nil {
				return nil, err
			}

			fieldCases = append(fieldCases, encodingCase)
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
