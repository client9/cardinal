package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol LessEqual

// @ExprPattern (_Integer, _Real)
func LessEqualInt64Float64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := float64(core.MustInt64(args[0]))
	y := core.MustFloat64(args[1])
	return core.NewBool(x <= y)
}

// @ExprPattern (_Real, _Integer)
func LessEqualFloat64Int64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustFloat64(args[0])
	y := float64(core.MustInt64(args[1]))
	return core.NewBool(x <= y)
}

// @ExprPattern (_Integer, _Integer)
func LessEqualInts(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustInt64(args[0])
	y := core.MustInt64(args[1])
	return core.NewBool(x <= y)
}

// @ExprPattern (_Real, _Real)
func LessEqualFloat64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustFloat64(args[0])
	y := core.MustFloat64(args[1])
	return core.NewBool(x <= y)
}
