package formulas

import (
	"github.com/aclements/go-z3/z3"
)

type OrFormula struct {
	ctx    *AnalysisContext
	values []Formula
}

func NewOrFormula(ctx *AnalysisContext) OrFormula {
	return OrFormula{
		ctx:    ctx,
		values: []Formula{},
	}
}

func (f *OrFormula) Add(formula Formula) {
	f.values = append(f.values, formula)
}

func (f OrFormula) Value() z3.Bool {
	if len(f.values) == 0 {
		return f.ctx.Z3ctx.FromBool(true)
	}

	result := f.values[0].Value()

	for _, formula := range f.values[1:] {
		result = result.Or(formula.Value())
	}

	return result
}
