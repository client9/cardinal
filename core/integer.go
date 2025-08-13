package core

import (
	"strconv"
)

type Integer int64

// New constructor functions for atomic types
func NewInteger(i int64) Integer { return Integer(i) }

// Integer type implementation
func (i Integer) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Integer) InputForm() string {
	return i.String()
}

func (i Integer) Head() string {
	return "Integer"
}

func (i Integer) Length() int64 {
	return 0
}

func (i Integer) Equal(rhs Expr) bool {
	if other, ok := rhs.(Integer); ok {
		return i == other
	}
	return false
}

func (i Integer) IsAtom() bool {
	return true
}
