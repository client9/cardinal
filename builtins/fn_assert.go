package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
	"github.com/client9/sexpr/stdlib"
)

func Assert(e *engine.Evaluator, c *engine.Context, arg core.Expr) core.Expr {
	result := e.Evaluate(arg)
	istrue := stdlib.TrueQExpr(result)
	if istrue {
		return result
	}
	return core.NewErrorExpr("Assertion", arg.InputForm(), []core.Expr{arg})
}
