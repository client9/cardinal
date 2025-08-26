package builtins

import (
	"fmt"
	"math"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Power
// @ExprAttributes  OneIdentity NumericFunction
// TODO: Error handling
//
// PowerInteger - if exp >= 0, then it returns an integer, if exp < 0 returns the float value or error
//
// @ExprPattern (_Integer, _Integer)
func PowerInteger(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	base, _ := core.ExtractInt64(args[0])
	exp, _ := core.ExtractInt64(args[1])
	if exp == 0 {
		return core.NewInteger(1)
	}
	if exp < 0 {
		val, err := powerFloat64(float64(base), float64(exp))
		if err != nil {
			// TODO ERRORS
			return core.NewSymbolNull()
		}
		return core.NewReal(val)
	}
	val, err := powerFloat64(float64(base), float64(exp))
	if err != nil {
		// TODO ERROR
		return core.NewSymbolNull()
	}
	return core.NewInteger(int64(val))
}

// PowerNumbers performs power operation on numeric arguments
// Returns (float64, error) for clear type safety
// TODO: Error handling
//
// @ExprPattern (_Real, _Real)
func PowerNumbers(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	base, _ := core.ExtractFloat64(args[0])
	exp, _ := core.ExtractFloat64(args[1])

	result, err := powerFloat64(base, exp)
	if err != nil {
		// TODO ERROR
		return core.NewSymbolNull()
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
