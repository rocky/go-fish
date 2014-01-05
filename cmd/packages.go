// Copyright 2013-2014 Rocky Bernstein.

package fishcmd

import (
	"reflect"
	"github.com/rocky/go-fish"
)

func init() {
	name := "packages"
	repl.Cmds[name] = &repl.CmdInfo{
		Fn: PackageCommand,
		Help: `packages [*package* [*package* ...] ]

Show information about imported packages.

If a package name is given, then detailed information is given about
that package import. Otherwise we give a list of imported packages.
`,

		Min_args: 0,
		Max_args: -1,  // Max_args < 0 means an arbitrary number
	}
	repl.AddToCategory("support", name)
	repl.AddAlias("pkg", name)
	repl.AddAlias("pkgs", name)
	repl.AddAlias("package", name)
}

func printReflectMap(title string, m map[string] reflect.Value) {
	if len(m) > 0 {
		list := []string {}
		for item := range m {
			list = append(list, item)
		}
		repl.PrintSorted(title, list)
	}
}

func printReflectTypeMap(title string, m map[string] reflect.Type) {
	if len(m) > 0 {
		list := []string {}
		for item := range m {
			list = append(list, item)
		}
		repl.PrintSorted(title, list)
	}
}

// PackageCommand implements the command:
//    package [*name* [name*...]]
// which shows information about a package or lists all packages.
func PackageCommand(args []string) {
	if len(args) > 1 {
		for _, pkg_name := range args[1:len(args)] {
			if pkg, ok := repl.Env.Pkgs[pkg_name]; ok {
				repl.Section("=== Package %s (\"%s\"): ===", pkg_name, pkg.Path)
				printReflectMap("Constants of "+pkg_name, pkg.Consts)
				printReflectMap("Functions of "+pkg_name, pkg.Funcs)
				printReflectTypeMap("Types of "+pkg_name, pkg.Types)
				printReflectMap("Variables of "+pkg_name, pkg.Vars)
			} else {
			repl.Errmsg("Package %s not imported", pkg_name)
			}
		}
	} else {
		pkgNames := []string {}
		for pkg := range repl.Env.Pkgs {
			pkgNames = append(pkgNames, pkg)
		}
		repl.PrintSorted("All imported packages", pkgNames)
	}
}
