go-fish - Yet another Go REPL (read, eval, print loop)
============================================================================

[![Build Status](https://travis-ci.org/rocky/go-fish.png)](https://travis-ci.org/rocky/go-fish)

This project provides an interactive environment for evaluating go
expressions.

Yeah, we know of some others. I think this one has more promise. The
heavy lifting for eval is provided by the Carl Chatfield's
[eval](https://github.com/0xfaded/eval) package.

Setup
-----

* Make sure our GO environment is setup, e.g. *$GOBIN*, *$GOPATH*, ...
* Make sure you have go a 1.2ish version installed.

```
   $ go get github.com/rocky/go-fish
   $ cd $GOPATH/src/github.com/rocky/go-fish
   $ make install  # or look at Makefile
```

If you have
[go-gnureadline](https://code.google.com/p/go-gnureadline/) installed
and want the GNU readline support. In addition to the above try:

```
   $ make go-fish-grl
   $ make install
```

If you have [remake](https://github.com/rocky/remake) installed, you can change *make* above to *remake -x* to see the simple *go* and shell commands that get run. (And *remake --tasks* is also your friend.)

Using
-----

Run `go-fish` or `go-fish-grl`. For now, we have only a static
environment provided and that's exactly the environment that *eval*
uses for itself. (In other words this is ideally suited to introspect
about itself). Since the *eval* package is a reasonable size program,
many of the packages like *os*, *fmt*, *strconv*, *errors*, etc. are
available. Look at the import list in file *eval_imports.go* for the
exact list.

Two global variables have been defined: *env*, and *results*. *env*
the environment that is defined, again largely by
*eval_imports.go*. As you enter expresions, the results are saved in
slice *results*. To quit, enter `Ctrl-D` (EOF) or the word `quit`.

Here's a sample session:

```console
$ ./go-fish
== A simple Go eval REPL ==

Results of expression are stored in variable slice "results".
The environment is stored in global variable "env".

Enter expressions to be evaluated at the "gofish>" prompt.

To see all results, type: "results".

To quit, enter: "quit" or Ctrl-D (EOF).
To get help, enter: "help".
gofish> 10+len("abc" + "def")
Kind = Type = int
results[0] = 16
gofish> os.Args[0]
Kind = Type = string
results[1] = "./go-fish"
gofish> help *
All command names:
help  packages  quit
gofish> packages
All imported packages:
ansi    binary  eval      fmt     math    rand     scanner  sync       time   
ast     bufio   exec      io      os      reflect  sort     syscall    token  
atomic  bytes   filepath  ioutil  parser  repl     strconv  tabwriter  unicode
big     errors  flag      log     pprof   runtime  strings  testing    utf8  
gofish> pkg ansi
=== Package ansi: ===
Constants of ansi:
  Reset
Functions of ansi:
  Color  ColorCode  ColorFunc  DisableColors
Variables of ansi:
  Color  ColorCode  ColorFunc  DisableColors
gofish> quit
go-fish: That's all folks...
$ 
```

See Also
--------

* [What's left to do?](https://github.com/rocky/go-fish/wiki/What%27s-left-to-do%3F)
* [go-play](http://code.google.com/p/go-play): A locally-run HTML5 web interface for experimenting with Go code
* [gub](https://github.com/rocky/ssa-interp): A Go debugger based on the SSA interpreter

[![endorse rocky](https://api.coderwall.com/rocky/endorsecount.png)](https://coderwall.com/rocky)
