package sexpr

import (
	"fmt"
	"math"
	"strings"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/stdlib"
)

// Evaluator represents the expression evaluator
type Evaluator struct {
	context *Context
}

// NewEvaluator creates a new evaluator with a fresh context
func NewEvaluator() *Evaluator {
	return &Evaluator{
		context: NewContext(),
	}
}

// NewEvaluatorWithContext creates an evaluator with a specific context
func NewEvaluatorWithContext(ctx *Context) *Evaluator {
	return &Evaluator{
		context: ctx,
	}
}

// Evaluate evaluates an expression in the current context
func (e *Evaluator) Evaluate(expr core.Expr) core.Expr {
	return e.evaluate(expr, e.context)
}

// evaluate is the main evaluation function
func (e *Evaluator) evaluate(expr core.Expr, ctx *Context) core.Expr {
	if expr == nil {
		return nil
	}

	// Push current expression to stack for recursion tracking
	exprStr := expr.String()
	if err := ctx.stack.Push("evaluate", exprStr); err != nil {
		// Return recursion error with stack trace
		return core.NewErrorExprWithStack("RecursionError", err.Error(), []core.Expr{expr}, ctx.stack.GetFrames())
	}
	defer ctx.stack.Pop()

	return e.evaluateExpr(expr, ctx)
}

func (e *Evaluator) evaluateExpr(expr core.Expr, ctx *Context) core.Expr {
	switch ex := expr.(type) {
	case core.Symbol:
		symbolName := string(ex)

		// Check for variable binding first
		if value, ok := ctx.Get(symbolName); ok {
			return value
		}

		// Check for built-in constants
		if constant, ok := e.getBuiltinConstant(symbolName); ok {
			return constant
		}

		// Return the symbol itself if not bound
		return ex
	case core.String, core.Integer, core.Real:
		// New atomic types evaluate to themselves
		return ex
	case core.List:
		return e.evaluateList(ex, ctx)
	default:
		// All other types (ByteArray, Association, ErrorExpr, etc.) evaluate to themselves
		return expr
	}
}

// evaluateList evaluates a list expression
func (e *Evaluator) evaluateList(list core.List, ctx *Context) core.Expr {
	if len(list.Elements) == 0 {
		return list
	}

	// Get the head (function name)
	head := list.Elements[0]
	args := list.Elements[1:]

	// Evaluate the head to get the function name
	evaluatedHead := e.evaluate(head, ctx)

	// Check if head is an error - propagate it
	if core.IsError(evaluatedHead) {
		return evaluatedHead
	}

	// Extract function name from evaluated head
	headName, ok := core.ExtractSymbol(evaluatedHead)
	if !ok {
		// Head is not a symbol, return unevaluated
		return list
	}

	// Apply attribute transformations before evaluation
	transformedList := e.applyAttributeTransformations(headName, list, ctx)
	if !listsEqual(transformedList, list) {
		// The list was transformed, re-evaluate it
		return e.evaluateList(transformedList, ctx)
	}

	// Handle OneIdentity attribute specially - it can return a non-List
	if ctx.symbolTable.HasAttribute(headName, OneIdentity) && len(list.Elements) == 2 {
		// OneIdentity: f(x) = x
		return e.evaluate(list.Elements[1], ctx)
	}

	// Check for special forms first (these don't follow normal evaluation rules)
	if specialResult := e.evaluateSpecialForm(headName, args, ctx); specialResult != nil {
		return specialResult
	}

	// Try pattern-based function resolution
	return e.evaluatePatternFunction(headName, args, ctx)
}

// evaluatePatternFunction evaluates a function using pattern-based dispatch
func (e *Evaluator) evaluatePatternFunction(headName string, args []core.Expr, ctx *Context) core.Expr {
	// Evaluate arguments based on hold attributes
	evaluatedArgs := e.evaluateArguments(headName, args, ctx)

	// Check for errors in evaluated arguments
	for _, arg := range evaluatedArgs {
		if core.IsError(arg) {
			return arg
		}
	}

	// Create the function call expression for pattern matching
	callExpr := core.NewList(headName, evaluatedArgs...)

	// Try to find a matching pattern in the function registry
	if result, found := ctx.functionRegistry.CallFunction(callExpr, ctx); found {
		// Check if result is an error and needs stack trace
		if core.IsError(result) {
			if errorExpr, ok := result.(*core.ErrorExpr); ok {
				// Add stack frame for this function call
				funcCallStr := headName + "(" + formatArgs(evaluatedArgs) + ")"
				if err := ctx.stack.Push(headName, funcCallStr); err == nil {
					ctx.stack.Pop() // Immediately pop since we're just adding to trace
					return core.NewErrorExprWithStack(errorExpr.ErrorType, errorExpr.Message, errorExpr.Args, ctx.stack.GetFrames())
				}
			}
		}
		// Re-evaluate function results until fixed point for proper symbolic computation
		// Only re-evaluate non-atomic expressions to avoid infinite recursion
		if !result.IsAtom() && !result.Equal(callExpr) {
			return e.evaluateToFixedPoint(result, ctx)
		}
		return result
	}

	// No pattern matched, return the unevaluated expression
	return callExpr
}

// evaluateToFixedPoint continues evaluating an expression until it reaches a fixed point
// (no more changes occur) or until a maximum number of iterations to prevent infinite loops
func (e *Evaluator) evaluateToFixedPoint(expr core.Expr, ctx *Context) core.Expr {
	const maxIterations = 100 // Prevent infinite loops
	current := expr

	for i := 0; i < maxIterations; i++ {
		next := e.evaluate(current, ctx)

		// Check if we've reached a fixed point (no more changes)
		if next.Equal(current) {
			return next
		}

		// Check for errors
		if core.IsError(next) {
			return next
		}

		// If the result is atomic, we can't evaluate further
		if next.IsAtom() {
			return next
		}

		current = next
	}

	// If we've hit the iteration limit, return what we have
	// This prevents infinite loops while still allowing significant evaluation
	return current
}

// evaluateArguments evaluates arguments based on hold attributes
func (e *Evaluator) evaluateArguments(headName string, args []core.Expr, ctx *Context) []core.Expr {
	evaluatedArgs := make([]core.Expr, len(args))

	holdAll := ctx.symbolTable.HasAttribute(headName, HoldAll)
	holdFirst := ctx.symbolTable.HasAttribute(headName, HoldFirst)
	holdRest := ctx.symbolTable.HasAttribute(headName, HoldRest)

	for i, arg := range args {
		if holdAll || (holdFirst && i == 0) || (holdRest && i > 0) {
			evaluatedArgs[i] = arg // Don't evaluate
		} else {
			evaluatedArgs[i] = e.evaluate(arg, ctx)
		}
	}

	return evaluatedArgs
}

// applyAttributeTransformations applies attribute-based transformations
func (e *Evaluator) applyAttributeTransformations(headName string, list core.List, ctx *Context) core.List {
	result := list

	// Apply Flat attribute (associativity)
	if ctx.symbolTable.HasAttribute(headName, Flat) {
		result = e.applyFlat(headName, result)
	}

	// Apply Orderless attribute (commutativity)
	if ctx.symbolTable.HasAttribute(headName, Orderless) {
		result = e.applyOrderless(result)
	}

	// Apply OneIdentity attribute
	if ctx.symbolTable.HasAttribute(headName, OneIdentity) {
		result = e.applyOneIdentity(result)
	}

	return result
}

// applyFlat implements the Flat attribute (associativity)
func (e *Evaluator) applyFlat(headName string, list core.List) core.List {
	if len(list.Elements) <= 1 {
		return list
	}

	head := list.Head()
	args := list.Elements[1:]

	newArgs := []core.Expr{}

	for _, arg := range args {
		// If the argument is the same function, flatten it
		if argList, ok := arg.(core.List); ok && len(argList.Elements) > 0 {
			if argHeadName, ok := core.ExtractSymbol(argList.Elements[0]); ok {
				if argHeadName == headName {
					// Flatten: f(a, f(b, c), d) → f(a, b, c, d)
					newArgs = append(newArgs, argList.Elements[1:]...)
					continue
				}
			}
		}
		newArgs = append(newArgs, arg)
	}

	return core.NewList(head, newArgs...)
}

// applyOrderless implements the Orderless attribute (commutativity)
func (e *Evaluator) applyOrderless(list core.List) core.List {
	// Use the stdlib Sort function for consistent ordering
	sorted := stdlib.Sort(list)
	if sortedList, ok := sorted.(core.List); ok {
		return sortedList
	}
	return list
}

// applyOneIdentity implements the OneIdentity attribute
func (e *Evaluator) applyOneIdentity(list core.List) core.List {
	// OneIdentity is now handled specially in evaluateList
	// This function is kept for consistency but doesn't transform anything
	return list
}

// evaluateSpecialForm handles special forms that don't follow normal evaluation rules
func (e *Evaluator) evaluateSpecialForm(headName string, args []core.Expr, ctx *Context) core.Expr {
	switch headName {
	case "If":
		return e.evaluateIf(args, ctx)
	case "Set":
		return e.evaluateSet(args, ctx)
	case "SetDelayed":
		return e.evaluateSetDelayed(args, ctx)
	case "Unset":
		return e.evaluateUnset(args, ctx)
	case "Hold":
		return e.evaluateHold(args, ctx)
	case "Evaluate":
		return e.evaluateEvaluate(args, ctx)
	case "CompoundExpression":
		return e.evaluateCompoundExpression(args, ctx)
	case "CompoundStatement":
		return e.evaluateCompoundExpression(args, ctx)
	case "And":
		return e.evaluateAnd(args, ctx)
	case "Or":
		return e.evaluateOr(args, ctx)
	case "SliceRange":
		return e.evaluateSliceRange(args, ctx)
	case "TakeFrom":
		return e.evaluateTakeFrom(args, ctx)
	case "PartSet":
		return e.evaluatePartSet(args, ctx)
	case "SliceSet":
		return e.evaluateSliceSet(args, ctx)
	case "Block":
		return e.evaluateBlock(args, ctx)
	case "Table":
		return e.evaluateTable(args, ctx)
	case "Do":
		return e.evaluateDo(args, ctx)
	case "Pattern":
		return e.evaluatePattern(args, ctx)
	default:
		return nil // Not a special form
	}
}

// getBuiltinConstant returns built-in constants
func (e *Evaluator) getBuiltinConstant(name string) (core.Expr, bool) {
	switch name {
	case "Pi":
		return core.NewReal(math.Pi), true
	case "E":
		return core.NewReal(math.E), true
	case "True":
		return core.NewBool(true), true
	case "False":
		return core.NewBool(false), true
	case "Null":
		return core.NewSymbolNull(), true
	}
	return nil, false
}

// Utility functions for numeric operations

// isNumeric checks if an expression is numeric
func isNumeric(expr core.Expr) bool {
	// Check new atomic types first
	switch expr.(type) {
	case core.Integer, core.Real:
		return true
	}
	return false
}

// getNumericValue extracts numeric value from an expression
func getNumericValue(expr core.Expr) (float64, bool) {
	// Check new atomic types first
	switch ex := expr.(type) {
	case core.Integer:
		return float64(ex), true
	case core.Real:
		return float64(ex), true
	}
	return 0, false
}

// createNumericResult creates appropriate numeric result (int if whole, float otherwise)
func createNumericResult(value float64) core.Expr {
	if value == float64(int64(value)) {
		return core.NewInteger(int64(value))
	}
	return core.NewReal(value)
}

// getBoolValue extracts boolean value from an expression
// NOTE: This wraps core.ExtractBool for backward compatibility
func getBoolValue(expr core.Expr) (bool, bool) {
	return core.ExtractBool(expr)
}

// isSymbol checks if an expression is a symbol
// NOTE: This wraps core.IsSymbol for backward compatibility
func isSymbol(expr core.Expr) bool {
	return core.IsSymbol(expr)
}

// patternsEqual compares two patterns for equivalence
// This ignores variable names and only compares pattern structure and types
func patternsEqual(pattern1, pattern2 core.Expr) bool {
	// Get pattern info for both patterns
	info1 := core.GetSymbolicPatternInfo(pattern1)
	info2 := core.GetSymbolicPatternInfo(pattern2)

	// If both are patterns, compare their structure (ignoring variable names)
	if (info1 != core.PatternInfo{} && info2 != core.PatternInfo{}) {
		return info1.Type == info2.Type && info1.TypeName == info2.TypeName
	}

	// For non-patterns or when one is a pattern and one isn't, do exact comparison
	switch p1 := pattern1.(type) {
	case core.Integer, core.Real, core.String:
		return pattern1.Equal(pattern2)
	case core.Symbol:
		if name2, ok := core.ExtractSymbol(pattern2); ok {
			// For symbol atoms that are pattern variables, ignore the variable name
			name1 := string(p1)
			// name2 already extracted above
			if core.IsPatternVariable(name1) && core.IsPatternVariable(name2) {
				info1 := core.ParsePatternInfo(name1)
				info2 := core.ParsePatternInfo(name2)
				return info1.Type == info2.Type && info1.TypeName == info2.TypeName
			}
			return name1 == name2
		}
		return false
	case core.List:
		if p2, ok := pattern2.(core.List); ok {
			if len(p1.Elements) != len(p2.Elements) {
				return false
			}
			for i := range p1.Elements {
				if !patternsEqual(p1.Elements[i], p2.Elements[i]) {
					return false
				}
			}
			return true
		}
		return false
	default:
		return false
	}
}

// GetContext returns the evaluator's context
func (e *Evaluator) GetContext() *Context {
	return e.context
}

// listsEqual checks if two lists are structurally equal
func listsEqual(list1, list2 core.List) bool {
	return list1.Equal(list2)
}

// formatArgs formats function arguments for stack traces
func formatArgs(args []core.Expr) string {
	if len(args) == 0 {
		return ""
	}

	argStrs := make([]string, len(args))
	for i, arg := range args {
		argStrs[i] = arg.String()
	}

	// Limit length for readability
	result := strings.Join(argStrs, ", ")
	if len(result) > 100 {
		result = result[:97] + "..."
	}

	return result
}

// Special form implementations

// evaluateIf implements the If special form
func (e *Evaluator) evaluateIf(args []core.Expr, ctx *Context) core.Expr {
	if len(args) < 2 || len(args) > 3 {
		return core.NewErrorExpr("ArgumentError", fmt.Sprintf("If expects 2 or 3 arguments, got %d", len(args)), args)
	}

	// Evaluate the condition
	condition := e.evaluate(args[0], ctx)
	if core.IsError(condition) {
		return condition
	}

	// Check the condition
	if boolVal, isBool := getBoolValue(condition); isBool {
		if boolVal {
			// Condition is true, evaluate and return the "then" branch
			return e.evaluate(args[1], ctx)
		} else {
			// Condition is false, evaluate and return the "else" branch if present
			if len(args) == 3 {
				return e.evaluate(args[2], ctx)
			} else {
				return core.NewSymbolNull()
			}
		}
	}

	// Condition is not a boolean, return an error
	return core.NewErrorExpr("TypeError", "If condition must be True or False", []core.Expr{condition})
}

// evaluateSet implements the Set special form (immediate assignment)
func (e *Evaluator) evaluateSet(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError", fmt.Sprintf("Set expects 2 arguments, got %d", len(args)), args)
	}

	// First argument should be a symbol (don't evaluate it)
	if symbolName, ok := core.ExtractSymbol(args[0]); ok {
		// Prevent direct assignment to $ variables
		if len(symbolName) > 0 && symbolName[0] == '$' {
			return core.NewErrorExpr("ProtectionError", "$ variables cannot be assigned directly", args)
		}

		// Evaluate the value
		value := e.evaluate(args[1], ctx)
		if core.IsError(value) {
			return value
		}

		// Set the variable
		if err := ctx.Set(symbolName, value); err != nil {
			return core.NewErrorExpr("ProtectionError", err.Error(), args)
		}

		return value
	}

	return core.NewErrorExpr("ArgumentError", "First argument to Set must be a symbol", args)
}

// evaluateSetDelayed implements the SetDelayed special form (delayed assignment)
func (e *Evaluator) evaluateSetDelayed(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError", fmt.Sprintf("SetDelayed expects 2 arguments, got %d", len(args)), args)
	}

	lhs := args[0]
	rhs := args[1] // Don't evaluate RHS for delayed assignment

	// Handle function definitions: f(x_) := body
	if list, ok := lhs.(core.List); ok && len(list.Elements) >= 1 {
		// This is a function definition
		headExpr := list.Elements[0]
		if functionName, ok := core.ExtractSymbol(headExpr); ok {

			// Register the pattern with the function registry
			err := ctx.functionRegistry.RegisterFunction(functionName, lhs, func(args []core.Expr, ctx *Context) core.Expr {
				// Create a new child context for function evaluation
				funcCtx := NewChildContext(ctx)

				// Pattern matching and variable binding happen in CallFunction
				// Just evaluate the RHS in the function context
				return e.evaluate(rhs, funcCtx)
			})

			if err != nil {
				return core.NewErrorExpr("DefinitionError", err.Error(), args)
			}

			return core.NewSymbolNull()
		}
	}

	// Handle simple variable assignment: x := value
	if symbolName, ok := core.ExtractSymbol(lhs); ok {
		// Prevent direct assignment to $ variables
		if len(symbolName) > 0 && symbolName[0] == '$' {
			return core.NewErrorExpr("ProtectionError", "$ variables cannot be assigned directly", args)
		}

		// For SetDelayed, store the unevaluated RHS
		if err := ctx.Set(symbolName, rhs); err != nil {
			return core.NewErrorExpr("ProtectionError", err.Error(), args)
		}

		return core.NewSymbolNull()
	}

	return core.NewErrorExpr("ArgumentError", "Invalid left-hand side for SetDelayed", args)
}

// evaluateUnset implements the Unset special form
func (e *Evaluator) evaluateUnset(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 1 {
		return core.NewErrorExpr("ArgumentError", fmt.Sprintf("Unset expects 1 argument, got %d", len(args)), args)
	}

	if symbolName, ok := core.ExtractSymbol(args[0]); ok {
		// Remove the variable binding
		delete(ctx.variables, symbolName)
		return core.NewSymbolNull()
	}

	return core.NewErrorExpr("ArgumentError", "Argument to Unset must be a symbol", args)
}

// evaluateHold implements the Hold special form
func (e *Evaluator) evaluateHold(args []core.Expr, ctx *Context) core.Expr {
	// Hold returns its arguments unevaluated wrapped in Hold
	return core.NewList("Hold", args...)
}

// evaluatePattern implements the Pattern special form
func (e *Evaluator) evaluatePattern(args []core.Expr, ctx *Context) core.Expr {
	// Pattern expressions should remain unevaluated during normal evaluation
	// They are only processed during pattern matching operations
	return core.NewList("Pattern", args...)
}

// evaluateEvaluate implements the Evaluate special form
func (e *Evaluator) evaluateEvaluate(args []core.Expr, ctx *Context) core.Expr {
	if len(args) == 0 {
		return core.NewSymbolNull()
	}

	if len(args) == 1 {
		// Evaluate the single argument
		return e.evaluate(args[0], ctx)
	}

	// Multiple arguments - evaluate all and return the last result
	var result core.Expr = core.NewSymbolNull()
	for _, arg := range args {
		result = e.evaluate(arg, ctx)
		if core.IsError(result) {
			return result
		}
	}
	return result
}

// evaluateCompoundExpression implements the CompoundExpression special form
func (e *Evaluator) evaluateCompoundExpression(args []core.Expr, ctx *Context) core.Expr {
	if len(args) == 0 {
		return core.NewSymbolNull()
	}

	var result core.Expr = core.NewSymbolNull()
	for _, arg := range args {
		result = e.evaluate(arg, ctx)
		if core.IsError(result) {
			return result
		}
	}
	return result
}

// evaluateAnd implements the And special form with short-circuit evaluation
func (e *Evaluator) evaluateAnd(args []core.Expr, ctx *Context) core.Expr {
	if len(args) == 0 {
		return core.NewBool(true)
	}

	var nonBooleanArgs []core.Expr

	for _, arg := range args {
		// Evaluate this argument
		result := e.evaluate(arg, ctx)
		if core.IsError(result) {
			return result
		}

		if boolVal, isBool := getBoolValue(result); isBool {
			if !boolVal {
				return core.NewBool(false) // Short-circuit on first False
			}
			// True values are eliminated (identity for And)
		} else {
			// Collect non-boolean values
			nonBooleanArgs = append(nonBooleanArgs, result)
		}
	}

	// Handle results based on remaining non-boolean arguments
	if len(nonBooleanArgs) == 0 {
		return core.NewBool(true) // All were True
	} else if len(nonBooleanArgs) == 1 {
		return nonBooleanArgs[0] // Single non-boolean argument
	} else {
		// Multiple non-boolean arguments, return simplified And expression
		return core.NewList("And", nonBooleanArgs...)
	}
}

// evaluateOr implements the Or special form with short-circuit evaluation
func (e *Evaluator) evaluateOr(args []core.Expr, ctx *Context) core.Expr {
	if len(args) == 0 {
		return core.NewBool(false)
	}

	var nonBooleanArgs []core.Expr

	for _, arg := range args {
		result := e.evaluate(arg, ctx)
		if core.IsError(result) {
			return result
		}

		if boolVal, isBool := getBoolValue(result); isBool {
			if boolVal {
				return core.NewBool(true) // Short-circuit on first True
			}
			// False values are eliminated (identity for Or)
		} else {
			// Collect non-boolean values
			nonBooleanArgs = append(nonBooleanArgs, result)
		}
	}

	// Handle results based on remaining non-boolean arguments
	if len(nonBooleanArgs) == 0 {
		return core.NewBool(false) // All were False
	} else if len(nonBooleanArgs) == 1 {
		return nonBooleanArgs[0] // Single non-boolean argument
	} else {
		// Multiple non-boolean arguments, return simplified Or expression
		return core.NewList("Or", nonBooleanArgs...)
	}
}

// evaluateSliceRange implements slice range syntax: expr[start:end]
func (e *Evaluator) evaluateSliceRange(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 3 {
		return core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("SliceRange expects 3 arguments (expr, start, end), got %d", len(args)), args)
	}

	// Evaluate the expression being sliced
	expr := e.evaluate(args[0], ctx)
	if core.IsError(expr) {
		return expr
	}

	// Check if the expression is sliceable
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewErrorExpr("TypeError",
			fmt.Sprintf("Expression of type %s is not sliceable", expr.Head()), []core.Expr{expr})
	}

	// Evaluate start and end indices
	startExpr := e.evaluate(args[1], ctx)
	if core.IsError(startExpr) {
		return startExpr
	}

	endExpr := e.evaluate(args[2], ctx)
	if core.IsError(endExpr) {
		return endExpr
	}

	// Extract integer values for start and end
	start, ok := core.ExtractInt64(startExpr)
	if !ok {
		return core.NewErrorExpr("TypeError",
			fmt.Sprintf("Slice start index must be an integer, got %s", startExpr.Head()), []core.Expr{startExpr})
	}

	end, ok := core.ExtractInt64(endExpr)
	if !ok {
		return core.NewErrorExpr("TypeError",
			fmt.Sprintf("Slice end index must be an integer, got %s", endExpr.Head()), []core.Expr{endExpr})
	}

	// Use the Sliceable interface to perform the slice operation
	return sliceable.Slice(start, end)
}

// evaluateTakeFrom implements slice syntax: expr[start:]
// If start is negative, uses Take for last n elements
// If start is positive, uses Drop for first n elements
func (e *Evaluator) evaluateTakeFrom(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("TakeFrom expects 2 arguments (expr, start), got %d", len(args)), args)
	}

	// Evaluate the expression being sliced
	expr := e.evaluate(args[0], ctx)
	if core.IsError(expr) {
		return expr
	}

	// Evaluate start index
	startExpr := e.evaluate(args[1], ctx)
	if core.IsError(startExpr) {
		return startExpr
	}

	// Extract integer value for start
	start, ok := core.ExtractInt64(startExpr)
	if !ok {
		return core.NewErrorExpr("TypeError",
			fmt.Sprintf("Slice start index must be an integer, got %s", startExpr.Head()), []core.Expr{startExpr})
	}

	if start < 0 {
		// Negative start: use Take to get last |start| elements
		// Take([1,2,3,4,5], -2) gives [4,5]
		return e.evaluate(core.NewList("Take", expr, core.NewInteger(start)), ctx)
	} else {
		// Positive start: use Drop to remove first (start-1) elements
		// Drop([1,2,3,4,5], 2) gives [3,4,5] (for start=3, 1-indexed)
		dropCount := start - 1
		return e.evaluate(core.NewList("Drop", expr, core.NewInteger(dropCount)), ctx)
	}
}

// evaluatePartSet implements slice assignment syntax: expr[index] = value
func (e *Evaluator) evaluatePartSet(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 3 {
		return core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("PartSet expects 3 arguments (expr, index, value), got %d", len(args)), args)
	}

	// Evaluate the expression being modified
	expr := e.evaluate(args[0], ctx)
	if core.IsError(expr) {
		return expr
	}

	// Check if the expression is sliceable
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewErrorExpr("TypeError",
			fmt.Sprintf("Expression of type %s is not sliceable", expr.Head()), []core.Expr{expr})
	}

	// Evaluate index
	indexExpr := e.evaluate(args[1], ctx)
	if core.IsError(indexExpr) {
		return indexExpr
	}

	// Extract integer value for index
	index, ok := core.ExtractInt64(indexExpr)
	if !ok {
		return core.NewErrorExpr("TypeError",
			fmt.Sprintf("Part index must be an integer, got %s", indexExpr.Head()), []core.Expr{indexExpr})
	}

	// Evaluate value
	value := e.evaluate(args[2], ctx)
	if core.IsError(value) {
		return value
	}

	// Use the Sliceable interface to perform the assignment
	return sliceable.SetElementAt(index, value)
}

// evaluateSliceSet implements slice assignment syntax: expr[start:end] = value
func (e *Evaluator) evaluateSliceSet(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 4 {
		return core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("SliceSet expects 4 arguments (expr, start, end, value), got %d", len(args)), args)
	}

	// Evaluate the expression being modified
	expr := e.evaluate(args[0], ctx)
	if core.IsError(expr) {
		return expr
	}

	// Check if the expression is sliceable
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewErrorExpr("TypeError",
			fmt.Sprintf("Expression of type %s is not sliceable", expr.Head()), []core.Expr{expr})
	}

	// Evaluate start index
	startExpr := e.evaluate(args[1], ctx)
	if core.IsError(startExpr) {
		return startExpr
	}

	// Extract integer value for start
	start, ok := core.ExtractInt64(startExpr)
	if !ok {
		return core.NewErrorExpr("TypeError",
			fmt.Sprintf("Slice start index must be an integer, got %s", startExpr.Head()), []core.Expr{startExpr})
	}

	// Evaluate end index
	endExpr := e.evaluate(args[2], ctx)
	if core.IsError(endExpr) {
		return endExpr
	}

	// Extract integer value for end (handle special case of -1 for "to end")
	var end int64
	if endValue, ok := core.ExtractInt64(endExpr); ok && endValue == -1 {
		// Special case: -1 means "to end of sequence"
		end = sliceable.(interface{ Length() int64 }).Length()
	} else if endValue, ok := core.ExtractInt64(endExpr); ok {
		end = endValue
	} else {
		return core.NewErrorExpr("TypeError",
			fmt.Sprintf("Slice end index must be an integer, got %s", endExpr.Head()), []core.Expr{endExpr})
	}

	// Evaluate value
	value := e.evaluate(args[3], ctx)
	if core.IsError(value) {
		return value
	}

	// Use the Sliceable interface to perform the slice assignment
	return sliceable.SetSlice(start, end, value)
}

// evaluateBlock implements the Block special form for dynamic scoping
// Block(List(vars...), body) temporarily changes variable values and evaluates body
func (e *Evaluator) evaluateBlock(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError", fmt.Sprintf("Block expects 2 arguments, got %d", len(args)), args)
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
	blockCtx := NewBlockContext(ctx, scopedVarNames)

	// Set initial values for Block variables
	for i := 1; i < len(varList.Elements); i++ { // Skip head element
		varExpr := varList.Elements[i]

		if _, ok := core.ExtractSymbol(varExpr); ok {
			// Simple variable: {x} - leave unset (will return undefined during lookup)
			// Don't set anything - the variable is scoped but has no value

		} else if assignment, ok := varExpr.(core.List); ok && len(assignment.Elements) == 3 {
			// Assignment: {x = value} - evaluate and set value
			varSymbol, _ := core.ExtractSymbol(assignment.Elements[1]) // Already validated above
			initialValue := e.evaluate(assignment.Elements[2], ctx)
			if core.IsError(initialValue) {
				return initialValue
			}
			if err := blockCtx.Set(varSymbol, initialValue); err != nil {
				return core.NewErrorExpr("ProtectionError", err.Error(), []core.Expr{varExpr})
			}
		}
	}

	// Evaluate body in the Block context
	result := e.evaluate(body, blockCtx)

	return result
}

// evaluateDo implements the Do special form for iteration without collecting results
// Do(expr, n) evaluates expr n times and returns Null
// Do(expr, core.List(i, start, end, increment)) iterates with variable binding
func (e *Evaluator) evaluateDo(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError", fmt.Sprintf("Do expects 2 arguments, got %d", len(args)), args)
	}

	expr := args[0] // Don't evaluate expr yet - Do has HoldAll
	spec := args[1] // Don't evaluate spec yet

	// Check if second argument is an integer (simple replication form)
	if n, ok := core.ExtractInt64(spec); ok {
		return e.evaluateDoSimple(expr, n, ctx)
	}

	// Check if second argument is a core.List (iterator form)
	if iterList, ok := spec.(core.List); ok {
		return e.evaluateDoIterator(expr, iterList, ctx)
	}

	return core.NewErrorExpr("ArgumentError", "Do second argument must be integer or core.List", args)
}

// evaluateDoSimple implements Do(expr, n) - evaluates expr n times and returns Null
func (e *Evaluator) evaluateDoSimple(expr core.Expr, n int64, ctx *Context) core.Expr {
	if n < 0 {
		return core.NewErrorExpr("ArgumentError", fmt.Sprintf("Do count must be non-negative, got %d", n), []core.Expr{core.NewInteger(n)})
	}

	// Evaluate expr n times (side effects only)
	for i := int64(0); i < n; i++ {
		result := e.evaluate(expr, ctx)
		if core.IsError(result) {
			return result // Return error immediately
		}
		// Discard result - Do is for side effects only
	}

	return core.NewSymbol("Null")
}

// evaluateDoIterator handles Do with iterator specification core.List(i, start, end, increment)
func (e *Evaluator) evaluateDoIterator(expr core.Expr, iterSpec core.List, ctx *Context) core.Expr {
	// Parse iterator specification into normalized form
	variable, start, end, increment, err := e.parseTableIteratorSpec(iterSpec, ctx)
	if err != nil {
		return err
	}

	current := start
	const maxIterations = 10000 // Prevent infinite loops

	for iteration := 0; iteration < maxIterations; iteration++ {
		// Check if we should continue iterating
		shouldContinue := e.evaluateIteratorCondition(current, end, increment, ctx)
		if !shouldContinue {
			break
		}

		// Evaluate expression with current iterator value (for side effects only)
		blockResult := e.evaluateWithIteratorBinding(expr, variable, current, ctx)
		if core.IsError(blockResult) {
			return blockResult // Return error immediately
		}
		// Discard result - Do is for side effects only

		// Increment for next iteration
		current = e.evaluateIteratorIncrement(current, increment, ctx)
		if core.IsError(current) {
			return current
		}
	}

	return core.NewSymbol("Null")
}

// evaluateTable implements the Table special form for list generation
// Table(expr, n) creates n copies of expr
// Table(expr, core.List(i, max)) will be implemented later for iterator forms
func (e *Evaluator) evaluateTable(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError", fmt.Sprintf("Table expects 2 arguments, got %d", len(args)), args)
	}

	expr := args[0] // Don't evaluate expr yet - Table has HoldAll
	spec := args[1] // Don't evaluate spec yet

	// Check if second argument is an integer (simple replication form)
	if n, ok := core.ExtractInt64(spec); ok {
		return e.evaluateTableSimple(expr, n, ctx)
	}

	// Check if second argument is a core.List (iterator form)
	if iterList, ok := spec.(core.List); ok {
		return e.evaluateTableIterator(expr, iterList, ctx)
	}

	return core.NewErrorExpr("ArgumentError", "Table second argument must be integer or core.List", args)
}

// evaluateTableSimple implements Table(expr, n) - creates n copies of expr
func (e *Evaluator) evaluateTableSimple(expr core.Expr, n int64, ctx *Context) core.Expr {
	if n < 0 {
		return core.NewErrorExpr("ArgumentError", fmt.Sprintf("Table count must be non-negative, got %d", n), []core.Expr{core.NewInteger(n)})
	}

	if n == 0 {
		return core.NewList("List")
	}

	// Create result list with proper capacity
	elements := make([]core.Expr, n)

	// Evaluate expr once for each position
	for i := 0; i < int(n); i++ {
		// Evaluate expr in current context for each iteration
		// This allows expressions with side effects to work correctly
		evaluated := e.evaluate(expr, ctx)
		if core.IsError(evaluated) {
			return evaluated
		}
		elements[i] = evaluated
	}

	return core.NewList("List", elements...)
}

// evaluateTableIterator implements Table(expr, core.List(i, start, end, increment))
// Handles all iterator forms using the general case with expression-based arithmetic
func (e *Evaluator) evaluateTableIterator(expr core.Expr, iterSpec core.List, ctx *Context) core.Expr {
	// Parse iterator specification into normalized form
	variable, start, end, increment, err := e.parseTableIteratorSpec(iterSpec, ctx)
	if err != nil {
		return err
	}

	var results []core.Expr

	current := start
	const maxIterations = 10000 // Prevent infinite loops

	for iteration := 0; iteration < maxIterations; iteration++ {
		// Check if we should continue iterating
		shouldContinue := e.evaluateIteratorCondition(current, end, increment, ctx)
		if !shouldContinue {
			break
		}

		// Use Block to bind iterator variable and evaluate expression
		blockResult := e.evaluateWithIteratorBinding(expr, variable, current, ctx)
		if core.IsError(blockResult) {
			return blockResult
		}
		results = append(results, blockResult)

		// Increment current value using expression arithmetic
		current = e.evaluateIteratorIncrement(current, increment, ctx)
		if core.IsError(current) {
			return current
		}
	}

	return core.NewList("List", results...)
}

// parseTableIteratorSpec parses iterator specifications and normalizes them
// core.List(i, max) → core.List(i, 1, max, 1)
// core.List(i, start, end) → core.List(i, start, end, 1)
// core.List(i, start, end, increment) → core.List(i, start, end, increment)
// IMPORTANT: Evaluates start, end, and increment expressions and validates they are numeric
func (e *Evaluator) parseTableIteratorSpec(iterSpec core.List, ctx *Context) (variable string, start, end, increment core.Expr, err core.Expr) {
	if len(iterSpec.Elements) < 3 || len(iterSpec.Elements) > 5 {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError",
			"Table iterator must be core.List(var, max), core.List(var, start, end), or core.List(var, start, end, step)", []core.Expr{iterSpec})
	}

	// Extract variable name
	if varSymbol, ok := core.ExtractSymbol(iterSpec.Elements[1]); ok {
		variable = varSymbol
	} else {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError", "Table iterator variable must be a symbol", []core.Expr{iterSpec.Elements[1]})
	}

	// Parse and evaluate based on number of arguments
	switch len(iterSpec.Elements) {
	case 3: // core.List(i, max) → core.List(i, 1, max, 1)
		start = core.NewInteger(1)
		end = e.evaluate(iterSpec.Elements[2], ctx)
		if core.IsError(end) {
			return "", nil, nil, nil, end
		}
		increment = core.NewInteger(1)

	case 4: // core.List(i, start, end) → core.List(i, start, end, 1)
		start = e.evaluate(iterSpec.Elements[2], ctx)
		if core.IsError(start) {
			return "", nil, nil, nil, start
		}
		end = e.evaluate(iterSpec.Elements[3], ctx)
		if core.IsError(end) {
			return "", nil, nil, nil, end
		}
		increment = core.NewInteger(1)

	case 5: // core.List(i, start, end, increment)
		start = e.evaluate(iterSpec.Elements[2], ctx)
		if core.IsError(start) {
			return "", nil, nil, nil, start
		}
		end = e.evaluate(iterSpec.Elements[3], ctx)
		if core.IsError(end) {
			return "", nil, nil, nil, end
		}
		increment = e.evaluate(iterSpec.Elements[4], ctx)
		if core.IsError(increment) {
			return "", nil, nil, nil, increment
		}

	default:
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError", "Invalid Table iterator specification", []core.Expr{iterSpec})
	}

	// Validate that arithmetic and comparison operations can be evaluated
	// Test if Plus(start, increment) evaluates to something different (not unevaluated)
	testPlus := core.NewList("Plus", start, increment)
	plusResult := e.evaluate(testPlus, ctx)
	if core.IsError(plusResult) {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("Table iterator arithmetic failed: %s", plusResult), []core.Expr{plusResult})
	}
	if plusResult.Equal(testPlus) && !plusResult.IsAtom() {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("Table iterator arithmetic unevaluated: Plus(%s, %s) - missing arithmetic definition", start, increment), []core.Expr{start, increment})
	}

	// Test if comparison operation evaluates
	testLessEqual := core.NewList("LessEqual", start, end)
	compareResult := e.evaluate(testLessEqual, ctx)
	if core.IsError(compareResult) {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("Table iterator comparison failed: %s", compareResult), []core.Expr{compareResult})
	}
	if compareResult.Equal(testLessEqual) && !compareResult.IsAtom() {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("Table iterator comparison unevaluated: LessEqual(%s, %s) - missing comparison definition", start, end), []core.Expr{start, end})
	}

	return variable, start, end, increment, nil
}

// evaluateIteratorCondition determines if iteration should continue
// Uses expression-based comparison with proper handling of increment direction
func (e *Evaluator) evaluateIteratorCondition(current, end, increment core.Expr, ctx *Context) bool {
	// Determine comparison operator based on increment sign
	var compSymbol string
	isNegative := e.isNegativeIncrement(increment, ctx)
	if isNegative {
		compSymbol = "GreaterEqual" // For negative increment, continue while current >= end
	} else {
		compSymbol = "LessEqual" // For positive increment, continue while current <= end
	}

	// Create and evaluate comparison expression
	compExpr := core.NewList(compSymbol, current, end)
	result := e.evaluate(compExpr, ctx)

	// Extract boolean result
	if boolVal, ok := core.ExtractBool(result); ok {
		return boolVal
	}

	// If comparison remains unevaluated, check if it's the exact same expression
	if result.Equal(compExpr) && !result.IsAtom() {
		// Comparison is unevaluated - this indicates missing comparison definition
		// This should have been caught during validation, but stop iteration safely
		return false
	}

	// If we get here, the comparison evaluated to something other than a boolean
	// This might be valid in some mathematical contexts, so be conservative
	return false
}

// evaluateIteratorIncrement adds increment to current value using expression arithmetic
func (e *Evaluator) evaluateIteratorIncrement(current, increment core.Expr, ctx *Context) core.Expr {
	plusExpr := core.NewList("Plus", current, increment)
	return e.evaluate(plusExpr, ctx)
}

// isNegativeIncrement determines if increment is negative using expression evaluation
func (e *Evaluator) isNegativeIncrement(increment core.Expr, ctx *Context) bool {
	// Create comparison: increment < 0
	zeroExpr := core.NewInteger(0)
	lessExpr := core.NewList("Less", increment, zeroExpr)
	result := e.evaluate(lessExpr, ctx)

	if boolVal, ok := core.ExtractBool(result); ok {
		return boolVal
	}

	// Default to positive if comparison fails
	return false
}

// evaluateWithIteratorBinding uses Block to bind iterator variable and evaluate expression
func (e *Evaluator) evaluateWithIteratorBinding(expr core.Expr, variable string, value core.Expr, ctx *Context) core.Expr {
	// Create Block(List(Set(variable, value)), expr)
	setExpr := core.NewList("Set", core.NewSymbol(variable), value)
	blockVars := core.NewList("List", setExpr)
	blockArgs := []core.Expr{blockVars, expr}

	return e.evaluateBlock(blockArgs, ctx)
}
