package basic_struct

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/NebulousLabs/encoding"
	"go.sia.tech/encodegen/pkg/encodegen"
	"testing"
)

var isTrue = true
var testInt = 511
var testSubMessage = SubMessage{
	Id:          -1,
	Description: "ABC",
	Strings:     []string{"A", "B", "C", "X", "Y", "Z"},
}
var msg = Message{
	Id:     1022,
	Name:   "name acc",
	Ints:   []int{1, 2, 5},
	Uint8s: []uint8{1, 2, 5},
	SubMessageX: &SubMessage{
		Id:          102,
		Description: "abcd",
	},
	MessagesX: []*SubMessage{
		{
			Id:          2102,
			Description: "abce",
		},
	},
	SubMessageY: SubMessage{
		Id:          3102,
		Description: "abcf",
	},
	MessagesY: []SubMessage{
		{
			Id:          5102,
			Description: "abcg",
		},
		{
			Id:          5106,
			Description: "abcgg",
		},
	},
	IsTrue:           &isTrue,
	Payload:          []byte{11, 1, 1, 123, 123, 123, 123},
	FixedBytes:       [9]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I'},
	FixedInts:        [5]int{4, 4, 4, 4, 4},
	FixedIntPointers: [40]*int{nil, &testInt, nil, &testInt, nil, nil, nil},
	FixedUint8s:      [40]uint8{1, 1, 1, 1, 1, 1, 1, 0},
	FixedSubMessage: [5]SubMessage{
		{
			Id:          500,
			Description: "AAAA",
			Strings:     nil,
		},
		{
			Id:          444,
			Description: "WWW",
			Strings:     []string{"W", "W", "W"},
		},
	},
	FixedPointerSubMessage: [5]*SubMessage{nil, &testSubMessage, &testSubMessage, &testSubMessage, nil},
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
