package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// isRuleOrRuleDelayed checks if an expression is a Rule or RuleDelayed
func isRuleOrRuleDelayed(expr core.Expr) bool {
	if symbolName, ok := core.ExtractSymbol(expr); ok {
		return symbolName == "Rule" || symbolName == "RuleDelayed"
	}
	if ruleList, ok := expr.(core.List); ok && len(ruleList.Elements) == 3 {
		head := ruleList.Elements[0]
		if symbolName, ok := core.ExtractSymbol(head); ok {
			return symbolName == "Rule" || symbolName == "RuleDelayed"
		}
	}

	return false
}

// applyRuleDelayedAware applies a rule (Rule or RuleDelayed) with proper handling for both types
func applyRuleDelayedAware(e *engine.Evaluator, c *engine.Context, expr core.Expr, rule core.Expr) core.Expr {
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
						// For RuleDelayed, evaluate RHS in a context with bindings
						ruleCtx := engine.NewChildContext(c)

						// Add pattern variable bindings to the rule context
						for varName, value := range bindings {
							ruleCtx.AddScopedVar(varName) // Keep bindings local
							if err := ruleCtx.Set(varName, value); err != nil {
								return core.NewErrorExpr("BindingError", err.Error(), []core.Expr{rule})
							}
						}

						// Evaluate replacement in the rule context
						return e.Evaluate(ruleCtx, replacement)
					}
				}
			}
		}
	} else if ruleDelayed, ok := rule.(core.RuleDelayedExpr); ok {
		// Handle direct RuleDelayedExpr
		matches, bindings := core.MatchWithBindings(ruleDelayed.Pattern, expr)
		if matches {
			// Create a new context with pattern variable bindings
			ruleCtx := engine.NewChildContext(c)

			// Add pattern variable bindings to the rule context
			for varName, value := range bindings {
				ruleCtx.AddScopedVar(varName) // Keep bindings local
				if err := ruleCtx.Set(varName, value); err != nil {
					return core.NewErrorExpr("BindingError", err.Error(), []core.Expr{rule})
				}
			}

			// Evaluate RHS in the rule context
			return e.Evaluate(ruleCtx, ruleDelayed.RHS)
		}
	}

	// If no match or invalid rule, return original expression
	return expr
}

// Replace,  supports both Rule and RuleDelayed
func Replace(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError", "Replace requires exactly 2 arguments", args)
	}

	expr := args[0]
	rule := args[1]
	// Handle single rule
	if isRuleOrRuleDelayed(rule) {
		return applyRuleDelayedAware(e, c, expr, rule)
	}

	// Handle List of rules
	if rulesList, ok := rule.(core.List); ok && len(rulesList.Elements) >= 1 {
		head := rulesList.Elements[0]
		if symbolName, ok := core.ExtractSymbol(head); ok && symbolName == "List" {
			// First, check if ALL elements (except head) are Rules or RuleDelayed
			for i := 1; i < len(rulesList.Elements); i++ {
				if !isRuleOrRuleDelayed(rulesList.Elements[i]) {
					return core.NewErrorExpr("ArgumentError", "Input was not a list of rules", args)
				}
			}

			// Only process as rule list if ALL elements are rules
			// Try each rule in order
			for i := 1; i < len(rulesList.Elements); i++ {
				ruleItem := rulesList.Elements[i]
				result := applyRuleDelayedAware(e, c, expr, ruleItem)
				if !result.Equal(expr) {
					return result
				}
			}
		}
	}

	// No rule matched or invalid rule format
	return expr
}
