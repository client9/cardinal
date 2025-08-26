package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol TrueQ
// TODO: Replace with direct atom comparison

// TrueQ check is an expression is explicity True
// @ExprPattern (_)
func TrueQ(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	// Check new Symbol type first
	if symbolName, ok := core.ExtractSymbol(expr); ok {
		return core.NewBool(symbolName == "True")
	}
	return core.NewBool(false)
}
