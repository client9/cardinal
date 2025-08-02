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

/*
// evaluateReplace implements enhanced Replace that supports both Rule and RuleDelayed
func (e *Evaluator) evaluateReplace(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError", "Replace requires exactly 2 arguments", args)
	}

	// Evaluate the expression to be replaced
	expr := e.evaluate(args[0], ctx)
	if core.IsError(expr) {
		return expr
	}

	// Evaluate the rule (but RuleDelayed RHS will remain unevaluated due to HoldRest)
	rule := e.evaluate(args[1], ctx)
	if core.IsError(rule) {
		return rule
	}

	// Handle single rule
	if e.isRuleOrRuleDelayed(rule) {
		result := e.applyRuleDelayedAware(expr, rule, ctx)
		// Re-evaluate the result to handle expressions like Plus(2, 2) -> 4
		return e.evaluate(result, ctx)
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
						// Rule matched and produced a different result
						// Re-evaluate the result to handle expressions like Plus(2, 2) -> 4
						return e.Evaluate(c, result)
					}
				}
		}
	}

	// No rule matched or invalid rule format
	return expr
}
*/

// Replace,  supports both Rule and RuleDelayed
func Replace(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError", "Replace requires exactly 2 arguments", args)
	}

	// Evaluate the expression to be replaced
	expr := e.Evaluate(c, args[0])
	if core.IsError(expr) {
		return expr
	}

	// Evaluate the rule (but RuleDelayed RHS will remain unevaluated due to HoldRest)
	rule := e.Evaluate(c, args[1])
	if core.IsError(rule) {
		return rule
	}

	// Handle single rule
	if isRuleOrRuleDelayed(rule) {
		result := applyRuleDelayedAware(e, c, expr, rule)
		// Re-evaluate the result to handle expressions like Plus(2, 2) -> 4
		return e.Evaluate(c, result)
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
					// Rule matched and produced a different result
					// Re-evaluate the result to handle expressions like Plus(2, 2) -> 4
					return e.Evaluate(c, result)
				}
			}
		}
	}

	// No rule matched or invalid rule format
	return expr
}

/*
// ReplaceWithRuleDelayedExpr applies a single rule (Rule or RuleDelayed) to an expression with evaluator access
func ReplaceWithRuleDelayedExpr(e *engine.Evaluator, c *engine.Context,expr core.Expr, rule core.Expr) core.Expr {
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

// ReplaceExpr applies rules to expressions: Replace(expr, rule)
func ReplaceExpr(e *engine.Evaluator, c *engine.Context,  expr, rule core.Expr) core.Expr {
	// Try to apply the rule to the expression (top-level only)
	return ReplaceWithRuleDelayedExpr(e, c, expr, rule)
}
// ReplaceAllExpr applies rules recursively: ReplaceAll(expr, rule)
func ReplaceAllExpr(e *engine.Evaluator, c *engine.Context, expr, rule core.Expr) core.Expr {
	// Check if rule is a List of Rules
	if ruleList, ok := rule.(core.List); ok && len(ruleList.Elements) > 1 {
		head := ruleList.Elements[0]
		if symbolName, ok := core.ExtractSymbol(head); ok && symbolName == "List" {
			// Validate that all elements are Rules
			for _, r := range ruleList.Elements[1:] {
				if rList, ok := r.(core.List); ok && len(rList.Elements) == 3 {
					if ruleHead, ok := core.ExtractSymbol(rList.Elements[0]); ok {
						if ruleHead != "Rule" && ruleHead != "RuleDelayed" {
							// Non-Rule element found, return unevaluated ReplaceAll
							return core.List{Elements: []core.Expr{
								core.NewSymbol("ReplaceAll"),
								expr,
								rule,
							}}
						}
					} else {
						// Non-Rule element found, return unevaluated ReplaceAll
						return core.List{Elements: []core.Expr{
							core.NewSymbol("ReplaceAll"),
							expr,
							rule,
						}}
					}
				} else {
					// Non-Rule element found, return unevaluated ReplaceAll
					return core.List{Elements: []core.Expr{
						core.NewSymbol("ReplaceAll"),
						expr,
						rule,
					}}
				}
			}

			// All elements are Rules, apply each rule in the list
			result := expr
			for _, r := range ruleList.Elements[1:] {
				result = replaceAllRecursive(e, c, result, r)
			}
			return result
		}
	}

	// Apply single rule recursively to the expression and all subexpressions
	return replaceAllRecursive(e, c, expr, rule)
}

// replaceAllRecursive recursively applies rules to an expression and its subexpressions
func replaceAllRecursive(e *engine.Evaluator, c *engine.Context, expr, rule core.Expr) core.Expr {
	// First try to apply the rule to the current expression
	replaced := ReplaceWithRuleDelayedExpr(e, c, expr, rule)
	if !expressionsEqual(replaced, expr) {
		// Rule applied, recursively apply to the result
		return replaceAllRecursive(e,c, replaced, rule)
	}

	// No rule applied, try to apply to subexpressions
	if list, ok := expr.(core.List); ok && len(list.Elements) > 0 {
		modified := false
		newElements := make([]core.Expr, len(list.Elements))

		for i, element := range list.Elements {
			newElement := replaceAllRecursive(e, c, element, rule)
			newElements[i] = newElement
			if !expressionsEqual(newElement, element) {
				modified = true
			}
		}

		if modified {
			return core.List{Elements: newElements}
		}
	}

	// Return original expression if no changes
	return expr
}

*/
