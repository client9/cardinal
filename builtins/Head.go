package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Head

// HeadExpr returns the head/type of an expression
// expr.Head() returns a string, but Head returns a symbol
// Head(Head("foo")) is Symbol("String")
//
// @ExprPattern (_)
func HeadExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewSymbol(args[0].Head())
}
