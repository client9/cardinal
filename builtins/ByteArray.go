package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol ByteArray
// @ExprAttributes Protected

// ByteArrayString is constructor of a ByteArray from a string
// @ExprPattern (_String)
func ByteArrayString(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewByteArrayFromString(string(args[0].(core.String)))
}
