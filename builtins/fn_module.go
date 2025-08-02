package builtins

import (
        "github.com/client9/sexpr/core"
        "github.com/client9/sexpr/engine"
)

// evaluateBlock implements the Block special form for dynamic scoping
// Block(List(vars...), body) temporarily changes variable values and evaluates body
func Module(e *engine.Evaluator, ctx *engine.Context, args []core.Expr) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError", "Module expects 2 arguments", args)
	}

	// First argument should be a list of variable specifications
	varSpec := args[0]
	body := args[1] // Don't evaluate body yet - Block has HoldAll

	// Parse variable specification
	varList, ok := varSpec.(core.List)
	if !ok {
		return core.NewErrorExpr("ArgumentError", "Block first argument must be a list of variables", []core.Expr{varSpec})
	}

	// Collect scoped variable names
	var scopedVarNames []string

	// Process variable specifications to extract scoped variable names
	for i := 1; i < len(varList.Elements); i++ { // Skip head element
		varExpr := varList.Elements[i]

		if varSymbol, ok := core.ExtractSymbol(varExpr); ok {
			// Simple variable: {x} - add to scoped variables
			scopedVarNames = append(scopedVarNames, varSymbol)

		} else if assignment, ok := varExpr.(core.List); ok && len(assignment.Elements) == 3 {
			// Assignment: {x = value} - add to scoped variables
			head := assignment.Elements[0]
			if headSymbol, ok := core.ExtractSymbol(head); ok && headSymbol == "Set" {
				if varSymbol, ok := core.ExtractSymbol(assignment.Elements[1]); ok {
					scopedVarNames = append(scopedVarNames, varSymbol)
				} else {
					return core.NewErrorExpr("ArgumentError", "Block variable assignment must use a symbol", []core.Expr{assignment.Elements[1]})
				}
			} else {
				return core.NewErrorExpr("ArgumentError", "Block variable specification must be Set expression", []core.Expr{varExpr})
			}
		} else {
			return core.NewErrorExpr("ArgumentError", "Block variable specification must be symbol or Set expression", []core.Expr{varExpr})
		}
	}

	// Create Block context with selective scoping
	blockCtx := engine.NewBlockContext(ctx, scopedVarNames)

	// Set initial values for Block variables
	for i := 1; i < len(varList.Elements); i++ { // Skip head element
		varExpr := varList.Elements[i]

		if _, ok := core.ExtractSymbol(varExpr); ok {
			// Simple variable: {x} - leave unset (will return undefined during lookup)
			// Don't set anything - the variable is scoped but has no value

		} else if assignment, ok := varExpr.(core.List); ok && len(assignment.Elements) == 3 {
			// Assignment: {x = value} - evaluate and set value
			varSymbol, _ := core.ExtractSymbol(assignment.Elements[1]) // Already validated above
			initialValue := e.Evaluate(ctx, assignment.Elements[2])
			if core.IsError(initialValue) {
				return initialValue
			}
			if err := blockCtx.Set(varSymbol, initialValue); err != nil {
				return core.NewErrorExpr("ProtectionError", err.Error(), []core.Expr{varExpr})
			}
		}
	}

	// Evaluate body in the Block context
	result := e.Evaluate(blockCtx, body)

	return result
}

