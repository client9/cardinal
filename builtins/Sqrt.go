package builtins

import (
	"math"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Sqrt
// @ExprAttributes Listable NumericFunction Protected

// @ExprPattern (_Real)
func Sqrt(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	r := args[0].(core.Real)
	if r.Prec() <= 53 {
		return core.NewReal(math.Sqrt(r.Float64()))
	}
	x := r.(core.BigFloat)
	return *new(core.BigFloat).Sqrt(&x)
}
