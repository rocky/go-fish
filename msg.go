// Copyright 2013-2014 Rocky Bernstein.

package repl

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mgutz/ansi"
	"code.google.com/p/go-columnize"
)

var	termReset, termBold, termHighlight string

func init() {
	termReset     = ansi.ColorCode("reset")
	termBold      = ansi.ColorCode("+b")
	termHighlight = ansi.ColorCode("+h")
}

func Errmsg(format string, a ...interface{}) (n int, err error) {
	if *Highlight {
		format = termHighlight + format + termReset + "\n"
	} else {
		format = "** " + format + "\n"
	}
	return fmt.Fprintf(os.Stdout, format, a...)
}

func MsgNoCr(format string, a ...interface{}) (n int, err error) {
	format = format
	return fmt.Fprintf(os.Stdout, format, a...)
}

func Msg(format string, a ...interface{}) (n int, err error) {
	format = format + "\n"
	return fmt.Fprintf(os.Stdout, format, a...)
}

// A more emphasized version of msg. For section headings.
func Section(format string, a ...interface{}) (n int, err error) {
	if *Highlight {
		format = termBold + format + termReset + "\n"
	} else {
		format = format + "\n"
	}
	return fmt.Fprintf(os.Stdout, format, a...)
}

func PrintSorted(title string, names []string) {
	Section(title + ":")
	sort.Strings(names)
	opts := columnize.DefaultOptions()
	opts.LinePrefix  = "  "
	opts.DisplayWidth = Maxwidth
	columnizedNames := strings.TrimRight(columnize.Columnize(names, opts),
		"\n")
	Msg(columnizedNames)

}
