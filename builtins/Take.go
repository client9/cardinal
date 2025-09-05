package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Take

// TakeList extracts the first or last n elements from a list
// Take(expr, n) takes first n elements; Take(expr, -n) takes last n elements
// @ExprPattern (_, _Integer)
func TakeList(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	n, _ := core.ExtractInt64(args[1])
	return core.Take(expr, n)
}

// TakeListSingle takes the nth element from a list and returns it as a single-element list
// Take(expr, [n]) - returns List(element_n)
// TODO error when n= 0?
// @ExprPattern (_, List(_Integer))
func TakeListSingle(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	n, _ := core.ExtractInt64(args[1].(core.List).Tail()[0])

	element := core.Part(expr, n)

	// Wrap the element in a list with the same head as the original list
	return core.ListFrom(expr.Head(), element)
}

// TakeListRange takes a range of elements from a list
// Take(expr, [n, m]) - takes elements from index n to m (inclusive)
// @ExprPattern (_, List(_Integer,_Integer))
func TakeListRange(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	expr := args[0]
	list := args[1].(core.List)
	return core.TakeRange(expr, list)
}
