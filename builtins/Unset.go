package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Unset
// @ExprAttributes HoldFirst

// Unset implements the Unset special form
// @ExprPattern (_)
func Unset(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) != 1 {
		return core.NewError("ArgumentError", "Unset expects 1 argument")
	}

	if symbolName, ok := core.ExtractSymbol(args[0]); ok {
		if err := c.Delete(symbolName); err != nil {
			return core.NewError("Protected", "unable to under protected symbol")
		}
		return core.NewSymbolNull()
	}

	return core.NewError("ArgumentError", "Argument to Unset must be a symbol")
}
