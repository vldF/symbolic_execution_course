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

	config := interpreter.InterpreterConfig{PathSelectorMode: interpreter.DFS}
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
		args, err := getArgsByState(
			resultState,
			initState,
			dynamicInterpreterCtx,
			solver,
			functionArgNames,
			functionArgTypes,
		)
		if err != nil {
			println(err)
			return nil
		}

		result = append(result, getTestCode(function.Package().Pkg.Path(), function.Name(), i, args))
	}

	return result
}

func getArgsByState(
	state *interpreter.State,
	initState *interpreter.State,
	ctx *interpreter.Context,
	solver *z3.Solver,
	args []string,
	artTypes []types.Type,
) ([]string, error) {
	solver.Reset()

	for _, c := range state.Constraints {
		solver.Assert(c.AsZ3Value().Value.(z3.Bool))
	}

	ok, err := solver.Check()
	if !ok {
		println(err)
		return nil, err
	}

	model := solver.Model()

	result := make([]string, 0)
	for i, arg := range args {
		argType := artTypes[i]

		argValue := initState.GetValueFromStack(arg)
		z3Value := model.Eval(argValue.AsZ3Value().Value, true)
		goValue := GetArg(state, ctx, model, z3Value, argType)
		result = append(result, goValue.String())
	}

	return result, nil
}

func GetArg(
	resultState *interpreter.State,
	ctx *interpreter.Context,
	model *z3.Model,
	value z3.Value,
	argType types.Type,
) fmt.Stringer {
	memory := resultState.Memory

	switch argType := argType.(type) {
	case *types.Basic:
		switch argType.Kind() {
		case types.Bool:
			val, isLiteral := value.(z3.Bool).AsBool()
			if isLiteral {
				switch {
				case val:
					return &PrimitiveArgument{stringValue: "true"}
				default:
					return &PrimitiveArgument{stringValue: "false"}
				}
			}
			panic("non-literal value")

		case types.Float32, types.Float64:
			v := model.Eval(value, true)
			val, isLiteral := v.(z3.Float).AsBigFloat()
			if isLiteral {
				return floatAsArg(val)
			}
			panic("non-literal value")
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
			v := model.Eval(value, true)
			val, isLiteral, _ := v.(z3.BV).AsInt64()
			if isLiteral {
				return &PrimitiveArgument{stringValue: fmt.Sprintf("%d", val)}
			}
			panic("non-literal value")
		case types.String:
			return &PrimitiveArgument{stringValue: "\"unsupported\""}
		case types.Complex64, types.Complex128:
			ptrValue := &interpreter.Z3Value{
				Context: ctx,
				Value:   value,
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

			panic("non-literal value")
		default:
			panic("unsupported basic type " + argType.String())
		}
	case *types.Pointer:
		return &PointerArgument{innerArgument: GetArg(resultState, ctx, model, value, argType.Elem())}
	case *types.Slice:
		elemType := argType.Elem()
		ptrValue := &interpreter.Z3Value{
			Context: ctx,
			Value:   value,
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
			arrayElements[i] = GetArg(resultState, ctx, model, elementSymbolicValue, elemType)
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
					Value:   value,
					Bits:    64,
				}
				structPtrValuePtr := memory.GetUnsafePointerToField(valueValue, fieldI, argTypeName)
				structPtrValue := memory.Load(structPtrValuePtr).AsZ3Value().Value

				elements[fieldName] = GetArg(resultState, ctx, model, structPtrValue, fieldType)
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

func getTestCode(funcPackage string, funcName string, suffix int, args []string) string {
	var sb strings.Builder
	sb.WriteString("func Test")
	sb.WriteString(funcName)
	sb.WriteString("_" + strconv.Itoa(suffix))
	sb.WriteString("(t *testing.T) {\n")
	sb.WriteString("  ")
	sb.WriteString(getFunctionCallString(funcPackage, funcName, args))
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
