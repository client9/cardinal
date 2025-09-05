package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol InputForm

// InputFormExpr returns the user-friendly InputForm representation of an expression
// @ExprPattern (_)
func InputFormExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	// For now, just return the InputForm representation
	// Pattern conversion logic will be added when patterns are moved to core
	return core.NewString(args[0].InputForm())
}
