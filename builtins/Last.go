package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Last

// LastExpr returns the last element of a list
// @ExprPattern (_)
func LastExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.Last(args[0])
}
