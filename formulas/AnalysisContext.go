package formulas

import (
	"github.com/aclements/go-z3/z3"
)

type AnalysisContext struct {
	Z3ctx *z3.Context
	Sorts Sorts

	Constraints []Formula
	ResultValue z3.Value
}

type Sorts struct {
	IntSort     z3.Sort
	FloatSort   z3.Sort
	UnknownSort z3.Sort
}
