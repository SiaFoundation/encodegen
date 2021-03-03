package medium_struct

type DoubleAliasSubMessage AliasSubMessage
type AliasSubMessage SubMessage

type AliasInt int
type DoubleAliasInt AliasInt

type Message struct {
	Id          int
	Name        string
	Names []string
	Ints        []*int
	SubMessageX *SubMessage
	MessagesX   []*SubMessage
	SubMessageY SubMessage
	MessagesY   []SubMessage
	IsTrue      *bool
	Payload     []byte
	Uint64      uint64
	// TODO
	AliasedSubmessage              AliasSubMessage
	ArrayAliasSubmessage           []AliasSubMessage
	DoubleAliasedSubmessage        DoubleAliasSubMessage
	PointerDoubleAliasedSubmessage *DoubleAliasSubMessage
	AliasInt                       AliasInt
	PointerAliasInt                *AliasInt
	ArrayAliasInt                  []AliasInt
	DoubleAliasInt                 DoubleAliasInt
}
