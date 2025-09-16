package core

import (
	"math"
	"strconv"
	"strings"

	"github.com/client9/cardinal/core/symbol"
)

type f64 float64

func NewReal(f float64) Real { return f64(f) }

// for testing?
func MustFloat64(e Expr) float64 {
	return e.(Real).Float64()
}

func (r f64) Prec() uint {
	return 53
}

// Real type implementation
func (r f64) String() string {
	str := strconv.FormatFloat(float64(r), 'f', -1, 64)
	if !strings.Contains(str, ".") {
		str += ".0"
	}
	return str
}

func (r f64) Neg() Real {
	return -r
}

func (r f64) InputForm() string {
	return r.String()
}

func (r f64) Head() Expr {
	return symbol.Real
}

func (r f64) Length() int64 {
	return 0
}

func (r f64) Equal(rhs Expr) bool {
	if other, ok := rhs.(f64); ok {
		return r == other
	}
	return false
}

func (r f64) IsAtom() bool {
	return true
}

func (r f64) IsFloat64() bool {
	return true
}

func (r f64) Float64() float64 {
	return float64(r)
}
func (r f64) AsBigFloat() BigFloat {
	return NewFloat(float64(r))
}

func (r f64) IsInt() bool {
	f := float64(r)
	return float64(f) == math.Round(f)
}

func (r f64) Sign() int {
	if r < 0.0 {
		return -1
	}
	if r > 0.0 {
		return 1
	}
	return 0
}
