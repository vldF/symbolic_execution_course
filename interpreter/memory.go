package interpreter

import "github.com/aclements/go-z3/z3"

type sortPtr string // type name is unique in the system

type Memory struct {
	ctx         *Context
	memoryLines map[sortPtr]z3.Array
	structures  map[sortPtr]*StructureDescriptor
}

type Pointer struct {
	ctx  *Context
	ptr  Value
	sPtr sortPtr
}

type StructureDescriptor struct {
	fields map[int]sortPtr
}

var basePtrs = make(map[sortPtr]int64)

func getNextPtr(sortPtr sortPtr) int64 {
	if _, ok := basePtrs[sortPtr]; !ok {
		basePtrs[sortPtr] = 1 // we start from 1 because 0 represents nil
	}

	return basePtrs[sortPtr]
}

func (ptr *Pointer) IsNil() BoolValue {
	zeroIntConst := ConcreteIntValue{
		ptr.ctx,
		int64(0),
	}

	return ptr.ptr.Eq(&zeroIntConst)
}

func (mem *Memory) NewPtr(typeName string) *Pointer {
	intPtr := getNextPtr(sortPtr(typeName))
	newPtr := &Pointer{
		ctx: mem.ctx,
		ptr: &ConcreteIntValue{
			mem.ctx,
			intPtr,
		},
		sPtr: sortPtr(typeName),
	}

	return newPtr
}

func (mem *Memory) NullPtr(typeName string) *Pointer {
	newPtr := &Pointer{
		ctx: mem.ctx,
		ptr: &ConcreteIntValue{
			mem.ctx,
			0,
		},
		sPtr: sortPtr(typeName),
	}

	return newPtr
}

func (mem *Memory) Store(ptr *Pointer, value Value) {
	sPtr := ptr.sPtr
	if _, ok := mem.memoryLines[sPtr]; !ok {
		var arrSort z3.Sort

		switch string(sPtr) {
		case "int":
			arrSort = mem.ctx.Z3Context.ArraySort(mem.ctx.TypesContext.Pointer, mem.ctx.TypesContext.IntSort)
		case "float":
			arrSort = mem.ctx.Z3Context.ArraySort(mem.ctx.TypesContext.Pointer, mem.ctx.TypesContext.FloatSort)
		}

		mem.memoryLines[sPtr] = mem.ctx.Z3Context.FreshConst(string(sPtr)+"-line", arrSort).(z3.Array)
	}

	line := mem.memoryLines[sPtr]
	mem.memoryLines[sPtr] = line.Store(ptr.ptr.AsZ3Value().Value, value.AsZ3Value().Value)
}

func (mem *Memory) Load(ptr *Pointer) Value {
	sPtr := ptr.sPtr
	context := mem.ctx
	if _, ok := mem.memoryLines[sPtr]; !ok {
		panic("no memory line for the pointer")
	}

	line := mem.memoryLines[sPtr]
	z3Value := line.Select(ptr.ptr.AsZ3Value().Value)

	return &Z3Value{
		Context: context,
		Value:   z3Value,
	}
}

func (mem *Memory) NewStruct(name string, fields map[int]string) {
	structSortPtr := sortPtr(name)
	if _, ok := mem.structures[structSortPtr]; ok {
		panic("struct already exists " + structSortPtr)
	}

	fieldsInDescriptor := make(map[int]sortPtr)
	for fieldName, fieldTypeName := range fields {
		fieldsInDescriptor[fieldName] = sortPtr(fieldTypeName)
	}

	structDescriptor := &StructureDescriptor{
		fields: fieldsInDescriptor,
	}

	mem.structures[structSortPtr] = structDescriptor
}

func (mem *Memory) StoreField(structPtr *Pointer, fieldIdx int, value Value) {
	sPtr := structPtr.sPtr
	structDescr := mem.structures[sPtr]
	if _, ok := structDescr.fields[fieldIdx]; !ok {
		var fieldSortPtr sortPtr
		switch castedValue := value.(type) {
		case *ConcreteIntValue:
			fieldSortPtr = "int"
		case *ConcreteFloatValue:
			fieldSortPtr = "float"
		case *ConcreteBoolValue:
			fieldSortPtr = "bool"
		case *Z3Value:
			switch {
			case castedValue.IsFloat():
				fieldSortPtr = "float"
			case castedValue.IsInteger():
				fieldSortPtr = "int"
			case castedValue.IsBool():
				fieldSortPtr = "bool"
			default:
				panic("unsupported value type")
			}
		default:
			panic("unsupported value type")
		}

		structDescr.fields[fieldIdx] = fieldSortPtr
	}

	line := mem.memoryLines[structDescr.fields[fieldIdx]]
	mem.memoryLines[sPtr] = line.Store(structPtr.ptr.AsZ3Value().Value, value.AsZ3Value().Value)
}

func (mem *Memory) LoadField(structPtr *Pointer, fieldIdx int) Value {
	sPtr := structPtr.sPtr
	if _, ok := mem.structures[sPtr]; !ok {
		panic("unknown structure " + sPtr)
	}

	structDescr := mem.structures[sPtr]
	if _, ok := structDescr.fields[fieldIdx]; !ok {
		panic("unknown field " + string(rune(fieldIdx)))
	}

	fieldSortPtr := structDescr.fields[fieldIdx]
	line := mem.memoryLines[fieldSortPtr]

	z3Value := line.Select(structPtr.ptr.AsZ3Value().Value)

	return &Z3Value{
		Context: mem.ctx,
		Value:   z3Value,
	}
}

func (ptr *Pointer) AsZ3Value() Z3Value {
	//TODO implement me
	panic("implement me")
}

func (ptr *Pointer) Eq(value Value) BoolValue {
	//TODO implement me
	panic("implement me")
}

func (ptr *Pointer) NotEq(value Value) BoolValue {
	//TODO implement me
	panic("implement me")
}

func (ptr *Pointer) IsFloat() bool {
	//TODO implement me
	panic("implement me")
}

func (ptr *Pointer) IsInteger() bool {
	//TODO implement me
	panic("implement me")
}

func (ptr *Pointer) IsBool() bool {
	//TODO implement me
	panic("implement me")
}

func (ptr *Pointer) And(value Value) Value {
	//TODO implement me
	panic("implement me")
}

func (ptr *Pointer) Or(value Value) Value {
	//TODO implement me
	panic("implement me")
}

func (ptr *Pointer) Xor(value Value) Value {
	//TODO implement me
	panic("implement me")
}

func (mem *Memory) GetFieldPointer(structPtr *Pointer, fieldIdx int) Value {
	if _, ok := mem.structures[structPtr.sPtr]; !ok {
		panic("unknown structure " + structPtr.sPtr)
	}

	structDescr := mem.structures[structPtr.sPtr]
	fieldSortPtr := structDescr.fields[fieldIdx]

	return &Pointer{
		ctx:  mem.ctx,
		ptr:  structPtr.ptr,
		sPtr: fieldSortPtr,
	}
}
