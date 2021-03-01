package basic_struct

type Message struct {
	Id          int
	Name        string
	Ints        []int
	SubMessageX *SubMessage
	MessagesX   []*SubMessage
	SubMessageY SubMessage
	MessagesY   []SubMessage
	IsTrue      *bool
	Payload     []byte
	// SQLNullString *sql.NullString
}
