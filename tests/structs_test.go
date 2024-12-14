package tests

import (
	"symbolic_execution_course/testdata"
	"testing"
)

func TestStruct_1(t *testing.T) {
	args := make(map[string]any)

	sArgs := make(map[int]any)
	sArgs[0] = 1
	sArgs[1] = 2.0

	args["s"] = StructArg{
		fields: sArgs,
	}

	expected := testdata.TestStruct(testdata.Struct1{IntField: 1, FloatField: 2.0})

	SymbolicMachineSatTest("structs", "TestStruct", args, expected, t)
	SymbolicMachineUnsatTest("structs", "TestStruct", args, expected+1, t)
}
