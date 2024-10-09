package smt

import (
	"github.com/aclements/go-z3/z3"
)

func (sCtx *SymContext) NewIntArgument(name string) z3.Int {
	result := sCtx.Ctx.IntConst(name)

	typesCtx := sCtx.TypesCtx
	minValueConst := sCtx.Ctx.FromInt(typesCtx.MinInt, sCtx.Ctx.IntSort()).(z3.Int)
	maxValueConst := sCtx.Ctx.FromInt(typesCtx.MaxInt, sCtx.Ctx.IntSort()).(z3.Int)

	// int max and int min are intentionally excluded
	sCtx.Solver.Assert(result.GT(minValueConst).And(result.LT(maxValueConst)))

	return result
}

func (sCtx *SymContext) NewFloat64Argument(name string) z3.Real {
	result := sCtx.Ctx.RealConst(name)

	typesCtx := sCtx.TypesCtx
	minValueConst := sCtx.Ctx.FromFloat64(typesCtx.MinFloat64, sCtx.Ctx.FloatSort(11, 53)).ToReal()
	maxValueConst := sCtx.Ctx.FromFloat64(typesCtx.MaxFloat64, sCtx.Ctx.FloatSort(11, 53)).ToReal()

	// int max and int min are intentionally excluded
	sCtx.Solver.Assert(result.GT(minValueConst).And(result.LT(maxValueConst)))

	return result
}
