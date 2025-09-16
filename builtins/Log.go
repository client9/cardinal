package builtins

import (
	"math"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Log

// @ExprPattern (_Real)
func Log(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	// HACK
	r := args[0].(core.Real).Float64()
	return core.NewReal(math.Log(r))
}

// @ExprPattern (_Real, 2)
func Log2(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	// HACK
	r := args[0].(core.Real).Float64()
	return core.NewReal(math.Log2(r))
}
