package builtins

import (
	"math"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/big"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Tan
// @ExprAttributes Listable NumericFunction Protected

// @ExprPattern (_Real)
func Tan(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	r := args[0].(core.Real)
	if r.Prec() <= 53 {
		return core.NewReal(math.Tan(r.Float64()))
	}
	return new(big.Float).SetPrec(r.Prec()).Tan(r.AsBigFloat())
}
