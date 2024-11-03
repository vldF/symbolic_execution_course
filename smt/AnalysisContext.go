package smt

import (
	"github.com/aclements/go-z3/z3"
	"golang.org/x/tools/go/ssa"
	"slices"
	"symbolic_execution_course/smt/memory"
)

type AnalysisContext struct {
	Z3ctx *z3.Context
	Sorts Sorts

	Constraints []Formula
	Args        map[string]z3.Value
	ResultValue z3.Value

	basicBlockHistory []*ssa.BasicBlock
	Memory            memory.Memory
}

type Sorts struct {
	IntSort     z3.Sort
	FloatSort   z3.Sort
	UnknownSort z3.Sort
	ComplexSort z3.Sort
	SymPtrSort  z3.Sort
	StructSort  z3.Sort
}

func (ctx *AnalysisContext) PushBasicBlock(bb *ssa.BasicBlock) {
	ctx.basicBlockHistory = append(ctx.basicBlockHistory, bb)
}

func (ctx *AnalysisContext) PopBasicBlock() {
	if len(ctx.basicBlockHistory) == 0 {
		return
	}

	ctx.basicBlockHistory = ctx.basicBlockHistory[:len(ctx.basicBlockHistory)-1]
}

func (ctx *AnalysisContext) HasBasicBlockInHistory(bb *ssa.BasicBlock) bool {
	return slices.Contains(ctx.basicBlockHistory, bb)
}
