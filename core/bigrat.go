package core

import (
	"github.com/client9/sexpr/core/symbol"
	"math/big"
)

type bigRat struct {
	val *big.Rat
}

func newBigRat(n *big.Rat) bigRat {
	return bigRat{
		val: n,
	}
}

func (i bigRat) IsInt() bool {
	return i.val.IsInt()
}

func (i bigRat) Denominator() Integer {
	return newBigInt(i.val.Denom())
}

func (i bigRat) Numerator() Integer {
	return newBigInt(i.val.Num())
}

// Integer type implementation
func (i bigRat) String() string {
	return i.val.String()
}

func (i bigRat) InputForm() string {
	return i.String()
}

func (i bigRat) Head() Expr {
	return symbol.Rational
}

func (i bigRat) Length() int64 {
	return 0
}

func (i bigRat) Equal(rhs Expr) bool {

	switch ratval := rhs.(type) {
	case bigRat:
		return i.val.Cmp(ratval.val) == 0
		/* TODO
		case rat64:
			return i.IsInt64() && i.Int64() == intval.Int64()
		*/
	default:
		return false
	}
}

func (i bigRat) IsAtom() bool {
	return true
}

func (i bigRat) Float64() float64 {
	val, _ := i.val.Float64()
	return val
}

func (i bigRat) Sign() int {
	return i.val.Sign()
}

func (i bigRat) Neg() Rational {
	return bigRat{
		val: new(big.Rat).Neg(i.val),
	}
}

func (i bigRat) Inv() Rational {
	return bigRat{
		val: new(big.Rat).Inv(i.val),
	}
}

// DOES NOT MAKE A COPY.  INTERNAL USE ONLY
func (i bigRat) asBigRat() bigRat {
	return i
}

// DESTRUCTIIVE
func (i *bigRat) add(n bigRat) {
	i.val.Add(i.val, n.val)
}

// DESTRUCTIVE
func (i *bigRat) times(n bigRat) {
	i.val.Mul(i.val, n.asBigRat().val)
}
