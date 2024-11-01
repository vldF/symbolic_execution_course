package main

import (
	"github.com/aclements/go-z3/z3"
	"go/constant"
	"go/token"
	ssa2 "golang.org/x/tools/go/ssa"
	"symbolic_execution_course/formulas"
)

func BuildAnalysisContext(function *ssa2.Function, z3ctx *z3.Context) *formulas.AnalysisContext {
	sorts := formulas.Sorts{
		IntSort:     z3ctx.IntSort(),
		FloatSort:   z3ctx.FloatSort(11, 53), // todo
		UnknownSort: z3ctx.UninterpretedSort("unknown"),
	}

	ctx := &formulas.AnalysisContext{
		Z3ctx:       z3ctx,
		Constraints: []formulas.Formula{},
		Sorts:       sorts,
		ResultValue: z3ctx.FreshConst("result", z3ctx.IntSort()), // todo
	}

	visitFunction(*function, ctx)

	return ctx
}

func visitFunction(node ssa2.Function, ctx *formulas.AnalysisContext) {
	if len(node.DomPreorder()) == 0 {
		return
	}

	entry := node.DomPreorder()[0]
	formula := visitBlock(entry, ctx)
	ctx.Constraints = append(ctx.Constraints, formula)
}

func visitBlock(block *ssa2.BasicBlock, ctx *formulas.AnalysisContext) formulas.Formula {
	result := formulas.NewAndFormula(ctx)
	for _, instr := range block.Instrs {
		if instrFormula := visitInstruction(instr, ctx); instrFormula != nil {
			result.Add(instrFormula)
		}
	}

	return &result
}

func visitValue(value ssa2.Value, ctx *formulas.AnalysisContext) z3.Value {
	switch value.(type) {
	case *ssa2.Alloc:
		return nil
	case *ssa2.Phi:
		return nil
	case *ssa2.Call:
		return nil
	case *ssa2.BinOp:
		return visitBinOp(value.(*ssa2.BinOp), ctx)
	case *ssa2.UnOp:
		return nil
	case *ssa2.ChangeType:
		return nil
	case *ssa2.Convert:
		return nil
	case *ssa2.MultiConvert:
		return nil
	case *ssa2.ChangeInterface:
		return nil
	case *ssa2.MakeClosure:
		return nil
	case *ssa2.MakeMap:
		return nil
	case *ssa2.MakeChan:
		return nil
	case *ssa2.MakeSlice:
		return nil
	case *ssa2.Slice:
		return nil
	case *ssa2.FieldAddr:
		return nil
	case *ssa2.Index:
		return nil
	case *ssa2.Lookup:
		return nil
	case *ssa2.Select:
		return nil
	case *ssa2.Range:
		return nil
	case *ssa2.Next:
		return nil
	case *ssa2.TypeAssert:
		return nil
	case *ssa2.Extract:
		return nil
	case *ssa2.Const:
		return visitConst(value.(*ssa2.Const), ctx)
	case *ssa2.Parameter:
		return visitParameter(value.(*ssa2.Parameter), ctx)
	default:
		println("unknown value", value.String())
		return nil

	}
}

func visitBinOp(value *ssa2.BinOp, ctx *formulas.AnalysisContext) z3.Value {
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
	}

	return nil
}

func visitParameter(parameter *ssa2.Parameter, ctx *formulas.AnalysisContext) z3.Value {
	return ctx.Z3ctx.Const(parameter.Name(), ctx.TypeToSort(parameter.Type()))
}

func visitConst(value *ssa2.Const, ctx *formulas.AnalysisContext) z3.Value {
	switch value.Value.Kind() {
	case constant.Int:
		return ctx.Z3ctx.FromInt(value.Int64(), ctx.Sorts.IntSort)
	case constant.Bool:
		return ctx.Z3ctx.FromBool(value.Int64() == 0)
	case constant.Float:
		return ctx.Z3ctx.FromFloat64(value.Float64(), ctx.Sorts.FloatSort)
	case constant.Complex:
	case constant.String:
	case constant.Unknown:
	}

	return ctx.Z3ctx.FreshConst("unknown", ctx.Sorts.UnknownSort)
}

func visitInstruction(instr ssa2.Instruction, ctx *formulas.AnalysisContext) formulas.Formula {
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
	default:
		println("unknown instruction", instr.String())
		return nil
	}

	return nil
}

func visitReturn(instr *ssa2.Return, ctx *formulas.AnalysisContext) formulas.Formula {
	result := instr.Results[0]
	resultExpr := visitValue(result, ctx)
	resultFormula := formulas.NewSimpleFormula(ctx)
	resultFormula.Add(ctx.Eq(resultExpr, ctx.ResultValue))

	return &resultFormula
}

func visitIf(instr *ssa2.If, ctx *formulas.AnalysisContext) formulas.Formula {
	cond := visitValue(instr.Cond, ctx)
	condBool := cond.(z3.Bool)

	mainBlock := instr.Block().Succs[0]
	elseBlock := instr.Block().Succs[1]

	mainBlockFormula := visitBlock(mainBlock, ctx)
	mainBlockFormulaWithCond := formulas.NewFormulaWithCond(ctx, condBool)
	mainBlockFormulaWithCond.Add(mainBlockFormula)

	elseBlockFormula := visitBlock(elseBlock, ctx)
	elseBlockFormulaWithCond := formulas.NewFormulaWithCond(ctx, condBool.Not())
	elseBlockFormulaWithCond.Add(elseBlockFormula)

	result := formulas.NewOrFormula(ctx)
	result.Add(mainBlockFormulaWithCond)
	result.Add(elseBlockFormulaWithCond)

	return &result
}
