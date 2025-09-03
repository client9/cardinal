package core

import (
	"math/big"
	"strconv"
)

type rat64 struct {
	a int64
	b int64
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

func (m rat64) Denominator() Integer {
	return newMachineInt(m.b)
}

func (m rat64) Numerator() Integer {
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

func (m rat64) HeadExpr() Symbol {
	return symbolRational
}

func (m rat64) asBigRat() bigRat {
	return newBigRat(big.NewRat(m.a, m.b))
}
