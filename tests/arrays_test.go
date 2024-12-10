package tests

import (
	"strconv"
	"symbolic_execution_course/testdata"
	"testing"
)

func TestCompareElement(t *testing.T) {
	argVariants := [][]int{{-1, 1}, {3, 1}, {1, 6}, {0, 2}, {1, 2}}

	for _, variant := range argVariants {
		testName := "index: " + strconv.Itoa(variant[0]) + ", value: " + strconv.Itoa(variant[1])
		t.Run(testName, func(t *testing.T) {
			args := make(map[string]any)

			args["array"] = ArrayArg{
				elements:    []any{1, 2, 3},
				elementType: "int",
			}

			args["index"] = variant[0]
			args["value"] = variant[1]

			expected := testdata.CompareElement([]int{1, 2, 3}, variant[0], variant[1])

			SymbolicMachineSatTest("arrays", "CompareElement", args, expected, t)
			SymbolicMachineUnsatTest("arrays", "CompareElement", args, expected+1, t)
		})
	}
}

func TestCompareAge_1(t *testing.T) {
	args := make(map[string]any)

	args["index"] = 1
	args["value"] = 10

	expected := 1

	SymbolicMachineSatTest("arrays", "compareAge", args, expected, t)
	SymbolicMachineUnsatTest("arrays", "compareAge", args, expected+1, t)
}

func TestCompareAge_2(t *testing.T) {
	args := make(map[string]any)

	args["index"] = 1
	args["value"] = 10

	expected := -1

	SymbolicMachineSatTest("arrays", "compareAge", args, expected, t)
}

func TestCompareAge_3(t *testing.T) {
	args := make(map[string]any)

	args["index"] = 1
	args["value"] = 10

	expected := -1

	SymbolicMachineSatTest("arrays", "compareAge", args, expected, t)
}
