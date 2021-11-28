package main

import "testing"

func testTokenize(t *testing.T, code string, expectedTokens []string) {
	actualTokens := Tokenize(code)

	if len(actualTokens) != len(expectedTokens) {
		t.Fatalf("expected len(actual) = len(expected), actual = %#v", actualTokens)
	}

	for i := range expectedTokens {
		if actualTokens[i] != expectedTokens[i] {
			t.Fatalf("!?: expected = %#v, actual = %#v",
				actualTokens[i],
				expectedTokens[i])
		}
	}
}

func TestTokenize(t *testing.T) {
	testTokenize(t, "3 4 +", []string{"3", "4", "+"})
	testTokenize(t, "123 456 +", []string{"123", "456", "+"})
	testTokenize(t, "   123  456      +", []string{"123", "456", "+"})
	testTokenize(t, "3 4 -", []string{"3", "4", "-"})
	testTokenize(t, "42", []string{"42"})
	testTokenize(t, "+", []string{"+"})
}
