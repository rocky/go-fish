// +build ignore

package main

// This is a simple REPL (read-eval-print loop) for GO.

// An intent (rocky) was also to have something that I can debug in
// the ssa-debugger tortoise/gub.sh. Right now that can't handle the
// unsafe package, pointers, and calls to C code. So that let's out
// go-gnureadline and lineedit.
import (
	"bufio"
	"fmt"
	"go/parser"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/rocky/go-repl"
	"github.com/0xfaded/go-interactive"
)

// Simple replacement for GNU readline
func readline(prompt string, in *bufio.Reader) (string, error) {
	fmt.Printf(prompt)
	line, err := in.ReadString('\n')
	if err == nil {
		line = strings.TrimRight(line, "\r\n")
	}
	return line, err
}

// The read-eval-print portion
func REPL(env *interactive.Env, results *([]interface{})) {

	var err error
	exprs := 0
	in := bufio.NewReader(os.Stdin)
	line, err := readline("go> ", in)
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
			fmt.Printf("Kind = %v\n", kind)
			typ  := value.Type().String()
			if typ != kind { fmt.Printf("Type = %v\n", typ) }
			if kind == "string" {
				fmt.Printf("results[%d] = %s\n", exprs,
					strconv.QuoteToASCII(value.String()))
			} else {
				fmt.Printf("results[%d] = %v\n", exprs, (value.Interface()))
			}
			exprs  += 1
			*results = append(*results, (*vals)[0].Interface())
		} else {
			sep := "("
			for _, v := range *vals {
				fmt.Printf("%s%v", sep, v.Interface())
			}
			fmt.Printf(")\n")
		}

		line, err = readline("go> ", in)
	}
}

func main() {
	// Set up the environment and then call REPL
	var vars   map[string] reflect.Value = make(map[string] reflect.Value)
	var consts map[string] reflect.Value = make(map[string] reflect.Value)
	var types  map[string] reflect.Type  = make(map[string] reflect.Type)

	var global_funcs map[string] reflect.Value = make(map[string] reflect.Value)
	var global_vars map[string]  reflect.Value = make(map[string] reflect.Value)

	// A place to store result values of expressions entered
	// interactively
	var results []interface{} = make([] interface{}, 0, 10)
	global_vars["results"] = reflect.ValueOf(&results)
	global_vars["arg0"] = reflect.ValueOf(os.Args[0])
	// global_funcs["Result"] = reflect.ValueOf(
	// 	func(i int) interface{} { return results[i] } )

	// What we have from the fmt package.
	var fmt_funcs    map[string] reflect.Value = make(map[string] reflect.Value)
	fmt_funcs["Println"] = reflect.ValueOf(fmt.Println)
	fmt_funcs["Printf"] = reflect.ValueOf(fmt.Printf)

	type Alice struct {
		Bob int
		Secret string
	}

	pkgs := map[string] interactive.Pkg {
			"fmt": &interactive.Env {
				Name:   "fmt",
				Vars:   vars,
				Consts: consts,
				Funcs:  fmt_funcs,
				Types:  types,
				Pkgs:   make(map[string] interactive.Pkg),
			}, "os": &interactive.Env {
				Name:   "os",
				Vars:   map[string] reflect.Value { "Stdout": reflect.ValueOf(&os.Stdout) },
				Consts: make(map[string] reflect.Value),
				Funcs:  make(map[string] reflect.Value),
				Types:  make(map[string] reflect.Type),
				Pkgs:   make(map[string] interactive.Pkg),
			},
		}


	repl.Extract_environment(pkgs)

	env := interactive.Env {
		Name:   ".",
		Vars:   global_vars,
		Consts: make(map[string] reflect.Value),
		Funcs:  global_funcs,
		Types:  map[string] reflect.Type{ "Alice": reflect.TypeOf(Alice{}) },
		Pkgs:   pkgs,
	}

	// Make this truly self-referential
	global_vars["env"] = reflect.ValueOf(&env)


	fmt.Printf(`=== A simple Go eval REPL ===

Results of expression are stored in variable slice "results".
Defined functions are: fmt.Println(), fmt.Printf().

The environment is stored in global variable "env".

Enter expressions to be evaluated at the "go>" prompt.

To see all results, type: "results".

To quit, enter: "quit" or Ctrl-D (EOF).
`)
	// And just when you thought we'd never get around to it...
	REPL(&env, &results)
}
