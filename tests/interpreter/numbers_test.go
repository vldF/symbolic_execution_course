package interpreter

import (
	"symbolic_execution_course/testdata"
	"symbolic_execution_course/tests"
	"testing"
)

func TestIntegerOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 1

	ctx := tests.PrepareTest("numbers", "IntegerOperations")

	expected := testdata.IntegerOperations(2, 1)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestIntegerOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 2

	ctx := tests.PrepareTest("numbers", "IntegerOperations")

	expected := testdata.IntegerOperations(1, 2)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestIntegerOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 2

	ctx := tests.PrepareTest("numbers", "IntegerOperations")

	expected := testdata.IntegerOperations(2, 2)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestFloatOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 5.0
	args["y"] = 2.0

	ctx := tests.PrepareTest("numbers", "FloatOperations")

	expected := testdata.FloatOperations(5, 2)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestFloatOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 2.0
	args["y"] = 5.0

	ctx := tests.PrepareTest("numbers", "FloatOperations")

	expected := testdata.FloatOperations(2, 5)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestFloatOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["x"] = 2.0
	args["y"] = 2.0

	ctx := tests.PrepareTest("numbers", "FloatOperations")

	expected := testdata.FloatOperations(2, 2)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestMixedOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 3.0

	ctx := tests.PrepareTest("numbers", "MixedOperations")

	expected := testdata.MixedOperations(2, 3)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestMixedOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 3.0

	ctx := tests.PrepareTest("numbers", "MixedOperations")

	expected := testdata.MixedOperations(3, 3)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestMixedOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 4.0

	ctx := tests.PrepareTest("numbers", "MixedOperations")

	expected := testdata.MixedOperations(3, 4)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestMixedOperations_4(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 9
	args["b"] = 4.0

	ctx := tests.PrepareTest("numbers", "MixedOperations")

	expected := testdata.MixedOperations(9, 4)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestMixedOperations_5(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 7.0

	ctx := tests.PrepareTest("numbers", "MixedOperations")

	expected := testdata.MixedOperations(4, 7)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestMixedOperations_6(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = -7.0

	ctx := tests.PrepareTest("numbers", "MixedOperations")

	expected := testdata.MixedOperations(3, -7)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestNestedConditions_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 2.0

	ctx := tests.PrepareTest("numbers", "NestedConditions")

	expected := testdata.NestedConditions(1, 2)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestNestedConditions_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = -1
	args["b"] = 2.0

	ctx := tests.PrepareTest("numbers", "NestedConditions")

	expected := testdata.NestedConditions(-1, 2)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestNestedConditions_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = -1
	args["b"] = -2.0

	ctx := tests.PrepareTest("numbers", "NestedConditions")

	expected := testdata.NestedConditions(-1, -2)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestBitwiseOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 4

	ctx := tests.PrepareTest("numbers", "BitwiseOperations")

	expected := testdata.BitwiseOperations(2, 4)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestBitwiseOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 5

	ctx := tests.PrepareTest("numbers", "BitwiseOperations")

	expected := testdata.BitwiseOperations(3, 5)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestBitwiseOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 5

	ctx := tests.PrepareTest("numbers", "BitwiseOperations")

	expected := testdata.BitwiseOperations(4, 5)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestAdvancedBitwise_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 5
	args["b"] = 4

	ctx := tests.PrepareTest("numbers", "AdvancedBitwise")

	expected := testdata.AdvancedBitwise(5, 4)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestAdvancedBitwise_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 5

	ctx := tests.PrepareTest("numbers", "AdvancedBitwise")

	expected := testdata.AdvancedBitwise(4, 5)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestAdvancedBitwise_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 4
	args["b"] = 4

	ctx := tests.PrepareTest("numbers", "AdvancedBitwise")

	expected := testdata.AdvancedBitwise(4, 4)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestCombinedBitwise_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 2
	args["b"] = 1

	ctx := tests.PrepareTest("numbers", "CombinedBitwise")

	expected := testdata.CombinedBitwise(2, 1)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestCombinedBitwise_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 3
	args["b"] = 1

	ctx := tests.PrepareTest("numbers", "CombinedBitwise")

	expected := testdata.CombinedBitwise(3, 1)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestCombinedBitwise_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 0b1111
	args["b"] = 0b101111

	ctx := tests.PrepareTest("numbers", "CombinedBitwise")

	expected := testdata.CombinedBitwise(0b1111, 0b101111)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestNestedBitwise_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = -1
	args["b"] = 0

	ctx := tests.PrepareTest("numbers", "NestedBitwise")

	expected := testdata.NestedBitwise(-1, 0)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestNestedBitwise_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 0b10101
	args["b"] = -1

	ctx := tests.PrepareTest("numbers", "NestedBitwise")

	expected := testdata.NestedBitwise(0b10101, -1)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestNestedBitwise_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 2

	ctx := tests.PrepareTest("numbers", "NestedBitwise")

	expected := testdata.NestedBitwise(1, 2)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestNestedBitwise_4(t *testing.T) {
	args := make(map[string]any)
	args["a"] = 1
	args["b"] = 3

	ctx := tests.PrepareTest("numbers", "NestedBitwise")

	expected := testdata.NestedBitwise(1, 3)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}
