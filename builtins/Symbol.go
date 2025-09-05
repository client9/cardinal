package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Symbol

// MMA resticts to valid identifier symbols (i.e. A-Za-z0-9).
// No restrictions here.

// NewSymbol creates a Symbol expressoin from a string.
// @ExprPattern (_String)
func NewSymbol(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	s := string(args[0].(core.String))
	return core.NewSymbol(s)
}
