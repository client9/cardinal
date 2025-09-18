package big

import (
	"github.com/client9/cardinal/core/symbol"
)

// Extends Int to include Expr interface
func (i *Int) InputForm() string {
	return i.String()
}

func (i *Int) Head() Expr {
	return symbol.Integer
}

func (i *Int) Length() int64 {
	return 0
}

func (i *Int) Equal(rhs Expr) bool {
	switch intval := rhs.(type) {
	//case machineInt:
	//		return i.IsInt64() && i.Int64() == intval.Int64()
	case *Int:
		return i.Cmp(intval) == 0
	default:
		return false
	}
}

func (i *Int) IsAtom() bool {
	return true
}
func (i *Int) AsBigInt() *Int {
	return i
}
func (i *Int) AsNeg() Expr {
	return new(Int).Neg(i)
}
func (i *Int) AsInv() Expr {
	// convert int to rat
	// then invert
	tmp := new(Rat).SetInt(i)
	return tmp.Inv(tmp)
}
