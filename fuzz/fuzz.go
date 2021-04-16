package fuzz

import (
	"gitlab.com/NebulousLabs/encoding"
	"go.sia.tech/encodegen/test"
)

func Fuzz(data []byte) int {
	var msg test.TestMessageSimple
	err := encoding.Unmarshal(data, &msg)
	if err == nil {
		return 1
	}
	return 0
}
