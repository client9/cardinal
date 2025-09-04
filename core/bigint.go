package core

import (
	"math/big"
)

type bigInt struct {
	val *big.Int
}

// mutable zero value
func bigIntZero() bigInt {
	return newBigInt(big.NewInt(0))
}

// mutable one value
func bigIntOne() bigInt {
	return newBigInt(big.NewInt(1))
}

func NewBigInt(n int64) bigInt {
	return bigInt{
		val: big.NewInt(n),
	}
}

func newBigInt(n *big.Int) bigInt {
	return bigInt{
		val: n,
	}
}

// Integer type implementation
func (i bigInt) String() string {
	return i.val.String()
}

func (i bigInt) InputForm() string {
	return i.String()
}

func (i bigInt) Head() Expr {
	return symbolInteger
}

func (i bigInt) Length() int64 {
	return 0
}

func (i bigInt) Equal(rhs Expr) bool {
	switch intval := rhs.(type) {
	case machineInt:
		return i.IsInt64() && i.Int64() == intval.Int64()
	case bigInt:
		return i.val.Cmp(intval.val) == 0
	default:
		return false
	}
}

func (i bigInt) IsAtom() bool {
	return true
}

func (i bigInt) Float64() float64 {
	return i.Float64()
}

func (i bigInt) IsInt64() bool {
	if i.val == nil {
		return true
	}
	return i.val.IsInt64()
}

func (i bigInt) Int64() int64 {
	if i.val == nil {
		return 0
	}
	return i.val.Int64()
}

func (i bigInt) Sign() int {
	if i.val == nil {
		return 0
	}
	return i.val.Sign()
}

func (i bigInt) Inv() Expr {
	num := big.NewInt(1)
	den := i.val
	return newBigRat(new(big.Rat).SetFrac(num, den))
}

func (i bigInt) Neg() Integer {
	return bigInt{
		val: new(big.Int).Neg(i.val),
	}
}

// DOES NOT MAKE A COPY.  INTERNAL USE ONLY
func (i bigInt) asBigInt() bigInt {
	return i
}

// DESTRUCTIIVE
func (i *bigInt) add(n bigInt) {
	i.val.Add(i.val, n.val)
}

// DESTRUCTIVE
func (i *bigInt) times(n bigInt) {
	i.val.Mul(i.val, n.asBigInt().val)
}
