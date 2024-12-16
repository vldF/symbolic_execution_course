package interpreter

import (
	"github.com/aclements/go-z3/z3"
	"symbolic_execution_course/heap"
)

type Context struct {
	Config       InterpreterConfig
	Z3Context    *z3.Context
	Solver       *z3.Solver
	TypesContext *TypesContext
	ReturnValue  Value
	States       *heap.Heap[*State]
	Results      []*State
	InitState    *State
}

type TypesContext struct {
	Int    PrimitiveIntDescr
	Int8   PrimitiveIntDescr
	Int16  PrimitiveIntDescr
	Int32  PrimitiveIntDescr
	Int64  PrimitiveIntDescr
	UInt   PrimitiveIntDescr
	UInt8  PrimitiveIntDescr
	UInt16 PrimitiveIntDescr
	UInt32 PrimitiveIntDescr
	UInt64 PrimitiveIntDescr

	Float   PrimitiveFloatDescr
	Float32 PrimitiveFloatDescr
	Float64 PrimitiveFloatDescr

	ArrayIndexSort z3.Sort

	Pointer     z3.Sort
	UnknownSort z3.Sort
}

type PrimitiveIntDescr struct {
	IsSigned bool
	Bits     int
	Sort     z3.Sort
}

type PrimitiveFloatDescr struct {
	EBits int
	MBits int
	Bits  int
	Sort  z3.Sort
}

func GetPrimitiveIntDescr(z3Ctx *z3.Context, isSigned bool, bits int) PrimitiveIntDescr {
	return PrimitiveIntDescr{
		IsSigned: isSigned,
		Bits:     bits,
		Sort:     z3Ctx.BVSort(bits),
	}
}

func GetPrimitiveFloatDescr(z3Ctx *z3.Context, eBits int, mBits int) PrimitiveFloatDescr {
	return PrimitiveFloatDescr{
		EBits: eBits,
		MBits: mBits,
		Bits:  eBits + mBits,
		Sort:  z3Ctx.FloatSort(eBits, mBits),
	}
}
