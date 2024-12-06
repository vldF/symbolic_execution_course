package interpreter

import (
	"golang.org/x/tools/go/ssa"
	"slices"
)

type State struct {
	Priority           int
	Constraints        []BoolValue
	StackFrames        []*StackFrame
	Statement          ssa.Instruction
	VisitedBasicBlocks []int
}

// StackFrame todo: save here callstack too
type StackFrame struct {
	Values map[string]Value
}

func (frame *StackFrame) copy() *StackFrame {
	newValues := make(map[string]Value, len(frame.Values))
	for k, v := range frame.Values {
		newValues[k] = v
	}

	return &StackFrame{
		Values: newValues,
	}
}

func (state *State) Copy() *State {
	constraints := make([]BoolValue, len(state.Constraints))
	copy(constraints, state.Constraints)

	newFrames := make([]*StackFrame, 0)
	for _, frame := range state.StackFrames {
		newFrames = append(newFrames, frame.copy())
	}

	blocks := make([]int, len(state.VisitedBasicBlocks))
	copy(blocks, state.VisitedBasicBlocks)

	return &State{
		Constraints:        constraints,
		StackFrames:        newFrames,
		Statement:          state.Statement,
		VisitedBasicBlocks: blocks,
	}
}

func (state *State) LastStackFrame() *StackFrame {
	return state.StackFrames[len(state.StackFrames)-1]
}

func (state *State) GetValueFromStack(name string) Value {
	for _, frame := range slices.Backward(state.StackFrames) {
		if res, ok := frame.Values[name]; ok {
			return res
		}
	}

	return nil
}

func (state *State) PushStackFrame() {
	newFrame := &StackFrame{
		Values: make(map[string]Value),
	}
	state.StackFrames = append(state.StackFrames, newFrame)
}

func (state *State) PopStackFrame() {
	state.StackFrames = state.StackFrames[:len(state.StackFrames)-2] // todo?
}
