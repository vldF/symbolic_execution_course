package smt

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func GetSsa(code string) *ssa.Package {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "analyzee.go", code, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	files := []*ast.File{f}
	pkg := types.NewPackage("analyzee", "")
	ssaPackage, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()}, fset, pkg, files, ssa.SanityCheckFunctions)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return ssaPackage
}
