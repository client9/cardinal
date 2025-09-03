package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Boole

// @ExprPattern (_)
func Boole(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	arg := args[0]
	if val, ok := core.ExtractBool(arg); ok {
		if val {
			return core.NewInteger(1)
		}
		return core.NewInteger(0)
	}

	// Return unchanged expression if not boolean (symbolic behavior)
	return core.ListFrom(symbol.Boole, arg)
}
