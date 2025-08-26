package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Minus

// MinusInteger returns the negation of an integer
// @ExprPattern (_Integer)
func MinusInteger(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x, _ := core.ExtractInt64(args[0])
	return core.NewInteger(-x)
}

// MinusReal returns the negation of a real number
// @ExprPattern (_Real)
func MinusReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x, _ := core.ExtractFloat64(args[0])
	return core.NewReal(-x)
}

// MinusExpr converts Minus(x) to Times(-1, x) as per Mathematica
// @ExprPattern (_)
func MinusExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := args[0]
	return core.NewList("Times", core.NewInteger(-1), x)
}
