package main

type TestType1 struct {
	A ***int64
	B uint64
	C []*string
	D TestType2
}

type TestType2 struct {
	A uint64
	D uint32
	E []TestType3
}

type TestType3 struct {
	F []byte
}
