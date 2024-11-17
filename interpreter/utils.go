package interpreter

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"go/types"
)

func (ctx *Context) TypeToSort(t types.Type) z3.Sort {
	switch t.(type) {
	case *types.Basic:
		switch t.(*types.Basic).Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Byte:
			return ctx.TypesContext.IntSort
		case types.UntypedFloat, types.Float32, types.Float64:
			return ctx.TypesContext.FloatSort
			//case types.UntypedComplex, types.Complex64, types.Complex128:
			//	return ctx.TypesContext.ComplexSort
		}
	case *types.Array:
		elemType := t.(*types.Array).Elem()
		return ctx.Z3Context.ArraySort(ctx.TypeToSort(elemType), ctx.TypesContext.IntSort)
	case *types.Slice:
		elemType := t.(*types.Slice).Elem()
		return ctx.Z3Context.ArraySort(ctx.TypeToSort(elemType), ctx.TypesContext.IntSort)
		//case *types.Struct:
		//	return ctx.TypesContext.StructSort
	}

	return ctx.TypesContext.UnknownSort
}

func FloatToString(f z3.Float) string {
	float, _ := f.AsBigFloat()

	return fmt.Sprintf("%d", float)
}

func (ctx *Context) GoToZ3Value(v any) Z3Value {
	switch casted := v.(type) {
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		bv := ctx.Z3Context.FromInt(int64(casted.(int)), ctx.TypesContext.IntSort).(z3.BV)
		return Z3Value{
			Context: ctx,
			Value:   bv,
		}
	case float64, float32:
		float := ctx.Z3Context.FromFloat64(casted.(float64), ctx.TypesContext.FloatSort)
		return Z3Value{
			Context: ctx,
			Value:   float,
		}
	case bool:
		b := ctx.Z3Context.FromBool(casted)
		return Z3Value{
			Context: ctx,
			Value:   b,
		}
	//case complex128:
	//	return ctx.NewComplex(casted)
	//case complex64:
	//	return ctx.NewComplex(complex128(casted))
	//case []int:
	//	arrId := ctx.NewArray(ctx.Sorts.IntSort, len(casted))
	//	arr := ctx.GetArrayValue(arrId)
	//
	//	for idx, val := range casted {
	//		z3Idx := ctx.Z3ctx.FromInt(int64(idx), ctx.Sorts.IntSort)
	//		z3Val := ctx.GoToZ3Value(val)
	//		arr = arr.Store(z3Idx, z3Val)
	//	}
	//
	//	ctx.Memory.Cells[arrId].Fields[arrayField] = arr
	//	return arrId
	default:
		panic("unsupported argument")
	}
}