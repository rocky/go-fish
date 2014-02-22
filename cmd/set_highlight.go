// Copyright 2013 Rocky Bernstein.

// set highlight - use terminal highlight?

package fishcmd

import (
	"github.com/rocky/go-fish"
)

func init() {
	parent := "set"
	repl.AddSubCommand(parent, &repl.SubcmdInfo{
		Fn: SetHighlightSubcmd,
		Help: `set highlight [on|off]

Sets whether terminal highlighting is to be used`,
		Min_args: 0,
		Max_args: 1,
		Short_help: "use terminal highlight",
		Name: "highlight",
	})
}

func SetHighlightSubcmd(args []string) {
	onoff := "on"
	if len(args) == 3 {
		onoff = args[2]
	}
	switch ParseOnOff(onoff) {
	case ONOFF_ON:
		if *repl.Highlight {
			repl.Errmsg("Highlight is already on")
		} else {
			repl.Msg("Setting highlight on")
			*repl.Highlight = true
		}
	case ONOFF_OFF:
		if !*repl.Highlight {
			repl.Errmsg("highight is already off")
		} else {
			repl.Msg("Setting highlight off")
			*repl.Highlight = false
		}
	case ONOFF_UNKNOWN:
		repl.Msg("Expecting 'on' or 'off', got '%s'; nothing done", onoff)
	}
}
