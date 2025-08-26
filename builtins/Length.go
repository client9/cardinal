package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Length

// LengthExpr returns the length of an expression
// @ExprPattern (_)
func LengthExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewInteger(args[0].Length())
}
