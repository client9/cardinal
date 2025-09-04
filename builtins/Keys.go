package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Keys

// @ExprPattern (_Association)
func Keys(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	assoc := args[0].(core.Association)
	return core.NewList("List", assoc.Keys()...)
}
