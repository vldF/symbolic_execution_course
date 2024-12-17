package testgen

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"golang.org/x/tools/go/ssa"
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

	functionArgs := make([]string, 0)
	for i := range function.Signature.Params().Len() {
		functionArgs = append(functionArgs, function.Signature.Params().At(i).Name())
	}

	resultStates := dynamicInterpreterCtx.Results
	for i, resultState := range resultStates {
		args, err := getArgsByState(resultState, initState, solver, functionArgs)
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
	solver *z3.Solver,
	args []string,
) ([]any, error) {
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
	result := make([]any, 0)
	for _, arg := range args {
		argValue := initState.GetValueFromStack(arg)
		z3Value := model.Eval(argValue.AsZ3Value().Value, true)
		goValue := decodeZ3Value(z3Value)
		result = append(result, goValue)
	}

	return result, nil
}

func decodeZ3Value(value z3.Value) any {
	switch v := value.(type) {
	case z3.Bool:
		val, isLiteral := v.AsBool()
		if isLiteral {
			return val
		}
		panic("non-literal value")
	case z3.Int:
		val, isLiteral, _ := v.AsInt64()
		if isLiteral {
			return val
		}
		panic("non-literal value")

	case z3.Float:
		val, isLiteral := v.AsBigFloat()
		if isLiteral {
			f, _ := val.Float64()
			return f
		}
		panic("non-literal value")
	case z3.BV:
		val, isLiteral, _ := v.AsInt64()
		if isLiteral {
			return val
		}
		panic("non-literal value")
	}

	panic("unsupported z3 value")
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
		return fmt.Sprintf("%f", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case int64:
		return fmt.Sprintf("%d", v)
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
