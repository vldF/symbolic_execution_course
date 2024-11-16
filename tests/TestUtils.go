package tests

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
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
	context := runAnalysisFor(fileName, funcName)
	for _, resState := range context.Results {
		solver := z3.NewSolver(context.Z3Context)
		addAsserts(resState, solver)
		addArgs(args, resState, solver, context)
		addResultConstraint(solver, expected, context)

		println("solver with test constraints:", solver.String())
		println()

		if ok, err := solver.Check(); !ok {
			println("Unsat!")
			println(err)
			t.Fail()
			return
		}

		println("Model:", solver.Model())
	}
}

func runAnalysisFor(fileName string, functionName string) *interpreter.Context {
	sourceFile, fileErr := os.Open("../testdata/" + fileName + ".go")
	if fileErr != nil {
		fmt.Printf("Error opening test file: %v\n", fileErr)
		return nil
	}
	code, readErr := io.ReadAll(sourceFile)
	if readErr != nil {
		fmt.Printf("Error reading test file: %v\n", readErr)
		return nil
	}

	ssa := ssa.GetSsa(string(code))
	fun := ssa.Func(functionName)

	println("function", functionName)
	return interpreter.Interpret(fun)
}

func addAsserts(state *interpreter.State, solver *z3.Solver) {
	for _, constraint := range state.Constraints {
		solver.Assert(constraint.Value)
	}

	println("Solver constraints:", solver.String())
	println()
}

func addArgs(args map[string]any, state *interpreter.State, solver *z3.Solver, ctx *interpreter.Context) {
	for argName, argValue := range args {
		argConst := state.Memory[argName]
		z3Value := ctx.GoToZ3Value(argValue)
		constraint := argConst.AsEq(&z3Value)

		solver.Assert(constraint.Value)
	}
}

func addResultConstraint(solver *z3.Solver, expectedResult any, ctx *interpreter.Context) {
	value := ctx.GoToZ3Value(expectedResult)
	solver.Assert(ctx.ReturnValue.AsEq(&value).Value)
}
