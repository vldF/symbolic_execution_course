package interpreter

import (
	"symbolic_execution_course/tests"
	"testing"
)

func TestSimpleAssume(t *testing.T) {
	ctx := tests.PrepareTest("intrinsics", "SimpleAssume")

	args := make(map[string]any)

	tests.SymbolicMachineSatTest(ctx, args, 2, t)
	tests.SymbolicMachineUnsatTest(ctx, args, 3, t)
}
