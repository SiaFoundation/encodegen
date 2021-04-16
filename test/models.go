package test

//go:generate encodegen -t TestMessageSimple,TestMessageEmbedded

import (
	imported "go.sia.tech/encodegen/test/imported"
	importedrename "go.sia.tech/encodegen/test/imported"
	importedchild "go.sia.tech/encodegen/test/imported/imported"
	importedchild1 "go.sia.tech/encodegen/test/imported/imported"
	importedchild2 "go.sia.tech/encodegen/test/imported/imported"
	importedchild3 "go.sia.tech/encodegen/test/imported/imported/imported"
	"go.sia.tech/encodegen/test/imported/subimported"
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
	M subimported.HashAlias
	O []importedrename.Imported
	P []importedchild.Imported
	Q *importedchild1.Imported
	R []importedchild2.Imported
	S []*importedchild3.Imported
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
	B []struct {
		A int
		B string
		C bool
		D imported.Imported
	}
	C []*struct {
		A int
		B string
		C bool
	}
}
