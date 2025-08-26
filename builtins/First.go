package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol First

// FirstExpr returns the first element of a list (after the head)
// @ExprPattern (_)
func FirstExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.First(args[0])
}
