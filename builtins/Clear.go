package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Clear
// @ExprAttributes HoldAll

// @ExprPattern (___Symbol)
func Clear(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	for _, arg := range args {
		c.Clear(arg.(core.Symbol))
	}
	return symbol.Null
}
