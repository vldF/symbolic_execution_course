package tests

import (
	"go/types"
	"symbolic_execution_course/testdata"
	"testing"
)

func TestStruct_1(t *testing.T) {
	args := make(map[string]any)

	sArgs := make(map[int]any)
	sArgs[0] = 1
	sArgs[1] = 2.0
	sTypes := make(map[int]types.BasicKind)
	sTypes[0] = types.Int
	sTypes[1] = types.Float64

	args["s"] = StructArg{
		typeName:    "Struct1",
		fields:      sArgs,
		fieldsTypes: sTypes,
	}

	expected := testdata.TestStruct(testdata.Struct1{IntField: 1, FloatField: 2.0})

	SymbolicMachineSatTest("structs", "TestStruct", args, expected, t)
	SymbolicMachineUnsatTest("structs", "TestStruct", args, expected+1, t)
}
