package main

import "log"
import "regexp"
import "fmt"
import "strconv"

type ForthInt int32

const intBitSize = 32

// Token
type Token struct {
	lit string
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%v}", t.lit)
}

// Parse
var tokenizer *regexp.Regexp

func init() {
	tokenizer = regexp.MustCompile(`\s*([\w.:;]+|[+-])`)
}

func Parse(code string) []Token {
	log.Printf("tokenize input code = %#v", code)

	tokens := make([]Token, 0)

	for _, group := range tokenizer.FindAllStringSubmatch(code, -1) {
		lit := group[1]
		token := Token{lit: lit}
		tokens = append(tokens, token)
	}

	return tokens
}

// Cell
type Cell interface {
	String() string
}

// Int
type Int struct {
	v ForthInt
}

func NewInt(v ForthInt) *Int {
	return &Int{v: v}
}

func (i Int) String() string {
	return fmt.Sprint(i.v)
}

type BuiltinProc func(*Env) error

// Proc
type Proc struct {
	fn BuiltinProc
}

func NewProc(fn func(*Env) error) *Proc {
	return &Proc{fn: fn}
}

func NewArrayProc(array []*Proc) *Proc {
	return &Proc{fn: func(env *Env) error {
		for _, proc := range array {
			if err := proc.Invoke(env); err != nil {
				return err
			}
		}
		return nil
	}}
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

// Buintin Procs
func printStack(env *Env) error {
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
}

func addTwoInt(env *Env) error {
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
	env.dictionary.Add(".s", NewProc(printStack))
	env.dictionary.Add("+", NewProc(addTwoInt))
	return env
}

func parseInt(token *Token) (ForthInt, bool) {
	v, err := strconv.ParseInt(token.lit, 10, intBitSize)
	if err != nil {
		return 0, false
	}
	return ForthInt(v), true
}

func (env *Env) compileWord(literal string) (*Proc, error) {
	cell, ok := env.dictionary.Get(literal)
	if !ok {
		return nil, fmt.Errorf("compile! undefined word: %v", literal)
	}
	proc, ok := cell.(*Proc)
	if !ok {
		return nil, fmt.Errorf("compile! it's not proc: %v", literal)
	}
	return proc, nil
}

func (env *Env) Compile(tokens []Token, pos int) (int, error) {
	name := tokens[pos].lit
	code := make([]*Proc, 0)

	for i := pos + 1; i < len(tokens); i++ {
		token := tokens[i]

		if token.lit == ";" {
			env.dictionary.Add(name, NewArrayProc(code))
			return i, nil
		}

		if v, ok := parseInt(&token); ok {
			code = append(code, NewProc(func(env *Env) error {
				env.stack.Push(NewInt(v))
				return nil
			}))
			continue
		}

		proc, err := env.compileWord(token.lit)
		if err != nil {
			return 0, err
		}
		code = append(code, proc)
	}

	return len(tokens), fmt.Errorf("compile! end of tokens")
}

func (env *Env) Execute(tokens []Token) error {
	stack := env.stack
	dictionary := env.dictionary

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if token.lit == ":" {
			pos, err := env.Compile(tokens, i+1)
			if err != nil {
				return err
			}
			i = pos
			continue
		}

		if v, ok := parseInt(&token); ok {
			stack.Push(NewInt(ForthInt(v)))
			continue
		}

		cell, ok := dictionary.Get(token.lit)
		if !ok {
			log.Fatalf("undefined word: %v", token.lit)
		}

		proc, ok := cell.(*Proc)
		if !ok {
			log.Fatalf("it's not word: %v", token.lit)
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
