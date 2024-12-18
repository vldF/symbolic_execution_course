package interpreter

import (
	"strconv"
	"symbolic_execution_course/tests"
	"testing"
)

func TestInvokeExternal(t *testing.T) {
	argVariants := []int{-10, 0, 5, 10}
	ctx := tests.PrepareTest("mocks", "InvokeExternal")

	for _, variant := range argVariants {
		t.Run(strconv.Itoa(variant), func(t *testing.T) {
			args := make(map[string]any)

			args["x"] = variant

			tests.SymbolicMachineSatTest(ctx, args, -1, t)
			tests.SymbolicMachineSatTest(ctx, args, 1, t)
			tests.SymbolicMachineUnsatTest(ctx, args, 2, t)
		})
	}
}
