package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// SetExpr evaluates immediate assignment: Set(lhs, rhs)
func SetExpr(e *engine.Evaluator, c *engine.Context, lhs, rhs core.Expr) core.Expr {
	// Evaluate the right-hand side immediately
	evalRhs := e.Evaluate(rhs)

	// Handle assignment to symbol
	if symbolName, ok := core.ExtractSymbol(lhs); ok {
		if err := c.Set(symbolName, evalRhs); err != nil {
			return core.NewError("Protected", err.Error())
		}
		return evalRhs
	}

	return core.NewError("SetError", "Invalid assignment target")
}
