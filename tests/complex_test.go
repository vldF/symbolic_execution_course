package tests

import "testing"

func TestBasicComplexOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = complex(2.0, 2.0)
	args["b"] = complex(1.0, 4.0)

	expected := complex(3, 6)

	SymbolicMachineSatTest("complex", "basicComplexOperations", args, expected, t)
}

func TestBasicComplexOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = complex(1.0, 4.0)
	args["b"] = complex(3.0, 2.0)

	expected := complex(-2.0, 2.0)

	SymbolicMachineSatTest("complex", "basicComplexOperations", args, expected, t)
}

func TestBasicComplexOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = complex(1.0, 2.0)
	args["b"] = complex(3.0, 4.0)

	expected := complex(-5, 10)

	SymbolicMachineSatTest("complex", "basicComplexOperations", args, expected, t)
}

func TestComplexMagnitude_1(t *testing.T) {
	args := make(map[string]any)
	c := complex(1.0, 2.0)
	args["a"] = c

	expected := real(c)*real(c) + imag(c)*imag(c)

	SymbolicMachineSatTest("complex", "complexMagnitude", args, expected, t)
}

func TestComplexMagnitude_2(t *testing.T) {
	args := make(map[string]any)
	c := complex(0.0, 0.0)
	args["a"] = c

	expected := real(c)*real(c) + imag(c)*imag(c)

	SymbolicMachineSatTest("complex", "complexMagnitude", args, expected, t)
}

func TestComplexMagnitude_3(t *testing.T) {
	args := make(map[string]any)
	c := complex(10.0, 10.0)
	args["a"] = c

	expected := real(c)*real(c) + imag(c)*imag(c)

	SymbolicMachineSatTest("complex", "complexMagnitude", args, expected, t)
}

func TestComplexOperations_1(t *testing.T) {
	args := make(map[string]any)
	a := complex(0.0, 0.0)
	b := complex(1.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := b

	SymbolicMachineSatTest("complex", "complexOperations", args, expected, t)
}

func TestComplexOperations_2(t *testing.T) {
	args := make(map[string]any)
	a := complex(1.0, 2.0)
	b := complex(0.0, 0.0)
	args["a"] = a
	args["b"] = b

	expected := a

	SymbolicMachineSatTest("complex", "complexOperations", args, expected, t)
}

func TestComplexOperations_3(t *testing.T) {
	args := make(map[string]any)
	a := complex(3.0, 2.0)
	b := complex(1.0, 4.0)
	args["a"] = a
	args["b"] = b

	expected := a / b

	SymbolicMachineSatTest("complex", "complexOperations", args, expected, t)
}

func TestComplexOperations_4(t *testing.T) {
	args := make(map[string]any)
	a := complex(1.0, 2.0)
	b := complex(3.0, 4.0)
	args["a"] = a
	args["b"] = b

	expected := a + b

	SymbolicMachineSatTest("complex", "complexOperations", args, expected, t)
}

func TestNestedComplexOperations_1(t *testing.T) {
	args := make(map[string]any)
	a := complex(-1.0, -1.0)
	b := complex(1.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := a * b

	SymbolicMachineSatTest("complex", "nestedComplexOperations", args, expected, t)
}

func TestNestedComplexOperations_2(t *testing.T) {
	args := make(map[string]any)
	a := complex(-1.0, 1.0)
	b := complex(1.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := a + b

	SymbolicMachineSatTest("complex", "nestedComplexOperations", args, expected, t)
}

func TestNestedComplexOperations_3(t *testing.T) {
	args := make(map[string]any)
	a := complex(11.0, 1.0)
	b := complex(2.0, -1.0)
	args["a"] = a
	args["b"] = b

	expected := a - b

	SymbolicMachineSatTest("complex", "nestedComplexOperations", args, expected, t)
}

func TestNestedComplexOperations_4(t *testing.T) {
	args := make(map[string]any)
	a := complex(11.0, 1.0)
	b := complex(10.0, 2.0)
	args["a"] = a
	args["b"] = b

	expected := a + b

	SymbolicMachineSatTest("complex", "nestedComplexOperations", args, expected, t)
}
