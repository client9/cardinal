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

// ReplaceExpr applies a single rule to an expression
// Replace(expr, Rule(pattern, replacement)) -> replacement if expr matches pattern, else expr
func ReplaceExpr(expr core.Expr, rule core.Expr) core.Expr {
	// Extract pattern and replacement from Rule(pattern, replacement)
	if ruleList, ok := rule.(core.List); ok && len(ruleList.Elements) == 3 {
		head := ruleList.Elements[0]
		if headSymbol, ok := head.(core.Symbol); ok && headSymbol.String() == "Rule" {
			pattern := ruleList.Elements[1]
			replacement := ruleList.Elements[2]

			// For now, implement exact matching only
			// TODO: Implement full pattern matching with variable binding
			if pattern.Equal(expr) {
				return replacement
			}
		}
	}

	// If no match or invalid rule, return original expression
	return expr
}
