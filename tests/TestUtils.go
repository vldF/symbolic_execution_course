package tests

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"go/types"
	"io"
	"os"
	"symbolic_execution_course/interpreter"
	"symbolic_execution_course/ssa"
	"testing"
)

func SymbolicMachineUnsatTest(
	fileName string,
	funcName string,
	args map[string]any,
	expected any,
	t *testing.T,
) {
	solver := symbolicMachineTest(fileName, funcName, args, expected)

	println("solver with test constraints:", solver.String())
	println()

	if ok, err := solver.Check(); !ok {
		println("Unsat! That's OK")
		println(err)
		return
	}

	model := solver.Model()
	println("Model:", model.String())
	println("test expected to be unsat")
	t.Fail()
}

func SymbolicMachineSatTest(
	fileName string,
	funcName string,
	args map[string]any,
	expected any,
	t *testing.T,
) {
	solver := symbolicMachineTest(fileName, funcName, args, expected)

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
}

func symbolicMachineTest(fileName string, funcName string, args map[string]any, expected any) *z3.Solver {
	context := runAnalysisFor(fileName, funcName)
	solver := z3.NewSolver(context.Z3Context)
	addAsserts(context.Results, solver)
	addArgs(args, context.Results[0], solver, context)
	addResultConstraint(solver, expected, context)

	return solver
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
		switch argCasted := argValue.(type) {
		case StructArg:
			//argSortPtr := ctx.Memory.StructToSortPtr[argCasted.name]
			argPtr := state.Stack[argName]
			for _, fieldValue := range argCasted.fields {
				switch castedFieldValue := fieldValue.(type) {
				case int64, int, int32, int16, int8:
					cell := ctx.Memory.Mem[interpreter.IntPtr].(*interpreter.PrimitiveValueCell)
					constraint := cell.Z3Arr.Select(argPtr.(interpreter.StructPointer).Ptr.AsZ3Value().Value).(z3.BV).Eq(ctx.GoToZ3Value(castedFieldValue).Value.(z3.BV))
					solver.Assert(constraint)
				case float32, float64:
					cell := ctx.Memory.Mem[interpreter.FloatPtr].(*interpreter.PrimitiveValueCell)
					float := cell.Z3Arr.Select(argPtr.(interpreter.StructPointer).Ptr.AsZ3Value().Value).(z3.Float)
					constraint := float.Eq(ctx.GoToZ3Value(castedFieldValue).Value.(z3.Float))
					solver.Assert(constraint)
				}
			}

			continue
		}

		argConst := state.Stack[argName]
		z3Value := ctx.GoToZ3Value(argValue)
		constraint := argConst.Eq(&z3Value).AsZ3Value().Value.(z3.Bool)

		res = append(res, constraint)
	}

	if len(res) == 0 {
		return
	}

	if len(res) == 1 {
		solver.Assert(res[0])
		return
	}

	solver.Assert(res[0].And(res[1:]...))
}

func addResultConstraint(solver *z3.Solver, expectedResult any, ctx *interpreter.Context) {
	value := ctx.GoToZ3Value(expectedResult)
	solver.Assert(ctx.ReturnValue.Eq(&value).AsZ3Value().Value.(z3.Bool))
}

type StructArg struct {
	name        string
	fields      map[int]any
	fieldsTypes map[int]types.BasicKind
}
