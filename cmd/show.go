// Copyright 2014 Rocky Bernstein.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// show command

package fishcmd

import (
	"github.com/rocky/go-fish"
)

func init() {
	name := "show"
	repl.Cmds[name] = &repl.CmdInfo{
		SubcmdMgr: &repl.SubcmdMgr{
			Name   : name,
			Subcmds: make(repl.SubcmdMap),
		},
		Fn: ShowCommand,
		Help: `Generic command for showing things about the debugger.

Type "set" for a list of "set" subcommands and what they do.
Type "help set *" for just a list of "info" subcommands.`,
		Min_args: 0,
		Max_args: 3,
	}
	repl.AddToCategory("support", name)
}

func init() {
	name := "show"
	repl.Cmds[name] = &repl.CmdInfo{
		SubcmdMgr: &repl.SubcmdMgr{
			Name   : name,
			Subcmds: make(repl.SubcmdMap),
		},
		Fn: ShowCommand,
		Help: `Show parts of the debugger environment.

Type "show" for a list of "show" subcommands and what they do.
Type "help show *" for just a list of "show" subcommands.`,
		Min_args: 0,
		Max_args: 3,
	}
	repl.AddToCategory("support", name)
}

func ShowOnOff(subcmdName string, on bool) {
	if on {
		repl.Msg("%s is on.", subcmdName)
	} else {
		repl.Msg("%s is off.", subcmdName)
	}
}

// show implements the debugger command:
//    show [*subcommand]
// which is a generic command for setting things about the debugged program.
func ShowCommand(args []string) {
	repl.SubcmdMgrCommand(args)
}
