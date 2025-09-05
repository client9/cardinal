package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Equal

// @ExprPattern (_Integer, _Real)
func EqualInt64Float64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := float64(core.MustInt64(args[0]))
	y := core.MustFloat64(args[1])
	return core.NewBool(x == y)
}

// @ExprPattern (_Real, _Integer)
func EqualFloat64Int64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustFloat64(args[0])
	y := float64(core.MustInt64(args[1]))
	return core.NewBool(x == y)
}

// @ExprPattern (_,_)
func Equal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := args[0]
	y := args[1]
	return core.NewBool(x.Equal(y))
}

/*
// @ExprPattern (_Integer, _Integer)
func EqualInts(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustInt64(args[0])
	y := core.MustInt64(args[1])
	return core.NewBool(x == y)
}
// @ExprPattern (_Real, _Real)
func EqualFloat64(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := core.MustFloat64(args[0])
	y := core.MustFloat64(args[1])
	return core.NewBool(x == y)
}
// @ExprPattern (_String, _String)
func EqualString(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x, _ := core.ExtractString(args[0])
	y, _ := core.ExtractString(args[1])
	return core.NewBool(x == y)
}
*/
