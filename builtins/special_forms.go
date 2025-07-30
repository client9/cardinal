package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/stdlib"
)

// Interfaces are defined in functional.go

// Special form functions that require evaluator access

// IfExpr evaluates conditional expressions: If(condition, then) or If(condition, then, else)
func IfExpr(evaluator Evaluator, args ...core.Expr) core.Expr {
	if len(args) < 2 || len(args) > 3 {
		return core.NewErrorExpr("ArgumentError", "If expects 2 or 3 arguments", args)
	}

	condition := args[0]
	thenExpr := args[1]
	var elseExpr core.Expr = core.NewSymbol("Null")

	if len(args) == 3 {
		elseExpr = args[2]
	}

	// Evaluate the condition
	evalCond := evaluator.Evaluate(condition)

	// Check if the condition is true using stdlib function
	if stdlib.TrueQExpr(evalCond) {
		return evaluator.Evaluate(thenExpr)
	} else {
		return evaluator.Evaluate(elseExpr)
	}
}

// SetExpr evaluates immediate assignment: Set(lhs, rhs)
func SetExpr(evaluator Evaluator, lhs, rhs core.Expr) core.Expr {
	// Evaluate the right-hand side immediately
	evalRhs := evaluator.Evaluate(rhs)

	// Handle assignment to symbol
	if symbolName, ok := core.ExtractSymbol(lhs); ok {
		ctx := evaluator.GetContext()
		ctx.Set(symbolName, evalRhs)
		return evalRhs
	}

	return core.NewErrorExpr("SetError", "Invalid assignment target", []core.Expr{lhs})
}

// SetDelayedExpr evaluates delayed assignment: SetDelayed(lhs, rhs)
func SetDelayedExpr(evaluator Evaluator, lhs, rhs core.Expr) core.Expr {
	ctx := evaluator.GetContext()

	// Handle function definitions: f(x_) := body
	if list, ok := lhs.(core.List); ok && len(list.Elements) >= 1 {
		// This is a function definition
		headExpr := list.Elements[0]
		if _, ok := core.ExtractSymbol(headExpr); ok {
			// Get the function registry from context
			registry := ctx.GetFunctionRegistry()

			// Register the pattern with the function registry
			err := registry.RegisterUserFunction(lhs, rhs)

			if err != nil {
				return core.NewErrorExpr("DefinitionError", err.Error(), []core.Expr{lhs, rhs})
			}

			return core.NewSymbol("Null")
		}
	}

	// Handle simple variable assignment: x := value
	if symbolName, ok := core.ExtractSymbol(lhs); ok {
		// Store the right-hand side without evaluation (delayed)
		ctx.Set(symbolName, rhs)
		return core.NewSymbol("Null")
	}

	return core.NewErrorExpr("SetDelayedError", "Invalid assignment target", []core.Expr{lhs})
}

// HoldExpr prevents evaluation of its arguments: Hold(expr1, expr2, ...)
func HoldExpr(evaluator Evaluator, args ...core.Expr) core.Expr {
	// Create a Hold expression with all the unevaluated arguments
	elements := make([]core.Expr, len(args)+1)
	elements[0] = core.NewSymbol("Hold")
	copy(elements[1:], args)
	return core.List{Elements: elements}
}

// EvaluateExpr forces evaluation: Evaluate(expr)
func EvaluateExpr(evaluator Evaluator, expr core.Expr) core.Expr {
	return evaluator.Evaluate(expr)
}

// CompoundExpressionExpr evaluates multiple expressions: CompoundExpression(expr1, expr2, ...)
func CompoundExpressionExpr(evaluator Evaluator, args ...core.Expr) core.Expr {
	var result core.Expr = core.NewSymbol("Null")

	// Evaluate each expression in sequence, return the last result
	for _, arg := range args {
		result = evaluator.Evaluate(arg)
	}

	return result
}

// AndExpr evaluates logical AND with short-circuiting: And(expr1, expr2, ...)
func AndExpr(evaluator Evaluator, args ...core.Expr) core.Expr {
	var unevaluatedArgs []core.Expr

	// Short-circuit evaluation: stop at first false, collect non-boolean true values
	for _, arg := range args {
		result := evaluator.Evaluate(arg)

		// Check if it's explicitly False
		if symbolName, ok := core.ExtractSymbol(result); ok && symbolName == "False" {
			return core.NewSymbol("False")
		}

		// Check if it's explicitly True - continue without adding to unevaluated
		if symbolName, ok := core.ExtractSymbol(result); ok && symbolName == "True" {
			continue
		}

		// For non-boolean values, collect them
		unevaluatedArgs = append(unevaluatedArgs, result)
	}

	// If no unevaluated args remain, all were True
	if len(unevaluatedArgs) == 0 {
		return core.NewSymbol("True")
	}

	// If only one arg remains, return it directly
	if len(unevaluatedArgs) == 1 {
		return unevaluatedArgs[0]
	}

	// Return And expression with remaining args
	elements := make([]core.Expr, len(unevaluatedArgs)+1)
	elements[0] = core.NewSymbol("And")
	copy(elements[1:], unevaluatedArgs)
	return core.List{Elements: elements}
}

// OrExpr evaluates logical OR with short-circuiting: Or(expr1, expr2, ...)
func OrExpr(evaluator Evaluator, args ...core.Expr) core.Expr {
	var nonFalseArgs []core.Expr

	// Evaluate arguments and collect non-False values
	for _, arg := range args {
		result := evaluator.Evaluate(arg)

		// Check if it's explicitly True - short-circuit
		if symbolName, ok := core.ExtractSymbol(result); ok && symbolName == "True" {
			return core.NewSymbol("True")
		}

		// Skip False values
		if symbolName, ok := core.ExtractSymbol(result); ok && symbolName == "False" {
			continue
		}

		// Collect non-False values
		nonFalseArgs = append(nonFalseArgs, result)
	}

	// If no non-False args remain, all were False
	if len(nonFalseArgs) == 0 {
		return core.NewSymbol("False")
	}

	// If only one non-False arg, return it directly
	if len(nonFalseArgs) == 1 {
		return nonFalseArgs[0]
	}

	// Return Or expression with remaining non-False args
	elements := make([]core.Expr, len(nonFalseArgs)+1)
	elements[0] = core.NewSymbol("Or")
	copy(elements[1:], nonFalseArgs)
	return core.List{Elements: elements}
}

// BlockExpr creates a lexical scope: Block({vars}, body)
func BlockExpr(evaluator Evaluator, vars core.Expr, body core.Expr) core.Expr {
	ctx := evaluator.GetContext()

	// Store original variable values to restore later
	savedVars := make(map[string]core.Expr)

	// Parse variable assignments from vars (should be a List)
	if varList, ok := vars.(core.List); ok && len(varList.Elements) > 0 {
		// Expect List(Set(x, value), Set(y, value), ...)
		for i := 1; i < len(varList.Elements); i++ {
			assignment := varList.Elements[i]
			if assignList, ok := assignment.(core.List); ok && len(assignList.Elements) == 3 {
				if symbolName, ok := core.ExtractSymbol(assignList.Elements[0]); ok && symbolName == "Set" {
					if varName, ok := core.ExtractSymbol(assignList.Elements[1]); ok {
						// Save old value
						if oldValue, exists := ctx.Get(varName); exists {
							savedVars[varName] = oldValue
						}
						// Set new value (evaluate the RHS)
						newValue := evaluator.Evaluate(assignList.Elements[2])
						ctx.Set(varName, newValue)
					}
				}
			}
		}
	}

	// Evaluate the body in the modified context
	result := evaluator.Evaluate(body)

	// Restore original values
	for varName, oldValue := range savedVars {
		ctx.Set(varName, oldValue)
	}

	// For variables that didn't exist before, remove them
	if varList, ok := vars.(core.List); ok && len(varList.Elements) > 0 {
		for i := 1; i < len(varList.Elements); i++ {
			assignment := varList.Elements[i]
			if assignList, ok := assignment.(core.List); ok && len(assignList.Elements) == 3 {
				if symbolName, ok := core.ExtractSymbol(assignList.Elements[0]); ok && symbolName == "Set" {
					if varName, ok := core.ExtractSymbol(assignList.Elements[1]); ok {
						if _, wasSaved := savedVars[varName]; !wasSaved {
							// This variable didn't exist before, remove it
							ctx.Delete(varName)
						}
					}
				}
			}
		}
	}

	return result
}

// TableExpr generates a table: Table(expr, {i, n}) or Table(expr, {i, start, end})
func TableExpr(evaluator Evaluator, expr core.Expr, iterator core.Expr) core.Expr {
	ctx := evaluator.GetContext()
	var results []core.Expr

	// Handle simple count: Table(expr, n)
	if intVal, ok := iterator.(core.Integer); ok {
		count := int(intVal)
		for i := 0; i < count; i++ {
			result := evaluator.Evaluate(expr)
			results = append(results, result)
		}
		elements := make([]core.Expr, len(results)+1)
		elements[0] = core.NewSymbol("List")
		copy(elements[1:], results)
		return core.List{Elements: elements}
	}

	// Handle iterator: List(i, n) or List(i, start, end)
	if iterList, ok := iterator.(core.List); ok && len(iterList.Elements) >= 2 {
		if symbolName, ok := core.ExtractSymbol(iterList.Elements[0]); ok && symbolName == "List" {
			if len(iterList.Elements) == 3 {
				// List(i, n) - iterate from 1 to n
				if varName, ok := core.ExtractSymbol(iterList.Elements[1]); ok {
					if countVal, ok := iterList.Elements[2].(core.Integer); ok {
						count := int(countVal)

						// Save original variable value
						var savedValue core.Expr
						var hadValue bool
						if oldValue, exists := ctx.Get(varName); exists {
							savedValue = oldValue
							hadValue = true
						}

						// Iterate from 1 to count
						for i := 1; i <= count; i++ {
							ctx.Set(varName, core.Integer(i))
							result := evaluator.Evaluate(expr)
							results = append(results, result)
						}

						// Restore original value
						if hadValue {
							ctx.Set(varName, savedValue)
						} else {
							ctx.Delete(varName)
						}
					}
				}
			} else if len(iterList.Elements) == 4 {
				// List(i, start, end) - iterate from start to end
				if varName, ok := core.ExtractSymbol(iterList.Elements[1]); ok {
					if startVal, ok := iterList.Elements[2].(core.Integer); ok {
						if endVal, ok := iterList.Elements[3].(core.Integer); ok {
							start := int(startVal)
							end := int(endVal)

							// Save original variable value
							var savedValue core.Expr
							var hadValue bool
							if oldValue, exists := ctx.Get(varName); exists {
								savedValue = oldValue
								hadValue = true
							}

							// Iterate from start to end
							for i := start; i <= end; i++ {
								ctx.Set(varName, core.Integer(i))
								result := evaluator.Evaluate(expr)
								results = append(results, result)
							}

							// Restore original value
							if hadValue {
								ctx.Set(varName, savedValue)
							} else {
								ctx.Delete(varName)
							}
						}
					}
				}
			}
		}
	}

	// Create List with results
	elements := make([]core.Expr, len(results)+1)
	elements[0] = core.NewSymbol("List")
	copy(elements[1:], results)
	return core.List{Elements: elements}
}

// DoExpr executes an expression multiple times: Do(expr, {i, n})
func DoExpr(evaluator Evaluator, expr core.Expr, iterator core.Expr) core.Expr {
	ctx := evaluator.GetContext()

	// Handle simple count: Do(expr, n)
	if intVal, ok := iterator.(core.Integer); ok {
		count := int(intVal)
		for i := 0; i < count; i++ {
			evaluator.Evaluate(expr)
		}
		return core.NewSymbol("Null")
	}

	// Handle iterator: List(i, n) or List(i, start, end)
	if iterList, ok := iterator.(core.List); ok && len(iterList.Elements) >= 2 {
		if symbolName, ok := core.ExtractSymbol(iterList.Elements[0]); ok && symbolName == "List" {
			if len(iterList.Elements) == 3 {
				// List(i, n) - iterate from 1 to n
				if varName, ok := core.ExtractSymbol(iterList.Elements[1]); ok {
					if countVal, ok := iterList.Elements[2].(core.Integer); ok {
						count := int(countVal)

						// Save original variable value
						var savedValue core.Expr
						var hadValue bool
						if oldValue, exists := ctx.Get(varName); exists {
							savedValue = oldValue
							hadValue = true
						}

						// Iterate from 1 to count
						for i := 1; i <= count; i++ {
							ctx.Set(varName, core.Integer(i))
							evaluator.Evaluate(expr)
						}

						// Restore original value
						if hadValue {
							ctx.Set(varName, savedValue)
						} else {
							ctx.Delete(varName)
						}
					}
				}
			} else if len(iterList.Elements) == 4 {
				// List(i, start, end) - iterate from start to end
				if varName, ok := core.ExtractSymbol(iterList.Elements[1]); ok {
					if startVal, ok := iterList.Elements[2].(core.Integer); ok {
						if endVal, ok := iterList.Elements[3].(core.Integer); ok {
							start := int(startVal)
							end := int(endVal)

							// Save original variable value
							var savedValue core.Expr
							var hadValue bool
							if oldValue, exists := ctx.Get(varName); exists {
								savedValue = oldValue
								hadValue = true
							}

							// Iterate from start to end
							for i := start; i <= end; i++ {
								ctx.Set(varName, core.Integer(i))
								evaluator.Evaluate(expr)
							}

							// Restore original value
							if hadValue {
								ctx.Set(varName, savedValue)
							} else {
								ctx.Delete(varName)
							}
						}
					}
				}
			}
		}
	}

	return core.NewSymbol("Null")
}

// FunctionExpr creates anonymous functions: Function({vars}, body) or Function(body)
func FunctionExpr(evaluator Evaluator, args ...core.Expr) core.Expr {
	// TODO: Implement proper function creation with lexical scoping
	// This is complex and requires function object representation
	return core.NewErrorExpr("NotImplemented", "Function not yet implemented in builtins", []core.Expr{})
}

// RuleDelayedExpr creates delayed rules: RuleDelayed(lhs, rhs)
func RuleDelayedExpr(evaluator Evaluator, lhs, rhs core.Expr) core.Expr {
	// Create a RuleDelayed expression - the actual rule application happens elsewhere
	return core.List{Elements: []core.Expr{
		core.NewSymbol("RuleDelayed"),
		lhs,
		rhs,
	}}
}

// ReplaceExpr applies rules to expressions: Replace(expr, rule)
func ReplaceExpr(evaluator Evaluator, expr, rule core.Expr) core.Expr {
	// Try to apply the rule to the expression (top-level only)
	return ReplaceWithRuleDelayedExpr(evaluator, expr, rule)
}

// ReplaceAllExpr applies rules recursively: ReplaceAll(expr, rule)
func ReplaceAllExpr(evaluator Evaluator, expr, rule core.Expr) core.Expr {
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
				result = replaceAllRecursive(evaluator, result, r)
			}
			return result
		}
	}

	// Apply single rule recursively to the expression and all subexpressions
	return replaceAllRecursive(evaluator, expr, rule)
}

// replaceAllRecursive recursively applies rules to an expression and its subexpressions
func replaceAllRecursive(evaluator Evaluator, expr, rule core.Expr) core.Expr {
	// First try to apply the rule to the current expression
	replaced := ReplaceWithRuleDelayedExpr(evaluator, expr, rule)
	if !expressionsEqual(replaced, expr) {
		// Rule applied, recursively apply to the result
		return replaceAllRecursive(evaluator, replaced, rule)
	}

	// No rule applied, try to apply to subexpressions
	if list, ok := expr.(core.List); ok && len(list.Elements) > 0 {
		modified := false
		newElements := make([]core.Expr, len(list.Elements))

		for i, element := range list.Elements {
			newElement := replaceAllRecursive(evaluator, element, rule)
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

// expressionsEqual compares two expressions for equality
func expressionsEqual(a, b core.Expr) bool {
	return a.String() == b.String()
}
