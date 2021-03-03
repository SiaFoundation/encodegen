package basic_struct

import (
	// "log"
	"github.com/stretchr/testify/assert"
	"gitlab.com/NebulousLabs/encoding"
	"go.sia.tech/encodegen/pkg/encodegen"
	"testing"
)

var isTrue = true
var msg = Message{
	Id:   1022,
	Name: "name acc",
	Ints: []int{1, 2, 5},
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
	IsTrue:  &isTrue,
	Payload: []byte(`"123"`),
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
