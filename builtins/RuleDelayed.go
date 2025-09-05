package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol RuleDelayed
// @ExprAttributes HoldRest

// RuleDelayed creates delayed rules: RuleDelayed(lhs, rhs)
// TODO: this seems wrong
//
// @ExprPattern (_,_)
func RuleDelayed(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	// Create a RuleDelayed expression - the actual rule application happens elsewhere
	return core.ListFrom(symbol.RuleDelayed, args[0], args[1])

}
