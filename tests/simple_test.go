package tests

import (
	"symbolic_execution_course/testdata"
	"testing"
)

func TestIdInt_1(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 1
	expected := 1

	SymbolicMachineTest("simple", "IdInt", args, expected, t)
}

func TestIdInt_2(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 2
	expected := 2

	SymbolicMachineTest("simple", "IdInt", args, expected, t)
}

func TestIdFloat_1(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 1.0
	expected := 1.0

	SymbolicMachineTest("simple", "IdFloat", args, expected, t)
}

func TestIdFloat_2(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 2.0
	expected := 2.0

	SymbolicMachineTest("simple", "IdFloat", args, expected, t)
}

func TestSimpleExpressionInt_1(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 1
	expected := testdata.SimpleExpressionInt(1)

	SymbolicMachineTest("simple", "SimpleExpressionInt", args, expected, t)
}

func TestSimpleExpressionInt_2(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 2
	expected := testdata.SimpleExpressionInt(2)

	SymbolicMachineTest("simple", "SimpleExpressionInt", args, expected, t)
}