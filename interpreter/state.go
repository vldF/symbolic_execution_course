package interpreter

import (
	"golang.org/x/tools/go/ssa"
)

type State struct {
	Constraints []BoolPredicate
	Memory      map[string]Value
	Statement   ssa.Instruction
}

func (state *State) Copy() *State {
	constraints := make([]BoolPredicate, len(state.Constraints))
	copy(constraints, state.Constraints)

	memory := make(map[string]Value)
	for k, v := range state.Memory {
		memory[k] = v
	}

	return &State{
		Constraints: constraints,
		Memory:      memory,
		Statement:   state.Statement,
	}
}
