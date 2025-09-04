package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Replace

func asRule(expr core.Expr) (a, b core.Expr, ok bool) {
	list, ok := expr.(core.List)
	if !ok {
		return nil, nil, false
	}
	head := list.Head()
	if head != symbol.Rule && head != symbol.RuleDelayed {
		return nil, nil, false
	}
	args := list.Tail()
	return args[0], args[1], true
}

// isRuleOrRuleDelayed checks if an expression is a Rule or RuleDelayed
func isRuleOrRuleDelayed(expr core.Expr) bool {
	if expr.Length() != 2 {
		return false
	}
	head := expr.Head()
	return head == symbol.Rule || head == symbol.RuleDelayed
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
		if matches, bindings := core.MatchWithBindings(expr, pattern); matches {
			return core.SubstituteBindings(replacement, bindings)
		}
	}

	// If no match or invalid rule, return original expression
	return expr
}

// Replace,  supports both Rule and RuleDelayed
//
// @ExprPattern (_,_)
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

	if !isRuleList(rule) {
		return core.NewError("ArgumentError", "Input was not a rule or list of rules")
	}

	// Handle List of rules
	ruleList, _ := rule.(core.List)
	ruleSlice := ruleList.Tail()

	// Only process as rule list if ALL elements are rules
	// Try each rule in order
	for _, ruleItem := range ruleSlice {
		result := applyRuleDelayedAware(expr, ruleItem)
		if !result.Equal(expr) {
			return result
		}
	}
	// No rule matched or invalid rule format
	return expr
}
