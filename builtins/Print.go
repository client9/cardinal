package builtins

import (
	"fmt"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Print

// Print outputs the expression and returns it unchanged
// This allows debugging intermediate values in compound statements
// @ExprPattern (_)
func Print(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	arg := args[0]
	fmt.Println(arg.String())
	return arg
}
