# Comments starting with #: below are remake GNU Makefile comments. See
# https://github.com/rocky/remake/wiki/Rake-tasks-for-gnu-make

.PHONY: all exports test check clean cmd

#: Same as: make go-fish
all: go-fish

#: The non-GNU Readline REPL front-end to the go-interactive evaluator
go-fish: repl_imports.go main.go repl.go cmd
	go build -o go-fish main.go

#: The GNU Readline REPL front-end to the go-interactive evaluator
go-fish-grl: repl_imports.go main_grl.go repl.go
	go build -o go-fish-grl main_grl.go

cmd:
	cd cmd && go build

main.go: repl_imports.go

#: Subsidiary program to import packages into go-fish
make_env: make_env.go
	@echo go build make_env.go

# Note: we have to create the next repl_imports.go to a new place
# either outside of this directory or to a non-go extension, otherwise
# make_env will try to read the file it is trying to create!

#: The recreated extracted imports by running make_env
repl_imports.go: make_env
	./make_env > repl_imports.next && mv repl_imports.next repl_imports.go

#: Check stuff
test: make_env.go
	go test -v

#: Same as: make test
check: test

#: Remove derived files.
clean:
	for file in make_env go-fish go-fish-grl repl_import.go ; do \
		if test -e "$$file" ; then rm $$file ; fi \
	done

#: Install this puppy
install:
	go install
	[ -x ./go-fish ] && cp ./go-fish $$GOBIN/go-fish
	[ -x ./go-fish-grl ] && cp ./go-fish $$GOBIN/go-fish-grl
