# Comments starting with #: below are remake GNU Makefile comments. See
# https://github.com/rocky/remake/wiki/Rake-tasks-for-gnu-make

.PHONY: all exports

#: Same as make go-repl
all: go-repl

#: The REPL front-end to the go-interactive evaluator
go-repl: extracted_imports.go main.go
	go build -o go-repl main.go

main.go: extracted_imports.go

#: Subsidiary program to import packages into go-repl
make_env: make_env.go
	go build make_env.go

#: Recreate extracted imports
imports: make_env
	./make_env > extracted_imports.go
