package interpreter

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"math/bits"
	"slices"
	"symbolic_execution_course/heap"
)

func Interpret(function *ssa.Function) *Context {
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
		make([]*State, 0),
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

	return &context
}

func processState(state *State, ctx *Context) {
	nextStates := visitInstruction(state.Statement, state, ctx)

	for _, nextState := range nextStates {
		if nextState == nil {
			continue
		}

		ctx.States.Insert(nextState)
	}
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

	constraints := make([]BoolValue, 0)
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
		Constraints:        constraints,
		Memory:             memory,
		Statement:          entry.Instrs[0],
		VisitedBasicBlocks: []int{entry.Index},
	}

	ctx.States.Insert(&initState)
}

func visitInstruction(instr ssa.Instruction, prevState *State, ctx *Context) []*State {
	switch casted := instr.(type) {
	case *ssa.Return:
		return []*State{visitReturn(casted, prevState, ctx)}
	case *ssa.If:
		return visitIf(casted, prevState, ctx)
	}

	return getNextStates(prevState)
}

func visitIf(casted *ssa.If, state *State, ctx *Context) []*State {
	nextStates := getNextStates(state)

	cond := visitValue(casted.Cond, state, ctx).(BoolValue)

	nextStates[0].Constraints = append(nextStates[0].Constraints, cond)
	nextStates[1].Constraints = append(nextStates[1].Constraints, cond.Not())

	return nextStates
}

func visitReturn(instr *ssa.Return, state *State, ctx *Context) *State {
	newState := state.Copy()
	newState.Statement = instr

	returnValue := visitValue(instr.Results[0], state, ctx)
	newState.Constraints = append(newState.Constraints, ctx.ReturnValue.Eq(returnValue))

	handleDone(newState, ctx)

	return nil
}

func handleDone(state *State, ctx *Context) {
	ctx.Results = append(ctx.Results, state)

	fmt.Println("handled!")
	fmt.Println("Memory:", state.Memory)
	fmt.Println("Constraints:", state.Constraints)
	fmt.Println("Statement", state.Statement)
}

func visitValue(val ssa.Value, state *State, ctx *Context) Value {
	name := val.Name()
	if precalculatedValue := state.Memory[name]; precalculatedValue != nil {
		return precalculatedValue
	}

	switch casted := val.(type) {
	case *ssa.BinOp:
		return visitBinOp(casted, state, ctx)
	case *ssa.Parameter:
		return (state.Memory)[casted.Name()]
	case *ssa.Const:
		return visitConst(casted, ctx)
	case *ssa.Phi:
		return visitPhi(casted, state, ctx)
	case *ssa.Convert:
		return visitConvert(casted, state, ctx)
	}

	panic("Unsupported value" + val.String())
}

func visitConvert(casted *ssa.Convert, state *State, ctx *Context) Value {
	value := visitValue(casted.X, state, ctx)
	switch tpe := casted.Type().(type) {
	case *types.Basic:
		switch tpe.Kind() {
		case types.Int, types.Int64, types.Int16, types.Int32, types.Int8, types.UntypedInt:
			return value.(ArithmeticValue)
		case types.Float64, types.Float32, types.UntypedFloat:
			return value.(ArithmeticValue)
		}
	}

	panic("Unsupported cast" + casted.String())
}

func visitPhi(casted *ssa.Phi, state *State, ctx *Context) Value {
	for idx, pred := range casted.Block().Preds {
		if slices.Index(state.VisitedBasicBlocks, pred.Index) != -1 {
			return visitValue(casted.Edges[idx], state, ctx)
		}
	}

	panic("can't determine the right edge")
}

func visitConst(value *ssa.Const, ctx *Context) Value {
	switch casted := value.Type().(type) {
	case *types.Basic:
		switch casted.Kind() {
		case types.Int, types.Int64, types.Int16, types.Int32, types.Int8,
			types.UntypedInt, types.Uint, types.Uint8, types.Uint16, types.Uint32:
			return &ConcreteIntValue{
				ctx,
				value.Int64(),
			}
		case types.Float64, types.Float32, types.UntypedFloat:
			return &ConcreteFloatValue{
				ctx,
				value.Float64(),
			}
		}
	}

	panic("Unsupported type" + value.String())
}

func visitBinOp(expr *ssa.BinOp, state *State, ctx *Context) Value {
	left := visitValue(expr.X, state, ctx)
	right := visitValue(expr.Y, state, ctx)

	var result Value
	switch expr.Op {
	case token.ADD:
		result = left.(ArithmeticValue).Add(right.(ArithmeticValue))
	case token.SUB:
		result = left.(ArithmeticValue).Sub(right.(ArithmeticValue))
	case token.MUL:
		result = left.(ArithmeticValue).Mul(right.(ArithmeticValue))
	case token.QUO:
		result = left.(ArithmeticValue).Div(right.(ArithmeticValue))
	case token.GTR:
		result = left.(ArithmeticValue).Gt(right.(ArithmeticValue))
	case token.GEQ:
		result = left.(ArithmeticValue).Ge(right.(ArithmeticValue))
	case token.LSS:
		result = left.(ArithmeticValue).Lt(right.(ArithmeticValue))
	case token.LEQ:
		result = left.(ArithmeticValue).Le(right.(ArithmeticValue))
	case token.REM:
		result = left.(ArithmeticValue).Rem(right.(ArithmeticValue))
	case token.EQL:
		result = left.(ArithmeticValue).Eq(right.(ArithmeticValue))
	case token.NEQ:
		result = left.(ArithmeticValue).NotEq(right.(ArithmeticValue))
	case token.OR:
		result = left.Or(right)
	case token.AND:
		result = left.And(right)
	case token.XOR:
		result = left.Xor(right)
	case token.SHL:
		result = left.(ArithmeticValue).Shl(right.(ArithmeticValue))
	case token.SHR:
		result = left.(ArithmeticValue).Shr(right.(ArithmeticValue))
	default:
		panic("unreachable" + expr.String())
	}

	rememberValue(expr.Name(), result, state)

	return result
}

func getNextStates(state *State) []*State {
	block := state.Statement.Block()
	idx := slices.Index(block.Instrs, state.Statement)

	result := make([]*State, 0)

	if idx+1 >= len(block.Instrs) {
		switch state.Statement.(type) {
		case *ssa.Return:
		case *ssa.If:
			state1 := state.Copy()
			state1.Statement = block.Succs[0].Instrs[0]
			state1.VisitedBasicBlocks = append(state1.VisitedBasicBlocks, block.Succs[0].Index)

			state2 := state.Copy()
			state2.Statement = block.Succs[1].Instrs[0]
			state2.VisitedBasicBlocks = append(state2.VisitedBasicBlocks, block.Succs[1].Index)

			result = append(result, state1, state2)
		default:
			nextState := state.Copy()
			nextState.Statement = block.Succs[0].Instrs[0]
			nextState.VisitedBasicBlocks = append(nextState.VisitedBasicBlocks, block.Succs[0].Index)
			result = append(result, nextState)
		}
	} else {
		nextState := state.Copy()
		nextState.Statement = block.Instrs[idx+1]
		result = append(result, nextState)
	}

	return result
}

func rememberValue(name string, value Value, state *State) {
	state.Memory[name] = value
}
