package interpreter

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
)

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

	AsInt(bits int) Value
	AsFloat(bits int) Value
	AsBool() Value

	asZ3IntValue() z3.BV
	asZ3FloatValue(bits int) z3.Float

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

type ConcreteIntValue struct {
	Context *Context
	Value   int64
	bits    int
}

type ConcreteFloatValue struct {
	Context *Context
	Value   float64
	bits    int
}

type ConcreteBoolValue struct {
	Context *Context
	Value   bool
}

func (left *ConcreteBoolValue) BoolAnd(value BoolValue) BoolValue {
	switch casterRight := value.(type) {
	case *ConcreteBoolValue:
		return &ConcreteBoolValue{
			Context: left.Context,
			Value:   left.Value && casterRight.Value,
		}
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.Bool).And(casterRight.Value.(z3.Bool)),
			1,
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
			1,
		}
	}

	panic("unsupported")
}

func (left *ConcreteIntValue) String() string {
	return fmt.Sprintf("ConcretIntValue(%v)", left.Value)
}

func (left *ConcreteFloatValue) String() string {
	return fmt.Sprintf("ConcretFloatValue(%v)", left.Value)
}

func (left *ConcreteBoolValue) String() string {
	return fmt.Sprintf("ConcretBoolValue(%v)", left.Value)
}

func (left *ConcreteIntValue) Add(v ArithmeticValue) ArithmeticValue {
	return numericValueBinop(
		left, v,

		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).Add(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Add(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value+right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value+right.Value, max(left.bits, right.bits))
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).Sub(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Sub(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value-right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value-right.Value, max(left.bits, right.bits))
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).Mul(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Mul(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value*right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value-right.Value, max(left.bits, right.bits))
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).SDiv(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Div(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value/right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value/right.Value, max(left.bits, right.bits))
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).Add(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Add(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value+right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value+right.Value, max(left.bits, right.bits))
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).Sub(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Sub(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value-right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value-right.Value, max(left.bits, right.bits))
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).Mul(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Mul(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value*right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value*right.Value, max(left.bits, right.bits))
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).SDiv(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Div(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value/right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value/right.Value, max(left.bits, right.bits))
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
			return floatConcreteValueOp(castedLeft.toFloat(castedRight.bits), castedRight)
		case *Z3Value:
			switch {
			case castedRight.IsFloat():
				leftNew := castedLeft.toFloat(castedLeft.bits).AsZ3Value()
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
			return floatConcreteValueOp(castedLeft, castedRight.toFloat(castedLeft.bits))
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
				leftNew := castedLeft.AsFloat(64).AsZ3Value()
				return floatZ3ValueOp(&leftNew, castedRight)
			case castedLeft.IsFloat() && castedRight.IsInteger():
				rightNew := castedLeft.AsFloat(64).AsZ3Value()
				return floatZ3ValueOp(castedLeft, &rightNew)
			case castedLeft.IsFloat() && castedRight.IsFloat():
				leftNew := castedLeft.AsFloat(64).AsZ3Value()
				rightNew := castedRight.AsFloat(64).AsZ3Value()
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
		return left.Context.CreateInt(left.Value<<castedRight.Value, left.bits)
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.BV).Lsh(castedRight.Value.(z3.BV)),
			left.bits,
		}
	}
	panic("unreachable")
}
func (left *ConcreteIntValue) Shr(right ArithmeticValue) ArithmeticValue {
	switch castedRight := right.(type) {
	case *ConcreteIntValue:
		return left.Context.CreateInt(left.Value>>castedRight.Value, left.bits)
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.BV).URsh(castedRight.Value.(z3.BV)),
			left.bits,
		}
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
			left.Bits,
		}
	case ArithmeticValue:
		return &Z3Value{
			left.Context,
			left.Value.(z3.BV).Lsh(castedRight.asZ3IntValue()),
			left.Bits,
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
			left.Bits,
		}
	case ArithmeticValue:
		return &Z3Value{
			left.Context,
			left.Value.(z3.BV).SRsh(castedRight.asZ3IntValue()),
			left.Bits,
		}
	}

	panic("unreachable")
}

func (left *ConcreteIntValue) toFloat(bits int) *ConcreteFloatValue {
	return left.Context.CreateFloat(float64(left.Value), bits)
}

func (left *ConcreteIntValue) AsZ3Value() Z3Value {
	return Z3Value{
		left.Context,
		left.Context.Z3Context.FromInt(
			left.Value,
			*left.Context.TypesContext.GetIntSort(left.bits)),
		left.bits,
	}
}

func (left *ConcreteFloatValue) AsZ3Value() Z3Value {
	return Z3Value{
		left.Context,
		left.Context.Z3Context.FromFloat64(
			left.Value,
			*left.Context.TypesContext.GetFloatSort(left.bits)),
		left.bits,
	}
}

type Z3Value struct {
	Context *Context
	Value   z3.Value
	Bits    int
}

func (left *ConcreteIntValue) asZ3IntValue() z3.BV {
	return left.Context.Z3Context.FromInt(
		left.Value,
		*left.Context.TypesContext.GetIntSort(left.bits)).(z3.BV)
}

func (left *ConcreteFloatValue) asZ3IntValue() z3.BV {
	return left.Context.Z3Context.FromFloat64(
		left.Value,
		*left.Context.TypesContext.GetFloatSort(left.bits)).
		ToReal().
		ToInt().
		ToBV(left.Context.TypesContext.Int.Bits)
}

func (left *ConcreteIntValue) asZ3FloatValue(bits int) z3.Float {
	return left.Context.Z3Context.FromInt(
		left.Value,
		*left.Context.TypesContext.GetIntSort(bits)).(z3.Int).
		ToReal().
		ToFloat(left.Context.TypesContext.Float64.Sort)
}

func (left *ConcreteFloatValue) asZ3FloatValue(bits int) z3.Float {
	return left.Context.Z3Context.FromFloat64(left.Value, *left.Context.TypesContext.GetFloatSort(bits))
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).Add(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Add(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value+right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value+right.Value, max(left.bits, right.bits))
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).Sub(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Sub(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value-right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value-right.Value, max(left.bits, right.bits))
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
			return &Z3Value{
				Context: left.Context,
				Value:   left.Value.(z3.BV).Mul(right.Value.(z3.BV)),
				Bits:    left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Mul(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value*right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value*right.Value, max(left.bits, right.bits))
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).SDiv(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) ArithmeticValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Div(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(*Z3Value, *Z3Value) ArithmeticValue {
			panic("unsupported")
		},

		func(left *ConcreteIntValue, right *ConcreteIntValue) ArithmeticValue {
			return left.Context.CreateInt(left.Value/right.Value, max(left.bits, right.bits))
		},
		func(left *ConcreteFloatValue, right *ConcreteFloatValue) ArithmeticValue {
			return left.Context.CreateFloat(left.Value/right.Value, max(left.bits, right.bits))
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
		return left.Value.(z3.Float).ToReal().ToInt().ToBV(left.Context.TypesContext.Int.Bits)
	case false:
		return left.Value.(z3.BV)
	}

	panic("unreachable")
}

func (left *Z3Value) asZ3FloatValue(bits int) z3.Float {
	switch left.IsFloat() {
	case true:
		return left.Value.(z3.Float)
	case false:
		return left.Value.(z3.BV).SToInt().ToReal().ToFloat(*left.Context.TypesContext.GetFloatSort(bits))
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).SGT(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).GT(right.Value.(z3.Float)),
				left.Bits,
			}
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).SGE(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).GE(right.Value.(z3.Float)),
				left.Bits,
			}
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).SLT(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).LT(right.Value.(z3.Float)),
				left.Bits,
			}
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).SLE(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).LE(right.Value.(z3.Float)),
				left.Bits,
			}
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).Eq(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).Eq(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Bool).Eq(right.Value.(z3.Bool)),
				left.Bits,
			}
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
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).NE(right.Value.(z3.BV)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Float).NE(right.Value.(z3.Float)),
				left.Bits,
			}
		},
		func(left *Z3Value, right *Z3Value) BoolValue {
			return &Z3Value{
				left.Context,
				left.Value.(z3.Bool).NE(right.Value.(z3.Bool)),
				left.Bits,
			}
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

func (left *ConcreteIntValue) AsInt(bits int) Value {
	return left.Context.CreateInt(left.Value, bits)
}

func (left *ConcreteFloatValue) AsInt(bits int) Value {
	return left.Context.CreateInt(int64(left.Value), bits)
}

func (left *Z3Value) AsInt(bits int) Value {
	switch {
	case left.IsFloat():
		return &Z3Value{
			left.Context,
			left.Value.(z3.Float).ToReal().ToInt().ToBV(bits),
			left.Bits,
		}
	case left.IsInteger():
		switch left.Bits {
		case bits:
			return left
		default:
			return &Z3Value{
				left.Context,
				left.Value.(z3.BV).SToInt().ToBV(bits),
				left.Bits,
			}
		}
	}

	panic("unreachable")
}

func (left *ConcreteIntValue) AsFloat(bits int) Value {
	return left.Context.CreateFloat(float64(left.Value), bits)
}

func (left *ConcreteFloatValue) AsFloat(bits int) Value {
	return left.Context.CreateFloat(left.Value, bits)
}

func (left *Z3Value) AsFloat(bits int) Value {
	switch {
	case left.IsFloat():
		switch left.Bits {
		case bits:
			return left
		default:
			return &Z3Value{
				Context: left.Context,
				Value:   left.Value.(z3.Float).ToReal().ToFloat(*left.Context.TypesContext.GetFloatSort(bits)),
			}
		}
	case left.IsInteger():
		return &Z3Value{
			left.Context,
			left.Value.(z3.BV).SToFloat(*left.Context.TypesContext.GetFloatSort(bits)),
			left.Bits,
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
		return left.Context.CreateInt(left.Value%castedRight.Value, left.bits)
	case *ConcreteFloatValue:
		panic("unsupported!")
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.asZ3IntValue().SRem(castedRight.Value.(z3.BV)),
			left.bits,
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
				leftBV.SRem(
					left.Context.Z3Context.FromInt(
						castedRight.Value, *left.Context.TypesContext.GetIntSort(castedRight.bits)).(z3.BV),
				),
				left.Bits,
			}
		case *ConcreteFloatValue:
			panic("unsupported!")
		case *Z3Value:
			return &Z3Value{
				left.Context,
				leftBV.SRem(castedRight.Value.(z3.BV)),
				left.Bits,
			}
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
func (left *ConcreteBoolValue) AsInt(bits int) Value {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) AsFloat() Value {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) asZ3IntValue() z3.BV {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) asZ3FloatValue(bits int) z3.Float {
	panic("Unsupported")
}
func (left *ConcreteBoolValue) Rem(Value) Value {
	panic("Unsupported")
}

func (left *ConcreteBoolValue) AsZ3Value() Z3Value {
	return Z3Value{
		left.Context,
		left.Context.Z3Context.FromBool(left.Value),
		1,
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
				1,
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
		return left.Context.CreateInt(left.Value&castedRight.Value, left.bits)
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.BV).And(castedRight.Value.(z3.BV)),
			left.bits,
		}
	}

	panic("unsupported!")
}
func (left *ConcreteIntValue) BitwiseOr(value ArithmeticValue) ArithmeticValue {
	switch castedRight := value.(type) {
	case *ConcreteIntValue:
		return left.Context.CreateInt(left.Value|castedRight.Value, left.bits)
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.BV).Or(castedRight.Value.(z3.BV)),
			left.bits,
		}
	}

	panic("unsupported!")
}
func (left *ConcreteIntValue) BitwiseXor(value ArithmeticValue) ArithmeticValue {
	switch castedRight := value.(type) {
	case *ConcreteIntValue:
		return left.Context.CreateInt(left.Value^castedRight.Value, left.bits)
	case *Z3Value:
		return &Z3Value{
			left.Context,
			left.AsZ3Value().Value.(z3.BV).Xor(castedRight.Value.(z3.BV)),
			left.bits,
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
			left.Bits,
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
			left.Bits,
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
			left.Bits,
		}
	}

	panic("unsupported!")
}

func (ctx *Context) CreateBool(value bool) *ConcreteBoolValue {
	return &ConcreteBoolValue{
		Context: ctx,
		Value:   value,
	}
}

func (ctx *Context) CreateInt(value int64, bits int) *ConcreteIntValue {
	return &ConcreteIntValue{
		Context: ctx,
		Value:   value,
		bits:    bits,
	}
}

func (ctx *Context) CreatePtrValue(value int) *ConcreteIntValue {
	return &ConcreteIntValue{
		Context: ctx,
		Value:   int64(value),
		bits:    ctx.TypesContext.Int64.Bits,
	}
}

func (ctx *Context) CreateFloat(value float64, bits int) *ConcreteFloatValue {
	return &ConcreteFloatValue{
		Context: ctx,
		Value:   value,
		bits:    bits,
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
