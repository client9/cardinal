package core

import (
	"strconv"
	"strings"
)

type Real float64

func NewReal(f float64) Real { return Real(f) }

func MustFloat64(e Expr) float64 {
	return float64(e.(Real))
}

// Real type implementation
func (r Real) String() string {
	str := strconv.FormatFloat(float64(r), 'f', -1, 64)
	if !strings.Contains(str, ".") {
		str += ".0"
	}
	return str
}

func (r Real) Neg() Real {
	return -r
}

func (r Real) InputForm() string {
	return r.String()
}

func (r Real) Head() Expr {
	return symbolReal
}

func (r Real) Length() int64 {
	return 0
}

func (r Real) Equal(rhs Expr) bool {
	if other, ok := rhs.(Real); ok {
		return r == other
	}
	return false
}

func (r Real) IsAtom() bool {
	return true
}
