// Copyright 2013-2015 Rocky Bernstein. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build ignore

// A tool for creating declarations that get fed into the REPL
package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
	"github.com/rocky/go-types"
	"github.com/rocky/go-loader"
)

// StartingImport is the import from which we start gathering
// package imports from.
const DefaultStartingImport = "github.com/rocky/go-fish"
const DefaultPackage = "repl"
// const DefaultStartingImport = "github.com/0xfaded/eval"

// MyImport is the import string name of this package that the output
// Go code will live in. We have to treat that special because we
// can't include it in an import which causes circular imports. Also,
// variables in this import should not have the package name included
// in it. For example we use refer to function SimpleReadline as
// SimpleReadLine, not repl.SimpleReadline
const MyImport = "github.com/rocky/go-fish"

// excludeSyscallConst excludes "syscall" constants from the output.
// "syscall" has architecture and/OS specific constants that make
// it hard to test this program automatically. So unless specifically
// asked for, we will exclude it by default.
var excludeSyscallConsts bool = true

// isExportedIdent returns false if e is an Ident with name "_".
// These identifers have no associated types.Object, and thus no type.
// isExportedIdent also returns false if identifier e doesn't start
// with an uppercase character and thus is not exported.
func isExportedIdent(e ast.Expr) bool {
	id, ok := e.(*ast.Ident)
	return !(ok && id.Name == "_") && unicode.IsUpper(rune(id.Name[0]))
}

func memberFromDecl(decl ast.Decl, program *loader.Program,
	consts []*string, funcs []*string, types []*string, vars []*string) (
	[]*string, []*string, []*string, []*string) {
	switch decl := decl.(type) {
	case *ast.GenDecl: // import, const, type or var
		switch decl.Tok {
		case token.CONST:
			for _, spec := range decl.Specs {
				for _, id := range spec.(*ast.ValueSpec).Names {
					if isExportedIdent(id) {
						consts = append(consts, &id.Name)
					}
				}
			}

		case token.VAR:
			for _, spec := range decl.Specs {
				for _, id := range spec.(*ast.ValueSpec).Names {
					if isExportedIdent(id) {
						vars = append(vars, &id.Name)
					}
				}
			}

		case token.TYPE:
			for _, spec := range decl.Specs {
				id := spec.(*ast.TypeSpec).Name
				if isExportedIdent(id) {
					types = append(types, &id.Name)
				}
			}
		}

	case *ast.FuncDecl:
		id := decl.Name
		if decl.Recv == nil && id.Name == "init" {
			return consts, funcs, types, vars
		}
		if isExportedIdent(id) && !strings.HasPrefix(id.Name, "Test") {
			// Can't handle receiver methods yet
			if decl.Recv == nil {
				filename := program.Fset.File(decl.Pos()).Name()
				if ! strings.HasSuffix(filename, "_test.go") {
					funcs = append(funcs, &id.Name)
				}
			}
		}
	}
	return consts, funcs, types, vars
}

func fullIdentName(path, pkg, ident string) (fullname string) {
	fullname = pkg + "." + ident
	if "repl" == pkg && MyImport == path {
		// This is me my package! We can't include repl
		fullname = ident
	}
	return fullname
}

func extractPackageSymbols(pkg_info *loader.PackageInfo,
	program *loader.Program) {

	typed_decls := map[string]string {
		"math.MaxInt64": "int64",
		"math.MaxUint16": "uint16",
		"math.MaxUint32": "uint32",
		"math.MaxUint64": "uint64",
		"math.MaxUint8": "uint8",
		"math.MinInt64": "int64",
		"syscall.CLONE_IO": "uint64",
		"syscall.IN_CLASSA_NET": "uint64",
		"syscall.IN_CLASSB_NET": "uint64",
		"syscall.IN_CLASSC_NET": "uint64",
		"syscall.IN_ONESHOT": "uint64",
		"syscall.LINUX_REBOOT_CMD_CAD_ON": "uint64",
		"syscall.LINUX_REBOOT_CMD_HALT": "uint64",
		"syscall.LINUX_REBOOT_CMD_RESTART2": "uint64",
		"syscall.LINUX_REBOOT_CMD_SW_SUSPEND": "uint64",
		"syscall.LINUX_REBOOT_MAGIC1": "uint64",
		"syscall.MS_MGC_MSK": "uint64",
		"syscall.MS_MGC_VAL": "uint64",
		"syscall.RTF_ADDRCLASSMASK": "uint64",
		"syscall.RTF_LOCAL": "uint64",
		"syscall.RT_TABLE_MAX": "uint64",
		"syscall.TIOCGDEV": "uint64",
		"syscall.TIOCGPTN": "uint64",
		"syscall.TUNGETFEATURES": "uint64",
		"syscall.TUNGETIFF": "uint64",
		"syscall.TUNGETSNDBUF": "uint64",
		"syscall.TUNGETVNETHDRSZ": "uint64",
		"syscall.WCLONE": "uint64",
	}

	name := pkg_info.Pkg.Name()
	path := pkg_info.Pkg.Path()
	if len(pkg_info.Files) > 0 {
		// Go source package.
		consts := make([]*string, 0, 10)
		vars   := make([]*string, 0, 10)
		types  := make([]*string, 0, 10)
		funcs  := make([]*string, 0, 10)
		for _, file := range pkg_info.Files {
			for _, decl := range file.Decls {
				consts, funcs, types, vars =
					memberFromDecl(decl, program, consts, funcs, types, vars)
			}
		}
		fmt.Println("\tconsts = make(map[string] reflect.Value)")

		if name == "syscall" && excludeSyscallConsts {
			fmt.Println("\t//syscall constants excluded")
		} else {
			for _, v := range consts {
				fullname := fullIdentName(path, name, *v)
				if typename, found := typed_decls[fullname]; found {
					fmt.Printf("\tconsts[\"%s\"] = reflect.ValueOf(%s(%s))\n",
						*v, typename, fullname)
				} else {
					fmt.Printf("\tconsts[\"%s\"] = reflect.ValueOf(%s)\n", *v, fullname)
				}
			}
		}

		fmt.Println("\n\tfuncs = make(map[string] reflect.Value)")
		for _, v := range funcs {
			fullname := fullIdentName(path, name, *v)
			fmt.Printf("\tfuncs[\"%s\"] = reflect.ValueOf(%s)\n", *v, fullname)
		}

		fmt.Println("\n\ttypes = make(map[string] reflect.Type)")
		for _, v := range types {
			fullname := fullIdentName(path, name, *v)
			fmt.Printf("\ttypes[\"%s\"] = reflect.TypeOf(new(%s)).Elem()\n",
				*v, fullname)
		}

		fmt.Println("\n\tvars = make(map[string] reflect.Value)")
		for _, v := range vars   {
			fullname := fullIdentName(path, name, *v)
			fmt.Printf("\tvars[\"%s\"] = reflect.ValueOf(&%s)\n", *v,
				fullname)
		}

		fmt.Printf(`	pkgs["%s"] = &eval.SimpleEnv {
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
	}
`, name)
	}

	// } else {
	// 	// GC-compiled binary package.
	// 	// No code.
	// 	// No position information.
	// 	scope := p.Object.Scope()
	// 	for _, name := range scope.Names() {
	// 		obj := scope.Lookup(name)
	// 		memberFromObject(p, obj, nil)
	// 		if obj, ok := obj.(*types.TypeName); ok {
	// 			named := obj.Type().(*types.Named)
	// 			for i, n := 0, named.NumMethods(); i < n; i++ {
	// 				memberFromObject(p, named.Method(i), nil)
	// 			}
	// 		}
	// 	}


}

// By is the type of a "less" function that defines the ordering of
// its Planet arguments.
type By func(p1, p2 *loader.PackageInfo) bool

// Sort is a method on the function type, By, that sorts the argument
// slice according to the function.
func (by By) Sort(pkg_infos []*loader.PackageInfo) {
	ps := &packageInfoSorter{
		pkg_infos: pkg_infos,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// packageInfoSorter joins a By function and a slice of
// importer.PackageInfos to be sorted.
type packageInfoSorter struct {
	pkg_infos []*loader.PackageInfo
	by        func(p1, p2 *loader.PackageInfo) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *packageInfoSorter) Len() int {
	return len(s.pkg_infos)
}

// Swap is part of sort.Interface.
// Swap is part of sort.Interface.
func (s *packageInfoSorter) Swap(i, j int) {
	s.pkg_infos[i], s.pkg_infos[j] = s.pkg_infos[j], s.pkg_infos[i]
}

// Less is part of sort.Interface. It is implemented by calling the
// "by" closure in the sorter.
func (s *packageInfoSorter) Less(i, j int) bool {
	return s.by(s.pkg_infos[i], s.pkg_infos[j])
}

// writePreamble prints the initial boiler-plate Go package code. That
// is it starts out:
//     package repl; import (... )
// Packages that end in _test are removed from the list of imported
// packages and this stripped down list is returned.
func writePreamble(pkg_infos map[*types.Package]*loader.PackageInfo,
	name string, startingImport string, pkgName string) []*loader.PackageInfo {
	path := func(p1, p2 *loader.PackageInfo) bool {
		return p1.Pkg.Path() < p2.Pkg.Path()
	}
	pkg_infos2 := []*loader.PackageInfo {}
	for _,pi := range(pkg_infos) {
		pkg_infos2 = append(pkg_infos2, pi)
	}
	By(path).Sort(pkg_infos2)
	fmt.Printf(`package %s

import (
`, pkgName)
	kept_pkgs := []*loader.PackageInfo {}
	for _, pkg_info := range pkg_infos2 {
		path := pkg_info.Pkg.Path()
		if !strings.HasSuffix(path, "_test") {
			if	MyImport != path {
				fmt.Printf("\t\"%s\"\n", path)
			}
			kept_pkgs = append(kept_pkgs, pkg_info)
		}
	}
	fmt.Printf(`)

// %sEnvironment adds to eval.Env those packages included
// with import "%s".

func %sEnvironment() *eval.SimpleEnv {
	var consts map[string] reflect.Value
	var vars   map[string] reflect.Value
	var types  map[string] reflect.Type
	var funcs  map[string] reflect.Value
	var pkgs   map[string] eval.Env = make(map[string] eval.Env)

`, name, startingImport, name)
	return kept_pkgs
}

// writePostamble finishes of the Go code
func writePostamble() {
	fmt.Printf(`
	mainEnv := eval.MakeSimpleEnv()
	mainEnv.Pkgs = pkgs
	return mainEnv
}`)
}

// main creates a Go program that adds to a github.com/0xfaded/eval
// environment (of type eval.Env) the transitive closure of imports
// for a given starting package. Here we use github.com/0xfaded/eval.
func main() {
	startingImport := DefaultStartingImport
	pkgName := DefaultPackage
	numArgs := len(os.Args)
	if numArgs == 2 {
		startingImport = os.Args[1]
	} else if numArgs == 3 {
		startingImport = os.Args[1]
		pkgName = os.Args[2]
	} else if numArgs > 3 {
		fmt.Printf("usage: %s [starting-import [package-name]]\n")
		os.Exit(1)
	}
	fmt.Printf("// starting import: \"%s\"\n", startingImport)

	importPkgs := map[string]bool{startingImport: false}


	config := loader.Config{
		Build: &build.Default,
		SourceImports: true,
		ImportPkgs: importPkgs,
	}

	var pkgs_string []string = make([] string, 0, 10)
	pkgs_string = append(pkgs_string, startingImport)

	program, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	var errpkgs []string

	pkg_infos := writePreamble(program.AllPackages, "Eval",
		startingImport, pkgName)

	for _, pkg_info := range pkg_infos {
		extractPackageSymbols(pkg_info, program)
	}
	if errpkgs != nil {
		log.Fatal("couldn't create these SSA packages due to type errors: %s",
			strings.Join(errpkgs, ", "))
	}

	writePostamble()

}
