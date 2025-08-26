package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Evaluate
//

// EvaluateExpr forces evaluation: Evaluate(expr)
//
// @ExprPattern (_)
func EvaluateExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return e.Evaluate(args[0])
}
