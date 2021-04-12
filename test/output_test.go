package test

import (
	"bytes"
	"gitlab.com/NebulousLabs/encoding"
	"go.sia.tech/encodegen/test/imported"
	"io"
	"reflect"
	"testing"
)

/*
Separate directory for tests so that generated source output does not pollute the main source directory

Functions for test types must be generated for these types to work.
*/

// all types that have generated code on them implement this interface
type SiaMarshaler interface {
	MarshalSia(w io.Writer) error
}

// implemented on pointers of objects
type SiaUnmarshaler interface {
	UnmarshalSia(r io.Reader) error
}

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
}

func TestSimpleMessage(t *testing.T) {

	messageBytes := compareMarshal(t, simpleMessage)
	compareUnmarshal(t, messageBytes, &TestMessageSimple{}, &TestMessageSimple{})

	messageBytes = compareMarshal(t, embeddedMessage)
	compareUnmarshal(t, messageBytes, &TestMessageEmbedded{}, &TestMessageEmbedded{})
}

func compareMarshal(t *testing.T, obj interface{}) []byte {
	impl, ok := obj.(SiaMarshaler)
	if !ok {
		t.Fatal("Type does not implement SiaMarshaler.  Make sure to generate the code before running these tests.")
	}

	bufferUnofficial := &bytes.Buffer{}
	impl.MarshalSia(bufferUnofficial)

	bytesOfficial := encoding.Marshal(obj)

	if !reflect.DeepEqual(bytesOfficial, bufferUnofficial.Bytes()) {
		t.Fatalf("Generated buffer (%+v) does not equal one generated by reflection (%+v)", bufferUnofficial.Bytes(), bytesOfficial)
	}
	return bytesOfficial
}

func compareUnmarshal(t *testing.T, marshaledBytes []byte, obj1 interface{}, obj2 interface{}) {
	impl, ok := obj1.(SiaUnmarshaler)
	if !ok {
		t.Fatal("Type does not implement SiaUnmarshaler.  Make sure to generate the code before running these tests.")
	}

	impl.UnmarshalSia(bytes.NewReader(marshaledBytes))
	encoding.Unmarshal(marshaledBytes, obj2)

	if !reflect.DeepEqual(obj2, impl) {
		t.Fatalf("Unmarshaled object (%+v) does not equal one generated by NebulousLabs/encoding (%+v).", impl, obj2)
	}
}

/*
Although the codegen approach still performs better and allocates less,
the reflection based library is getting a "free ride" and having its performance
significantly overstated because it can simply call (Un)MarshalSia methods on the
types (which it checks for in https://gitlab.com/NebulousLabs/encoding/-/blob/master/marshal.go#L155)
and not have to do most of the reflection it normally would have to.

Requiring people to run only the reflection based benchmarks, then generate the
code, then run the codegen based benchmarks seems hacky.  I tried aliasing the
types passed into the reflection based library but ran into some errors and
realized that that fields in the struct (like the field for an imported.Imported struct)
would still have their MarshalSia methods invoked.  Not sure how to solve this.
*/

func BenchmarkSimpleMessageCodegenMarshal(b *testing.B) {
	benchmarkCodegenMarshal(b, simpleMessage)
}
func BenchmarkSimpleMessageCodegenUnmarshal(b *testing.B) {
	benchmarkCodegenUnmarshal(b, &simpleMessage, &TestMessageSimple{})
}

func BenchmarkEmbeddedMessageCodegenMarshal(b *testing.B) {
	benchmarkCodegenMarshal(b, embeddedMessage)
}
func BenchmarkEmbeddedMessageCodegenUnmarshal(b *testing.B) {
	benchmarkCodegenUnmarshal(b, &embeddedMessage, &TestMessageEmbedded{})
}

func BenchmarkSimpleMessageOfficialMarshal(b *testing.B) {
	benchmarkOfficialMarshal(b, simpleMessage)
}

func BenchmarkSimpleMessageOfficialUnmarshal(b *testing.B) {
	benchmarkOfficialUnmarshal(b, &simpleMessage, &TestMessageSimple{})
}

func BenchmarkEmbeddedMessageOfficialMarshal(b *testing.B) {
	benchmarkOfficialMarshal(b, simpleMessage)
}

func BenchmarkEmbeddedMessageOfficialUnmarshal(b *testing.B) {
	benchmarkOfficialUnmarshal(b, &embeddedMessage, &TestMessageEmbedded{})
}

func benchmarkCodegenMarshal(b *testing.B, obj interface{}) {
	impl, ok := obj.(SiaMarshaler)
	if !ok {
		b.Fatal("Type does not implement SiaMarshaler.  Make sure to generate the code before running these tests.")
	}
	buffer := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		impl.MarshalSia(buffer)
		buffer.Reset()
	}
}

func benchmarkCodegenUnmarshal(b *testing.B, obj interface{}, dst interface{}) {
	impl, ok := obj.(SiaUnmarshaler)
	if !ok {
		b.Fatal("Type does not implement SiaMarshaler.  Make sure to generate the code before running these tests.")
	}
	data := encoding.Marshal(obj)
	for i := 0; i < b.N; i++ {
		impl.UnmarshalSia(bytes.NewReader(data))
	}
}

func benchmarkOfficialMarshal(b *testing.B, obj interface{}) {
	buffer := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		// I'm not sure this benchmark is fair because we have to allocate a new encoder object (which just contains the io.writer) every iteration (the encoder object does not allow you to reset the buffer).
		// However the buffer itself is reused so the difference should not be large.
		encoding.NewEncoder(buffer).Encode(obj)
		buffer.Reset()
	}
}

func benchmarkOfficialUnmarshal(b *testing.B, obj interface{}, dst interface{}) {
	data := encoding.Marshal(simpleMessage)
	msg := TestMessageSimple{}
	for i := 0; i < b.N; i++ {
		encoding.Unmarshal(data, &msg)
	}
}
