package interpreter

import "github.com/aclements/go-z3/z3"

func hasSolution(state *State, context *Context) bool {
	if len(state.Constraints) == 0 {
		return true
	}

	stateRes := state.Constraints[0].AsBool().AsZ3Value().Value.(z3.Bool)
	for _, constraint := range state.Constraints[1:] {
		asBool := constraint.AsBool().AsZ3Value().Value.(z3.Bool)
		stateRes = stateRes.And(asBool)
	}

	solver := context.Solver
	solver.Reset()
	solver.Assert(stateRes)

	//println("Solver constraints:", solver.String())

	check, err := solver.Check()
	if err != nil {
		println("error!", err.Error())
		return false
	}

	//println("is sat", check)
	//println()

	return check
}
