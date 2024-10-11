package smt

import "github.com/aclements/go-z3/z3"

type SymComplex struct {
	re z3.Float
	im z3.Float
}

func (sCtx *SymContext) NewComplexConst(name string) SymComplex {
	reConst := sCtx.Ctx.Const(name+".real", sCtx.Ctx.FloatSort(11, 53)).(z3.Float)
	imConst := sCtx.Ctx.Const(name+".imag", sCtx.Ctx.FloatSort(11, 53)).(z3.Float)

	sCtx.Solver.Assert(reConst.IsZero().And(reConst.IsPositive()).Not())
	sCtx.Solver.Assert(imConst.IsZero().And(imConst.IsPositive()).Not())

	sCtx.Solver.Assert(reConst.IsNaN().Not())
	sCtx.Solver.Assert(imConst.IsNaN().Not())

	return SymComplex{
		reConst,
		imConst,
	}
}

func (complex SymComplex) Real() z3.Float { return complex.re }
func (complex SymComplex) Imag() z3.Float { return complex.im }