package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Abs

// @ExprPattern (_Integer)
func AbsInteger(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	i := args[0].(core.Integer)
	if i.Sign() == -1 {
		return i.Neg()
	}
	return i
}

// @ExprPattern (_Rational)
func AbsRational(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	i := args[0].(core.Rational)
	if i.Sign() == -1 {
		return i.Neg()
	}
	return i
}

// @ExprPattern (_Real)
func AbsReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	i := args[0].(core.Real)
	if i.Sign() == -1 {
		return i.Neg()
	}
	return i
}

// @ExprPattern (Times(-1,_))
func AbsTimes(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	// Abs(Times(-1, x)) is x
	return args[0].(core.List).Tail()[1]
}
