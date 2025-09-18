package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/big"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol SetPrecision

// TrueQ check is an expression is explicity True
// @ExprPattern (_Real, n_Integer)
func SetPrecisionReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	r := args[0].(core.Real)
	prec := args[1].(core.Integer).Int64()
	if prec <= 53 {
		return core.NewReal(r.Float64())
	}
	return new(big.Float).Set(r.AsBigFloat()).SetPrec(uint(prec))
}
