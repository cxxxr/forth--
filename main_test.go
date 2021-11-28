package main

import "testing"

func testParse(t *testing.T, code string, expectedTokens []Word) {
	actualTokens := Parse(code)

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

func newIntWord(lit string) Word {
	return Word{lit: lit, tok: INT}
}

func newPlusWord() Word {
	return Word{lit: "+", tok: PLUS}
}

func newMinusWord() Word {
	return Word{lit: "-", tok: MINUS}
}

func TestParse(t *testing.T) {
	testParse(t, "3 4 +",
		[]Word{newIntWord("3"), newIntWord("4"), newPlusWord()})

	testParse(t, "123 456 +",
		[]Word{newIntWord("123"), newIntWord("456"), newPlusWord()})

	testParse(t, "   123  456      +",
		[]Word{newIntWord("123"), newIntWord("456"), newPlusWord()})

	testParse(t, "3 4 -",
		[]Word{newIntWord("3"), newIntWord("4"), newMinusWord()})

	testParse(t, "42",
		[]Word{newIntWord("42")})

	testParse(t, "+",
		[]Word{newPlusWord()})
}
