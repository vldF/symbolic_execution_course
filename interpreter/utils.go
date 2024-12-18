package interpreter

import (
	"fmt"
	"github.com/aclements/go-z3/z3"
	"go/types"
	"strconv"
)

func (ctx *Context) TypeToSort(t types.Type) z3.Sort {
	switch casted := t.(type) {
	case *types.Basic:
		switch t.(*types.Basic).Kind() {
		case types.UntypedComplex, types.Complex64, types.Complex128:
			return ctx.TypesContext.Pointer
		case types.String:
			return ctx.TypesContext.UnknownSort
		default:
			return *ctx.TypesContext.GetPrimitiveTypeSortOrNil(casted.Name())
		}

	case *types.Array:
		elemType := t.(*types.Array).Elem()
		return ctx.Z3Context.ArraySort(ctx.TypesContext.ArrayIndexSort, ctx.TypeToSort(elemType))
	case *types.Slice:
		elemType := t.(*types.Slice).Elem()
		return ctx.Z3Context.ArraySort(ctx.TypesContext.ArrayIndexSort, ctx.TypeToSort(elemType))
	case *types.Named:
		return ctx.TypesContext.Pointer
	case *types.Pointer:
		return ctx.TypesContext.Pointer
	}

	panic("can't get sort")
}

func FloatToString(f z3.Float) string {
	float, _ := f.AsBigFloat()

	return fmt.Sprintf("%d", float)
}

func GetStructureFields(s *types.Named) map[int]string {
	castedStruct := s.Underlying().(*types.Struct)
	fieldsCount := castedStruct.NumFields()
	result := make(map[int]string, fieldsCount)
	for i := 0; i < fieldsCount; i++ {
		result[i] = GetTypeName(castedStruct.Field(i).Type())
	}

	return result
}

func GetTypeName(s types.Type) string {
	switch castedType := s.(type) {
	case *types.Named:
		return castedType.String()
	case *types.Basic:
		return castedType.String()
	case *types.Pointer:
		return GetTypeName(castedType.Elem())
	}

	panic("can't get type " + s.String())
}

func (ctx *TypesContext) GetPrimitiveTypeBits(typeName string) int {
	switch typeName {
	case "int":
		return ctx.Int.Bits
	case "int8":
		return ctx.Int8.Bits
	case "int16":
		return ctx.Int16.Bits
	case "int32":
		return ctx.Int32.Bits
	case "int64":
		return ctx.Int64.Bits
	case "uint":
		return ctx.UInt.Bits
	case "uint8":
		return ctx.UInt8.Bits
	case "uint16":
		return ctx.UInt16.Bits
	case "uint32":
		return ctx.UInt32.Bits
	case "uint64":
		return ctx.UInt64.Bits
	case "float":
		return ctx.Float.Bits
	case "float32":
		return ctx.Float32.Bits
	case "float64":
		return ctx.Float64.Bits
	default:
		panic("unsupported type")
	}
}

func (ctx *TypesContext) GetPrimitiveTypeSortOrNil(typeName string) *z3.Sort {
	switch typeName {
	case "int":
		return &ctx.Int.Sort
	case "int8":
		return &ctx.Int8.Sort
	case "int16":
		return &ctx.Int16.Sort
	case "int32":
		return &ctx.Int32.Sort
	case "int64":
		return &ctx.Int64.Sort
	case "uint":
		return &ctx.UInt.Sort
	case "uint8":
		return &ctx.UInt8.Sort
	case "uint16":
		return &ctx.UInt16.Sort
	case "uint32":
		return &ctx.UInt32.Sort
	case "uint64":
		return &ctx.UInt64.Sort
	case "float":
		return &ctx.Float.Sort
	case "float32":
		return &ctx.Float32.Sort
	case "float64":
		return &ctx.Float64.Sort
	default:
		return nil
	}
}

func (ctx *TypesContext) GetIntSortDescr(bits int) *PrimitiveIntDescr {
	switch bits {
	case 8:
		return &ctx.Int8
	case 16:
		return &ctx.Int16
	case 32:
		return &ctx.Int32
	case 64:
		return &ctx.Int64
	default:
		panic("invalid Bits")
	}
}

func (ctx *TypesContext) GetIntSort(bits int) *z3.Sort {
	return &ctx.GetIntSortDescr(bits).Sort
}

func (ctx *TypesContext) GetFloatSort(bits int) *z3.Sort {
	switch bits {
	case 32:
		return &ctx.Float32.Sort
	case 64:
		return &ctx.Float64.Sort
	default:
		panic("invalid Bits " + strconv.Itoa(bits))
	}
}
