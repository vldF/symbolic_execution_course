package tests

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"io"
	"os"
	"symbolic_execution_course/formulas"
	"symbolic_execution_course/smt"
	"testing"
)

func SymbolicMachineTest(
	fileName string,
	funcName string,
	args map[string]any,
	expected any,
	t *testing.T,
) {
	ctx, solver := runAnalysisFor(fileName, funcName)

	addTestConstraints(ctx, solver, args, expected)

	res, err := solver.Check()

	if err != nil {
		t.Fatal(err)
	}
	if res != true {
		t.Error("unsat!")
	}
}

func addTestConstraints(ctx *formulas.AnalysisContext, solver *z3.Solver, args map[string]any, expected any) {
	for argName, argValue := range args {
		argConst := ctx.Args[argName]
		solver.Assert(ctx.Eq(ctx.GoToZ3Value(argValue), argConst))
	}

	resultConst := ctx.ResultValue
	solver.Assert(ctx.Eq(resultConst, ctx.GoToZ3Value(expected)))
}

func runAnalysisFor(fileName string, functionName string) (*formulas.AnalysisContext, *z3.Solver) {
	sourceFile, fileErr := os.Open("../testdata/" + fileName + ".go")
	if fileErr != nil {
		fmt.Printf("Error opening test file: %v\n", fileErr)
		return nil, nil
	}
	code, readErr := io.ReadAll(sourceFile)
	if readErr != nil {
		fmt.Printf("Error reading test file: %v\n", readErr)
		return nil, nil
	}

	ssa := smt.GetSsa(string(code))
	fun := ssa.Func(functionName)

	solver, z3ctx := smt.CreateSolver()
	analysisCtx := smt.BuildAnalysisContext(fun, z3ctx)

	putConstraintsToSolver(solver, analysisCtx.Constraints)

	println("function", functionName)
	println(" ", "SMT is:")
	println(" ", solver.String())

	return analysisCtx, solver
}

func putConstraintsToSolver(solver *z3.Solver, constraints []formulas.Formula) {
	for _, constraint := range constraints {
		solver.Assert(constraint.Value())
	}
}
