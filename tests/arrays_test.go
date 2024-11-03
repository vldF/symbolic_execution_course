package tests

import "testing"

func TestCompareElement_1(t *testing.T) {
	args := make(map[string]any)

	arr := make([]int, 3)
	arr[0] = 1
	arr[1] = 2
	arr[2] = 3
	args["array"] = arr

	args["index"] = -1
	args["value"] = 1

	expected := -1

	SymbolicMachineTest("arrays", "compareElement", args, expected, t)
}

func TestCompareElement_2(t *testing.T) {
	args := make(map[string]any)

	arr := make([]int, 3)
	arr[0] = 1
	arr[1] = 2
	arr[2] = 3
	args["array"] = arr

	args["index"] = len(arr) + 1
	args["value"] = 1

	expected := -1

	SymbolicMachineTest("arrays", "compareElement", args, expected, t)
}

func TestCompareElement_3(t *testing.T) {
	args := make(map[string]any)

	arr := make([]int, 3)
	arr[0] = 1
	arr[1] = 2
	arr[2] = 3
	args["array"] = arr

	args["index"] = 1
	args["value"] = 0

	expected := 1

	SymbolicMachineTest("arrays", "compareElement", args, expected, t)
}

func TestCompareElement_4(t *testing.T) {
	args := make(map[string]any)

	arr := make([]int, 3)
	arr[0] = 1
	arr[1] = 2
	arr[2] = 3
	args["array"] = arr

	args["index"] = 0
	args["value"] = 2

	expected := -1

	SymbolicMachineTest("arrays", "compareElement", args, expected, t)
}

func TestCompareElement_5(t *testing.T) {
	args := make(map[string]any)

	arr := make([]int, 3)
	arr[0] = 1
	arr[1] = 2
	arr[2] = 3
	args["array"] = arr

	args["index"] = 1
	args["value"] = 2

	expected := 0

	SymbolicMachineTest("arrays", "compareElement", args, expected, t)
}
