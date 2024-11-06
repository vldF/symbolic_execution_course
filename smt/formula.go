package smt

import (
	"github.com/aclements/go-z3/z3"
)

type Formula interface {
	Value() z3.Bool
}
