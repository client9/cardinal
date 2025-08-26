package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol AtomQ

// AtomQ checks if an expression is an atom
// @ExprPattern (_)
func AtomQ(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewBool(args[0].IsAtom())
}
