package interpreter

import "github.com/aclements/go-z3/z3"

type sortPtr string // type name is unique in the system

type Memory struct {
	ctx         *Context
	memoryLines map[sortPtr]z3.Array
	structures  map[sortPtr]*StructureDescriptor
}

type Pointer struct {
	ctx              *Context
	ptr              Value
	sPtr             sortPtr
	arrayElemSortPtr sortPtr
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

		// todo
		switch string(sPtr) {
		case "int":
			arrSort = mem.ctx.Z3Context.ArraySort(mem.ctx.TypesContext.Pointer, mem.ctx.TypesContext.IntSort)
		case "float":
			arrSort = mem.ctx.Z3Context.ArraySort(mem.ctx.TypesContext.Pointer, mem.ctx.TypesContext.FloatSort)
		default:
			// non default type
			arrSort = mem.ctx.Z3Context.ArraySort(mem.ctx.TypesContext.Pointer, mem.ctx.TypesContext.Pointer)
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
		return
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

func (mem *Memory) GetUnsafePointerToField(ptr Value, fieldIdx int, structName string) *Pointer {
	fieldSort := sortPtr(structName)
	if descr, ok := mem.structures[fieldSort]; ok {
		return &Pointer{
			ctx:  mem.ctx,
			ptr:  ptr,
			sPtr: descr.fields[fieldIdx],
		}
	}

	panic("unknown structure " + structName)
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

func (mem *Memory) initMemoryWrapper() {
	fields := make(map[int]string)
	fields[0] = "arrays-pointer"
	fields[1] = "int"

	mem.NewStruct("array-wrapper", fields)
	arrSort := mem.ctx.Z3Context.ArraySort(mem.ctx.TypesContext.Pointer, mem.ctx.TypesContext.IntSort)
	mem.memoryLines["int"] = mem.ctx.Z3Context.FreshConst("int-line", arrSort).(z3.Array)
}

func (mem *Memory) AllocateArray(elementType string) *Pointer {
	mem.initMemoryWrapper()

	wrapperPtr := mem.NewPtr("array-wrapper")
	elementSortPtr := sortPtr(elementType)
	elementsLineSortPtr := sortPtr("array-" + elementType)
	wrapperPtr.arrayElemSortPtr = elementSortPtr

	z3Ctx := mem.ctx.Z3Context
	elementSort := mem.ctx.TypesContext.Pointer

	lineSort := z3Ctx.ArraySort(
		mem.ctx.TypesContext.Pointer,
		z3Ctx.ArraySort(mem.ctx.TypesContext.ArrayIndexSort, elementSort),
	)
	lineArray := z3Ctx.Const(string(elementsLineSortPtr), lineSort).(z3.Array)

	mem.memoryLines[elementsLineSortPtr] = lineArray

	return wrapperPtr
}

func (mem *Memory) StoreAtArrayIndex(arrayPtr *Pointer, index int, value Value) {
	elementSortPtr := arrayPtr.arrayElemSortPtr
	elementsLineSortPtr := "array-" + elementSortPtr
	if _, ok := mem.memoryLines[elementSortPtr]; !ok {
		panic("unknown array of " + elementSortPtr)
	}

	valuePointer := mem.NewPtr(string(arrayPtr.arrayElemSortPtr))
	mem.Store(valuePointer, value)

	line := mem.memoryLines[elementsLineSortPtr]
	array := line.Select(arrayPtr.ptr.AsZ3Value().Value).(z3.Array)

	indexValue := mem.ctx.Z3Context.FromInt(int64(index), mem.ctx.TypesContext.ArrayIndexSort)
	array = array.Store(indexValue, valuePointer.AsZ3Value().Value)
	line = line.Store(arrayPtr.ptr.AsZ3Value().Value, array)
	mem.memoryLines[elementsLineSortPtr] = line
}

func (mem *Memory) LoadByArrayIndex(arrayPtr *Pointer, index Value) Value {
	elementSortPtr := arrayPtr.arrayElemSortPtr
	valuePointerValue := mem.GetArrayElementPointer(arrayPtr, index)

	valuePtr := &Pointer{
		ctx:  mem.ctx,
		ptr:  valuePointerValue,
		sPtr: elementSortPtr,
	}

	value := mem.Load(valuePtr)
	return value
}

func (mem *Memory) GetArrayElementPointer(arrayPtr *Pointer, index Value) *Pointer {
	elementSortPtr := arrayPtr.arrayElemSortPtr
	elementsLineSortPtr := "array-" + elementSortPtr
	if _, ok := mem.memoryLines[elementsLineSortPtr]; !ok {
		panic("unknown array of " + elementSortPtr)
	}

	line := mem.memoryLines[elementsLineSortPtr]
	array := line.Select(arrayPtr.ptr.AsZ3Value().Value).(z3.Array)

	valuePtrZ3Value := array.Select(index.AsZ3Value().Value).(z3.Value)
	valuePtrValue := &Z3Value{
		Context: mem.ctx,
		Value:   valuePtrZ3Value,
	}
	valuePtr := &Pointer{
		ctx:  mem.ctx,
		ptr:  valuePtrValue,
		sPtr: elementSortPtr,
	}

	return valuePtr
}

func (mem *Memory) SetArrayLen(arrayPtr *Pointer, len Value) {
	mem.StoreField(arrayPtr, 1, len)
}

func (mem *Memory) GetArrayLen(arrayPtr *Pointer) Value {
	return mem.LoadField(arrayPtr, 1)
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
