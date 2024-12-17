package testgen

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"golang.org/x/tools/go/ssa"
	"math"
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
	functionArgTypes := make([]string, 0)
	for i := range function.Signature.Params().Len() {
		param := function.Signature.Params().At(i)
		functionArgNames = append(functionArgNames, param.Name())
		functionArgTypes = append(functionArgTypes, param.Type().String())
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
	artTypes []string,
) ([]any, error) {
	solver.Reset()

	for _, c := range state.Constraints {
		solver.Assert(c.AsZ3Value().Value.(z3.Bool))
	}

	println("=====")
	println(solver.String())
	println("=====")

	ok, err := solver.Check()
	if !ok {
		println(err)
		return nil, err
	}

	model := solver.Model()
	println(model.String())

	result := make([]any, 0)
	for i, arg := range args {
		argType := artTypes[i]

		argValue := initState.GetValueFromStack(arg)
		z3Value := model.Eval(argValue.AsZ3Value().Value, true)
		goValue := Z3ToGoValue(state, ctx, model, z3Value, argType)
		result = append(result, goValue)
	}

	return result, nil
}

func anyAsString(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float32:
		if math.IsNaN(float64(v)) {
			return "math.NaN()"
		}

		if math.IsInf(float64(v), 1) {
			return "math.Inf(1)"
		}

		if math.IsInf(float64(v), -1) {
			return "math.Inf(-1)"
		}

		return fmt.Sprintf("%E", v)
	case float64:
		if math.IsNaN(v) {
			return "math.NaN()"
		}

		if math.IsInf(v, 1) {
			return "math.Inf(1)"
		}

		if math.IsInf(v, -1) {
			return "math.Inf(-1)"
		}

		println("123: ", fmt.Sprintf("%E", v))

		return fmt.Sprintf("%E", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case complex64:
		return fmt.Sprintf("complex(%s, %s)", anyAsString(real(v)), anyAsString(imag(v)))
	case complex128:
		return fmt.Sprintf("complex(%s, %s)", anyAsString(real(v)), anyAsString(imag(v)))
	case ArrayArgDescriptor:
		var sb strings.Builder
		sb.WriteString("[]")
		sb.WriteString(v.elementTypeName)
		sb.WriteString("{")
		for _, element := range v.elements {
			sb.WriteString(anyAsString(element))
			sb.WriteString(", ")
		}
		sb.WriteString("}")
		return sb.String()
	}

	panic("unsupported type")
}

func getTestCode(funcPackage string, funcName string, suffix int, args []any) string {
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

func getFunctionCallString(funcPackage string, funcName string, args []any) string {
	var sb strings.Builder

	sb.WriteString(funcPackage)
	sb.WriteString(".")
	sb.WriteString(funcName)
	sb.WriteString("(")

	for _, a := range args {
		sb.WriteString(anyAsString(a))
		sb.WriteString(", ")
	}

	sb.WriteString(")")

	return sb.String()
}

func Z3ToGoValue(
	state *interpreter.State,
	ctx *interpreter.Context,
	model *z3.Model,
	v z3.Value,
	typeName string,
) any {
	switch typeName {
	case "bool":
		val, isLiteral := v.(z3.Bool).AsBool()
		if isLiteral {
			return val
		}
		panic("non-literal value")
	case "int", "int8", "int16", "int32", "int64":
		val, isLiteral, _ := v.(z3.BV).AsInt64()
		if isLiteral {
			return val
		}
		panic("non-literal value " + v.String())

	case "float", "float32", "float64":
		val, isLiteral := v.(z3.Float).AsBigFloat()
		if isLiteral {
			if val == nil {
				return math.NaN()
			}

			f, _ := val.Float64()
			return f
		}
		panic("non-literal value")
	case "complex64", "complex128":
		ptrValue := &interpreter.Z3Value{
			Context: ctx,
			Value:   v,
			Bits:    64,
		}

		rPtr := state.Memory.GetUnsafePointerToField(ptrValue, 0, "complex")
		r := state.Memory.Load(rPtr).AsZ3Value().Value
		iPtr := state.Memory.GetUnsafePointerToField(ptrValue, 1, "complex")
		i := state.Memory.Load(iPtr).AsZ3Value().Value

		println("I:", model.Eval(i, true).String())
		println("R:", model.Eval(r, true).String())

		realFloat := Z3ToGoValue(state, ctx, model, model.Eval(r, true), "float64").(float64)
		imagFloat := Z3ToGoValue(state, ctx, model, model.Eval(i, true), "float64").(float64)

		return complex(realFloat, imagFloat)
	}

	switch {
	case strings.HasPrefix(typeName, "[]"):
		elementTypeName := strings.TrimPrefix(typeName, "[]")
		arrayPtrValue := &interpreter.Z3Value{
			Context: ctx,
			Value:   v,
			Bits:    64,
		}
		arrayPtr := state.Memory.GetUnsafeArrayPointer(arrayPtrValue, elementTypeName)

		arrayLenValue := state.Memory.GetArrayLen(arrayPtr)
		arrayLenZ3Value := model.Eval(arrayLenValue.AsZ3Value().Value, true)
		arrayLen := Z3ToGoValue(state, ctx, model, arrayLenZ3Value, "int").(int64)

		elements := make([]any, arrayLen)
		for i := range arrayLen {
			idx := ctx.CreateInt(i, 64)
			elementValue := state.Memory.LoadByArrayIndex(arrayPtr, idx).AsZ3Value().Value
			elementZ3Value := model.Eval(elementValue, true)
			elements[i] = Z3ToGoValue(state, ctx, model, elementZ3Value, elementTypeName)
		}

		return ArrayArgDescriptor{elements: elements, elementTypeName: elementTypeName}
	}

	panic("unsupported type")
}

type ArrayArgDescriptor struct {
	elements        []any
	elementTypeName string
}
