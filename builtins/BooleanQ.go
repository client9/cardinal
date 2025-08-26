package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol BooleanQ
//
// BooleanQExpr checks if an expression is a boolean (True/False symbol)

// @ExprPattern (_)
func BooleanQExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewBool(core.IsBool(args[0]))
}
