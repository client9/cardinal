package sexpr

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/client9/sexpr/core"
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
func (e *Evaluator) Evaluate(expr Expr) Expr {
	return e.evaluate(expr, e.context)
}

// evaluate is the main evaluation function
func (e *Evaluator) evaluate(expr Expr, ctx *Context) Expr {
	if expr == nil {
		return nil
	}

	// Push current expression to stack for recursion tracking
	exprStr := expr.String()
	if err := ctx.stack.Push("evaluate", exprStr); err != nil {
		// Return recursion error with stack trace
		return NewErrorExprWithStack("RecursionError", err.Error(), []Expr{expr}, ctx.stack.GetFrames())
	}
	defer ctx.stack.Pop()

	return e.evaluateExpr(expr, ctx)
}

func (e *Evaluator) evaluateExpr(expr Expr, ctx *Context) Expr {
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
	case List:
		return e.evaluateList(ex, ctx)
	default:
		// All other types (ByteArray, Association, ErrorExpr, etc.) evaluate to themselves
		return expr
	}
}

// evaluateList evaluates a list expression
func (e *Evaluator) evaluateList(list List, ctx *Context) Expr {
	if len(list.Elements) == 0 {
		return list
	}

	// Get the head (function name)
	head := list.Elements[0]
	args := list.Elements[1:]

	// Evaluate the head to get the function name
	evaluatedHead := e.evaluate(head, ctx)

	// Check if head is an error - propagate it
	if IsError(evaluatedHead) {
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
func (e *Evaluator) evaluatePatternFunction(headName string, args []Expr, ctx *Context) Expr {
	// Evaluate arguments based on hold attributes
	evaluatedArgs := e.evaluateArguments(headName, args, ctx)

	// Check for errors in evaluated arguments
	for _, arg := range evaluatedArgs {
		if IsError(arg) {
			return arg
		}
	}

	// Create the function call expression for pattern matching
	callExpr := NewList(append([]Expr{core.NewSymbol(headName)}, evaluatedArgs...)...)

	// Try to find a matching pattern in the function registry
	if result, found := ctx.functionRegistry.CallFunction(callExpr, ctx); found {
		// Check if result is an error and needs stack trace
		if IsError(result) {
			if errorExpr, ok := result.(*ErrorExpr); ok {
				// Add stack frame for this function call
				funcCallStr := headName + "(" + formatArgs(evaluatedArgs) + ")"
				if err := ctx.stack.Push(headName, funcCallStr); err == nil {
					ctx.stack.Pop() // Immediately pop since we're just adding to trace
					return NewErrorExprWithStack(errorExpr.ErrorType, errorExpr.Message, errorExpr.Args, ctx.stack.GetFrames())
				}
			}
		}
		return result
	}

	// No pattern matched, return the unevaluated expression
	return callExpr
}

// evaluateArguments evaluates arguments based on hold attributes
func (e *Evaluator) evaluateArguments(headName string, args []Expr, ctx *Context) []Expr {
	evaluatedArgs := make([]Expr, len(args))

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
func (e *Evaluator) applyAttributeTransformations(headName string, list List, ctx *Context) List {
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
func (e *Evaluator) applyFlat(headName string, list List) List {
	if len(list.Elements) <= 1 {
		return list
	}

	head := list.Elements[0]
	args := list.Elements[1:]
	newArgs := []Expr{}

	for _, arg := range args {
		// If the argument is the same function, flatten it
		if argList, ok := arg.(List); ok && len(argList.Elements) > 0 {
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

	return NewList(append([]Expr{head}, newArgs...)...)
}

// applyOrderless implements the Orderless attribute (commutativity)
func (e *Evaluator) applyOrderless(list List) List {
	if len(list.Elements) <= 2 {
		return list
	}

	head := list.Elements[0]
	args := list.Elements[1:]

	// Sort arguments by their string representation for canonical ordering
	sort.Slice(args, func(i, j int) bool {
		return args[i].String() < args[j].String()
	})

	return NewList(append([]Expr{head}, args...)...)
}

// applyOneIdentity implements the OneIdentity attribute
func (e *Evaluator) applyOneIdentity(list List) List {
	// OneIdentity is now handled specially in evaluateList
	// This function is kept for consistency but doesn't transform anything
	return list
}

// evaluateSpecialForm handles special forms that don't follow normal evaluation rules
func (e *Evaluator) evaluateSpecialForm(headName string, args []Expr, ctx *Context) Expr {
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
	default:
		return nil // Not a special form
	}
}

// getBuiltinConstant returns built-in constants
func (e *Evaluator) getBuiltinConstant(name string) (Expr, bool) {
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
func isNumeric(expr Expr) bool {
	// Check new atomic types first
	switch expr.(type) {
	case core.Integer, core.Real:
		return true
	}
	return false
}

// getNumericValue extracts numeric value from an expression
func getNumericValue(expr Expr) (float64, bool) {
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
func createNumericResult(value float64) Expr {
	if value == float64(int64(value)) {
		return core.NewInteger(int64(value))
	}
	return core.NewReal(value)
}

// getBoolValue extracts boolean value from an expression
// NOTE: This wraps core.ExtractBool for backward compatibility
func getBoolValue(expr Expr) (bool, bool) {
	return core.ExtractBool(expr)
}

// isSymbol checks if an expression is a symbol
// NOTE: This wraps core.IsSymbol for backward compatibility
func isSymbol(expr Expr) bool {
	return core.IsSymbol(expr)
}

// patternsEqual compares two patterns for equivalence
// This ignores variable names and only compares pattern structure and types
func patternsEqual(pattern1, pattern2 Expr) bool {
	// Get pattern info for both patterns
	info1 := getSymbolicPatternInfo(pattern1)
	info2 := getSymbolicPatternInfo(pattern2)

	// If both are patterns, compare their structure (ignoring variable names)
	if (info1 != PatternInfo{} && info2 != PatternInfo{}) {
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
			if isPatternVariable(name1) && isPatternVariable(name2) {
				info1 := parsePatternInfo(name1)
				info2 := parsePatternInfo(name2)
				return info1.Type == info2.Type && info1.TypeName == info2.TypeName
			}
			return name1 == name2
		}
		return false
	case List:
		if p2, ok := pattern2.(List); ok {
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
func listsEqual(list1, list2 List) bool {
	return list1.Equal(list2)
}

// formatArgs formats function arguments for stack traces
func formatArgs(args []Expr) string {
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
func (e *Evaluator) evaluateIf(args []Expr, ctx *Context) Expr {
	if len(args) < 2 || len(args) > 3 {
		return NewErrorExpr("ArgumentError", fmt.Sprintf("If expects 2 or 3 arguments, got %d", len(args)), args)
	}

	// Evaluate the condition
	condition := e.evaluate(args[0], ctx)
	if IsError(condition) {
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
	return NewErrorExpr("TypeError", "If condition must be True or False", []Expr{condition})
}

// evaluateSet implements the Set special form (immediate assignment)
func (e *Evaluator) evaluateSet(args []Expr, ctx *Context) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError", fmt.Sprintf("Set expects 2 arguments, got %d", len(args)), args)
	}

	// First argument should be a symbol (don't evaluate it)
	if symbolName, ok := core.ExtractSymbol(args[0]); ok {

		// Evaluate the value
		value := e.evaluate(args[1], ctx)
		if IsError(value) {
			return value
		}

		// Set the variable
		ctx.Set(symbolName, value)
		return value
	}

	return NewErrorExpr("ArgumentError", "First argument to Set must be a symbol", args)
}

// evaluateSetDelayed implements the SetDelayed special form (delayed assignment)
func (e *Evaluator) evaluateSetDelayed(args []Expr, ctx *Context) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError", fmt.Sprintf("SetDelayed expects 2 arguments, got %d", len(args)), args)
	}

	lhs := args[0]
	rhs := args[1] // Don't evaluate RHS for delayed assignment

	// Handle function definitions: f(x_) := body
	if list, ok := lhs.(List); ok && len(list.Elements) >= 1 {
		// This is a function definition
		headExpr := list.Elements[0]
		if functionName, ok := core.ExtractSymbol(headExpr); ok {

			// Register the pattern with the function registry
			err := ctx.functionRegistry.RegisterFunction(functionName, lhs, func(args []Expr, ctx *Context) Expr {
				// Create a new child context for function evaluation
				funcCtx := NewChildContext(ctx)

				// Pattern matching and variable binding happen in CallFunction
				// Just evaluate the RHS in the function context
				return e.evaluate(rhs, funcCtx)
			})

			if err != nil {
				return NewErrorExpr("DefinitionError", err.Error(), args)
			}

			return core.NewSymbolNull()
		}
	}

	// Handle simple variable assignment: x := value
	if symbolName, ok := core.ExtractSymbol(lhs); ok {
		// For SetDelayed, store the unevaluated RHS
		ctx.Set(symbolName, rhs)
		return core.NewSymbolNull()
	}

	return NewErrorExpr("ArgumentError", "Invalid left-hand side for SetDelayed", args)
}

// evaluateUnset implements the Unset special form
func (e *Evaluator) evaluateUnset(args []Expr, ctx *Context) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError", fmt.Sprintf("Unset expects 1 argument, got %d", len(args)), args)
	}

	if symbolName, ok := core.ExtractSymbol(args[0]); ok {
		// Remove the variable binding
		delete(ctx.variables, symbolName)
		return core.NewSymbolNull()
	}

	return NewErrorExpr("ArgumentError", "Argument to Unset must be a symbol", args)
}

// evaluateHold implements the Hold special form
func (e *Evaluator) evaluateHold(args []Expr, ctx *Context) Expr {
	// Hold returns its arguments unevaluated wrapped in Hold
	return NewList(append([]Expr{core.NewSymbol("Hold")}, args...)...)
}

// evaluateEvaluate implements the Evaluate special form
func (e *Evaluator) evaluateEvaluate(args []Expr, ctx *Context) Expr {
	if len(args) == 0 {
		return core.NewSymbolNull()
	}

	if len(args) == 1 {
		// Evaluate the single argument
		return e.evaluate(args[0], ctx)
	}

	// Multiple arguments - evaluate all and return the last result
	var result Expr = core.NewSymbolNull()
	for _, arg := range args {
		result = e.evaluate(arg, ctx)
		if IsError(result) {
			return result
		}
	}
	return result
}

// evaluateCompoundExpression implements the CompoundExpression special form
func (e *Evaluator) evaluateCompoundExpression(args []Expr, ctx *Context) Expr {
	if len(args) == 0 {
		return core.NewSymbolNull()
	}

	var result Expr = core.NewSymbolNull()
	for _, arg := range args {
		result = e.evaluate(arg, ctx)
		if IsError(result) {
			return result
		}
	}
	return result
}

// evaluateAnd implements the And special form with short-circuit evaluation
func (e *Evaluator) evaluateAnd(args []Expr, ctx *Context) Expr {
	if len(args) == 0 {
		return core.NewBool(true)
	}

	var nonBooleanArgs []Expr

	for _, arg := range args {
		// Evaluate this argument
		result := e.evaluate(arg, ctx)
		if IsError(result) {
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
		return NewList(append([]Expr{core.NewSymbol("And")}, nonBooleanArgs...)...)
	}
}

// evaluateOr implements the Or special form with short-circuit evaluation
func (e *Evaluator) evaluateOr(args []Expr, ctx *Context) Expr {
	if len(args) == 0 {
		return core.NewBool(false)
	}

	var nonBooleanArgs []Expr

	for _, arg := range args {
		result := e.evaluate(arg, ctx)
		if IsError(result) {
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
		return NewList(append([]Expr{core.NewSymbol("Or")}, nonBooleanArgs...)...)
	}
}

// evaluateSliceRange implements slice range syntax: expr[start:end]
func (e *Evaluator) evaluateSliceRange(args []Expr, ctx *Context) Expr {
	if len(args) != 3 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("SliceRange expects 3 arguments (expr, start, end), got %d", len(args)), args)
	}

	// Evaluate the expression being sliced
	expr := e.evaluate(args[0], ctx)
	if IsError(expr) {
		return expr
	}

	// Check if the expression is sliceable
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Expression of type %s is not sliceable", expr.Type()), []Expr{expr})
	}

	// Evaluate start and end indices
	startExpr := e.evaluate(args[1], ctx)
	if IsError(startExpr) {
		return startExpr
	}

	endExpr := e.evaluate(args[2], ctx)
	if IsError(endExpr) {
		return endExpr
	}

	// Extract integer values for start and end
	start, ok := core.ExtractInt64(startExpr)
	if !ok {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Slice start index must be an integer, got %s", startExpr.Type()), []Expr{startExpr})
	}

	end, ok := core.ExtractInt64(endExpr)
	if !ok {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Slice end index must be an integer, got %s", endExpr.Type()), []Expr{endExpr})
	}

	// Use the Sliceable interface to perform the slice operation
	return sliceable.Slice(start, end)
}

// evaluateTakeFrom implements slice syntax: expr[start:]
// If start is negative, uses Take for last n elements
// If start is positive, uses Drop for first n elements
func (e *Evaluator) evaluateTakeFrom(args []Expr, ctx *Context) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("TakeFrom expects 2 arguments (expr, start), got %d", len(args)), args)
	}

	// Evaluate the expression being sliced
	expr := e.evaluate(args[0], ctx)
	if IsError(expr) {
		return expr
	}

	// Evaluate start index
	startExpr := e.evaluate(args[1], ctx)
	if IsError(startExpr) {
		return startExpr
	}

	// Extract integer value for start
	start, ok := core.ExtractInt64(startExpr)
	if !ok {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Slice start index must be an integer, got %s", startExpr.Type()), []Expr{startExpr})
	}

	if start < 0 {
		// Negative start: use Take to get last |start| elements
		// Take([1,2,3,4,5], -2) gives [4,5]
		return e.evaluate(NewList(core.NewSymbol("Take"), expr, core.NewInteger(start)), ctx)
	} else {
		// Positive start: use Drop to remove first (start-1) elements
		// Drop([1,2,3,4,5], 2) gives [3,4,5] (for start=3, 1-indexed)
		dropCount := start - 1
		return e.evaluate(NewList(core.NewSymbol("Drop"), expr, core.NewInteger(dropCount)), ctx)
	}
}

// evaluatePartSet implements slice assignment syntax: expr[index] = value
func (e *Evaluator) evaluatePartSet(args []Expr, ctx *Context) Expr {
	if len(args) != 3 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("PartSet expects 3 arguments (expr, index, value), got %d", len(args)), args)
	}

	// Evaluate the expression being modified
	expr := e.evaluate(args[0], ctx)
	if IsError(expr) {
		return expr
	}

	// Check if the expression is sliceable
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Expression of type %s is not sliceable", expr.Type()), []Expr{expr})
	}

	// Evaluate index
	indexExpr := e.evaluate(args[1], ctx)
	if IsError(indexExpr) {
		return indexExpr
	}

	// Extract integer value for index
	index, ok := core.ExtractInt64(indexExpr)
	if !ok {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Part index must be an integer, got %s", indexExpr.Type()), []Expr{indexExpr})
	}

	// Evaluate value
	value := e.evaluate(args[2], ctx)
	if IsError(value) {
		return value
	}

	// Use the Sliceable interface to perform the assignment
	return sliceable.SetElementAt(index, value)
}

// evaluateSliceSet implements slice assignment syntax: expr[start:end] = value
func (e *Evaluator) evaluateSliceSet(args []Expr, ctx *Context) Expr {
	if len(args) != 4 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("SliceSet expects 4 arguments (expr, start, end, value), got %d", len(args)), args)
	}

	// Evaluate the expression being modified
	expr := e.evaluate(args[0], ctx)
	if IsError(expr) {
		return expr
	}

	// Check if the expression is sliceable
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Expression of type %s is not sliceable", expr.Type()), []Expr{expr})
	}

	// Evaluate start index
	startExpr := e.evaluate(args[1], ctx)
	if IsError(startExpr) {
		return startExpr
	}

	// Extract integer value for start
	start, ok := core.ExtractInt64(startExpr)
	if !ok {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Slice start index must be an integer, got %s", startExpr.Type()), []Expr{startExpr})
	}

	// Evaluate end index
	endExpr := e.evaluate(args[2], ctx)
	if IsError(endExpr) {
		return endExpr
	}

	// Extract integer value for end (handle special case of -1 for "to end")
	var end int64

	// TODO: UGLY
	if endAtom, ok := endExpr.(core.Integer); ok && int(endAtom) == -1 {
		// Special case: -1 means "to end of sequence"
		end = sliceable.(interface{ Length() int64 }).Length()
	} else {
		var ok bool
		end, ok = core.ExtractInt64(endExpr)
		if !ok {
			return NewErrorExpr("TypeError",
				fmt.Sprintf("Slice end index must be an integer, got %s", endExpr.Type()), []Expr{endExpr})
		}
	}

	// Evaluate value
	value := e.evaluate(args[3], ctx)
	if IsError(value) {
		return value
	}

	// Use the Sliceable interface to perform the slice assignment
	return sliceable.SetSlice(start, end, value)
}
