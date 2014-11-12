// Copyright 2013-2014 Rocky Bernstein.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// set command

package fishcmd

import (
	"github.com/rocky/go-fish"
)

func init() {
	name := "set"
	repl.Cmds[name] = &repl.CmdInfo{
		SubcmdMgr: &repl.SubcmdMgr{
			Name   : name,
			Subcmds: make(repl.SubcmdMap),
		},
		Fn: SetCommand,
		Help: `Modifies parts of the REPL environment.

Type "set" for a list of "set" subcommands and what they do.
Type "help set *" for just a list of "set" subcommands.
`,
		Min_args: 0,
		Max_args: 3,
	}
	repl.AddToCategory("support", name)
}


type onoff uint8
const (
	ONOFF_ON = iota
	ONOFF_OFF
	ONOFF_UNKNOWN
)

func ParseOnOff(onoff string) onoff {
	switch onoff {
	case "on", "1", "yes":
		return ONOFF_ON
	case "off", "0", "none":
		return ONOFF_OFF
	default:
		return ONOFF_UNKNOWN
	}
}

func init() {
	parent := "set"
	repl.AddSubCommand(parent, &repl.SubcmdInfo{
		Fn: SetHighlightSubcmd,
		Help: `Modifies parts of the debugger environment.

You can give unique prefix of the name of a subcommand to get
information about just that subcommand.

Type "set" for a list of "set" subcommands and what they do.`,
		Min_args: 0,
		Max_args: 3,
	})
}

// setCommand implements the debugger command:
//    set [*subcommand*]
// which modifies parts of the debugger environment.
func SetCommand(args []string) {
	repl.SubcmdMgrCommand(args)
}
