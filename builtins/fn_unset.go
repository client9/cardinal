package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// Unset implements the Unset special form
func Unset(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) != 1 {
		return core.NewErrorExpr("ArgumentError", "Unset expects 1 argument", args)
	}

	if symbolName, ok := core.ExtractSymbol(args[0]); ok {
		if err := c.Delete(symbolName); err != nil {
			return core.NewErrorExpr("Protected", "unable to under protected symbol", args)
		}
		return core.NewSymbolNull()
	}

	return core.NewErrorExpr("ArgumentError", "Argument to Unset must be a symbol", args)
}
