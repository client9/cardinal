package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Most

// MostExpr returns a new list with the last element removed
// @ExprPattern (_)
func MostExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.Most(args[0])
}
