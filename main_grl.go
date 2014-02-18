// +build ignore
// Copyright 2013-2014 Rocky Bernstein

package main

// This simple REPL (read-eval-print loop) for Go using GNU Readline

import (
	"fmt"
	"os"
	"reflect"

	"code.google.com/p/go-gnureadline"
	"github.com/davecgh/go-spew/spew"
	"github.com/rocky/go-fish"
	"github.com/rocky/go-fish/cmd"
)

func intro_text() {
	repl.Section("== A Go eval REPL with GNU Readline support ==")
	fmt.Printf(`
Results of expression are stored in variable slice "results".
The environment is stored in global variable "env".
Short form assignment, e.g. a, b := 1, 2, is supported.

Enter expressions to be evaluated at the "gofish>" prompt.

To see all results, type: "results".

To quit, enter: "quit" or Ctrl-D (EOF).
To get help, enter: "help".
`)

}

// history_file is file name where history entries were and are to be saved. If
// the empty string, no history is saved and no history read in initially.
var historyFile string

// term is the current environment TERM value, e.g. "gnome", "xterm", or "vt100"
var term string

// gnuReadLineSetup is boilerplate initialization for GNU Readline.
func gnuReadLineSetup() {
	term = os.Getenv("TERM")
	historyFile = repl.HistoryFile(".go-fish")
	if historyFile != "" {
		gnureadline.ReadHistory(historyFile)
	}
	// Set maximum number of history entries
	gnureadline.StifleHistory(100)
}

// gnuReadLineTermination has GNU Readline Termination tasks:
// save history file if ane, and reset the terminal.
func gnuReadLineTermination() {
	if historyFile != "" {
		gnureadline.WriteHistory(historyFile)
	}
	if term != "" {
		gnureadline.Rl_reset_terminal(term)
	}
}

func spewInspect(a ...interface{}) string {
	value := a[0].(reflect.Value)
	return spew.Sdump(value.Interface())
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
	gnuReadLineSetup()

	defer gnuReadLineTermination()

	// Initialize REPL commands
	fishcmd.Init()

	repl.REPL(env, gnureadline.Readline, spewInspect)
	os.Exit(repl.ExitCode)
}
