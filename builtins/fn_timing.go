package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
	"time"
)

func Timing(e *engine.Evaluator, c *engine.Context, arg core.Expr) core.Expr {
	start := time.Now()
	result := EvaluateExpr(e, c, arg)
	elapsed := time.Since(start)
	return core.NewList("List", result, core.NewInteger(int64(elapsed)))
}
