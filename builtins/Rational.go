package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Rational
// @ExprAttributes Protected
//
//

// @ExprPattern (_Integer,_Integer)
func Rationalize(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewRational(args[0].(core.Integer).Int64(), args[1].(core.Integer).Int64())
}
