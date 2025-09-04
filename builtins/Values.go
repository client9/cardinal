package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Values

// @ExprPattern (_Association)
func Values(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	assoc := args[0].(core.Association)
	values := assoc.Values()
	return core.NewList(symbol.List, values...)
}
