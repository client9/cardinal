package builtins

import (
	"math"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol N
// @ExprAttributes Protected
//

// @ExprPattern (E)
func N_E(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(math.E)
}

// @ExprPattern (_Integer)
func N_Rational(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(args[0].(core.Integer).Float64())
}

// @ExprPattern (_Rational)
func N_Integer(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(args[0].(core.Rational).Float64())
}

// @ExprPattern (_)
func N(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return evalNum(e, c, args[0])
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
