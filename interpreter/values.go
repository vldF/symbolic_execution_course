package interpreter

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
)

func (left *ConcreteBoolValue) String() string {
	return fmt.Sprintf("ConcreteBoolValue(%v)", left.Value)
}

type Value interface {
	AsZ3Value() Z3Value

	Eq(Value) BoolValue
	NotEq(Value) BoolValue

	IsFloat() bool
	IsInteger() bool
	IsBool() bool

	And(Value) Value
	Or(Value) Value
	Xor(Value) Value
}

type ArithmeticValue interface {
	Add(ArithmeticValue) ArithmeticValue
	Sub(ArithmeticValue) ArithmeticValue
	Mul(ArithmeticValue) ArithmeticValue
	Div(ArithmeticValue) ArithmeticValue

	AsZ3Value() Z3Value
	Gt(ArithmeticValue) BoolValue
	Ge(ArithmeticValue) BoolValue
	Lt(ArithmeticValue) BoolValue
	Le(ArithmeticValue) BoolValue

	Eq(Value) BoolValue
	NotEq(Value) BoolValue

	AsInt() Value
	AsFloat() Value
	AsBool() Value

	asZ3IntValue() z3.BV
	asZ3FloatValue() z3.Float

	IsBool() bool
	IsFloat() bool
	IsInteger() bool

	Rem(Value) Value

	Shl(ArithmeticValue) ArithmeticValue
	Shr(ArithmeticValue) ArithmeticValue

	BitwiseAnd(ArithmeticValue) ArithmeticValue
	BitwiseOr(ArithmeticValue) ArithmeticValue
	BitwiseXor(ArithmeticValue) ArithmeticValue

	And(Value) Value
	Or(Value) Value
	Xor(Value) Value
}

type BoolValue interface {
	AsZ3Value() Z3Value

	Eq(Value) BoolValue
	NotEq(Value) BoolValue

	AsBool() Value

	IsFloat() bool
	IsInteger() bool
	IsBool() bool

	BoolAnd(BoolValue) BoolValue
	BoolOr(BoolValue) BoolValue
	BoolXor(BoolValue) BoolValue

	And(Value) Value
	Or(Value) Value
	Xor(Value) Value

	Not() BoolValue
}

type ConcreteValue[V int64 | bool | float64] struct {
	Context *Context
	Value   V
}

type ConcreteIntValue ConcreteValue[int64]
type ConcreteFloatValue ConcreteValue[float64]

func (left *ConcreteBoolValue) BoolAnd(value BoolValue) BoolValue {
	switch casterRight := value.(type) {
	case *ConcreteBoolValue:
		return &ConcreteBoolValue{
			left.Context,
			left.Value && casterRight.Value,
		}
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.Bool).And(casterRight.Value.(z3.Bool)),
		}
	}

	panic("unsupported")
}

func (left *ConcreteBoolValue) BoolOr(value BoolValue) BoolValue {
	switch casterRight := value.(type) {
	case *ConcreteBoolValue:
		return &ConcreteBoolValue{
			left.Context,
			left.Value && casterRight.Value,
		}
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.Bool).Or(casterRight.Value.(z3.Bool)),
		}
	}

	panic("unsupported")
}

func (left *ConcreteValue[T]) String() string {
	return fmt.Sprintf("ConcretValue(%v)", left.Value)
}

func (left *ConcreteIntValue) Add(v ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, v,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).Add(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Add(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value + right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value + right.Value}
		},

		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func (left *ConcreteIntValue) Sub(v ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, v,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).Sub(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Sub(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value - right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value - right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func (left *ConcreteIntValue) Mul(v ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, v,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).Mul(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Mul(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value * right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value * right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func (left *ConcreteIntValue) Div(v ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, v,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).SDiv(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Div(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value / right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value / right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func (left *ConcreteFloatValue) Add(v ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, v,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).Add(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Add(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value + right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value + right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func (left *ConcreteFloatValue) Sub(v ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, v,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).Sub(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Sub(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value - right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value - right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func (left *ConcreteFloatValue) Mul(v ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, v,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).Mul(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Mul(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value * right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value * right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func (left *ConcreteFloatValue) Div(v ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, v,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).SDiv(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Div(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value / right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value / right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func numericValueBinop[T any](
	left Value,
	right Value,

	intZ3ValueOp func(*Z3Value, *Z3Value) T,
	floatZ3ValueOp func(*Z3Value, *Z3Value) T,
	boolZ3ValueOp func(*Z3Value, *Z3Value) T,

	intConcreteValueOp func(*ConcreteIntValue, *ConcreteIntValue) T,
	floatConcreteValueOp func(*ConcreteFloatValue, *ConcreteFloatValue) T,
	boolConcreteValueOp func(*ConcreteBoolValue, *ConcreteBoolValue) T) T {
	switch castedLeft := left.(type) {
	case *ConcreteIntValue:
		switch castedRight := right.(type) {
		case *ConcreteIntValue:
			return intConcreteValueOp(castedLeft, castedRight)
		case *ConcreteFloatValue:
			return floatConcreteValueOp(castedLeft.toFloat(), castedRight)
		case *Z3Value:
			switch {
			case castedRight.IsFloat():
				leftNew := castedLeft.toFloat().AsZ3Value()
				return floatZ3ValueOp(&leftNew, castedRight)
			case castedRight.IsInteger():
				leftNew := castedLeft.AsZ3Value()
				return intZ3ValueOp(&leftNew, castedRight)
			default:
				panic("unreachable")
			}
		default:
			panic("unreachable")
		}
	case *ConcreteFloatValue:
		switch castedRight := right.(type) {
		case *ConcreteIntValue:
			return floatConcreteValueOp(castedLeft, castedRight.toFloat())
		case *ConcreteFloatValue:
			return floatConcreteValueOp(castedLeft, castedRight)
		case *Z3Value:
			newLeft := castedLeft.AsZ3Value()
			newRight := castedRight.AsZ3Value()
			return floatZ3ValueOp(&newLeft, &newRight)
		default:
			panic("unreachable")
		}
	case *ConcreteBoolValue:
		switch castedRight := right.(type) {
		case *ConcreteIntValue:
			panic("unsupported")
		case *ConcreteFloatValue:
			panic("unsupported")
		case *Z3Value:
			newLeft := castedLeft.AsZ3Value()
			newRight := castedRight.AsZ3Value()
			return boolZ3ValueOp(&newLeft, &newRight)
		default:
			panic("unreachable")
		}
	case *Z3Value:
		switch castedRight := right.(type) {
		case *ConcreteIntValue:
			newRight := castedRight.AsZ3Value()
			return intZ3ValueOp(castedLeft, &newRight)
		case *ConcreteFloatValue:
			newRight := castedRight.AsZ3Value()
			return floatZ3ValueOp(castedLeft, &newRight)
		case *Z3Value:
			switch {
			case castedLeft.IsInteger() && castedRight.IsInteger():
				return intZ3ValueOp(castedLeft, castedRight)
			case castedLeft.IsInteger() && castedRight.IsFloat():
				leftNew := castedLeft.AsFloat().AsZ3Value()
				return floatZ3ValueOp(&leftNew, castedRight)
			case castedLeft.IsFloat() && castedRight.IsInteger():
				rightNew := castedLeft.AsFloat().AsZ3Value()
				return floatZ3ValueOp(castedLeft, &rightNew)
			case castedLeft.IsFloat() && castedRight.IsFloat():
				leftNew := castedLeft.AsFloat().AsZ3Value()
				rightNew := castedRight.AsFloat().AsZ3Value()
				return floatZ3ValueOp(&leftNew, &rightNew)
			case castedLeft.IsBool() && castedRight.IsBool():
				leftNew := castedLeft.AsBool().AsZ3Value()
				rightNew := castedRight.AsBool().AsZ3Value()
				return boolZ3ValueOp(&leftNew, &rightNew)
			default:
				panic("unreachable")
			}
		}
	default:
		panic("unreachable")
	}

	panic("unreachable")
}

func (left *ConcreteIntValue) Shl(right ArithmeticValue) ArithmeticValue {
	switch castedRight := right.(type) {
	case *ConcreteIntValue:
		return &ConcreteIntValue{
			left.Context,
			left.Value << castedRight.Value}
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.BV).Lsh(castedRight.Value.(z3.BV))}
	}
	panic("unreachable")
}
func (left *ConcreteIntValue) Shr(right ArithmeticValue) ArithmeticValue {
	switch castedRight := right.(type) {
	case *ConcreteIntValue:
		return &ConcreteIntValue{
			left.Context,
			left.Value >> castedRight.Value}
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.BV).URsh(castedRight.Value.(z3.BV))}
	}
	panic("unreachable")
}

func (left *ConcreteFloatValue) Shl(right ArithmeticValue) ArithmeticValue {
	panic("unsupported")
}
func (left *ConcreteFloatValue) Shr(right ArithmeticValue) ArithmeticValue {
	panic("unsupported")
}

func (left *Z3Value) Shl(right ArithmeticValue) ArithmeticValue {
	switch castedRight := right.(type) {
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.Value.(z3.BV).Lsh(castedRight.Value.(z3.BV)),
		}
	case ArithmeticValue:
		return &Z3Value{
			left.Context,
			left.Value.(z3.BV).Lsh(castedRight.asZ3IntValue()),
		}
	}

	panic("unreachable")
}

func (left *Z3Value) Shr(right ArithmeticValue) ArithmeticValue {
	switch castedRight := right.(type) {
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.Value.(z3.BV).SRsh(castedRight.Value.(z3.BV)),
		}
	case ArithmeticValue:
		return &Z3Value{
			left.Context,
			left.Value.(z3.BV).SRsh(castedRight.asZ3IntValue()),
		}
	}

	panic("unreachable")
}

func (left *ConcreteIntValue) toFloat() *ConcreteFloatValue {
	return &ConcreteFloatValue{
		left.Context,
		float64(left.Value),
	}
}

func (left *Z3Value) toFloat() *Z3Value {
	return &Z3Value{
		left.Context,
		left.asZ3FloatValue(),
	}
}

func (left *ConcreteFloatValue) toInt() *ConcreteIntValue {
	return &ConcreteIntValue{
		left.Context,
		int64(left.Value),
	}
}

func (left *ConcreteIntValue) AsZ3Value() Z3Value {
	return Z3Value{
		left.Context,
		left.Context.Z3Context.FromInt(
			left.Value,
			left.Context.TypesContext.IntSort).(z3.BV),
	}
}

func (left *ConcreteFloatValue) AsZ3Value() Z3Value {
	return Z3Value{
		left.Context,
		left.Context.Z3Context.FromFloat64(
			left.Value,
			left.Context.TypesContext.FloatSort),
	}
}

type ConcreteBoolValue ConcreteValue[bool]

type Z3Value struct {
	Context *Context
	Value   z3.Value
}

func (left *ConcreteIntValue) asZ3IntValue() z3.BV {
	return left.Context.Z3Context.FromInt(left.Value, left.Context.TypesContext.IntSort).(z3.BV)
}

func (left *ConcreteFloatValue) asZ3IntValue() z3.BV {
	return left.Context.Z3Context.FromFloat64(left.Value, left.Context.TypesContext.FloatSort).ToReal().ToInt().
		ToBV(left.Context.TypesContext.IntBits)
}

func (left *ConcreteIntValue) asZ3FloatValue() z3.Float {
	return left.Context.Z3Context.FromInt(left.Value, left.Context.TypesContext.IntSort).(z3.Int).ToReal().
		ToFloat(left.Context.TypesContext.FloatSort)
}

func (left *ConcreteFloatValue) asZ3FloatValue() z3.Float {
	return left.Context.Z3Context.FromFloat64(left.Value, left.Context.TypesContext.FloatSort)
}

func (left *ConcreteIntValue) IsInteger() bool {
	return true
}

func (left *ConcreteFloatValue) IsInteger() bool {
	return false
}

func (left *ConcreteBoolValue) IsInteger() bool {
	return false
}

func (left *ConcreteIntValue) IsFloat() bool {
	return false
}

func (left *ConcreteFloatValue) IsFloat() bool {
	return true
}

func (left *ConcreteBoolValue) IsFloat() bool {
	return true
}

func (left *ConcreteFloatValue) IsBool() bool {
	return false
}

func (left *ConcreteIntValue) IsBool() bool {
	return false
}

func (left *ConcreteBoolValue) IsBool() bool {
	return true
}

func (left *Z3Value) IsFloat() bool {
	if _, ok := left.Value.(z3.Float); ok {
		return true
	}

	return false
}

func (left *Z3Value) Add(right ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, right,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).Add(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Add(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value + right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value + right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func (left *Z3Value) Sub(right ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, right,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).Sub(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Sub(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value - right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value - right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func (left *Z3Value) Mul(right ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, right,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).Mul(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Mul(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value * right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value * right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func (left *Z3Value) Div(right ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, right,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).SDiv(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Div(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return &ConcreteIntValue{left.Context, left.Value / right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return &ConcreteFloatValue{left.Context, left.Value / right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) ArithmeticValue {
			panic("unsupported")
		},
	)
}

func (left *Z3Value) AsZ3Value() Z3Value {
	return *left
}

func (left *Z3Value) asZ3IntValue() z3.BV {
	switch left.IsFloat() {
	case true:
		return left.Value.(z3.Float).ToReal().ToInt().ToBV(left.Context.TypesContext.IntBits)
	case false:
		return left.Value.(z3.BV)
	}

	panic("unreachable")
}

func (left *Z3Value) asZ3FloatValue() z3.Float {
	switch left.IsFloat() {
	case true:
		return left.Value.(z3.Float)
	case false:
		return left.Value.(z3.BV).SToInt().ToReal().ToFloat(left.Context.TypesContext.FloatSort)
	}

	panic("unreachable")
}

func (left *Z3Value) IsInteger() bool {
	if _, ok := left.Value.(z3.BV); ok {
		return true
	}

	return false
}

func (left *Z3Value) IsBool() bool {
	if _, ok := left.Value.(z3.Bool); ok {
		return true
	}

	return false
}

func (ctx *Context) eq(left z3.Value, right z3.Value) z3.Bool {
	switch left.(type) {
	case z3.BV:
		return left.(z3.BV).Eq(right.(z3.BV))
	case z3.Bool:
		return left.(z3.Bool).Eq(right.(z3.Bool))
	case z3.Float:
		return left.(z3.Float).Eq(right.(z3.Float))
	case z3.Array:
		return left.(z3.Array).Eq(right.(z3.Array))
	}
	panic("can't build eq")
}

func (left *Z3Value) Gt(right ArithmeticValue) BoolValue {
	return gt(left, right)
}

func (left *ConcreteIntValue) Gt(right ArithmeticValue) BoolValue {
	return gt(left, right)
}

func (left *ConcreteFloatValue) Gt(right ArithmeticValue) BoolValue {
	return gt(left, right)
}

func gt(left ArithmeticValue, right ArithmeticValue) BoolValue {
	return numericValueBinop(
		left, right,

		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).SGT(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).GT(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) BoolValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value > right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value > right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) BoolValue {
			panic("unsupported")
		},
	)
}

func (left *Z3Value) Ge(right ArithmeticValue) BoolValue {
	return ge(left, right)
}

func (left *ConcreteIntValue) Ge(right ArithmeticValue) BoolValue {
	return ge(left, right)
}

func (left *ConcreteFloatValue) Ge(right ArithmeticValue) BoolValue {
	return ge(left, right)
}

func ge(left ArithmeticValue, right ArithmeticValue) BoolValue {
	return numericValueBinop(
		left, right,

		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).SGE(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).GE(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) BoolValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value >= right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value >= right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) BoolValue {
			panic("unsupported")
		},
	)
}

func (left *Z3Value) Lt(right ArithmeticValue) BoolValue {
	return lt(left, right)
}

func (left *ConcreteIntValue) Lt(right ArithmeticValue) BoolValue {
	return lt(left, right)
}

func (left *ConcreteFloatValue) Lt(right ArithmeticValue) BoolValue {
	return lt(left, right)
}

func lt(left ArithmeticValue, right ArithmeticValue) BoolValue {
	return numericValueBinop(
		left, right,

		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).SLT(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).LT(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) BoolValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value < right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value < right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) BoolValue {
			panic("unsupported")
		},
	)
}

func (left *Z3Value) Le(right ArithmeticValue) BoolValue {
	return le(left, right)
}

func (left *ConcreteIntValue) Le(right ArithmeticValue) BoolValue {
	return le(left, right)
}

func (left *ConcreteFloatValue) Le(right ArithmeticValue) BoolValue {
	return le(left, right)
}

func le(left ArithmeticValue, right ArithmeticValue) BoolValue {
	return numericValueBinop(
		left, right,

		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).SLE(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).LE(right.Value.(z3.Float))}
		},
		func(*Z3Value, *Z3Value) BoolValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value <= right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value <= right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) BoolValue {
			panic("unsupported")
		},
	)
}

func (left *Z3Value) Eq(right Value) BoolValue {
	return eq(left, right)
}

func (left *ConcreteIntValue) Eq(right Value) BoolValue {
	return eq(left, right)
}

func (left *ConcreteFloatValue) Eq(right Value) BoolValue {
	return eq(left, right)
}

func (left *ConcreteBoolValue) Eq(right Value) BoolValue {
	return eq(left, right)
}

func eq(left Value, right Value) BoolValue {
	return numericValueBinop(
		left, right,

		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).Eq(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).Eq(right.Value.(z3.Float))}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Bool).Eq(right.Value.(z3.Bool))}
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value == right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value == right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value == right.Value}
		},
	)
}

func (left *Z3Value) NotEq(right Value) BoolValue {
	return notEq(left, right)
}

func (left *ConcreteIntValue) NotEq(right Value) BoolValue {
	return notEq(left, right)
}

func (left *ConcreteFloatValue) NotEq(right Value) BoolValue {
	return notEq(left, right)
}

func (left *ConcreteBoolValue) NotEq(right Value) BoolValue {
	return notEq(left, right)
}

func notEq(left Value, right Value) BoolValue {
	return numericValueBinop(
		left, right,

		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.BV).NE(right.Value.(z3.BV))}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{left.Context, left.Value.(z3.Float).NE(right.Value.(z3.Float))}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Bool).NE(right.Value.(z3.Bool))}
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value != right.Value}
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value != right.Value}
		},
		func(left *ConcreteBoolValue, right *ConcreteBoolValue) BoolValue {
			return &ConcreteBoolValue{left.Context, left.Value != right.Value}
		},
	)
}

func (left *ConcreteIntValue) AsInt() Value {
	return left
}

func (left *ConcreteFloatValue) AsInt() Value {
	return &ConcreteIntValue{
		left.Context,
		int64(left.Value),
	}
}

func (left *Z3Value) AsInt() Value {
	switch {
	case left.IsFloat():
		return &Z3Value{
			left.Context,
			left.Value.(z3.Float).ToReal().ToInt(),
		}
	case left.IsInteger():
		return left
	}

	panic("unreachable")
}

func (left *ConcreteIntValue) AsFloat() Value {
	return &ConcreteFloatValue{
		left.Context,
		float64(left.Value),
	}
}

func (left *ConcreteFloatValue) AsFloat() Value {
	return left
}

func (left *Z3Value) AsFloat() Value {
	switch {
	case left.IsFloat():
		return left
	case left.IsInteger():
		return &Z3Value{
			left.Context,
			left.Value.(z3.BV).SToFloat(left.Context.TypesContext.FloatSort),
		}
	}

	panic("unreachable")
}

func (left *ConcreteIntValue) AsBool() Value {
	panic("Unsupported")
}

func (left *ConcreteFloatValue) AsBool() Value {
	panic("Unsupported")
}

func (left *Z3Value) AsBool() Value {
	switch {
	case left.IsFloat():
		panic("Unsupported")
	case left.IsInteger():
		panic("Unsupported")
	case left.IsBool():
		return left
	}

	panic("unreachable")
}

func (left *ConcreteIntValue) Rem(right Value) Value {
	switch castedRight := right.(type) {
	case *ConcreteIntValue:
		return &ConcreteIntValue{
			left.Context,
			left.Value % castedRight.Value,
		}
	case *ConcreteFloatValue:
		panic("unsupported!")
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.asZ3IntValue().SRem(castedRight.Value.(z3.BV)),
		}
	}

	panic("unreachable")
}

func (left *ConcreteFloatValue) Rem(right Value) Value {
	panic("unsupported!")
}

func (left *Z3Value) Rem(right Value) Value {
	switch {
	case left.IsFloat():
		panic("unsupported!")
	case left.IsInteger():
		leftBV := left.asZ3IntValue()

		switch castedRight := right.(type) {
		case *ConcreteIntValue:
			return &Z3Value{
				left.Context,
				leftBV.SRem(left.Context.Z3Context.FromInt(castedRight.Value, left.Context.TypesContext.IntSort).(z3.BV)),
			}
		case *ConcreteFloatValue:
			panic("unsupported!")
		}
	}

	panic("unsupported!")
}

func (left *ConcreteBoolValue) Not() BoolValue {
	return &ConcreteBoolValue{
		Context: left.Context,
		Value:   !left.Value,
	}
}

func (left *Z3Value) Not() BoolValue {
	return &Z3Value{
		Context: left.Context,
		Value:   left.Value.(z3.Bool).Not(),
	}
}

func (left *ConcreteBoolValue) Add(right Value) Value {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) Sub(right Value) Value {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) Mul(right Value) Value {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) Div(right Value) Value {
	panic("Unsupported")
}

func (left *ConcreteBoolValue) Gt(Value) Value {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) Ge(Value) Value {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) Lt(Value) Value {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) Le(Value) Value {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) AsInt() Value {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) AsFloat() Value {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) asZ3IntValue() z3.BV {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) asZ3FloatValue() z3.Float {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) Rem(Value) Value {
	panic("Unsupported")
}

func (left *ConcreteBoolValue) AsZ3Value() Z3Value {
	return Z3Value{
		left.Context,
		left.Context.Z3Context.FromBool(left.Value),
	}
}

func (left *ConcreteBoolValue) AsBool() Value {
	return left
}

func (left *ConcreteBoolValue) BoolXor(value BoolValue) BoolValue {
	switch castedRight := value.(type) {
	case *ConcreteBoolValue:
		return &ConcreteBoolValue{
			left.Context,
			left.Value != castedRight.Value, // xor
		}
	case *Z3Value:
		switch {
		case value.IsBool():
			return &Z3Value{
				left.Context,
				left.Context.Z3Context.FromBool(left.Value).Xor(castedRight.Value.(z3.Bool)),
			}
		}
	}
	panic("unsupported!")
}

func (left *Z3Value) BoolAnd(value BoolValue) BoolValue {
	return value.BoolAnd(left)
}

func (left *Z3Value) BoolOr(value BoolValue) BoolValue {
	return value.BoolOr(left)
}

func (left *Z3Value) BoolXor(value BoolValue) BoolValue {
	return value.BoolXor(left)
}

func (left *ConcreteIntValue) And(value Value) Value {
	return left.BitwiseAnd(value.(ArithmeticValue))
}
func (left *ConcreteIntValue) Or(value Value) Value {
	return left.BitwiseOr(value.(ArithmeticValue))
}
func (left *ConcreteIntValue) Xor(value Value) Value {
	return left.BitwiseXor(value.(ArithmeticValue))
}

func (left *ConcreteFloatValue) And(value Value) Value {
	panic("unsupported")
}
func (left *ConcreteFloatValue) Or(value Value) Value {
	panic("unsupported")
}
func (left *ConcreteFloatValue) Xor(value Value) Value {
	panic("unsupported")
}

func (left *ConcreteBoolValue) And(value Value) Value {
	return left.BoolAnd(value.(BoolValue))
}
func (left *ConcreteBoolValue) Or(value Value) Value {
	return left.BoolOr(value.(BoolValue))
}
func (left *ConcreteBoolValue) Xor(value Value) Value {
	return left.BoolXor(value.(BoolValue))
}

func (left *Z3Value) And(value Value) Value {
	switch {
	case left.IsBool():
		return left.BoolAnd(value.(BoolValue))
	default:
		return left.BitwiseAnd(value.(ArithmeticValue))
	}
}
func (left *Z3Value) Or(value Value) Value {
	switch {
	case left.IsBool():
		return left.BoolOr(value.(BoolValue))
	default:
		return left.BitwiseOr(value.(ArithmeticValue))
	}
}
func (left *Z3Value) Xor(value Value) Value {
	switch {
	case left.IsBool():
		return left.BoolXor(value.(BoolValue))
	default:
		return left.BitwiseXor(value.(ArithmeticValue))
	}
}

func (left *ConcreteIntValue) BitwiseAnd(value ArithmeticValue) ArithmeticValue {
	switch castedRight := value.(type) {
	case *ConcreteIntValue:
		return &ConcreteIntValue{
			left.Context,
			left.Value & castedRight.Value,
		}
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.BV).And(castedRight.Value.(z3.BV)),
		}
	}

	panic("unsupported!")
}
func (left *ConcreteIntValue) BitwiseOr(value ArithmeticValue) ArithmeticValue {
	switch castedRight := value.(type) {
	case *ConcreteIntValue:
		return &ConcreteIntValue{
			left.Context,
			left.Value | castedRight.Value,
		}
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.BV).Or(castedRight.Value.(z3.BV)),
		}
	}

	panic("unsupported!")
}
func (left *ConcreteIntValue) BitwiseXor(value ArithmeticValue) ArithmeticValue {
	switch castedRight := value.(type) {
	case *ConcreteIntValue:
		return &ConcreteIntValue{
			left.Context,
			left.Value ^ castedRight.Value,
		}
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.BV).Xor(castedRight.Value.(z3.BV)),
		}
	}

	panic("unsupported!")
}

func (left *ConcreteFloatValue) BitwiseAnd(value ArithmeticValue) ArithmeticValue {
	panic("unsupported!")
}
func (left *ConcreteFloatValue) BitwiseOr(value ArithmeticValue) ArithmeticValue {
	panic("unsupported!")
}
func (left *ConcreteFloatValue) BitwiseXor(value ArithmeticValue) ArithmeticValue {
	panic("unsupported!")
}

func (left *Z3Value) BitwiseAnd(value ArithmeticValue) ArithmeticValue {
	switch castedRight := value.(type) {
	case *ConcreteIntValue:
		return castedRight.BitwiseAnd(left)
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.Value.(z3.BV).And(castedRight.Value.(z3.BV)),
		}
	}

	panic("unsupported!")
}
func (left *Z3Value) BitwiseOr(value ArithmeticValue) ArithmeticValue {
	switch castedRight := value.(type) {
	case *ConcreteIntValue:
		return castedRight.BitwiseOr(left)
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.Value.(z3.BV).Or(castedRight.Value.(z3.BV)),
		}
	}

	panic("unsupported!")
}
func (left *Z3Value) BitwiseXor(value ArithmeticValue) ArithmeticValue {
	switch castedRight := value.(type) {
	case *ConcreteIntValue:
		return castedRight.BitwiseXor(left)
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.Value.(z3.BV).Xor(castedRight.Value.(z3.BV)),
		}
	}

	panic("unsupported!")
}
