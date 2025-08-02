package builtins

import (
	"math/rand/v2"
	"github.com/client9/sexpr/engine"
	"github.com/client9/sexpr/core"
)

func RandomReal(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	return core.NewReal(rand.Float64())
}
