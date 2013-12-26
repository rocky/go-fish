# Comments starting with #: below are remake GNU Makefile comments. See
# https://github.com/rocky/remake/wiki/Rake-tasks-for-gnu-make

.PHONY: all exports

#: Same as make go-fish
all: go-fish

#: The REPL front-end to the go-interactive evaluator
go-fish: extracted_imports.go main.go repl.go
	go build -o go-fish main.go

main.go: extracted_imports.go

#: Subsidiary program to import packages into go-fish
make_env: make_env.go
	go build make_env.go

#: Recreate extracted imports
imports: make_env
	./make_env > extracted_imports.go
