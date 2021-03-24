package alias_struct

import (
	// "log"
	"github.com/stretchr/testify/assert"
	"gitlab.com/NebulousLabs/encoding"
	"go.sia.tech/encodegen/pkg/encodegen"
	"testing"
)

var testInt int = 999
var testByte byte = 0x99
var testAliasInt AliasInt = AliasInt(9001)
var doubleAliasSubMessageField = DoubleAliasSubMessage(AliasSubMessage(SubMessage{Id: 888, Description: "bbb"}))
var testSubMessage = SubMessage{Id: 8888, Description: "AAAAA", Strings: []string{"WASDF", "WAASDSAD"}}

var testUint16 uint16 = 5

var msg = Message{
	Id: 1022,
	// Sub:                               SubMessage{Id: 9001, Description: "AAA"},
	// AliasSubMessageField:              AliasSubMessage(SubMessage{Id: 9}),
	// ArrayAliasSubMessageField:         []AliasSubMessage{AliasSubMessage(SubMessage{Id: 8, Description: "aaaa", Strings: []string{"9"}})},
	// DoubleAliasSubMessageField:        doubleAliasSubMessageField,
	// PointerDoubleAliasSubMessageField: &doubleAliasSubMessageField,
	// AliasIntField:                     testAliasInt,
	// PointerAliasIntField:              &testAliasInt,
	// AliasIntArrayField:                AliasIntArray([]int{5, 5, 5, 5, 999995, 5, 5, 5, 9, -5}),
	// AliasIntPointerArrayField:         AliasIntPointerArray([]*int{nil, &testInt, nil, nil, &testInt}),
	// AliasSubMessageArrayField: AliasSubMessageArray([]SubMessage{
	// 	{Id: 4000, Description: "111", Strings: []string{"999", "888"}},
	// 	{Id: 8111, Description: "111", Strings: []string{}},
	// 	{Id: 8111, Description: "111", Strings: nil},
	// 	{Id: 8111, Strings: nil},
	// 	{Strings: nil},
	// }),
	// AliasSubMessagePointerArrayField: AliasSubMessagePointerArray([]*SubMessage{
	// 	&testSubMessage,
	// 	nil,
	// 	&testSubMessage,
	// 	&testSubMessage,
	// 	&testSubMessage,
	// 	nil,
	// }),
	// ArrayAliasSubMessagePointerArrayField: []AliasSubMessagePointerArray{AliasSubMessagePointerArray([]*SubMessage{
	// 	nil,
	// 	nil,
	// 	nil,
	// 	&testSubMessage,
	// 	nil,
	// 	&testSubMessage,
	// 	&testSubMessage,
	// 	&testSubMessage,
	// 	nil,
	// }), nil, nil},
	ByteSlice:                   []byte{1, 123, 123, 123, 132},
	AliasByteSliceField:         AliasByteSlice(nil),
	AliasFixedByteArrayField:    AliasFixedByteArray([40]byte{1, 1, 1, 1, 'A', 'B', 'C'}),
	AliasFixedPointerArrayField: AliasFixedPointerArray([3]*uint16{&testUint16, nil, &testUint16}),
	AliasFixedSubMessageArrayField: AliasFixedSubMessageArray([3]SubMessage{
		{
			Id:          1,
			Description: "A",
		},
		{
			Id:          2,
			Description: "B",
		},
		{
			Id:          3,
			Description: "C",
		},
	}),
	AliasFixedSubMessagePointerArrayField: AliasFixedSubMessagePointerArray([3]*SubMessage{
		&testSubMessage, nil, &testSubMessage,
	}),
}

func TestMessage(t *testing.T) {
	bufferOfficial := encoding.Marshal(msg)

	bufferUnofficial := &encodegen.ObjBuffer{}
	msg.MarshalBuffer(bufferUnofficial)

	assert.Equal(t, bufferOfficial, bufferUnofficial.Bytes(), "Generated buffer does not equal one generated by reflection")

	newMessageUnofficial := Message{}
	(*Message)(&newMessageUnofficial).UnmarshalBuffer(bufferUnofficial)

	newMessageOfficial := Message{}
	encoding.Unmarshal(bufferOfficial, &newMessageOfficial)

	assert.Equal(t, newMessageOfficial, newMessageUnofficial, "Unmarshaled struct does not equal one generated by reflection")
}

func BenchmarkMessageMarshalReflect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encoding.Marshal(msg)
	}
}

func BenchmarkMessageMarshalCodegen(b *testing.B) {
	buffer := &encodegen.ObjBuffer{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		msg.MarshalBuffer(buffer)
	}
	assert.Equal(b, encoding.Marshal(msg), buffer.Bytes(), "Codegen didn't produce valid output, disregard this benchmark.")
}

func BenchmarkMessageUnmarshalReflect(b *testing.B) {
	data := encoding.Marshal(msg)
	message := Message{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoding.Unmarshal(data, &message)
	}
}

func BenchmarkMessageUnmarshalCodegen(b *testing.B) {
	buffer := &encodegen.ObjBuffer{}
	msg.MarshalBuffer(buffer)

	message := Message{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		(*Message)(&message).UnmarshalBuffer(buffer)
		buffer.Rewind()
	}

	messageOfficial := Message{}
	encoding.Unmarshal(buffer.Bytes(), &messageOfficial)

	assert.Equal(b, messageOfficial, message, "Codegen didn't produce valid output, disregard this benchmark.")
}
