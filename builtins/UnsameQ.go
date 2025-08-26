package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol UnsameQ

// @ExprPattern (_,_)
func UnsameQ(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := args[0]
	y := args[1]
	return core.NewBool(!x.Equal(y))
}
