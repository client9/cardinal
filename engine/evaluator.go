package engine

import (
	"fmt"
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

// GetContext returns the evaluator's current context
func (e *Evaluator) GetContext() *Context {
	return e.context
}

// Evaluate evaluates an expression in the current context
func (e *Evaluator) Evaluate(c *Context, expr core.Expr) core.Expr {
	return e.evaluate(c, expr)
}

// evaluate is the main evaluation function
func (e *Evaluator) evaluate(ctx *Context, expr core.Expr) core.Expr {
	// Push current expression to stack for recursion tracking

	if err := ctx.stack.Push("evaluate", expr); err != nil {
		// Return recursion error with stack trace
		return core.NewErrorExprWithStack("RecursionError", err.Error(), []core.Expr{expr}, ctx.stack.GetFrames())
	}
	defer ctx.stack.Pop()
	return e.evaluateToFixedPoint(ctx, expr)
	//return e.evaluateExpr(ctx, expr)
}

// evaluate is the main evaluation function
func (e *Evaluator) evaluateOnce(ctx *Context, expr core.Expr) core.Expr {
	if err := ctx.stack.Push("evaluate", expr); err != nil {
		// Return recursion error with stack trace
		return core.NewErrorExprWithStack("RecursionError", err.Error(), []core.Expr{expr}, ctx.stack.GetFrames())
	}
	defer ctx.stack.Pop()
	return e.evaluateExpr(ctx, expr)
}

func (e *Evaluator) evaluateExpr(ctx *Context, expr core.Expr) core.Expr {
	switch ex := expr.(type) {
	case core.Symbol:
		symbolName := string(ex)

		// Check for variable binding first
		if value, ok := ctx.Get(symbolName); ok {
			return value
		}
		// Return the symbol itself if not bound
		return ex
	case core.String, core.Integer, core.Real:
		// New atomic types evaluate to themselves
		return ex
	case core.List:
		return e.evaluateList(ctx, ex)
	default:
		// All other types (ByteArray, Association, ErrorExpr, etc.) evaluate to themselves
		return expr
	}
}

// evaluateList evaluates a list expression
func (e *Evaluator) evaluateList(ctx *Context, list core.List) core.Expr {

	if len(list.Elements) == 0 {
		return list
	}

	// Get the head (function name)
	head := list.Elements[0]
	args := list.Elements[1:]

	// Evaluate the head to get the function name
	evaluatedHead := e.evaluate(ctx, head)

	// Check if head is an error - propagate it
	if core.IsError(evaluatedHead) {
		return evaluatedHead
	}

	// Check if head is a function expression (function application)
	if funcExpr, ok := evaluatedHead.(core.FunctionExpr); ok {
		return e.applyFunction(funcExpr, args, ctx)
	}

	// Extract function name from evaluated head
	headName, ok := core.ExtractSymbol(evaluatedHead)
	if !ok {

		// Head is not a symbol, return unevaluated
		return list
	}

	// Apply attribute transformations before evaluation
	transformedList := e.applyAttributeTransformations(headName, list, ctx)

	if !transformedList.Equal(list) {
		// The list was transformed, re-evaluate it
		return e.evaluateList(ctx, transformedList)
	}
	// Handle OneIdentity attribute specially - it can return a non-List
	if ctx.symbolTable.HasAttribute(headName, OneIdentity) && len(list.Elements) == 2 {
		// OneIdentity: f(x) = x
		result := e.evaluate(ctx, list.Elements[1])
		return result
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
	if result, found := ctx.functionRegistry.CallFunction(callExpr, ctx, e); found {
		// Check if result is an error and needs stack trace
		if core.IsError(result) {
			if errorExpr, ok := result.(*core.ErrorExpr); ok {
				// Add stack frame for this function call
				if err := ctx.stack.Push(headName, callExpr); err == nil {
					ctx.stack.Pop() // Immediately pop since we're just adding to trace
					return core.NewErrorExprWithStack(errorExpr.ErrorType, errorExpr.Message, errorExpr.Args, ctx.stack.GetFrames())
				}
			}
		}
		/*
			// Re-evaluate function results until fixed point for proper symbolic computation
			// Only re-evaluate non-atomic expressions to avoid infinite recursion
			if !result.IsAtom() && !result.Equal(callExpr) {
				return e.evaluateToFixedPoint(ctx, result)
			}
		*/
		return result
	}

	// No pattern matched, return the unevaluated expression
	return callExpr
}

// evaluateToFixedPoint continues evaluating an expression until it reaches a fixed point
// (no more changes occur) or until a maximum number of iterations to prevent infinite loops
func (e *Evaluator) evaluateToFixedPoint(ctx *Context, expr core.Expr) core.Expr {
	// TODO get value from config
	const maxIterations = 100 // Prevent infinite loops
	current := expr

	for i := 0; i < maxIterations; i++ {
		next := e.evaluateOnce(ctx, current)

		// If the result is atomic, we can't evaluate further
		if next.IsAtom() {
			return next
		}
		// Check if we've reached a fixed point (no more changes)
		if next.Equal(current) {
			return next
		}

		// Check for errors
		if core.IsError(next) {
			return next
		}

		current = next
	}

	// If we've hit the iteration limit, return what we have
	// This prevents infinite loops while still allowing significant evaluation
	// TODO WRONG
	return current
}

// evaluateArguments evaluates arguments based on hold attributes
func (e *Evaluator) evaluateArguments(headName string, args []core.Expr, ctx *Context) []core.Expr {
	evaluatedArgs := make([]core.Expr, len(args))

	// TODO -- one lookup
	holdAll := ctx.symbolTable.HasAttribute(headName, HoldAll)
	holdFirst := ctx.symbolTable.HasAttribute(headName, HoldFirst)
	holdRest := ctx.symbolTable.HasAttribute(headName, HoldRest)

	for i, arg := range args {
		if holdAll || (holdFirst && i == 0) || (holdRest && i > 0) {
			evaluatedArgs[i] = arg // Don't evaluate
		} else {
			evaluatedArgs[i] = e.evaluate(ctx, arg)
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
	if len(list.Elements) <= 2 {
		// Not enough elements to sort (need at least head + 2 args)
		return list
	}

	head := list.Elements[0]
	args := make([]core.Expr, len(list.Elements)-1)
	copy(args, list.Elements[1:])

	// Sort arguments using canonical ordering
	sort.Slice(args, func(i, j int) bool {
		return core.CanonicalCompare(args[i], args[j])
	})

	// Reconstruct the list with sorted arguments
	resultElements := make([]core.Expr, len(list.Elements))
	resultElements[0] = head
	copy(resultElements[1:], args)

	return core.List{Elements: resultElements}
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
	// Special forms that are not yet moved to builtins (complex implementation)
	case "Function":
		return e.evaluateFunction(args, ctx)
	case "SliceRange":
		return e.evaluateSliceRange(args, ctx)
	case "TakeFrom":
		return e.evaluateTakeFrom(args, ctx)
	case "PartSet":
		return e.evaluatePartSet(args, ctx)
	case "SliceSet":
		return e.evaluateSliceSet(args, ctx)
	default:
		return nil // Not a special form, or handled by pattern-based system
	}
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

// extractImmediateSlotNumbers extracts only immediate slot numbers (not nested in Function calls)
func (e *Evaluator) extractImmediateSlotNumbers(expr core.Expr) []int {
	var slots []int
	slotSet := make(map[int]bool)

	e.extractImmediateSlotsRecursive(expr, slotSet)

	// Convert set to sorted slice
	for slot := range slotSet {
		slots = append(slots, slot)
	}

	// Sort slots for consistent ordering
	for i := 0; i < len(slots); i++ {
		for j := i + 1; j < len(slots); j++ {
			if slots[i] > slots[j] {
				slots[i], slots[j] = slots[j], slots[i]
			}
		}
	}

	return slots
}

// extractImmediateSlotsRecursive extracts only immediate slots, stops at nested Functions
func (e *Evaluator) extractImmediateSlotsRecursive(expr core.Expr, slotSet map[int]bool) {
	switch exprTyped := expr.(type) {
	case core.Symbol:
		symbolName := string(exprTyped)
		if len(symbolName) >= 1 && symbolName[0] == '$' {
			// Parse slot number
			slotStr := symbolName[1:]
			if slotStr == "" {
				// Bare $ is slot 1
				slotSet[1] = true
			} else {
				// Parse number
				slotNum := 0
				for _, ch := range slotStr {
					if ch >= '0' && ch <= '9' {
						slotNum = slotNum*10 + int(ch-'0')
					} else {
						// Not a pure number, ignore (e.g., $name)
						return
					}
				}
				if slotNum > 0 {
					slotSet[slotNum] = true
				}
			}
		}
	case core.List:
		// Check if this is a Function call - if so, don't recurse into it
		if len(exprTyped.Elements) > 0 {
			if head, isSymbol := exprTyped.Elements[0].(core.Symbol); isSymbol && string(head) == "Function" {
				// This is a nested Function call - don't extract slots from it
				return
			}
		}

		// For non-Function lists, recurse into all elements
		for _, elem := range exprTyped.Elements {
			e.extractImmediateSlotsRecursive(elem, slotSet)
		}
	case core.FunctionExpr:
		// Don't recurse into nested FunctionExpr bodies - they have their own slots
	}
}

// partiallyEvaluateForFunction evaluates nested Function calls but preserves slot variables
func (e *Evaluator) partiallyEvaluateForFunction(expr core.Expr, ctx *Context) core.Expr {
	switch exprTyped := expr.(type) {
	case core.Symbol:
		symbolName := string(exprTyped)
		// Preserve slot variables ($, $1, $2, etc.)
		if len(symbolName) >= 1 && symbolName[0] == '$' {
			return expr // Don't evaluate slot variables
		}
		// Evaluate other symbols normally
		return e.evaluate(ctx, expr)
	case core.List:
		// Check if this is a Function call - if so, evaluate it
		if len(exprTyped.Elements) > 0 {
			if head, isSymbol := exprTyped.Elements[0].(core.Symbol); isSymbol && string(head) == "Function" {
				// This is a nested Function call - evaluate it
				return e.evaluate(ctx, expr)
			}
		}

		// For other lists, recursively partially evaluate elements
		newElements := make([]core.Expr, len(exprTyped.Elements))
		for i, elem := range exprTyped.Elements {
			newElements[i] = e.partiallyEvaluateForFunction(elem, ctx)
		}
		return core.List{Elements: newElements}
	default:
		// For other expression types, evaluate normally
		return e.evaluate(ctx, expr)
	}
}

// evaluateFunction implements the Function special form
// Function(x, body) or Function([x, y], body)
// Also handles slot-based functions: Function($1 + $2)
func (e *Evaluator) evaluateFunction(args []core.Expr, ctx *Context) core.Expr {
	if len(args) < 1 || len(args) >  2 {
		return core.NewErrorExpr("ArgumentError", "Function requires 1 or 2 arguments", args)
	}
	if len(args) == 1 {
		// Single argument: could be slot-based function like Function($1 + $2) or constant function like Function(42)
		body := args[0]

		// Partially evaluate the body to handle nested Function calls while preserving slots
		body = e.partiallyEvaluateForFunction(body, ctx)

		slots := e.extractImmediateSlotNumbers(body)

		if len(slots) == 0 {
			// No slot variables: this is a constant function like Function(42)
			return core.NewFunction([]string{}, body)
		}

		// Generate regular parameter names for all slots up to highest number
		maxSlot := slots[len(slots)-1] // slots are sorted
		var parameters []string
		for i := 1; i <= maxSlot; i++ {
			parameters = append(parameters, fmt.Sprintf("slot%d", i))
		}

		return core.NewFunction(parameters, body)
	}


	// First argument: parameters (held unevaluated)
	paramArg := args[0]
	body := args[1] // Body is held unevaluated

	var parameters []string

	// Parse parameters: either Symbol or List of Symbols
	if paramList, ok := paramArg.(core.List); ok {
		// Multiple parameters: Function([x, y], body) or zero parameters: Function([], body)

		// Skip the "List" head, process actual parameters
		for i := 1; i < len(paramList.Elements); i++ {
			if paramName, ok := core.ExtractSymbol(paramList.Elements[i]); ok {
				parameters = append(parameters, paramName)
			} else {
				return core.NewErrorExpr("ArgumentError", "Function parameters must be symbols", args)
			}
		}
	} else if paramName, ok := core.ExtractSymbol(paramArg); ok {
		// Single parameter: Function(x, body)
		parameters = []string{paramName}
	} else {
		return core.NewErrorExpr("ArgumentError", "Function parameters must be symbols or list of symbols", args)
	}

	// Create and return the FunctionExpr
	return core.NewFunction(parameters, body)
}

// applyFunction applies a FunctionExpr to arguments
func (e *Evaluator) applyFunction(funcExpr core.FunctionExpr, args []core.Expr, ctx *Context) core.Expr {
	// Evaluate all arguments first
	evaluatedArgs := make([]core.Expr, len(args))
	for i, arg := range args {
		evaluatedArgs[i] = e.evaluate(ctx, arg)
		// Check for errors in arguments
		if core.IsError(evaluatedArgs[i]) {
			return evaluatedArgs[i]
		}
	}

	// Check argument count
	if len(args) != len(funcExpr.Parameters) {
		return core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("Function expects %d arguments, got %d", len(funcExpr.Parameters), len(args)),
			args)
	}

	// Check if this is a slot-based function (parameters named like "slot1", "slot2")
	isSlotBased := len(funcExpr.Parameters) > 0 && len(funcExpr.Parameters[0]) > 4 && funcExpr.Parameters[0][:4] == "slot"

	if isSlotBased {
		// For slot-based functions, substitute $1, $2, etc. directly in the body
		substitutedBody := e.substituteSlots(funcExpr.Body, evaluatedArgs)
		// Evaluate the substituted body in the original context (no new bindings needed)
		return e.evaluate(ctx, substitutedBody)
	} else {
		// Create a new child context for regular function evaluation
		funcCtx := NewChildContext(ctx)

		// Bind parameters to arguments in the function context
		for i, param := range funcExpr.Parameters {
			funcCtx.AddScopedVar(param) // Keep parameter local to function
			if err := funcCtx.Set(param, evaluatedArgs[i]); err != nil {
				return core.NewErrorExpr("BindingError", err.Error(), args)
			}
		}

		// Evaluate the function body in the new context
		return e.evaluate(funcCtx, funcExpr.Body)
	}
}

// substituteSlots replaces slot variables ($1, $2, etc.) with corresponding argument values
func (e *Evaluator) substituteSlots(expr core.Expr, args []core.Expr) core.Expr {
	switch exprTyped := expr.(type) {
	case core.Symbol:
		symbolName := string(exprTyped)
		if len(symbolName) >= 1 && symbolName[0] == '$' {
			// Parse slot number
			slotStr := symbolName[1:]
			var slotNum int
			if slotStr == "" {
				// Bare $ is slot 1
				slotNum = 1
			} else {
				// Parse number
				for _, ch := range slotStr {
					if ch >= '0' && ch <= '9' {
						slotNum = slotNum*10 + int(ch-'0')
					} else {
						// Not a pure number, return as-is (e.g., $name)
						return expr
					}
				}
			}
			// Replace with corresponding argument (1-indexed)
			if slotNum >= 1 && slotNum <= len(args) {
				return args[slotNum-1]
			}
		}
		return expr
	case core.List:
		// Recursively substitute in all elements
		newElements := make([]core.Expr, len(exprTyped.Elements))
		for i, elem := range exprTyped.Elements {
			newElements[i] = e.substituteSlots(elem, args)
		}
		return core.List{Elements: newElements}
	case core.FunctionExpr:
		// Recursively substitute in function body
		newBody := e.substituteSlots(exprTyped.Body, args)
		return core.FunctionExpr{
			Parameters: exprTyped.Parameters,
			Body:       newBody,
		}
	default:
		// For other types (numbers, strings, etc.), return as-is
		return expr
	}
}

// evaluateSliceRange implements slice range syntax: expr[start:end]
func (e *Evaluator) evaluateSliceRange(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 3 {
		return core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("SliceRange expects 3 arguments (expr, start, end), got %d", len(args)), args)
	}

	// Evaluate the expression being sliced
	expr := e.evaluate(ctx, args[0])
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
	startExpr := e.evaluate(ctx, args[1])
	if core.IsError(startExpr) {
		return startExpr
	}

	endExpr := e.evaluate(ctx, args[2])
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
	expr := e.evaluate(ctx, args[0])
	if core.IsError(expr) {
		return expr
	}

	// Evaluate start index
	startExpr := e.evaluate(ctx, args[1])
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
		return e.evaluate(ctx, core.NewList("Take", expr, core.NewInteger(start)))
	} else {
		// Positive start: use Drop to remove first (start-1) elements
		// Drop([1,2,3,4,5], 2) gives [3,4,5] (for start=3, 1-indexed)
		dropCount := start - 1
		return e.evaluate(ctx, core.NewList("Drop", expr, core.NewInteger(dropCount)))
	}
}

// evaluatePartSet implements slice assignment syntax: expr[index] = value
func (e *Evaluator) evaluatePartSet(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 3 {
		return core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("PartSet expects 3 arguments (expr, index, value), got %d", len(args)), args)
	}

	// Evaluate the expression being modified
	expr := e.evaluate(ctx, args[0])
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
	indexExpr := e.evaluate(ctx, args[1])
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
	value := e.evaluate(ctx, args[2])
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
	expr := e.evaluate(ctx, args[0])
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
	startExpr := e.evaluate(ctx, args[1])
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
	endExpr := e.evaluate(ctx, args[2])
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
	value := e.evaluate(ctx, args[3])
	if core.IsError(value) {
		return value
	}

	// Use the Sliceable interface to perform the slice assignment
	return sliceable.SetSlice(start, end, value)
}
