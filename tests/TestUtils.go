package tests

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"io"
	"os"
	"reflect"
	"symbolic_execution_course/interpreter"
	"symbolic_execution_course/ssa"
	"testing"
)

func SymbolicMachineUnsatTest(
	ctx *interpreter.Context,
	args map[string]any,
	unexpected any,
	t *testing.T,
) {
	solver := symbolicMachineTest(ctx, args, unexpected)

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

func PrepareTest(
	fileName string,
	funcName string,
) *interpreter.Context {
	return runAnalysisFor(fileName, funcName)
}

func SymbolicMachineSatTest(
	ctx *interpreter.Context,
	args map[string]any,
	expected any,
	t *testing.T,
) {
	solver := symbolicMachineTest(ctx, args, expected)

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

func symbolicMachineTest(context *interpreter.Context, args map[string]any, expected any) *z3.Solver {
	solver := context.Solver
	solver.Reset()
	addArgs(args, context.InitState, solver, context)
	addAsserts(context.Results, solver, expected, context)

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

func addAsserts(
	resultStates []*interpreter.State,
	solver *z3.Solver,
	expectedResult any,
	context *interpreter.Context,
) {
	results := make([]z3.Bool, 0)
	for _, resultState := range resultStates {
		if len(resultState.Constraints) == 0 {
			continue
		}

		stateRes := resultState.Constraints[0].AsBool().AsZ3Value().Value.(z3.Bool)
		for _, constraint := range resultState.Constraints[1:] {
			asBool := constraint.AsBool().AsZ3Value().Value.(z3.Bool)
			stateRes = stateRes.And(asBool)
		}

		retConstraint := getResultConstraint(expectedResult, context, resultState)
		stateRes = stateRes.And(retConstraint)

		results = append(results, stateRes)
	}

	if len(results) == 0 {
		return
	}

	solver.Assert(results[0].Or(results[1:]...))

	println("Solver constraints:", solver.String())
	println()
}

func addArgs(
	args map[string]any,
	initState *interpreter.State,
	solver *z3.Solver,
	ctx *interpreter.Context,
) {
	initialStackFrame := initState.StackFrames[0]
	initMemory := initState.Memory

	for argName, argValue := range args {
		switch argCasted := argValue.(type) {
		case StructArg:
			actualArgPtr := initialStackFrame.Values[argName].(*interpreter.Pointer)

			for fieldIdx, expectedVal := range argCasted.fields {
				actualVal := initMemory.LoadField(actualArgPtr, fieldIdx)
				expectedZ3Val := GoToZ3Value(ctx, expectedVal)
				solver.Assert(actualVal.Eq(&expectedZ3Val).AsZ3Value().Value.(z3.Bool))
			}
		case complex128:
			argPtr := initialStackFrame.Values[argName].(*interpreter.Pointer)
			r := real(argCasted)
			expectedRealValue := GoToZ3Value(ctx, r)
			i := imag(argCasted)
			expectedImagValue := GoToZ3Value(ctx, i)

			actualRealValue := initMemory.LoadField(argPtr, 0)
			actualImagValue := initMemory.LoadField(argPtr, 1)

			solver.Assert(expectedImagValue.Eq(actualImagValue).AsZ3Value().Value.(z3.Bool))
			solver.Assert(expectedRealValue.Eq(actualRealValue).AsZ3Value().Value.(z3.Bool))
		case ArrayArg:
			argPtr := initialStackFrame.Values[argName].(*interpreter.Pointer)
			for idx, element := range argCasted.elements {
				idxValue := GoToZ3Value(ctx, idx)
				valueInMemory := initMemory.LoadByArrayIndex(argPtr, &idxValue)
				expectedValue := GoToZ3Value(ctx, element)
				solver.Assert(valueInMemory.Eq(&expectedValue).AsZ3Value().Value.(z3.Bool))
			}

			actualLenValue := initMemory.GetArrayLen(argPtr)
			expectedLenValue := GoToZ3Value(ctx, len(argCasted.elements))
			solver.Assert(actualLenValue.Eq(&expectedLenValue).AsZ3Value().Value.(z3.Bool))
		default:
			argConst := initialStackFrame.Values[argName]
			z3Value := GoToZ3Value(ctx, argValue)
			constraint := argConst.Eq(&z3Value).AsZ3Value().Value.(z3.Bool)

			solver.Assert(constraint)
		}

	}
}

func getResultConstraint(
	expectedResult any,
	ctx *interpreter.Context,
	resultState *interpreter.State,
) z3.Bool {
	switch argCasted := expectedResult.(type) {
	case complex128:
		resultPtrValue := ctx.ReturnValue

		r := real(argCasted)
		expectedRealValue := GoToZ3Value(ctx, r)

		i := imag(argCasted)
		expectedImagValue := GoToZ3Value(ctx, i)

		actualRealValuePtr := resultState.Memory.GetUnsafePointerToField(resultPtrValue, 0, "complex")
		actualRealValue := resultState.Memory.Load(actualRealValuePtr)
		actualImagValuePtr := resultState.Memory.GetUnsafePointerToField(resultPtrValue, 1, "complex")
		actualImagValue := resultState.Memory.Load(actualImagValuePtr)

		realEq := expectedRealValue.Eq(actualRealValue).AsZ3Value().Value.(z3.Bool)
		imagEq := expectedImagValue.Eq(actualImagValue).AsZ3Value().Value.(z3.Bool)

		return realEq.And(imagEq)
	default:
		value := GoToZ3Value(ctx, expectedResult)

		return ctx.ReturnValue.Eq(&value).AsZ3Value().Value.(z3.Bool)
	}
}

type StructArg struct {
	fields map[int]any
}

type ArrayArg struct {
	elements []any
}

func GoToZ3Value(ctx *interpreter.Context, v any) interpreter.Z3Value {
	switch casted := v.(type) {
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		sort := ctx.TypesContext.GetPrimitiveTypeSortOrNil(reflect.TypeOf(v).String())
		bits := ctx.TypesContext.GetPrimitiveTypeBits(reflect.TypeOf(v).String())
		bv := ctx.Z3Context.FromInt(int64(casted.(int)), *sort).(z3.BV)
		return interpreter.Z3Value{
			Context: ctx,
			Value:   bv,
			Bits:    bits,
		}
	case float64, float32:
		sort := ctx.TypesContext.GetPrimitiveTypeSortOrNil(reflect.TypeOf(v).String())
		bits := ctx.TypesContext.GetPrimitiveTypeBits(reflect.TypeOf(v).String())
		float := ctx.Z3Context.FromFloat64(casted.(float64), *sort)
		return interpreter.Z3Value{
			Context: ctx,
			Value:   float,
			Bits:    bits,
		}
	case bool:
		b := ctx.Z3Context.FromBool(casted)
		return interpreter.Z3Value{
			Context: ctx,
			Value:   b,
			Bits:    1,
		}
	default:
		panic("unsupported argument")
	}
}
