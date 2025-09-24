package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/big"
	"github.com/client9/cardinal/engine"
	"math"
)

// @ExprSymbol Ceiling
// @ExprAttributes Protected
//
//

// @ExprPattern (_Integer)
func CeilingInteger(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0]
}

// @ExprPattern (_Real)
func CeilingReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	r := args[0].(core.Real)
	if r.IsFloat64() {
		return core.NewInteger(int64(math.Ceil(r.Float64())))
	}

	// Big Real
	return new(big.Float).Ceil(r.AsBigFloat()).Int()
}

/*
// @ExprPattern (_Rational)
func FloorRational(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0].(core.Rational).AsNum()
}
*/
