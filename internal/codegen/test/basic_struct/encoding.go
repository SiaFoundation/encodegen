// Code generated by encodegen. DO NOT EDIT.

package basic_struct

import (
	"go.sia.tech/encodegen/pkg/encodegen"
)

// MarshalBuffer implements MarshalerBuffer

func (m *Message) MarshalBuffer(b *encodegen.ObjBuffer) {
	if m != nil {

		b.WriteUint64(uint64(m.Id))

		b.WritePrefixedBytes([]byte(m.Name))

		b.WriteUint64(uint64(len(m.Ints)))

		for i := range m.Ints {

			b.WriteUint64(uint64(m.Ints[i]))

		}

		if m.SubMessageX != nil {
			b.WriteBool(true)
			(*SubMessage)(m.SubMessageX).MarshalBuffer(b)
		} else {
			b.WriteBool(false)
		}

		b.WriteUint64(uint64(len(m.MessagesX)))

		for i := range m.MessagesX {

			if m.MessagesX[i] != nil {
				b.WriteBool(true)

				m.MessagesX[i].MarshalBuffer(b)

			} else {
				b.WriteBool(false)
			}

		}

		(*SubMessage)(&m.SubMessageY).MarshalBuffer(b)

		b.WriteUint64(uint64(len(m.MessagesY)))

		for i := range m.MessagesY {

			m.MessagesY[i].MarshalBuffer(b)

		}

		if m.IsTrue != nil {
			b.WriteBool(true)

			b.WriteBool((*m.IsTrue))

		} else {
			b.WriteBool(false)
		}

		b.WriteUint64(uint64(len(m.Payload)))

		b.Write(m.Payload)

		b.WriteUint64(uint64(len(m.Strings)))

		for i := range m.Strings {

			b.WritePrefixedBytes([]byte(m.Strings[i]))

		}

	}

}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (m *Message) UnmarshalBuffer(b *encodegen.ObjBuffer) error {

	if m != nil {

		var length int = 0

		m.Id = int(b.ReadUint64())

		m.Name = string(b.ReadPrefixedBytes())

		length = int(b.ReadUint64())
		if length > 0 {

			m.Ints = make([]int, length)

			for i := range m.Ints {

				m.Ints[i] = int(b.ReadUint64())

			}

		}

		if b.ReadBool() {
			if m.SubMessageX == nil {
				m.SubMessageX = new(SubMessage)
			}
			(*SubMessage)(m.SubMessageX).UnmarshalBuffer(b)
		}

		length = int(b.ReadUint64())
		if length > 0 {
			m.MessagesX = make([]*SubMessage, length)

			for i := range m.MessagesX {

				if b.ReadBool() {
					m.MessagesX[i] = new(SubMessage)
					m.MessagesX[i].UnmarshalBuffer(b)
				}

			}
		}

		(*SubMessage)(&m.SubMessageY).UnmarshalBuffer(b)

		length = int(b.ReadUint64())
		if length > 0 {
			m.MessagesY = make([]SubMessage, length)

			for i := range m.MessagesY {

				(*SubMessage)(&m.MessagesY[i]).UnmarshalBuffer(b)

			}
		}

		if b.ReadBool() {
			if m.IsTrue == nil {
				m.IsTrue = new(bool)
			}
			*m.IsTrue = bool(b.ReadBool())
		}

		length = int(b.ReadUint64())
		if length > 0 {

			m.Payload = make([]byte, length)

			b.Read(m.Payload)

		}

		length = int(b.ReadUint64())
		if length > 0 {

			m.Strings = make([]string, length)

			for i := range m.Strings {

				m.Strings[i] = string(b.ReadPrefixedBytes())

			}

		}

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer

func (m *SubMessage) MarshalBuffer(b *encodegen.ObjBuffer) {
	if m != nil {

		b.WriteUint64(uint64(m.Id))

		b.WritePrefixedBytes([]byte(m.Description))

		b.WriteUint64(uint64(len(m.Strings)))

		for i := range m.Strings {

			b.WritePrefixedBytes([]byte(m.Strings[i]))

		}

	}

}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (m *SubMessage) UnmarshalBuffer(b *encodegen.ObjBuffer) error {

	if m != nil {

		var length int = 0

		m.Id = int(b.ReadUint64())

		m.Description = string(b.ReadPrefixedBytes())

		length = int(b.ReadUint64())
		if length > 0 {

			m.Strings = make([]string, length)

			for i := range m.Strings {

				m.Strings[i] = string(b.ReadPrefixedBytes())

			}

		}

	}
	return b.Err()
}
