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

var tokenizer *regexp.Regexp

func init() {
	tokenizer = regexp.MustCompile(`\s*([\w.:;]+|[+-])`)
}

func Parse(code string) []Token {
	tokens := make([]Token, 0)
	for _, group := range tokenizer.FindAllStringSubmatch(code, -1) {
		token := Token{Lit: group[1]}
		tokens = append(tokens, token)
	}
	return tokens
}
