// Code generated by encodegen. DO NOT EDIT.
package embedded_struct

import (
	"go.sia.tech/encodegen/pkg/encodegen"
)

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
func (m *Message) MarshalBuffer(b *encodegen.ObjBuffer) {
	if m != nil {

		b.WriteUint64(uint64(m.Id))

		b.WriteUint64(uint64(m.Anonymous.IntegerField))

		b.WritePrefixedBytes(encodegen.StringToBytes(m.Anonymous.StringField))

		b.WriteUint64(uint64(len(m.Anonymous.IntegerSliceField)))
		for i1 := range m.Anonymous.IntegerSliceField {
			b.WriteUint64(uint64(m.Anonymous.IntegerSliceField[i1]))
		}

		(*SubMessage)(&m.Anonymous.Sub).MarshalBuffer(b)

		(*AliasSubMessage)(&m.Anonymous.AliasSub).MarshalBuffer(b)

		b.WriteUint64(uint64(m.Anonymous.AnonymousSub1.A))

		b.WriteUint64(uint64(m.Anonymous.AnonymousSub1.AnonymousSub2.A))

		if m.Anonymous2 != nil {
			b.WriteBool(true)

			b.WriteUint64(uint64(m.Anonymous2.IntegerField))

			if m.Anonymous2.Anonymous2Sub1 != nil {
				b.WriteBool(true)

				b.WriteUint64(uint64(m.Anonymous2.Anonymous2Sub1.A))

			} else {
				b.WriteBool(false)
			}

		} else {
			b.WriteBool(false)
		}

		b.WriteUint64(uint64(len(m.Anonymous4)))
		for i := range m.Anonymous4 {

			b.WriteUint64(uint64(m.Anonymous4[i].A))

			b.WriteUint64(uint64(len(m.Anonymous4[i].IntegerSliceField)))
			for i1 := range m.Anonymous4[i].IntegerSliceField {
				b.WriteUint64(uint64(m.Anonymous4[i].IntegerSliceField[i1]))
			}

			b.WriteUint64(uint64(len(m.Anonymous4[i].Anonymous5)))
			for i1 := range m.Anonymous4[i].Anonymous5 {

				b.WriteUint64(uint64(len(m.Anonymous4[i].Anonymous5[i1].A)))
				for i2 := range m.Anonymous4[i].Anonymous5[i1].A {
					b.WriteUint64(uint64(m.Anonymous4[i].Anonymous5[i1].A[i2]))
				}

				if m.Anonymous4[i].Anonymous5[i1].B != nil {
					b.WriteBool(true)

					b.WriteUint64(uint64(m.Anonymous4[i].Anonymous5[i1].B.A))

				} else {
					b.WriteBool(false)
				}

			}

		}

		b.WriteUint64(uint64(len(m.Anonymous5)))
		for i := range m.Anonymous5 {

			b.WriteUint64(uint64(m.Anonymous5[i].A))

			b.WriteUint64(uint64(len(m.Anonymous5[i].B)))
			for i1 := range m.Anonymous5[i].B {
				m.Anonymous5[i].B[i1].MarshalBuffer(b)
			}

		}

		b.WriteUint64(uint64(len(m.Anonymous6)))
		for i := range m.Anonymous6 {
			if m.Anonymous6[i] != nil {
				b.WriteBool(true)

				b.WriteUint64(uint64(m.Anonymous6[i].A))

			} else {
				b.WriteBool(false)
			}
		}

		b.WriteUint64(uint64(m.End))

	}
}

// UnmarshalBuffer implements encodegen's UnmarshalerBuffer
func (m *Message) UnmarshalBuffer(b *encodegen.ObjBuffer) error {
	if m != nil {
		var length int = 0

		m.Id = int(b.ReadUint64())

		m.Anonymous.IntegerField = int(b.ReadUint64())

		m.Anonymous.StringField = string(b.ReadPrefixedBytes())

		length = int(b.ReadUint64())
		if length > 0 {
			if len(m.Anonymous.IntegerSliceField) < length {
				m.Anonymous.IntegerSliceField = make([]int, length)
			}
			for i1 := range m.Anonymous.IntegerSliceField {
				if i1 == length {
					break
				}
				m.Anonymous.IntegerSliceField[i1] = int((b.ReadUint64()))
			}
		}

		(*SubMessage)(&m.Anonymous.Sub).UnmarshalBuffer(b)

		(*AliasSubMessage)(&m.Anonymous.AliasSub).UnmarshalBuffer(b)

		m.Anonymous.AnonymousSub1.A = int(b.ReadUint64())

		m.Anonymous.AnonymousSub1.AnonymousSub2.A = int(b.ReadUint64())

		if b.ReadBool() {
			if m.Anonymous2 == nil {
				m.Anonymous2 = new(struct {
					IntegerField   int
					Anonymous2Sub1 *struct{ A int }
				})
			}

			m.Anonymous2.IntegerField = int(b.ReadUint64())

			if b.ReadBool() {
				if m.Anonymous2.Anonymous2Sub1 == nil {
					m.Anonymous2.Anonymous2Sub1 = new(struct{ A int })
				}

				m.Anonymous2.Anonymous2Sub1.A = int(b.ReadUint64())

			}

		}

		length = int(b.ReadUint64())
		if length > 0 {
			if len(m.Anonymous4) < length {
				m.Anonymous4 = make([]struct {
					A                 int
					IntegerSliceField []int
					Anonymous5        []struct {
						A []int
						B *struct{ A int }
					}
				}, length)
			}
			for i := range m.Anonymous4 {
				if i == length {
					break
				}

				m.Anonymous4[i].A = int(b.ReadUint64())

				length = int(b.ReadUint64())
				if length > 0 {
					if len(m.Anonymous4[i].IntegerSliceField) < length {
						m.Anonymous4[i].IntegerSliceField = make([]int, length)
					}
					for i1 := range m.Anonymous4[i].IntegerSliceField {
						if i1 == length {
							break
						}
						m.Anonymous4[i].IntegerSliceField[i1] = int((b.ReadUint64()))
					}
				}

				length = int(b.ReadUint64())
				if length > 0 {
					if len(m.Anonymous4[i].Anonymous5) < length {
						m.Anonymous4[i].Anonymous5 = make([]struct {
							A []int
							B *struct{ A int }
						}, length)
					}
					for i1 := range m.Anonymous4[i].Anonymous5 {
						if i1 == length {
							break
						}

						length = int(b.ReadUint64())
						if length > 0 {
							if len(m.Anonymous4[i].Anonymous5[i1].A) < length {
								m.Anonymous4[i].Anonymous5[i1].A = make([]int, length)
							}
							for i2 := range m.Anonymous4[i].Anonymous5[i1].A {
								if i2 == length {
									break
								}
								m.Anonymous4[i].Anonymous5[i1].A[i2] = int((b.ReadUint64()))
							}
						}

						if b.ReadBool() {
							if m.Anonymous4[i].Anonymous5[i1].B == nil {
								m.Anonymous4[i].Anonymous5[i1].B = new(struct{ A int })
							}

							m.Anonymous4[i].Anonymous5[i1].B.A = int(b.ReadUint64())

						}

					}
				}

			}
		}

		length = int(b.ReadUint64())
		if length > 0 {
			if len(m.Anonymous5) < length {
				m.Anonymous5 = make([]struct {
					A int
					B []AliasSubMessage
				}, length)
			}
			for i := range m.Anonymous5 {
				if i == length {
					break
				}

				m.Anonymous5[i].A = int(b.ReadUint64())

				length = int(b.ReadUint64())
				if length > 0 {
					if len(m.Anonymous5[i].B) < length {
						m.Anonymous5[i].B = make([]AliasSubMessage, length)
					}
					for i1 := range m.Anonymous5[i].B {
						if i1 == length {
							break
						}
						(*AliasSubMessage)(&m.Anonymous5[i].B[i1]).UnmarshalBuffer(b)
					}
				}

			}
		}

		length = int(b.ReadUint64())
		if length > 0 {
			if len(m.Anonymous6) < length {
				m.Anonymous6 = make([]*struct{ A int }, length)
			}
			for i := range m.Anonymous6 {
				if i == length {
					break
				}
				if b.ReadBool() {
					if m.Anonymous6[i] == nil {
						m.Anonymous6[i] = new(struct{ A int })
					}

					m.Anonymous6[i].A = int(b.ReadUint64())

				}
			}
		}

		m.End = int(b.ReadUint64())

	}
	return b.Err()
}

// MarshalBuffer implements MarshalerBuffer
func (m *SubMessage) MarshalBuffer(b *encodegen.ObjBuffer) {
	if m != nil {

		b.WriteUint64(uint64(m.Id))

		b.WritePrefixedBytes(encodegen.StringToBytes(m.Description))

		b.WriteUint64(uint64(len(m.Strings)))
		for i := range m.Strings {
			b.WritePrefixedBytes(encodegen.StringToBytes(m.Strings[i]))
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
			if len(m.Strings) < length {
				m.Strings = make([]string, length)
			}
			for i := range m.Strings {
				if i == length {
					break
				}
				m.Strings[i] = string(encodegen.BytesToString(b.ReadPrefixedBytes()))
			}
		}

	}
	return b.Err()
}
