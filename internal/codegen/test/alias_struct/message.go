package alias_struct

// import 	"gitlab.com/NebulousLabs/Sia/crypto"

type DoubleAliasSubMessage AliasSubMessage
type AliasSubMessage SubMessage

type AliasInt int
type DoubleAliasInt AliasInt
type AliasIntArray []int
type AliasIntPointerArray []*int

type AliasSubMessageArray []SubMessage
type AliasSubMessagePointerArray []*SubMessage

type Message struct {
	Id                                    int
	Sub                                   SubMessage
	AliasSubMessageField                  AliasSubMessage
	ArrayAliasSubMessageField             []AliasSubMessage
	DoubleAliasSubMessageField            DoubleAliasSubMessage
	PointerDoubleAliasSubMessageField     *DoubleAliasSubMessage
	AliasIntField                         AliasInt
	PointerAliasIntField                  *AliasInt
	AliasIntArrayField                    AliasIntArray
	AliasIntPointerArrayField             AliasIntPointerArray
	AliasSubMessageArrayField             AliasSubMessageArray
	AliasSubMessagePointerArrayField      AliasSubMessagePointerArray
	ArrayAliasSubMessagePointerArrayField []AliasSubMessagePointerArray
	DoubleAliasIntField                   DoubleAliasInt
}
