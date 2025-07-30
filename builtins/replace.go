package builtins

import (
	"github.com/client9/sexpr/core"
)

// ReplaceWithRuleDelayedExpr applies a single rule (Rule or RuleDelayed) to an expression with evaluator access
func ReplaceWithRuleDelayedExpr(evaluator Evaluator, expr core.Expr, rule core.Expr) core.Expr {
	// Handle both Rule and RuleDelayed
	if ruleList, ok := rule.(core.List); ok && len(ruleList.Elements) == 3 {
		head := ruleList.Elements[0]
		if symbolName, ok := core.ExtractSymbol(head); ok {
			if symbolName == "Rule" || symbolName == "RuleDelayed" {
				pattern := ruleList.Elements[1]
				replacement := ruleList.Elements[2]

				// Use pattern matching with variable binding
				matches, bindings := core.MatchWithBindings(pattern, expr)
				if matches {
					if symbolName == "Rule" {
						// For Rule, substitute directly (current behavior)
						return core.SubstituteBindings(replacement, bindings)
					} else {
						// For RuleDelayed, substitute bindings but don't evaluate yet
						// The evaluation happens when the result is used
						substituted := core.SubstituteBindings(replacement, bindings)
						return substituted
					}
				}
			}
		}
	} else if ruleDelayed, ok := rule.(core.RuleDelayedExpr); ok {
		// Handle direct RuleDelayedExpr if this type exists
		matches, bindings := core.MatchWithBindings(ruleDelayed.Pattern, expr)
		if matches {
			// TODO: Implement proper context-based evaluation for RuleDelayed
			// This requires creating child contexts and evaluating RHS with bindings
			return core.SubstituteBindings(ruleDelayed.RHS, bindings)
		}
	}

	// If no match or invalid rule, return original expression
	return expr
}
