package stdlib

import (
	"github.com/client9/sexpr/core"
)

// NotExpr performs logical negation on boolean expressions
func NotExpr(expr core.Expr) core.Expr {
	// Check if the expression is a boolean value (True/False symbol)
	if core.IsBool(expr) {
		val, _ := core.ExtractBool(expr)
		return core.NewBool(!val)
	}

	// Return unchanged expression if not boolean (symbolic behavior)
	return core.NewList("Not", expr)
}
