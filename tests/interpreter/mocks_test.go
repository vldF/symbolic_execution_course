package interpreter

import (
	"symbolic_execution_course/tests"
	"testing"
)

func TestInvokeExternal(t *testing.T) {
	ctx := tests.PrepareTest("mocks", "InvokeExternal")

	args := make(map[string]any)

	tests.SymbolicMachineSatTest(ctx, args, 1, t)
	tests.SymbolicMachineSatTest(ctx, args, 2, t)
	tests.SymbolicMachineSatTest(ctx, args, 3, t)
	tests.SymbolicMachineSatTest(ctx, args, -1, t)
	tests.SymbolicMachineUnsatTest(ctx, args, -2, t)
}
