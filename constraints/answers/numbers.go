package main

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"github.com/vldF/symbolic_execution_course/constraints/smt"
)

func solveNumbers() {
	solveIntegerOperations()
	solveFloatOperations()
	solveMixedOperations()
	solveNestedConditions()
	solveBitwiseOperations()
	solveAdvancedBitwise()
	solveCombinedBitwise()
	solveNestedBitwise()
}

//	func integerOperations(a int, b int) int {
//		if a > b {
//			return a + b  							(1)
//		} else if a < b {
//			return a - b							(2)
//		} else {
//			return a * b							(3)
//		}
//	}
func solveIntegerOperations() {
	fmt.Println("func integerOperations(a int, b int)")

	runForCase(integerOperators1)
	runForCase(integerOperators2)
	runForCase(integerOperators3)
}

func integerOperators1(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	cond := argA.GT(argB)
	sCtx.Solver.Assert(cond)

	return "if a > b"
}

func integerOperators2(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	prevCond := argA.GT(argB)
	sCtx.Solver.Assert(prevCond.Not())

	cond := argA.LT(argB)
	sCtx.Solver.Assert(cond)

	return "if a < b"
}

func integerOperators3(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	prevCond1 := argA.GT(argB)
	sCtx.Solver.Assert(prevCond1.Not())

	prevCond2 := argA.LT(argB)
	sCtx.Solver.Assert(prevCond2.Not())

	return "else"
}

//	func floatOperations(x float64, y float64) float64 {
//		if x > y {
//			return x / y							(1)
//		} else if x < y {
//			return x * y							(2)
//		}
//		return 0.0									(3)
//	}
func solveFloatOperations() {
	fmt.Println("func floatOperations(x float64, y float64) float64")
	runForCase(floatOperations1)
	runForCase(floatOperations2)
	runForCase(floatOperations3)
}

func floatOperations1(sCtx *smt.SymContext) string {
	argX := sCtx.NewFloat64Argument("x")
	argY := sCtx.NewFloat64Argument("y")

	cond := argX.GT(argY)
	sCtx.Solver.Assert(cond)

	return "if x > y"
}

func floatOperations2(sCtx *smt.SymContext) string {
	argX := sCtx.NewFloat64Argument("x")
	argY := sCtx.NewFloat64Argument("y")

	prevCond := argX.GT(argY)
	sCtx.Solver.Assert(prevCond.Not())

	cond := argX.LT(argY)
	sCtx.Solver.Assert(cond)

	return "if x < y"
}

func floatOperations3(sCtx *smt.SymContext) string {
	argX := sCtx.NewFloat64Argument("x")
	argY := sCtx.NewFloat64Argument("y")

	prevCond1 := argX.GT(argY)
	sCtx.Solver.Assert(prevCond1.Not())

	prevCond2 := argX.LT(argY)
	sCtx.Solver.Assert(prevCond2.Not())

	return "else"
}

//	func mixedOperations(a int, b float64) float64 {
//		var result float64
//
//		if a%2 == 0 {
//			result = float64(a) + b			(1)
//		} else {
//			result = float64(a) - b			(2)
//		}
//
//		if result < 10 {
//			result *= 2						(3)
//		} else {
//			result /= 2						(4)
//		}
//
//		return result
//	}

func solveMixedOperations() {
	fmt.Println("func mixedOperations(a int, b float64) float64")
	runForCase(mixedOperations13)
	runForCase(mixedOperations14)
	runForCase(mixedOperations23)
	runForCase(mixedOperations24)
}

func mixedOperations13(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewFloat64Argument("b")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	twoIntConst := sCtx.Ctx.FromInt(2, sCtx.Ctx.IntSort()).(z3.Int)
	sCtx.Solver.Assert(argA.Mod(twoIntConst).Eq(zeroIntConst))

	floatSort := sCtx.Ctx.FloatSort(11, 53)
	result := argA.ToReal().ToFloat(floatSort).Add(argB)
	tenRealConst := sCtx.Ctx.FromFloat64(10.0, sCtx.Ctx.FloatSort(11, 53))
	sCtx.Solver.Assert(result.LT(tenRealConst))

	return "a%2 == 0 && result < 10"
}

func mixedOperations14(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewFloat64Argument("b")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	twoIntConst := sCtx.Ctx.FromInt(2, sCtx.Ctx.IntSort()).(z3.Int)
	sCtx.Solver.Assert(argA.Mod(twoIntConst).Eq(zeroIntConst))

	floatSort := sCtx.Ctx.FloatSort(11, 53)
	result := argA.ToReal().ToFloat(floatSort).Add(argB)
	tenRealConst := sCtx.Ctx.FromFloat64(10.0, sCtx.Ctx.FloatSort(11, 53))
	sCtx.Solver.Assert(result.LT(tenRealConst).Not())

	return "a%2 == 0 && !(result < 10)"
}

func mixedOperations23(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewFloat64Argument("b")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	twoIntConst := sCtx.Ctx.FromInt(2, sCtx.Ctx.IntSort()).(z3.Int)
	sCtx.Solver.Assert(argA.Mod(twoIntConst).Eq(zeroIntConst).Not())

	floatSort := sCtx.Ctx.FloatSort(11, 53)
	result := argA.ToReal().ToFloat(floatSort).Sub(argB)
	tenRealConst := sCtx.Ctx.FromFloat64(10.0, sCtx.Ctx.FloatSort(11, 53))
	sCtx.Solver.Assert(result.LT(tenRealConst))

	return "!(a%2 == 0) && result < 10"
}

func mixedOperations24(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewFloat64Argument("b")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	twoIntConst := sCtx.Ctx.FromInt(2, sCtx.Ctx.IntSort()).(z3.Int)
	sCtx.Solver.Assert(argA.Mod(twoIntConst).Eq(zeroIntConst).Not())

	floatSort := sCtx.Ctx.FloatSort(11, 53)
	result := argA.ToReal().ToFloat(floatSort).Sub(argB)
	tenRealConst := sCtx.Ctx.FromFloat64(10.0, sCtx.Ctx.FloatSort(11, 53))
	sCtx.Solver.Assert(result.LT(tenRealConst).Not())

	return "!(a%2 == 0) && !(result < 10)"
}

//	func nestedConditions(a int, b float64) float64 {
//		if a < 0 {
//			if b < 0 {
//				return float64(a*-1) + b 	(1)
//			}
//			return float64(a*-1) - b		(2)
//		}
//		return float64(a) + b				(3)
//	}

func solveNestedConditions() {
	fmt.Println("func nestedConditions(a int, b float64) float64")
	runForCase(nestedConditions1)
	runForCase(nestedConditions2)
	runForCase(nestedConditions3)
}

func nestedConditions1(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewFloat64Argument("b")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	zeroRealConst := sCtx.Ctx.FromFloat64(0.0, sCtx.Ctx.FloatSort(11, 53))
	sCtx.Solver.Assert(argA.LT(zeroIntConst))
	sCtx.Solver.Assert(argB.LT(zeroRealConst))

	return "a < 0 && b < 0.0"
}

func nestedConditions2(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewFloat64Argument("b")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	zeroRealConst := sCtx.Ctx.FromFloat64(0.0, sCtx.Ctx.FloatSort(11, 53))
	sCtx.Solver.Assert(argA.LT(zeroIntConst))
	sCtx.Solver.Assert(argB.LT(zeroRealConst).Not())

	return "a < 0 && !(b < 0.0)"
}

func nestedConditions3(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	_ = sCtx.NewFloat64Argument("b")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	sCtx.Solver.Assert(argA.LT(zeroIntConst).Not())

	return "!(a < 0)"
}

//	func bitwiseOperations(a int, b int) int {
//		if a&1 == 0 && b&1 == 0 {
//			return a | b					(1)
//		} else if a&1 == 1 && b&1 == 1 {
//			return a & b					(2)
//		}
//		return a ^ b						(3)
//	}
func solveBitwiseOperations() {
	fmt.Println("func bitwiseOperations(a int, b int) int")
	runForCase(bitwiseOperations1)
	runForCase(bitwiseOperations2)
	runForCase(bitwiseOperations3)
}

func bitwiseOperations1(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	oneBVConst := sCtx.Ctx.FromInt(1, sCtx.Ctx.IntSort()).(z3.Int).ToBV(32)
	zeroBVConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int).ToBV(32)

	sCtx.Solver.Assert(argA.ToBV(32).And(oneBVConst).Eq(zeroBVConst).And(argB.ToBV(32).And(oneBVConst).Eq(zeroBVConst)))

	return "a&1 == 0 && b&1 == 0"
}

func bitwiseOperations2(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	oneBVConst := sCtx.Ctx.FromInt(1, sCtx.Ctx.IntSort()).(z3.Int).ToBV(32)
	zeroBVConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int).ToBV(32)

	prevCond := argA.ToBV(32).And(oneBVConst).Eq(zeroBVConst).And(argB.ToBV(32).And(oneBVConst).Eq(zeroBVConst))
	sCtx.Solver.Assert(prevCond.Not())

	newCond := argA.ToBV(32).And(oneBVConst).Eq(oneBVConst).And(argB.ToBV(32).And(oneBVConst).Eq(oneBVConst))
	sCtx.Solver.Assert(newCond)

	return "a&1 == 0 && b&1 == 0"
}

func bitwiseOperations3(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	oneBVConst := sCtx.Ctx.FromInt(1, sCtx.Ctx.IntSort()).(z3.Int).ToBV(32)
	zeroBVConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int).ToBV(32)

	prevCond1 := argA.ToBV(32).And(oneBVConst).Eq(zeroBVConst).And(argB.ToBV(32).And(oneBVConst).Eq(zeroBVConst))
	sCtx.Solver.Assert(prevCond1.Not())

	prevCond2 := argA.ToBV(32).And(oneBVConst).Eq(oneBVConst).And(argB.ToBV(32).And(oneBVConst).Eq(oneBVConst))
	sCtx.Solver.Assert(prevCond2.Not())

	return "a&1 == 0 && b&1 == 0"
}

//	func advancedBitwise(a int, b int) int {
//		if a > b {
//			return a << 1				(1)
//		} else if a < b {
//			return b >> 1				(2)
//		}
//		return a ^ b					(3)
//	}
func solveAdvancedBitwise() {
	fmt.Println("func advancedBitwise(a int, b int) int")
	runForCase(advancedBitwise1)
	runForCase(advancedBitwise2)
	runForCase(advancedBitwise3)
}

func advancedBitwise1(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	cond := argA.GT(argB)
	sCtx.Solver.Assert(cond)

	return "a > b"
}

func advancedBitwise2(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	prevCond := argA.GT(argB)
	sCtx.Solver.Assert(prevCond.Not())

	cond := argA.LT(argB)
	sCtx.Solver.Assert(cond)

	return "!(a > b) && (a < b)"
}

func advancedBitwise3(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	prevCond1 := argA.GT(argB)
	sCtx.Solver.Assert(prevCond1.Not())

	prevCond2 := argA.LT(argB)
	sCtx.Solver.Assert(prevCond2.Not())

	return "!(a > b) && !(a < b)"
}

//	func combinedBitwise(a int, b int) int {
//		if a&b == 0 {
//			return a | b				(1)
//		} else {
//			result := a & b
//			if result > 10 {
//				return result ^ b		(2)
//			}
//			return result				(3)
//		}
//	}
func solveCombinedBitwise() {
	fmt.Println("func combinedBitwise(a int, b int) int")
	runForCase(combinedBitwise1)
	runForCase(combinedBitwise2)
	runForCase(combinedBitwise3)
}

func combinedBitwise1(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	zeroBVConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int).ToBV(32)

	cond := argA.ToBV(32).And(argB.ToBV(32)).Eq(zeroBVConst)
	sCtx.Solver.Assert(cond)

	return "a&b == 0"
}

func combinedBitwise2(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	zeroBVConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int).ToBV(32)

	argAAndArgB := argA.ToBV(32).And(argB.ToBV(32))
	prevCond := argAAndArgB.Eq(zeroBVConst)
	sCtx.Solver.Assert(prevCond.Not())

	const10 := sCtx.Ctx.FromInt(10, sCtx.Ctx.IntSort()).(z3.Int)
	cond := argAAndArgB.SToInt().GT(const10)
	sCtx.Solver.Assert(cond)

	return "!(a&b == 0) && (a&b) > 10"
}

func combinedBitwise3(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	zeroBVConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int).ToBV(32)

	argAAndArgB := argA.ToBV(32).And(argB.ToBV(32))
	prevCond1 := argAAndArgB.Eq(zeroBVConst)
	sCtx.Solver.Assert(prevCond1.Not())

	const10 := sCtx.Ctx.FromInt(10, sCtx.Ctx.IntSort()).(z3.Int)
	prevCond2 := argAAndArgB.SToInt().GT(const10)
	sCtx.Solver.Assert(prevCond2.Not())

	return "!(a&b == 0) && !((a&b) > 10)"
}

//	func nestedBitwise(a int, b int) int {
//		if a < 0 {
//			return -1 				(1)
//		}
//
//		if b < 0 {
//			return a ^ 0			(2)
//		}
//
//		if a&b == 0 {
//			return a | b			(3)
//		} else {
//			return a & b			(4)
//		}
//	}
func solveNestedBitwise() {
	fmt.Println("func nestedBitwise(a int, b int) int")
	runForCase(nestedBitwise1)
	runForCase(nestedBitwise2)
	runForCase(nestedBitwise3)
	runForCase(nestedBitwise4)
}

func nestedBitwise1(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	_ = sCtx.NewIntArgument("b")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)

	cond := argA.LT(zeroIntConst)
	sCtx.Solver.Assert(cond)

	return "a < 0"
}

func nestedBitwise2(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)

	prevCond := argA.LT(zeroIntConst)
	sCtx.Solver.Assert(prevCond.Not())

	cond := argB.LT(zeroIntConst)
	sCtx.Solver.Assert(cond)

	return "!(a < 0) && (b < 0)"
}

func nestedBitwise3(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)

	prevCond1 := argA.LT(zeroIntConst)
	sCtx.Solver.Assert(prevCond1.Not())

	prevCond2 := argB.LT(zeroIntConst)
	sCtx.Solver.Assert(prevCond2.Not())

	typesCtx := sCtx.TypesCtx
	cond := argA.ToBV(typesCtx.IntSize).And(argB.ToBV(typesCtx.IntSize)).Eq(zeroIntConst.ToBV(typesCtx.IntSize))
	sCtx.Solver.Assert(cond)

	return "!(a < 0) && !(b < 0) && (a&b == 0)"
}

func nestedBitwise4(sCtx *smt.SymContext) string {
	argA := sCtx.NewIntArgument("a")
	argB := sCtx.NewIntArgument("b")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)

	prevCond1 := argA.LT(zeroIntConst)
	sCtx.Solver.Assert(prevCond1.Not())

	prevCond2 := argB.LT(zeroIntConst)
	sCtx.Solver.Assert(prevCond2.Not())

	typesCtx := sCtx.TypesCtx
	prevCond3 := argA.ToBV(typesCtx.IntSize).And(argB.ToBV(typesCtx.IntSize)).Eq(zeroIntConst.ToBV(typesCtx.IntSize))
	sCtx.Solver.Assert(prevCond3.Not())

	return "!(a < 0) && !(b < 0) && !(a&b == 0)"
}
