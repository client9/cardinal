package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
	"math/rand/v2"
)

func RandomReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(rand.Float64())
}
