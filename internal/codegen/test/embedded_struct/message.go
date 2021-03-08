package embedded_struct

type Message struct {
	Id        int
	Anonymous struct {
		IntegerField      int
		StringField       string
		IntegerSliceField []int
		Sub               SubMessage
	}
	End int
}
