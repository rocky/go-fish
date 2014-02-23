// Copyright 2014 Rocky Bernstein.

// set width - set width of line

package fishcmd

import (
	"github.com/rocky/go-fish"
)

func init() {
	parent := "set"
	repl.AddSubCommand(parent, &repl.SubcmdInfo{
		Fn: SetWidthSubcmd,
		Help: `set width num

Sets the line length the REPL thinks we have`,
		Min_args: 1,
		Max_args: 1,
		Short_help: "set line width",
		Name: "width",
	})
}

func SetWidthSubcmd(args []string) {
	i, err := repl.GetInt(args[2], "line width", 0, 10000)
	if err != nil { return }
	repl.Maxwidth = i
	ShowWidthSubcmd(args)
}
