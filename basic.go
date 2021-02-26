package main

type Message struct {
	Id   int
	Name string
	// Price       float64
	Ints []int
	// Floats      []float32
	SubMessageX *SubMessage
	MessagesX   []*SubMessage
	SubMessageY SubMessage
	MessagesY   []SubMessage
	IsTrue      *bool
	Payload     []byte
}

type SubMessage struct {
	Id          int
	Description string
}
