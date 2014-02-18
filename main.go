// +build ignore
// Copyright 2013-2014 Rocky Bernstein

package main

// This is a simple REPL (read-eval-print loop) for Go.

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
	repl.Section("== A simple Go eval REPL ==")
	fmt.Printf(`
Results of expression are stored in variable slice "results".
The environment is stored in global variable "env".

Enter expressions to be evaluated at the "gofish>" prompt.

To see all results, type: "results".

To quit, enter: "quit" or Ctrl-D (EOF).
To get help, enter: "help".
`)

}

// Set up the Go package, function, constant, variable environment; then REPL
// (Read, Eval, Print, and Loop).
func main() {

	// A place to store result values of expressions entered
	// interactively
	env := repl.MakeEvalEnv()

	// Make this truly self-referential
	env.Vars["env"] = reflect.ValueOf(env)

	intro_text()

	repl.Input = bufio.NewReader(os.Stdin)

	// Initialize REPL commands
	fishcmd.Init()

	repl.REPL(env, repl.SimpleReadLine, repl.SimpleInspect)
	os.Exit(repl.ExitCode)
}
