package interpreter

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"math/bits"
	"symbolic_execution_course/heap"
)

func Interpret(function *ssa.Function) {
	z3Config := z3.NewContextConfig()
	z3Context := z3.NewContext(z3Config)

	typesContext := TypesContext{
		IntBits:     bits.UintSize,
		IntSort:     z3Context.BVSort(bits.UintSize),
		FloatSort:   z3Context.FloatSort(11, 53),
		UnknownSort: z3Context.UninterpretedSort("unknown"),
	}

	states := heap.HeapInit[*State](func(state *State, state2 *State) bool {
		return true // todo
	})

	context := Context{
		z3Context,
		&typesContext,
		nil,
		states,
	}

	ret := getReturnConst(function, &context)
	context.ReturnValue = &Z3Value{
		&context,
		ret,
	}

	addInitState(function, &context)

	for {
		nextState, ok := context.States.Pop()
		if !ok || nextState == nil {
			break
		}

		processState(nextState, &context)
	}
}

func processState(state *State, ctx *Context) {
	nextState := visitInstruction(state.Statement, state, ctx)

	ctx.States.Insert(nextState)
}

func getReturnConst(function *ssa.Function, ctx *Context) z3.Value {
	switch t := function.Signature.Results().At(0).Type().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int, types.Int64, types.Int16, types.Int32, types.Int8:
			return ctx.Z3Context.FreshConst("return", ctx.TypesContext.IntSort).(z3.BV)
		case types.Float64, types.Float32:
			return ctx.Z3Context.FreshConst("return", ctx.TypesContext.FloatSort)
		}
	}

	panic("unsupported type")
}

func addInitState(function *ssa.Function, ctx *Context) {
	if len(function.DomPreorder()) == 0 {
		return
	}

	constraints := make([]BoolPredicate, 0)
	memory := make(map[string]Value)
	entry := function.DomPreorder()[0]

	for _, param := range function.Params {
		name := param.Name()
		sort := ctx.TypeToSort(param.Type())
		val := Z3Value{
			Context: ctx,
			Value:   ctx.Z3Context.Const(name, sort),
		}

		memory[name] = &val
	}

	initState := State{
		Constraints: constraints,
		Memory:      memory,
		Statement:   entry.Instrs[0],
	}

	ctx.States.Insert(&initState)
}

func visitInstruction(instr ssa.Instruction, prevState *State, ctx *Context) *State {
	switch casted := instr.(type) {
	case *ssa.Return:
		return visitReturn(casted, prevState, ctx)
	}

	return prevState
}

func visitReturn(instr *ssa.Return, state *State, ctx *Context) *State {
	newState := state.Copy()
	newState.Statement = instr

	returnValue := visitValue(instr.Results[0], state, ctx)
	newState.Constraints = append(newState.Constraints, ctx.ReturnValue.AsEq(returnValue))

	handleDone(newState, returnValue, ctx)

	return nil
}

func handleDone(state *State, returnValue Value, ctx *Context) {
	fmt.Println("handled!")

	fmt.Println("Memory:", state.Memory)
	fmt.Println("Constraints:", state.Constraints)
	fmt.Println("Statement", state.Statement)
}

func visitValue(val ssa.Value, state *State, ctx *Context) Value {
	switch casted := val.(type) {
	case *ssa.BinOp:
		return visitBinOp(casted, state, ctx)
	case *ssa.Parameter:
		return (state.Memory)[casted.Name()]
	}

	panic("Unsupported value")
}

func visitBinOp(expr *ssa.BinOp, state *State, ctx *Context) Value {
	left := visitValue(expr.X, state, ctx).(NumericValue)
	right := visitValue(expr.Y, state, ctx).(NumericValue)

	switch expr.Op {
	case token.ADD:
		return left.Add(right).(Value)
	case token.SUB:
		return left.Add(right).(Value)
	case token.MUL:
		return left.Add(right).(Value)
	case token.QUO:
		return left.Add(right).(Value)
	default:
		panic("unreachable")
	}
}
