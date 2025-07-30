package engine

import (
	"github.com/client9/sexpr/core"

	"fmt"
)

// AttributesExpr gets the attributes of a symbol
func AttributesExpr(expr core.Expr, ctx *Context) core.Expr {
	// The argument should be a symbol
	if symbolName, ok := core.ExtractSymbol(expr); ok {

		// Get the attributes from the symbol table
		attrs := ctx.symbolTable.Attributes(symbolName)

		// Convert attributes to a list of symbols
		attrElements := make([]core.Expr, len(attrs)+1)
		attrElements[0] = core.NewSymbol("List")

		for i, attr := range attrs {
			attrElements[i+1] = core.NewSymbol(attr.String())
		}

		return core.List{Elements: attrElements}
	}

	return core.NewErrorExpr("ArgumentError",
		"Attributes expects a symbol as argument", []core.Expr{expr})
}

// ReplaceWithRuleDelayed applies a single rule (Rule or RuleDelayed) to an expression with evaluator access
func ReplaceWithRuleDelayed(expr core.Expr, rule core.Expr, evaluator *Evaluator, ctx *Context) core.Expr {
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
						ruleCtx := NewChildContext(ctx)

						// Add pattern variable bindings to the rule context
						for varName, value := range bindings {
							ruleCtx.AddScopedVar(varName) // Keep bindings local
							if err := ruleCtx.Set(varName, value); err != nil {
								return core.NewErrorExpr("BindingError", err.Error(), []core.Expr{rule})
							}
						}

						// Evaluate replacement in the rule context
						return evaluator.evaluate(replacement, ruleCtx)
					}
				}
			}
		}
	} else if ruleDelayed, ok := rule.(core.RuleDelayedExpr); ok {
		// Handle direct RuleDelayedExpr
		matches, bindings := core.MatchWithBindings(ruleDelayed.Pattern, expr)
		if matches {
			// Create a new context with pattern variable bindings
			ruleCtx := NewChildContext(ctx)

			// Add pattern variable bindings to the rule context
			for varName, value := range bindings {
				ruleCtx.AddScopedVar(varName) // Keep bindings local
				if err := ruleCtx.Set(varName, value); err != nil {
					return core.NewErrorExpr("BindingError", err.Error(), []core.Expr{rule})
				}
			}

			// Evaluate RHS in the rule context
			return evaluator.evaluate(ruleDelayed.RHS, ruleCtx)
		}
	}

	// If no match or invalid rule, return original expression
	return expr
}

// WrapAttributesExpr is a clean wrapper for Attributes that uses the business logic function
func WrapAttributesExpr(args []core.Expr, ctx *Context) core.Expr {
	// Validate argument count
	if len(args) != 1 {
		return core.NewErrorExpr("ArgumentError",
			"Attributes expects 1 argument", args)
	}

	// Check for errors in arguments first
	if core.IsError(args[0]) {
		return args[0]
	}

	// Call business logic function
	return AttributesExpr(args[0], ctx)
}

// SetAttributesSingle sets a single attribute on a symbol
func SetAttributesSingle(symbol core.Expr, attr core.Expr, ctx *Context) core.Expr {
	// The first argument should be a symbol
	if symbolName, ok := core.ExtractSymbol(symbol); ok {

		// The second argument should be an attribute symbol
		if attrName, ok := core.ExtractSymbol(attr); ok {

			// Convert string to Attribute
			if attribute, ok := StringToAttribute(attrName); ok {
				// Set the attribute on the symbol
				ctx.symbolTable.SetAttributes(symbolName, []Attribute{attribute})
				return core.NewSymbolNull()
			}

			return core.NewErrorExpr("ArgumentError",
				fmt.Sprintf("Unknown attribute: %s", attrName), []core.Expr{attr})
		}
	}

	return core.NewErrorExpr("ArgumentError",
		"SetAttributes expects (symbol, attribute)", []core.Expr{symbol, attr})
}

// SetAttributesList sets multiple attributes on a symbol
func SetAttributesList(symbol core.Expr, attrList core.List, ctx *Context) core.Expr {
	// The first argument should be a symbol
	if symbolName, ok := core.ExtractSymbol(symbol); ok {

		// Process each attribute in the list (skip head at index 0)
		var attributes []Attribute
		for i := 1; i < len(attrList.Elements); i++ {
			attrExpr := attrList.Elements[i]

			if attrName, ok := core.ExtractSymbol(attrExpr); ok {

				// Convert string to Attribute
				if attribute, ok := StringToAttribute(attrName); ok {
					attributes = append(attributes, attribute)
				} else {
					return core.NewErrorExpr("ArgumentError",
						fmt.Sprintf("Unknown attribute: %s", attrName), []core.Expr{attrExpr})
				}
			} else {
				return core.NewErrorExpr("ArgumentError",
					"Attributes list must contain symbols", []core.Expr{attrExpr})
			}
		}

		// Set all attributes on the symbol
		ctx.symbolTable.SetAttributes(symbolName, attributes)
		return core.NewSymbolNull()
	}

	return core.NewErrorExpr("ArgumentError",
		"SetAttributes expects (symbol, attribute list)", []core.Expr{symbol, attrList})
}

// ClearAttributesSingle clears a single attribute from a symbol
func ClearAttributesSingle(symbol core.Expr, attr core.Expr, ctx *Context) core.Expr {
	// The first argument should be a symbol
	if symbolName, ok := core.ExtractSymbol(symbol); ok {

		// The second argument should be an attribute symbol
		if attrName, ok := core.ExtractSymbol(attr); ok {

			// Convert string to Attribute
			if attribute, ok := StringToAttribute(attrName); ok {
				// Clear the attribute from the symbol
				ctx.symbolTable.ClearAttributes(symbolName, []Attribute{attribute})
				return core.NewSymbolNull()
			}

			return core.NewErrorExpr("ArgumentError",
				fmt.Sprintf("Unknown attribute: %s", attrName), []core.Expr{attr})
		}
	}

	return core.NewErrorExpr("ArgumentError",
		"ClearAttributes expects (symbol, attribute)", []core.Expr{symbol, attr})
}

// ClearAttributesList clears multiple attributes from a symbol
func ClearAttributesList(symbol core.Expr, attrList core.List, ctx *Context) core.Expr {
	// The first argument should be a symbol
	if symbolName, ok := core.ExtractSymbol(symbol); ok {

		// Process each attribute in the list (skip head at index 0)
		for i := 1; i < len(attrList.Elements); i++ {
			attrExpr := attrList.Elements[i]

			if attrName, ok := core.ExtractSymbol(attrExpr); ok {

				// Convert string to Attribute
				if attribute, ok := StringToAttribute(attrName); ok {
					// Clear this attribute from the symbol
					ctx.symbolTable.ClearAttributes(symbolName, []Attribute{attribute})
				} else {
					return core.NewErrorExpr("ArgumentError",
						fmt.Sprintf("Unknown attribute: %s", attrName), []core.Expr{attrExpr})
				}
			} else {
				return core.NewErrorExpr("ArgumentError",
					"Attributes list must contain symbols", []core.Expr{attrExpr})
			}
		}

		return core.NewSymbolNull()
	}

	return core.NewErrorExpr("ArgumentError",
		"ClearAttributes expects (symbol, attribute list)", []core.Expr{symbol, attrList})
}

// Clean wrappers for SetAttributes and ClearAttributes

// WrapSetAttributesSingle is a clean wrapper for SetAttributes with single attribute
func WrapSetAttributesSingle(args []core.Expr, ctx *Context) core.Expr {
	// Validate argument count
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError",
			"SetAttributes expects 2 arguments", args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if core.IsError(arg) {
			return arg
		}
	}

	// Call business logic function
	return SetAttributesSingle(args[0], args[1], ctx)
}

// WrapSetAttributesList is a clean wrapper for SetAttributes with attribute list
func WrapSetAttributesList(args []core.Expr, ctx *Context) core.Expr {
	// Validate argument count
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError",
			"SetAttributes expects 2 arguments", args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if core.IsError(arg) {
			return arg
		}
	}

	// Extract list from second argument
	attrList, ok := args[1].(core.List)
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"Second argument must be a list", args)
	}

	// Call business logic function
	return SetAttributesList(args[0], attrList, ctx)
}

// WrapClearAttributesSingle is a clean wrapper for ClearAttributes with single attribute
func WrapClearAttributesSingle(args []core.Expr, ctx *Context) core.Expr {
	// Validate argument count
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError",
			"ClearAttributes expects 2 arguments", args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if core.IsError(arg) {
			return arg
		}
	}

	// Call business logic function
	return ClearAttributesSingle(args[0], args[1], ctx)
}

// WrapClearAttributesList is a clean wrapper for ClearAttributes with attribute list
func WrapClearAttributesList(args []core.Expr, ctx *Context) core.Expr {
	// Validate argument count
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError",
			"ClearAttributes expects 2 arguments", args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if core.IsError(arg) {
			return arg
		}
	}

	// Extract list from second argument
	attrList, ok := args[1].(core.List)
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"Second argument must be a list", args)
	}

	// Call business logic function
	return ClearAttributesList(args[0], attrList, ctx)
}

// PatternSpecificityExpr calculates the specificity of a pattern expression for debugging
func PatternSpecificityExpr(pattern core.Expr, ctx *Context) core.Expr {
	// Calculate specificity directly from the pattern expression
	specificity := core.GetPatternSpecificity(pattern)
	return core.NewInteger(int64(specificity))
}

// ShowPatternsExpr lists all registered patterns for a function name
func ShowPatternsExpr(functionName core.Expr, ctx *Context) core.Expr {
	if funcName, ok := core.ExtractSymbol(functionName); ok {

		// Get function definitions from the registry
		definitions := ctx.functionRegistry.GetFunctionDefinitions(funcName)
		if definitions == nil {
			return core.NewErrorExpr("ArgumentError",
				fmt.Sprintf("No patterns found for function: %s", funcName), []core.Expr{functionName})
		}

		// Create a list of pattern information
		elements := make([]core.Expr, len(definitions)+1)
		elements[0] = core.NewSymbol("List")

		for i, def := range definitions {
			// Create a rule showing pattern -> specificity
			patternStr := def.Pattern.String()
			specificityStr := fmt.Sprintf("%d", def.Specificity)

			ruleElements := []core.Expr{
				core.NewSymbol("Rule"),
				core.NewString(patternStr),
				core.NewString(specificityStr),
			}

			elements[i+1] = core.List{Elements: ruleElements}
		}

		return core.List{Elements: elements}
	}

	return core.NewErrorExpr("ArgumentError",
		"ShowPatterns expects a symbol", []core.Expr{functionName})
}

// WrapPatternSpecificity is a clean wrapper for PatternSpecificity
func WrapPatternSpecificity(args []core.Expr, ctx *Context) core.Expr {
	// Validate argument count
	if len(args) != 1 {
		return core.NewErrorExpr("ArgumentError",
			"PatternSpecificity expects 1 argument", args)
	}

	// Check for errors in arguments first
	if core.IsError(args[0]) {
		return args[0]
	}

	// Call business logic function
	return PatternSpecificityExpr(args[0], ctx)
}

// WrapShowPatterns is a clean wrapper for ShowPatterns
func WrapShowPatterns(args []core.Expr, ctx *Context) core.Expr {
	// Validate argument count
	if len(args) != 1 {
		return core.NewErrorExpr("ArgumentError",
			"ShowPatterns expects 1 argument", args)
	}

	// Check for errors in arguments first
	if core.IsError(args[0]) {
		return args[0]
	}

	// Call business logic function
	return ShowPatternsExpr(args[0], ctx)
}

// MapExpr applies a function to each element of a list
// Map(f, {a, b, c}) -> {f(a), f(b), f(c)}
func MapExpr(function core.Expr, list core.Expr, evaluator *Evaluator, ctx *Context) core.Expr {
	// Check if the second argument is a list
	listExpr, ok := list.(core.List)
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"Map expects a list as the second argument", []core.Expr{list})
	}

	// If the list is empty or only has a head, return it unchanged
	if len(listExpr.Elements) <= 1 {
		return listExpr
	}

	// Extract head and elements
	head := listExpr.Elements[0]
	elements := listExpr.Elements[1:]

	// Apply the function to each element
	resultElements := make([]core.Expr, len(elements)+1)
	resultElements[0] = head // Keep the same head

	for i, element := range elements {
		// Create function application: function(element)
		applicationElements := []core.Expr{function, element}
		application := core.List{Elements: applicationElements}

		// Evaluate the function application
		result := evaluator.evaluate(application, ctx)
		resultElements[i+1] = result
	}

	return core.List{Elements: resultElements}
}

// WrapMapExpr is a clean wrapper for Map
func WrapMapExpr(args []core.Expr, ctx *Context) core.Expr {
	// Validate argument count
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError",
			"Map expects 2 arguments", args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if core.IsError(arg) {
			return arg
		}
	}

	// Get evaluator from context - we need this for function application
	// For now, we'll create a new evaluator instance
	// TODO: This is not ideal, but we need access to evaluation
	evaluator := NewEvaluator()

	// Copy the current context state to the new evaluator
	evaluator.context = ctx

	// Call business logic function
	return MapExpr(args[0], args[1], evaluator, ctx)
}

// ApplyExpr applies a function to a list of arguments as separate parameters
// Apply(f, {a, b, c}) -> f(a, b, c)
func ApplyExpr(function core.Expr, list core.Expr, evaluator *Evaluator, ctx *Context) core.Expr {
	// Check if the second argument is a list
	listExpr, ok := list.(core.List)
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"Apply expects a list as the second argument", []core.Expr{list})
	}

	// If the list is empty or only has a head, apply function with no arguments
	if len(listExpr.Elements) <= 1 {
		// Create function application with no arguments: function()
		applicationElements := []core.Expr{function}
		application := core.List{Elements: applicationElements}

		// Evaluate the function application
		return evaluator.evaluate(application, ctx)
	}

	// Extract elements (skip the head)
	elements := listExpr.Elements[1:]

	// Create function application: function(arg1, arg2, arg3, ...)
	applicationElements := make([]core.Expr, len(elements)+1)
	applicationElements[0] = function
	copy(applicationElements[1:], elements)

	application := core.List{Elements: applicationElements}

	// Evaluate the function application
	return evaluator.evaluate(application, ctx)
}

// WrapApplyExpr is a clean wrapper for Apply
func WrapApplyExpr(args []core.Expr, ctx *Context) core.Expr {
	// Validate argument count
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError",
			"Apply expects 2 arguments", args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if core.IsError(arg) {
			return arg
		}
	}

	// Get evaluator from context - we need this for function application
	evaluator := NewEvaluator()

	// Copy the current context state to the new evaluator
	evaluator.context = ctx

	// Call business logic function
	return ApplyExpr(args[0], args[1], evaluator, ctx)
}
