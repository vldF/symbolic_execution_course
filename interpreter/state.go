package interpreter

import (
	"golang.org/x/tools/go/ssa"
)

type State struct {
	Priority           int
	Constraints        []BoolValue
	StackFrames        []*StackFrame
	Statement          ssa.Instruction
	VisitedBasicBlocks []int
}

type StackFrame struct {
	Initiator *ssa.Call
	Values    map[string]Value
}

func (frame *StackFrame) copy() *StackFrame {
	newValues := make(map[string]Value, len(frame.Values))
	for k, v := range frame.Values {
		newValues[k] = v
	}

	return &StackFrame{
		Initiator: frame.Initiator,
		Values:    newValues,
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
	if len(state.StackFrames) == 0 {
		return nil
	}

	return state.StackFrames[len(state.StackFrames)-1]
}

func (state *State) GetValueFromStack(name string) Value {
	return state.LastStackFrame().Values[name]
}

func (state *State) PushStackFrame(initiator *ssa.Call) {
	newFrame := &StackFrame{
		Values:    make(map[string]Value),
		Initiator: initiator,
	}

	state.StackFrames = append(state.StackFrames, newFrame)
}

func (state *State) PopStackFrame() {
	state.StackFrames = state.StackFrames[:len(state.StackFrames)-1]
}
