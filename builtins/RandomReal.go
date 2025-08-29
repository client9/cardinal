package builtins

import (
	"math/rand/v2"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

func rand01() float64 {
	return rand.Float64()
}

// @ExprSymbol RandomReal

// @ExprPattern (___)
func RandomReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) == 0 {
		return core.NewReal(rand01())
	}

	count := 1
	min := 0.0
	max := 1.0

	if val, ok := core.GetNumericValue(args[0]); ok {
		max = val
	}
	if list, ok := args[0].(core.List); ok && list.Length() == 2 {
		argok := 0
		largs := list.Tail()
		if val, ok := core.GetNumericValue(largs[0]); ok {
			min = val
			argok += 1
		}
		if val, ok := core.GetNumericValue(largs[1]); ok {
			max = val
			argok += 1
		}
		if argok != 2 {
			return core.NewError("ArgumentError", "expected n or {m,n} as numeric values")
		}
	}
	if len(args) == 2 {
		val, ok := core.GetNumericValue(args[1])
		if !ok || val < 0.0 {
			return core.NewError("ArgumentError", "expected n or {m,n} as numeric values")
		}
		count = int(val)
	}

	if count == 1 {
		return core.NewReal(rand01()*(max-min) + min)
	}

	out := make([]core.Expr, count+1)
	out[0] = symbol.List
	for i := 1; i <= count; i++ {
		out[i] = core.NewReal(rand01()*(max-min) + min)
	}
	return core.NewListFromExprs(out...)
}
