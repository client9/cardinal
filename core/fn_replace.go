package core

// ReplaceExpr applies a single rule to an expression
// Replace(expr, Rule(pattern, replacement)) -> replacement if expr matches pattern, else expr
func ReplaceExpr(expr Expr, rule Expr) Expr {
	// Extract pattern and replacement from Rule(pattern, replacement)
	if ruleList, ok := rule.(List); ok && len(ruleList.Elements) == 3 {
		head := ruleList.Elements[0]
		if symbolName, ok := ExtractSymbol(head); ok && symbolName == "Rule" {
			pattern := ruleList.Elements[1]
			replacement := ruleList.Elements[2]

			// Use pattern matching with variable binding
			matches, bindings := MatchWithBindings(pattern, expr)
			if matches {
				// If pattern matches, substitute variables in replacement and return it
				return SubstituteBindings(replacement, bindings)
			}
		}
	}

	// If no match or invalid rule, return original expression
	return expr
}

// ReplaceWithRules applies a list of rules to an expression
// Replace(expr, List(Rule1, Rule2, ...)) -> replacement from first matching rule, else expr
func ReplaceWithRules(expr Expr, rulesList Expr) Expr {
	// Extract List(Rule1, Rule2, ...)
	if list, ok := rulesList.(List); ok && len(list.Elements) >= 1 {
		head := list.Elements[0]
		if symbolName, ok := ExtractSymbol(head); ok && symbolName == "List" {
			// Iterate through each rule in order
			for i := 1; i < len(list.Elements); i++ {
				rule := list.Elements[i]

				// Validate that this element is actually a Rule
				if ruleExpr, ok := rule.(List); ok && len(ruleExpr.Elements) == 3 {
					if ruleName, ok := ExtractSymbol(ruleExpr.Elements[0]); ok && ruleName == "Rule" {
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
