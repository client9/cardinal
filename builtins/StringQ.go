package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol StringQ

// StringQ checks if an expression is a string
//
// @ExprPattern (_)
func StringQ(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	// Check new String type first
	_, ok := expr.(core.String)
	return core.NewBool(ok)
}
