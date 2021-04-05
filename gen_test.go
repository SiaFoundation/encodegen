package main

import "testing"

func TestGenerate(t *testing.T) {
	src, err := Generate(".", "Array", "Hash")
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(src)
}
