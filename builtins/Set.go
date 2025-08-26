package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Set
// @ExprAttributes HoldFirst

// SetExpr evaluates immediate assignment: Set(lhs, rhs)
// @ExprPattern (_,_)
func SetExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	lhs := args[0]
	rhs := args[1]
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
