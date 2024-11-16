package tests

import (
	"fmt"
	"io"
	"os"
	"symbolic_execution_course/interpreter"
	"symbolic_execution_course/ssa"
	"testing"
)

func SymbolicMachineTest(
	fileName string,
	funcName string,
	args map[string]any,
	expected any,
	t *testing.T,
) {
	runAnalysisFor(fileName, funcName)
}

func runAnalysisFor(fileName string, functionName string) {
	sourceFile, fileErr := os.Open("../testdata/" + fileName + ".go")
	if fileErr != nil {
		fmt.Printf("Error opening test file: %v\n", fileErr)
		return
	}
	code, readErr := io.ReadAll(sourceFile)
	if readErr != nil {
		fmt.Printf("Error reading test file: %v\n", readErr)
		return
	}

	ssa := ssa.GetSsa(string(code))
	fun := ssa.Func(functionName)

	println("function", functionName)
	interpreter.Interpret(fun)
}
