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

	for argName, argValue := range args {
		switch argCasted := argValue.(type) {
		case StructArg:
			argSortPtr := ctx.Memory.TypeToSortPtr[argCasted.typeName]
			argCell := ctx.Memory.Mem[argSortPtr].(interpreter.StructValueCell)
			argPtr := state.Stack[argName].(interpreter.StructPointer)
			for i, fieldValue := range argCasted.fields {
				fieldSortPtr := argCell.Fields[i]
				cell := ctx.Memory.Mem[fieldSortPtr].(*interpreter.PrimitiveValueCell)
				switch castedFieldValue := fieldValue.(type) {
				case int64, int, int32, int16, int8:
					value := cell.Z3Arr.Select(argPtr.Ptr.AsZ3Value().Value).(z3.BV)
					constraint := value.Eq(ctx.GoToZ3Value(castedFieldValue).Value.(z3.BV))
					solver.Assert(constraint)
				case float32, float64:
					value := cell.Z3Arr.Select(argPtr.Ptr.AsZ3Value().Value).(z3.Float)
					constraint := value.Eq(ctx.GoToZ3Value(castedFieldValue).Value.(z3.Float))
					solver.Assert(constraint)
				}
			}

			continue
		case complex128:
			typeName := "complex"
			complexSortPtr := ctx.Memory.TypeToSortPtr[typeName]
			argCell := ctx.Memory.Mem[complexSortPtr].(interpreter.StructValueCell)
			argPtr := state.Stack[argName].(interpreter.StructPointer)

			realCell := ctx.Memory.Mem[argCell.Fields[0]].(*interpreter.PrimitiveValueCell)
			realValue := real(argCasted)
			realSymbolicConst := realCell.Z3Arr.Select(argPtr.Ptr.AsZ3Value().Value).(z3.Float)
			constraint := realSymbolicConst.Eq(ctx.GoToZ3Value(realValue).Value.(z3.Float))
			solver.Assert(constraint)

			imagCell := ctx.Memory.Mem[argCell.Fields[1]].(*interpreter.PrimitiveValueCell)
			imagValue := imag(argCasted)
			imagSymbolicConst := imagCell.Z3Arr.Select(argPtr.Ptr.AsZ3Value().Value).(z3.Float)
			constraint = imagSymbolicConst.Eq(ctx.GoToZ3Value(imagValue).Value.(z3.Float))
			solver.Assert(constraint)
			continue
		case ArrayArg:
			typeName := argCasted.elementTypeName
			sortPtr := ctx.Memory.TypeToSortPtr[typeName+"-array-wrapper"]
			wrapperSortPtr := ctx.Memory.Mem[sortPtr].(interpreter.ArrayWrapperCell)

			wrapperPtr := state.Stack[argName].(interpreter.StructPointer)

			lenConst := wrapperSortPtr.GetLen(wrapperPtr.Ptr, ctx).AsZ3Value().Value.(z3.BV)
			expectedLen := ctx.Z3Context.FromInt(int64(len(argCasted.elements)), ctx.TypesContext.IntSort).(z3.BV)
			constraint := lenConst.Eq(expectedLen)
			solver.Assert(constraint)

			for i, element := range argCasted.elements {
				value := ctx.GoToZ3Value(element)
				val := value.AsZ3Value()

				idx := ctx.Z3Context.FromInt(int64(i), ctx.TypesContext.IntSort).(z3.BV)

				arrayValue := interpreter.Z3Value{
					Context: ctx,
					Value:   wrapperSortPtr.GetValue(wrapperPtr.Ptr, ctx).AsZ3Value().Value.(z3.Array).Select(idx),
				}

				solver.Assert(arrayValue.Eq(&val).AsZ3Value().Value.(z3.Bool))
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
	switch castedExpectedResult := expectedResult.(type) {
	case complex128:
		resultPtr := ctx.ReturnValue
		typeName := "complex"
		complexSortPtr := ctx.Memory.TypeToSortPtr[typeName]
		argCell := ctx.Memory.Mem[complexSortPtr].(interpreter.StructValueCell)
		realSortPtr := argCell.Fields[0]
		imagSortPtr := argCell.Fields[1]

		realFields := ctx.Memory.Mem[realSortPtr].(*interpreter.PrimitiveValueCell)
		expectedRealComponent := real(castedExpectedResult)
		constraint := realFields.Z3Arr.Select(resultPtr.Value).(z3.Float).Eq(ctx.GoToZ3Value(expectedRealComponent).Value.(z3.Float))
		solver.Assert(constraint)

		imagFields := ctx.Memory.Mem[imagSortPtr].(*interpreter.PrimitiveValueCell)
		expectedImagComponent := imag(castedExpectedResult)
		constraint = imagFields.Z3Arr.Select(resultPtr.Value).(z3.Float).Eq(ctx.GoToZ3Value(expectedImagComponent).Value.(z3.Float))
		solver.Assert(constraint)
	default:
		value := ctx.GoToZ3Value(expectedResult)
		solver.Assert(ctx.ReturnValue.Eq(&value).AsZ3Value().Value.(z3.Bool))
	}
}

type StructArg struct {
	typeName    string
	fields      map[int]any
	fieldsTypes map[int]types.BasicKind
}

type ArrayArg struct {
	elements        []any
	elementTypeName string
	elementType     types.Type
}
