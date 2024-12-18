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

func Interpret(
	function *ssa.Function,
	interConfig InterpreterConfig,
) *Context {
	z3Config := z3.NewContextConfig()
	z3Context := z3.NewContext(z3Config)

	typesContext := TypesContext{
		Int:    GetPrimitiveIntDescr(z3Context, true, bits.UintSize),
		Int8:   GetPrimitiveIntDescr(z3Context, true, 8),
		Int16:  GetPrimitiveIntDescr(z3Context, true, 16),
		Int32:  GetPrimitiveIntDescr(z3Context, true, 32),
		Int64:  GetPrimitiveIntDescr(z3Context, true, 64),
		UInt:   GetPrimitiveIntDescr(z3Context, false, bits.UintSize),
		UInt8:  GetPrimitiveIntDescr(z3Context, false, 8),
		UInt16: GetPrimitiveIntDescr(z3Context, false, 16),
		UInt32: GetPrimitiveIntDescr(z3Context, false, 32),
		UInt64: GetPrimitiveIntDescr(z3Context, false, 64),

		Float:   GetPrimitiveFloatDescr(z3Context, 8, 24),
		Float32: GetPrimitiveFloatDescr(z3Context, 8, 24),
		Float64: GetPrimitiveFloatDescr(z3Context, 11, 53),

		ArrayIndexSort: z3Context.BVSort(64),
		Pointer:        z3Context.BVSort(bits.UintSize),
		UnknownSort:    z3Context.UninterpretedSort("unknown"),
	}

	states := heap.HeapInit[*State](func(state *State, state2 *State) bool {
		return state.Priority > state2.Priority
	})

	solver := z3.NewSolver(z3Context)

	context := Context{
		Config:       interConfig,
		Z3Context:    z3Context,
		Solver:       solver,
		TypesContext: &typesContext,
		States:       states,
		Results:      make([]*State, 0),
	}

	context.ReturnValue = getReturnConst(function, &context)

	addInitState(function, &context)

	for {
		nextState, ok := context.States.Pop()
		if !ok || nextState == nil {
			break
		}

		if len(nextState.VisitedBasicBlocks) > 100 {
			println("potentially infinitive loop found, skip it")
			continue
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

		nextState.Priority = nextState.GetPriority(ctx.Config.PathSelectorMode)
		ctx.States.Insert(nextState)
	}
}

func getReturnConst(function *ssa.Function, ctx *Context) Value {
	var resSort z3.Sort
	var resBits int

	switch t := function.Signature.Results().At(0).Type().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.UntypedComplex, types.Complex64, types.Complex128:
			resSort = ctx.TypesContext.Pointer
			resBits = 64
		default:
			resSort = *ctx.TypesContext.GetPrimitiveTypeSortOrNil(t.Name())
			resBits = ctx.TypesContext.GetPrimitiveTypeBits(t.Name())
		}
	default:
		panic("unsupported type")
	}

	c := ctx.Z3Context.FreshConst("return", resSort)
	return &Z3Value{
		ctx,
		c,
		resBits,
	}
}

func addInitState(function *ssa.Function, ctx *Context) {
	if len(function.DomPreorder()) == 0 {
		return
	}

	constraints := make([]BoolValue, 0)

	stackFrames := make([]*StackFrame, 0)
	initialFrame := &StackFrame{
		Initiator: nil,
		Values:    map[string]Value{},
	}
	stackFrames = append(stackFrames, initialFrame)
	memory := &Memory{
		ctx:         ctx,
		memoryLines: make(map[sortPtr]z3.Array),
		structures:  make(map[sortPtr]*StructureDescriptor),
	}

	entry := function.DomPreorder()[0]

	for _, param := range function.Params {
		name := param.Name()
		t := param.Type()
		val := newValueOfType(ctx, t, memory, name)
		initialFrame.Values[name] = val
	}

	mocker := &Mocker{
		history: make(map[string][]MockDescriptor),
	}

	initState := &State{
		Priority:           0,
		Constraints:        constraints,
		StackFrames:        stackFrames,
		Statement:          entry.Instrs[0],
		VisitedBasicBlocks: []int{entry.Index},
		Memory:             memory,
		Mocker:             mocker,
	}

	ctx.States.Insert(initState)
	ctx.InitState = initState
}

func newValueOfType(ctx *Context, tpe types.Type, memory *Memory, name string) Value {
	sort := ctx.TypeToSort(tpe)
	switch casted := tpe.(type) {
	case *types.Basic:
		switch casted.Kind() {
		case types.Complex128, types.Complex64, types.UntypedComplex:
			typeName := "complex"

			fields := map[int]string{
				0: "float64",
				1: "float64",
			}

			memory.NewStruct(typeName, fields)
			ptr := memory.NewPtr(typeName)
			return ptr
		default:
			bits := ctx.TypesContext.GetPrimitiveTypeBits(casted.String())
			val := Z3Value{
				Context: ctx,
				Value:   ctx.Z3Context.FreshConst(name, sort),
				Bits:    bits,
			}

			return &val
		}
	case *types.Named:
		typeName := GetTypeName(casted)
		struct_ := casted.Underlying().(*types.Struct)

		fields := make(map[int]string)
		fieldsCount := struct_.NumFields()
		for i := 0; i < fieldsCount; i++ {
			fields[i] = GetTypeName(struct_.Field(i).Type())
		}

		memory.NewStruct(typeName, fields)
		return memory.NewPtr(typeName)
	case *types.Slice:
		typeName := GetTypeName(casted.Elem())
		arrPtr := memory.AllocateArray(typeName, true)

		elemType := casted.Elem()
		switch castedElemType := elemType.(type) {
		case *types.Named:
			memory.NewStruct(GetTypeName(castedElemType), GetStructureFields(castedElemType))
		case *types.Pointer:
			elemType = castedElemType.Elem()
			memory.NewStruct(GetTypeName(castedElemType), GetStructureFields(elemType.(*types.Named)))
		}

		return arrPtr
	}
	return nil
}

func visitInstruction(instr ssa.Instruction, prevState *State, ctx *Context) []*State {
	switch casted := instr.(type) {
	case *ssa.Return:
		return []*State{visitReturn(casted, prevState, ctx)}
	case *ssa.If:
		return visitIfInstr(casted, prevState, ctx)
	case *ssa.Store:
		return visitStoreInstr(casted, prevState, ctx)
	case *ssa.BinOp:
		return visitBinOpInstr(casted, prevState, ctx)
	case *ssa.UnOp:
		return visitUnOpInstr(casted, prevState, ctx)
	case *ssa.Phi:
		return visitPhiInstr(casted, prevState, ctx)
	case *ssa.Convert:
		return visitConvertInstr(casted, prevState, ctx)
	case *ssa.Alloc:
		return visitAllocInstr(casted, prevState)
	case *ssa.Jump:
		return visitJumpInstr(casted, prevState, ctx)
	case *ssa.Slice:
		return visitSliceInstr(casted, prevState, ctx)
	case *ssa.IndexAddr:
		return visitIndexAddrInstr(casted, prevState, ctx)
	case *ssa.MakeSlice:
		return visitMakeSliceInstr(casted, prevState, ctx)
	case *ssa.Call:
		return visitCallInstr(casted, prevState, ctx)
	}

	//panic("unknown instruction " + instr.String())

	return createPossibleNextStates(prevState)
}

func visitMakeSliceInstr(casted *ssa.MakeSlice, state *State, ctx *Context) []*State {
	newState := createPossibleNextStates(state)[0]
	elementType := GetTypeName(casted.Type().(*types.Slice).Elem())
	arrayPtr := newState.Memory.AllocateArray(elementType, false)
	saveToStack(casted.Name(), arrayPtr, newState)

	arrayLen := visitValue(casted.Len, newState, ctx)
	newState.Memory.SetArrayLen(arrayPtr, arrayLen)

	return []*State{newState}
}

func visitIndexAddrInstr(casted *ssa.IndexAddr, state *State, ctx *Context) []*State {
	newState := createPossibleNextStates(state)[0]
	arrayPtr := visitValue(casted.X, newState, ctx).(*Pointer)
	index := visitValue(casted.Index, newState, ctx)
	arrayElementPointer := newState.Memory.GetArrayElementPointer(arrayPtr, index)
	arrayLen := newState.Memory.GetArrayLen(arrayPtr)

	saveToStack(casted.Name(), arrayElementPointer, newState)
	arrayLenIsGreaterThanIndexConstr := arrayLen.(ArithmeticValue).Gt(index.(ArithmeticValue))
	newState.Constraints = append(state.Constraints, arrayLenIsGreaterThanIndexConstr)

	return []*State{newState}
}

func visitSliceInstr(casted *ssa.Slice, state *State, ctx *Context) []*State {
	newState := createPossibleNextStates(state)[0]
	result := visitValue(casted.X, newState, ctx)
	saveToStack(casted.Name(), result, newState)

	return []*State{newState}
}

func visitJumpInstr(_ *ssa.Jump, state *State, _ *Context) []*State {
	return createPossibleNextStates(state)
}

func visitAllocInstr(casted *ssa.Alloc, state *State) []*State {
	newState := createPossibleNextStates(state)[0]
	elem := casted.Type().(*types.Pointer).Elem()
	var result Value
	switch elem := elem.(type) {
	case *types.Array:
		elementType := GetTypeName(elem.Elem())
		result = newState.Memory.AllocateArray(elementType, false)
	case *types.Named:
		structName := GetTypeName(elem)
		newState.Memory.NewStruct(structName, GetStructureFields(elem))
		result = newState.Memory.NewPtr(structName)
	default:
		panic("unknown alloc type " + elem.String())
	}

	saveToStack(casted.Name(), result, newState)
	return []*State{newState}
}

func visitConvertInstr(casted *ssa.Convert, state *State, ctx *Context) []*State {
	newState := createPossibleNextStates(state)[0]

	value := visitValue(casted.X, newState, ctx)
	if ptr, ok := value.(*Pointer); ok {
		value = newState.Memory.Load(ptr)
	}
	var result Value

	switch tpe := casted.Type().(type) {
	case *types.Basic:
		switch tpe.Kind() {
		case types.Int:
			result = value.(ArithmeticValue).AsInt(ctx.TypesContext.Int.Bits)
		case types.Int8:
			result = value.(ArithmeticValue).AsInt(8)
		case types.Int16:
			result = value.(ArithmeticValue).AsInt(16)
		case types.Int32:
			result = value.(ArithmeticValue).AsInt(32)
		case types.Int64:
			result = value.(ArithmeticValue).AsInt(64)

		case types.Float32:
			result = value.(ArithmeticValue).AsFloat(32)
		case types.Float64:
			result = value.(ArithmeticValue).AsFloat(64)
		}
	default:
		panic("Unsupported cast" + casted.String())
	}

	saveToStack(casted.Name(), result, newState)

	return []*State{newState}
}

func visitPhiInstr(casted *ssa.Phi, state *State, ctx *Context) []*State {
	newState := createPossibleNextStates(state)[0]

	block := casted.Block()
	predBlocksIndexes := make([]int, 0)

	for _, pred := range block.Preds {
		predBlocksIndexes = append(predBlocksIndexes, pred.Index)
	}

	for _, visitedBlockIdx := range slices.Backward(newState.VisitedBasicBlocks) {
		if i := slices.Index(predBlocksIndexes, visitedBlockIdx); i != -1 {
			res := visitValue(casted.Edges[i], newState, ctx)
			saveToStack(casted.Name(), res, newState)

			return []*State{newState}
		}
	}

	panic("can't determine the right edge")
}

func visitUnOpInstr(casted *ssa.UnOp, state *State, ctx *Context) []*State {
	newState := createPossibleNextStates(state)[0]

	name := casted.Name()
	var result Value

	switch casted.Op {
	case token.MUL: // todo: unary +, -
		result = visitDereference(casted, newState, ctx)
	}

	saveToStack(name, result, newState)
	return []*State{newState}
}

func visitBinOpInstr(casted *ssa.BinOp, state *State, ctx *Context) []*State {
	newState := createPossibleNextStates(state)[0]

	left := visitValue(casted.X, newState, ctx)
	right := visitValue(casted.Y, newState, ctx)

	var result Value
	switch t := casted.Type().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Complex128, types.Complex64, types.UntypedComplex:
			result = visitComplexBinOp(casted, newState, ctx)
		default:
			switch casted.Op {
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
				panic("unreachable" + casted.String())
			}
		}
	}

	saveToStack(casted.Name(), result, newState)

	return []*State{newState}
}

func visitStoreInstr(casted *ssa.Store, state *State, ctx *Context) []*State {
	newState := createPossibleNextStates(state)[0]
	storeValue := visitValue(casted.Val, newState, ctx)

	if _, ok := casted.Addr.Type().(*types.Pointer); ok {
		_, isParam := casted.Val.(*ssa.Parameter)
		_, isStruct := casted.Val.Type().Underlying().(*types.Struct)
		if isParam && isStruct {
			newState.LastStackFrame().Values[casted.Addr.Name()] = storeValue
		} else {
			ptr := visitValue(casted.Addr, newState, ctx).(*Pointer)
			newState.Memory.Store(ptr, storeValue)
		}
	} else {
		newState.LastStackFrame().Values[casted.Addr.Name()] = storeValue
	}

	return []*State{newState}
}

func visitIfInstr(casted *ssa.If, state *State, ctx *Context) []*State {
	possibleStates := createPossibleNextStates(state)
	result := make([]*State, 0)

	cond := visitValue(casted.Cond, state, ctx).(BoolValue)

	trueState := possibleStates[0]
	trueState.Constraints = append(trueState.Constraints, cond)

	if hasSolution(trueState, ctx) {
		result = append(result, trueState)
	}

	falseState := possibleStates[1]
	falseState.Constraints = append(falseState.Constraints, cond.Not())
	if hasSolution(falseState, ctx) {
		result = append(result, falseState)
	}

	return result
}

func visitReturn(instr *ssa.Return, state *State, ctx *Context) *State {
	newState := state.Copy()
	newState.Statement = instr

	frame := newState.LastStackFrame()
	returnValue := visitValue(instr.Results[0], newState, ctx)
	if frame.Initiator == nil {
		// return from the main function
		switch castedReturnValue := returnValue.(type) {
		case *Pointer:
			newState.Constraints = append(newState.Constraints, ctx.ReturnValue.Eq(castedReturnValue.ptr))
		default:
			newState.Constraints = append(newState.Constraints, ctx.ReturnValue.Eq(returnValue))
		}

		handleDone(newState, ctx)

		return nil
	}

	// return from the function call
	newState.Statement = frame.Initiator
	newState = createPossibleNextStates(newState)[0]
	newState.PopStackFrame()
	saveToStack(frame.Initiator.Name(), returnValue, newState)

	return newState
}

func handleDone(state *State, ctx *Context) {
	ctx.Results = append(ctx.Results, state)

	fmt.Println("handled!")
	fmt.Println("Stack:", state.LastStackFrame().Values)
	fmt.Println("Constraints:", state.Constraints)
	fmt.Println("Statement", state.Statement)
}

func visitValue(val ssa.Value, state *State, ctx *Context) Value {
	name := val.Name()
	if precalculatedValue := state.GetValueFromStack(name); precalculatedValue != nil {
		return precalculatedValue
	}

	switch casted := val.(type) {
	//case *ssa.Parameter:
	//	return visitParameter(casted, state, ctx)
	case *ssa.Const:
		return visitConst(casted, ctx, state)
	case *ssa.FieldAddr:
		return visitFieldAddr(casted, state, ctx)
	case *ssa.IndexAddr:
		return visitIndexAddr(casted, state, ctx)
	}

	panic("Unsupported value " + val.String())
}

func visitIndexAddr(casted *ssa.IndexAddr, state *State, ctx *Context) Value {
	arrayPtr := visitValue(casted.X, state, ctx).(*Pointer)
	index := visitValue(casted.Index, state, ctx)

	return state.Memory.LoadByArrayIndex(arrayPtr, index)
}

func visitComplexBinOp(expr *ssa.BinOp, state *State, ctx *Context) Value {
	left := visitValue(expr.X, state, ctx).(*Pointer)
	right := visitValue(expr.Y, state, ctx).(*Pointer)

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
				Bits:    ctx.TypesContext.Float64.Bits,
			}

			imagComponent := &Z3Value{
				Context: ctx,
				Value:   imagComponentSymb,
				Bits:    ctx.TypesContext.Float64.Bits,
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

	result := complexBinOp(left, right, op, state)

	return result
}

func complexBinOp(
	left *Pointer,
	right *Pointer,
	op func(
		leftReal ArithmeticValue,
		leftImag ArithmeticValue,
		rightReal ArithmeticValue,
		rightImag ArithmeticValue) (Value, Value),
	state *State,
) *Pointer {
	leftReal, leftImag := getRealAndImagValues(left, state)
	rightReal, rightImag := getRealAndImagValues(right, state)

	resultRealVal, resultImagVal := op(leftReal, leftImag, rightReal, rightImag)

	resultPtr := state.Memory.NewPtr("complex")
	state.Memory.StoreField(resultPtr, 0, resultRealVal)
	state.Memory.StoreField(resultPtr, 1, resultImagVal)

	return resultPtr
}

func getRealAndImagValues(ptr *Pointer, state *State) (ArithmeticValue, ArithmeticValue) {
	realValue := state.Memory.LoadField(ptr, 0)
	imagValue := state.Memory.LoadField(ptr, 1)

	return realValue.(ArithmeticValue), imagValue.(ArithmeticValue)
}

func visitCallInstr(call *ssa.Call, state *State, ctx *Context) []*State {
	newState := createPossibleNextStates(state)[0]

	switch castedCallValue := call.Call.Value.(type) {
	case *ssa.Builtin:
		switch castedCallValue.Name() {
		case "real":
			argPtr := visitValue(call.Call.Args[0], newState, ctx).(*Pointer)
			saveToStack(call.Name(), state.Memory.LoadField(argPtr, 0), newState)
		case "imag":
			argPtr := visitValue(call.Call.Args[0], newState, ctx).(*Pointer)
			saveToStack(call.Name(), state.Memory.LoadField(argPtr, 1), newState)
		case "len":
			arrayPtr := visitValue(call.Call.Args[0], newState, ctx).(*Pointer)
			saveToStack(call.Name(), state.Memory.GetArrayLen(arrayPtr), newState)
		}
	case *ssa.Function:
		newState = visitFunctionCall(call, newState, ctx)
	}

	return []*State{newState}
}

func visitFunctionCall(call *ssa.Call, state *State, ctx *Context) *State {
	function := call.Call.Value.(*ssa.Function)

	if isOverriddenFunctionCall(call) {
		return visitOverriddenFunctionCall(call, state, ctx)
	} else if isIntrinsicFunctionCall(call) {
		return visitIntrinsicFunctionCall(call, state, ctx)
	} else if function.Package().Pkg.Name() != ctx.Config.MainPackage {
		return visitMockFunctionCall(call, state, ctx, function)
	} else {
		return visitRealFunctionCall(call, call.Call.Args, state, ctx, function)
	}
}

func visitRealFunctionCall(
	initiator *ssa.Call,
	args []ssa.Value,
	state *State, ctx *Context,
	function *ssa.Function,
) *State {
	newState := state.Copy()
	newState.PushStackFrame(initiator)

	functionArgs := function.Signature.Params()
	for i := range functionArgs.Len() {
		arg := functionArgs.At(i)
		saveToStack(arg.Name(), visitValue(args[i], state, ctx), newState)
	}

	newState.Statement = function.Blocks[0].Instrs[0]

	return newState
}

func visitAssumeFunctionCall(call *ssa.Call, state *State, ctx *Context) *State {
	newState := state.Copy()

	args := call.Call.Args

	for _, arg := range args {
		if _, ok := arg.Type().(*types.Basic); !ok {
			continue
		}

		castedType := arg.Type().(*types.Basic)

		if castedType.Kind() != types.Bool {
			continue
		}

		predicate := visitValue(arg, newState, ctx).(BoolValue)
		newState.Constraints = append(newState.Constraints, predicate)
	}

	return newState
}

func visitMakeSymbolicFunctionCall(call *ssa.Call, state *State, ctx *Context) *State {
	newState := state.Copy()

	symbolType := call.Type()
	name := call.Name()
	val := newValueOfType(ctx, symbolType.Underlying(), state.Memory, name)
	saveToStack(name, val, newState)

	return newState
}

var overriddenFunctions = map[string]string{
	"math.Sqrt": "math_Sqrt",
}

func isOverriddenFunctionCall(call *ssa.Call) bool {
	funcName := call.Call.Value.(*ssa.Function).String()
	_, ok := overriddenFunctions[funcName]

	return ok
}

func visitOverriddenFunctionCall(call *ssa.Call, state *State, ctx *Context) *State {
	originalFuncName := call.Call.Value.(*ssa.Function).String()
	overriddenFuncName := overriddenFunctions[originalFuncName]
	overriddenFunc := call.Parent().Package().Func(overriddenFuncName)

	return visitRealFunctionCall(call, call.Call.Args, state, ctx, overriddenFunc)
}

var intrinsicFunctions = []string{
	"Assume", "MakeSymbolic",
}

func isIntrinsicFunctionCall(call *ssa.Call) bool {
	funcName := call.Call.Value.(*ssa.Function).Object().Name()
	return slices.Contains(intrinsicFunctions, funcName)
}

func visitIntrinsicFunctionCall(call *ssa.Call, state *State, ctx *Context) *State {
	function := call.Call.Value.(*ssa.Function)
	switch function.Object().Name() {
	case "Assume":
		return visitAssumeFunctionCall(call, state, ctx)
	case "MakeSymbolic":
		return visitMakeSymbolicFunctionCall(call, state, ctx)
	default:
		panic("Unknown intrinsic")
	}
}

func visitMockFunctionCall(
	call *ssa.Call,
	state *State,
	ctx *Context,
	function *ssa.Function,
) *State {
	newState := state.Copy()

	symbolType := call.Type()
	name := call.Name()
	tpe := call.Type()
	val := newValueOfType(ctx, symbolType.Underlying(), state.Memory, name)
	saveToStack(name, val, newState)

	newState.Mocker.Add(function.String(), val, newState.Memory, tpe)

	return newState
}

func visitFieldAddr(casted *ssa.FieldAddr, state *State, ctx *Context) Value {
	value := visitValue(casted.X, state, ctx)
	fieldIdx := casted.Field

	switch castedValue := value.(type) {
	case *Pointer:
		field := state.Memory.GetFieldPointer(castedValue, fieldIdx)
		return field
	default:
		typeName := casted.X.Type().(*types.Pointer).Elem().(*types.Named).String()
		return state.Memory.GetUnsafePointerToField(value, fieldIdx, typeName)
	}
}

func visitDereference(casted *ssa.UnOp, state *State, ctx *Context) Value {
	switch casted.X.Type().(type) {
	case *types.Pointer:
		ptr := visitValue(casted.X, state, ctx).(*Pointer)
		return state.Memory.Load(ptr)
	default:
		return visitValue(casted.X, state, ctx)
	}
}

func visitConst(value *ssa.Const, ctx *Context, state *State) Value {
	switch casted := value.Type().(type) {
	case *types.Basic:
		switch casted.Kind() {
		case types.Int:
			return ctx.CreateInt(value.Int64(), ctx.TypesContext.Int.Bits)
		case types.Int8:
			return ctx.CreateInt(value.Int64(), 8)
		case types.Int16:
			return ctx.CreateInt(value.Int64(), 16)
		case types.Int32:
			return ctx.CreateInt(value.Int64(), 32)
		case types.Int64:
			return ctx.CreateInt(value.Int64(), 64)
		case types.Uint:
			return ctx.CreateInt(value.Int64(), ctx.TypesContext.UInt.Bits)
		case types.Uint8:
			return ctx.CreateInt(value.Int64(), 8)
		case types.Uint16:
			return ctx.CreateInt(value.Int64(), 16)
		case types.Uint32:
			return ctx.CreateInt(value.Int64(), 32)
		case types.Uint64:
			return ctx.CreateInt(value.Int64(), 64)
		case types.UntypedInt:
			return ctx.CreateInt(value.Int64(), ctx.TypesContext.Int.Bits)

		case types.Float32:
			return ctx.CreateFloat(value.Float64(), 32)
		case types.Float64:
			return ctx.CreateFloat(value.Float64(), 64)
		case types.UntypedFloat:
			return ctx.CreateFloat(value.Float64(), 64)

		case types.Complex128:
			compx := value.Complex128()
			bits := ctx.TypesContext.Float64.Bits
			ptr := state.Memory.NewPtr("complex")
			state.Memory.StoreField(ptr, 0, ctx.CreateFloat(real(compx), bits))
			state.Memory.StoreField(ptr, 1, ctx.CreateFloat(imag(compx), bits))

			return ptr
		}
	}

	panic("Unsupported type" + value.String())
}

func createPossibleNextStates(state *State) []*State {
	block := state.Statement.Block()
	idx := slices.Index(block.Instrs, state.Statement)

	result := make([]*State, 0)

	if idx+1 >= len(block.Instrs) {
		switch state.Statement.(type) {
		case *ssa.Return, *ssa.Panic:
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

func saveToStack(name string, value Value, state *State) {
	state.LastStackFrame().Values[name] = value
}
