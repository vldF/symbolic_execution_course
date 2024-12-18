package testgen

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"math"
	"math/big"
	"strconv"
	"strings"
	"symbolic_execution_course/interpreter"
)

func GenerateTests(function *ssa.Function) []string {
	result := make([]string, 0)

	config := interpreter.InterpreterConfig{
		PathSelectorMode: interpreter.DFS,
		MainPackage:      function.Package().Pkg.Name(),
	}
	dynamicInterpreterCtx := interpreter.Interpret(function, config)

	initState := dynamicInterpreterCtx.InitState
	solver := dynamicInterpreterCtx.Solver

	functionArgNames := make([]string, 0)
	functionArgTypes := make([]types.Type, 0)
	for i := range function.Signature.Params().Len() {
		param := function.Signature.Params().At(i)
		functionArgNames = append(functionArgNames, param.Name())
		functionArgTypes = append(functionArgTypes, param.Type())
	}

	resultStates := dynamicInterpreterCtx.Results
	for i, resultState := range resultStates {
		solver.Reset()

		for _, c := range resultState.Constraints {
			solver.Assert(c.AsZ3Value().Value.(z3.Bool))
		}

		ok, err := solver.Check()
		if !ok {
			println(err)
			continue
		}

		model := solver.Model()

		args, err := getArgsByState(
			resultState,
			initState,
			dynamicInterpreterCtx,
			model,
			functionArgNames,
			functionArgTypes,
		)
		if err != nil {
			println(err)
			return nil
		}

		arrangeCode := getArrangeCode(resultState, dynamicInterpreterCtx, model)

		funcPackage := function.Package().Pkg.Path()
		funcName := function.Name()
		assertCode := getFunctionCallString(funcPackage, funcName, args)

		fullTestCode := getTestCode(funcName, i, arrangeCode, assertCode)
		result = append(result, fullTestCode)
	}

	return result
}

func getArgsByState(
	state *interpreter.State,
	initState *interpreter.State,
	ctx *interpreter.Context,
	model *z3.Model,
	args []string,
	artTypes []types.Type,
) ([]string, error) {

	result := make([]string, 0)
	for i, arg := range args {
		argType := artTypes[i]

		argValue := initState.GetValueFromStack(arg)
		z3Value := model.Eval(argValue.AsZ3Value().Value, true)
		goValue := GetConcreteValue(state.Memory, ctx, model, z3Value, argType)
		result = append(result, goValue.String())
	}

	return result, nil
}

func getArrangeCode(
	state *interpreter.State,
	ctx *interpreter.Context,
	model *z3.Model,
) string {
	mockedValues := state.Mocker.GetAll()
	if mockedValues == nil || len(mockedValues) == 0 {
		return ""
	}

	var result strings.Builder
	for funcName, values := range mockedValues {
		for i, descriptor := range values {
			z3Value := descriptor.Value.AsZ3Value().Value
			mem := descriptor.Memory
			typ := descriptor.Type

			result.WriteString("    // ")
			result.WriteString("mock needed: ")
			result.WriteString(funcName)
			result.WriteString("(...)")
			result.WriteString("#")
			result.WriteString(strconv.Itoa(i))
			result.WriteString(": ")
			result.WriteString(GetConcreteValue(mem, ctx, model, z3Value, typ).String())
			result.WriteString("\n")
		}
	}

	return result.String()
}

func GetConcreteValue(
	memory *interpreter.Memory,
	ctx *interpreter.Context,
	model *z3.Model,
	symbolicValue z3.Value,
	valueType types.Type,
) fmt.Stringer {
	switch argType := valueType.(type) {
	case *types.Basic:
		switch argType.Kind() {
		case types.Bool:
			val, isLiteral := symbolicValue.(z3.Bool).AsBool()
			if isLiteral {
				switch {
				case val:
					return &PrimitiveArgument{stringValue: "true"}
				default:
					return &PrimitiveArgument{stringValue: "false"}
				}
			}
			panic("non-literal symbolicValue")

		case types.Float32, types.Float64:
			v := model.Eval(symbolicValue, true)
			val, isLiteral := v.(z3.Float).AsBigFloat()
			if isLiteral {
				return floatAsArg(val)
			}
			panic("non-literal symbolicValue")
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
			v := model.Eval(symbolicValue, true)
			val, isLiteral, _ := v.(z3.BV).AsInt64()
			if isLiteral {
				return &PrimitiveArgument{stringValue: fmt.Sprintf("%d", val)}
			}
			panic("non-literal symbolicValue")
		case types.String:
			return &PrimitiveArgument{stringValue: "\"unsupported\""}
		case types.Complex64, types.Complex128:
			ptrValue := &interpreter.Z3Value{
				Context: ctx,
				Value:   symbolicValue,
				Bits:    64,
			}

			rPtr := memory.GetUnsafePointerToField(ptrValue, 0, "complex")
			rValue := memory.Load(rPtr).AsZ3Value().Value
			r, rIsLiteral := model.Eval(rValue, true).(z3.Float).AsBigFloat()
			iPtr := memory.GetUnsafePointerToField(ptrValue, 1, "complex")
			iValue := memory.Load(iPtr).AsZ3Value().Value
			i, iIsLiteral := model.Eval(iValue, true).(z3.Float).AsBigFloat()

			if rIsLiteral && iIsLiteral {
				return &ComplexArgument{
					real: floatAsArg(r).stringValue,
					imag: floatAsArg(i).stringValue,
				}
			}

			panic("non-literal symbolicValue")
		default:
			panic("unsupported basic type " + argType.String())
		}
	case *types.Pointer:
		return &PointerArgument{innerArgument: GetConcreteValue(memory, ctx, model, symbolicValue, argType.Elem())}
	case *types.Slice:
		elemType := argType.Elem()
		ptrValue := &interpreter.Z3Value{
			Context: ctx,
			Value:   symbolicValue,
			Bits:    64,
		}
		elemTypeName := elemType.String()
		if strings.HasPrefix(elemTypeName, "*") {
			elemTypeName = strings.TrimPrefix(elemTypeName, "*")
		}

		arrayPtr := memory.GetUnsafeArrayPointer(ptrValue, elemTypeName)
		arrayLenValue := memory.GetArrayLen(arrayPtr).AsZ3Value().Value
		arrayLen, isLiteral, _ := model.Eval(arrayLenValue, true).(z3.BV).AsInt64()
		if !isLiteral {
			panic("unknown array len")
		}

		arrayElements := make([]fmt.Stringer, arrayLen)
		for i := range arrayLen {
			indexValue := ctx.CreateInt(i, 64)
			elementSymbolicValue := memory.LoadByArrayIndex(arrayPtr, indexValue).AsZ3Value().Value
			arrayElements[i] = GetConcreteValue(memory, ctx, model, elementSymbolicValue, elemType)
		}

		return &ArrayArgument{elementType: elemType.String(), values: arrayElements}
	case *types.Named:
		argTypeName := argType.String()
		argValue := argType.Underlying()
		switch argValue := argValue.(type) {
		case *types.Struct:
			fieldsCount := argValue.NumFields()
			elements := make(map[string]fmt.Stringer, fieldsCount)
			for fieldI := range fieldsCount {
				field := argValue.Field(fieldI)
				fieldName := field.Name()
				fieldType := field.Type()
				valueValue := &interpreter.Z3Value{
					Context: ctx,
					Value:   symbolicValue,
					Bits:    64,
				}
				structPtrValuePtr := memory.GetUnsafePointerToField(valueValue, fieldI, argTypeName)
				structPtrValue := memory.Load(structPtrValuePtr).AsZ3Value().Value

				elements[fieldName] = GetConcreteValue(memory, ctx, model, structPtrValue, fieldType)
			}

			return &StructArgument{
				name:     argTypeName,
				elements: elements,
			}
		}
	}

	panic("unsupported type")
}

func floatAsArg(float *big.Float) *PrimitiveArgument {
	if float == nil {
		return &PrimitiveArgument{stringValue: "math.NaN()"}
	}
	valAsFloat, _ := float.Float64()

	switch {
	case math.IsInf(valAsFloat, 1):
		return &PrimitiveArgument{stringValue: "math.Inf(1)"}
	case math.IsInf(valAsFloat, -1):
		return &PrimitiveArgument{stringValue: "math.Inf(-11)"}
	case math.IsNaN(valAsFloat):
		return &PrimitiveArgument{stringValue: "math.NaN()"}
	}
	return &PrimitiveArgument{stringValue: fmt.Sprintf("%E", valAsFloat)}
}

func getTestCode(
	funcName string,
	suffix int,
	arrangeCode string,
	assertCode string,
) string {
	var sb strings.Builder
	sb.WriteString("func Test")
	sb.WriteString(funcName)
	sb.WriteString("_" + strconv.Itoa(suffix))
	sb.WriteString("(t *testing.T) {\n")
	if len(arrangeCode) != 0 {
		sb.WriteString(arrangeCode)
		sb.WriteString("\n")
	}
	sb.WriteString("    ")
	sb.WriteString(assertCode)
	sb.WriteString("\n")
	sb.WriteString("}")

	return sb.String()
}

func getFunctionCallString(funcPackage string, funcName string, args []string) string {
	var sb strings.Builder

	sb.WriteString(funcPackage)
	sb.WriteString(".")
	sb.WriteString(funcName)
	sb.WriteString("(")

	for _, a := range args {
		sb.WriteString(a)
		sb.WriteString(", ")
	}

	sb.WriteString(")")

	return sb.String()
}
