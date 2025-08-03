package stdlib

import (
	"github.com/client9/sexpr/core"
)

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
