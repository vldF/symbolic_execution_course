package main

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"github.com/vldF/symbolic_execution_course/constraints/smt"
	"math"
	"math/bits"
)

type Z3AwareFunction func(sCtx *smt.SymContext) string

func runForCase(function Z3AwareFunction) {
	sCtx := CreateSymContext()
	solver := sCtx.Solver

	caseName := function(&sCtx)
	fmt.Println("===================")
	fmt.Println(caseName)

	fmt.Println("constraints: ", solver.String())

	check, err := solver.Check()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("is satisfied: ", check)
	if !check {
		unsatCore := sCtx.Solver.GetUnsatCore()
		for i := range unsatCore {
			fmt.Println(unsatCore[i])
		}

		return
	}

	fmt.Println(solver.Model().String())
}

func CreateSymContext() smt.SymContext {
	config := z3.Config{}
	ctx := z3.NewContext(&config)
	solver := z3.NewSolver(ctx)

	typesCtx := smt.TypesContext{
		MaxInt:      math.MaxInt32,
		MinInt:      math.MinInt32,
		IntSize:     bits.UintSize,
		MaxFloat64:  math.MaxFloat64,
		MinFloat64:  -math.MaxFloat64,
		Float64Size: 64,
	}

	sCtx := smt.SymContext{
		Solver:   solver,
		Ctx:      ctx,
		TypesCtx: typesCtx,
	}
	return sCtx
}
