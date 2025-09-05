package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Head

// Head returns the head/type of an expression
// Head(Foo(1,2,3)) is Foo (e.g. Symbol(Foo))
// Head(Head("Foo")) is Symbol("String")
//
// @ExprPattern (_)
func Head(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return args[0].Head()
}
