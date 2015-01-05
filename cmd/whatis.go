// Copyright 2013-2015 Rocky Bernstein.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// whatis command

package fishcmd

import (
	"strings"
	"go/parser"
	"github.com/rocky/go-fish"
	"github.com/rocky/eval"
)

func init() {
	name := "whatis"
	repl.Cmds[name] = &repl.CmdInfo{
		Fn: WhatisCommand,
		Help: `whatis expression

Shows the type checker information for an expression. As a special
cases
   if expression is a package name, we'll confirm that.
   if the expression is a type, we'll show reflect.Kind information
 `,

		Min_args: 0,
		Max_args: -1,
	}
	repl.AddToCategory("data", name)
}

func WhatisCommand(args []string) {
	if len(args) == 2 {
		arg := args[1]
		if _, ok := repl.Env.Pkg(arg).(*eval.SimpleEnv); ok  {
			repl.Msg("`%s' is a package", arg)
			return
		}
		ids := strings.Split(arg, ".")
		if len(ids) == 1 {
			name := ids[0]
			if typ := repl.Env.Type(name); typ != nil  {
				repl.Msg("%s is a type: %s", typ.String())
				return
			}
		}
		if len(ids) == 2 {
			pkgName  := ids[0]
			name     := ids[1]
			if pkg, ok := repl.Env.Pkg(pkgName).(*eval.SimpleEnv); ok  {
				if typ := pkg.Type(name); typ != nil  {
					repl.Msg("%s is a kind: %s", arg, typ.Kind())
					repl.Msg("%s is a type: %v", arg, typ)
					return
				}
			}
		}
	}
	line := repl.CmdLine[len(args[0]):len(repl.CmdLine)]
	if expr, err := parser.ParseExpr(line); err != nil {
		if pair := eval.FormatErrorPos(line, err.Error()); len(pair) == 2 {
			repl.Msg(pair[0])
			repl.Msg(pair[1])
		}
		repl.Errmsg("parse error: %s\n", err)
	} else if cexpr, errs := eval.CheckExpr(expr, repl.Env); len(errs) != 0 {
		for _, cerr := range errs {
			repl.Msg("%v", cerr)
		}
	} else {
		repl.Section(cexpr.String())
		if cexpr.IsConst() {
			repl.Msg("constant:\t%s", cexpr.Const())
		}
		knownTypes := cexpr.KnownType()
		if len(knownTypes) == 1{
			repl.Msg("type:\t%s", knownTypes[0])
		} else {
			for i, v := range knownTypes {
				repl.Msg("type[%d]:\t%s", i, v)
			}
		}
	}
}
