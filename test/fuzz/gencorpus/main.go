package main

import (
	"bytes"
	"github.com/dvyukov/go-fuzz/gen"
	"github.com/erggo/datafiller"
	"go.sia.tech/encodegen/test"
)

func main() {
	for {
		simpleMessage := test.TestMessageSimple{}
		datafiller.Fill(&simpleMessage)
		buffer := new(bytes.Buffer)
		simpleMessage.MarshalSia(buffer)
		gen.Emit(buffer.Bytes(), nil, true)
	}
}
