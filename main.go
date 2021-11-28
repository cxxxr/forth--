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

// Cell
type Cell interface {
	String() string
}

// Int
type Int struct {
	v int
}

func NewInt(v int) *Int {
	return &Int{v: v}
}

func (i Int) String() string {
	return fmt.Sprint(i.v)
}

// Proc
type Proc struct {
	fn func(*Env) error
}

func NewProc(fn func(*Env) error) *Proc {
	return &Proc{fn: fn}
}

func (p Proc) String() string {
	return fmt.Sprintf("[Proc %p]", p.fn)
}

func (p *Proc) Invoke(env *Env) error {
	return p.fn(env)
}

// Stack
type Stack struct {
	data []Cell
}

func NewStack() *Stack {
	return &Stack{data: make([]Cell, 0)}
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

// Dictionary
type Dictionary struct {
	data map[string]Cell
}

func NewDictionary() *Dictionary {
	return &Dictionary{data: make(map[string]Cell)}
}

func (dict *Dictionary) Add(name string, cell Cell) {
	dict.data[name] = cell
}

func (dict *Dictionary) Get(name string) (Cell, bool) {
	cell, ok := dict.data[name]
	return cell, ok
}

// Env
type Env struct {
	stack      *Stack
	dictionary *Dictionary
}

func NewEnv() *Env {
	env := new(Env)
	env.stack = NewStack()
	env.dictionary = NewDictionary()
	env.dictionary.Add(".s", NewProc(func(env *Env) error {
		fmt.Print("[")
		for i, v := range env.stack.data {
			if i == 0 {
				fmt.Printf("%v", v)
			} else {
				fmt.Printf(" %v", v)
			}
		}
		fmt.Print("]")
		fmt.Println()

		return nil
	}))
	env.dictionary.Add("+", NewProc(func(env *Env) error {
		stack := env.stack
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

		return nil
	}))
	return env
}

func parseInt(word *Token) (int, bool) {
	v, err := strconv.ParseInt(word.lit, 10, 32) // REVIEW: bitは32でいいらしい
	if err != nil {
		return 0, false
	}
	return int(v), true
}

func (env *Env) Execute(tokens []Token) error {
	stack := env.stack
	dictionary := env.dictionary

	for _, word := range tokens {

		if v, ok := parseInt(&word); ok {
			stack.Push(NewInt(int(v)))
			continue
		}

		cell, ok := dictionary.Get(word.lit)
		if !ok {
			log.Fatalf("undefined word: %v", word.lit)
		}

		proc, ok := cell.(*Proc)
		if !ok {
			log.Fatalf("it's not word: %v", word.lit)
		}

		if err := proc.Invoke(env); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func main() {
	log.SetPrefix("forth: ")
	log.SetFlags(0)

	tokens := Parse("100 200 300 400 + .s")
	log.Print(tokens)

	env := NewEnv()
	if err := env.Execute(tokens); err != nil {
		log.Fatal(err)
	}
}
