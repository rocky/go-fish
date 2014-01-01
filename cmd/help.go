// Copyright 2013-2014 Rocky Bernstein.

package fishcmd

import (
	"sort"
	"strings"
	"code.google.com/p/go-columnize"
	"github.com/rocky/go-fish"
)

func init() {
	name := "help"
	repl.Cmds[name] = &repl.CmdInfo{
		Fn: HelpCommand,
		Help: `help [*command* | * ]

To evaluate an expression, just type the expression.

If the first word of the line starts with a gofish command, then that
takes precendence. For example, "help" is a gofish command.

Typing "help *" will print a list of available gofish commands.

When "help and an argument is given, if it is '*' a list of repl
commands is shown. Otherwise the argument is checked to see if it is
command name. For example 'help quit' gives help on the 'quit'
debugger command.

`,

		Min_args: 0,
		Max_args: 2,
	}
	repl.AddToCategory("support", name)
	repl.AddAlias("?", name)
	// Down the line we'll have abbrevs
	repl.AddAlias("h", name)
}

// HelpCommand implements the command:
//    help [*name* |* ]
// which gives help.
func HelpCommand(args []string) {
	if len(args) == 1 {
		repl.Msg(repl.Cmds["help"].Help)
	} else {
		what := args[1]
		cmd := repl.LookupCmd(what)
		if what == "*" {
			var names []string
			for k, _ := range repl.Cmds {
				names = append(names, k)
			}
			repl.Section("All command names:")
			sort.Strings(names)
			opts := columnize.DefaultOptions()
			opts.LinePrefix  = "  "
			opts.DisplayWidth = repl.Maxwidth
			mems := strings.TrimRight(columnize.Columnize(names, opts),
				"\n")
			repl.Msg(mems)
		} else if what == "categories" {
			repl.Section("Categories")
			for k, _ := range repl.Categories {
				repl.Msg("\t %s", k)
			}
		} else if info := repl.Cmds[cmd]; info != nil {
			// if len(args) > 2 {
			// 	if info.SubcmdMgr != nil {
			// 		repl.HelpSubCommand(info.SubcmdMgr, args)
			// 		return
			// 	}
			// }
			repl.Msg(info.Help)
			if len(info.Aliases) > 0 {
				repl.Msg("Aliases: %s",
					strings.Join(info.Aliases, ", "))
			}
		} else if cmds := repl.Categories[what]; len(cmds) > 0 {
			repl.Section("Commands in class: %s", what)
			sort.Strings(cmds)
			opts := columnize.DefaultOptions()
			opts.DisplayWidth = repl.Maxwidth
			mems := strings.TrimRight(columnize.Columnize(cmds, opts),
				"\n")
			repl.Msg(mems)
		} else {
			repl.Errmsg("Can't find help for %s", what)
		}
	}
}
