package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Precision

// @ExprPattern (_Integer)
func PrecisionInteger(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return symbol.Infinity
}

// @ExprPattern (_Rational)
func PrecisionRational(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return symbol.Infinity
}

// @ExprPattern (_Real)
func PrecisionReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	r := args[0].(core.Real)
	return core.NewInteger(int64(r.Prec()))
}
