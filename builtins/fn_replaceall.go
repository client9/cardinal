package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// evaluateReplaceAll implements enhanced ReplaceAll that supports both Rule and RuleDelayed
func ReplaceAll(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) != 2 {
		return core.NewError("ArgumentError", "ReplaceAll requires exactly 2 arguments")
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
	if !result.Equal(expr) {
		return result
	}

	// If no changes were made and rule is a list with non-rules, return unevaluated ReplaceAll
	if rulesList, ok := rule.(core.List); ok && rulesList.Length() > 0 && rulesList.Head() == "List" {
		for _, r := range rulesList.Tail() {
			if !isRuleOrRuleDelayed(r) {
				return core.NewError("ArgumentError", "Input was not a list of rules")
			}
		}
	}
	return expr
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
	} else if rulesList, ok := rule.(core.List); ok && rulesList.Length() > 0 && rulesList.Head() == "List" {
		// Handle List of rules
		// First, check if ALL elements (except head) are Rules or RuleDelayed
		allAreRules := true
		rulesSlice := rulesList.Tail()
		for _, r := range rulesSlice {
			if !isRuleOrRuleDelayed(r) {
				allAreRules = false
				break
			}
		}

		// Only process as rule list if ALL elements are rules
		if allAreRules {
			// Try each rule in order
			for _, ruleItem := range rulesSlice {
				result = applyRuleDelayedAware(expr, ruleItem)
				if !result.Equal(expr) {
					// Rule matched at this level
					return result
				}
			}
		}
	}

	// No match at this level, recursively apply to subexpressions
	if list, ok := expr.(core.List); ok {
		// Create new list with transformed elements
		newElements := make([]core.Expr, list.Length()+1)
		changed := false

		for i, element := range list.AsSlice() {
			newElement := replaceAllRecursive(element, rule)
			newElements[i] = newElement
			if !newElement.Equal(element) {
				changed = true
			}
		}

		if changed {
			return core.NewListFromExprs(newElements...)
		}
	}

	// Return original expression if no changes
	return expr
}
