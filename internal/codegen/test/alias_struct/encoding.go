// Code generated by encodegen. DO NOT EDIT.
package alias_struct

import (
	importedtype "go.sia.tech/encodegen/internal/codegen/test/importedtype"
	encodegen "go.sia.tech/encodegen/pkg/encodegen"
)

// MarshalBuffer implements MarshalerBuffer
func (s *AliasByteSlice) MarshalBuffer(b *encodegen.ObjBuffer) {
	if s != nil {

		b.WriteUint64(uint64(len(*s)))
		b.Write([]byte(*s))

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (s *AliasByteSlice) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if s != nil {
		var length int = 0
		_ = length

		length = int(b.ReadUint64())
		if length > 0 {
			if len(*s) < length {
				*s = make([]byte, length)
			}
			(*s) = (*s)[:length]
			b.Read(*s)
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (a *AliasFixedByteArray) MarshalBuffer(b *encodegen.ObjBuffer) {
	if a != nil {

		temp := [40]byte(*a)
		b.Write([]byte(temp[:]))

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (a *AliasFixedByteArray) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if a != nil {
		var length int = 0
		_ = length

		temp := [40]byte(*a)
		b.Read(temp[:])
		*a = temp

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (a *AliasFixedImportedTypeArray) MarshalBuffer(b *encodegen.ObjBuffer) {
	if a != nil {

		temp := [3]importedtype.Imported(*a)
		for i1 := range temp {
			(*importedtype.Imported)(&temp[i1]).MarshalBuffer(b)
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (a *AliasFixedImportedTypeArray) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if a != nil {
		var length int = 0
		_ = length

		for i1 := range *a {
			(*importedtype.Imported)(&(*a)[i1]).UnmarshalBuffer(b)
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (a *AliasFixedImportedTypePointerArray) MarshalBuffer(b *encodegen.ObjBuffer) {
	if a != nil {

		temp := [3]*importedtype.Imported(*a)
		for i1 := range temp {
			if temp[i1] != nil {
				b.WriteBool(true)
				(*importedtype.Imported)(temp[i1]).MarshalBuffer(b)
			} else {
				b.WriteBool(false)
			}
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (a *AliasFixedImportedTypePointerArray) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if a != nil {
		var length int = 0
		_ = length

		for i1 := range *a {
			if b.ReadBool() {
				if (*a)[i1] == nil {
					(*a)[i1] = new(importedtype.Imported)
				}
				(*importedtype.Imported)((*a)[i1]).UnmarshalBuffer(b)
			}
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (a *AliasFixedPointerArray) MarshalBuffer(b *encodegen.ObjBuffer) {
	if a != nil {

		temp := [3]*uint16(*a)
		for i1 := range temp {
			if temp[i1] != nil {
				b.WriteBool(true)
				b.WriteUint64(uint64(*temp[i1]))
			} else {
				b.WriteBool(false)
			}
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (a *AliasFixedPointerArray) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if a != nil {
		var length int = 0
		_ = length

		for i1 := range *a {
			if b.ReadBool() {
				if (*a)[i1] == nil {
					(*a)[i1] = new(uint16)
				}
				*(*a)[i1] = uint16((b.ReadUint64()))
			}
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (a *AliasFixedSubMessageArray) MarshalBuffer(b *encodegen.ObjBuffer) {
	if a != nil {

		temp := [3]SubMessage(*a)
		for i1 := range temp {
			(*SubMessage)(&temp[i1]).MarshalBuffer(b)
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (a *AliasFixedSubMessageArray) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if a != nil {
		var length int = 0
		_ = length

		for i1 := range *a {
			(*SubMessage)(&(*a)[i1]).UnmarshalBuffer(b)
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (a *AliasFixedSubMessagePointerArray) MarshalBuffer(b *encodegen.ObjBuffer) {
	if a != nil {

		temp := [3]*SubMessage(*a)
		for i1 := range temp {
			if temp[i1] != nil {
				b.WriteBool(true)
				(*SubMessage)(temp[i1]).MarshalBuffer(b)
			} else {
				b.WriteBool(false)
			}
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (a *AliasFixedSubMessagePointerArray) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if a != nil {
		var length int = 0
		_ = length

		for i1 := range *a {
			if b.ReadBool() {
				if (*a)[i1] == nil {
					(*a)[i1] = new(SubMessage)
				}
				(*SubMessage)((*a)[i1]).UnmarshalBuffer(b)
			}
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (t *AliasImportedType) MarshalBuffer(b *encodegen.ObjBuffer) {
	if t != nil {

		(*importedtype.Imported)(t).MarshalBuffer(b)

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (t *AliasImportedType) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if t != nil {

		(*importedtype.Imported)(t).UnmarshalBuffer(b)

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (s *AliasImportedTypePointerSlice) MarshalBuffer(b *encodegen.ObjBuffer) {
	if s != nil {

		b.WriteUint64(uint64(len(*s)))
		temp := []*importedtype.Imported(*s)
		for i1 := range temp {
			if temp[i1] != nil {
				b.WriteBool(true)
				(*importedtype.Imported)(temp[i1]).MarshalBuffer(b)
			} else {
				b.WriteBool(false)
			}
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (s *AliasImportedTypePointerSlice) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if s != nil {
		var length int = 0
		_ = length

		length = int(b.ReadUint64())
		if length > 0 {
			if len(*s) < length {
				*s = make([]*importedtype.Imported, length)
			}
			(*s) = (*s)[:length]
			for i1 := range *s {
				if b.ReadBool() {
					if (*s)[i1] == nil {
						(*s)[i1] = new(importedtype.Imported)
					}
					(*importedtype.Imported)((*s)[i1]).UnmarshalBuffer(b)
				}
			}
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (s *AliasImportedTypeSlice) MarshalBuffer(b *encodegen.ObjBuffer) {
	if s != nil {

		b.WriteUint64(uint64(len(*s)))
		temp := []importedtype.Imported(*s)
		for i1 := range temp {
			(*importedtype.Imported)(&temp[i1]).MarshalBuffer(b)
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (s *AliasImportedTypeSlice) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if s != nil {
		var length int = 0
		_ = length

		length = int(b.ReadUint64())
		if length > 0 {
			if len(*s) < length {
				*s = make([]importedtype.Imported, length)
			}
			(*s) = (*s)[:length]
			for i1 := range *s {
				(*importedtype.Imported)(&(*s)[i1]).UnmarshalBuffer(b)
			}
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (i *AliasInt) MarshalBuffer(b *encodegen.ObjBuffer) {
	if i != nil {

		b.WriteUint64(uint64(int(*i)))

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (i *AliasInt) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if i != nil {

		*i = AliasInt(int((b.ReadUint64())))

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (a *AliasIntArray) MarshalBuffer(b *encodegen.ObjBuffer) {
	if a != nil {

		b.WriteUint64(uint64(len(*a)))
		temp := []int(*a)
		for i1 := range temp {
			b.WriteUint64(uint64(temp[i1]))
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (a *AliasIntArray) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if a != nil {
		var length int = 0
		_ = length

		length = int(b.ReadUint64())
		if length > 0 {
			if len(*a) < length {
				*a = make([]int, length)
			}
			(*a) = (*a)[:length]
			for i1 := range *a {
				(*a)[i1] = int((b.ReadUint64()))
			}
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (a *AliasIntPointerArray) MarshalBuffer(b *encodegen.ObjBuffer) {
	if a != nil {

		b.WriteUint64(uint64(len(*a)))
		temp := []*int(*a)
		for i1 := range temp {
			if temp[i1] != nil {
				b.WriteBool(true)
				b.WriteUint64(uint64(*temp[i1]))
			} else {
				b.WriteBool(false)
			}
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (a *AliasIntPointerArray) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if a != nil {
		var length int = 0
		_ = length

		length = int(b.ReadUint64())
		if length > 0 {
			if len(*a) < length {
				*a = make([]*int, length)
			}
			(*a) = (*a)[:length]
			for i1 := range *a {
				if b.ReadBool() {
					if (*a)[i1] == nil {
						(*a)[i1] = new(int)
					}
					*(*a)[i1] = int((b.ReadUint64()))
				}
			}
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (m *AliasSubMessage) MarshalBuffer(b *encodegen.ObjBuffer) {
	if m != nil {

		(*SubMessage)(m).MarshalBuffer(b)

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (m *AliasSubMessage) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if m != nil {

		(*SubMessage)(m).UnmarshalBuffer(b)

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (a *AliasSubMessageArray) MarshalBuffer(b *encodegen.ObjBuffer) {
	if a != nil {

		b.WriteUint64(uint64(len(*a)))
		temp := []SubMessage(*a)
		for i1 := range temp {
			(*SubMessage)(&temp[i1]).MarshalBuffer(b)
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (a *AliasSubMessageArray) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if a != nil {
		var length int = 0
		_ = length

		length = int(b.ReadUint64())
		if length > 0 {
			if len(*a) < length {
				*a = make([]SubMessage, length)
			}
			(*a) = (*a)[:length]
			for i1 := range *a {
				(*SubMessage)(&(*a)[i1]).UnmarshalBuffer(b)
			}
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (a *AliasSubMessagePointerArray) MarshalBuffer(b *encodegen.ObjBuffer) {
	if a != nil {

		b.WriteUint64(uint64(len(*a)))
		temp := []*SubMessage(*a)
		for i1 := range temp {
			if temp[i1] != nil {
				b.WriteBool(true)
				(*SubMessage)(temp[i1]).MarshalBuffer(b)
			} else {
				b.WriteBool(false)
			}
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (a *AliasSubMessagePointerArray) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if a != nil {
		var length int = 0
		_ = length

		length = int(b.ReadUint64())
		if length > 0 {
			if len(*a) < length {
				*a = make([]*SubMessage, length)
			}
			(*a) = (*a)[:length]
			for i1 := range *a {
				if b.ReadBool() {
					if (*a)[i1] == nil {
						(*a)[i1] = new(SubMessage)
					}
					(*SubMessage)((*a)[i1]).UnmarshalBuffer(b)
				}
			}
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (t *DoubleAliasImportedType) MarshalBuffer(b *encodegen.ObjBuffer) {
	if t != nil {

		(*AliasImportedType)(t).MarshalBuffer(b)

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (t *DoubleAliasImportedType) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if t != nil {

		(*AliasImportedType)(t).UnmarshalBuffer(b)

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (i *DoubleAliasInt) MarshalBuffer(b *encodegen.ObjBuffer) {
	if i != nil {

		(*AliasInt)(i).MarshalBuffer(b)

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (i *DoubleAliasInt) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if i != nil {

		(*AliasInt)(i).UnmarshalBuffer(b)

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (m *DoubleAliasSubMessage) MarshalBuffer(b *encodegen.ObjBuffer) {
	if m != nil {

		(*AliasSubMessage)(m).MarshalBuffer(b)

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (m *DoubleAliasSubMessage) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if m != nil {

		(*AliasSubMessage)(m).UnmarshalBuffer(b)

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (i *Integer) MarshalBuffer(b *encodegen.ObjBuffer) {
	if i != nil {

		b.WriteUint64(uint64(len(*i)))
		temp := []*int(*i)
		for i1 := range temp {
			if temp[i1] != nil {
				b.WriteBool(true)
				b.WriteUint64(uint64(*temp[i1]))
			} else {
				b.WriteBool(false)
			}
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (i *Integer) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if i != nil {
		var length int = 0
		_ = length

		length = int(b.ReadUint64())
		if length > 0 {
			if len(*i) < length {
				*i = make([]*int, length)
			}
			(*i) = (*i)[:length]
			for i1 := range *i {
				if b.ReadBool() {
					if (*i)[i1] == nil {
						(*i)[i1] = new(int)
					}
					*(*i)[i1] = int((b.ReadUint64()))
				}
			}
		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (m *Message) MarshalBuffer(b *encodegen.ObjBuffer) {
	if m != nil {

		b.WriteUint64(uint64(m.Id))

		(*SubMessage)(&m.Sub).MarshalBuffer(b)

		(*AliasSubMessage)(&m.AliasSubMessageField).MarshalBuffer(b)

		b.WriteUint64(uint64(len(m.ArrayAliasSubMessageField)))
		for i := range m.ArrayAliasSubMessageField {
			m.ArrayAliasSubMessageField[i].MarshalBuffer(b)
		}

		(*DoubleAliasSubMessage)(&m.DoubleAliasSubMessageField).MarshalBuffer(b)

		if m.PointerDoubleAliasSubMessageField != nil {
			b.WriteBool(true)
			(*DoubleAliasSubMessage)(m.PointerDoubleAliasSubMessageField).MarshalBuffer(b)
		} else {
			b.WriteBool(false)
		}

		(*AliasInt)(&m.AliasIntField).MarshalBuffer(b)

		if m.PointerAliasIntField != nil {
			b.WriteBool(true)
			(*AliasInt)(m.PointerAliasIntField).MarshalBuffer(b)
		} else {
			b.WriteBool(false)
		}

		(*AliasIntArray)(&m.AliasIntArrayField).MarshalBuffer(b)

		(*AliasIntPointerArray)(&m.AliasIntPointerArrayField).MarshalBuffer(b)

		(*AliasSubMessageArray)(&m.AliasSubMessageArrayField).MarshalBuffer(b)

		(*AliasSubMessagePointerArray)(&m.AliasSubMessagePointerArrayField).MarshalBuffer(b)

		b.WriteUint64(uint64(len(m.ArrayAliasSubMessagePointerArrayField)))
		for i := range m.ArrayAliasSubMessagePointerArrayField {
			m.ArrayAliasSubMessagePointerArrayField[i].MarshalBuffer(b)
		}

		(*DoubleAliasInt)(&m.DoubleAliasIntField).MarshalBuffer(b)

		b.WriteUint64(uint64(len(m.ByteSlice)))
		b.Write(m.ByteSlice)

		(*AliasByteSlice)(&m.AliasByteSliceField).MarshalBuffer(b)

		(*AliasFixedByteArray)(&m.AliasFixedByteArrayField).MarshalBuffer(b)

		(*AliasFixedPointerArray)(&m.AliasFixedPointerArrayField).MarshalBuffer(b)

		(*AliasFixedSubMessageArray)(&m.AliasFixedSubMessageArrayField).MarshalBuffer(b)

		(*AliasFixedSubMessagePointerArray)(&m.AliasFixedSubMessagePointerArrayField).MarshalBuffer(b)

		for i := range m.AliasFixedByteArrayArrayField {
			m.AliasFixedByteArrayArrayField[i].MarshalBuffer(b)
		}

		for i := range m.AliasFixedPointerArrayArrayField {
			m.AliasFixedPointerArrayArrayField[i].MarshalBuffer(b)
		}

		for i := range m.AliasFixedSubMessageArrayArrayField {
			m.AliasFixedSubMessageArrayArrayField[i].MarshalBuffer(b)
		}

		for i := range m.AliasFixedSubMessagePointerArrayArrayField {
			m.AliasFixedSubMessagePointerArrayArrayField[i].MarshalBuffer(b)
		}

		(*Integer)(&m.IntegerField).MarshalBuffer(b)

		(*AliasImportedType)(&m.AliasImportedTypeField).MarshalBuffer(b)

		(*DoubleAliasImportedType)(&m.DoubleAliasImportedTypeField).MarshalBuffer(b)

		if m.PointerAliasImportedTypeField != nil {
			b.WriteBool(true)
			(*AliasImportedType)(m.PointerAliasImportedTypeField).MarshalBuffer(b)
		} else {
			b.WriteBool(false)
		}

		(*AliasImportedTypeSlice)(&m.AliasImportedTypeSliceField).MarshalBuffer(b)

		(*AliasFixedImportedTypeArray)(&m.AliasFixedImportedTypeArrayField).MarshalBuffer(b)

		(*AliasImportedTypePointerSlice)(&m.AliasImportedTypePointerSliceField).MarshalBuffer(b)

		(*AliasFixedImportedTypePointerArray)(&m.AliasFixedImportedTypePointerArrayField).MarshalBuffer(b)

		(*importedtype.Hash)(&m.Hash).MarshalBuffer(b)

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (m *Message) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if m != nil {
		var length int = 0
		_ = length

		m.Id = int(b.ReadUint64())

		(*SubMessage)(&m.Sub).UnmarshalBuffer(b)

		(*AliasSubMessage)(&m.AliasSubMessageField).UnmarshalBuffer(b)

		length = int(b.ReadUint64())
		if length > 0 {
			if len(m.ArrayAliasSubMessageField) < length {
				m.ArrayAliasSubMessageField = make([]AliasSubMessage, length)
			}
			m.ArrayAliasSubMessageField = m.ArrayAliasSubMessageField[:length]
			for i := range m.ArrayAliasSubMessageField {
				(*AliasSubMessage)(&m.ArrayAliasSubMessageField[i]).UnmarshalBuffer(b)
			}
		}

		(*DoubleAliasSubMessage)(&m.DoubleAliasSubMessageField).UnmarshalBuffer(b)

		if b.ReadBool() {
			if m.PointerDoubleAliasSubMessageField == nil {
				m.PointerDoubleAliasSubMessageField = new(DoubleAliasSubMessage)
			}
			(*DoubleAliasSubMessage)(m.PointerDoubleAliasSubMessageField).UnmarshalBuffer(b)
		}

		(*AliasInt)(&m.AliasIntField).UnmarshalBuffer(b)

		if b.ReadBool() {
			if m.PointerAliasIntField == nil {
				m.PointerAliasIntField = new(AliasInt)
			}
			(*AliasInt)(m.PointerAliasIntField).UnmarshalBuffer(b)
		}

		(*AliasIntArray)(&m.AliasIntArrayField).UnmarshalBuffer(b)

		(*AliasIntPointerArray)(&m.AliasIntPointerArrayField).UnmarshalBuffer(b)

		(*AliasSubMessageArray)(&m.AliasSubMessageArrayField).UnmarshalBuffer(b)

		(*AliasSubMessagePointerArray)(&m.AliasSubMessagePointerArrayField).UnmarshalBuffer(b)

		length = int(b.ReadUint64())
		if length > 0 {
			if len(m.ArrayAliasSubMessagePointerArrayField) < length {
				m.ArrayAliasSubMessagePointerArrayField = make([]AliasSubMessagePointerArray, length)
			}
			m.ArrayAliasSubMessagePointerArrayField = m.ArrayAliasSubMessagePointerArrayField[:length]
			for i := range m.ArrayAliasSubMessagePointerArrayField {
				(*AliasSubMessagePointerArray)(&m.ArrayAliasSubMessagePointerArrayField[i]).UnmarshalBuffer(b)
			}
		}

		(*DoubleAliasInt)(&m.DoubleAliasIntField).UnmarshalBuffer(b)

		length = int(b.ReadUint64())
		if length > 0 {
			if len(m.ByteSlice) < length {
				m.ByteSlice = make([]byte, length)
			}
			m.ByteSlice = m.ByteSlice[:length]
			b.Read(m.ByteSlice)
		}

		(*AliasByteSlice)(&m.AliasByteSliceField).UnmarshalBuffer(b)

		(*AliasFixedByteArray)(&m.AliasFixedByteArrayField).UnmarshalBuffer(b)

		(*AliasFixedPointerArray)(&m.AliasFixedPointerArrayField).UnmarshalBuffer(b)

		(*AliasFixedSubMessageArray)(&m.AliasFixedSubMessageArrayField).UnmarshalBuffer(b)

		(*AliasFixedSubMessagePointerArray)(&m.AliasFixedSubMessagePointerArrayField).UnmarshalBuffer(b)

		for i := range m.AliasFixedByteArrayArrayField {
			(*AliasFixedByteArray)(&m.AliasFixedByteArrayArrayField[i]).UnmarshalBuffer(b)
		}

		for i := range m.AliasFixedPointerArrayArrayField {
			(*AliasFixedPointerArray)(&m.AliasFixedPointerArrayArrayField[i]).UnmarshalBuffer(b)
		}

		for i := range m.AliasFixedSubMessageArrayArrayField {
			(*AliasFixedSubMessageArray)(&m.AliasFixedSubMessageArrayArrayField[i]).UnmarshalBuffer(b)
		}

		for i := range m.AliasFixedSubMessagePointerArrayArrayField {
			(*AliasFixedSubMessagePointerArray)(&m.AliasFixedSubMessagePointerArrayArrayField[i]).UnmarshalBuffer(b)
		}

		(*Integer)(&m.IntegerField).UnmarshalBuffer(b)

		(*AliasImportedType)(&m.AliasImportedTypeField).UnmarshalBuffer(b)

		(*DoubleAliasImportedType)(&m.DoubleAliasImportedTypeField).UnmarshalBuffer(b)

		if b.ReadBool() {
			if m.PointerAliasImportedTypeField == nil {
				m.PointerAliasImportedTypeField = new(AliasImportedType)
			}
			(*AliasImportedType)(m.PointerAliasImportedTypeField).UnmarshalBuffer(b)
		}

		(*AliasImportedTypeSlice)(&m.AliasImportedTypeSliceField).UnmarshalBuffer(b)

		(*AliasFixedImportedTypeArray)(&m.AliasFixedImportedTypeArrayField).UnmarshalBuffer(b)

		(*AliasImportedTypePointerSlice)(&m.AliasImportedTypePointerSliceField).UnmarshalBuffer(b)

		(*AliasFixedImportedTypePointerArray)(&m.AliasFixedImportedTypePointerArrayField).UnmarshalBuffer(b)

		(*importedtype.Hash)(&m.Hash).UnmarshalBuffer(b)

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (m *SubMessage) MarshalBuffer(b *encodegen.ObjBuffer) {
	if m != nil {

		b.WriteUint64(uint64(m.Id))

		b.WriteString((m.Description))

		b.WriteUint64(uint64(len(m.Strings)))
		for i := range m.Strings {
			b.WriteString((m.Strings[i]))
		}

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (m *SubMessage) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if m != nil {
		var length int = 0
		_ = length

		m.Id = int(b.ReadUint64())

		m.Description = string(b.ReadString())

		length = int(b.ReadUint64())
		if length > 0 {
			if len(m.Strings) < length {
				m.Strings = make([]string, length)
			}
			m.Strings = m.Strings[:length]
			for i := range m.Strings {
				m.Strings[i] = string((b.ReadString()))
			}
		}

	}
	return b.Err()
}
