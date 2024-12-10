package tests

import (
	"symbolic_execution_course/testdata"
	"testing"
)

func TestStructAllocateStoreRead(t *testing.T) {
	aVariants := []int{0, 1, 2, 5, 10}

	for _, variant := range aVariants {
		args := make(map[string]any)
		args["a"] = variant

		expected := testdata.StructAllocateStoreRead(variant)

		SymbolicMachineSatTest("memory", "StructAllocateStoreRead", args, expected, t)
		SymbolicMachineUnsatTest("memory", "StructAllocateStoreRead", args, expected+1, t)
		SymbolicMachineUnsatTest("memory", "StructAllocateStoreRead", args, -2.0, t)
	}
}

func TestArrayAllocateStoreRead(t *testing.T) {
	aVariants := []int{0, 1, 2, 5, 10}

	for _, variant := range aVariants {
		args := make(map[string]any)
		args["a"] = variant

		expected := testdata.ArrayAllocateStoreRead(variant)

		SymbolicMachineSatTest("memory", "ArrayAllocateStoreRead", args, expected, t)
		SymbolicMachineUnsatTest("memory", "ArrayAllocateStoreRead", args, expected+1, t)
	}
}
