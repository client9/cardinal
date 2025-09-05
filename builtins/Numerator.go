package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Numerator
// @ExprAttributes Protected
//
//

// @ExprPattern (_Integer)
func NumeratorInteger(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0]
}

// @ExprPattern (_Real)
func NumeratorReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0]
}

// @ExprPattern (_Rational)
func NumeratorRational(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0].(core.Rational).Numerator()
}
