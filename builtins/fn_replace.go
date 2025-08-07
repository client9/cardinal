package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

func asRule(expr core.Expr) (a, b core.Expr, ok bool) {
	list, ok := expr.(core.List)
	if !ok {
		return nil, nil, false
	}
	head := list.Head()
	if head != "Rule" && head != "RuleDelayed" {
		return nil, nil, false
	}
	args := list.Tail()
	return args[0], args[1], true
}

// isRuleOrRuleDelayed checks if an expression is a Rule or RuleDelayed
func isRuleOrRuleDelayed(expr core.Expr) bool {

	//if symbolName, ok := core.ExtractSymbol(expr); ok {
	//	return symbolName == "Rule" || symbolName == "RuleDelayed"
	//}
	if expr.Length() != 2 {
		return false
	}
	head := expr.Head()
	return head == "Rule" || head == "RuleDelayed"
}

func isRuleList(expr core.Expr) bool {
	list, ok := expr.(core.List)
	if !ok {
		return false
	}
	for _, arg := range list.Tail() {
		if !isRuleOrRuleDelayed(arg) {
			return false
		}
	}
	return true
}

// applyRuleDelayedAware applies a rule (Rule or RuleDelayed) with proper handling for both types
func applyRuleDelayedAware(expr core.Expr, rule core.Expr) core.Expr {
	// Handle both Rule and RuleDelayed

	if pattern, replacement, ok := asRule(rule); ok {
		// Use pattern matching with variable binding
		if matches, bindings := core.MatchWithBindings(pattern, expr); matches {
			return core.SubstituteBindings(replacement, bindings)
			/*
				if rule.Head() == "Rule" {
					// For Rule, substitute directly (current behavior)
					return core.SubstituteBindings(replacement, bindings)
				} else {
					// For RuleDelayed, evaluate RHS in a context with bindings
					ruleCtx := engine.NewChildContext(c)

					// Add pattern variable bindings to the rule context
					for varName, value := range bindings {
						ruleCtx.AddScopedVar(varName) // Keep bindings local
						if err := ruleCtx.Set(varName, value); err != nil {
							return core.NewError("BindingError", err.Error())
						}
					}
					// Evaluate replacement in the rule context
					return e.Evaluate(ruleCtx, replacement)
				}
			*/
		}
	}
	/*
		if ruleDelayed, ok := rule.(core.RuleDelayedExpr); ok {
			// Handle direct RuleDelayedExpr
			matches, bindings := core.MatchWithBindings(ruleDelayed.Pattern, expr)
			if matches {
				// Create a new context with pattern variable bindings
				ruleCtx := engine.NewChildContext(c)

				// Add pattern variable bindings to the rule context
				for varName, value := range bindings {
					ruleCtx.AddScopedVar(varName) // Keep bindings local
					if err := ruleCtx.Set(varName, value); err != nil {
						return core.NewError("BindingError", err.Error())
					}
				}

				// Evaluate RHS in the rule context
				return e.Evaluate(ruleCtx, ruleDelayed.RHS)
			}
		}
	*/
	// If no match or invalid rule, return original expression
	return expr
}

// Replace,  supports both Rule and RuleDelayed
func Replace(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) != 2 {
		return core.NewError("ArgumentError", "Replace requires exactly 2 arguments")
	}

	expr := args[0]
	rule := args[1]
	// Handle single rule
	if isRuleOrRuleDelayed(rule) {
		return applyRuleDelayedAware(expr, rule)
	}

	// Handle List of rules
	if rulesList, ok := rule.(core.List); ok && rulesList.Length() > 0 {
		if rulesList.Head() == "List" {
			ruleSlice := rulesList.Tail()
			// First, check if ALL elements (except head) are Rules or RuleDelayed
			for _, r := range ruleSlice {
				if !isRuleOrRuleDelayed(r) {
					return core.NewError("ArgumentError", "Input was not a list of rules")
				}
			}

			// Only process as rule list if ALL elements are rules
			// Try each rule in order
			for _, ruleItem := range ruleSlice {
				result := applyRuleDelayedAware(expr, ruleItem)
				if !result.Equal(expr) {
					return result
				}
			}
		}
	}

	// No rule matched or invalid rule format
	return expr
}
