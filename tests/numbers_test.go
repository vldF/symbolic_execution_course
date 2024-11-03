package tests

import (
	"testing"
)

func TestIntegerOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 1

	expected := 3

	SymbolicMachineTest("numbers", "integerOperations", args, expected, t)
}

func TestIntegerOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 2

	expected := -1

	SymbolicMachineTest("numbers", "integerOperations", args, expected, t)
}

func TestIntegerOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 2

	expected := 4

	SymbolicMachineTest("numbers", "integerOperations", args, expected, t)
}

func TestFloatOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 5.0
	args["y"] = 2.0

	expected := 2.5

	SymbolicMachineTest("numbers", "floatOperations", args, expected, t)
}

func TestFloatOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 2.0
	args["y"] = 5.0

	expected := 10.0

	SymbolicMachineTest("numbers", "floatOperations", args, expected, t)
}

func TestFloatOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 2.0
	args["y"] = 2.0

	expected := 0.0

	SymbolicMachineTest("numbers", "floatOperations", args, expected, t)
}

func TestMixedOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 3.0

	expected := 10.0

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func TestMixedOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 3.0

	expected := 0.0

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func TestMixedOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 4.0

	expected := -2.0

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func TestMixedOperations_4(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 9
	args["b"] = 4.0

	expected := 10.0

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func TestMixedOperations_5(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 7.0

	expected := 5.5

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func TestMixedOperations_6(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = -7.0

	expected := 5.0

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func TestNestedConditions_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 2.0

	expected := 3.0

	SymbolicMachineTest("numbers", "nestedConditions", args, expected, t)
}

func TestNestedConditions_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = -1
	args["b"] = 2.0

	expected := -1.0

	SymbolicMachineTest("numbers", "nestedConditions", args, expected, t)
}

func TestNestedConditions_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = -1
	args["b"] = -2.0

	expected := -1.0

	SymbolicMachineTest("numbers", "nestedConditions", args, expected, t)
}

func TestBitwiseOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 4

	expected := 2 | 4

	SymbolicMachineTest("numbers", "bitwiseOperations", args, expected, t)
}

func TestBitwiseOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 5

	expected := 3 & 5

	SymbolicMachineTest("numbers", "bitwiseOperations", args, expected, t)
}

func TestBitwiseOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 5

	expected := 4 ^ 5

	SymbolicMachineTest("numbers", "bitwiseOperations", args, expected, t)
}

func TestAdvancedBitwise_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 5
	args["b"] = 4

	expected := 5 << 1

	SymbolicMachineTest("numbers", "advancedBitwise", args, expected, t)
}

func TestAdvancedBitwise_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 5

	expected := 5 >> 1

	SymbolicMachineTest("numbers", "advancedBitwise", args, expected, t)
}

func TestAdvancedBitwise_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 4

	expected := 4 ^ 4

	SymbolicMachineTest("numbers", "advancedBitwise", args, expected, t)
}

func TestCombinedBitwise_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 1

	expected := 2 | 1

	SymbolicMachineTest("numbers", "combinedBitwise", args, expected, t)
}

func TestCombinedBitwise_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 1

	expected := 1

	SymbolicMachineTest("numbers", "combinedBitwise", args, expected, t)
}

func TestCombinedBitwise_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 0b1111
	args["b"] = 0b101111

	expected := 0b1111 ^ 0b101111

	SymbolicMachineTest("numbers", "combinedBitwise", args, expected, t)
}

func TestNestedBitwise_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = -1
	args["b"] = 0

	expected := -1

	SymbolicMachineTest("numbers", "nestedBitwise", args, expected, t)
}

func TestNestedBitwise_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 0b10101
	args["b"] = -1

	expected := 0b10101 ^ 0

	SymbolicMachineTest("numbers", "nestedBitwise", args, expected, t)
}

func TestNestedBitwise_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 2

	expected := 1 | 2

	SymbolicMachineTest("numbers", "nestedBitwise", args, expected, t)
}

func TestNestedBitwise_4(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 3

	expected := 1 & 3

	SymbolicMachineTest("numbers", "nestedBitwise", args, expected, t)
}