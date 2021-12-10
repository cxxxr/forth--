package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cxxxr/forth--/forth"
)

type ForthInt int32

const intBitSize = 32

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

func parseInt(token *forth.Token) (ForthInt, bool) {
	v, err := strconv.ParseInt(token.Lit, 10, intBitSize)
	if err != nil {
		return 0, false
	}
	return ForthInt(v), true
}

func (env *Env) compileWord(literal string) (*Proc, error) {
	cell, ok := env.dictionary.Get(literal)
	if !ok {
		return nil, fmt.Errorf("undefined word: %v", literal)
	}
	proc, ok := cell.(*Proc)
	if !ok {
		return nil, fmt.Errorf("it's not proc: %v", literal)
	}
	return proc, nil
}

func (env *Env) Compile(tokens []forth.Token, start int) (*Proc, int, error) {
	code := make([]*Proc, 0)

	for i := start; i < len(tokens); i++ {
		token := tokens[i]

		if token.Lit == ":" {
			name := tokens[i+1].Lit
			compiled, pos, err := env.Compile(tokens, i+2)
			if err != nil {
				return nil, pos, nil
			}
			if pos == len(tokens) {
				return nil, pos, fmt.Errorf("end of tokens")
			}
			env.dictionary.Add(name, compiled)
			i = pos
			continue
		}

		if token.Lit == ";" {
			return NewArrayProc(code), i, nil
		}

		if v, ok := parseInt(&token); ok {
			code = append(code, NewProc(func(env *Env) error {
				env.stack.Push(NewInt(v))
				return nil
			}))
			continue
		}

		proc, err := env.compileWord(token.Lit)
		if err != nil {
			return nil, i, err
		}
		code = append(code, proc)
	}

	return NewArrayProc(code), len(tokens), nil
}

func (env *Env) Execute(tokens []forth.Token) error {
	compiled, _, err := env.Compile(tokens, 0)
	if err != nil {
		return err
	}
	return compiled.Invoke(env)
}

func prompt(scanner *bufio.Scanner) (string, bool) {
	fmt.Print("> ")
	if !scanner.Scan() {
		return "", false
	}
	return scanner.Text(), true
}

func interp() {
	scanner := bufio.NewScanner(os.Stdin)

	env := NewEnv()

	for {
		line, ok := prompt(scanner)
		if !ok {
			break
		}

		tokens := forth.Parse(line)
		log.Printf("tokens = %#v\n", tokens)

		if err := env.Execute(tokens); err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	log.SetPrefix("forth: ")
	log.SetFlags(0)

	interp()
}
