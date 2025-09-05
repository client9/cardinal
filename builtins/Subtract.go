package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Subtract

// SubtractIntegers performs integer subtraction
// @ExprPattern (_Integer, _Integer)
func SubtractIntegers(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x, _ := core.ExtractInt64(args[0])
	y, _ := core.ExtractInt64(args[1])
	return core.NewInteger(x - y)
}

// SubtractNumbers performs mixed numeric subtraction (returns float64)
// @ExprPattern (_Real, _Real)
func SubtractNumbers(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x, _ := core.ExtractFloat64(args[0])
	y, _ := core.ExtractFloat64(args[1])
	return core.NewReal(x - y)
}
