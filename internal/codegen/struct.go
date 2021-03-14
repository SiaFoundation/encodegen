package codegen

import (
	"strconv"
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

	decodingCases, err := s.generateFieldDecoding(structInfo.Fields(), "", "")
	if err != nil {
		return "", err
	}

	encodingCases, err := s.generateFieldEncoding(structInfo.Fields(), "", "")
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

func getNextIterator(currentIdentifier string) string {
	// this function allows the generated the code to iterate over slices of structs that have slices within them without having iteration identifiers conflict (i.e., there'd be multiple "range i := r.Fields"s)
	idSplit := strings.Split(currentIdentifier, "i")
	if len(idSplit) != 2 {
		return "i"
	}
	if idSplit[1] != "" {
		num, err := strconv.Atoi(idSplit[1])
		if err != nil {
			return "i"
		}
		return fmt.Sprintf("i%d", num+1)
	} else {
		return "i1"
	}
}

func (s *Struct) generateFieldDecoding(fields []*toolbox.FieldInfo, anonymousPrefix string, currentIterator string) ([]string, error) {
	fieldCases := []string{}
	for i := range fields {
		templateKey := -1

		if fields[i] == nil {
			continue
		}
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
		field.Iterator = getNextIterator(currentIterator)

		if fieldTypeInfo != nil {
			err = s.generateStructCode(fieldTypeInfo.Name)
			if err != nil {
				return nil, err
			}
		}

		if len(field.AnonymousChildFields) > 0 {
			oldPrefix := anonymousPrefix
			anonymousPrefix = fieldCopy.Name

			if field.IsPointer {
				fieldCases = append(fieldCases, fmt.Sprintf(`
				if b.ReadBool() {
					if %s == nil {
						%s = new(%s)
					}
				`, field.Accessor, field.Accessor, noPointer(field.Type)))
			}
			if field.IsSlice {
				fieldCases = append(fieldCases, fmt.Sprintf(`
				length = int(b.ReadUint64())
				if length > 0 {
					%s = make(%s, length) 
					for %s := range %s {

				`, field.Accessor, field.Type, field.Iterator, field.Accessor))
				anonymousPrefix = fieldCopy.Name + "[" + field.Iterator + "]"
				if field.IsPointerComponent {
					fieldCases = append(fieldCases, fmt.Sprintf(`
					if b.ReadBool() {
						%s[%s] = new(%s)
					`, field.Accessor, field.Iterator, field.ComponentType))
				}
			}


			anonymousCases, err := s.generateFieldDecoding(fieldCopy.AnonymousChildFields, anonymousPrefix, field.Iterator)
			if err != nil {
				return nil, err
			}
			fieldCases = append(fieldCases, anonymousCases...)

			if field.IsPointer {
				fieldCases = append(fieldCases, "}")
			}
			if field.IsSlice {
				fieldCases = append(fieldCases, "}")
				fieldCases = append(fieldCases, "}")
				if field.IsPointerComponent {
					fieldCases = append(fieldCases, `}`)

				}
			}

			anonymousPrefix = oldPrefix

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

func (s *Struct) generateFieldEncoding(fields []*toolbox.FieldInfo, anonymousPrefix string, currentIterator string) ([]string, error) {
	fieldCases := []string{}
	for i := range fields {
		templateKey := -1

		if fields[i] == nil {
			continue
		}
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
		field.Iterator = getNextIterator(currentIterator)
		// if we have an anonymous struct
		if len(field.AnonymousChildFields) > 0 {
			oldPrefix := anonymousPrefix
			anonymousPrefix = fieldCopy.Name

			if field.IsPointer {
				fieldCases = append(fieldCases, fmt.Sprintf(`
				if %s != nil {
					b.WriteBool(true)
				`, field.Accessor))
			}
			if field.IsSlice {
				fieldCases = append(fieldCases, fmt.Sprintf(`
				b.WriteUint64(uint64(len(%s)))
				for %s := range %s {
			`, field.Accessor, field.Iterator, field.Accessor))
				anonymousPrefix = fieldCopy.Name + "[" + field.Iterator + "]"

				if field.IsPointerComponent {
					fieldCases = append(fieldCases, fmt.Sprintf(`
					if %s[%s] != nil {
						b.WriteBool(true)
					`, field.Accessor, field.Iterator))
				}

			}

			anonymousCases, err := s.generateFieldEncoding(fieldCopy.AnonymousChildFields, anonymousPrefix, field.Iterator)
			if err != nil {
				return nil, err
			}
			fieldCases = append(fieldCases, anonymousCases...)



			if field.IsSlice {
				if field.IsPointerComponent {
					fieldCases = append(fieldCases, fmt.Sprintf(`
					} else {
						b.WriteBool(false)
					}`))
				}
				fieldCases = append(fieldCases, `}`)

			}
			if field.IsPointer {
				fieldCases = append(fieldCases, `
				} else {
					b.WriteBool(false)
				}`)
			}

			anonymousPrefix = oldPrefix
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
