package main

type TestType1 struct {
	A  int64
	B  *uint64
	Z  []string
	ZZ []*string
	D  TestType2
}

type TestType2 struct {
	A uint64
	D uint32
	G []*TestType3
}

type TestType3 struct {
	F []byte
}
