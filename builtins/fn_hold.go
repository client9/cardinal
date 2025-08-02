package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)
/*
 // evaluateHold implements the Hold special form
func (e *Evaluator) evaluateHold(args []core.Expr, ctx *Context) core.Expr {
        // Hold returns its arguments unevaluated wrapped in Hold
        return core.NewList("Hold", args...)
}
*/
// HoldExpr prevents evaluation of its arguments: Hold(expr1, expr2, ...)
func HoldExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	// Hold returns its arguments unevaluated wrapped in Hold
	return core.NewList("Hold", args...)
/*
	// Create a Hold expression with all the unevaluated arguments
	elements := make([]core.Expr, len(args)+1)
	elements[0] = core.NewSymbol("Hold")
	copy(elements[1:], args)
	return core.List{Elements: elements}
*/
}
