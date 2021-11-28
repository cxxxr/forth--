package main

import "testing"

func TestRegressionBuiltinProc(t *testing.T) {
	words := Parse("100 200 +")

	env := NewEnv()
	if err := env.Execute(words); err != nil {
		t.Fatal(err)
	}

	actual, err := env.stack.Peek()
	if err != nil {
		t.Fatal(err)
	}

	v := actual.(*Int)
	if v.v != 300 {
		t.Fatalf("wrong something: %v", v)
	}

	println(v.v)
}
