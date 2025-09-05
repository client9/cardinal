package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Unequal

// @ExprPattern (_Integer, _Real)
func UnequalInt64Float64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := float64(core.MustInt64(args[0]))
	y := core.MustFloat64(args[1])
	return core.NewBool(x != y)
}

// @ExprPattern (_Real, _Integer)
func UnequalFloat64Int64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustFloat64(args[0])
	y := float64(core.MustInt64(args[1]))
	return core.NewBool(x != y)
}

// @ExprPattern (_,_)
func Unequal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := args[0]
	y := args[1]
	return core.NewBool(!x.Equal(y))
}
