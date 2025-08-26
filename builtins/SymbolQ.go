package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol SymbolQ

// SymbolQExpr checks if an expression is a symbol
// @ExprPattern (_)
func SymbolQExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	return core.NewBool(core.IsSymbol(args[0]))
}
