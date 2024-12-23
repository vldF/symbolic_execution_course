package interpreter

type Pointer struct {
	ctx              *Context
	ptr              Value
	sPtr             sortPtr
	arrayElemSortPtr sortPtr
}

func (ptr *Pointer) AsZ3Value() Z3Value {
	return ptr.ptr.AsZ3Value()
}

func (ptr *Pointer) Eq(value Value) BoolValue {
	switch value := value.(type) {
	case *Pointer:
		if ptr.sPtr != value.sPtr {
			return &ConcreteBoolValue{
				ptr.ctx,
				false,
			}
		}

		return ptr.ptr.Eq(value.ptr)
	}

	return &ConcreteBoolValue{
		ptr.ctx,
		false,
	}
}

func (ptr *Pointer) NotEq(value Value) BoolValue {
	return ptr.Eq(value).Not()
}

func (ptr *Pointer) IsFloat() bool {
	return false
}

func (ptr *Pointer) IsInteger() bool {
	return false
}

func (ptr *Pointer) IsBool() bool {
	return false
}

func (ptr *Pointer) And(Value) Value {
	panic("unsupported")
}

func (ptr *Pointer) Or(Value) Value {
	panic("unsupported")
}

func (ptr *Pointer) Xor(Value) Value {
	panic("unsupported")
}
