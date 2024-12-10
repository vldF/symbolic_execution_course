package tests

import (
	"fmt"
	"symbolic_execution_course/testdata"
	"testing"
)

func TestConstantLoop(t *testing.T) {
	args := []int{1, -1, 2, 15, 99, 1000}
	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.ConstantLoop(a)

			SymbolicMachineSatTest("loops", "ConstantLoop", args, expected, t)
			SymbolicMachineUnsatTest("loops", "ConstantLoop", args, -1, t)
		})
	}
}

func TestDynamicLoop(t *testing.T) {
	args := []int{1, -1, 2, 20}
	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.DynamicLoop(a)

			SymbolicMachineSatTest("loops", "DynamicLoop", args, expected, t)
			SymbolicMachineUnsatTest("loops", "DynamicLoop", args, -1, t)
		})
	}
}
