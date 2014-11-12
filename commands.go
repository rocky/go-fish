// Copyright 2013-2014 Rocky Bernstein.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repl

type CmdFunc func([]string)

type CmdInfo struct {
	Help string
	Category string
	Min_args int
	Max_args int
	Fn CmdFunc
	Aliases []string
	SubcmdMgr *SubcmdMgr
}

// Cmds contains a list of the top-level REPL commands we implement.
// For example, "quit", and "help" are REPL commands.
var Cmds map[string]*CmdInfo  = make(map[string]*CmdInfo)


// Aliases maps a name to its underlying gofish command name.
// For example, "?" is an alias for "help".
var	Aliases map[string]string = make(map[string]string)

// Categories maps a REPL category name into the list of
// REPL commands in that category.
var	Categories map[string] []string = make(map[string] []string)

// AddAlias adds "alias" for a command name "cmdname"
func AddAlias(alias string, cmdname string) bool {
	if unalias := Aliases[alias]; unalias != "" {
		return false
	}
	Aliases[alias] = cmdname
	Cmds[cmdname].Aliases = append(Cmds[cmdname].Aliases, alias)
	return true
}

// AddToCategory adds "cmdname" into general category "category".
func AddToCategory(category string, cmdname string) {
	Categories[category] = append(Categories[category], cmdname)
	// Cmds[cmdname].category = category
}


// LookupCmd canonicalize parameter cmd, by changing it to the underlying
// gofish command if it is an alias.
func LookupCmd(cmd string) (string) {
	if Cmds[cmd] == nil {
		cmd = Aliases[cmd];
	}
	return cmd
}
