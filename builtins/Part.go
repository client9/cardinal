package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Part

// PartList extracts an element from a list by integer index (1-based)
// @ExprPattern (_, _Integer)
func PartList(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	n, _ := core.ExtractInt64(args[1])
	return core.Part(expr, n)
}
