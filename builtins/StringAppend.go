package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol StringAppend

// StringAppend appends a string to another string.
// Note: Append normally adds an element of list to a list
//
//	We don't have "character" or "rune" type, and characters
//	are single character *strings*.
//
// So while the intention of Append is Append("foo", "d")
// it can be used for string joining Append("foo", "bar")
//
// This is added for the principle of least surprise.
//
// @ExprPattern (_String, _String)
func StringAppend(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	lhs, _ := core.ExtractString(args[0])
	rhs, _ := core.ExtractString(args[1])
	return core.NewString(lhs + rhs)
}
