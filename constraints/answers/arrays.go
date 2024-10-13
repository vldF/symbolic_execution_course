package main

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"github.com/vldF/symbolic_execution_course/constraints/smt"
)

func solveArrays() {
	solveCompareElements()
	solveCompareAges()
}

//	func compareElement(array []int, index int, value int) int {
//	   if index < 0 || index >= len(array) {
//	       return -1 // Индекс вне границ		(1)
//	   }
//	   element := array[index]
//	   if element > value {
//	       return 1 // Элемент больше			(2)
//	   } else if element < value {
//	       return -1 // Элемент меньше			(3)
//	   }
//	   return 0 // Элемент равен				(4)
//	}
func solveCompareElements() {
	fmt.Println("func compareElement(array []int, index int, value int) int")
	runForCase(compareElement1)
	runForCase(compareElement2)
	runForCase(compareElement3)
	runForCase(compareElement4)
}

func compareElement1(sCtx *smt.SymContext) string {
	argArray := sCtx.NewIntArray("array")
	argIndex := sCtx.Ctx.IntConst("index")
	_ = sCtx.Ctx.IntConst("value")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	cond := argIndex.LT(zeroIntConst).Or(argIndex.GT(argArray.Len()).Or(argIndex.Eq(argArray.Len())))
	sCtx.Solver.Assert(cond)

	return "(index < 0 || index >= len(array))"
}

func compareElement2(sCtx *smt.SymContext) string {
	argArray := sCtx.NewIntArray("array")
	argIndex := sCtx.Ctx.IntConst("index")
	argValue := sCtx.Ctx.IntConst("value")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	prevCond1 := argIndex.LT(zeroIntConst).Or(argIndex.GT(argArray.Len()).Or(argIndex.Eq(argArray.Len())))
	sCtx.Solver.Assert(prevCond1.Not())

	elementVar := argArray.Arr().Select(argIndex).(z3.Int)
	cond := elementVar.GT(argValue)
	sCtx.Solver.Assert(cond)

	return "!(index < 0 || index >= len(array)) && (array[index] > argValue)"
}

func compareElement3(sCtx *smt.SymContext) string {
	argArray := sCtx.NewIntArray("array")
	argIndex := sCtx.Ctx.IntConst("index")
	argValue := sCtx.Ctx.IntConst("value")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	prevCond1 := argIndex.LT(zeroIntConst).Or(argIndex.GT(argArray.Len()).Or(argIndex.Eq(argArray.Len())))
	sCtx.Solver.Assert(prevCond1.Not())

	elementVar := argArray.Arr().Select(argIndex).(z3.Int)
	cond := elementVar.LT(argValue)
	sCtx.Solver.Assert(cond)

	return "!(index < 0 || index >= len(array)) && !(array[index] > argValue) && (array[index] < argValue)"
}

func compareElement4(sCtx *smt.SymContext) string {
	argArray := sCtx.NewIntArray("array")
	argIndex := sCtx.Ctx.IntConst("index")
	argValue := sCtx.Ctx.IntConst("value")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	prevCond1 := argIndex.LT(zeroIntConst).Or(argIndex.GT(argArray.Len()).Or(argIndex.Eq(argArray.Len())))
	sCtx.Solver.Assert(prevCond1.Not())

	elementVar := argArray.Arr().Select(argIndex).(z3.Int)
	prevCond2 := elementVar.LT(argValue)
	sCtx.Solver.Assert(prevCond2.Not())

	return "!(index < 0 || index >= len(array)) && !(array[index] > argValue) && !(array[index] > argValue)"
}

//	type Person struct {
//		Name string
//		Age  int
//	}
//
//	func compareAge(people []*Person, index int, value int) int {
//		if index < 0 || index >= len(people) {
//			return -1 // Индекс вне границ							(1)
//		}
//		age := people[index].Age // Достаем возраст по индексу
//
//		if age > value {
//			return 1 // Возраст больше								(2)
//		} else if age < value {
//			return -1 // Возраст меньше								(3)
//		}
//		return 0 // Возраст равен									(4)
//	}
func solveCompareAges() {
	fmt.Println("func compareAges")

	runForCase(compareAge1)
	runForCase(compareAge2)
	runForCase(compareAge3)
	runForCase(compareAge4)
}

func compareAge1(sCtx *smt.SymContext) string {
	personStructDescriptor := map[string]z3.Sort{
		"Name": sCtx.Ctx.UninterpretedSort("string"),
		"Age":  sCtx.Ctx.IntSort(),
	}

	argPeople := sCtx.NewStructArray("people", personStructDescriptor)
	argIndex := sCtx.Ctx.IntConst("index")
	_ = sCtx.Ctx.IntConst("value")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	cond := argIndex.LT(zeroIntConst).Or(argIndex.GT(argPeople.Len()).Or(argIndex.Eq(argPeople.Len())))
	sCtx.Solver.Assert(cond)

	return "(index < 0 || index >= len(people))"
}

func compareAge2(sCtx *smt.SymContext) string {
	personStructDescriptor := map[string]z3.Sort{
		"Name": sCtx.Ctx.UninterpretedSort("string"),
		"Age":  sCtx.Ctx.IntSort(),
	}

	argPeople := sCtx.NewStructArray("people", personStructDescriptor)
	argIndex := sCtx.Ctx.IntConst("index")
	argValue := sCtx.Ctx.IntConst("value")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	prevCond := argIndex.LT(zeroIntConst).Or(argIndex.GT(argPeople.Len()).Or(argIndex.Eq(argPeople.Len())))
	sCtx.Solver.Assert(prevCond.Not())

	cond := argPeople.GetStructure(argIndex)["Age"].(z3.Int).GT(argValue)
	sCtx.Solver.Assert(cond)

	return "!(index < 0 || index >= len(people)) && (age > value)"
}

func compareAge3(sCtx *smt.SymContext) string {
	personStructDescriptor := map[string]z3.Sort{
		"Name": sCtx.Ctx.UninterpretedSort("string"),
		"Age":  sCtx.Ctx.IntSort(),
	}

	argPeople := sCtx.NewStructArray("people", personStructDescriptor)
	argIndex := sCtx.Ctx.IntConst("index")
	argValue := sCtx.Ctx.IntConst("value")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	prevCond1 := argIndex.LT(zeroIntConst).Or(argIndex.GT(argPeople.Len()).Or(argIndex.Eq(argPeople.Len())))
	sCtx.Solver.Assert(prevCond1.Not())

	prevCond2 := argPeople.GetStructure(argIndex)["Age"].(z3.Int).GT(argValue)
	sCtx.Solver.Assert(prevCond2.Not())

	cond := argPeople.GetStructure(argIndex)["Age"].(z3.Int).LT(argValue)
	sCtx.Solver.Assert(cond)

	return "!(index < 0 || index >= len(people)) && !(age > value) && (age < value)"
}

func compareAge4(sCtx *smt.SymContext) string {
	personStructDescriptor := map[string]z3.Sort{
		"Name": sCtx.Ctx.UninterpretedSort("string"),
		"Age":  sCtx.Ctx.IntSort(),
	}

	argPeople := sCtx.NewStructArray("people", personStructDescriptor)
	argIndex := sCtx.Ctx.IntConst("index")
	argValue := sCtx.Ctx.IntConst("value")

	zeroIntConst := sCtx.Ctx.FromInt(0, sCtx.Ctx.IntSort()).(z3.Int)
	prevCond1 := argIndex.LT(zeroIntConst).Or(argIndex.GT(argPeople.Len()).Or(argIndex.Eq(argPeople.Len())))
	sCtx.Solver.Assert(prevCond1.Not())

	prevCond2 := argPeople.GetStructure(argIndex)["Age"].(z3.Int).GT(argValue)
	sCtx.Solver.Assert(prevCond2.Not())

	prevCond3 := argPeople.GetStructure(argIndex)["Age"].(z3.Int).LT(argValue)
	sCtx.Solver.Assert(prevCond3.Not())

	return "!(index < 0 || index >= len(people)) && !(age > value) && !(age < value)"
}
