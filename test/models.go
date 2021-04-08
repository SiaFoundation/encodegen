package test

type Imported struct {
	A int
	B string
	C bool
}

type Hash [32]byte

type TestMessageSimple struct {
	A uint64
	B uint32
	C uint16
	D uint8
	E bool
	F [32]Hash
	G []byte
}

type TestMessageEmbedded struct {
	A struct {
		A uint64
		B uint32
		C uint16
		D uint8
		E int8
		F bool
		G []byte
		H struct {
			I []int
		}
	}
}
