package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// IfExpr evaluates conditional expressions: If(condition, then) or If(condition, then, else)
func IfExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	if len(args) < 2 || len(args) > 3 {
		return core.NewErrorExpr("ArgumentError",
			"If expects 2 or 3 arguments", args)
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
		return core.NewSymbolNull()
	}

	// Condition is not a boolean, return an error
	return core.NewErrorExpr("TypeError", "If condition must be True or False", []core.Expr{condition})
}
