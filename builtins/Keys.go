package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Keys

// @ExprPattern (_Association)
func Keys(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	assoc := args[0].(core.Association)
	return core.NewList(symbol.List, assoc.Keys()...)
}
