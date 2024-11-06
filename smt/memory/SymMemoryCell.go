package memory

import "github.com/aclements/go-z3/z3"

const (
	Array = iota
	Struct
)

type SymMemoryPtr = z3.Uninterpreted

type SymMemoryCell struct {
	Kind   int
	Fields map[string]z3.Value
	Sorts  map[string]z3.Sort
}
