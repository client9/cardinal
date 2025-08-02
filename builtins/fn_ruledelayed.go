package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// RuleDelayed creates delayed rules: RuleDelayed(lhs, rhs)
func RuleDelayed(e *engine.Evaluator, c *engine.Context, lhs, rhs core.Expr) core.Expr {
	// Create a RuleDelayed expression - the actual rule application happens elsewhere

	// TODO List constructor
	return core.List{Elements: []core.Expr{
		core.NewSymbol("RuleDelayed"),
		lhs,
		rhs,
	}}
}
