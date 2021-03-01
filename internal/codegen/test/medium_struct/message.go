package medium_struct

type AliasSubMessage SubMessage

type Message struct {
	Id          int
	Name        string
	Ints        []*int
	SubMessageX *SubMessage
	MessagesX   []*SubMessage
	SubMessageY SubMessage
	MessagesY   []SubMessage
	IsTrue      *bool
	Payload     []byte
	Uint64      uint64
	// TODO
	// AliasedSubmessage AliasSubMessage
}
