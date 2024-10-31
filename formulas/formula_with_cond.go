package formulas

import (
	"github.com/aclements/go-z3/z3"
)

type FormulaWithCond struct {
	ctx    *AnalysisContext
	cond   z3.Bool
	values []Formula
}

func NewFormulaWithCond(ctx *AnalysisContext, cond z3.Bool) FormulaWithCond {
	return FormulaWithCond{
		ctx:    ctx,
		cond:   cond,
		values: make([]Formula, 0),
	}
}

func (f *FormulaWithCond) Add(formula Formula) {
	f.values = append(f.values, formula)
}

func (f FormulaWithCond) Value() z3.Bool {
	result := f.cond
	for _, value := range f.values {
		result = result.And(value.Value())
	}

	return result
}
