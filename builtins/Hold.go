package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Hold
// @ExprAttributes HoldAll

// HoldExpr prevents evaluation of its arguments: Hold(expr1, expr2, ...)
//
// @ExprPattern (___)
func HoldExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	// Hold returns its arguments unevaluated wrapped in Hold
	return core.ListFrom(symbol.Hold, args...)
}
