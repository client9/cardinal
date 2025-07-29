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
		return core.NewBool(!val)
	}

	// Return unchanged expression if not boolean (symbolic behavior)
	return core.NewList("Not", expr)
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
		if symbolName, ok := core.ExtractSymbol(head); ok && symbolName == "Rule" {
			pattern := ruleList.Elements[1]
			replacement := ruleList.Elements[2]

			// Use pattern matching with variable binding
			matches, bindings := core.MatchWithBindings(pattern, expr)
			if matches {
				// If pattern matches, substitute variables in replacement and return it
				return core.SubstituteBindings(replacement, bindings)
			}
		}
	}

	// If no match or invalid rule, return original expression
	return expr
}

// ReplaceWithRules applies a list of rules to an expression
// Replace(expr, List(Rule1, Rule2, ...)) -> replacement from first matching rule, else expr
func ReplaceWithRules(expr core.Expr, rulesList core.Expr) core.Expr {
	// Extract List(Rule1, Rule2, ...)
	if list, ok := rulesList.(core.List); ok && len(list.Elements) >= 1 {
		head := list.Elements[0]
		if symbolName, ok := core.ExtractSymbol(head); ok && symbolName == "List" {
			// Iterate through each rule in order
			for i := 1; i < len(list.Elements); i++ {
				rule := list.Elements[i]

				// Validate that this element is actually a Rule
				if ruleExpr, ok := rule.(core.List); ok && len(ruleExpr.Elements) == 3 {
					if ruleName, ok := core.ExtractSymbol(ruleExpr.Elements[0]); ok && ruleName == "Rule" {
						// Try to apply this rule using existing ReplaceExpr logic
						result := ReplaceExpr(expr, rule)
						if !result.Equal(expr) {
							// Rule matched and produced a different result - return it
							return result
						}
					}
				}
				// If this element is not a valid Rule, continue to next element
			}
		}
	}

	// No rules matched or invalid list structure, return original expression
	return expr
}

// ReplaceAllExpr applies a single rule to all subexpressions recursively
// ReplaceAll(expr, Rule(pattern, replacement)) -> expr with all matching subexpressions replaced
func ReplaceAllExpr(expr core.Expr, rule core.Expr) core.Expr {
	// First try to apply the rule to the current expression
	result := ReplaceExpr(expr, rule)

	// If the rule matched at this level, we're done (don't recurse into replacement)
	if !result.Equal(expr) {
		return result
	}

	// If no match at this level, recursively apply to subexpressions
	if list, ok := expr.(core.List); ok && len(list.Elements) > 0 {
		// Create new list with transformed elements
		newElements := make([]core.Expr, len(list.Elements))
		changed := false

		for i, element := range list.Elements {
			newElement := ReplaceAllExpr(element, rule)
			newElements[i] = newElement
			if !newElement.Equal(element) {
				changed = true
			}
		}

		if changed {
			return core.NewListFromExprs(newElements...)
		}
	}

	// No changes made, return original expression
	return expr
}

// ReplaceAllWithRules applies a list of rules to all subexpressions recursively
// ReplaceAll(expr, List(Rule1, Rule2, ...)) -> expr with all matching subexpressions replaced
func ReplaceAllWithRules(expr core.Expr, rulesList core.Expr) core.Expr {
	// First try to apply rules to the current expression
	result := ReplaceWithRules(expr, rulesList)

	// If a rule matched at this level, we're done (don't recurse into replacement)
	if !result.Equal(expr) {
		return result
	}

	// If no match at this level, recursively apply to subexpressions
	if list, ok := expr.(core.List); ok && len(list.Elements) > 0 {
		// Create new list with transformed elements
		newElements := make([]core.Expr, len(list.Elements))
		changed := false

		for i, element := range list.Elements {
			newElement := ReplaceAllWithRules(element, rulesList)
			newElements[i] = newElement
			if !newElement.Equal(element) {
				changed = true
			}
		}

		if changed {
			return core.NewListFromExprs(newElements...)
		}
	}

	// No changes made, return original expression
	return expr
}
