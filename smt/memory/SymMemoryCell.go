package memory

import "github.com/aclements/go-z3/z3"

const (
	Array = iota
)

type SymMemoryCell struct {
	Kind   int
	Fields map[string]z3.Value
}
