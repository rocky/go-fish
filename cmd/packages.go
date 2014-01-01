// Copyright 2013-2014 Rocky Bernstein.

package fishcmd

import (
	"sort"
	"strings"
	"code.google.com/p/go-columnize"
	"github.com/rocky/go-fish"
)

func init() {
	name := "packages"
	repl.Cmds[name] = &repl.CmdInfo{
		Fn: PackageCommand,
		Help: `packages [*package* ]

Show information about imported packages.

If a package name is given, then detailed information is given about
that package import. Otherwise we give a list of imported packages.
`,

		Min_args: 0,
		Max_args: 1,
	}
	repl.AddToCategory("support", name)
	repl.AddAlias("pkg", name)
	repl.AddAlias("pkgs", name)
	repl.AddAlias("package", name)
}

func PackageCommand(args []string) {
	repl.Section("All imported packages:")
	opts := columnize.DefaultOptions()
	opts.DisplayWidth = repl.Maxwidth
	pkgNames := []string {}
	if len(args) > 1 {
		repl.Errmsg("Sorry, information about single package not done yet")
	} else {
		for pkg := range repl.Env.Pkgs {
			pkgNames = append(pkgNames, pkg)
		}
		sort.Strings(pkgNames)
		columnizedNames := strings.TrimRight(columnize.Columnize(pkgNames, opts),
			"\n")
		repl.Msg(columnizedNames)
	}
}
