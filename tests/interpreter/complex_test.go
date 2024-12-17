package interpreter

import (
	"fmt"
	"symbolic_execution_course/testdata"
	"symbolic_execution_course/tests"
	"testing"
)

func TestComplexReal_1(t *testing.T) {
	ctx := tests.PrepareTest("complex", "ComplexReal")

	args := make(map[string]any)
	args["a"] = complex(2.0, 3.0)

	expected := testdata.ComplexReal(complex(2.0, 3.0))

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestComplexImag_1(t *testing.T) {
	ctx := tests.PrepareTest("complex", "ComplexImag")

	args := make(map[string]any)
	args["a"] = complex(2.0, 3.0)

	expected := testdata.ComplexImag(complex(2.0, 3.0))

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestComplexId_1(t *testing.T) {
	ctx := tests.PrepareTest("complex", "ComplexId")

	args := make(map[string]any)
	args["a"] = complex(2.0, 3.0)

	expected := testdata.ComplexId(complex(2.0, 3.0))

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+complex(1.0, 1.0), t)
}

func TestBasicComplexOperations_1(t *testing.T) {
	ctx := tests.PrepareTest("complex", "BasicComplexOperations")

	args := make(map[string]any)
	a := complex(2.0, 2.0)
	b := complex(1.0, 4.0)

	args["a"] = a
	args["b"] = b

	expected := testdata.BasicComplexOperations(a, b)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+complex(1.0, 1.0), t)
}

func TestBasicComplexOperations_2(t *testing.T) {
	ctx := tests.PrepareTest("complex", "BasicComplexOperations")

	args := make(map[string]any)
	a := complex(1.0, 4.0)
	b := complex(3.0, 2.0)

	args["a"] = a
	args["b"] = b

	expected := testdata.BasicComplexOperations(a, b)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+complex(1.0, 1.0), t)
}

func TestBasicComplexOperations_3(t *testing.T) {
	ctx := tests.PrepareTest("complex", "BasicComplexOperations")

	args := make(map[string]any)
	a := complex(1.0, 2.0)
	b := complex(3.0, 4.0)

	args["a"] = a
	args["b"] = b

	expected := testdata.BasicComplexOperations(a, b)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+complex(1.0, 1.0), t)
}

func TestComplexMagnitude_1(t *testing.T) {
	ctx := tests.PrepareTest("complex", "ComplexMagnitude")

	args := make(map[string]any)
	c := complex(1.0, 2.0)
	args["a"] = c

	expected := testdata.ComplexMagnitude(c)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestComplexMagnitude_2(t *testing.T) {
	ctx := tests.PrepareTest("complex", "ComplexMagnitude")

	args := make(map[string]any)
	c := complex(0.0, 0.0)
	args["a"] = c

	expected := testdata.ComplexMagnitude(c)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestComplexMagnitude_3(t *testing.T) {
	ctx := tests.PrepareTest("complex", "ComplexMagnitude")

	args := make(map[string]any)
	c := complex(10.0, 10.0)
	args["a"] = c

	expected := testdata.ComplexMagnitude(c)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestComplexOperations_1(t *testing.T) {
	ctx := tests.PrepareTest("complex", "ComplexOperations")

	args := make(map[string]any)
	a := complex(0.0, 0.0)
	b := complex(1.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.ComplexOperations(a, b)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+complex(1, 1), t)
}

func TestComplexOperations_2(t *testing.T) {
	ctx := tests.PrepareTest("complex", "ComplexOperations")

	args := make(map[string]any)
	a := complex(1.0, 2.0)
	b := complex(0.0, 0.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.ComplexOperations(a, b)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestComplexOperations_3(t *testing.T) {
	ctx := tests.PrepareTest("complex", "ComplexOperations")

	args := make(map[string]any)
	a := complex(3.0, 2.0)
	b := complex(1.0, 4.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.ComplexOperations(a, b)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestComplexOperations_4(t *testing.T) {
	ctx := tests.PrepareTest("complex", "ComplexOperations")

	args := make(map[string]any)
	a := complex(1.0, 2.0)
	b := complex(3.0, 4.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.ComplexOperations(a, b)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestNestedComplexOperations_1(t *testing.T) {
	ctx := tests.PrepareTest("complex", "NestedComplexOperations")

	args := make(map[string]any)
	a := complex(-1.0, -1.0)
	b := complex(1.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.NestedComplexOperations(a, b)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestNestedComplexOperations_2(t *testing.T) {
	ctx := tests.PrepareTest("complex", "NestedComplexOperations")

	args := make(map[string]any)
	a := complex(-1.0, 1.0)
	b := complex(1.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.NestedComplexOperations(a, b)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestNestedComplexOperations_3(t *testing.T) {
	ctx := tests.PrepareTest("complex", "NestedComplexOperations")

	args := make(map[string]any)
	a := complex(11.0, 1.0)
	b := complex(2.0, -1.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.NestedComplexOperations(a, b)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestNestedComplexOperations_4(t *testing.T) {
	ctx := tests.PrepareTest("complex", "NestedComplexOperations")

	args := make(map[string]any)
	a := complex(11.0, 1.0)
	b := complex(10.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.NestedComplexOperations(a, b)

	tests.SymbolicMachineSatTest(ctx, args, expected, t)
	tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}

func TestComplexComparison(t *testing.T) {
	t.Skipped()
	return

	args := [][]complex128{
		{complex(1, 1), complex(1, 1)},
		{complex(1, 1), complex(0, 0)},
		{complex(10, -5), complex(3, 4)},
	}

	ctx := tests.PrepareTest("complex", "ComplexComparison")
	for i, variant := range args {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = variant[0]
			args["b"] = variant[1]

			expected := testdata.ComplexComparison(variant[0], variant[1])

			tests.SymbolicMachineSatTest(ctx, args, expected, t)
			tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
		})
	}
}
