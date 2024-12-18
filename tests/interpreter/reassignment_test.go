package interpreter

import (
	"symbolic_execution_course/testdata"
	"symbolic_execution_course/tests"
	"testing"
)

func TestArrayReassignment_1(t *testing.T) {
	args := make(map[string]any)
	args["arr"] = tests.ArrayArg{Elements: []any{0}}

	expected := testdata.ArrayReassignment([]int{0})

	ctx := tests.PrepareTest("reassignment", "ArrayReassignment")

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
	tests.SymbolicMachineUnsatTest(ctx, args, -1, t)
}

func TestArrayReassignment_2(t *testing.T) {
	args := make(map[string]any)
	args["arr"] = tests.ArrayArg{Elements: []any{1, 2, 3, 4}}

	expected := testdata.ArrayReassignment([]int{1, 2, 3, 4})

	ctx := tests.PrepareTest("reassignment", "ArrayReassignment")

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
	tests.SymbolicMachineUnsatTest(ctx, args, -1, t)
}

func TestArrayReassignment_3(t *testing.T) {
	args := make(map[string]any)

	ctx := tests.PrepareTest("reassignment", "ArrayReassignment")

	// always returns 1, -1 is impossible
	tests.SymbolicMachineSatTest(ctx, args, 1, t)
	tests.SymbolicMachineUnsatTest(ctx, args, -1, t)
}
