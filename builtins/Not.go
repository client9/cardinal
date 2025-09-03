package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Not

// Not performs logical negation on boolean expressions
//
// @ExprPattern (_)
func Not(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]

	if val, ok := core.ExtractBool(expr); ok {
		return core.NewBool(!val)
	}

	// Return unchanged expression if not boolean (symbolic behavior)
	return core.ListFrom(symbol.Not, expr)
}
