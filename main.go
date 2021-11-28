package main

import "log"
import "regexp"
import "fmt"

// Token
type Token int

const (
	INT Token = iota

	operator_begin
	PLUS
	MINUS
	operator_end
)

func token(literal string) Token {
	switch literal {
	case "+":
		return PLUS
	case "-":
		return MINUS
	default:
		return INT
	}
}

// Word
type Word struct {
	tok Token
	lit string
}

func (w Word) String() string {
	return fmt.Sprintf("Word{%v,%v}", w.tok, w.lit)
}

// Tokenizer
var tokenizer *regexp.Regexp

func init() {
	tokenizer = regexp.MustCompile(`\s*(\w+|[+-])`)
}

func Tokenize(code string) []Word {
	log.Printf("tokenize input code = %#v", code)

	words := make([]Word, 0)

	for _, group := range tokenizer.FindAllStringSubmatch(code, -1) {
		lit := group[1]
		word := Word{tok: token(lit), lit: lit}
		words = append(words, word)
	}

	return words
}

func main() {
	log.SetPrefix("forth: ")
	log.SetFlags(0)

	tokens := Tokenize("100 200 +")
	log.Print(tokens)

	tokens = Tokenize("100 200 -")
	log.Print(tokens)
}
