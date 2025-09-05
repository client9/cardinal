package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol SymbolName

// MMA resticts to valid identifier symbols (i.e. A-Za-z0-9).
// No restrictions here.

// NewSymbol creates a Symbol expressoin from a string.
// @ExprPattern (_Symbol)
func SymbolName(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewString(args[0].String())
}
