package interpreter

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
)

type Value interface {
	AsEq(Value) BoolPredicate
}

type BoolPredicate struct {
	Context *Context
	Value   z3.Bool
}

func (left *BoolPredicate) String() string {
	return fmt.Sprintf("BoolPredicate(%v)", left.Value)
}

func (left *ConcreteIntValue) AsEq(right Value) BoolPredicate {
	switch castedRight := right.(type) {
	case *ConcreteIntValue:
		return BoolPredicate{
			left.Context,
			left.asZ3IntValue().Eq(castedRight.asZ3IntValue()),
		}

	case *ConcreteFloatValue:
		return BoolPredicate{
			left.Context,
			left.asZ3FloatValue().Eq(castedRight.asZ3FloatValue()),
		}

	case *Z3Value:
		return right.AsEq(left)
	}

	panic("unreachable")
}

func (left *ConcreteFloatValue) AsEq(right Value) BoolPredicate {
	switch castedRight := right.(type) {
	case *ConcreteIntValue:
		return BoolPredicate{
			left.Context,
			left.asZ3IntValue().Eq(castedRight.asZ3IntValue()),
		}

	case *ConcreteFloatValue:
		return BoolPredicate{
			left.Context,
			left.asZ3FloatValue().Eq(castedRight.asZ3FloatValue()),
		}
	}

	panic("unreachable")
}

type NumericValue interface {
	Add(NumericValue) NumericValue
	Sub(NumericValue) NumericValue
	Mul(NumericValue) NumericValue
	Div(NumericValue) NumericValue

	AsZ3Value() Z3Value

	asZ3IntValue() z3.BV
	asZ3FloatValue() z3.Float

	IsFloat() bool
	IsInteger() bool
}

type ConcreteValue[V int64 | bool | float64] struct {
	Context *Context
	Value   V
}

type ConcreteIntValue ConcreteValue[int64]
type ConcreteFloatValue ConcreteValue[float64]

func (left *ConcreteValue[T]) String() string {
	return fmt.Sprintf("ConcretValue(%v)", left.Value)
}

func (left *ConcreteIntValue) Add(v NumericValue) NumericValue {
	switch castedRight := v.(type) {
	case *ConcreteIntValue:
		return &ConcreteIntValue{
			left.Context,
			left.Value + castedRight.Value,
		}
	case *ConcreteFloatValue:
		return &ConcreteFloatValue{
			left.Context,
			float64(left.Value) + castedRight.Value,
		}
	}

	return nil
}

func (left *ConcreteIntValue) AsZ3Value() Z3Value {
	return Z3Value{
		left.Context,
		left.Context.Z3Context.FromInt(
			left.Value,
			left.Context.TypesContext.IntSort).(z3.Int).ToBV(left.Context.TypesContext.IntBits),
	}
}

func (left *ConcreteFloatValue) AsZ3Value() Z3Value {
	return Z3Value{
		left.Context,
		left.Context.Z3Context.FromFloat64(
			left.Value,
			left.Context.TypesContext.IntSort),
	}
}

func (left *ConcreteIntValue) Sub(v NumericValue) NumericValue {
	switch castedRight := v.(type) {
	case *ConcreteIntValue:
		return &ConcreteIntValue{
			left.Context,
			left.Value - castedRight.Value,
		}
	case *ConcreteFloatValue:
		return &ConcreteFloatValue{
			left.Context,
			float64(left.Value) - castedRight.Value}
	}

	return nil
}

func (left *ConcreteIntValue) Mul(v NumericValue) NumericValue {
	switch castedRight := v.(type) {
	case *ConcreteIntValue:
		return &ConcreteIntValue{
			left.Context,
			left.Value * castedRight.Value,
		}
	case *ConcreteFloatValue:
		return &ConcreteFloatValue{
			left.Context,
			float64(left.Value) * castedRight.Value,
		}
	}

	return nil
}

func (left *ConcreteIntValue) Div(v NumericValue) NumericValue {
	switch castedRight := v.(type) {
	case *ConcreteIntValue:
		return &ConcreteIntValue{
			left.Context,
			left.Value / castedRight.Value,
		}
	case *ConcreteFloatValue:
		return &ConcreteFloatValue{
			left.Context,
			float64(left.Value) / castedRight.Value,
		}
	}

	return nil
}

func (left *ConcreteFloatValue) Add(v NumericValue) NumericValue {
	switch castedRight := v.(type) {
	case *ConcreteIntValue:
		return &ConcreteFloatValue{
			left.Context,
			left.Value + float64(castedRight.Value),
		}
	case *ConcreteFloatValue:
		return &ConcreteFloatValue{
			left.Context,
			left.Value + castedRight.Value,
		}
	}

	return nil
}

func (left *ConcreteFloatValue) Sub(v NumericValue) NumericValue {
	switch castedRight := v.(type) {
	case *ConcreteIntValue:
		return &ConcreteFloatValue{
			left.Context,
			left.Value - float64(castedRight.Value),
		}
	case *ConcreteFloatValue:
		return &ConcreteFloatValue{
			left.Context,
			left.Value - castedRight.Value,
		}
	}

	return nil
}

func (left *ConcreteFloatValue) Mul(v NumericValue) NumericValue {
	switch castedRight := v.(type) {
	case *ConcreteIntValue:
		return &ConcreteFloatValue{
			left.Context,
			left.Value * float64(castedRight.Value),
		}
	case *ConcreteFloatValue:
		return &ConcreteFloatValue{
			left.Context,
			left.Value * castedRight.Value,
		}
	}

	return nil
}

func (left *ConcreteFloatValue) Div(v NumericValue) NumericValue {
	switch castedRight := v.(type) {
	case *ConcreteIntValue:
		return &ConcreteFloatValue{
			left.Context,
			left.Value / float64(castedRight.Value),
		}
	case *ConcreteFloatValue:
		return &ConcreteFloatValue{
			left.Context,
			left.Value / castedRight.Value,
		}
	}

	return nil
}

type ConcreteBoolValue ConcreteValue[bool]

type Z3Value struct {
	Context *Context
	Value   z3.Value
}

func (left *Z3Value) AsEq(right Value) BoolPredicate {
	switch {
	case left.IsInteger():
		switch castedRight := right.(type) {
		case *ConcreteIntValue:
			return BoolPredicate{
				left.Context,
				left.asZ3IntValue().Eq(castedRight.asZ3IntValue()),
			}
		case *ConcreteFloatValue:
			return BoolPredicate{
				left.Context,
				left.asZ3FloatValue().Eq(castedRight.asZ3FloatValue()),
			}
		case *Z3Value:
			switch {
			case castedRight.IsInteger():
				return BoolPredicate{
					left.Context,
					left.asZ3IntValue().Eq(castedRight.asZ3IntValue()),
				}
			case castedRight.IsFloat():
				return BoolPredicate{
					left.Context,
					left.asZ3FloatValue().Eq(castedRight.asZ3FloatValue()),
				}
			}
		}
	case left.IsFloat():
		switch castedRight := right.(type) {
		case *ConcreteIntValue:
			return BoolPredicate{
				left.Context,
				left.asZ3FloatValue().Eq(castedRight.asZ3FloatValue()),
			}
		case *ConcreteFloatValue:
			return BoolPredicate{
				left.Context,
				left.asZ3FloatValue().Eq(castedRight.asZ3FloatValue()),
			}
		case *Z3Value:
			switch {
			case castedRight.IsInteger():
				return BoolPredicate{
					left.Context,
					left.asZ3FloatValue().Eq(castedRight.asZ3FloatValue()),
				}
			case castedRight.IsFloat():
				return BoolPredicate{
					left.Context,
					left.asZ3FloatValue().Eq(castedRight.asZ3FloatValue()),
				}
			}
		}
	}

	panic("unreachable")
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

func (left *ConcreteIntValue) IsFloat() bool {
	return false
}

func (left *ConcreteFloatValue) IsFloat() bool {
	return true
}

func (left *Z3Value) IsFloat() bool {
	if _, ok := left.Value.(z3.Float); ok {
		return true
	}

	return false
}

func (left *Z3Value) Add(right NumericValue) NumericValue {
	switch left.IsFloat() {
	case true:
		return &Z3Value{
			left.Context,
			left.Value.(z3.Float).Add(right.asZ3FloatValue())}
	case false:
		switch right.IsFloat() {
		case true:
			return &Z3Value{
				left.Context,
				left.asZ3FloatValue().Add(right.asZ3FloatValue())}
		case false:
			return &Z3Value{
				left.Context,
				left.asZ3IntValue().Add(right.asZ3IntValue())}
		}
	}

	return nil
}

func (left *Z3Value) Sub(right NumericValue) NumericValue {
	switch left.IsFloat() {
	case true:
		return &Z3Value{
			left.Context,
			left.Value.(z3.Float).Sub(right.asZ3FloatValue())}
	case false:
		switch right.IsFloat() {
		case true:
			return &Z3Value{
				left.Context,
				left.asZ3FloatValue().Sub(right.asZ3FloatValue())}
		case false:
			return &Z3Value{
				left.Context,
				left.asZ3IntValue().Sub(right.asZ3IntValue())}
		}
	}

	return nil
}

func (left *Z3Value) Mul(right NumericValue) NumericValue {
	switch left.IsFloat() {
	case true:
		return &Z3Value{
			left.Context,
			left.Value.(z3.Float).Mul(right.asZ3FloatValue())}
	case false:
		switch right.IsFloat() {
		case true:
			return &Z3Value{
				left.Context,
				left.asZ3FloatValue().Mul(right.asZ3FloatValue())}
		case false:
			return &Z3Value{
				left.Context,
				left.asZ3IntValue().Mul(right.asZ3IntValue())}
		}
	}

	return nil
}

func (left *Z3Value) Div(right NumericValue) NumericValue {
	switch left.IsFloat() {
	case true:
		return &Z3Value{
			left.Context,
			left.asZ3FloatValue().Div(right.asZ3FloatValue())}
	case false:
		switch right.IsFloat() {
		case true:
			return &Z3Value{
				left.Context,
				left.asZ3FloatValue().Div(right.asZ3FloatValue())}
		case false:
			return &Z3Value{
				left.Context,
				left.asZ3IntValue().SDiv(right.asZ3IntValue())}
		}
	}

	return nil
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
