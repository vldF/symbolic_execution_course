package tests

import (
	"testing"
)

func TestIdInt_1(t *testing.T) {
	args := make(map[string]any)

	args["x"] = 1
	expected := 1

	SymbolicMachineTest("simple", "idInt", args, expected, t)
}
