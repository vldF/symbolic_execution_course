package smt

import (
	"github.com/aclements/go-z3/z3"
	"symbolic_execution_course/smt/memory"
)

const (
	arrayField    = "array"
	arrayLenField = "length"
)

func (ctx *AnalysisContext) NewArray(sort z3.Sort, length int) *memory.SymMemoryPtr {
	id := ctx.Z3ctx.FreshConst("arr-wrapper", ctx.Z3ctx.UninterpretedSort("array-id")).(memory.SymMemoryPtr)

	indexSort := ctx.Sorts.IntSort
	arraySort := ctx.Z3ctx.ArraySort(indexSort, sort)

	arrayConst := ctx.Z3ctx.FreshConst("array-content", arraySort)

	var lenConst z3.Value
	if length >= 0 {
		lenConst = ctx.Z3ctx.FromInt(int64(length), indexSort)
	} else {
		lenConst = ctx.Z3ctx.FreshConst("array-length", indexSort)
	}

	symCell := memory.SymMemoryCell{
		Kind: memory.Array,
		Fields: map[string]z3.Value{
			arrayField:    arrayConst,
			arrayLenField: lenConst,
		},
	}
	ctx.Memory.Cells[&id] = &symCell

	return &id
}

func (ctx *AnalysisContext) GetArrayValue(id *z3.Uninterpreted) z3.Array {
	v := ctx.Memory.Cells[id].Fields[arrayField]
	return v.(z3.Array)
}

func (ctx *AnalysisContext) GetArrayLen(id *z3.Uninterpreted) z3.Int {
	v := ctx.Memory.Cells[id].Fields[arrayLenField]
	return v.(z3.Int)
}
