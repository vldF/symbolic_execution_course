package tests

import (
	"symbolic_execution_course/testdata"
	"testing"
)

func TestIdInt_1(t *testing.T) {
	args := make(map[string]any)

	ctx := PrepareTest("simple", "IdInt")

	args["x"] = 1
	expected := 1

	SymbolicMachineSatTest(ctx, args, expected, t)
	SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestIdInt_2(t *testing.T) {
	args := make(map[string]any)

	ctx := PrepareTest("simple", "IdInt")

	args["x"] = 2
	expected := 2

	SymbolicMachineSatTest(ctx, args, expected, t)
	SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestIdFloat_1(t *testing.T) {
	args := make(map[string]any)

	ctx := PrepareTest("simple", "IdFloat")

	args["x"] = 1.0
	expected := 1.0

	SymbolicMachineSatTest(ctx, args, expected, t)
	SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestIdFloat_2(t *testing.T) {
	args := make(map[string]any)

	ctx := PrepareTest("simple", "IdFloat")

	args["x"] = 2.0
	expected := 2.0

	SymbolicMachineSatTest(ctx, args, expected, t)
	SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestSimpleExpressionInt_1(t *testing.T) {
	args := make(map[string]any)

	ctx := PrepareTest("simple", "SimpleExpressionInt")

	args["x"] = 1
	expected := testdata.SimpleExpressionInt(1)

	SymbolicMachineSatTest(ctx, args, expected, t)
	SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestSimpleExpressionInt_2(t *testing.T) {
	args := make(map[string]any)

	ctx := PrepareTest("simple", "SimpleExpressionInt")

	args["x"] = 2
	expected := testdata.SimpleExpressionInt(2)

	SymbolicMachineSatTest(ctx, args, expected, t)
	SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}
