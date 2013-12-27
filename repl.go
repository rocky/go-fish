package repl

// The common guts of is a simple REPL (read-eval-print loop) for GO.

// We separate this from the main package so that the main package
// can provide its own readline function. This could be, for example,
// GNU Readline, lineedit or something else.
import (
	"bufio"
	"fmt"
	"go/parser"
	"io"
	"strconv"
	"strings"

	"github.com/0xfaded/eval"
)

type ReadLineFnType func(prompt string, add_history ... bool) (string, error)
var  readLineFn ReadLineFnType

func SetReadLineFn(fn ReadLineFnType) {
	readLineFn = fn
}

func GetReadLineFn() ReadLineFnType {
	return readLineFn
}

// FIXME: The GNU readline interface, doesn't have an I/O parameter.
// Our SimpleReadLine needs a bufio.Reader. We'll use a global
// variable here to get around the interface mismatch problem.
var Input *bufio.Reader

// Simple replacement for GNU readline
func SimpleReadLine(prompt string, add_history ... bool) (string, error) {
	fmt.Printf(prompt)
	line, err := Input.ReadString('\n')
	if err == nil {
		line = strings.TrimRight(line, "\r\n")
	}
	return line, err
}

func init() {
	readLineFn = SimpleReadLine
}

// The read-eval-print portion
func REPL(env *interactive.Env, results *([]interface{})) {

	var err error
	exprs := 0
	line, err := readLineFn("go> ", true)
	for line != "quit" {
		if err != nil {
			if err == io.EOF { break }
			panic(err)
		}
		ctx := &interactive.Ctx{line}
		if expr, err := parser.ParseExpr(line); err != nil {
			fmt.Printf("parse error: %s\n", err)
		} else if cexpr, errs := interactive.CheckExpr(ctx, expr, env); len(errs) != 0 {
			for _, cerr := range errs {
				fmt.Printf("%v\n", cerr)
			}
		} else if vals, _, err := interactive.EvalExpr(ctx, cexpr, env); err != nil {
			fmt.Printf("eval error: %s\n", err)
		} else if vals == nil {
			fmt.Printf("nil\n")
		} else if len(*vals) == 0 {
			fmt.Printf("void\n")
		} else if len(*vals) == 1 {
			value := (*vals)[0]
			kind := value.Kind().String()
			typ  := value.Type().String()
			if typ != kind {
				fmt.Printf("Kind = %v\n", kind)
				fmt.Printf("Type = %v\n", typ)
			} else {
				fmt.Printf("Kind = Type = %v\n", kind)
			}
			if kind == "string" {
				fmt.Printf("results[%d] = %s\n", exprs,
					strconv.QuoteToASCII(value.String()))
			} else {
				fmt.Printf("results[%d] = %v\n", exprs, (value.Interface()))
			}
			exprs  += 1
			*results = append(*results, (*vals)[0].Interface())
		} else {
			fmt.Printf("(")
			size := len(*vals)
			for i, v := range *vals {
				fmt.Printf("%v", v.Interface())
				if i < size-1 { fmt.Printf(", ") }
			}
			fmt.Printf(")\n")
			exprs  += 1
			*results = append(*results, (*vals))
		}

		line, err = readLineFn("go> ", true)
	}
}
