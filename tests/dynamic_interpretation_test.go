package tests

import (
	"fmt"
	"symbolic_execution_course/testdata"
	"testing"
)

func TestImpossibleBranch(t *testing.T) {
	args := []int{-5, -4, 0, 2, 5}
	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.ImpossibleBranch(a)

			SymbolicMachineSatTest("dynamic_interpretation", "ImpossibleBranch", args, expected, t)
			SymbolicMachineUnsatTest("dynamic_interpretation", "ImpossibleBranch", args, -1, t)
		})
	}
}

func TestRepeatingConditions(t *testing.T) {
	args := []int{-5, -4, 0, 2, 5}
	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.RepeatingConditions(a)

			SymbolicMachineSatTest("dynamic_interpretation", "RepeatingConditions", args, expected, t)
			SymbolicMachineUnsatTest("dynamic_interpretation", "RepeatingConditions", args, -1, t)
		})
	}
}

func TestImpossibleCondition(t *testing.T) {
	args := []int{-5, -4, 0, 2, 5}
	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.ImpossibleCondition(a)

			SymbolicMachineSatTest("dynamic_interpretation", "ImpossibleCondition", args, expected, t)
			SymbolicMachineUnsatTest("dynamic_interpretation", "ImpossibleCondition", args, -1, t)
		})
	}
}
