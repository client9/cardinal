package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol TrueQ

// TrueQ check is an expression is explicity True
// @ExprPattern (_)
func TrueQ(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewBool(args[0] == symbol.True)
}
