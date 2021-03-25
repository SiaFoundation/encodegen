package basic_struct

import (
	"go.sia.tech/encodegen/internal/codegen/test/alias_struct"
	importedtyperename "go.sia.tech/encodegen/internal/codegen/test/importedtype"
)

// this is a test to ensure that alias_struct is not imported in encoding.go (only importedtype should be imported because it is the only one used in a struct field)
var _ = alias_struct.Message{}

type Message struct {
	Id                     int
	Name                   string
	Ints                   []int
	Uint8s                 []uint8
	SubMessageX            *SubMessage
	MessagesX              []*SubMessage
	SubMessageY            SubMessage
	MessagesY              []SubMessage
	IsTrue                 *bool
	Payload                []byte
	Strings                []string
	FixedBytes             [9]byte
	FixedInts              [5]int
	FixedIntPointers       [40]*int
	FixedUint8s            [40]uint8
	FixedSubMessage        [5]SubMessage
	FixedPointerSubMessage [5]*SubMessage
	Imported               importedtyperename.Imported
	ImportedPointer        *importedtyperename.Imported
	ImportedSlice          []importedtyperename.Imported
	ImportedPointerSlice   []*importedtyperename.Imported
	FixedImported          [3]importedtyperename.Imported
	FixedPointerImported   [3]*importedtyperename.Imported
}
