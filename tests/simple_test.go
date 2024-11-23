package tests

import (
	"symbolic_execution_course/testdata"
	"testing"
)

func TestIdInt_1(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 1
	expected := 1

	SymbolicMachineSatTest("simple", "IdInt", args, expected, t)
	SymbolicMachineUnsatTest("simple", "IdInt", args, expected+1, t)
}

func TestIdInt_2(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 2
	expected := 2

	SymbolicMachineSatTest("simple", "IdInt", args, expected, t)
	SymbolicMachineUnsatTest("simple", "IdInt", args, expected+1, t)
}

func TestIdFloat_1(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 1.0
	expected := 1.0

	SymbolicMachineSatTest("simple", "IdFloat", args, expected, t)
	SymbolicMachineUnsatTest("simple", "IdFloat", args, expected+1, t)
}

func TestIdFloat_2(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 2.0
	expected := 2.0

	SymbolicMachineSatTest("simple", "IdFloat", args, expected, t)
	SymbolicMachineUnsatTest("simple", "IdFloat", args, expected+1, t)
}

func TestSimpleExpressionInt_1(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 1
	expected := testdata.SimpleExpressionInt(1)

	SymbolicMachineSatTest("simple", "SimpleExpressionInt", args, expected, t)
	SymbolicMachineUnsatTest("simple", "SimpleExpressionInt", args, expected+1, t)
}

func TestSimpleExpressionInt_2(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 2
	expected := testdata.SimpleExpressionInt(2)

	SymbolicMachineSatTest("simple", "SimpleExpressionInt", args, expected, t)
	SymbolicMachineUnsatTest("simple", "SimpleExpressionInt", args, expected+1, t)
}
