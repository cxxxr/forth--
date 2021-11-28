package main

import "log"

func Tokenize(code string) []string {
	log.Printf("tokenize input code = %#v", code)
	tokens := []string{"3", "4", "+"}
	return tokens
}

func main() {
	log.SetPrefix("forth: ")
	log.SetFlags(0)

	tokens := Tokenize("3 4 +")
	log.Print(tokens)
}
