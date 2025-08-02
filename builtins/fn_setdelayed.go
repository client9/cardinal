package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// SetDelayedExpr evaluates delayed assignment: SetDelayed(lhs, rhs)
func SetDelayedExpr(e *engine.Evaluator, c *engine.Context, lhs, rhs core.Expr) core.Expr {

	// Handle function definitions: f(x_) := body
	if list, ok := lhs.(core.List); ok && len(list.Elements) >= 1 {
		// This is a function definition
		headExpr := list.Elements[0]
		if _, ok := core.ExtractSymbol(headExpr); ok {
			// Get the function registry from context
			registry := c.GetFunctionRegistry()

			// Register the pattern with the function registry
			err := registry.RegisterUserFunction(lhs, rhs)

			if err != nil {
				return core.NewErrorExpr("DefinitionError", err.Error(), []core.Expr{lhs, rhs})
			}

			return core.NewSymbol("Null")
		}
	}

	// Handle simple variable assignment: x := value
	if symbolName, ok := core.ExtractSymbol(lhs); ok {
		// Store the right-hand side without evaluation (delayed)
		if err := c.Set(symbolName, rhs); err != nil {
			return core.NewErrorExpr("Protected", err.Error(), []core.Expr{lhs})
		}
		return core.NewSymbol("Null")
	}

	return core.NewErrorExpr("SetDelayedError", "Invalid assignment target", []core.Expr{lhs})
}

