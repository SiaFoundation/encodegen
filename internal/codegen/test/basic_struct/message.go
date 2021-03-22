package basic_struct

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
}
