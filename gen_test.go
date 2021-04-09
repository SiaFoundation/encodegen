package main

import "testing"

func TestGenerate(t *testing.T) {
	src, err := Generate("./test", "TestMessageSimple", "TestMessageSecond")
	if err != nil {
		t.Fatal(err)
	}
	//t.FailNow()
	t.Fatal(src)
}
