package interpreter

import (
	"golang.org/x/tools/go/ssa"
	"math/rand"
)

func (state *State) GetPriority(mode PathSelectorMode) int {
	switch mode {
	case Random:
		return state.getRandomPriority()
	case DFS:
		return state.getDfsPriority()
	case NURS:
		return state.getNursPriority()
	}
	panic("unsupported path selector mode")
}

func (state *State) getRandomPriority() int {
	return rand.Int()
}

func (state *State) getDfsPriority() int {
	return state.Priority + 1
}

func (state *State) getNursPriority() int {
	// priority(state) := alpha * depth(state) + beta * newBranches(state),
	// where alpha and beta are empirical coefficients,
	// depth — depth on the current state
	// newBranches — count of branches just after the state

	alpha := 2
	beta := 5

	depth := state.getDepth()
	branches := state.getBranches()

	return alpha*depth + beta*branches
}

func (state *State) getDepth() int {
	currentBlock := state.Statement.Block()
	return getDepth(currentBlock)
}

func getDepth(bb *ssa.BasicBlock) int {
	maxPrevDepth := 0
	for _, pred := range bb.Preds {
		newDepth := getDepth(pred)
		if newDepth > maxPrevDepth {
			maxPrevDepth = newDepth
		}
	}

	return 1 + maxPrevDepth
}

func (state *State) getBranches() int {
	statement := state.Statement
	switch statement.(type) {
	case *ssa.If:
		return 2
	case *ssa.Return:
		return 0
	default:
		return 1
	}
}
