package main

import (
	"context"
	"github.com/urfave/cli/v3"
	"golang.org/x/tools/go/ssa"
	"os"
	"path"
	"path/filepath"
	"strings"
	ssa2 "symbolic_execution_course/ssa"
	"symbolic_execution_course/testgen"
)

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "inputPath",
				Required: true,
				Usage:    "set a path with input files in go lang",
			},
			&cli.StringFlag{
				Name:     "inputPackage",
				Required: true,
				Usage:    "set a package of input files",
			},

			&cli.StringFlag{
				Name:     "outputPackage",
				Required: true,
				Usage:    "set a package of generated go files",
			},
			&cli.StringFlag{
				Name:     "outputPath",
				Required: true,
				Usage:    "set a package of input files",
			},
		},

		Action: run,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		println(err)
	}
}

func run(ctx context.Context, command *cli.Command) error {
	targetDirPath := command.String("outputPath")
	baseDirPath := command.String("inputPath")
	inputPackage := command.String("inputPackage")
	outputPackage := command.String("outputPackage")

	err := os.RemoveAll(targetDirPath)
	if err != nil {
	}

	err = os.MkdirAll(targetDirPath, os.ModePerm)
	if err != nil {
		return err
	}

	inputFiles, err := os.ReadDir(baseDirPath)
	if err != nil {
		return err
	}

	for _, file := range inputFiles {
		testMethods := make([]string, 0)
		inf, err := file.Info()
		if err != nil {
			return err
		}

		pkg := ssa2.FromFile(filepath.Join(baseDirPath, inf.Name()), outputPackage)

		functions := getAllFunctions(pkg)
		hasMathImport := false
		for _, functionName := range functions {
			funcSsa := pkg.Func(functionName)
			res := testgen.GenerateTests(funcSsa)
			for _, resMethod := range res {
				if strings.Contains(resMethod, "math") {
					hasMathImport = true
				}
			}

			testMethods = append(testMethods, res...)
		}

		if len(testMethods) == 0 {
			continue
		}

		fileName := strings.TrimSuffix(file.Name(), ".go")
		targetFilePath := path.Join(targetDirPath, fileName+"_test.go")

		var resultFileText strings.Builder
		resultFileText.WriteString("package generated\n")
		resultFileText.WriteString("\n")
		resultFileText.WriteString("import (\n")
		resultFileText.WriteString("    \"")
		resultFileText.WriteString(inputPackage)
		resultFileText.WriteString("\"\n")
		resultFileText.WriteString("    \"testing\"\n")
		if hasMathImport {
			resultFileText.WriteString("    \"math\"\n")
		}
		resultFileText.WriteString(")\n")
		resultFileText.WriteString("\n")

		for _, method := range testMethods {
			resultFileText.WriteString(method + "\n")
		}

		err = os.WriteFile(targetFilePath, []byte(resultFileText.String()), 0644)
		if err != nil {
			println(err.Error())
			continue
		}
	}

	return nil
}

func getAllFunctions(pkg *ssa.Package) []string {
	result := make([]string, 0)
	for _, member := range pkg.Members {
		switch castedMember := member.(type) {
		case *ssa.Function:
			firstLetter := string([]rune(castedMember.Name())[0])
			if castedMember.Synthetic == "" && strings.ToUpper(firstLetter) == firstLetter {
				result = append(result, member.(*ssa.Function).Name())
			}
		}
	}

	return result
}
