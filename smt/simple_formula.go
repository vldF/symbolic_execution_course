package smt

import (
	"github.com/aclements/go-z3/z3"
)

type SimpleFormula struct {
	ctx    *AnalysisContext
	values []z3.Bool
}

func NewSimpleFormula(ctx *AnalysisContext) SimpleFormula {
	return SimpleFormula{
		ctx:    ctx,
		values: make([]z3.Bool, 0),
	}
}

func (f *SimpleFormula) Add(value z3.Bool) {
	f.values = append(f.values, value)
}

func (f SimpleFormula) Value() z3.Bool {
	if len(f.values) == 0 {
		return f.ctx.Z3ctx.FromBool(true)
	}

	result := f.values[0]
	for _, value := range f.values[1:] {
		result = result.And(value)
	}

	return result
}
