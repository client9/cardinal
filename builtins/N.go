package builtins

import (
	"math"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/big"
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

// @ExprPattern (E, _Number)
func N_E_Prec(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	prec := args[1].(core.Integer).Int64()
	if prec <= 53 {
		return core.NewReal(math.E)
	}
	return new(big.Float).SetPrec(uint(prec)).E()
}

/*
// Optimized algorithms for E^
// @ExprPattern (Power(E, _Number))
func N_Pow_E(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	powargs := args[0].(core.List).Tail()
	// base = powargs[0]
	exp := powargs[1].(core.Number).Float64()

	// machine level precision
	val, err := core.PowerFloat64(math.E, exp)
	if err != nil {
		return core.NewError("MATH", "MATH")
	}
	return core.NewReal(val)
}
*/

// Optimized algorithms for E^exp
// It's done here in N, otherwise generic algorithm for Power is used
//
// generic N(E^2, 200) --> Power(N(E, 200), N(2,200)) -> mpfr.Pow(...)
// special N(E^2, 200) --> mpfr.Exp(...)
//
// @ExprPattern (Power(E, _Number), _Number)
func N_Pow_E_Prec(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	powargs := args[0].(core.List).Tail()
	prec := args[1].(core.Integer).Int64()
	exp := powargs[1].(core.Number)

	// set exponent default precision
	//   will get adjused by ToBigFloat
	expflt := new(big.Float).SetPrec(uint(prec))
	core.ToBigFloat(expflt, exp)

	z := new(big.Float).SetPrec(core.UMin(uint(prec), expflt.Prec())).Exp(expflt)
	return z
}

// @ExprPattern (Pi)
func N_Pi(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(math.Pi)
}

// @ExprPattern (Pi, _Number)
func N_Pi_Prec(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	prec := args[1].(core.Integer).Int64()
	if prec <= 53 {
		return core.NewReal(math.Pi)
	}
	return new(big.Float).SetPrec(uint(prec)).Pi()
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
	return new(big.Float).SetPrec(uint(prec)).SetRat(i.AsBigRat())
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
	return new(big.Float).SetPrec(uint(prec)).SetInt(i.AsBigInt())
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
