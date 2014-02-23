// Copyright 2014 Rocky Bernstein.

// set width - set width of line

package fishcmd

import (
	"github.com/rocky/go-fish"
)

func init() {
	parent := "show"
	repl.AddSubCommand(parent, &repl.SubcmdInfo{
		Fn: ShowWidthSubcmd,
		Help: `show width

Show the line length the REPL thinks we have`,
		Min_args: 0,
		Max_args: 0,
		Short_help: "show line width",
		Name: "width",
	})
}

func ShowWidthSubcmd(args []string) {
	repl.Msg("Line width is %d", repl.Maxwidth)
}
