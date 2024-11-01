package main

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	ssa2 "golang.org/x/tools/go/ssa"
	"io"
	"os"
	"symbolic_execution_course/formulas"
	"symbolic_execution_course/smt"
)

func main() {
	runNumbers()
}

func runNumbers() {
	runAnalysisFor("numbers", "integerOperations")
	runAnalysisFor("numbers", "floatOperations")
	runAnalysisFor("numbers", "mixedOperations")
	runAnalysisFor("numbers", "nestedConditions")
	runAnalysisFor("numbers", "bitwiseOperations")
	runAnalysisFor("numbers", "advancedBitwise")
	runAnalysisFor("numbers", "combinedBitwise")
	runAnalysisFor("numbers", "nestedBitwise")
}

func runAnalysisFor(fileName string, functionName string) {
	sourceFile, fileErr := os.Open("testdata/" + fileName + ".go")
	if fileErr != nil {
		fmt.Printf("Error opening test file: %v\n", fileErr)
		return
	}
	code, readErr := io.ReadAll(sourceFile)
	if readErr != nil {
		fmt.Printf("Error reading test file: %v\n", readErr)
		return
	}

	ssa := GetSsa(string(code))
	fun := ssa.Func(functionName)
	println("function", fun.Signature.String())

	runForFunction(fun)
}

func runForFunction(fun *ssa2.Function) {
	solver, z3ctx := smt.CreateSolver()
	analysisCtx := BuildAnalysisContext(fun, z3ctx)

	putConstraintsToSolver(solver, analysisCtx.Constraints)

	println(" ", "SMT is:")
	println(" ", solver.String())
	println("===")

	res, solverErr := solver.Check()
	if solverErr != nil {
		fmt.Printf("Error checking solver: %v\n", solverErr)
		return
	}

	if !res {
		println("Unsat!")
		return
	}

	println("model is")
	println(solver.Model().String())
}

func putConstraintsToSolver(solver *z3.Solver, constraints []formulas.Formula) {
	for _, constraint := range constraints {
		solver.Assert(constraint.Value())
	}
}
