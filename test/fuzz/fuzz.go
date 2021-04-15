package test

import (
	"bytes"
	"go.sia.tech/encodegen/test"
)

func Fuzz(data []byte) int {
	simpleMessage := test.TestMessageSimple{}
	err := simpleMessage.UnmarshalSia(bytes.NewReader(data))
	if err == nil {
		// function must return 1 if the fuzzer should increase priority of the given input during subsequent fuzzing (for example, the input is lexically correct and was parsed successfully);
		return 1
	}
	// 0 otherwise
	return 0
}
