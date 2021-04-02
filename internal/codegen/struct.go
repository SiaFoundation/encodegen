package codegen

import (
	"go.sia.tech/encodegen/internal/toolbox"
	"strings"
)

type Struct struct {
	*toolbox.TypeInfo
	*Generator
	Alias               string
	ReuseMemory         bool
	EmptyMarshalledSize int
}

func NewStruct(info *toolbox.TypeInfo, generator *Generator, reuseMemory bool) *Struct {
	return &Struct{
		TypeInfo:    info,
		Generator:   generator,
		Alias:       extractReceiverAlias(info.Name),
		ReuseMemory: reuseMemory,
	}
}

//Generate generates decoderCode + encoderCode
func (s *Struct) Generate() (string, error) {
	return s.generateEncoding(s.TypeInfo)
}

func (s *Struct) generateEncoding(structInfo *toolbox.TypeInfo) (string, error) {
	hasSlice := fieldsHaveSlice(structInfo.Fields()) || structInfo.IsSlice
	decodingCases, encodingCases, err := s.generateFieldMethods(structInfo.Fields(), "", "")

	// if we have an alias type
	if structInfo.IsDerived {
		decodingCases, encodingCases, err = s.generateAliasCases(structInfo)
		if err != nil {
			return "", err
		}
	} else if structInfo.ComponentType != "" {
		// alias of a slice
		decodingCases, encodingCases, err = s.generateAliasCases(structInfo)
		if err != nil {
			return "", err
		}
		hasSlice = true
	}

	data := struct {
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

func (s *Struct) generateFieldMethods(fields []*toolbox.FieldInfo, anonymousPrefix string, currentIterator string) ([]string, []string, error) {
	var decodeTemplateKey int
	var encodeTemplateKey int
	var decodingCases []string
	var encodingCases []string
	encodingCases = nil

	for _, field := range fields {
		decodeTemplateKey = -1
		encodeTemplateKey = -1

		if field == nil {
			continue
		}

		// dont modify the original
		fieldCopy := *field
		fieldTypeInfo := s.Type(fieldCopy.TypeName)
		if len(anonymousPrefix) > 0 {
			fieldCopy.Name = anonymousPrefix + "." + fieldCopy.Name
		}

		field, err := NewField(s, &fieldCopy, fieldTypeInfo)
		if err != nil {
			return nil, nil, err
		}
		field.ReuseMemory = s.ReuseMemory

		field.Iterator = getNextIterator(currentIterator)

		// if we are working with a struct type, the code for it is generated too so we can call its MarshalBuffer/UnmarshalBuffer functions
		if fieldTypeInfo != nil {
			err = s.generateStructCode(Type{Name: fieldTypeInfo.Name, ReuseMemory: s.ReuseMemory})
			if err != nil {
				return nil, nil, err
			}
		}

		// if we have an anonymous struct
		if len(field.AnonymousChildFields) > 0 {
			oldPrefix := anonymousPrefix
			anonymousPrefix = fieldCopy.Name
			newDecodingCases, newEncodingCases, err := s.generateAnonymousStructCases(field, anonymousPrefix, currentIterator)
			if err != nil {
				return nil, nil, err
			}
			decodingCases = append(decodingCases, newDecodingCases...)
			encodingCases = append(encodingCases, newEncodingCases...)
			anonymousPrefix = oldPrefix
			continue
		} else if strings.Contains(field.Type, "struct{") {
			// if we have an anonymous struct at this point it means it has no fields and thus there is nothing else for us to do
			continue
		}

	main:
		switch {
		case isPrimitiveString(field.Type):
			decodeTemplateKey = decodeBaseType
			encodeTemplateKey = encodeBaseType
		case field.IsSlice && isPrimitiveString(field.ComponentType):
			decodeTemplateKey = decodeBaseTypeSlice
			encodeTemplateKey = encodeBaseTypeSlice
			// see comments regarding treatment of uint8s in base_types.go
			if field.ComponentType == "uint8" {
				field.PrimitiveFunctions = UInt8SlicePrimitiveFunctions
			}
		default:
			if fieldTypeInfo != nil {
				if !(field.IsSlice || fieldTypeInfo.IsSlice) {
					decodeTemplateKey = decodeStruct
					encodeTemplateKey = encodeStruct
					break main
				}

				// if we are working with a struct type (in this case a pointer or slice to it), the code for it is generated too so we can call its MarshalBuffer/UnmarshalBuffer functions
				err := s.generateStructCode(Type{Name: field.ComponentType, ReuseMemory: s.ReuseMemory})
				if err != nil {
					return nil, nil, err
				}

				if field.IsSlice {
					decodeTemplateKey = decodeStructSlice
					encodeTemplateKey = encodeStructSlice
				} else {
					decodeTemplateKey = decodeStruct
					encodeTemplateKey = encodeStruct
				}

				break main
			} else {
				// imported or alias type
				if field.IsSlice {
					decodeTemplateKey = decodeStructSlice
					encodeTemplateKey = encodeStructSlice
				} else {
					decodeTemplateKey = decodeStruct
					encodeTemplateKey = encodeStruct
				}

				s.enableTypeImported(field.ComponentType)

				err := s.generateStructCode(Type{Name: field.ComponentType, ReuseMemory: s.ReuseMemory})
				if err != nil {
					return nil, nil, err
				}
			}
		}
		if decodeTemplateKey != -1 {
			decodingCase, err := expandFieldTemplate(decodeTemplateKey, field)
			if err != nil {
				return nil, nil, err
			}
			decodingCases = append(decodingCases, decodingCase)
		}
		if encodeTemplateKey != -1 {
			encodingCase, err := expandFieldTemplate(encodeTemplateKey, field)
			if err != nil {
				return nil, nil, err
			}
			encodingCases = append(encodingCases, encodingCase)
		}

	}
	return decodingCases, encodingCases, nil
}

func (s *Struct) generateAnonymousStructCases(field *Field, anonymousPrefix string, currentIterator string) ([]string, []string, error) {
	var decodeTemplateKey int = -1
	var encodeTemplateKey int = -1

	if field.IsSlice {
		anonymousPrefix = anonymousPrefix + "[" + field.Iterator + "]" // a.b.c -> a.b.c[i]
		decodeTemplateKey = decodeAnonymousStructSlice
		encodeTemplateKey = encodeAnonymousStructSlice
	} else if field.IsPointer {
		decodeTemplateKey = decodeAnonymousStructPointer
		encodeTemplateKey = encodeAnonymousStructPointer
	}

	newDecodingCases, newEncodingCases, err := s.generateFieldMethods(field.AnonymousChildFields, anonymousPrefix, field.Iterator)
	if err != nil {
		return nil, nil, err
	}
	// if its not a pointer or a slice just return the array of cases for all the fields
	if !field.IsSlice && !field.IsPointer {
		return newDecodingCases, newEncodingCases, nil
	}
	data := struct {
		Accessor           string
		Type               string
		ComponentType      string
		Cases              string
		Iterator           string
		IsPointerComponent bool
		IsFixed            bool
		FixedSize          int
		ReuseMemory        bool
	}{
		field.Accessor,
		field.Type,
		field.ComponentType,
		strings.Join(newDecodingCases, "\n"),
		field.Iterator,
		field.IsPointerComponent,
		field.IsFixed,
		field.FixedSize,
		s.ReuseMemory,
	}
	decodingCase, err := expandFieldTemplate(decodeTemplateKey, data)
	if err != nil {
		return nil, nil, err
	}

	data.Cases = strings.Join(newEncodingCases, "\n")
	encodingCase, err := expandFieldTemplate(encodeTemplateKey, data)
	if err != nil {
		return nil, nil, err
	}

	return []string{decodingCase}, []string{encodingCase}, nil
}

func (s *Struct) generateAliasCases(structInfo *toolbox.TypeInfo) ([]string, []string, error) {
	var err error
	var decodeKey int
	var encodeKey int
	var newStructInfo = Field{
		Accessor:           s.Alias,
		Derived:            structInfo.Derived,
		Name:               structInfo.Name,
		ComponentType:      structInfo.ComponentType,
		IsPointerComponent: structInfo.IsPointerComponentType,
		ReuseMemory:        s.ReuseMemory,
		IsFixed:            structInfo.IsFixed,
		FixedSize:          structInfo.FixedSize,
		Iterator:           "i1",
	}
	if structInfo.IsPointerComponentType {
		newStructInfo.ComponentType = "*" + structInfo.ComponentType
	}

	if isPrimitiveString(structInfo.Derived) || isPrimitiveString(structInfo.ComponentType) {
		// alias to primitive type or array of primitive type
		if structInfo.IsSlice {
			newStructInfo.PrimitiveFunctions = supportedPrimitives[structInfo.ComponentType]
			decodeKey = decodeAliasBaseTypeSlice
			encodeKey = encodeAliasBaseTypeSlice
			if structInfo.ComponentType == "uint8" {
				newStructInfo.PrimitiveFunctions = UInt8SlicePrimitiveFunctions
			}
		} else {
			newStructInfo.PrimitiveFunctions = supportedPrimitives[structInfo.Derived]
			decodeKey = decodeAliasBaseType
			encodeKey = encodeAliasBaseType
		}
	} else {
		// alias of struct or struct slice
		if structInfo.IsSlice {
			decodeKey = decodeAliasStructSlice
			encodeKey = encodeAliasStructSlice
		} else {
			decodeKey = decodeAliasStruct
			encodeKey = encodeAliasStruct
		}
		err := s.generateStructCode(Type{
			Name:        structInfo.ComponentType,
			ReuseMemory: s.ReuseMemory,
		})
		if err != nil {
			return nil, nil, err
		}

		if structInfo.IsImport {
			s.enableTypeImported(structInfo.Derived)
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

func (s *Struct) enableTypeImported(selectorName string) {
	typeSplit := strings.Split(selectorName, ".")
	if len(typeSplit) > 1 {
		importData, ok := s.imports[typeSplit[0]]
		if ok {
			importData.Enabled = true
			importData.Types = append(importData.Types, typeSplit[1])
			s.imports[typeSplit[0]] = importData
		}
	}
}
