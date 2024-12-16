package tests

import (
	"symbolic_execution_course/testdata"
	"testing"
)

func TestArrayReassignment_1(t *testing.T) {
	args := make(map[string]any)
	args["arr"] = ArrayArg{elements: []any{0}}

	expected := testdata.ArrayReassignment([]int{0})

	ctx := PrepareTest("reassignment", "ArrayReassignment")

	SymbolicMachineSatTest(ctx, args, expected, t)
	SymbolicMachineUnsatTest(ctx, args, expected+1, t)
	SymbolicMachineUnsatTest(ctx, args, -1, t)
}

func TestArrayReassignment_2(t *testing.T) {
	args := make(map[string]any)
	args["arr"] = ArrayArg{elements: []any{1, 2, 3, 4}}

	expected := testdata.ArrayReassignment([]int{1, 2, 3, 4})

	ctx := PrepareTest("reassignment", "ArrayReassignment")

	SymbolicMachineSatTest(ctx, args, expected, t)
	SymbolicMachineUnsatTest(ctx, args, expected+1, t)
	SymbolicMachineUnsatTest(ctx, args, -1, t)
}

func TestArrayReassignment_3(t *testing.T) {
	args := make(map[string]any)

	ctx := PrepareTest("reassignment", "ArrayReassignment")

	// always returns 1, -1 is impossible
	SymbolicMachineSatTest(ctx, args, 1, t)
	SymbolicMachineUnsatTest(ctx, args, -1, t)
}
