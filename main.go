package main

import (
	"golang.org/x/tools/go/ssa"
	"os"
	"path"
	"path/filepath"
	"strings"
	ssa2 "symbolic_execution_course/ssa"
	"symbolic_execution_course/testgen"
)

func main() {
	packageName := "testdata"
	baseDirPath := "/Users/vfeofilaktov/GolandProjects/symbolic_execution_course/testgen-inp/"
	targetDirPath := "/Users/vfeofilaktov/GolandProjects/symbolic_execution_course/generated/"

	err := os.RemoveAll(targetDirPath)
	if err != nil {
	}

	err = os.MkdirAll(targetDirPath, os.ModePerm)
	if err != nil {
		println(err.Error())
	}

	inputFiles, err := os.ReadDir(baseDirPath)
	if err != nil {
		panic(err.Error())
	}

	for _, file := range inputFiles {
		testMethods := make([]string, 0)
		inf, err := file.Info()
		if err != nil {
			panic(err.Error())
		}

		pkg := ssa2.FromFile(filepath.Join(baseDirPath, inf.Name()), packageName)

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
		resultFileText.WriteString("    \"symbolic_execution_course/testdata\"\n")
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
