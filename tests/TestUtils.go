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

func SymbolicMachineUnsatTest(
	fileName string,
	funcName string,
	args map[string]any,
	unexpected any,
	t *testing.T,
) {
	solver := symbolicMachineTest(fileName, funcName, args, unexpected)

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
	solver := context.Solver
	solver.Reset()
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
	config := interpreter.InterpreterConfig{PathSelectorMode: interpreter.NURS}

	println("function", functionName)
	return interpreter.Interpret(fun, config)
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
	initialStackFrame := state.StackFrames[0]

	for argName, argValue := range args {
		switch argCasted := argValue.(type) {
		case StructArg:
			actualArgPtr := initialStackFrame.Values[argName].(*interpreter.Pointer)

			for fieldIdx, expectedVal := range argCasted.fields {
				actualVal := ctx.Memory.LoadField(actualArgPtr, fieldIdx)
				expectedZ3Val := ctx.GoToZ3Value(expectedVal)
				solver.Assert(actualVal.Eq(&expectedZ3Val).AsZ3Value().Value.(z3.Bool))
			}

			continue
		case complex128:
			argPtr := initialStackFrame.Values[argName].(*interpreter.Pointer)
			r := real(argCasted)
			expectedRealValue := ctx.GoToZ3Value(r)
			i := imag(argCasted)
			expectedImagValue := ctx.GoToZ3Value(i)

			actualRealValue := ctx.Memory.LoadField(argPtr, 0)
			actualImagValue := ctx.Memory.LoadField(argPtr, 1)

			solver.Assert(expectedImagValue.Eq(actualImagValue).AsZ3Value().Value.(z3.Bool))
			solver.Assert(expectedRealValue.Eq(actualRealValue).AsZ3Value().Value.(z3.Bool))
			continue
		case ArrayArg:
			argPtr := initialStackFrame.Values[argName].(*interpreter.Pointer)
			for idx, element := range argCasted.elements {
				idxValue := ctx.GoToZ3Value(idx)
				valueInMemory := ctx.Memory.LoadByArrayIndex(argPtr, &idxValue)
				expectedValue := ctx.GoToZ3Value(element)
				solver.Assert(valueInMemory.Eq(&expectedValue).AsZ3Value().Value.(z3.Bool))
			}

			actualLenValue := ctx.Memory.GetArrayLen(argPtr)
			expectedLenValue := ctx.GoToZ3Value(len(argCasted.elements))
			solver.Assert(actualLenValue.Eq(&expectedLenValue).AsZ3Value().Value.(z3.Bool))
			continue
		}

		argConst := initialStackFrame.Values[argName]
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
	switch argCasted := expectedResult.(type) {
	case complex128:
		resultPtrValue := ctx.ReturnValue

		r := real(argCasted)
		expectedRealValue := ctx.GoToZ3Value(r)

		i := imag(argCasted)
		expectedImagValue := ctx.GoToZ3Value(i)

		actualRealValuePtr := ctx.Memory.GetUnsafePointerToField(resultPtrValue, 0, "complex")
		actualRealValue := ctx.Memory.Load(actualRealValuePtr)
		actualImagValuePtr := ctx.Memory.GetUnsafePointerToField(resultPtrValue, 1, "complex")
		actualImagValue := ctx.Memory.Load(actualImagValuePtr)

		solver.Assert(expectedImagValue.Eq(actualImagValue).AsZ3Value().Value.(z3.Bool))
		solver.Assert(expectedRealValue.Eq(actualRealValue).AsZ3Value().Value.(z3.Bool))
	default:
		value := ctx.GoToZ3Value(expectedResult)
		solver.Assert(ctx.ReturnValue.Eq(&value).AsZ3Value().Value.(z3.Bool))
	}
}

type StructArg struct {
	fields map[int]any
}

type ArrayArg struct {
	elements []any
}
