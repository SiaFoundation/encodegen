//+build ignore

package main

import (
	"github.com/dvyukov/go-fuzz/gen"
	"github.com/erggo/datafiller"
	"gitlab.com/NebulousLabs/encoding"
	"go.sia.tech/encodegen/internal/test"
)

func main() {
	for {
		var msg test.TestMessageSimple
		datafiller.Fill(&msg)
		gen.Emit(encoding.Marshal(msg), nil, true)
	}
}
