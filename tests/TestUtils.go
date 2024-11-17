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
	solver := z3.NewSolver(context.Z3Context)
	addAsserts(context.Results, solver)
	addArgs(args, context.Results[0], solver, context)
	addResultConstraint(solver, expected, context)

	println("solver with test constraints:", solver.String())
	println()

	if ok, err := solver.Check(); !ok {
		println("Unsat!")
		println(err)
		t.Fail()
		return
	}

	model := solver.Model()
	println("Model:", model.String())
	//println("Result:", model.Eval(context.ReturnValue.Value, true).(z3.BV).String())
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

func addAsserts(states []*interpreter.State, solver *z3.Solver) {
	results := make([]z3.Bool, 0)
	for _, state := range states {
		if len(state.Constraints) == 0 {
			continue
		}

		stateRes := state.Constraints[0].AsBool().AsZ3Value().Value.(z3.Bool)
		for _, constraint := range state.Constraints[1:] {
			asBool := constraint.AsBool().AsZ3Value().Value.(z3.Bool)
			stateRes = stateRes.And(asBool)
		}

		results = append(results, stateRes)
	}

	if len(results) == 0 {
		return
	}

	solver.Assert(results[0].Or(results[1:]...))

	println("Solver constraints:", solver.String())
	println()
}

func addArgs(args map[string]any, state *interpreter.State, solver *z3.Solver, ctx *interpreter.Context) {
	res := make([]z3.Bool, 0)

	for argName, argValue := range args {
		argConst := state.Memory[argName]
		z3Value := ctx.GoToZ3Value(argValue)
		constraint := argConst.Eq(&z3Value).AsZ3Value().Value.(z3.Bool)

		res = append(res, constraint)
	}

	solver.Assert(res[0].And(res[1:]...))
}

func addResultConstraint(solver *z3.Solver, expectedResult any, ctx *interpreter.Context) {
	value := ctx.GoToZ3Value(expectedResult)
	solver.Assert(ctx.ReturnValue.Eq(&value).AsZ3Value().Value.(z3.Bool))
}
