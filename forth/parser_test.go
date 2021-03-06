package forth

import (
	"testing"
)

type testCase struct {
	code   string
	tokens []Token
}

func testParse(t *testing.T, tc testCase) {
	actual := Parse(tc.code)
	expected := tc.tokens

	if len(actual) != len(expected) {
		t.Fatalf(
			"len(actual) != len(expected)\nactual: %v" +
			"expected: %v\n",
			actual,
			expected,
		)
	}

	for i := range actual {
		if !actual[i].Eq(&expected[i]) {
			t.Fatalf("%v != %v", actual[i], expected[i])
		}
	}
}

func TestParse(t *testing.T) {
	testCases := []testCase{
		{
			code: "1 2 +",
			tokens: []Token{
				{Lit: "1"},
				{Lit: "2"},
				{Lit: "+"},
			},
		},
		{
			code: "100   200 +",
			tokens: []Token{
				{Lit: "100"},
				{Lit: "200"},
				{Lit: "+"},
			},
		},
		{
			code: "1 2 -",
			tokens: []Token{
				{Lit: "1"},
				{Lit: "2"},
				{Lit: "-"},
			},
		},
		{
			code: "  foo bar hoge",
			tokens: []Token{
				{Lit: "foo"},
				{Lit: "bar"},
				{Lit: "hoge"},
			},
		},
		{
			code: ": 2x 2 * ;",
			tokens: []Token{
				{Lit: ":"},
				{Lit: "2x"},
				{Lit: "2"},
				{Lit: "*"},
				{Lit: ";"},
			},
		},
		{
			code: ": 2+ 2 + ;",
			tokens: []Token{
				{Lit: ":"},
				{Lit: "2+"},
				{Lit: "2"},
				{Lit: "+"},
				{Lit: ";"},
			},
		},
	}

	for _, tc := range testCases {
		testParse(t, tc)
	}
}
