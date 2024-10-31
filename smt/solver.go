package smt

import "github.com/aclements/go-z3/z3"

func CreateSolver() (*z3.Solver, *z3.Context) {
	config := z3.Config{}
	ctx := z3.NewContext(&config)
	return z3.NewSolver(ctx), ctx
}
