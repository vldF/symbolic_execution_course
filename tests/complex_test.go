package tests

import "testing"

func TestBasicComplexOperations_1(t *testing.T) {
	args := make(map[string]any)
	args["a"] = complex(2.0, 2.0)
	args["b"] = complex(1.0, 4.0)

	expected := complex(3, 6)

	SymbolicMachineTest("complex", "basicComplexOperations", args, expected, t)
}

func TestBasicComplexOperations_2(t *testing.T) {
	args := make(map[string]any)
	args["a"] = complex(1.0, 4.0)
	args["b"] = complex(3.0, 2.0)

	expected := complex(-2.0, 2.0)

	SymbolicMachineTest("complex", "basicComplexOperations", args, expected, t)
}

func TestBasicComplexOperations_3(t *testing.T) {
	args := make(map[string]any)
	args["a"] = complex(1.0, 2.0)
	args["b"] = complex(3.0, 4.0)

	expected := complex(-5, 10)

	SymbolicMachineTest("complex", "basicComplexOperations", args, expected, t)
}
