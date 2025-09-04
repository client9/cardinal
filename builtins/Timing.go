package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
	"time"
)

// @ExprSymbol Timing
// @ExprAttributes HoldAll
//

// @ExprPattern (_)
func Timing(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	start := time.Now()
	result := EvaluateExpr(e, c, args)
	elapsed := time.Since(start)
	return core.NewList(symbol.List, result, core.NewReal(float64(elapsed)/1.0e9))
}
