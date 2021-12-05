package main

import (
	"testing"

	"github.com/cxxxr/forth--/forth"
)

func testExecute(t *testing.T, code string, expectedPeek int) {
	words := forth.Parse(code)

	env := NewEnv()
	if err := env.Execute(words); err != nil {
		t.Fatal(err)
	}

	actual, err := env.stack.Peek()
	if err != nil {
		t.Fatal(err)
	}

	v := actual.(*Int)
	if v.v != ForthInt(expectedPeek) {
		t.Fatalf("wrong something: %v", v)
	}

	println(v.v)
}

func TestRegressionBuiltinProc(t *testing.T) {
	testExecute(t, "100 200 +", 300)
}

func TestRegressionUserDefinedProc(t *testing.T) {
	testExecute(t, ": 2+ 2 + ; 10 2+ .s", 12)
}
