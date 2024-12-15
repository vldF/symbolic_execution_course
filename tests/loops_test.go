package tests

import (
	"fmt"
	"symbolic_execution_course/testdata"
	"testing"
)

func TestConstantLoop(t *testing.T) {
	args := []int{1, -1, 2, 15, 99, 1000}

	ctx := PrepareTest("loops", "ConstantLoop")
	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.ConstantLoop(a)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, -1, t)
		})
	}
}

func TestDynamicLoop(t *testing.T) {
	args := []int{1, -1, 2, 20}

	ctx := PrepareTest("loops", "DynamicLoop")
	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.DynamicLoop(a)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, -1, t)
		})
	}
}
