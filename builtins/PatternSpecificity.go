package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol PatternSpecificity

// PatternSpecificity calculates the specificity of a pattern expression for debugging
// TODO: could move to core and directly use
// @ExprPattern (_)
func PatternSpecificity(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	specificity := core.GetPatternSpecificity(args[0])
	return core.NewInteger(int64(specificity))
}
