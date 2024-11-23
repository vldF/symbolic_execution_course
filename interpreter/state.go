package interpreter

import (
	"golang.org/x/tools/go/ssa"
)

type State struct {
	Constraints        []BoolValue
	Stack              map[string]Value
	Statement          ssa.Instruction
	VisitedBasicBlocks []int
}

func (state *State) Copy() *State {
	constraints := make([]BoolValue, len(state.Constraints))
	copy(constraints, state.Constraints)

	memory := make(map[string]Value)
	for k, v := range state.Stack {
		memory[k] = v
	}

	blocks := make([]int, len(state.VisitedBasicBlocks))
	copy(blocks, state.VisitedBasicBlocks)

	return &State{
		Constraints:        constraints,
		Stack:              memory,
		Statement:          state.Statement,
		VisitedBasicBlocks: blocks,
	}
}
