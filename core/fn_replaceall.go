package core

// ReplaceAllExpr applies a single rule to all subexpressions recursively
// ReplaceAll(expr, Rule(pattern, replacement)) -> expr with all matching subexpressions replaced
func ReplaceAllExpr(expr Expr, rule Expr) Expr {
	// First try to apply the rule to the current expression
	result := ReplaceExpr(expr, rule)

	// If the rule matched at this level, we're done (don't recurse into replacement)
	if !result.Equal(expr) {
		return result
	}

	// If no match at this level, recursively apply to subexpressions
	if list, ok := expr.(List); ok {
		// Create new list with transformed elements
		newElements := make([]Expr, list.Length()+1)
		changed := false

		for i, element := range list.AsSlice() {
			newElement := ReplaceAllExpr(element, rule)
			newElements[i] = newElement
			if !newElement.Equal(element) {
				changed = true
			}
		}

		if changed {
			return NewListFromExprs(newElements...)
		}
	}

	// No changes made, return original expression
	return expr
}

// ReplaceAllWithRules applies a list of rules to all subexpressions recursively
// ReplaceAll(expr, List(Rule1, Rule2, ...)) -> expr with all matching subexpressions replaced
func ReplaceAllWithRules(expr Expr, rulesList Expr) Expr {
	// First try to apply rules to the current expression
	result := ReplaceWithRules(expr, rulesList)

	// If a rule matched at this level, we're done (don't recurse into replacement)
	if !result.Equal(expr) {
		return result
	}

	// If no match at this level, recursively apply to subexpressions
	if list, ok := expr.(List); ok {
		// Create new list with transformed elements
		newElements := make([]Expr, list.Length()+1)
		changed := false

		for i, element := range list.AsSlice() {
			newElement := ReplaceAllWithRules(element, rulesList)
			newElements[i] = newElement
			if !newElement.Equal(element) {
				changed = true
			}
		}

		if changed {
			return NewListFromExprs(newElements...)
		}
	}

	// No changes made, return original expression
	return expr
}
