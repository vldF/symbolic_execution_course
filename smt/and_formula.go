package smt

import (
	"github.com/aclements/go-z3/z3"
)

type AndFormula struct {
	ctx    *AnalysisContext
	values []Formula
}

func NewAndFormula(ctx *AnalysisContext) AndFormula {
	return AndFormula{
		ctx:    ctx,
		values: []Formula{},
	}
}

func (f *AndFormula) Add(formula Formula) {
	f.values = append(f.values, formula)
}

func (f AndFormula) Value() z3.Bool {
	if len(f.values) == 0 {
		return f.ctx.Z3ctx.FromBool(true)
	}

	result := f.values[0].Value()

	for _, formula := range f.values[1:] {
		result = result.And(formula.Value())
	}

	return result
}
