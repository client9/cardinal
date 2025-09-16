package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol IntegerPart

// @ExprPattern (_Integer)
func IntegerPart(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0]
}

// @ExprPattern (_Rational)
func IntegerPartRational(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewInteger(int64(args[0].(core.Rational).Float64()))
}

// @ExprPattern (_Real)
func IntegerPartReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewInteger(int64(args[0].(core.Real).Float64()))
}
