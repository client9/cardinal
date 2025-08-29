package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol If
// @ExprAttributes HoldRest

// IfExpr evaluates conditional expressions: If(condition, then) or If(condition, then, else)
// @ExprPattern (___)
func IfExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	if len(args) < 2 || len(args) > 3 {
		return core.NewError("ArgumentError",
			"If expects 2 or 3 arguments")
	}

	// Evaluate the condition
	condition := e.Evaluate(args[0])
	if core.IsError(condition) {
		return condition
	}

	// Check the condition
	if boolVal, ok := core.ExtractBool(condition); ok {
		if boolVal {
			// Condition is true, evaluate and return the "then" branch
			return e.Evaluate(args[1])
		}
		// Condition is false, evaluate and return the "else" branch if present
		if len(args) == 3 {
			return e.Evaluate(args[2])
		}
		return symbol.Null
	}

	// Condition is not a boolean, return an error
	return core.NewError("TypeError", "If condition must be True or False")
}
