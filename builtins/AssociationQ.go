package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol AssociationQ

// @ExprPattern (_)
func AssociationQ(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	_, ok := args[0].(core.Association)
	return core.NewBool(ok)
}
