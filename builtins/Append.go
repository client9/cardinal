package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
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

// @ExprPattern (_String, _String)
func AppendString(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	a := args[0].(core.String)
	b := args[1].(core.String)
	return core.NewString(string(a) + string(b))
}
