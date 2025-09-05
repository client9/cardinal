package core

import (
	"github.com/client9/cardinal/core/symbol"
)

// ReplaceExpr applies a single rule to an expression
// Replace(expr, Rule(pattern, replacement)) -> replacement if expr matches pattern, else expr
func ReplaceExpr(expr Expr, rule Expr) Expr {
	// Extract pattern and replacement from Rule(pattern, replacement)
	if ruleList, ok := rule.(List); ok && ruleList.Length() == 2 && ruleList.Head() == symbol.Rule {
		e := ruleList.Tail()
		if matches, bindings := MatchWithBindings(expr, e[0]); matches {
			return SubstituteBindings(e[1], bindings)
		}
	}

	// If no match or invalid rule, return original expression
	return expr
}

// ReplaceWithRules applies a list of rules to an expression
// Replace(expr, List(Rule1, Rule2, ...)) -> replacement from first matching rule, else expr
func ReplaceWithRules(expr Expr, rulesList Expr) Expr {
	list, ok := rulesList.(List)
	if !ok {
		return expr
	}
	if list.Head() != symbol.List {
		return expr
	}
	// Iterate through each rule in order
	for _, rule := range list.Tail() {
		if ruleExpr, ok := rule.(List); ok && ruleExpr.Length() == 2 && ruleExpr.Head() == symbol.Rule {
			// Try to apply this rule using existing ReplaceExpr logic
			result := ReplaceExpr(expr, rule)
			if !result.Equal(expr) {
				// Rule matched and produced a different result - return it
				return result
			}
		}
		// If this element is not a valid Rule, continue to next element
	}

	// No rules matched or invalid list structure, return original expression
	return expr
}
