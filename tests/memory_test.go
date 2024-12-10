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

func TestStructOfStructAllocateStoreRead(t *testing.T) {
	aVariants := []int{0, 1, 2, 5, 10}

	for _, variant := range aVariants {
		args := make(map[string]any)
		args["a"] = variant

		expected := testdata.StructOfStructAllocateStoreRead(variant)

		SymbolicMachineSatTest("memory", "StructOfStructAllocateStoreRead", args, expected, t)
		SymbolicMachineUnsatTest("memory", "StructOfStructAllocateStoreRead", args, expected+1, t)
		SymbolicMachineUnsatTest("memory", "StructOfStructAllocateStoreRead", args, -2.0, t)
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

func TestArrayAllocateStoreReadDynamic(t *testing.T) {
	argVariants := [][]int{{0, 0}, {1, 2}, {5, 10}, {10, 0}, {10, 10}}

	for _, variant := range argVariants {
		args := make(map[string]any)
		args["a"] = variant[0]
		args["idx"] = variant[1]

		expected := testdata.ArrayAllocateStoreReadDynamic(variant[0], variant[1])

		SymbolicMachineSatTest("memory", "ArrayAllocateStoreReadDynamic", args, expected, t)
		SymbolicMachineUnsatTest("memory", "ArrayAllocateStoreReadDynamic", args, expected+1, t)
	}
}

func TestArrayAllocateStoreReadStore(t *testing.T) {
	args := make(map[string]any)

	expected := testdata.ArrayAllocateStoreReadStore()

	SymbolicMachineSatTest("memory", "ArrayAllocateStoreReadStore", args, expected, t)
	SymbolicMachineUnsatTest("memory", "ArrayAllocateStoreReadStore", args, expected+1, t)
	SymbolicMachineUnsatTest("memory", "ArrayAllocateStoreReadStore", args, -1, t)
}
