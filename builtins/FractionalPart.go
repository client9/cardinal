package builtins

import (
	//	"fmt"
	"math"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/big"
	"github.com/client9/cardinal/engine"
)

// Given a n-bit precision number (mantissa)
//  and e such that n * 2^e
// @ExprSymbol FractionalPart

// @ExprPattern (_Integer)
func FractionalPart(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(0)
}

// @ExprPattern (_Rational)
func FractionalPartRational(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(math.Trunc(args[0].(core.Rational).Float64()))
}

// @ExprPattern (_Real)
func FractionalPartReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	arg := args[0].(core.Real)

	if arg.IsFloat64() {
		flt := args[0].(core.Real).Float64()
		mantissa, exp := math.Frexp(flt)
		if exp >= 52 {
			return core.NewInteger(0)
		}
		if exp >= 0 {
			tmp := args[0].(core.Real).Float64()
			return core.NewReal(tmp - math.Trunc(tmp))
		}
		// it's a small number
		// This is
		m := big.NewFloat(mantissa)
		z := new(big.Float).SetMantExp(m, exp)
		return z.Frac(z)
	}
	z := arg.AsBigFloat()
	return z.Frac(z)
}
