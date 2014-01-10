// Copyright 2014 Rocky Bernstein.

package fishcmd

import (
	"reflect"
	"strings"
	"github.com/rocky/go-fish"
)

func init() {
	name := "method"
	repl.Cmds[name] = &repl.CmdInfo{
		Fn: MethodCommand,
		Help: `method [*typeorvalue* [*typeorvalue* ...] ]

Show information about methods of a type or value.

If a type name is given, then information is given about
the methods of that type. If a value is given methods of that value are given
`,

		Min_args: 0,
		Max_args: -1,  // Max_args < 0 means an arbitrary number
	}
	repl.AddToCategory("support", name)
	repl.AddAlias("fn", name)
	repl.AddAlias("func", name)
}

func printMethodsOf(fullname string) {
	pkgName  := "."
	name     :=  fullname
	names   := strings.Split(fullname, ".")
	if len(names) > 1 {
		pkgName = names[0]
		name    = names[1]
		}
	pkgs := repl.Env.Pkgs
	pkg, ok  := pkgs[pkgName]
	if !ok || pkg == nil {
		repl.Errmsg("Can't find package %s", pkgName)
		return
	}
	if v, ok := pkg.Vars[name]; ok {
		methods := map[string] reflect.Value {}
		for i:= 0; i < v.Type().NumMethod(); i++ {
			meth := v.Method(i)
			name := meth.Type().Name()
			if name != "" {
				methods[name] = v
			}
		}
		if len(methods) == 0 {
			repl.Msg("No methods found for variable %s", fullname)
		} else {
			printReflectMap("Methods for variable " + fullname, methods)
		}
	} else if v, ok := pkg.Types[name]; ok {
		if v == nil {
			repl.Errmsg("Don't have method info recorded for type %s", fullname)
			return
		}
		methods := map[string] reflect.Type {}
		for i:= 0; i < v.NumMethod(); i++ {
			meth := v.Method(i)
			name := meth.Name
			methods[name] = v
		}
		if len(methods) == 0 {
			repl.Msg("No methods found for type %s", fullname)
		} else {
			printReflectTypeMap("Methods for type " + fullname, methods)
		}
	} else {
		repl.Errmsg("Can't find member %s in package %s", name, pkgName)
	}
}

// MethodCommand implements the command:
//    method [*name* [name*...]]
// which shows information about a package or lists all packages.
func MethodCommand(args []string) {
	if len(args) > 2 {
		for _, name := range args[1:len(args)] {
			printMethodsOf(name)
		}
	} else {
		printMethodsOf(args[1])
	}
}
