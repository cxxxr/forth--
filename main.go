package main

import "log"
import "regexp"
import "fmt"
import "strconv"

// TokenOld
type TokenOld int

const (
	INT TokenOld = iota

	operator_begin
	PLUS
	MINUS
	PRINT
	operator_end
)

func token(literal string) TokenOld {
	switch literal {
	case "+":
		return PLUS
	case "-":
		return MINUS
	case ".s":
		return PRINT
	default:
		return INT
	}
}

// Token
// リテラルから見たときと実行時に見たときでTokenの扱いは違うはずなので
// 本当は分けたほうが良いかもしれないが今は簡単さを優先する
type Token struct {
	tok TokenOld
	lit string
}

func (w Token) String() string {
	return fmt.Sprintf("Token{%v}", w.lit)
}

// Parse
var tokenizer *regexp.Regexp

func init() {
	tokenizer = regexp.MustCompile(`\s*([\w.]+|[+-])`)
}

// TODO: error handling
func Parse(code string) []Token {
	log.Printf("tokenize input code = %#v", code)

	tokens := make([]Token, 0)

	for _, group := range tokenizer.FindAllStringSubmatch(code, -1) {
		lit := group[1]
		word := Token{lit: lit}
		tokens = append(tokens, word)
	}

	return tokens
}

// Execute

type Cell interface {
	String() string
}

type Int struct {
	v int
}

func (i Int) String() string {
	return fmt.Sprint(i.v)
}

type Stack struct {
	data []Cell
}

type Env struct {
	stack *Stack
}

func NewInt(v int) *Int {
	return &Int{v: v}
}

func (stack *Stack) Peek() (Cell, error) {
	if len(stack.data) == 0 {
		return nil, fmt.Errorf("stack underflow: %#v", stack.data)
	}
	v := stack.data[len(stack.data)-1]
	return v, nil
}

func (stack *Stack) Pop() (Cell, error) {
	v, err := stack.Peek()
	if err != nil {
		return nil, err
	}
	stack.data = stack.data[:len(stack.data)-1]
	log.Printf("poped: %#v", stack.data)
	return v, nil
}

func (stack *Stack) Push(c Cell) {
	stack.data = append(stack.data, c)
	log.Printf("pushed: %#v", stack.data)
}

func (env *Env) Execute(tokens []Token) error {
	stack := env.stack

	for _, word := range tokens {
		switch word.tok {
		case PLUS:
			// operator: + x y
			// pop: y rhs
			// pop: x lhs
			rhs, err := stack.Pop()
			if err != nil {
				return err
			}

			lhs, err := stack.Pop()
			if err != nil {
				return err
			}

			y, ok := rhs.(*Int)
			if !ok {
				return fmt.Errorf("It's not Int: %#v", rhs)
			}

			x, ok := lhs.(*Int)
			if !ok {
				return fmt.Errorf("It's not Int: %#v", lhs)
			}

			stack.Push(NewInt(x.v + y.v))
		case MINUS:
			log.Fatal("TODO")
		case INT:
			v, err := strconv.ParseInt(word.lit, 10, 32) // REVIEW: bitは32でいいらしい
			if err != nil {
				log.Fatalf("unexpected int: %s", word.lit)
			}
			stack.Push(NewInt(int(v)))
		case PRINT:
			fmt.Print("[")
			for i, v := range stack.data {
				if i == 0 {
					fmt.Printf("%v", v)
				} else {
					fmt.Printf(" %v", v)
				}
			}
			fmt.Print("]")
			fmt.Println()
		}
	}

	return nil
}

func main() {
	log.SetPrefix("forth: ")
	log.SetFlags(0)

	tokens := Parse("100 200 300 400 + .s")
	log.Print(tokens)

	env := Env{stack: &Stack{}}
	if err := env.Execute(tokens); err != nil {
		log.Fatal(err)
	}

	fmt.Println(env.stack)
}
