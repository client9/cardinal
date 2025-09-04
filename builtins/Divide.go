package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Divide
// TODO: symbolic  a/b --> a * b^-1
// TODO: fixed error type
/*
//
// DivideIntegers performs integer division on int64 arguments
// Returns (int64, error) for clear type safety using Go's integer division
//
// @ExprPattern (_Integer, _Integer)
func DivideIntegers(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x, _ := core.ExtractInt64(args[0])
	y, _ := core.ExtractInt64(args[1])
	if y == 0 {
		return NewError("DivisionByZero", "Division by zero"),
	}

	return core.NewRational(x,y) //Integer(x / y)
}

// DivideNumbers performs division on numeric arguments
// Returns (float64, error) for clear type safety
//
// @ExprPattern (_Real, _Real)
func DivideReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x, _ := core.ExtractFloat64(args[0])
	y, _ := core.ExtractFloat64(args[1])
	if y == 0 {
		return NewError("DivisionByZero", "Division by zero"),
	}

	return core.NewReal(x / y)
}

// @ExprPattern (_Number, _Number)
func DivideNumber(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x, _ := core.GetNumericValue(args[0])
	y, _ := core.GetNumericValue(args[1])
	if y == 0 {
		return NewError("DivisionByZero", "Division by zero"),
	}

	return core.NewReal(x / y)
}
*/
// @ExprPattern (_,_)
func DivideAny(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.ListFrom(symbol.Times, args[0], core.ListFrom(symbol.Power, args[1], core.NewInteger(-1)))
}
