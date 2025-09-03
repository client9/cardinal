package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Denominator
// @ExprAttributes Protected
//
//

// @ExprPattern (_Integer)
func DenominatorInteger(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewInteger(1)
}

// @ExprPattern (_Real)
func DenominatorReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewInteger(1)
}

// @ExprPattern (_Rational)
func DenominatorRational(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0].(core.Rational).Denominator()
}
