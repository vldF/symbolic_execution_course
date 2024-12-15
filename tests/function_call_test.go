package tests

import (
	"strconv"
	"symbolic_execution_course/testdata"
	"testing"
)

func TestTwice(t *testing.T) {
	argVariants := []int{1, 2, -1, -2, 10}

	ctx := PrepareTest("function_call", "Twice")
	for _, variant := range argVariants {
		t.Run(strconv.Itoa(variant), func(t *testing.T) {
			args := make(map[string]any)

			args["a"] = variant

			expected := testdata.Twice(variant)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, expected+1, t)
		})
	}
}

func TestTwiceComplex(t *testing.T) {
	argVariants := []complex128{complex(1, 1), complex(2, 2), complex(3, 3)}

	ctx := PrepareTest("function_call", "TwiceComplex")
	for i, variant := range argVariants {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			args := make(map[string]any)

			args["a"] = variant

			expected := testdata.TwiceComplex(variant)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, expected+1, t)
		})
	}
}

func TestTwiceStruct(t *testing.T) {
	argVariants := []int{1, 2, 5, 10}

	ctx := PrepareTest("function_call", "TwiceStruct")
	for i, variant := range argVariants {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			args := make(map[string]any)

			args["a"] = variant

			expected := testdata.TwiceStruct(variant)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, expected+1, t)
		})
	}
}

func TestAddRecursive(t *testing.T) {
	argVariants := [][]int{{1, 1}, {3, 3}, {10, 10}, {20, 20}}

	ctx := PrepareTest("function_call", "AddRecursive")
	for i, variant := range argVariants {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			args := make(map[string]any)

			args["a"] = variant[0]
			args["b"] = variant[1]

			expected := testdata.AddRecursive(variant[0], variant[1])

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, expected+1, t)
		})
	}
}

func TestFib(t *testing.T) {
	args := make(map[string]any)
	ctx := PrepareTest("function_call", "Fib")

	args["n"] = 5

	expected := testdata.Fib(5)

	SymbolicMachineSatTest(ctx, args, expected, t)
	SymbolicMachineUnsatTest(ctx, args, expected+1, t)
}
