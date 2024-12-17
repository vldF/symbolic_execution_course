package interpreter

import (
	"strconv"
	"symbolic_execution_course/testdata"
	"symbolic_execution_course/tests"
	"testing"
)

func TestCompareElement(t *testing.T) {
	argVariants := [][]int{{-1, 1}, {3, 1}, {1, 6}, {0, 2}, {1, 2}}
	ctx := tests.PrepareTest("arrays", "CompareElement")

	for _, variant := range argVariants {
		testName := "index: " + strconv.Itoa(variant[0]) + ", value: " + strconv.Itoa(variant[1])
		t.Run(testName, func(t *testing.T) {
			args := make(map[string]any)

			args["array"] = tests.ArrayArg{
				Elements: []any{1, 2, 3},
			}

			args["index"] = variant[0]
			args["value"] = variant[1]

			expected := testdata.CompareElement([]int{1, 2, 3}, variant[0], variant[1])

			tests.SymbolicMachineSatTest(ctx, args, expected, t)
			tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
		})
	}
}

func TestCompareAges(t *testing.T) {
	argVariants := [][]int{{-1, 1}, {3, 1}, {1, 6}, {0, 2}, {1, 2}}

	people := []*testdata.Person{
		{Name: "name1", Age: 1},
		{Name: "name2", Age: 5},
		{Name: "name3", Age: 10},
		{Name: "name4", Age: 15},
	}
	ctx := tests.PrepareTest("arrays", "CompareAge")

	for _, variant := range argVariants {
		testName := "index: " + strconv.Itoa(variant[0]) + ", value: " + strconv.Itoa(variant[1])
		t.Run(testName, func(t *testing.T) {
			args := make(map[string]any)

			args["index"] = variant[0]
			args["value"] = variant[1]

			expected := testdata.CompareAge(people, variant[0], variant[1])

			tests.SymbolicMachineSatTest(ctx, args, expected, t)
		})
	}
}
