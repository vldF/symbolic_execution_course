package memory

import "github.com/aclements/go-z3/z3"

type Memory struct {
	Cells map[*z3.Uninterpreted]*SymMemoryCell
}
