package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Assert
// TODO: direct atom check

// @ExprPattern (_)
func Assert(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	arg := args[0]
	result := e.Evaluate(arg)
	if symbolName, ok := core.ExtractSymbol(result); ok && symbolName == "True" {
		return result
	}
	return core.NewError("AssertionFailed", arg.InputForm())
}
