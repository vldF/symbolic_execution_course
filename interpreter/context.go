package interpreter

import (
	"github.com/aclements/go-z3/z3"
	"symbolic_execution_course/heap"
)

type Context struct {
	Z3Context    *z3.Context
	TypesContext *TypesContext
	ReturnValue  *Z3Value
	States       *heap.Heap[*State]
	Results      []*State
}

type TypesContext struct {
	IntBits   int
	IntSort   z3.Sort
	FloatSort z3.Sort

	UnknownSort z3.Sort
}
