// Copyright 2013 Rocky Bernstein.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// quit command

package fishcmd

import (
	"strconv"
	"github.com/rocky/go-fish"
)

func init() {
	name := "quit"
	repl.Cmds[name] = &repl.CmdInfo{
		Fn: QuitCommand,
		Help: `quit [exit-code]

Terminates program. If an exit code is given, that is the exit code
for the program. Zero (normal termination) is used if no
termintation code.
`,

		Min_args: 0,
		Max_args: 1,
	}
	repl.AddToCategory("support", name)
	repl.Aliases["q"] = name
}

func QuitCommand(args []string) {
	rc := 0
	if len(args) == 2 {
		new_rc, ok := strconv.Atoi(args[1])
		if ok == nil { rc = new_rc } else {
			repl.Errmsg("Expecting integer return code; got %s.",
				args[1])
			return
		}
	}
	repl.Msg("go-fish: That's all folks...")

	repl.LeaveREPL = true
	repl.ExitCode = rc
}
