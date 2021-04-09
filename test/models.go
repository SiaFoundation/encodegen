package test

import (
	"go.sia.tech/encodegen/test/imported"
)

type Hash [32]byte

type TestMessageSimple struct {
	A uint64
	B uint32
	C uint16
	D uint8
	E bool
	F [32]Hash
	G []byte
	H imported.Imported
	I *uint64
	J *imported.Imported
	K ***uint64
	L []*imported.Imported
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
