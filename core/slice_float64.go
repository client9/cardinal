package core

import (
	"fmt"
	"slices"
	"strings"
)

type SliceFloat64 []float64

func (s SliceFloat64) Length() int64 {
	return int64(len(s))
}

func (s SliceFloat64) String() string {
	parts := make([]string, len(s))
	for i, val := range s {
		parts[i] = fmt.Sprintf("%g", val)
	}
	return "List(" + strings.Join(parts, ", ") + ")"
}

func (s SliceFloat64) InputForm() string {
	return s.String()
}

func (s SliceFloat64) Head() Expr {
	return symbolList
}

func (s SliceFloat64) IsAtom() bool {
	return true
}

func (s SliceFloat64) Equal(a Expr) bool {
	if rhs, ok := a.(SliceFloat64); ok {
		return slices.Equal(s, rhs)
	}
	return false
}

func (s SliceFloat64) ElementAt(n int64) Expr {
	out, err := ElementAt(s, int(n))
	if err != nil {
		return NewError("PartError", err.Error())
	}
	return NewReal(out)
}

func (s SliceFloat64) Slice(i, j int64) Expr {
	out, err := Slice(s, int(i), int(j))
	if err != nil {
		return NewError("PartError", err.Error())
	}
	return out
}

func (s SliceFloat64) Join(other Expr) Expr {
	rhs, ok := other.(SliceFloat64)
	if !ok {
		return NewError("PartError", "Not a list of float64")
	}
	out, err := Join(s, rhs)
	if err != nil {
		return NewError("PartError", err.Error())
	}
	return out
}

func (s SliceFloat64) SetElementAt(n int64, val float64) Expr {
	out, err := SetElementAt(s, int(n), val)
	if err != nil {
		return NewError("PartError", err.Error())
	}
	return out
}

func (s SliceFloat64) SetSlice(i, j int64, values Expr) Expr {
	rhs, ok := values.(SliceFloat64)
	if !ok {
		return NewError("PartError", "Not a list of float64")
	}
	out, err := Replace(s, int(i), int(j), rhs)
	if err != nil {
		return NewError("PartError", err.Error())
	}
	return out
}
