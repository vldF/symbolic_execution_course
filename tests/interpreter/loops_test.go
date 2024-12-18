package interpreter

import (
	"fmt"
	"symbolic_execution_course/testdata"
	"symbolic_execution_course/tests"
	"testing"
)

func TestConstantLoop(t *testing.T) {
	args := []int{1, -1, 2, 15, 99, 1000}

	ctx := tests.PrepareTest("loops", "ConstantLoop")
	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.ConstantLoop(a)

			tests.SymbolicMachineSatTest(ctx, args, expected, t)
			tests.SymbolicMachineUnsatTest(ctx, args, -1, t)
		})
	}
}

func TestDynamicLoop(t *testing.T) {
	args := []int{1, -1, 2, 20}

	ctx := tests.PrepareTest("loops", "DynamicLoop")
	for _, a := range args {
		t.Run(fmt.Sprintf("%d", a), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = a

			expected := testdata.DynamicLoop(a)

			tests.SymbolicMachineSatTest(ctx, args, expected, t)
			tests.SymbolicMachineUnsatTest(ctx, args, -1, t)
		})
	}
}
