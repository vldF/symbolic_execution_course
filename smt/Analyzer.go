package smt

import (
	"github.com/aclements/go-z3/z3"
	"go/constant"
	"go/token"
	"go/types"
	ssa2 "golang.org/x/tools/go/ssa"
	"symbolic_execution_course/smt/memory"
)

func BuildAnalysisContext(function *ssa2.Function, z3ctx *z3.Context) *AnalysisContext {
	sorts := Sorts{
		IntSort:     z3ctx.IntSort(),
		FloatSort:   z3ctx.FloatSort(11, 53), // todo
		UnknownSort: z3ctx.UninterpretedSort("unknown"),
		SymPtrSort:  z3ctx.UninterpretedSort("sym-ptr"),
	}

	sorts.ComplexSort = z3ctx.ArraySort(sorts.IntSort, sorts.FloatSort)

	memoryInst := memory.Memory{
		Cells: make(map[*z3.Uninterpreted]*memory.SymMemoryCell),
	}

	ctx := &AnalysisContext{
		Z3ctx:       z3ctx,
		Constraints: []Formula{},
		Sorts:       sorts,
		Args:        make(map[string]z3.Value),
		Memory:      memoryInst,
	}

	returnType := function.Signature.Results().At(0).Type()
	if returnType == nil {
		return nil
	}

	initializeArgs(function, z3ctx, ctx)

	ctx.ResultValue = z3ctx.FreshConst("result", ctx.TypeToSort(returnType))

	visitFunction(*function, ctx)

	return ctx
}

func initializeArgs(function *ssa2.Function, z3ctx *z3.Context, ctx *AnalysisContext) {
	for i := 0; i < function.Signature.Params().Len(); i++ {
		argName := function.Signature.Params().At(i).Name()
		argType := function.Signature.Params().At(i).Type()

		switch argType.(type) {
		case *types.Array:
			arrType := argType.(*types.Array)
			elemType := arrType.Elem()
			elemSort := ctx.TypeToSort(elemType)
			ctx.Args[argName] = ctx.NewArray(elemSort, -1)
		case *types.Slice:
			arrType := argType.(*types.Slice)
			elemType := arrType.Elem()
			elemSort := ctx.TypeToSort(elemType)
			ctx.Args[argName] = ctx.NewArray(elemSort, -1)
		default:
			argSort := ctx.TypeToSort(argType)
			ctx.Args[argName] = z3ctx.Const(argName, argSort)
		}
	}
}

func visitFunction(node ssa2.Function, ctx *AnalysisContext) {
	if len(node.DomPreorder()) == 0 {
		return
	}

	entry := node.DomPreorder()[0]
	formula := visitBlock(entry, ctx)
	ctx.Constraints = append(ctx.Constraints, formula)
}

func visitBlock(block *ssa2.BasicBlock, ctx *AnalysisContext) Formula {
	ctx.PushBasicBlock(block)

	result := NewAndFormula(ctx)
	for _, instr := range block.Instrs {
		if instrFormula := visitInstruction(instr, ctx); instrFormula != nil {
			result.Add(instrFormula)
		}
	}

	ctx.PopBasicBlock()
	return &result
}

func visitValue(value ssa2.Value, ctx *AnalysisContext) z3.Value {
	switch value.(type) {
	case *ssa2.Phi:
		return visitPhi(value.(*ssa2.Phi), ctx)
	case *ssa2.Call:
		call := value.(*ssa2.Call)
		switch callValue := call.Call.Value.(type) {
		case *ssa2.Builtin:
			switch callValue.Name() {
			case "real":
				arg := call.Call.Args[0]
				return visitComplexReal(arg, ctx)

			case "imag":
				arg := call.Call.Args[0]
				return visitComplexImag(arg, ctx)

			case "len":
				arg := visitValue(call.Call.Args[0], ctx)
				return ctx.GetArrayLen(arg.(*z3.Uninterpreted))
			}
		}
	case *ssa2.BinOp:
		return visitBinOp(value.(*ssa2.BinOp), ctx)
	case *ssa2.UnOp:
		return visitUnOp(value.(*ssa2.UnOp), ctx)
	case *ssa2.Convert:
		return visitValue(value.(*ssa2.Convert).X, ctx)
	case *ssa2.Const:
		return visitConst(value.(*ssa2.Const), ctx)
	case *ssa2.Parameter:
		return visitParameter(value.(*ssa2.Parameter), ctx)
	}

	println("unknown value", value.String())
	return nil
}

func visitUnOp(op *ssa2.UnOp, ctx *AnalysisContext) z3.Value {
	switch op.Op {
	case token.MUL:
		return visitDerefValue(op, ctx)
	}

	return nil
}

func visitDerefValue(op *ssa2.UnOp, ctx *AnalysisContext) z3.Value {
	indexAddr := op.X.(*ssa2.IndexAddr)
	arrayId := visitValue(indexAddr.X, ctx)
	indexValue := visitValue(indexAddr.Index, ctx)

	return ctx.GetArrayValue(arrayId.(*z3.Uninterpreted)).Select(indexValue)
}

func visitComplexReal(arg ssa2.Value, ctx *AnalysisContext) z3.Value {
	return ctx.ComplexGetReal(visitValue(arg, ctx))
}

func visitComplexImag(arg ssa2.Value, ctx *AnalysisContext) z3.Value {
	return ctx.ComplexGetImag(visitValue(arg, ctx))
}

func visitBinOp(value *ssa2.BinOp, ctx *AnalysisContext) z3.Value {
	switch value.Op {
	case token.ADD:
		return ctx.Add(visitValue(value.X, ctx), visitValue(value.Y, ctx))
	case token.MUL:
		return ctx.Mul(visitValue(value.X, ctx), visitValue(value.Y, ctx))
	case token.QUO:
		return ctx.Div(visitValue(value.X, ctx), visitValue(value.Y, ctx))
	case token.EQL:
		return ctx.Eq(visitValue(value.X, ctx), visitValue(value.Y, ctx))
	case token.NEQ:
		return ctx.Ne(visitValue(value.X, ctx), visitValue(value.Y, ctx))
	case token.LSS:
		return ctx.Lt(visitValue(value.X, ctx), visitValue(value.Y, ctx))
	case token.LEQ:
		return ctx.Le(visitValue(value.X, ctx), visitValue(value.Y, ctx))
	case token.GTR:
		return ctx.Gt(visitValue(value.X, ctx), visitValue(value.Y, ctx))
	case token.GEQ:
		return ctx.Ge(visitValue(value.X, ctx), visitValue(value.Y, ctx))
	case token.SUB:
		return ctx.Sub(visitValue(value.X, ctx), visitValue(value.Y, ctx))
	case token.REM:
		return visitValue(value.X, ctx).(z3.Int).Rem(visitValue(value.Y, ctx).(z3.Int))
	case token.AND:
		return visitValue(value.X, ctx).(z3.Int).ToBV(64).And(visitValue(value.Y, ctx).(z3.Int).ToBV(64)).SToInt()
	case token.OR:
		return visitValue(value.X, ctx).(z3.Int).ToBV(64).Or(visitValue(value.Y, ctx).(z3.Int).ToBV(64)).SToInt()
	case token.XOR:
		return visitValue(value.X, ctx).(z3.Int).ToBV(64).Xor(visitValue(value.Y, ctx).(z3.Int).ToBV(64)).SToInt()
	case token.SHL:
		return visitValue(value.X, ctx).(z3.Int).ToBV(64).Lsh(visitValue(value.Y, ctx).(z3.Int).ToBV(64)).SToInt()
	case token.SHR:
		return visitValue(value.X, ctx).(z3.Int).ToBV(64).SRsh(visitValue(value.Y, ctx).(z3.Int).ToBV(64)).SToInt()
	default:
		println("unsupported binop", value.String())
		return nil
	}
}

func visitParameter(parameter *ssa2.Parameter, ctx *AnalysisContext) z3.Value {
	return ctx.Args[parameter.Name()]
}

func visitConst(value *ssa2.Const, ctx *AnalysisContext) z3.Value {
	switch value.Value.Kind() {
	case constant.Int:
		return ctx.Z3ctx.FromInt(value.Int64(), ctx.Sorts.IntSort)
	case constant.Bool:
		return ctx.Z3ctx.FromBool(value.Int64() == 0)
	case constant.Float:
		return ctx.Z3ctx.FromFloat64(value.Float64(), ctx.Sorts.FloatSort)
	case constant.Complex:
		return ctx.NewComplex(value.Complex128())
	case constant.String:
	case constant.Unknown:
	}

	return ctx.Z3ctx.FreshConst("unknown", ctx.Sorts.UnknownSort)
}

func visitInstruction(instr ssa2.Instruction, ctx *AnalysisContext) Formula {
	switch instr.(type) {
	case *ssa2.If:
		return visitIf(instr.(*ssa2.If), ctx)
	case *ssa2.Return:
		return visitReturn(instr.(*ssa2.Return), ctx)
	case *ssa2.RunDefers:
		return nil
	case *ssa2.Panic:
		return nil
	case *ssa2.Go:
		return nil
	case *ssa2.Defer:
		return nil
	case *ssa2.Send:
		return nil
	case *ssa2.Store:
		return nil
	case *ssa2.MapUpdate:
		return nil
	case *ssa2.DebugRef:
		return nil
	case ssa2.CallInstruction:
		return nil
	case *ssa2.IndexAddr:
		return nil
	case *ssa2.Jump:
		return visitBlock(instr.(*ssa2.Jump).Block().Succs[0], ctx)
	default:
		println("unknown instruction", instr.String())
		return nil
	}

	return nil
}

func visitPhi(phi *ssa2.Phi, ctx *AnalysisContext) z3.Value {
	block := phi.Block()
	var predIdx int

	for i, pred := range block.Preds {
		if !ctx.HasBasicBlockInHistory(pred) {
			continue
		}

		predIdx = i
		break
	}

	return visitValue(phi.Edges[predIdx], ctx)
}

func visitReturn(instr *ssa2.Return, ctx *AnalysisContext) Formula {
	result := instr.Results[0]
	resultExpr := visitValue(result, ctx)
	resultFormula := NewSimpleFormula(ctx)
	resultFormula.Add(ctx.Eq(resultExpr, ctx.ResultValue))

	return &resultFormula
}

func visitIf(instr *ssa2.If, ctx *AnalysisContext) Formula {
	cond := visitValue(instr.Cond, ctx)
	condBool := cond.(z3.Bool)

	mainBlock := instr.Block().Succs[0]
	elseBlock := instr.Block().Succs[1]

	mainBlockFormula := visitBlock(mainBlock, ctx)
	mainBlockFormulaWithCond := NewFormulaWithCond(ctx, condBool)
	mainBlockFormulaWithCond.Add(mainBlockFormula)

	elseBlockFormula := visitBlock(elseBlock, ctx)
	elseBlockFormulaWithCond := NewFormulaWithCond(ctx, condBool.Not())
	elseBlockFormulaWithCond.Add(elseBlockFormula)

	result := NewOrFormula(ctx)
	result.Add(mainBlockFormulaWithCond)
	result.Add(elseBlockFormulaWithCond)

	return &result
}
