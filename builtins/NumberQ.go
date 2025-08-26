package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol NumberQ

// NumberQ checks if an expression is numeric (int or float)
//
// @ExprPattern (_)
func NumberQ(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewBool(core.IsNumeric(args[0]))
}
