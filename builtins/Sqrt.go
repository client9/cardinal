package builtins

import (
	"fmt"
	"math"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/big"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Sqrt
// @ExprAttributes Listable NumericFunction Protected

// @ExprPattern (_Real)
func SqrtReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	r := args[0].(core.Real)
	if r.Prec() <= 53 {
		return core.NewReal(math.Sqrt(r.Float64()))
	}
	fmt.Println("SQRT PREC=", r.Prec())
	return new(big.Float).SetPrec(r.Prec()).Sqrt(r.AsBigFloat())
}

// Sqrt is symbolically converted to Power(x, Rational(1,2))
//
// @ExprPattern (_)
func Sqrt(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.ListFrom(symbol.Power, args[0], core.NewRational(1, 2))
}
