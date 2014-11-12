// Copyright 2013-2014 Rocky Bernstein.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// command argument-validation routines

package repl

import "strconv"

func ArgCountOK(min int, max int, args [] string) bool {
	l := len(args)-1 // strip command name from count
	if l < min {
		Errmsg("Too few args; need at least %d, got %d", min, l)
		return false
	} else if max > 0 && l > max {
		Errmsg("Too many args; need at most %d, got %d", max, l)
		return false
	}
	return true
}

type NumError struct {
	bogus bool
}

func (e *NumError) Error() string {
	return "generic error"
}
var genericError = &NumError{bogus: true}

func GetInt(arg string, what string, min int, max int) (int, error) {
	errmsg_fmt := "Expecting integer " + what + "; got '%s'."
	i, err := strconv.Atoi(arg)
	if err != nil {
		Errmsg(errmsg_fmt, arg)
		return 0, err
	}
	if i < min {
		Errmsg("Expecting integer value %s to be at least %d; got %d.",
			what, min, i)
        return 0, genericError
	} else if max > 0 && i > max {
        Errmsg("Expecting integer value %s to be at most %d; got %d.",
			what, max, i)
        return 0, genericError
	}
	return i, nil
}


func GetUInt(arg string, what string, min uint64, max uint64) (uint64, error) {
	errmsg_fmt := "Expecting integer " + what + "; got '%s'."
	i, err := strconv.ParseUint(arg, 10, 0)
	if err != nil {
		Errmsg(errmsg_fmt, arg)
		return 0, err
	}
	if i < min {
		Errmsg("Expecting integer value %s to be at least %d; got %d.",
			what, min, i)
        return 0, genericError
	} else if max > 0 && i > max {
        Errmsg("Expecting integer value %s to be at most %d; got %d.",
			what, max, i)
        return 0, genericError
	}
	return i, nil
}
