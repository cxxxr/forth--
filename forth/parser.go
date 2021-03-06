package forth

import (
	"fmt"
	"regexp"
)

type Token struct {
	Lit string // TODO: unexport
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%v}", t.Lit)
}

func (t *Token) Eq(other *Token) bool {
	return t.Lit == other.Lit
}

var tokenizer *regexp.Regexp

func init() {
	tokenizer = regexp.MustCompile(`\s*(\S+)`)
}

func Parse(code string) []Token {
	tokens := make([]Token, 0)
	for _, group := range tokenizer.FindAllStringSubmatch(code, -1) {
		lit := group[1]
		tokens = append(tokens, Token{lit})
	}
	return tokens
}
