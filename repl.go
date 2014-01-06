// Copyright 2013-2014 Rocky Bernstein.

// Package repl is a simple REPL (read-eval-print loop) for GO using
// http://github.com/0xfaded/eval to the heavy lifting to implement
// the eval() part.
//
// Inside this package we provide two front-ends, one which uses GNU
// Readline (http://code.google.com/p/go-gnureadline) and one which doesn't.
// Feel free to add patches to support other kinds of readline support.
//
package repl

// We separate this from the main package so that the main package
// can provide its own readline function. This could be, for example,
// GNU Readline, lineedit or something else.
import (
	"bufio"
	"flag"
	"fmt"
	"go/parser"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/0xfaded/eval"
)

var Highlight = flag.Bool("highlight", true, `use syntax highlighting in output`)

// Maxwidth is the size of the line. We will try to wrap text that is
// longer than this. It like the COLUMNS environment variable
var Maxwidth int

// ReadLineFnType is function signature for a common read line
// interface that we support.
type ReadLineFnType func(prompt string, add_history ... bool) (string, error)
var  readLineFn ReadLineFnType

var initial_cwd string

// GOFISH_RESTART_CMD is a string that was used to invoke gofish.
//If we want to restart gofish, this is what we'll use.
var GOFISH_RESTART_CMD string


// HistoryFile returns a string file name to use for saving command
// history entries.
func HistoryFile(history_basename string) string {
	home_dir := os.Getenv("HOME")
	if home_dir == "" {
		// FIXME: also try ~ ?
		fmt.Println("ignoring history file; environment variable HOME not set")
		return ""
	}
	history_file := filepath.Join(home_dir, history_basename)
	if fi, err := os.Stat(history_file); err != nil {
		fmt.Println("No history file found to read in: ", err.Error())
	} else {
		if fi.IsDir() {
			Errmsg("Ignoring history file %s; is a directory, should be a file",
				history_file)
			return ""
		}
	}
	return history_file
}

// SetReadLineFn is used to set a specific readline function to be used
// as the "read" part of the read/eval/print loop.
func SetReadLineFn(fn ReadLineFnType) {
	readLineFn = fn
}

// GetReadLineFn returns the current readline function in effect for
// the "read" part of the read/eval/print loop.
func GetReadLineFn() ReadLineFnType {
	return readLineFn
}

// Input is a workaround for the fact that ReadLineFnType doesn't have
// an input parameter, but SimpleReadLine below needs a
// *bufioReader. So set this global variable beforehand if you are using
// SimpleReadLine.
var Input *bufio.Reader

// SimpleReadLine is simple replacement for GNU readline.
// prompt is the command prompt to print before reading input.
// add_history is ignored, but provided as a parameter to match
// those readline interfaces that do support saving command history.
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
	widthstr := os.Getenv("COLUMNS")
	initial_cwd, _ = os.Getwd()
	GOFISH_RESTART_CMD = os.Getenv("GOFISH_RESTART_CMD")
	if len(widthstr) == 0 {
		Maxwidth = 80
	} else if i, err := strconv.Atoi(widthstr); err == nil {
		Maxwidth = i
	}
}

// MakeEvalEnv creates an environment to use in evaluation.  The
// environment is exactly that environment needed by eval
// automatically extracted from the package eval
// (http://github.com/0xfaded/eval).
func MakeEvalEnv() eval.Env {
	var pkgs map[string] eval.Pkg = make(map[string] eval.Pkg)
	EvalEnvironment(pkgs)

	env := eval.Env {
		Name:   ".",
		Vars:   make(map[string] reflect.Value),
		Consts: make(map[string] reflect.Value),
		Funcs:  make(map[string] reflect.Value),
		Types:  make(map[string] reflect.Type),
		Pkgs:   pkgs,
	}
	return env
}

// LeaveREPL is set when we want to quit.
var LeaveREPL bool = false

// ExitCode is the exit code this program will set on exit.
var ExitCode  int  = 0

// Env is the evaluation environment we are working with.
var Env *eval.Env

// REPL is the read, eval, and print loop.
func REPL(env *eval.Env) {

	var err error

	// A place to store result values of expressions entered
	// interactively
	results := make([] interface{}, 0, 10)
	env.Vars["results"] = reflect.ValueOf(&results)

	Env = env
	exprs := 0
	line, err := readLineFn("gofish> ", true)
	for true {
		if err != nil {
			if err == io.EOF { break }
			panic(err)
		}
		if wasProcessed(line) {
			if LeaveREPL {break}
			line, err = readLineFn("gofish> ", true)
			continue
		}
		ctx := &eval.Ctx{line}
		if expr, err := parser.ParseExpr(line); err != nil {
			if pair := eval.FormatErrorPos(line, err.Error()); len(pair) == 2 {
				Msg(pair[0])
				Msg(pair[1])
			}
			Errmsg("parse error: %s", err)
		} else if cexpr, errs := eval.CheckExpr(ctx, expr, env); len(errs) != 0 {
			for _, cerr := range errs {
				Errmsg("%v", cerr)
			}
		} else if vals, _, err := eval.EvalExpr(ctx, cexpr, env); err != nil {
			Errmsg("eval error: %s", err)
		} else if vals == nil {
			Msg("Kind=nil\nnil")
		} else if len(*vals) == 0 {
			Msg("Kind=Slice\nvoid")
		} else if len(*vals) == 1 {
			value := (*vals)[0]
			if value.IsValid() {
				kind := value.Kind().String()
				typ  := value.Type().String()
				if typ != kind {
					Msg("Kind = %v", kind)
					Msg("Type = %v", typ)
				} else {
					Msg("Kind = Type = %v", kind)
				}
				Msg("results[%d] = %s", exprs, eval.Inspect(value))
				exprs += 1
				results = append(results, (*vals)[0].Interface())
			} else {
				Msg("%s", value)
			}
		} else {
			Msg("Kind = Multi-Value")
			size := len(*vals)
			for i, v := range *vals {
				fmt.Printf("%s", eval.Inspect(v))
				if i < size-1 { fmt.Printf(", ") }
			}
			Msg("")
			exprs += 1
			results = append(results, (*vals))
		}

		line, err = readLineFn("gofish> ", true)
	}
}
