package main

import "testing"

func TestTokenize(t *testing.T) {
	actual := Tokenize("3 4 +")
	expected := []string{"3", "4", "+"}

	if len(actual) != len(expected) {
		t.Fatalf("expected len(actual) = len(expected), actual = %#v", actual)
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("!?: expected = %#v, actual = %#v", actual[i], expected[i])
		}
	}
}
