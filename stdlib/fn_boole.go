package stdlib

import (
	"github.com/client9/sexpr/core"
)

func Boole(expr core.Expr) core.Expr {
	if val, ok := core.ExtractBool(expr); ok {
		if val {
			return core.NewInteger(1)
		}
		return core.NewInteger(0)
	}

	// Return unchanged expression if not boolean (symbolic behavior)
	return core.NewList("Boole", expr)
}
