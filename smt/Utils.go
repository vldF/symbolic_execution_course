package smt

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"go/types"
)

func (ctx *AnalysisContext) Eq(left z3.Value, right z3.Value) z3.Bool {
	switch left.(type) {
	case z3.Int:
		return left.(z3.Int).Eq(right.(z3.Int))
	case z3.Bool:
		return left.(z3.Bool).Eq(right.(z3.Bool))
	case z3.Float:
		return left.(z3.Float).Eq(right.(z3.Float))
	}

	switch {
	case ctx.IsComplex(left) && ctx.IsComplex(right):
		return ctx.ComplexEq(left.(Z3Complex), right.(Z3Complex))
	}

	panic("can't build Eq")
}

func (ctx *AnalysisContext) Ne(left z3.Value, right z3.Value) z3.Bool {
	return ctx.Eq(left, right).Not()
}

func (ctx *AnalysisContext) Add(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.Add(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.Add(second) },
		ctx.ComplexAdd)
}

func (ctx *AnalysisContext) Sub(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.Sub(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.Sub(second) },
		ctx.ComplexSub)
}

func (ctx *AnalysisContext) Mul(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.Mul(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.Mul(second) },
		ctx.ComplexMul)
}

func (ctx *AnalysisContext) Div(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.Div(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.Div(second) },
		nil)
}

func (ctx *AnalysisContext) Lt(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.LT(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.LT(second) },
		nil)
}

func (ctx *AnalysisContext) Le(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.LE(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.LE(second) },
		nil)
}

func (ctx *AnalysisContext) Gt(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.GT(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.GT(second) },
		nil)
}

func (ctx *AnalysisContext) Ge(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.GE(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.GE(second) },
		nil)
}

func (ctx *AnalysisContext) arithOp(
	left z3.Value,
	right z3.Value,
	floatOp func(z3.Float, z3.Float) z3.Value,
	intOp func(z3.Int, z3.Int) z3.Value,
	complexOp func(Z3Complex, Z3Complex) Z3Complex) z3.Value {
	switch left.(type) {
	case z3.Int:
		switch right.(type) {
		case z3.Int:
			return intOp(left.(z3.Int), right.(z3.Int))
		case z3.Float:
			return floatOp(left.(z3.Int).ToReal().ToFloat(ctx.Sorts.FloatSort), right.(z3.Float))
		}
	case z3.Float:
		switch right.(type) {
		case z3.Float:
			return floatOp(left.(z3.Float), right.(z3.Float))
		case z3.Int:
			return floatOp(left.(z3.Float), right.(z3.Int).ToReal().ToFloat(ctx.Sorts.FloatSort))
		}
	case Z3Complex:
		return complexOp(left.(Z3Complex), right.(Z3Complex))
	}

	panic("unsupported arguments")
}

func (ctx *AnalysisContext) TypeToSort(t types.Type) z3.Sort {
	switch t.(type) {
	case *types.Basic:
		switch t.(*types.Basic).Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Byte:
			return ctx.Sorts.IntSort
		case types.UntypedFloat, types.Float32, types.Float64:
			return ctx.Sorts.FloatSort
		case types.UntypedComplex, types.Complex64, types.Complex128:
			return ctx.Sorts.ComplexSort
		}
	}

	return ctx.Sorts.UnknownSort
}

func FloatToString(f z3.Float) string {
	float, _ := f.AsBigFloat()

	return fmt.Sprintf("%d", float)
}

func (ctx *AnalysisContext) GoToZ3Value(v any) z3.Value {
	switch casted := v.(type) {
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return ctx.Z3ctx.FromInt(int64(casted.(int)), ctx.Sorts.IntSort)
	case float64, float32:
		return ctx.Z3ctx.FromFloat64(float64(casted.(float64)), ctx.Sorts.FloatSort)
	case bool:
		return ctx.Z3ctx.FromBool(casted)
	case complex128:
		return ctx.NewComplex(casted)
	case complex64:
		return ctx.NewComplex(complex128(casted))
	default:
		panic("unsupported argument")
	}
}
