package main

type TestType1 struct {
	A int
	B uint64
	C string
	D []byte
	E TestType2
}

type TestType2 struct {
	F uint64
	G []uint64
	H TestType3
}

type TestType3 struct {
	I string
}
