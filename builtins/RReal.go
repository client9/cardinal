package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/atom"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol RReal

// @ExprPattern ()
func RRealDefault(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(rand01())
}

// @ExprPattern (_Number)
func RRealMax(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	max, _ := core.GetNumericValue(args[0])
	return core.NewReal(rand01() * max)
}

// @ExprPattern (List(_Number, _Number))
func RRealMinMax(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	list, _ := args[0].(core.List)
	largs := list.Tail()
	min, _ := core.GetNumericValue(largs[0])
	max, _ := core.GetNumericValue(largs[1])
	return core.NewReal(rand01()*(max-min) + min)
}

// @ExprPattern (_Number, _Number)
func RRealMaxCount(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	list, _ := args[0].(core.List)
	largs := list.Tail()
	max, _ := core.GetNumericValue(largs[0])
	nval, _ := core.GetNumericValue(largs[1])
	return minMaxCount(0.0, max, int64(nval))
}

// @ExprPattern (List(_Number, _Number), _Number)
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
	out[0] = core.SymbolFor(atom.List)
	for i := int64(1); i <= n; i++ {
		out[i] = core.NewReal(rand01()*(max-min) + min)
	}
	return core.NewListFromExprs(out...)
}
