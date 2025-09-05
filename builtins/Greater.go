package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol GreaterEqual

// @ExprPattern (_Integer, _Real)
func GreaterEqualInt64Float64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := float64(core.MustInt64(args[0]))
	y := core.MustFloat64(args[1])
	return core.NewBool(x >= y)
}

// @ExprPattern (_Real, _Integer)
func GreaterEqualFloat64Int64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustFloat64(args[0])
	y := float64(core.MustInt64(args[1]))
	return core.NewBool(x >= y)
}

// @ExprPattern (_Integer, _Integer)
func GreaterEqualInts(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustInt64(args[0])
	y := core.MustInt64(args[1])
	return core.NewBool(x >= y)
}

// @ExprPattern (_Real, _Real)
func GreaterEqualFloat64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustFloat64(args[0])
	y := core.MustFloat64(args[1])
	return core.NewBool(x >= y)
}
