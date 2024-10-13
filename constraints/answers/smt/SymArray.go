package smt

import "github.com/aclements/go-z3/z3"

type SymSimpleArray struct {
	len z3.Int
	arr z3.Array
}

func (sCtx *SymContext) NewIntArray(name string) SymSimpleArray {
	lenVal := sCtx.Ctx.IntConst(name + "." + "len")
	zeroConst := sCtx.Ctx.FromInt(-1, sCtx.Ctx.IntSort()).(z3.Int)
	sCtx.Solver.Assert(lenVal.GT(zeroConst))

	arrSort := sCtx.Ctx.ArraySort(sCtx.Ctx.IntSort(), sCtx.Ctx.IntSort())
	arr := sCtx.Ctx.Const(name+"."+"array", arrSort).(z3.Array)

	return SymSimpleArray{len: lenVal, arr: arr}
}

func (arr *SymSimpleArray) Len() z3.Int {
	return arr.len
}

func (arr *SymSimpleArray) Arr() z3.Array {
	return arr.arr
}

type SymStructArray struct {
	len    z3.Int
	arrays map[string]z3.Array
}

type SymStructure = map[string]z3.Value

func (sCtx *SymContext) NewStructArray(name string, elementDesc map[string]z3.Sort) SymStructArray {
	lenVal := sCtx.Ctx.IntConst(name + "." + "len")
	zeroConst := sCtx.Ctx.FromInt(-1, sCtx.Ctx.IntSort()).(z3.Int)
	sCtx.Solver.Assert(lenVal.GT(zeroConst))

	arrays := make(map[string]z3.Array)

	for fieldName, fieldSort := range elementDesc {
		arrSort := sCtx.Ctx.ArraySort(sCtx.Ctx.IntSort(), fieldSort)
		arrays[fieldName] = sCtx.Ctx.Const(name+"."+fieldName+".array", arrSort).(z3.Array)
	}

	return SymStructArray{len: lenVal, arrays: arrays}
}

func (arr *SymStructArray) Len() z3.Int {
	return arr.len
}

func (arr *SymStructArray) GetStructure(index z3.Int) SymStructure {
	resultFields := make(SymStructure)

	for fieldName, innerArray := range arr.arrays {
		resultFields[fieldName] = innerArray.Select(index)
	}

	return resultFields
}
