// +build ignore
// Copyright 2013-2014 Rocky Bernstein

package main

// This is a simple REPL (read-eval-print loop) for GO.

// (rocky) My intent here is to have something that I can debug in
// the ssa-debugger tortoise/gub.sh. Right now that can't handle the
// unsafe package, pointers, and calls to C code. So that let's out
// go-gnureadline and lineedit.
// See also main_gr.go for GNU readline code.
import (
	"bufio"
	"fmt"
	"os"
	"reflect"

	"github.com/rocky/go-fish"
	"github.com/rocky/go-fish/cmd"
)

func intro_text() {
	fmt.Printf(`=== A simple Go eval REPL ===

Results of expression are stored in variable slice "results".
The environment is stored in global variable "env".

Enter expressions to be evaluated at the "go>" prompt.

To see all results, type: "results".

To quit, enter: "quit" or Ctrl-D (EOF).
`)

}

func main() {
	// Set up the environment and then call REPL
	// A place to store result values of expressions entered
	// interactively
	var results []interface{} = make([] interface{}, 0, 10)
	var global_vars map[string]  reflect.Value = make(map[string] reflect.Value)
	global_vars["results"] = reflect.ValueOf(&results)

	env := repl.MakeEvalEnv()
	env.Vars = global_vars

	// Make this truly self-referential
	global_vars["env"] = reflect.ValueOf(&env)

	intro_text()

	repl.Input = bufio.NewReader(os.Stdin)

	// Initialize REPL commands
	fishcmd.Init()

	// And just when you thought we'd never get around to it...
	repl.REPL(&env, &results)
}
