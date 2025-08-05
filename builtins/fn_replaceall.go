package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// evaluateReplaceAll implements enhanced ReplaceAll that supports both Rule and RuleDelayed
func ReplaceAll(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError", "ReplaceAll requires exactly 2 arguments", args)
	}

	// Evaluate the expression to be replaced
	expr := e.Evaluate(args[0])
	if core.IsError(expr) {
		return expr
	}

	// Evaluate the rule (but RuleDelayed RHS will remain unevaluated due to HoldRest)
	rule := e.Evaluate(args[1])
	if core.IsError(rule) {
		return rule
	}

	result := replaceAllRecursive(expr, rule)
	// If no changes were made and rule is a list with non-rules, return unevaluated ReplaceAll
	if result.Equal(expr) {
		if rulesList, ok := rule.(core.List); ok && len(rulesList.Elements) >= 1 {
			head := rulesList.Elements[0]
			if symbolName, ok := core.ExtractSymbol(head); ok && symbolName == "List" {
				// Check if this list contains non-rules
				for i := 1; i < len(rulesList.Elements); i++ {
					if !isRuleOrRuleDelayed(rulesList.Elements[i]) {
						return core.NewErrorExpr("ArgumentError", "Input was not a list of rules", args)
					}
				}
			}
		}
	}
	return result
	// Re-evaluate the final result to handle expressions like Plus(2, 2) -> 4
	//return e.Evaluate(c, result)
}

// replaceAllRecursive recursively applies rules to all subexpressions
func replaceAllRecursive(expr core.Expr, rule core.Expr) core.Expr {
	// First try to apply the rule at this level
	var result core.Expr

	// Handle single rule
	if isRuleOrRuleDelayed(rule) {
		result = applyRuleDelayedAware(expr, rule)
		if !result.Equal(expr) {
			// Rule matched at this level, return the result (don't recurse into replacement)
			return result
		}
	} else if rulesList, ok := rule.(core.List); ok && len(rulesList.Elements) >= 1 {
		// Handle List of rules
		head := rulesList.Elements[0]
		if symbolName, ok := core.ExtractSymbol(head); ok && symbolName == "List" {
			// First, check if ALL elements (except head) are Rules or RuleDelayed
			allAreRules := true
			for i := 1; i < len(rulesList.Elements); i++ {
				if !isRuleOrRuleDelayed(rulesList.Elements[i]) {
					allAreRules = false
					break
				}
			}

			// Only process as rule list if ALL elements are rules
			if allAreRules {
				// Try each rule in order
				for i := 1; i < len(rulesList.Elements); i++ {
					ruleItem := rulesList.Elements[i]
					result = applyRuleDelayedAware(expr, ruleItem)
					if !result.Equal(expr) {
						// Rule matched at this level
						return result
					}
				}
			}
		}
	}

	// No match at this level, recursively apply to subexpressions
	if list, ok := expr.(core.List); ok && len(list.Elements) > 0 {
		// Create new list with transformed elements
		newElements := make([]core.Expr, len(list.Elements))
		changed := false

		for i, element := range list.Elements {
			newElement := replaceAllRecursive(element, rule)
			newElements[i] = newElement
			if !newElement.Equal(element) {
				changed = true
			}
		}

		if changed {
			return core.List{Elements: newElements}
		}
	}

	// Return original expression if no changes
	return expr
}
