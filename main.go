package main

import "log"
import "regexp"

func Tokenize(code string) []string {
	log.Printf("tokenize input code = %#v", code)

	tokenizer := regexp.MustCompile(`\s*(\w+|[+-])`)

	tokens := make([]string, 0)

	for _, group := range tokenizer.FindAllStringSubmatch(code, -1) {
		token := group[1]
		tokens = append(tokens, token)
	}

	return tokens
}

func main() {
	log.SetPrefix("forth: ")
	log.SetFlags(0)

	tokens := Tokenize("100 200 +")
	log.Print(tokens)

	tokens = Tokenize("100 200 -")
	log.Print(tokens)
}
