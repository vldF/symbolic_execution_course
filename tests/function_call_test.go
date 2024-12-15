package tests

import (
	"strconv"
	"symbolic_execution_course/testdata"
	"testing"
)

func TestTwice(t *testing.T) {
	argVariants := []int{1, 2, -1, -2, 10}

	for _, variant := range argVariants {
		t.Run(strconv.Itoa(variant), func(t *testing.T) {
			args := make(map[string]any)

			args["a"] = variant

			expected := testdata.Twice(variant)

			SymbolicMachineSatTest("function_call", "Twice", args, expected, t)
			SymbolicMachineUnsatTest("function_call", "Twice", args, expected+1, t)
		})
	}
}

func TestTwiceComplex(t *testing.T) {
	argVariants := []complex128{complex(1, 1), complex(2, 2), complex(3, 3)}

	for i, variant := range argVariants {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			args := make(map[string]any)

			args["a"] = variant

			expected := testdata.TwiceComplex(variant)

			SymbolicMachineSatTest("function_call", "TwiceComplex", args, expected, t)
			SymbolicMachineUnsatTest("function_call", "TwiceComplex", args, expected+1, t)
		})
	}
}

func TestTwiceStruct(t *testing.T) {
	argVariants := []int{1, 2, 5, 10}

	for i, variant := range argVariants {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			args := make(map[string]any)

			args["a"] = variant

			expected := testdata.TwiceStruct(variant)

			SymbolicMachineSatTest("function_call", "TwiceStruct", args, expected, t)
			SymbolicMachineUnsatTest("function_call", "TwiceStruct", args, expected+1, t)
		})
	}
}

func TestAddRecursive(t *testing.T) {
	argVariants := [][]int{{1, 1}, {3, 3}, {10, 10}, {20, 20}}

	for i, variant := range argVariants {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			args := make(map[string]any)

			args["a"] = variant[0]
			args["b"] = variant[1]

			expected := testdata.AddRecursive(variant[0], variant[1])

			SymbolicMachineSatTest("function_call", "AddRecursive", args, expected, t)
			SymbolicMachineUnsatTest("function_call", "AddRecursive", args, expected+1, t)
		})
	}
}

// this test is really slow because the interpreter build all
// states with depth ~100 regardless of the function argument

func TestFib(t *testing.T) {
	args := make(map[string]any)

	args["n"] = 5

	expected := testdata.Fib(5)

	SymbolicMachineSatTest("function_call", "Fib", args, expected, t)
	SymbolicMachineUnsatTest("function_call", "Fib", args, expected+1, t)
}
