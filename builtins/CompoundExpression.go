package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol CompoundExpression
// @ExprAttributes HoldAll

// CompoundExpression implements the CompoundExpression special form
// TODO wrapgen doesn't like args []core.Expr.. forcing a wrapper
//
// @ExprPattern (___)
func CompoundExpression(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	var result core.Expr = symbol.Null

	for _, arg := range args {
		result = e.Evaluate(arg)
		if core.IsError(result) {
			return result
		}
	}

	return result
}
