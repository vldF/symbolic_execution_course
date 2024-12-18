package interpreter

import (
	"strconv"
	"symbolic_execution_course/testdata"
	"symbolic_execution_course/tests"
	"testing"
)

func TestDistance(t *testing.T) {
	argVariants := [][]float64{{0, 0, 1, 1}, {1, 2, 3, 4}, {2, 3, 4, 5}}

	ctx := tests.PrepareTest("functions_overriding", "Distance")
	for i, variant := range argVariants {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			args := make(map[string]any)

			args["x1"] = variant[0]
			args["y1"] = variant[1]
			args["x2"] = variant[2]
			args["y2"] = variant[3]

			expected := testdata.Distance(variant[0], variant[1], variant[2], variant[3])

			tests.SymbolicMachineUnsatTest(ctx, args, expected+1, t)
		})
	}
}
