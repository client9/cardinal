package builtins

import (
	"math"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol N
// @ExprAttributes Protected
//

// @ExprPattern (E)
func N_E(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(math.E)
}

// @ExprPattern (_Rational)
func N_Rational(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	i := args[0].(core.Rational)
	return core.NewReal(i.Float64())
}

// @ExprPattern (_Rational, _Integer)
func N_RationalPrec(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	i := args[0].(core.Rational)
	prec := args[1].(core.Integer).Int64()
	if prec <= 53 {
		return core.NewReal(i.Float64())
	}
	x := i.(core.BigRat)
	return *(new(core.BigFloat).SetPrec(uint(prec)).SetRat(&x))
}

// @ExprPattern (_Integer)
func N_Integer(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(args[0].(core.Integer).Float64())
}

// @ExprPattern (_Integer, _Integer)
func N_IntegerPrec(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	i := args[0].(core.Integer)
	prec := args[1].(core.Integer).Int64()
	if prec <= 53 {
		return core.NewReal(i.Float64())
	}
	x := i.AsBigInt()
	return *(new(core.BigFloat).SetPrec(uint(prec)).SetInt(&x))
}

// @ExprPattern (_)
func N(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return evalNum(e, c, args[0])
}

// @ExprPattern (_,_Integer)
func NPrec(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	prec := args[1].(core.Integer)
	return evalNumPrec(e, c, args[0], prec)
}

func evalNumPrec(e *engine.Evaluator, c *engine.Context, arg core.Expr, prec core.Expr) core.Expr {
	switch num := arg.(type) {
	case core.List:
		elem := make([]core.Expr, 1+int(num.Length()))
		elem[0] = num.Head()
		for i, listarg := range num.Tail() {
			elem[i+1] = e.Evaluate(core.ListFrom(symbol.N, listarg, prec))
		}
		return core.NewListFromExprs(elem...)
	default:
		return arg
	}
}
func evalNum(e *engine.Evaluator, c *engine.Context, arg core.Expr) core.Expr {
	switch num := arg.(type) {
	case core.List:
		elem := make([]core.Expr, 1+int(num.Length()))
		elem[0] = num.Head()
		for i, listarg := range num.Tail() {
			elem[i+1] = e.Evaluate(core.ListFrom(symbol.N, listarg))
		}
		return core.NewListFromExprs(elem...)
	default:
		return arg
	}
}
