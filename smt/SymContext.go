package smt

import "github.com/aclements/go-z3/z3"

type SymContext struct {
	Solver   *z3.Solver
	Ctx      *z3.Context
	TypesCtx TypesContext
}

type TypesContext struct {
	MaxInt  int64
	MinInt  int64
	IntSize int

	MaxFloat64  float64
	MinFloat64  float64
	Float64Size int64
}
