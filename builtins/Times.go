package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Times
// @ExprAttributes Flat Orderless OneIdentity NumericFunction

// TimesExpr performs multiplication with light simplification
// Combines all numeric values and leaves symbolic terms unchanged
// @ExprPattern (___)
func TimesExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) == 0 {
		return core.NewInteger(1) // Times() = 1
	}
	return core.TimesList(args)
}
