package ssa

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"io"
	"os"
	"strings"
)

func FromCode(
	code string,
	intrinsicsPath string,
) *ssa.Package {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "analyzee.go", code, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	intrinsicsFile, err := parser.ParseFile(
		fset,
		intrinsicsPath,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	files := []*ast.File{f, intrinsicsFile}
	pkg := types.NewPackage("analyzee", "")
	ssaPackage, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()}, fset, pkg, files, ssa.SanityCheckFunctions|ssa.PrintFunctions)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return ssaPackage
}

func FromFile(
	filePath string,
	packageName string,
	intrinsicsPath string,
) *ssa.Package {
	pathSep := string(os.PathSeparator)
	pathParts := strings.Split(filePath, pathSep)
	fileName := pathParts[len(pathParts)-1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	code, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, code, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	intrinsicsFile, err := parser.ParseFile(fset, intrinsicsPath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	files := []*ast.File{f, intrinsicsFile}
	pkg := types.NewPackage(packageName, "")
	ssaPackage, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()},
		fset,
		pkg,
		files,
		ssa.SanityCheckFunctions|ssa.PrintFunctions,
	)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return ssaPackage
}
