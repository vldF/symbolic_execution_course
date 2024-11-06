package smt

import "github.com/aclements/go-z3/z3"

type Z3Complex = z3.Array

func (ctx *AnalysisContext) NewComplex(from complex128) Z3Complex {
	realEntry := ctx.Z3ctx.FromFloat64(real(from), ctx.Sorts.FloatSort)
	imagEntry := ctx.Z3ctx.FromFloat64(imag(from), ctx.Sorts.FloatSort)

	return ctx.newComplex(realEntry, imagEntry)
}

func (ctx *AnalysisContext) newComplex(real z3.Float, imag z3.Float) Z3Complex {
	// [real, imag]
	complexConst := ctx.Z3ctx.FreshConst("complex", ctx.Sorts.ComplexSort).(z3.Array)
	complexConst = complexConst.Store(ctx.realIdx(), real)
	complexConst = complexConst.Store(ctx.imagIdx(), imag)

	return complexConst
}

func (ctx *AnalysisContext) ComplexGetReal(value z3.Value) z3.Float {
	if !ctx.IsComplex(value) {
		panic("not a complex value")
	}

	return value.(z3.Array).Select(ctx.realIdx()).(z3.Float)
}

func (ctx *AnalysisContext) ComplexGetImag(value z3.Value) z3.Float {
	if !ctx.IsComplex(value) {
		panic("not a complex value")
	}

	return value.(z3.Array).Select(ctx.imagIdx()).(z3.Float)
}

func (ctx *AnalysisContext) realIdx() z3.Value {
	return ctx.Z3ctx.FromInt(0, ctx.Sorts.IntSort)
}

func (ctx *AnalysisContext) imagIdx() z3.Value {
	return ctx.Z3ctx.FromInt(1, ctx.Sorts.IntSort)
}

func (ctx *AnalysisContext) IsComplex(value z3.Value) bool {
	return value.Sort().AsAST().ID() == ctx.Sorts.ComplexSort.AsAST().ID()
}

func (ctx *AnalysisContext) ComplexAdd(left Z3Complex, right Z3Complex) Z3Complex {
	return ctx.newComplex(
		ctx.ComplexGetReal(left).Add(ctx.ComplexGetReal(right)),
		ctx.ComplexGetImag(left).Add(ctx.ComplexGetImag(right)))
}

func (ctx *AnalysisContext) ComplexSub(left Z3Complex, right Z3Complex) Z3Complex {
	return ctx.newComplex(
		ctx.ComplexGetReal(left).Sub(ctx.ComplexGetReal(right)),
		ctx.ComplexGetImag(left).Sub(ctx.ComplexGetImag(right)))
}

func (ctx *AnalysisContext) ComplexMul(left Z3Complex, right Z3Complex) Z3Complex {
	realEntry := ctx.ComplexGetReal(left).Mul(ctx.ComplexGetReal(right)).Sub(ctx.ComplexGetImag(left).Mul(ctx.ComplexGetImag(right)))
	imagEntry := ctx.ComplexGetReal(left).Mul(ctx.ComplexGetImag(right)).Add(ctx.ComplexGetImag(left).Mul(ctx.ComplexGetReal(right)))

	return ctx.newComplex(realEntry, imagEntry)
}

func (ctx *AnalysisContext) ComplexEq(left Z3Complex, right Z3Complex) z3.Bool {
	return ctx.ComplexGetReal(left).Eq(ctx.ComplexGetReal(right)).And(ctx.ComplexGetImag(left).Eq(ctx.ComplexGetImag(right)))
}

func (ctx *AnalysisContext) ComplexDiv(n Z3Complex, m Z3Complex) Z3Complex {
	realN := ctx.ComplexGetReal(n)
	realM := ctx.ComplexGetReal(m)
	imagN := ctx.ComplexGetImag(n)
	imagM := ctx.ComplexGetImag(m)

	ratio1 := imagM.Div(realM)
	denom1 := realM.Add(ratio1.Mul(imagM))
	e1 := realN.Add(imagN.Mul(ratio1)).Div(denom1)
	f1 := imagN.Sub(realN.Mul(ratio1)).Div(denom1)

	ratio2 := realM.Div(imagM)
	denom2 := imagM.Add(ratio2.Mul(realM))
	e2 := realN.Mul(ratio2).Add(imagN).Div(denom2)
	f2 := imagN.Mul(ratio2).Sub(realN).Div(denom2)

	return ctx.newComplex(
		(realN.Abs().GE(imagM.Abs())).IfThenElse(e1, e2).(z3.Float),
		(realN.Abs().GE(imagM.Abs())).IfThenElse(f1, f2).(z3.Float))
}
