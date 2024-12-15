package tests

import (
	"fmt"
	"symbolic_execution_course/testdata"
	"testing"
)

func TestImpossibleBranch(t *testing.T) {
	args := []int{-5, -4, 0, 2, 5}
	ctx := PrepareTest("dynamic_interpretation", "ImpossibleBranch")
	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.ImpossibleBranch(a)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, -1, t)
		})
	}
}

func TestRepeatingConditions(t *testing.T) {
	args := []int{-5, -4, 0, 2, 5}
	ctx := PrepareTest("dynamic_interpretation", "RepeatingConditions")

	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.RepeatingConditions(a)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, -1, t)
		})
	}
}

func TestImpossibleCondition(t *testing.T) {
	args := []int{-5, -4, 0, 2, 5}
	ctx := PrepareTest("dynamic_interpretation", "ImpossibleCondition")

	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.ImpossibleCondition(a)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, -1, t)
		})
	}
}
