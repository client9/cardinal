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

// MatchQExprs checks if an expression matches a pattern (pure test, no variable binding)
func MatchQExprs(expr, pattern core.Expr) bool {
	// Use the pure pattern matcher from core (no Context needed for pure testing)
	matcher := core.NewPatternMatcher()
	return matcher.TestMatch(pattern, expr)
}
