// starting import: "github.com/rocky/go-fish"
package repl

import (
	"bufio"
	"github.com/0xfaded/eval"
	"reflect"
)

// EvalEnvironment adds to eval.Pkg those packages included
// with import "github.com/rocky/go-fish".

func EvalEnvironment() *eval.SimpleEnv {
	var consts map[string] reflect.Value
	var vars   map[string] reflect.Value
	var types  map[string] reflect.Type
	var funcs  map[string] reflect.Value
	var pkgs   map[string] eval.Env = make(map[string] eval.Env)

	consts = make(map[string] reflect.Value)
	consts["MaxScanTokenSize"] = reflect.ValueOf(bufio.MaxScanTokenSize)

	funcs = make(map[string] reflect.Value)
	funcs["NewReaderSize"] = reflect.ValueOf(bufio.NewReaderSize)
	funcs["NewReader"] = reflect.ValueOf(bufio.NewReader)
	funcs["NewWriterSize"] = reflect.ValueOf(bufio.NewWriterSize)
	funcs["NewWriter"] = reflect.ValueOf(bufio.NewWriter)
	funcs["NewReadWriter"] = reflect.ValueOf(bufio.NewReadWriter)
	funcs["NewScanner"] = reflect.ValueOf(bufio.NewScanner)
	funcs["ScanBytes"] = reflect.ValueOf(bufio.ScanBytes)
	funcs["ScanRunes"] = reflect.ValueOf(bufio.ScanRunes)
	funcs["ScanLines"] = reflect.ValueOf(bufio.ScanLines)
	funcs["ScanWords"] = reflect.ValueOf(bufio.ScanWords)

	types = make(map[string] reflect.Type)
	types["Reader"] = reflect.TypeOf(*new(bufio.Reader))
	types["Writer"] = reflect.TypeOf(*new(bufio.Writer))
	types["ReadWriter"] = reflect.TypeOf(*new(bufio.ReadWriter))
	types["Scanner"] = reflect.TypeOf(*new(bufio.Scanner))
	types["SplitFunc"] = reflect.TypeOf(*new(bufio.SplitFunc))

	vars = make(map[string] reflect.Value)
	vars["ErrInvalidUnreadByte"] = reflect.ValueOf(&bufio.ErrInvalidUnreadByte)
	vars["ErrInvalidUnreadRune"] = reflect.ValueOf(&bufio.ErrInvalidUnreadRune)
	vars["ErrBufferFull"] = reflect.ValueOf(&bufio.ErrBufferFull)
	vars["ErrNegativeCount"] = reflect.ValueOf(&bufio.ErrNegativeCount)
	vars["ErrTooLong"] = reflect.ValueOf(&bufio.ErrTooLong)
	vars["ErrNegativeAdvance"] = reflect.ValueOf(&bufio.ErrNegativeAdvance)
	vars["ErrAdvanceTooFar"] = reflect.ValueOf(&bufio.ErrAdvanceTooFar)
	pkgs["bufio"] = &eval.SimpleEnv {
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "bufio",
	}

	mainEnv := eval.MakeSimpleEnv()
	mainEnv.Pkgs = pkgs
	return mainEnv
}
