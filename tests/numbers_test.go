package tests

import (
	"testing"
)

func Test_Numbers_IntegerOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 1

	expected := 3

	SymbolicMachineTest("numbers", "integerOperations", args, expected, t)
}

func Test_Numbers_IntegerOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 2

	expected := -1

	SymbolicMachineTest("numbers", "integerOperations", args, expected, t)
}

func Test_Numbers_IntegerOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 2

	expected := 4

	SymbolicMachineTest("numbers", "integerOperations", args, expected, t)
}

func Test_Numbers_FloatOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 5.0
	args["y"] = 2.0

	expected := 2.5

	SymbolicMachineTest("numbers", "floatOperations", args, expected, t)
}

func Test_Numbers_FloatOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 2.0
	args["y"] = 5.0

	expected := 10.0

	SymbolicMachineTest("numbers", "floatOperations", args, expected, t)
}

func Test_Numbers_FloatOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 2.0
	args["y"] = 2.0

	expected := 0.0

	SymbolicMachineTest("numbers", "floatOperations", args, expected, t)
}

func Test_Numbers_MixedOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 3.0

	expected := 10.0

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func Test_Numbers_MixedOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 3.0

	expected := 0.0

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func Test_Numbers_MixedOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 4.0

	expected := -2.0

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func Test_Numbers_MixedOperations_4(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 9
	args["b"] = 4.0

	expected := 10.0

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func Test_Numbers_MixedOperations_5(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 7.0

	expected := 5.5

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func Test_Numbers_MixedOperations_6(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = -7.0

	expected := 5.0

	SymbolicMachineTest("numbers", "mixedOperations", args, expected, t)
}

func Test_Numbers_NestedConditions_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 2.0

	expected := 3.0

	SymbolicMachineTest("numbers", "nestedConditions", args, expected, t)
}

func Test_Numbers_NestedConditions_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = -1
	args["b"] = 2.0

	expected := -1.0

	SymbolicMachineTest("numbers", "nestedConditions", args, expected, t)
}

func Test_Numbers_NestedConditions_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = -1
	args["b"] = -2.0

	expected := -1.0

	SymbolicMachineTest("numbers", "nestedConditions", args, expected, t)
}

func Test_Numbers_BitwiseOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 4

	expected := 2 | 4

	SymbolicMachineTest("numbers", "bitwiseOperations", args, expected, t)
}

func Test_Numbers_BitwiseOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 5

	expected := 3 & 5

	SymbolicMachineTest("numbers", "bitwiseOperations", args, expected, t)
}

func Test_Numbers_BitwiseOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 5

	expected := 4 ^ 5

	SymbolicMachineTest("numbers", "bitwiseOperations", args, expected, t)
}
