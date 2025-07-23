package stdlib

import (
	"github.com/client9/sexpr/core"
)

// Logical functions

// NotExpr performs logical negation on boolean expressions
func NotExpr(expr core.Expr) core.Expr {
	// Check if the expression is a boolean value (True/False symbol)
	if core.IsBool(expr) {
		val, _ := core.ExtractBool(expr)
		return core.NewBoolAtom(!val)
	}

	// Return unchanged expression if not boolean (symbolic behavior)
	return core.NewList(core.NewSymbolAtom("Not"), expr)
}

// MatchQExprs checks if an expression matches a pattern
// func MatchQExprs(expr, pattern Expr, ctx *Context) bool {
// 	// Convert string-based pattern to symbolic if needed
// 	symbolicPattern := convertToSymbolicPattern(pattern)

// 	// Create a temporary context for pattern matching (don't pollute original context)
// 	tempCtx := NewChildContext(ctx)

// 	// Use the existing pattern matching logic
// 	return matchPatternForMatchQ(symbolicPattern, expr, tempCtx)
// }
