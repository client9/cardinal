package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
	"math/rand/v2"
)

func rand01() float64 {
	return rand.Float64()
}

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
	out[0] = core.NewSymbol("List")
	for i := 1; i <= count; i++ {
		out[i] = core.NewReal(rand01()*(max-min) + min)
	}
	return core.NewListFromExprs(out...)
}

func RRealDefault(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(rand01())
}
func RRealMax(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	max, _ := core.GetNumericValue(args[0])
	return core.NewReal(rand01() * max)
}

func RRealMinMax(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	list, _ := args[0].(core.List)
	largs := list.Tail()
	min, _ := core.GetNumericValue(largs[0])
	max, _ := core.GetNumericValue(largs[1])
	return core.NewReal(rand01()*(max-min) + min)
}

func RRealMaxCount(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	max, _ := core.GetNumericValue(args[0])
	nval, _ := core.GetNumericValue(args[1])
	return minMaxCount(0.0, max, int64(nval))
}

func RRealMinMaxCount(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	list, _ := args[0].(core.List)
	largs := list.Tail()
	min, _ := core.GetNumericValue(largs[0])
	max, _ := core.GetNumericValue(largs[1])
	nval, _ := core.GetNumericValue(args[1])
	return minMaxCount(min, max, int64(nval))
}

func minMaxCount(min, max float64, n int64) core.Expr {
	out := make([]core.Expr, n+1)
	out[0] = core.NewSymbol("List")
	for i := int64(1); i <= n; i++ {
		out[i] = core.NewReal(rand01()*(max-min) + min)
	}
	return core.NewListFromExprs(out...)
}
