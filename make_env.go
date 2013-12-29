// +build ignore

// A tool for creating declarations that get fed into the REPL
package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"log"
	"sort"
	"strings"
	"unicode"
	"code.google.com/p/go.tools/importer"
)

// StartingImport is the import from which we start gathering
// package imports from.
const StartingImport = "github.com/0xfaded/eval"

// isExportedIdent returns false if e is an Ident with name "_".
// These identifers have no associated types.Object, and thus no type.
// isExportedIdent also returns false if identifier e doesn't start
// with an uppercase character and thus is not exported.
func isExportedIdent(e ast.Expr) bool {
	id, ok := e.(*ast.Ident)
	return !(ok && id.Name == "_") && unicode.IsUpper(rune(id.Name[0]))
}

func memberFromDecl(decl ast.Decl, imp *importer.Importer,
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

		// case token.TYPE:
		// 	for _, spec := range decl.Specs {
		// 		id := spec.(*ast.TypeSpec).Name
		// 		if isExportedIdent(id) {
		// 			types = append(types, &id.Name)
		// 		}
		// 	}
		}

	case *ast.FuncDecl:
		id := decl.Name
		if decl.Recv == nil && id.Name == "init" {
			return consts, funcs, types, vars
		}
		if isExportedIdent(id) && !strings.HasPrefix(id.Name, "Test") {
			// Can't handle receiver methods yet
			if decl.Recv == nil {
				filename := imp.Fset.File(decl.Pos()).Name()
				if ! strings.HasSuffix(filename, "_test.go") {
					funcs = append(funcs, &id.Name)
				}
			}
		}
	}
	return consts, funcs, types, vars
}

func extractPackageSymbols(pkg_info *importer.PackageInfo, imp *importer.Importer) {

	var typed_decls = map[string]string {
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
	if len(pkg_info.Files) > 0 {
		// Go source package.
		consts := make([]*string, 0, 10)
		vars   := make([]*string, 0, 10)
		types  := make([]*string, 0, 10)
		funcs  := make([]*string, 0, 10)
		for _, file := range pkg_info.Files {
			for _, decl := range file.Decls {
				consts, funcs, types, vars =
					memberFromDecl(decl, imp, consts, funcs, types, vars)
			}
		}
		fmt.Println("\tconsts = make(map[string] reflect.Value)")
		for _, v := range consts {
			fullname := name + "." + *v
			if typename, found := typed_decls[fullname]; found {
				fmt.Printf("\tconsts[\"%s\"] = reflect.ValueOf(%s(%s))\n",
					*v, typename, fullname)
			} else {
				fmt.Printf("\tconsts[\"%s\"] = reflect.ValueOf(%s)\n", *v, fullname)
			}
		}

		fmt.Println("\n\tfuncs = make(map[string] reflect.Value)")
		for _, v := range funcs {
			fmt.Printf("\tfuncs[\"%s\"] = reflect.ValueOf(%s.%s)\n", *v, name, *v)
		}

		fmt.Println("\n\ttypes = make(map[string] reflect.Type)")
		for _, v := range types {
			fmt.Printf("\ttypes[\"%s\"] = reflect.TypeOf(%s.%s){}\n", *v, name, *v)
		}

		fmt.Println("\n\tvars = make(map[string] reflect.Value)")
		for _, v := range vars   {
			fullname := name + "." + *v
			fmt.Printf("\tvars[\"%s\"] = reflect.ValueOf(&%s)\n", *v, fullname)
		}

		/****
		for _, pkg := range pkg_info.Pkg.Imports() {
			// fmt.Printf("%d %v\n", j, pkg)
			name := pkg.Name()
			if name == "testing" { continue }
			fmt.Printf("pkgs[\"%s\"] = %s\n", name, name)
		}
        ****/

		fmt.Printf(`	pkgs["%s"] = &eval.Env {
		Name: "%s",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
	}
`, name, name)
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
type By func(p1, p2 *importer.PackageInfo) bool

// Sort is a method on the function type, By, that sorts the argument
// slice according to the function.
func (by By) Sort(pkg_infos []*importer.PackageInfo) {
	ps := &packageInfoSorter{
		pkg_infos: pkg_infos,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// packageInfoSorter joins a By function and a slice of
// importer.PackageInfos to be sorted.
type packageInfoSorter struct {
	pkg_infos []*importer.PackageInfo
	by      func(p1, p2 *importer.PackageInfo) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *packageInfoSorter) Len() int {
	return len(s.pkg_infos)
}

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
func writePreamble(pkg_infos []*importer.PackageInfo, name string) {
	path := func(p1, p2 *importer.PackageInfo) bool {
		return p1.Pkg.Path() < p2.Pkg.Path()
	}
	By(path).Sort(pkg_infos)
	fmt.Println(`package repl

import (`)
	for _, pkg_info := range pkg_infos {
		fmt.Printf("\t\"%s\"\n", pkg_info.Pkg.Path())
	}
	fmt.Printf(`)

type pkgType map[string] eval.Pkg

// %sEnvironment adds to the eval.Pkg those imports from the package
// eval (https://%s).
func %sEnvironment(pkgs pkgType) {
	var consts map[string] reflect.Value
	var vars   map[string] reflect.Value
	var types  map[string] reflect.Type
	var funcs  map[string] reflect.Value

`, name, name, StartingImport)
}

// writePostamble finishes of the Go code
func writePostamble() {
	fmt.Println("}")
}

// main creates a Go program that adds to a github.com/0xfaded/eval
// environment (of type eval.Env) the transitive closure of imports
// for a given starting package. Here we use github.com/0xfaded/eval.
func main() {
	impctx := importer.Config{Build: &build.Default}

	// Load, parse and type-check the program.
	imp := importer.New(&impctx)

	var pkgs_string []string = make([] string, 0, 10)
	pkgs_string = append(pkgs_string, StartingImport)
	//pkgs_string = append(pkgs_string, "fmt")

	pkg_infos, _, err := imp.LoadInitialPackages(pkgs_string)
	if err != nil {
		log.Fatal(err)
	}

	pkg_infos = imp.AllPackages()
	var errpkgs []string

	writePreamble(pkg_infos, "Eval")

	for _, pkg_info := range pkg_infos {
		if pkg_info.Err != nil {
			errpkgs = append(errpkgs, pkg_info.Pkg.Path())
		} else {
			extractPackageSymbols(pkg_info, imp)
		}
	}
	if errpkgs != nil {
		log.Fatal("couldn't create these SSA packages due to type errors: %s",
			strings.Join(errpkgs, ", "))
	}

	writePostamble()

}
