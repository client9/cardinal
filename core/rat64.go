package core

import (
	"strconv"

	"github.com/client9/cardinal/core/big"
	"github.com/client9/cardinal/core/symbol"
)

type rat64 struct {
	a int64
	b int64
}

func (m rat64) Int() Integer {
	return newMachineInt(m.a)
}

var rat64Zero = rat64{0, 1}

var rat64One = rat64{1, 1}

func (m rat64) IsInt() bool {
	return m.b == 1
}

func (m rat64) Sign() int {
	n := m.a
	if n == 0 {
		return 0
	}
	if n > 0 {
		return 1
	}
	return -1
}

func (m rat64) StandardForm() Expr {
	if m.a == 0 {
		return newMachineInt(0)
	}

	if m.b < 0 {
		m.a = -m.a
		m.b = -m.b
	}

	if m.b == 1 {
		return newMachineInt(m.a)
	}

	if m.a == 1 || m.a == -1 {
		return rat64{m.a, m.b}
	}

	g := gcd(m.a, m.b)
	den := m.b / g
	if den == 1 {
		return newMachineInt(m.a / g)
	}
	return rat64{m.a / g, den}
}

func (m rat64) Denom() machineInt {
	return newMachineInt(m.b)
}

func (m rat64) Num() machineInt {
	return newMachineInt(m.a)
}

func (m rat64) Neg() Rational {
	return rat64{-m.a, m.b}
}

func (m rat64) Inv() Rational {
	return rat64{m.b, m.a}
}

func (m rat64) Float64() float64 {
	return float64(m.a) / float64(m.b)
}

func (m rat64) Length() int64 {
	return 0
}

func (m rat64) IsAtom() bool {
	return true
}

func (m rat64) Equal(rhs Expr) bool {
	if rhs, ok := rhs.(rat64); ok {
		return m.a == rhs.a && m.b == rhs.b
	}
	return false
}

func (m rat64) String() string {
	return strconv.FormatInt(m.a, 10) + "/" + strconv.FormatInt(m.b, 10)
}

func (m rat64) InputForm() string {
	return m.String()
}

func (m rat64) Head() Expr {
	return symbol.Rational
}

func (m rat64) AsBigRat() *big.Rat {
	return big.NewRat(m.a, m.b)
}
func (m rat64) AsNum() Expr {
	return newMachineInt(m.a)
}
func (m rat64) AsDenom() Expr {
	return newMachineInt(m.b)
}
func (m rat64) AsNeg() Expr {
	return rat64{-m.a, m.b}
}
func (m rat64) AsInv() Expr {
	// TODO ERROR
	return rat64{m.b, m.a}
}
