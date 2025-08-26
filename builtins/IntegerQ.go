package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol IntegerQ

// IntegerQExpr checks if an expression is an integer
//
// @ExprPattern (_)
func IntegerQ(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	_, ok := args[0].(core.Integer)
	return core.NewBool(ok)
}
