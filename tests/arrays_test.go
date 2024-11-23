package tests

import (
	"testing"
)

func TestCompareElement_1(t *testing.T) {
	args := make(map[string]any)

	args["array"] = ArrayArg{
		elements:        []any{1, 2, 3},
		elementTypeName: "int",
	}

	args["index"] = -1
	args["value"] = 1

	expected := -1

	SymbolicMachineSatTest("arrays", "compareElement", args, expected, t)
	SymbolicMachineUnsatTest("arrays", "compareElement", args, expected+1, t)
}

func TestCompareElement_2(t *testing.T) {
	args := make(map[string]any)

	args["array"] = ArrayArg{
		elements:        []any{1, 2, 3},
		elementTypeName: "int",
	}

	args["index"] = 3
	args["value"] = 1

	expected := -1

	SymbolicMachineSatTest("arrays", "compareElement", args, expected, t)
	SymbolicMachineUnsatTest("arrays", "compareElement", args, expected+1, t)
}

func TestCompareElement_3(t *testing.T) {
	args := make(map[string]any)

	args["array"] = ArrayArg{
		elements:        []any{1, 2, 3},
		elementTypeName: "int",
	}

	args["index"] = 1
	args["value"] = 0

	expected := 1

	SymbolicMachineSatTest("arrays", "compareElement", args, expected, t)
	SymbolicMachineUnsatTest("arrays", "compareElement", args, expected+1, t)
}

func TestCompareElement_4(t *testing.T) {
	args := make(map[string]any)

	args["array"] = ArrayArg{
		elements:        []any{1, 2, 3},
		elementTypeName: "int",
	}

	args["index"] = 0
	args["value"] = 2

	expected := -1

	SymbolicMachineSatTest("arrays", "compareElement", args, expected, t)
	SymbolicMachineUnsatTest("arrays", "compareElement", args, expected+1, t)
}

func TestCompareElement_5(t *testing.T) {
	args := make(map[string]any)

	args["array"] = ArrayArg{
		elements:        []any{1, 2, 3},
		elementTypeName: "int",
	}

	args["index"] = 1
	args["value"] = 2

	expected := 0

	SymbolicMachineSatTest("arrays", "compareElement", args, expected, t)
	SymbolicMachineUnsatTest("arrays", "compareElement", args, expected+1, t)
}

//func TestCompareAge_1(t *testing.T) {
//	args := make(map[string]any)
//
//	args["index"] = 1
//	args["value"] = 10
//
//	expected := 1
//
//	SymbolicMachineSatTest("arrays", "compareAge", args, expected, t)
//}
//
//func TestCompareAge_2(t *testing.T) {
//	args := make(map[string]any)
//
//	args["index"] = 1
//	args["value"] = 10
//
//	expected := -1
//
//	SymbolicMachineSatTest("arrays", "compareAge", args, expected, t)
//}
//
//func TestCompareAge_3(t *testing.T) {
//	args := make(map[string]any)
//
//	args["index"] = 1
//	args["value"] = 10
//
//	expected := -1
//
//	SymbolicMachineSatTest("arrays", "compareAge", args, expected, t)
//}
