package tests

import (
	"symbolic_execution_course/testdata"
	"testing"
)

func TestIntegerOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 1

	expected := testdata.IntegerOperations(2, 1)

	SymbolicMachineSatTest("numbers", "IntegerOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "IntegerOperations", args, expected+1, t)
}

func TestIntegerOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 2

	expected := testdata.IntegerOperations(1, 2)

	SymbolicMachineSatTest("numbers", "IntegerOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "IntegerOperations", args, expected+1, t)
}

func TestIntegerOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 2

	expected := testdata.IntegerOperations(2, 2)

	SymbolicMachineSatTest("numbers", "IntegerOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "IntegerOperations", args, expected+1, t)
}

func TestFloatOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 5.0
	args["y"] = 2.0

	expected := testdata.FloatOperations(5, 2)

	SymbolicMachineSatTest("numbers", "FloatOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "FloatOperations", args, expected+1, t)
}

func TestFloatOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 2.0
	args["y"] = 5.0

	expected := testdata.FloatOperations(2, 5)

	SymbolicMachineSatTest("numbers", "FloatOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "FloatOperations", args, expected+1, t)
}

func TestFloatOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 2.0
	args["y"] = 2.0

	expected := testdata.FloatOperations(2, 2)

	SymbolicMachineSatTest("numbers", "FloatOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "FloatOperations", args, expected+1, t)
}

func TestMixedOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 3.0

	expected := testdata.MixedOperations(2, 3)

	SymbolicMachineSatTest("numbers", "MixedOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "MixedOperations", args, expected+1, t)
}

func TestMixedOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 3.0

	expected := testdata.MixedOperations(3, 3)

	SymbolicMachineSatTest("numbers", "MixedOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "MixedOperations", args, expected+1, t)
}

func TestMixedOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 4.0

	expected := testdata.MixedOperations(3, 4)

	SymbolicMachineSatTest("numbers", "MixedOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "MixedOperations", args, expected+1, t)
}

func TestMixedOperations_4(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 9
	args["b"] = 4.0

	expected := testdata.MixedOperations(9, 4)

	SymbolicMachineSatTest("numbers", "MixedOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "MixedOperations", args, expected+1, t)
}

func TestMixedOperations_5(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 7.0

	expected := testdata.MixedOperations(4, 7)

	SymbolicMachineSatTest("numbers", "MixedOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "MixedOperations", args, expected+1, t)
}

func TestMixedOperations_6(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = -7.0

	expected := testdata.MixedOperations(3, -7)

	SymbolicMachineSatTest("numbers", "MixedOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "MixedOperations", args, expected+1, t)
}

func TestNestedConditions_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 2.0

	expected := testdata.NestedConditions(1, 2)

	SymbolicMachineSatTest("numbers", "NestedConditions", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "NestedConditions", args, expected+1, t)
}

func TestNestedConditions_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = -1
	args["b"] = 2.0

	expected := testdata.NestedConditions(-1, 2)

	SymbolicMachineSatTest("numbers", "NestedConditions", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "NestedConditions", args, expected+1, t)
}

func TestNestedConditions_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = -1
	args["b"] = -2.0

	expected := testdata.NestedConditions(-1, -2)

	SymbolicMachineSatTest("numbers", "NestedConditions", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "NestedConditions", args, expected+1, t)
}

func TestBitwiseOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 4

	expected := testdata.BitwiseOperations(2, 4)

	SymbolicMachineSatTest("numbers", "BitwiseOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "BitwiseOperations", args, expected+1, t)
}

func TestBitwiseOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 5

	expected := testdata.BitwiseOperations(3, 5)

	SymbolicMachineSatTest("numbers", "BitwiseOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "BitwiseOperations", args, expected+1, t)
}

func TestBitwiseOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 5

	expected := testdata.BitwiseOperations(4, 5)

	SymbolicMachineSatTest("numbers", "BitwiseOperations", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "BitwiseOperations", args, expected+1, t)
}

func TestAdvancedBitwise_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 5
	args["b"] = 4

	expected := testdata.AdvancedBitwise(5, 4)

	SymbolicMachineSatTest("numbers", "AdvancedBitwise", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "AdvancedBitwise", args, expected+1, t)
}

func TestAdvancedBitwise_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 5

	expected := testdata.AdvancedBitwise(4, 5)

	SymbolicMachineSatTest("numbers", "AdvancedBitwise", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "AdvancedBitwise", args, expected+1, t)
}

func TestAdvancedBitwise_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 4

	expected := testdata.AdvancedBitwise(4, 4)

	SymbolicMachineSatTest("numbers", "AdvancedBitwise", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "AdvancedBitwise", args, expected+1, t)
}

func TestCombinedBitwise_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 1

	expected := testdata.CombinedBitwise(2, 1)

	SymbolicMachineSatTest("numbers", "CombinedBitwise", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "CombinedBitwise", args, expected+1, t)
}

func TestCombinedBitwise_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 1

	expected := testdata.CombinedBitwise(3, 1)

	SymbolicMachineSatTest("numbers", "CombinedBitwise", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "CombinedBitwise", args, expected+1, t)
}

func TestCombinedBitwise_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 0b1111
	args["b"] = 0b101111

	expected := testdata.CombinedBitwise(0b1111, 0b101111)

	SymbolicMachineSatTest("numbers", "CombinedBitwise", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "CombinedBitwise", args, expected+1, t)
}

func TestNestedBitwise_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = -1
	args["b"] = 0

	expected := testdata.NestedBitwise(-1, 0)

	SymbolicMachineSatTest("numbers", "NestedBitwise", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "NestedBitwise", args, expected+1, t)
}

func TestNestedBitwise_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 0b10101
	args["b"] = -1

	expected := testdata.NestedBitwise(0b10101, -1)

	SymbolicMachineSatTest("numbers", "NestedBitwise", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "NestedBitwise", args, expected+1, t)
}

func TestNestedBitwise_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 2

	expected := testdata.NestedBitwise(1, 2)

	SymbolicMachineSatTest("numbers", "NestedBitwise", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "NestedBitwise", args, expected+1, t)
}

func TestNestedBitwise_4(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 3

	expected := testdata.NestedBitwise(1, 3)

	SymbolicMachineSatTest("numbers", "NestedBitwise", args, expected, t)
	SymbolicMachineUnsatTest("numbers", "NestedBitwise", args, expected+1, t)
}
