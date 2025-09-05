package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol FullForm

// FullFormExpr returns the full string representation of an expression
// @ExprPattern (_)
func FullFormExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	// For now, just return the string representation
	// Pattern conversion logic will be added when patterns are moved to core
	return core.NewString(args[0].String())
}
