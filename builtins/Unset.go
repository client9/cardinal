package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Unset
// @ExprAttributes HoldFirst

// Unset implements the Unset special form
// @ExprPattern (_)
func Unset(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if err := c.Delete(args[0].(core.Symbol)); err != nil {
		return core.NewError("Protected", "unable to under protected symbol")
	}
	return symbol.Null
}
