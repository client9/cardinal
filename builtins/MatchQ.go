package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol MatchQ

// MatchQExprs checks if an expression matches a pattern (pure test, no variable binding)
// @ExprPattern (_,_)
func MatchQ(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	pattern := args[1]

	ok, _ := core.MatchWithBindings(expr, pattern)
	return core.NewBool(ok)
}
