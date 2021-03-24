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

type AliasByteSlice []byte

type AliasFixedByteArray [40]byte
type AliasFixedPointerArray [3]*uint16

type AliasFixedSubMessageArray [3]SubMessage
type AliasFixedSubMessagePointerArray [3]*SubMessage

type Integer []*int

type Message struct {
	Id                                         int
	Sub                                        SubMessage
	AliasSubMessageField                       AliasSubMessage
	ArrayAliasSubMessageField                  []AliasSubMessage
	DoubleAliasSubMessageField                 DoubleAliasSubMessage
	PointerDoubleAliasSubMessageField          *DoubleAliasSubMessage
	AliasIntField                              AliasInt
	PointerAliasIntField                       *AliasInt
	AliasIntArrayField                         AliasIntArray
	AliasIntPointerArrayField                  AliasIntPointerArray
	AliasSubMessageArrayField                  AliasSubMessageArray
	AliasSubMessagePointerArrayField           AliasSubMessagePointerArray
	ArrayAliasSubMessagePointerArrayField      []AliasSubMessagePointerArray
	DoubleAliasIntField                        DoubleAliasInt
	ByteSlice                                  []byte
	AliasByteSliceField                        AliasByteSlice
	AliasFixedByteArrayField                   AliasFixedByteArray
	AliasFixedPointerArrayField                AliasFixedPointerArray
	AliasFixedSubMessageArrayField             AliasFixedSubMessageArray
	AliasFixedSubMessagePointerArrayField      AliasFixedSubMessagePointerArray
	AliasFixedByteArrayArrayField              [3]AliasFixedByteArray
	AliasFixedPointerArrayArrayField           [3]AliasFixedPointerArray
	AliasFixedSubMessageArrayArrayField        [3]AliasFixedSubMessageArray
	AliasFixedSubMessagePointerArrayArrayField [3]AliasFixedSubMessagePointerArray
	IntegerField                               Integer
}
