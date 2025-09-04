package builtins

import (
	"fmt"
	"math"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"

	"github.com/client9/sexpr/engine"
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

// @ExprPattern (_Integer, -1)
func PowerIntegerInv(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	arg := args[0].(core.Integer)
	if arg.Sign() == 0 {

		return core.NewError("DivisionByZero", "Division by zero")
	}
	return arg.Inv()
}

// @ExprPattern (_Rational, -1)
func PowerRationalInv(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	arg := args[0].(core.Rational)
	if arg.Sign() == 0 {

		return core.NewError("DivisionByZero", "Division by zero")
	}
	return arg.Inv()
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
		return core.PowerInteger(x, y.Neg()).Inv()
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
		core.ListFrom(symbol.Power, x.Numerator(), n),
		core.ListFrom(symbol.Power, x.Denominator(), n.Neg()),
	)

}

// PowerNumbers performs power operation on numeric arguments
// Returns (float64, error) for clear type safety
// TODO: Error handling
//
// @ExprPattern (_Number, _Number)
func PowerNumbers(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	base, _ := core.GetNumericValue(args[0])
	exp, _ := core.GetNumericValue(args[1])

	result, err := powerFloat64(base, exp)
	if err != nil {
		return core.NewError("DivisionByZero", "Division by zero")
	}
	return core.NewReal(result)
}

func powerFloat64(base, exp float64) (float64, error) {
	result := math.Pow(base, exp)

	// Check for invalid results (NaN, Inf)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, fmt.Errorf("MathematicalError")
	}

	return result, nil
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
