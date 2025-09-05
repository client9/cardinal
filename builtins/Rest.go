package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Rest

// RestExpr returns a new list with the first element after head removed
// @ExprPattern (_)
func RestExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.Rest(args[0])
}
