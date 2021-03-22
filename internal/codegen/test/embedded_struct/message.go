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
		AnonymousSub1     struct {
			A             int
			AnonymousSub2 struct {
				A int
			}
		}
	}
	Anonymous2 *struct {
		IntegerField   int
		Anonymous2Sub1 *struct {
			A int
		}
	}
	Anonymous3 struct{}
	Anonymous4 []struct {
		A                 int
		IntegerSliceField []int
		Anonymous5        []struct {
			A []int
			B *struct {
				A int
			}
		}
	}
	Anonymous5 []struct {
		A int
		B []AliasSubMessage
	}
	Anonymous6            []*struct{ A int }
	AnonymousFixed        [5]struct{ A int }
	AnonymousPointerFixed [5]*struct{ A int }
	End                   int
}
