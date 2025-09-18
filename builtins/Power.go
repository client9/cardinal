package builtins

import (
	"fmt"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/big"
	"github.com/client9/cardinal/core/symbol"

	"github.com/client9/cardinal/engine"
)

// Power(_Integer, _Rational) simplification is non-obvious
//  4^(1/2) --> (2^2)^(1/2) --> 2^1 --> 2
//  4^(3/5) --> (2^2)^(3/5) --> 2^(6/5) --> 2 * 2^(1/5)
// 12^(3/5) --> (2^2 * 3)^(3/5) --> 2*2^(1/5) * 3^(3/5)
//

// @ExprSymbol Power
// @ExprAttributes  OneIdentity NumericFunction
// TODO: Error handling
//

// @ExprPattern (_, 1)
func PowerXOne(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0]
}

// @ExprPattern (1, _)
func PowerOneX(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0]
}

// @ExprPattern (1.0,_)
func PowerOneRealX(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0]
}

// @ExprPattern (_Number, -1.0)
func PowerNumberInvReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return PowerNumberInv(e, c, args)
}

// @ExprPattern (_Number, -1)
func PowerNumberInv(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	fmt.Println("IN POWER -1")
	arg := args[0].(core.Number)
	fmt.Println("In power number inv: ", arg)
	if arg.Sign() == 0 {

		return core.NewError("DivisionByZero", "Division by zero")
	}
	result := arg.AsInv()
	fmt.Println(result)
	return arg.AsInv()
}

// @ExprPattern (_Integer, _Integer)
func PowerInteger(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := args[0].(core.Integer)
	y := args[1].(core.Integer)

	if x.Sign() == 0 && y.Sign() == -1 {
		return core.NewError("DivisionByZero", "Division by zero")
	}

	switch y.Sign() {
	case 0:
		return core.NewInteger(1)
	case -1:
		// x ^ -y == 1/ (x^y)
		return core.PowerInteger(x, y.AsNeg().(core.Integer)).AsInv()
	default:
		return core.PowerInteger(x, y)
	}
}

// @ExprPattern (_Rational, _Integer)
func PowerRatInt(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := args[0].(core.Rational)
	n := args[1].(core.Integer)

	if x.Sign() == 0 && n.Sign() == -1 {
		return core.NewError("DivisionByZero", "Division by zero")
	}

	// (x/y)^n = x^n/y^n = x^n * y^-n
	return core.ListFrom(symbol.Times,
		core.ListFrom(symbol.Power, x.AsNum(), n),
		core.ListFrom(symbol.Power, x.AsDenom(), n.AsNeg()),
	)

}

// @ExprPattern (_Integer,_Real)
func PowerIntToReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	base := args[0].(core.Integer)
	exp := args[1].(core.Real)

	if exp.IsFloat64() {
		result, err := core.PowerFloat64(base.Float64(), exp.Float64())
		if err != nil {
			return core.NewError("DivisionByZero", "Division by zero")
		}
		return core.NewReal(result)
	}

	expflt := exp.AsBigFloat()
	baseflt := core.ToBigFloat(new(big.Float).SetPrec(exp.Prec()), base)

	return new(big.Float).SetPrec(exp.Prec()).Pow(baseflt, expflt)
}

// @ExprPattern (_Real, _Integer)
func PowerRealToInt(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	base := args[0].(core.Real)
	exp := args[1].(core.Integer)

	if base.IsFloat64() {
		result, err := core.PowerFloat64(base.Float64(), exp.Float64())
		if err != nil {
			return core.NewError("DivisionByZero", "Division by zero")
		}
		return core.NewReal(result)
	}

	// base is BigFloat
	baseflt := base.AsBigFloat()
	prec := uint(baseflt.Prec())

	// set exponent default precision
	//   will get adjused by ToBigFloat
	expflt := new(big.Float).SetPrec(prec)
	core.ToBigFloat(expflt, exp)

	prec = core.UMin(prec, expflt.Prec())

	z := new(big.Float).SetPrec(prec).Pow(baseflt, expflt)
	return z

}

// PowerNumbers performs power operation on numeric arguments
// Returns (float64, error) for clear type safety
// TODO: Error handling
//
// @ExprPattern (_Real, _Real)
func PowerNumbers(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	base := args[0].(core.Real)
	exp := args[1].(core.Real)

	if base.IsFloat64() || exp.IsFloat64() {
		result, err := core.PowerFloat64(base.Float64(), exp.Float64())
		if err != nil {
			return core.NewError("DivisionByZero", "Division by zero")
		}
		return core.NewReal(result)
	}

	// both are Big Floats
	b := base.AsBigFloat()
	return new(big.Float).SetPrec(b.Prec()).Pow(b, exp.AsBigFloat())
}

// (a^(x))^(y) = a^(x*y) if y is an integer only
// unclear on restriction
// but (x ^2.0) ^ 3.0 does not simplify in MMA
//
// @ExprPattern (Power(_, _), _Integer)
func PowerPowerInt(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	pow := args[0].(core.List).Tail()
	a := pow[0]
	x := pow[1]
	y := args[1]
	return core.ListFrom(symbol.Power, a, core.ListFrom(symbol.Times, x, y))
}
