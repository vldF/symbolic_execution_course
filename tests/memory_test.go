package tests

import (
	"strconv"
	"symbolic_execution_course/testdata"
	"testing"
)

func TestStructAllocateStoreRead(t *testing.T) {
	aVariants := []int{0, 1, 2, 5, 10}

	ctx := PrepareTest("memory", "StructAllocateStoreRead")
	for _, variant := range aVariants {
		t.Run(strconv.Itoa(variant), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = variant

			expected := testdata.StructAllocateStoreRead(variant)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, expected+1, t)
			SymbolicMachineUnsatTest(ctx, args, -2.0, t)
		})
	}
}

func TestStructOfStructAllocateStoreRead(t *testing.T) {
	aVariants := []int{0, 1, 2, 5, 10}

	ctx := PrepareTest("memory", "StructOfStructAllocateStoreRead")
	for _, variant := range aVariants {
		t.Run(strconv.Itoa(variant), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = variant

			expected := testdata.StructOfStructAllocateStoreRead(variant)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, expected+1, t)
			SymbolicMachineUnsatTest(ctx, args, -2.0, t)
		})
	}
}

func TestArrayAllocateStoreRead(t *testing.T) {
	aVariants := []int{0, 1, 2, 5, 10}

	ctx := PrepareTest("memory", "ArrayAllocateStoreRead")
	for _, variant := range aVariants {
		t.Run(strconv.Itoa(variant), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = variant

			expected := testdata.ArrayAllocateStoreRead(variant)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, expected+1, t)
		})
	}
}

func TestArrayAllocateStoreReadDynamic(t *testing.T) {
	argVariants := [][]int{{0, 0}, {1, 2}, {5, 10}, {10, 0}, {10, 10}}

	ctx := PrepareTest("memory", "ArrayAllocateStoreReadDynamic")
	for _, variant := range argVariants {
		t.Run(strconv.Itoa(variant[0])+"-"+strconv.Itoa(variant[1]), func(t *testing.T) {
			args := make(map[string]any)
			args["a"] = variant[0]
			args["idx"] = variant[1]

			expected := testdata.ArrayAllocateStoreReadDynamic(variant[0], variant[1])

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, expected+1, t)
		})
	}
}

func TestArrayAllocateStoreReadStore(t *testing.T) {
	args := make(map[string]any)

	ctx := PrepareTest("memory", "ArrayAllocateStoreReadStore")

	expected := testdata.ArrayAllocateStoreReadStore()

	SymbolicMachineSatTest(ctx, args, expected, t)
	SymbolicMachineUnsatTest(ctx, args, expected+1, t)
	SymbolicMachineUnsatTest(ctx, args, -1, t)
}
