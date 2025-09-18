package big

import (
	"github.com/client9/cardinal/core/symbol"
)

// Extends Int to include Expr interface
func (i *Rat) InputForm() string {
	return i.String()
}

func (i *Rat) Head() Expr {
	return symbol.Rational
}

func (i *Rat) Length() int64 {
	return 0
}

func (i *Rat) Equal(rhs Expr) bool {
	switch intval := rhs.(type) {
	case *Rat:
		return i.Cmp(intval) == 0
	default:
		return false
	}
}

func (i *Rat) IsAtom() bool {
	return true
}
func (i *Rat) AsBigRat() *Rat {
	return i
}

func (i *Rat) AsDenom() Expr {
	return i.Denom()
}

func (i *Rat) AsNum() Expr {
	return i.Num()
}
func (i *Rat) AsNeg() Expr {
	return new(Rat).Neg(i)
}
func (i *Rat) AsInv() Expr {
	return new(Rat).Inv(i)
}
