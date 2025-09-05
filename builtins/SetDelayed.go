package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol SetDelayed
// @ExprAttributes HoldAll
//

// SetDelayedExpr evaluates delayed assignment: SetDelayed(lhs, rhs)
// @ExprPattern (_,_)
func SetDelayedExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	lhs := args[0]
	rhs := args[1]
	// Handle function definitions: f(x_) := body
	if list, ok := lhs.(core.List); ok && list.Length() > 0 {
		// This is a function definition
		headExpr := list.Head()
		if _, ok := core.ExtractSymbol(headExpr); ok {
			// Get the function registry from context
			registry := c.GetFunctionRegistry()

			// Register the pattern with the function registry
			err := registry.RegisterUserFunction(lhs, rhs)

			if err != nil {
				return core.NewError("DefinitionError", err.Error())
			}

			return symbol.Null
		}
	}

	// Handle simple variable assignment: x := value
	if symbolName, ok := lhs.(core.Symbol); ok {
		// Store the right-hand side without evaluation (delayed)
		if err := c.Set(symbolName, rhs); err != nil {
			return core.NewError("Protected", err.Error())
		}
		return symbol.Null
	}

	return core.NewError("SetDelayedError", "Invalid assignment target")
}
