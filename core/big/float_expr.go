package big

import (
	"github.com/client9/bignum/mpfr"
	"github.com/client9/cardinal/core/symbol"
)

// Extends Int to include Expr interface
func (i *Float) InputForm() string {
	return i.String()
}

func (i *Float) Head() Expr {
	return symbol.Real
}

func (i *Float) Length() int64 {
	return 0
}

func (i *Float) Equal(rhs Expr) bool {
	if other, ok := rhs.(*Float); ok {
		return i.Cmp(other) == 0
	}
	if other, ok := rhs.(symbol.NumberExpr); ok {
		return i.Float64() == other.Float64()
	}
	return false
}

func (i *Float) IsAtom() bool {
	return true
}

func (i Float) IsFloat64() bool {
	return false
}
func (i *Float) AsBigFloat() *Float {
	return i
}
func (i *Float) AsNeg() Expr {
	return new(Float).Neg(i)
}

func (i *Float) AsInv() Expr {
	// Go big.Float and MPFR doens't have an Inv for floats.
	//  But MPFR does have a sepcialization of (int64 / BigFloat)
	z := new(Float).SetPrec(i.Prec())
	mpfr.UiDiv(z.ptr, 1, i.ptr, i.mode)
	return z
}

func (z *Float) Pow(x, y *Float) *Float {

	if z.ptr == nil {
		z.init()
	}
	if z.Prec() == 0 {
		z.SetPrec(x.Prec())
	}

	mpfr.Pow(z.ptr, x.ptr, y.ptr, z.mode)
	return z
}
func (z *Float) Sin(x *Float) *Float {
	if z.ptr == nil {
		z.init()
	}
	if z.Prec() == 0 {
		z.SetPrec(x.Prec())
	}
	mpfr.Sin(z.ptr, x.ptr, z.mode)
	return z
}

func (z *Float) Cos(x *Float) *Float {
	if z.ptr == nil {
		z.init()
	}
	if z.Prec() == 0 {
		z.SetPrec(x.Prec())
	}
	mpfr.Cos(z.ptr, x.ptr, z.mode)
	return z
}
func (z *Float) Tan(x *Float) *Float {
	if z.ptr == nil {
		z.init()
	}
	if z.Prec() == 0 {
		z.SetPrec(x.Prec())
	}
	mpfr.Tan(z.ptr, x.ptr, z.mode)
	return z
}

func (z *Float) Pi() *Float {
	if z.ptr == nil {
		z.init()
	}
	mpfr.ConstPi(z.ptr, z.mode)
	return z
}

// TODO: Cache
func (z *Float) E() *Float {
	if z.ptr == nil {
		z.init()
	}
	mpfr.Exp(z.ptr, NewFloat(1.0).ptr, z.mode)
	return z
}

// Exp is Power(E, op)
func (z *Float) Exp(x *Float) *Float {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if z.Prec() == 0 {
		z.SetPrec(x.Prec())
	}
	mpfr.Exp(z.ptr, x.ptr, z.mode)
	return z
}
