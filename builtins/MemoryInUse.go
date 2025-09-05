package builtins

import (
	"runtime"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol MemoryInUse
// @ExprAttributes Protected
//
//

// @ExprPattern ()
func MemoryInUse(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	return core.NewInteger(int64(m.HeapAlloc))
}
