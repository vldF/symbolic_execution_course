package tests

import (
	"strconv"
	"symbolic_execution_course/testdata"
	"testing"
)

func TestShortOverflow(t *testing.T) {
	aVariants := [][]int{{0, 0}, {127, 1}, {127, 5}, {64, 64}, {127, 127}}

	ctx := PrepareTest("overflow", "ShortOverflow")
	for _, variant := range aVariants {
		t.Run(strconv.Itoa(variant[0])+"-"+strconv.Itoa(variant[1]), func(t *testing.T) {
			args := make(map[string]any)
			args["x"] = variant[0]
			args["y"] = variant[1]

			expected := testdata.ShortOverflow(variant[0], variant[1])

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, expected+1, t)
		})
	}
}

func TestOverflowInLoop(t *testing.T) {
	aVariants := []int{0, 10, 20, 30, 40}

	ctx := PrepareTest("overflow", "OverflowInLoop")
	for _, variant := range aVariants {
		t.Run(strconv.Itoa(variant), func(t *testing.T) {
			args := make(map[string]any)
			args["x"] = variant

			expected := testdata.OverflowInLoop(variant)

			SymbolicMachineSatTest(ctx, args, expected, t)
			SymbolicMachineUnsatTest(ctx, args, expected+1, t)
		})
	}
}
