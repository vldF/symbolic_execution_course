package main

import (
	"fmt"
	"github.com/vldF/symbolic_execution_course/constraints/smt"
)

func solveComplex() {
	solveBasicComplexOperations()
	solveComplexMagnitude()
	solveComplexComparison()
	solveComplexOperations()
	solveNestedComplexOperations()
}

//	func basicComplexOperations(a complex128, b complex128) complex128 {
//		if real(a) > real(b) {
//			return a + b							(1)
//		} else if imag(a) > imag(b) {
//			return a - b							(2)
//		}
//		return a * b								(3)
//	}
func solveBasicComplexOperations() {
	fmt.Println("func basicComplexOperations(a complex128, b complex128) complex128")

	runForCase(basicComplexOperations1)
	runForCase(basicComplexOperations2)
	runForCase(basicComplexOperations3)
}

func basicComplexOperations1(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	argB := ctx.NewComplexConst("b")

	cond := argA.Real().GT(argB.Real())
	ctx.Solver.Assert(cond)

	return "(real(a) > real(b))"
}

func basicComplexOperations2(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	argB := ctx.NewComplexConst("b")

	prevCond := argA.Real().GT(argB.Real())
	ctx.Solver.Assert(prevCond.Not())

	cond := argA.Imag().GT(argB.Imag())
	ctx.Solver.Assert(cond)

	return "!(real(a) > real(b)) && (imag(a) > real(b))"
}

func basicComplexOperations3(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	argB := ctx.NewComplexConst("b")

	prevCond1 := argA.Real().GT(argB.Real())
	ctx.Solver.Assert(prevCond1.Not())

	prevCond2 := argA.Imag().GT(argB.Imag())
	ctx.Solver.Assert(prevCond2.Not())

	return "!(real(a) > real(b)) && !(imag(a) > real(b))"
}

//	func complexMagnitude(a complex128) float64 {
//		magnitude := real(a)*real(a) + imag(a)*imag(a)
//		return magnitude (1)
//	}
func solveComplexMagnitude() {
	fmt.Println("func complexMagnitude(a complex128) float64")

	runForCase(complexMagnitude1)
}

func complexMagnitude1(ctx *smt.SymContext) string {
	_ = ctx.NewComplexConst("a")
	_ = ctx.NewComplexConst("b")

	return "*"
}

//	func complexComparison(a complex128, b complex128) string {
//		magA := complexMagnitude(a)
//		magB := complexMagnitude(b)
//
//		if magA > magB {
//			return "Magnitude of a is greater than b"		(1)
//		} else if magA < magB {
//			return "Magnitude of b is greater than a"		(2)
//		}
//		return "Magnitudes are equal"						(3)
//	}

func solveComplexComparison() {
	fmt.Println("func complexComparison(a complex128, b complex128) string")
	runForCase(complexComparison1)
	runForCase(complexComparison2)
	runForCase(complexComparison3)
}

func complexComparison1(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	argB := ctx.NewComplexConst("b")

	// inlined call of complexMagnitude
	magnitudeArgA := argA.Real().Mul(argA.Real()).Add(argA.Imag().Mul(argA.Imag()))
	magnitudeArgB := argB.Real().Mul(argB.Real()).Add(argB.Imag().Mul(argB.Imag()))

	cond := magnitudeArgA.GT(magnitudeArgB)
	ctx.Solver.Assert(cond)

	return "magA > magB"
}

func complexComparison2(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	argB := ctx.NewComplexConst("b")

	// inlined call of complexMagnitude
	magnitudeArgA := argA.Real().Mul(argA.Real()).Add(argA.Imag().Mul(argA.Imag()))
	magnitudeArgB := argB.Real().Mul(argB.Real()).Add(argB.Imag().Mul(argB.Imag()))

	prevCond := magnitudeArgA.GT(magnitudeArgB)
	ctx.Solver.Assert(prevCond.Not())

	cond := magnitudeArgA.LT(magnitudeArgB)
	ctx.Solver.Assert(cond)

	return "!(magA > magB) && (magA < magB)"
}

func complexComparison3(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	argB := ctx.NewComplexConst("b")

	// inlined call of complexMagnitude
	magnitudeArgA := argA.Real().Mul(argA.Real()).Add(argA.Imag().Mul(argA.Imag()))
	magnitudeArgB := argB.Real().Mul(argB.Real()).Add(argB.Imag().Mul(argB.Imag()))

	prevCond1 := magnitudeArgA.GT(magnitudeArgB)
	ctx.Solver.Assert(prevCond1.Not())

	prevCond2 := magnitudeArgA.LT(magnitudeArgB)
	ctx.Solver.Assert(prevCond2.Not())

	return "!(magA > magB) && !(magA < magB)"
}

//	func complexOperations(a complex128, b complex128) complex128 {
//		if real(a) == 0 && imag(a) == 0 {
//			return b								(1)
//		} else if real(b) == 0 && imag(b) == 0 {
//			return a								(2)
//		} else if real(a) > real(b) {
//			return a / b							(3)
//		}
//		return a + b								(4)
//	}
func solveComplexOperations() {
	fmt.Println("func solveComplexOperations(a complex128, b complex128) complex128")

	runForCase(complexOperations1)
	runForCase(complexOperations2)
	runForCase(complexOperations3)
	runForCase(complexOperations4)
}

func complexOperations1(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	_ = ctx.NewComplexConst("b")

	realZeroConst := ctx.Ctx.FromFloat32(0.0, ctx.Ctx.FloatSort(11, 53))
	cond := argA.Real().Eq(realZeroConst).And(argA.Imag().Eq(realZeroConst))
	ctx.Solver.Assert(cond)

	return "(real(a) == 0 && imag(a) == 0)"
}

func complexOperations2(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	argB := ctx.NewComplexConst("b")

	realZeroConst := ctx.Ctx.FromFloat32(0.0, ctx.Ctx.FloatSort(11, 53))
	prevCond := argA.Real().Eq(realZeroConst).And(argA.Imag().Eq(realZeroConst))
	ctx.Solver.Assert(prevCond.Not())

	cond := argB.Real().Eq(realZeroConst).And(argB.Imag().Eq(realZeroConst))
	ctx.Solver.Assert(cond)

	return "!(real(a) == 0 && imag(a) == 0) && (real(b) == 0 && imag(b) == 0)"
}

func complexOperations3(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	argB := ctx.NewComplexConst("b")

	realZeroConst := ctx.Ctx.FromFloat32(0.0, ctx.Ctx.FloatSort(11, 53))
	prevCond1 := argA.Real().Eq(realZeroConst).And(argA.Imag().Eq(realZeroConst))
	ctx.Solver.Assert(prevCond1.Not())

	prevCond2 := argB.Real().Eq(realZeroConst).And(argB.Imag().Eq(realZeroConst))
	ctx.Solver.Assert(prevCond2.Not())

	cond := argA.Real().GT(argB.Real())
	ctx.Solver.Assert(cond)

	return "!(real(a) == 0 && imag(a) == 0) && !(real(b) == 0 && imag(b) == 0) && (real(a) > real(b))"
}

func complexOperations4(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	argB := ctx.NewComplexConst("b")

	realZeroConst := ctx.Ctx.FromFloat32(0.0, ctx.Ctx.FloatSort(11, 53))
	prevCond1 := argA.Real().Eq(realZeroConst).And(argA.Imag().Eq(realZeroConst))
	ctx.Solver.Assert(prevCond1.Not())

	prevCond2 := argB.Real().Eq(realZeroConst).And(argB.Imag().Eq(realZeroConst))
	ctx.Solver.Assert(prevCond2.Not())

	prevCond3 := argB.Real().GT(argB.Real())
	ctx.Solver.Assert(prevCond3.Not())

	return "!(real(a) == 0 && imag(a) == 0) && !(real(b) == 0 && imag(b) == 0) && !(real(a) > real(b))"
}

//	func nestedComplexOperations(a complex128, b complex128) complex128 {
//		if real(a) < 0 {
//			if imag(a) < 0 {
//				return a * b		(1)
//			}
//			return a + b			(2)
//		}
//
//		if imag(b) < 0 {
//			return a - b			(3)
//		}
//		return a + b				(4)
//	}
func solveNestedComplexOperations() {
	fmt.Println("func nestedComplexOperations(a complex128, b complex128) complex128")
	runForCase(nestedComplexOperations1)
	runForCase(nestedComplexOperations2)
	runForCase(nestedComplexOperations3)
	runForCase(nestedComplexOperations4)
}

func nestedComplexOperations1(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	_ = ctx.NewComplexConst("b")

	realZeroConst := ctx.Ctx.FromFloat32(0.0, ctx.Ctx.FloatSort(11, 53))
	cond1 := argA.Real().LT(realZeroConst)
	cond2 := argA.Imag().LT(realZeroConst)
	ctx.Solver.Assert(cond1)
	ctx.Solver.Assert(cond2)

	return "(real(a) < 0) && (imag(a) < 0)"
}

func nestedComplexOperations2(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	_ = ctx.NewComplexConst("b")

	realZeroConst := ctx.Ctx.FromFloat32(0.0, ctx.Ctx.FloatSort(11, 53))
	cond1 := argA.Real().LT(realZeroConst)
	cond2 := argA.Imag().LT(realZeroConst)
	ctx.Solver.Assert(cond1)
	ctx.Solver.Assert(cond2.Not())

	return "(real(a) < 0) && !(imag(a) < 0)"
}

func nestedComplexOperations3(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	argB := ctx.NewComplexConst("b")

	realZeroConst := ctx.Ctx.FromFloat32(0.0, ctx.Ctx.FloatSort(11, 53))
	prevCond := argA.Real().LT(realZeroConst)
	ctx.Solver.Assert(prevCond.Not())

	cond2 := argB.Imag().LT(realZeroConst)
	ctx.Solver.Assert(cond2)

	return "!(real(a) < 0) && (imag(b) < 0)"
}

func nestedComplexOperations4(ctx *smt.SymContext) string {
	argA := ctx.NewComplexConst("a")
	argB := ctx.NewComplexConst("b")

	realZeroConst := ctx.Ctx.FromFloat32(0.0, ctx.Ctx.FloatSort(11, 53))
	prevCond := argA.Real().LT(realZeroConst)
	ctx.Solver.Assert(prevCond.Not())

	cond2 := argB.Imag().LT(realZeroConst)
	ctx.Solver.Assert(cond2.Not())

	return "!(real(a) < 0) && !(imag(b) < 0)"
}
