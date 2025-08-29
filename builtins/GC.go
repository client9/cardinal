package builtins

import (
	"runtime"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol GC
// @ExprAttributes Protected
//
//

// @ExprPattern ()
func GC(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	runtime.GC()
	return symbol.Null
}
