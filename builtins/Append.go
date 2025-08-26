package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Append
// TODO: List only, or anything?

// ListAppends adds an expression to the end of List type
// @ExprPattern (_List, _)
func Append(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	list := args[0].(core.List)
	element := args[1]
	return list.Append(element)
}
