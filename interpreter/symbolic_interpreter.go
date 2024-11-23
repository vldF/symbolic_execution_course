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
		IntBits:       bits.UintSize,
		IntSort:       z3Context.BVSort(bits.UintSize),
		FloatSort:     z3Context.FloatSort(11, 53),
		StructPointer: z3Context.BVSort(bits.UintSize),
		UnknownSort:   z3Context.UninterpretedSort("unknown"),
	}

	states := heap.HeapInit[*State](func(state *State, state2 *State) bool {
		return true // todo
	})

	context := Context{
		Z3Context:    z3Context,
		TypesContext: &typesContext,
		States:       states,
		Results:      make([]*State, 0),
		Memory:       nil,
	}
	context.Memory = &Memory{
		context:         &context,
		Mem:             make(map[sortPointer]interface{}),
		StructToSortPtr: make(map[string]sortPointer),
	}

	intArrSort := z3Context.ArraySort(typesContext.IntSort, typesContext.IntSort)
	floatArrSort := z3Context.ArraySort(typesContext.IntSort, typesContext.FloatSort)
	context.Memory.Mem[intPtr] = &PrimitiveValueCell{z3Context.Const("ints", intArrSort).(z3.Array)}
	context.Memory.Mem[floatPtr] = &PrimitiveValueCell{z3Context.Const("floats", floatArrSort).(z3.Array)}

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
			return ctx.Z3Context.FreshConst("return", ctx.TypesContext.IntSort)
		case types.Float64, types.Float32:
			return ctx.Z3Context.FreshConst("return", ctx.TypesContext.FloatSort)
		case types.UntypedComplex, types.Complex64, types.Complex128:
			return ctx.Z3Context.FreshConst("return", ctx.TypesContext.StructPointer)
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
		switch casted := param.Type().(type) {
		case *types.Basic:
			switch casted.Kind() {
			case types.Complex128, types.Complex64, types.UntypedComplex:
				typeName := "complex"

				fields := map[int]types.BasicKind{
					0: types.Float64,
					1: types.Float64,
				}

				ctx.Memory.NewStruct(typeName, fields)

				memory[name] = StructPointer{
					ctx,
					ctx.Memory.StructToSortPtr[typeName],
					ctx.Memory.AllocateStruct(),
					typeName,
				}
			default:
				val := Z3Value{
					Context: ctx,
					Value:   ctx.Z3Context.Const(name, sort),
				}

				memory[name] = &val
			}
		case *types.Named:
			typeName := casted.Obj().Name()
			struct_ := casted.Underlying().(*types.Struct)

			fields := make(map[int]types.BasicKind)
			fieldsCount := struct_.NumFields()
			for i := 0; i < fieldsCount; i++ {
				fields[i] = struct_.Field(i).Type().(*types.Basic).Kind()
			}

			ctx.Memory.NewStruct(typeName, fields)

			memory[name] = StructPointer{
				ctx,
				ctx.Memory.StructToSortPtr[typeName],
				ctx.Memory.AllocateStruct(),
				typeName,
			}
		}
	}

	initState := State{
		Constraints:        constraints,
		Stack:              memory,
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
	case *ssa.Store:
		return visitStore(casted, prevState, ctx)
	}

	return getNextStates(prevState)
}

func visitStore(casted *ssa.Store, state *State, ctx *Context) []*State {
	newState := getNextStates(state)[0]
	//addrValue := visitValue(casted.Addr, state, ctx)
	storeValue := visitValue(casted.Val, state, ctx)

	//newState.Constraints = append(newState.Constraints, addrValue.Eq(storeValue))

	newState.Stack[casted.Addr.Name()] = storeValue

	return []*State{newState}
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

	returnValue := visitValue(instr.Results[0], newState, ctx)
	switch castedReturnValue := returnValue.(type) {
	case StructPointer:
		newState.Constraints = append(newState.Constraints, ctx.ReturnValue.Eq(castedReturnValue.Ptr.value))
	default:
		newState.Constraints = append(newState.Constraints, ctx.ReturnValue.Eq(returnValue))
	}

	handleDone(newState, ctx)

	return nil
}

func handleDone(state *State, ctx *Context) {
	ctx.Results = append(ctx.Results, state)

	fmt.Println("handled!")
	fmt.Println("Stack:", state.Stack)
	fmt.Println("Constraints:", state.Constraints)
	fmt.Println("Statement", state.Statement)
}

func visitValue(val ssa.Value, state *State, ctx *Context) Value {
	name := val.Name()
	if precalculatedValue := state.Stack[name]; precalculatedValue != nil {
		return precalculatedValue
	}

	switch casted := val.(type) {
	case *ssa.BinOp:
		switch t := casted.Type().(type) {
		case *types.Basic:
			switch t.Kind() {
			case types.Complex128, types.Complex64, types.UntypedComplex:
				return visitComplexBinOp(casted, state, ctx)
			}
		}
		return visitSimpleBinOp(casted, state, ctx)
	case *ssa.Parameter:
		return visitParameter(casted, state, ctx)
	case *ssa.Const:
		return visitConst(casted, ctx)
	case *ssa.Phi:
		return visitPhi(casted, state, ctx)
	case *ssa.Convert:
		return visitConvert(casted, state, ctx)
	case *ssa.UnOp:
		return visitUnary(casted, state, ctx)
	case *ssa.FieldAddr:
		return visitFieldAddr(casted, state, ctx)
	case *ssa.Alloc:
		return visitAlloc(casted, state, ctx)
	case *ssa.Call:
		return visitCall(casted, state, ctx)
	}

	panic("Unsupported value " + val.String())
}

func visitComplexBinOp(expr *ssa.BinOp, state *State, ctx *Context) Value {
	left := visitValue(expr.X, state, ctx).(StructPointer)
	right := visitValue(expr.Y, state, ctx).(StructPointer)

	var op func(
		leftReal ArithmeticValue,
		leftImag ArithmeticValue,
		rightReal ArithmeticValue,
		rightImag ArithmeticValue,
	) (Value, Value)

	switch expr.Op {
	case token.ADD:
		op = func(leftReal ArithmeticValue, leftImag ArithmeticValue, rightReal ArithmeticValue, rightImag ArithmeticValue) (Value, Value) {
			return leftReal.Add(rightReal), leftImag.Add(rightImag)
		}
	case token.SUB:
		op = func(leftReal ArithmeticValue, leftImag ArithmeticValue, rightReal ArithmeticValue, rightImag ArithmeticValue) (Value, Value) {
			return leftReal.Sub(rightReal), leftImag.Sub(rightImag)
		}
	case token.MUL:
		op = func(leftReal ArithmeticValue, leftImag ArithmeticValue, rightReal ArithmeticValue, rightImag ArithmeticValue) (Value, Value) {
			realComponent := leftReal.Mul(rightReal).Sub(leftImag.Mul(rightImag))
			imagComponent := leftReal.Mul(rightImag).Add(leftImag.Mul(rightReal))

			return realComponent, imagComponent
		}
	case token.QUO:
		op = func(leftReal ArithmeticValue, leftImag ArithmeticValue, rightReal ArithmeticValue, rightImag ArithmeticValue) (Value, Value) {
			ratio1 := rightImag.Div(rightReal)
			denom1 := rightReal.Add(ratio1.Mul(rightImag))
			e1 := leftReal.Add(leftImag.Mul(ratio1)).Div(denom1).AsZ3Value().Value
			f1 := leftImag.Sub(leftReal.Mul(ratio1)).Div(denom1).AsZ3Value().Value

			ratio2 := rightReal.Div(rightImag)
			denom2 := rightImag.Add(ratio2.Mul(rightReal))
			e2 := leftReal.Mul(ratio2).Add(leftImag).Div(denom2).AsZ3Value().Value
			f2 := leftImag.Mul(ratio2).Sub(leftReal).Div(denom2).AsZ3Value().Value

			realComponentSymb := leftReal.AsZ3Value().Value.(z3.Float).GE(rightImag.AsZ3Value().Value.(z3.Float)).IfThenElse(e1, e2).(z3.Float)
			imagComponentSymb := leftReal.AsZ3Value().Value.(z3.Float).GE(rightImag.AsZ3Value().Value.(z3.Float)).IfThenElse(f1, f2).(z3.Float)

			realComponent := &Z3Value{
				Context: ctx,
				Value:   realComponentSymb,
			}

			imagComponent := &Z3Value{
				Context: ctx,
				Value:   imagComponentSymb,
			}

			return realComponent, imagComponent
		}
	case token.GTR:
	case token.GEQ:
	case token.LSS:
	case token.LEQ:
	case token.REM:
	case token.EQL:
	case token.NEQ:
	default:
		panic("unreachable" + expr.String())
	}

	result := complexBinOp(
		left,
		right,
		op,
		state,
		ctx,
	)

	rememberValue(expr.Name(), result, state)

	return result
}

func complexBinOp(
	left StructPointer,
	right StructPointer,
	op func(leftReal ArithmeticValue, leftImag ArithmeticValue, rightReal ArithmeticValue, rightImag ArithmeticValue) (Value, Value),
	state *State,
	ctx *Context,
) StructPointer {
	if left.structName != right.structName {
		panic("unsupported operation")
	}

	complexSortName := left.structName
	complexSortPtr := ctx.Memory.StructToSortPtr[complexSortName]
	complexSortCell := ctx.Memory.Mem[complexSortPtr].(StructValueCell)

	leftReal, leftImag := getRealAndImagValues(left.Ptr, complexSortCell, ctx)
	rightReal, rightImag := getRealAndImagValues(right.Ptr, complexSortCell, ctx)

	resultRealVal, resultImagVal := op(leftReal, leftImag, rightReal, rightImag)
	resultPtr := ctx.Memory.getNextPtr()
	resultReal, resultImag := getRealAndImagValues(resultPtr, complexSortCell, ctx)
	state.Constraints = append(state.Constraints, resultRealVal.Eq(resultReal))
	state.Constraints = append(state.Constraints, resultImagVal.Eq(resultImag))

	return StructPointer{
		context:    ctx,
		SortPtr:    complexSortPtr,
		Ptr:        resultPtr,
		structName: complexSortName,
	}
}

func getRealAndImagValues(ptr ValuePointer, complexSortCell StructValueCell, ctx *Context) (ArithmeticValue, ArithmeticValue) {
	realSortPtr := complexSortCell.Fields[0]
	realSortCell := ctx.Memory.Mem[realSortPtr].(*PrimitiveValueCell)
	realValue := realSortCell.getValue(ptr, ctx)

	imagSortPtr := complexSortCell.Fields[1]
	imagSortCell := ctx.Memory.Mem[imagSortPtr].(*PrimitiveValueCell)
	imagValue := imagSortCell.getValue(ptr, ctx)

	return realValue.(ArithmeticValue), imagValue.(ArithmeticValue)
}

func visitCall(call *ssa.Call, state *State, ctx *Context) Value {
	switch castedCallValue := call.Call.Value.(type) {
	case *ssa.Builtin:
		switch castedCallValue.Name() {
		case "real":
			arg := visitValue(call.Call.Args[0], state, ctx).(StructPointer)
			return ctx.Memory.GetStructField(arg, 0)
		case "imag":
			arg := visitValue(call.Call.Args[0], state, ctx).(StructPointer)
			return ctx.Memory.GetStructField(arg, 1)
		}
	}

	panic("unsupported call")
}

func visitAlloc(casted *ssa.Alloc, state *State, ctx *Context) Value {
	structName := casted.Type().(*types.Pointer).Elem().(*types.Named).Obj().Name()
	if state.Stack[structName] != nil {
		return state.Stack[structName]
	}

	res := &StructPointer{
		context:    ctx,
		SortPtr:    ctx.Memory.StructToSortPtr[structName],
		Ptr:        ctx.Memory.AllocateStruct(),
		structName: structName,
	}

	rememberValue(structName, res, state)
	return res
}

func visitFieldAddr(casted *ssa.FieldAddr, state *State, ctx *Context) Value {
	casted.X.Type().(*types.Pointer).Elem().(*types.Named).Obj().Name()
	structPtr := visitValue(casted.X, state, ctx)
	fieldIdx := casted.Field
	field := ctx.Memory.GetStructField(structPtr.(StructPointer), fieldIdx)
	return field
}

func visitUnary(casted *ssa.UnOp, state *State, ctx *Context) Value {
	name := casted.Name()
	var result Value

	switch casted.Op {
	case token.MUL:
		result = visitDereference(casted, state, ctx)
	}

	rememberValue(name, result, state)
	return result
}

func visitDereference(casted *ssa.UnOp, state *State, ctx *Context) Value {
	return visitValue(casted.X, state, ctx)
}

func visitParameter(casted *ssa.Parameter, state *State, ctx *Context) Value {
	return (state.Stack)[casted.Name()]
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

func visitSimpleBinOp(expr *ssa.BinOp, state *State, ctx *Context) Value {
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
	state.Stack[name] = value
}
