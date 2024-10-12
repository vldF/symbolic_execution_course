package main

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
)

//	func pushPopIncrementality(j int) int {
//	    result := j
//
//	    for i := 1; i <= 10; i++ {
//	        result += i
//	    }
//
//	    if result%2 == 0 {
//	        result++
//	    }
//	}
func solvePushPop() {
	fmt.Println("func pushPopIncrementality(j int) int")
	sCtx := CreateSymContext()
	solver := sCtx.Solver
	ctx := sCtx.Ctx

	argJ := ctx.IntConst("j")
	resultVar := ctx.IntConst("result")
	solver.Assert(resultVar.Eq(argJ))

	// loop unrolling
	for i := 1; i <= 10; i++ {
		intConst := ctx.FromInt(int64(i), ctx.IntSort()).(z3.Int)
		resultVar = resultVar.Add(intConst)
	}

	solver.Push()

	// encode the state inside the 'if'
	intConst0 := ctx.FromInt(0, ctx.IntSort()).(z3.Int)
	intConst2 := ctx.FromInt(2, ctx.IntSort()).(z3.Int)
	solver.Assert(resultVar.Mod(intConst2).Eq(intConst0))
	res1, err1 := solver.Check()
	if err1 != nil {
		panic(err1)
	}

	fmt.Println("is satisfiable:", res1)
	fmt.Println("formula is:", solver.String())
	fmt.Println(solver.Model().String())

	solver.Pop()
	// encode the state outside the 'if'
	solver.Assert(resultVar.Mod(intConst2).NE(intConst0))
	res2, err2 := solver.Check()
	if err2 != nil {
		panic(err2)
	}

	fmt.Println("is satisfiable:", res2)
	fmt.Println("formula is:", solver.String())
	fmt.Println(solver.Model().String())
}
