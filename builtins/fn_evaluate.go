package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// EvaluateExpr forces evaluation: Evaluate(expr)
func EvaluateExpr(e *engine.Evaluator, c *engine.Context,  arg core.Expr) core.Expr {
	return e.Evaluate(c, arg)
}

