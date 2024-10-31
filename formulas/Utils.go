package formulas

import (
	"github.com/aclements/go-z3/z3"
	"go/types"
)

func (ctx AnalysisContext) Eq(left z3.Value, right z3.Value) z3.Bool {
	switch left.(type) {
	case z3.Int:
		return left.(z3.Int).Eq(right.(z3.Int))
	case z3.Bool:
		return left.(z3.Bool).Eq(right.(z3.Bool))
	case z3.Float:
		return left.(z3.Float).Eq(right.(z3.Float))
	}

	return ctx.Z3ctx.FromBool(false)
}

func (ctx AnalysisContext) Ne(left z3.Value, right z3.Value) z3.Bool {
	return ctx.Eq(left, right).Not()
}

func (ctx AnalysisContext) Add(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.Add(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.Add(second) })
}

func (ctx AnalysisContext) Sub(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.Sub(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.Sub(second) })
}

func (ctx AnalysisContext) Mul(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.Mul(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.Mul(second) })
}

func (ctx AnalysisContext) Div(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.Div(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.Div(second) })
}

func (ctx AnalysisContext) Lt(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.LT(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.LT(second) })
}

func (ctx AnalysisContext) Le(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.LE(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.LE(second) })
}

func (ctx AnalysisContext) Gt(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.GT(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.GT(second) })
}

func (ctx AnalysisContext) Ge(left z3.Value, right z3.Value) z3.Value {
	return ctx.arithOp(
		left,
		right,
		func(first z3.Float, second z3.Float) z3.Value { return first.GE(second) },
		func(first z3.Int, second z3.Int) z3.Value { return first.GE(second) })
}

func (ctx AnalysisContext) arithOp(
	left z3.Value,
	right z3.Value,
	floatOp func(z3.Float, z3.Float) z3.Value,
	intOp func(z3.Int, z3.Int) z3.Value) z3.Value {
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
	}

	panic("unsupported arguments")
}

func (ctx AnalysisContext) TypeToSort(t types.Type) z3.Sort {
	switch t.(type) {
	case *types.Basic:
		switch t.(*types.Basic).Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Byte:
			return ctx.Sorts.IntSort
		case types.UntypedFloat, types.Float32, types.Float64:
			return ctx.Sorts.FloatSort
		}
	}

	return ctx.Sorts.UnknownSort
}
