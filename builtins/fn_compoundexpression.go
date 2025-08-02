package builtins

import (
        "github.com/client9/sexpr/core"
        "github.com/client9/sexpr/engine"
)

// CompoundExpression implements the CompoundExpression special form
// TODO wrapgen doesn't like args []core.Expr.. forcing a wrapper
func CompoundExpression(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	var result core.Expr = core.NewSymbolNull()

	for _, arg := range args {
		result = e.Evaluate(c, arg)
		if core.IsError(result) {
			return result
		}
	}

	return result
}
