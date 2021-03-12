package embedded_struct

type AliasSubMessage SubMessage

type Message struct {
	Id        int
	Anonymous struct {
		IntegerField      int
		StringField       string
		IntegerSliceField []int
		Sub               SubMessage
		AliasSub          AliasSubMessage
		Anonymous3        struct {
			A          int
			Anonymous4 struct {
				B          int
				Anonymous5 *struct {
					WWW        *int
					Anonymous6 struct{}
				}
			}
			C int
		}
	}
	Anonymous2 *struct {
		IntegerField int
		Anonymous3   *struct {
			A int
		}
	}
	Empty struct{}
	End   int
}
