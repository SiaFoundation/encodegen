package alias_struct

import (
	// "log"
	"github.com/stretchr/testify/assert"
	"gitlab.com/NebulousLabs/encoding"
	"go.sia.tech/encodegen/pkg/encodegen"
	"testing"
)

var testInt int = 999
var testAliasInt AliasInt = AliasInt(9001)
var doubleAliasSubMessageField = DoubleAliasSubMessage(AliasSubMessage(SubMessage{Id: 888, Description: "bbb"}))
var testSubMessage = SubMessage{Id: 8888, Description: "AAAAA", Strings: []string{"WASDF", "WAASDSAD"}}

var msg = Message{
	Id:                                1022,
	Sub:                               SubMessage{Id: 9001, Description: "AAA"},
	AliasSubMessageField:              AliasSubMessage(SubMessage{Id: 9}),
	ArrayAliasSubMessageField:         []AliasSubMessage{AliasSubMessage(SubMessage{Id: 8, Description: "aaaa", Strings: []string{"9"}})},
	DoubleAliasSubMessageField:        doubleAliasSubMessageField,
	PointerDoubleAliasSubMessageField: &doubleAliasSubMessageField,
	AliasIntField:                     testAliasInt,
	PointerAliasIntField:              &testAliasInt,
	AliasIntArrayField:                AliasIntArray([]int{5, 5, 5, 5, 999995, 5, 5, 5, 9, -5}),
	AliasIntPointerArrayField:         AliasIntPointerArray([]*int{nil, &testInt, nil, nil, &testInt}),
	AliasSubMessageArrayField: AliasSubMessageArray([]SubMessage{
		{Id: 4000, Description: "111", Strings: []string{"999", "888"}},
		{Id: 8111, Description: "111", Strings: []string{}},
		{Id: 8111, Description: "111", Strings: nil},
		{Id: 8111, Strings: nil},
		{Strings: nil},
	}),
	AliasSubMessagePointerArrayField: AliasSubMessagePointerArray([]*SubMessage{
		&testSubMessage,
		nil,
		&testSubMessage,
		&testSubMessage,
		&testSubMessage,
		nil,
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