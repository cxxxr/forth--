package forth

import (
	"fmt"
	"regexp"
	"strconv"
)

 // TODO: main.goの定義との重複を解消する
const intBitSize = 32
type ForthInt int32

type Token struct {
	Lit string // TODO: unexport
	intValue ForthInt // TODO: Litがintの場合だけ値が入るがこれは良くないので直す
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

func parseToken(lit string) Token {
	if v, err := strconv.ParseInt(lit, 10, intBitSize); err == nil {
		return Token{Lit: lit, intValue: ForthInt(v)}
	} else {
		return Token{Lit: lit}
	}
}

func Parse(code string) []Token {
	tokens := make([]Token, 0)
	for _, group := range tokenizer.FindAllStringSubmatch(code, -1) {
		lit := group[1]
		tokens = append(tokens, parseToken(lit))
	}
	return tokens
}
