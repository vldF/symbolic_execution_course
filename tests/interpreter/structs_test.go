package interpreter

import (
	"symbolic_execution_course/testdata"
	"symbolic_execution_course/tests"
	"testing"
)

func TestStruct_1(t *testing.T) {
	args := make(map[string]any)

	ctx := tests.PrepareTest("structs", "TestStruct")

	sArgs := make(map[int]any)
	sArgs[0] = 1
	sArgs[1] = 2.0

	args["s"] = tests.StructArg{
		Fields: sArgs,
	}

	expected := testdata.TestStruct(testdata.Struct1{IntField: 1, FloatField: 2.0})

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}
