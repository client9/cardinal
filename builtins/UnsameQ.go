package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol UnsameQ

// @ExprPattern (_,_)
func UnsameQ(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	x := args[0]
	y := args[1]
	return core.NewBool(!x.Equal(y))
}
