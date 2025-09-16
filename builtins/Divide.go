package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
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
	fmt.Println("IN DIVIDE")
	x, _ := core.ExtractInt64(args[0])
	y, _ := core.ExtractInt64(args[1])
	if y == 0 {
		return NewError("DivisionByZero", "Division by zero"),
	}
	if y == 1 {
		return core.NewInt(1)
	}

	tmp := core.NewRational(x,y)
	if tmp.IsInt() {
		return tmp.Num()
	}
	return tmp
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

// @ExprPattern (_Real, _Real)
func DivideReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := args[0].(core.Real)
	y := args[1].(core.Real)
	if y.Sign() == 0 {
		return core.NewError("DivisionByZero", "Division by zero")
	}
	return core.DivReal(x, y)
}

// @ExprPattern (_,_)
func DivideAny(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.ListFrom(symbol.Times, args[0], core.ListFrom(symbol.Power, args[1], core.NewInteger(-1)))
}
