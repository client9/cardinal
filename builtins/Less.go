package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Less

// @ExprPattern (_Integer, _Real)
func LessInt64Float64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := float64(core.MustInt64(args[0]))
	y := core.MustFloat64(args[1])
	return core.NewBool(x < y)
}

// @ExprPattern (_Real, _Integer)
func LessFloat64Int64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustFloat64(args[0])
	y := float64(core.MustInt64(args[1]))
	return core.NewBool(x < y)
}

// @ExprPattern (_Integer, _Integer)
func LessInts(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustInt64(args[0])
	y := core.MustInt64(args[1])
	return core.NewBool(x < y)
}

// @ExprPattern (_Real, _Real)
func LessFloat64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustFloat64(args[0])
	y := core.MustFloat64(args[1])
	return core.NewBool(x < y)
}
