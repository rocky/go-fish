// +build ignore

package main

// A tool for creating declarations that get fed into the REPL

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"log"
	"strings"
	"unicode"
	"code.google.com/p/go.tools/importer"
)

// isBlankIdent returns true iff e is an Ident with name "_".
// They have no associated types.Object, and thus no type.
//
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
			fmt.Printf("\tconsts[\"%s\"] = reflect.ValueOf(%s.%s)\n", *v, name, *v)
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
			fmt.Printf("\tvars[\"%s\"] = reflect.ValueOf(&%s.%s)\n", *v, name, *v)
		}

		/****
		for _, pkg := range pkg_info.Pkg.Imports() {
			// fmt.Printf("%d %v\n", j, pkg)
			name := pkg.Name()
			if name == "testing" { continue }
			fmt.Printf("pkgs[\"%s\"] = %s\n", name, name)
		}
        ****/

		fmt.Printf(`	pkgs["%s"] = &interactive.Env {
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

func extractPackagesSymbols(pkg_infos []*importer.PackageInfo, imp *importer.Importer) {
	fmt.Println(`
package repl

import (
	"reflect"
	"github.com/0xfaded/go-interactive"
)

type pkgType map[string] interactive.Pkg

func Extract_environment(pkgs pkgType) {
	var consts map[string] reflect.Value
	var vars   map[string] reflect.Value
	var types  map[string] reflect.Type
	var funcs  map[string] reflect.Value
`)
	for _, pkg_info := range pkg_infos {
		extractPackageSymbols(pkg_info, imp)
	}
	fmt.Println("}")
}

func main() {
	impctx := importer.Config{Build: &build.Default}

	// Load, parse and type-check the program.
	imp := importer.New(&impctx)

	var pkgs_string []string = make([] string, 0, 10)
	pkgs_string = append(pkgs_string, "github.com/0xfaded/go-interactive")
	//pkgs_string = append(pkgs_string, "fmt")

	pkg_infos, _, err := imp.LoadInitialPackages(pkgs_string)
	if err != nil {
		log.Fatal(err)
	}
	extractPackagesSymbols(pkg_infos, imp)
}
