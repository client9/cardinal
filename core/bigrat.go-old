package core

import (
	"github.com/client9/cardinal/core/symbol"
	"math/big"
)

type BigRat struct {
	val *big.Rat
}

func NewRat(a, b int64) BigRat {
	return BigRat{
		val: big.NewRat(a, b),
	}
}

func newBigRat(n *big.Rat) BigRat {
	return BigRat{
		val: n,
	}
}

func (i BigRat) IsInt() bool {
	return i.val.IsInt()
}

func (i BigRat) Denominator() Integer {
	return newBigInt(i.val.Denom())
}

func (i BigRat) Numerator() Integer {
	return newBigInt(i.val.Num())
}

// Integer type implementation
func (i BigRat) String() string {
	return i.val.String()
}

func (i BigRat) InputForm() string {
	return i.String()
}

func (i BigRat) Head() Expr {
	return symbol.Rational
}

func (i BigRat) Length() int64 {
	return 0
}

func (i BigRat) Equal(rhs Expr) bool {

	switch ratval := rhs.(type) {
	case BigRat:
		return i.val.Cmp(ratval.val) == 0
		/* TODO
		case rat64:
			return i.IsInt64() && i.Int64() == intval.Int64()
		*/
	default:
		return false
	}
}

func (i BigRat) IsAtom() bool {
	return true
}

func (i BigRat) Float64() float64 {
	val, _ := i.val.Float64()
	return val
}

func (i BigRat) Sign() int {
	return i.val.Sign()
}

func (i BigRat) Neg() Rational {
	return BigRat{
		val: new(big.Rat).Neg(i.val),
	}
}

func (i BigRat) Inv() Rational {
	return BigRat{
		val: new(big.Rat).Inv(i.val),
	}
}

// DOES NOT MAKE A COPY.  INTERNAL USE ONLY
func (i BigRat) AsBigRat() BigRat {
	return i
}

func (z *BigRat) Add(x, y *BigRat) *BigRat {
	z.val.Add(x.val, y.val)
	return z
}
func (z *BigRat) AddInt(x *BigRat, y *BigInt) *BigRat {
	return z.Add(x, new(BigRat).SetInt(y))
}

func (z *BigRat) Mul(x, y *BigRat) *BigRat {
	z.val.Mul(x.val, y.val)
	return z
}

func (z *BigRat) MulInt(x *BigRat, y *BigInt) *BigRat {
	return z.Mul(x, new(BigRat).SetInt(y))
}

func (z *BigRat) Set(x *BigRat) *BigRat {
	z.val.Set(x.val)
	return z
}

func (z *BigRat) SetInt(x *BigInt) *BigRat {
	z.val.SetInt(x.val)
	return z
}
