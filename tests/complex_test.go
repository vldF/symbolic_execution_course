package tests

import (
	"symbolic_execution_course/testdata"
	"testing"
)

func TestComplexReal_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = complex(2.0, 3.0)

	expected := testdata.ComplexReal(complex(2.0, 3.0))

	SymbolicMachineSatTest("complex", "ComplexReal", args, expected, t)
	SymbolicMachineUnsatTest("complex", "ComplexReal", args, expected+1, t)
}

func TestComplexImag_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = complex(2.0, 3.0)

	expected := testdata.ComplexImag(complex(2.0, 3.0))

	SymbolicMachineSatTest("complex", "ComplexImag", args, expected, t)
	SymbolicMachineUnsatTest("complex", "ComplexImag", args, expected+1, t)
}

func TestComplexId_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = complex(2.0, 3.0)

	expected := testdata.ComplexId(complex(2.0, 3.0))

	SymbolicMachineSatTest("complex", "ComplexId", args, expected, t)
	SymbolicMachineUnsatTest("complex", "ComplexId", args, expected+complex(1.0, 1.0), t)
}

func TestBasicComplexOperations_1(t *testing.T) {
	args := make(map[string]any)
	a := complex(2.0, 2.0)
	b := complex(1.0, 4.0)

	args["a"] = a
	args["b"] = b

	expected := testdata.BasicComplexOperations(a, b)

	SymbolicMachineSatTest("complex", "BasicComplexOperations", args, expected, t)
	SymbolicMachineUnsatTest("complex", "BasicComplexOperations", args, expected+complex(1.0, 1.0), t)
}

func TestBasicComplexOperations_2(t *testing.T) {
	args := make(map[string]any)
	a := complex(1.0, 4.0)
	b := complex(3.0, 2.0)

	args["a"] = a
	args["b"] = b

	expected := testdata.BasicComplexOperations(a, b)

	SymbolicMachineSatTest("complex", "BasicComplexOperations", args, expected, t)
	SymbolicMachineUnsatTest("complex", "BasicComplexOperations", args, expected+complex(1.0, 1.0), t)
}

func TestBasicComplexOperations_3(t *testing.T) {
	args := make(map[string]any)
	a := complex(1.0, 2.0)
	b := complex(3.0, 4.0)

	args["a"] = a
	args["b"] = b

	expected := testdata.BasicComplexOperations(a, b)

	SymbolicMachineSatTest("complex", "BasicComplexOperations", args, expected, t)
	SymbolicMachineUnsatTest("complex", "BasicComplexOperations", args, expected+complex(1.0, 1.0), t)
}

func TestComplexMagnitude_1(t *testing.T) {
	args := make(map[string]any)
	c := complex(1.0, 2.0)
	args["a"] = c

	expected := testdata.ComplexMagnitude(c)

	SymbolicMachineSatTest("complex", "ComplexMagnitude", args, expected, t)
	SymbolicMachineUnsatTest("complex", "ComplexMagnitude", args, expected+1, t)
}

func TestComplexMagnitude_2(t *testing.T) {
	args := make(map[string]any)
	c := complex(0.0, 0.0)
	args["a"] = c

	expected := testdata.ComplexMagnitude(c)

	SymbolicMachineSatTest("complex", "ComplexMagnitude", args, expected, t)
	SymbolicMachineUnsatTest("complex", "ComplexMagnitude", args, expected+1, t)
}

func TestComplexMagnitude_3(t *testing.T) {
	args := make(map[string]any)
	c := complex(10.0, 10.0)
	args["a"] = c

	expected := testdata.ComplexMagnitude(c)

	SymbolicMachineSatTest("complex", "ComplexMagnitude", args, expected, t)
	SymbolicMachineUnsatTest("complex", "ComplexMagnitude", args, expected+1, t)
}

func TestComplexOperations_1(t *testing.T) {
	args := make(map[string]any)
	a := complex(0.0, 0.0)
	b := complex(1.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.ComplexOperations(a, b)

	SymbolicMachineSatTest("complex", "ComplexOperations", args, expected, t)
	SymbolicMachineUnsatTest("complex", "ComplexOperations", args, expected+complex(1, 1), t)
}

func TestComplexOperations_2(t *testing.T) {
	args := make(map[string]any)
	a := complex(1.0, 2.0)
	b := complex(0.0, 0.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.ComplexOperations(a, b)

	SymbolicMachineSatTest("complex", "ComplexOperations", args, expected, t)
	SymbolicMachineUnsatTest("complex", "ComplexOperations", args, expected+1, t)
}

func TestComplexOperations_3(t *testing.T) {
	args := make(map[string]any)
	a := complex(3.0, 2.0)
	b := complex(1.0, 4.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.ComplexOperations(a, b)

	SymbolicMachineSatTest("complex", "ComplexOperations", args, expected, t)
	SymbolicMachineUnsatTest("complex", "ComplexOperations", args, expected+1, t)
}

func TestComplexOperations_4(t *testing.T) {
	args := make(map[string]any)
	a := complex(1.0, 2.0)
	b := complex(3.0, 4.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.ComplexOperations(a, b)

	SymbolicMachineSatTest("complex", "ComplexOperations", args, expected, t)
	SymbolicMachineUnsatTest("complex", "ComplexOperations", args, expected+1, t)
}

func TestNestedComplexOperations_1(t *testing.T) {
	args := make(map[string]any)
	a := complex(-1.0, -1.0)
	b := complex(1.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.NestedComplexOperations(a, b)

	SymbolicMachineSatTest("complex", "NestedComplexOperations", args, expected, t)
	SymbolicMachineUnsatTest("complex", "NestedComplexOperations", args, expected+1, t)
}

func TestNestedComplexOperations_2(t *testing.T) {
	args := make(map[string]any)
	a := complex(-1.0, 1.0)
	b := complex(1.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.NestedComplexOperations(a, b)

	SymbolicMachineSatTest("complex", "NestedComplexOperations", args, expected, t)
	SymbolicMachineUnsatTest("complex", "NestedComplexOperations", args, expected+1, t)
}

func TestNestedComplexOperations_3(t *testing.T) {
	args := make(map[string]any)
	a := complex(11.0, 1.0)
	b := complex(2.0, -1.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.NestedComplexOperations(a, b)

	SymbolicMachineSatTest("complex", "NestedComplexOperations", args, expected, t)
	SymbolicMachineUnsatTest("complex", "NestedComplexOperations", args, expected+1, t)
}

func TestNestedComplexOperations_4(t *testing.T) {
	args := make(map[string]any)
	a := complex(11.0, 1.0)
	b := complex(10.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := testdata.NestedComplexOperations(a, b)

	SymbolicMachineSatTest("complex", "NestedComplexOperations", args, expected, t)
	SymbolicMachineUnsatTest("complex", "NestedComplexOperations", args, expected+1, t)
}
