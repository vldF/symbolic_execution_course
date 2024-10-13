package main

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
)

//	func compareAndIncrement(a, b int) int {
//	    if a > b {
//	        c := a + 1
//
//	        if (c > b) {
//	            return 1			(1)
//	        } else {
//	            return -1			(2)
//	        }
//	    }
//
//	    return 42					(3)
//	}
func solveSelfconstraints() {
	fmt.Println("func compareAndIncrement(a, b int) int")
	selfconstraints1()
	selfconstraints2()
}

func selfconstraints1() {
	fmt.Println("(a > b) && (c > b)")
	sCtx := CreateSymContext()

	argA := sCtx.Ctx.IntConst("a")
	argB := sCtx.Ctx.IntConst("b")

	cond1 := argA.GT(argB)
	sCtx.Solver.Assert(cond1)

	intConst1 := sCtx.Ctx.FromInt(1, sCtx.Ctx.IntSort()).(z3.Int)
	cVar := argA.Add(intConst1)
	cond2 := cVar.GT(argB)
	sCtx.Solver.Assert(cond2)
}

func selfconstraints2() {
	fmt.Println("(a > b) && !(c > b)")
	sCtx := CreateSymContext()
	sCtx.Ctx.Config().SetBool("unsat_core", true)

	assumptions := make(map[string]z3.Bool, 0)
	softConstraints := make(map[string]z3.Bool, 0)

	argA := sCtx.Ctx.IntConst("a")
	argB := sCtx.Ctx.IntConst("b")

	cond1 := argA.GT(argB)
	sCtx.Solver.Assert(cond1)

	intConst1 := sCtx.Ctx.FromInt(1, sCtx.Ctx.IntSort()).(z3.Int)
	cVar := argA.Add(intConst1)
	cond2 := cVar.GT(argB)
	assumptionVar := sCtx.Ctx.BoolConst("assumption1")
	assumptions[assumptionVar.String()] = assumptionVar
	softConstraints[assumptionVar.String()] = cond2.Not()
	// may be other assumptions
	// ...

	sCtx.Solver.Push()

	for assumptionName, assumption := range assumptions {
		constraint := softConstraints[assumptionName]
		sCtx.Solver.AssertAndTrack(constraint, assumption)
	}

	for {
		checkResult, _ := sCtx.Solver.Check()
		if checkResult {
			fmt.Println("success!")
			fmt.Println(sCtx.Solver.Model().String())
			break
		}

		if len(assumptions) == 0 {
			fmt.Println("can't satisfy the formula")
			break
		}

		unsatCore := sCtx.Solver.GetUnsatCore()
		for _, unsatBoolMarker := range unsatCore {
			assumptionName := unsatBoolMarker.String()
			fmt.Println("remove", assumptionName, "constraint")
			delete(assumptions, assumptionName)
			delete(softConstraints, assumptionName)
			sCtx.Solver.Pop()
			sCtx.Solver.Push()

			// add remaining assumptions back
			for assumptionName, assumption := range assumptions {
				constraint := softConstraints[assumptionName]
				sCtx.Solver.AssertAndTrack(constraint, assumption)
			}
		}
	}
}

func contains(arr []z3.Bool, elem z3.Bool) bool {
	for _, e := range arr {
		if e.String() == elem.String() {
			return true
		}
	}

	return false
}
