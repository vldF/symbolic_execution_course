package interpreter

import (
	"github.com/aclements/go-z3/z3"
	"go/types"
	"math"
)

var intPtr = sortPointer{
	ConcreteIntValue{
		Context: nil,
		Value:   math.MaxInt64 - 1,
	},
}
var floatPtr = sortPointer{
	ConcreteIntValue{
		Context: nil,
		Value:   math.MaxInt64 - 2,
	},
}

var basePtrCounter int64 = 0
var baseSortPtrCounter int64 = 100000000

type Memory struct {
	context       *Context
	Mem           map[sortPointer]interface{}
	TypeToSortPtr map[string]sortPointer
}

func (memory *Memory) AllocateInt() ValuePointer {
	return memory.getNextPtr()
}

func (memory *Memory) AllocateFloat() ValuePointer {
	return memory.getNextPtr()
}

func (memory *Memory) AllocateStruct() ValuePointer {
	return memory.getNextPtr()
}

func (memory *Memory) NewArray(elementName string, elements types.BasicKind) {
	z3Context := memory.context.Z3Context

	basePtrCounter++
	newSortPtr := sortPointer{
		ConcreteIntValue{
			memory.context,
			basePtrCounter,
		},
	}
	memory.TypeToSortPtr[elementName+"-array-wrapper"] = newSortPtr

	basePtrCounter++
	lenSortPtr := sortPointer{
		ConcreteIntValue{
			memory.context,
			basePtrCounter,
		},
	}
	intArrSort := z3Context.ArraySort(memory.context.TypesContext.IntSort, memory.context.TypesContext.IntSort)
	memory.Mem[lenSortPtr] = &PrimitiveValueCell{
		z3Context.FreshConst(elementName+"-array-wrapper-len", intArrSort).(z3.Array),
	}

	basePtrCounter++
	elementSortPtr := sortPointer{
		ConcreteIntValue{
			memory.context,
			basePtrCounter,
		},
	}

	switch elements {
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
		intArrSort := z3Context.ArraySort(memory.context.TypesContext.IntSort, memory.context.TypesContext.IntSort)
		arrayOfIntArrSort := z3Context.ArraySort(memory.context.TypesContext.IntSort, intArrSort)
		memory.Mem[elementSortPtr] = &PrimitiveValueCell{
			z3Context.FreshConst(elementName+"-ints", arrayOfIntArrSort).(z3.Array),
		}
	case types.Float32, types.Float64, types.UntypedFloat:
		intArrSort := z3Context.ArraySort(memory.context.TypesContext.IntSort, memory.context.TypesContext.IntSort)
		arrayOfIntArrSort := z3Context.ArraySort(memory.context.TypesContext.IntSort, intArrSort)
		memory.Mem[elementSortPtr] = &PrimitiveValueCell{
			z3Context.FreshConst(elementName+"-ints", arrayOfIntArrSort).(z3.Array),
		}
	}

	memory.Mem[newSortPtr] = ArrayWrapperCell{
		memory:   memory,
		Lens:     lenSortPtr,
		Elements: elementSortPtr,
	}

}

func (memory *Memory) NewStruct(name string, fields map[int]types.BasicKind) {
	basePtrCounter++
	newSortPtr := sortPointer{
		ConcreteIntValue{
			memory.context,
			basePtrCounter,
		},
	}
	memory.TypeToSortPtr[name] = newSortPtr

	fieldsInCell := make(map[int]sortPointer)

	for fieldName, fieldType := range fields {
		baseSortPtrCounter++
		fieldsPtr := sortPointer{
			ConcreteIntValue{
				Context: memory.context,
				Value:   baseSortPtrCounter,
			},
		}

		z3Context := memory.context.Z3Context
		typesContext := memory.context.TypesContext
		switch fieldType {
		case types.Int8, types.Int16, types.Int32, types.Int64, types.Int, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			intArrSort := z3Context.ArraySort(typesContext.IntSort, typesContext.IntSort)
			memory.Mem[fieldsPtr] = &PrimitiveValueCell{z3Context.FreshConst(name+"-ints", intArrSort).(z3.Array)}
		case types.Float32, types.Float64:
			floatArrSort := z3Context.ArraySort(typesContext.IntSort, typesContext.FloatSort)
			memory.Mem[fieldsPtr] = &PrimitiveValueCell{z3Context.FreshConst(name+"-floats", floatArrSort).(z3.Array)}
		default:
			panic("unsupported type")
		}

		fieldsInCell[fieldName] = fieldsPtr
	}

	memory.Mem[newSortPtr] = StructValueCell{
		memory: memory,
		Fields: fieldsInCell,
	}
}

func (memory *Memory) getNextPtr() ValuePointer {
	basePtrCounter++

	return ValuePointer{
		memory.context,
		&ConcreteIntValue{
			memory.context,
			basePtrCounter,
		},
	}
}

func (memory *Memory) GetIntValue(pointer ValuePointer) Value {
	cell := memory.Mem[intPtr].(PrimitiveValueCell)

	return cell.getValue(pointer, memory.context)
}

func (memory *Memory) GetFloatValue(pointer ValuePointer) Value {
	cell := memory.Mem[floatPtr].(PrimitiveValueCell)

	return cell.getValue(pointer, memory.context)
}

func (memory *Memory) GetStructField(structPtr StructPointer, fieldIdx int) Value {
	structure := memory.Mem[structPtr.SortPtr].(StructValueCell)
	fieldsPtr := structure.Fields[fieldIdx]
	fields := memory.Mem[fieldsPtr]
	valueCell := fields.(*PrimitiveValueCell)

	return valueCell.getValue(structPtr.Ptr, memory.context)
}

type sortPointer struct {
	Value ConcreteIntValue
}

type ValuePointer struct {
	context *Context
	value   Value
}

type StructPointer struct {
	context    *Context
	SortPtr    sortPointer
	Ptr        ValuePointer
	structName string
}

func (s StructPointer) AsZ3Value() Z3Value {
	//TODO implement me
	panic("implement me")
}

func (s StructPointer) Eq(value Value) BoolValue {
	switch castedValue := value.(type) {
	case StructPointer:
		return &Z3Value{
			s.context,
			s.Ptr.AsZ3Value().Value.(z3.BV).Eq(castedValue.Ptr.AsZ3Value().Value.(z3.BV)).
				And(s.SortPtr.Value.AsZ3Value().Value.(z3.BV).Eq(s.SortPtr.Value.AsZ3Value().Value.(z3.BV))),
		}
	default:
		panic("unsupported type")
	}
}

func (s StructPointer) NotEq(value Value) BoolValue {
	return s.Eq(value).Not()
}

func (s StructPointer) IsFloat() bool {
	return false
}

func (s StructPointer) IsInteger() bool {
	return false
}

func (s StructPointer) IsBool() bool {
	return false
}

func (s StructPointer) And(value Value) Value {
	//TODO implement me
	panic("implement me")
}

func (s StructPointer) Or(value Value) Value {
	//TODO implement me
	panic("implement me")
}

func (s StructPointer) Xor(value Value) Value {
	//TODO implement me
	panic("implement me")
}

func (v ValuePointer) AsZ3Value() Z3Value {
	return v.value.AsZ3Value()
}

func (v ValuePointer) Eq(value Value) BoolValue {
	switch castedValue := value.(type) {
	case ValuePointer:
		return &Z3Value{
			v.context,
			v.value.AsZ3Value().Value.(z3.BV).Eq(castedValue.value.AsZ3Value().Value.(z3.BV)),
		}
	default:
		panic("unsupported type")
	}
}

func (v ValuePointer) NotEq(value Value) BoolValue {
	return v.Eq(value).Not()
}

func (v ValuePointer) IsFloat() bool {
	return false
}

func (v ValuePointer) IsInteger() bool {
	return false
}

func (v ValuePointer) IsBool() bool {
	return false
}

func (v ValuePointer) And(value Value) Value {
	//TODO implement me
	panic("implement me")
}

func (v ValuePointer) Or(value Value) Value {
	//TODO implement me
	panic("implement me")
}

func (v ValuePointer) Xor(value Value) Value {
	//TODO implement me
	panic("implement me")
}

type PrimitiveValueCell struct {
	Z3Arr z3.Array
}

func (cell *PrimitiveValueCell) getValue(index ValuePointer, context *Context) Value {
	return &Z3Value{
		Context: context,
		Value:   cell.Z3Arr.Select(index.value.AsZ3Value().Value),
	}
}

type StructValueCell struct {
	memory *Memory
	Fields map[int]sortPointer
}

type ArrayWrapperCell struct {
	memory   *Memory
	Lens     sortPointer
	Elements sortPointer
}

func (cell *ArrayWrapperCell) GetValue(index ValuePointer, context *Context) Value {
	z3Val := cell.memory.Mem[cell.Elements].(*PrimitiveValueCell).Z3Arr.Select(index.value.AsZ3Value().Value.(z3.BV))

	return &Z3Value{
		Context: context,
		Value:   z3Val,
	}
}

func (cell *ArrayWrapperCell) GetLen(index ValuePointer, context *Context) Value {
	z3Val := cell.memory.Mem[cell.Lens].(*PrimitiveValueCell).Z3Arr.Select(index.value.AsZ3Value().Value.(z3.BV))

	return &Z3Value{
		Context: context,
		Value:   z3Val,
	}
}
