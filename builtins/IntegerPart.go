package builtins

import (
	"fmt"
	"math"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/big"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol IntegerPart

// @ExprPattern (_Integer)
func IntegerPart(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0]
}

// @ExprPattern (_Rational)
func IntegerPartRational(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewInteger(int64(args[0].(core.Rational).Float64()))
}

// @ExprPattern (_Real)
func IntegerPartReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	arg := args[0].(core.Real)

	if arg.IsFloat64() {
		flt := args[0].(core.Real).Float64()
		mantissa, exp := math.Frexp(flt)
		if exp <= 0 {
			return core.NewInteger(0)
		}
		if exp <= 52 {
			return core.NewInteger(int64(args[0].(core.Real).Float64()))
		}
		// it's a big number
		m := big.NewFloat(mantissa)
		z := new(big.Float).SetMantExp(m, exp)
		fmt.Println("PREC on Final IS ", z.Prec())
		return z.Int()
	}
	flt := arg.AsBigFloat()
	return flt.Int()
}
