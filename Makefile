# Comments starting with #: below are remake GNU Makefile comments. See
# https://github.com/rocky/remake/wiki/Rake-tasks-for-gnu-make

.PHONY: all

#: Same as make go-epl
all: go-repl

#: The REPL front-end to the go-interactive evaluator
go-repl:
	go build -o go-repl main.go

#: Subsidiary program to import packages into go-repl
make_env:
	go build make_env
