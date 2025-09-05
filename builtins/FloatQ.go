package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol FloatQ
// TODO: use helper not direct type cast

// FloatQ checks if an expression is a float
//
// @ExprPattern (_)
func FloatQ(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	_, ok := args[0].(core.Real)
	return core.NewBool(ok)
}
