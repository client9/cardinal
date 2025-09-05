package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Rule
// @ExprAttributes HoldRest Protected
//
// Rule
// TODO: this seems wrong
//
// @ExprPattern (_,_)
func Rule(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if core.ExprHasNamedPattern(args[0]) {
		// Create a RuleDelayed expression - the actual rule application happens elsewhere
		return core.ListFrom(symbol.RuleDelayed, args[0], args[1])
	}

	// doesn't have a named pattern, we can evaluate immediately
	return core.ListFrom(symbol.Rule, args[0], e.Evaluate(args[1]))
}
