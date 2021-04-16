package test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	"gitlab.com/NebulousLabs/encoding"
	"go.sia.tech/encodegen/internal/test/imported"
)

var simpleMessage = TestMessageSimple{
	A: 5,
	B: 4,
	C: 3,
	D: 2,
	E: true,
	F: [32]Hash{},
	G: []byte{1, 1, 1, 1, 0},
	H: imported.Imported{A: 5, B: "AAA", C: true},
	I: new(uint64),
	J: &imported.Imported{A: 555, B: "AAA", C: false},
	K: new(**uint64),
	L: []*imported.Imported{nil, {A: 999}, nil, nil, {C: true}},
	M: [32]byte{5},
}

var embeddedMessage = TestMessageEmbedded{
	A: struct {
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
	}{
		A: 1,
		B: 1,
		C: 1,
		D: 1,
		E: 1,
		F: false,
		G: []byte{0, 1, 2, 3, 4},
		H: struct{ I []int }{I: []int{99, 0, 99}},
	},
	B: []struct {
		A int
		B string
		C bool
		D imported.Imported
	}{
		{A: 5, B: "X", C: true, D: imported.Imported{A: 5, B: "AAA", C: true}},
		{A: 1, B: "AAA", C: false},
	},
	C: []*struct {
		A int
		B string
		C bool
	}{
		{A: 5, B: "X", C: true},
		{A: 1, B: "AAA", C: false},
	},
}

var update = flag.Bool("update", false, "update .golden files")

func TestGolden(t *testing.T) {
	tests := []struct {
		name string
		obj  interface{}
		typ  interface{}
	}{
		{"simple", simpleMessage, new(TestMessageSimple)},
		{"embedded", embeddedMessage, new(TestMessageEmbedded)},
	}

	if *update {
		for _, test := range tests {
			path := fmt.Sprintf("testdata/%v.golden", test.name)
			err := ioutil.WriteFile(path, encoding.Marshal(test.obj), 0660)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	for _, test := range tests {
		path := fmt.Sprintf("testdata/%v.golden", test.name)
		golden, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(encoding.Marshal(test.obj), golden) {
			t.Errorf("encoded %T did not match golden file", test.obj)
		}
		if err := encoding.Unmarshal(golden, test.typ); err != nil {
			t.Errorf("decoding into %T failed: %v", test.typ, err)
		}
	}
}

func BenchmarkMarshal(b *testing.B) {
	objs := []interface{}{
		simpleMessage,
		embeddedMessage,
	}
	for _, obj := range objs {
		b.Run(fmt.Sprintf("%T", obj), func(b *testing.B) {
			b.SetBytes(int64(len(encoding.Marshal(obj))))
			for i := 0; i < b.N; i++ {
				encoding.Marshal(obj)
			}
		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	tests := []struct {
		obj interface{}
		typ interface{}
	}{
		{simpleMessage, new(TestMessageSimple)},
		{embeddedMessage, new(TestMessageEmbedded)},
	}
	for _, test := range tests {
		b.Run(fmt.Sprintf("%T", test.obj), func(b *testing.B) {
			data := encoding.Marshal(test.obj)
			b.SetBytes(int64(len(data)))
			for i := 0; i < b.N; i++ {
				encoding.Unmarshal(data, test.typ)
			}
		})
	}
}
