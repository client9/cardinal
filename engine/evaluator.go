package engine

import (
	"fmt"
	"sort"

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

// GetContext returns the evaluator's current context
func (e *Evaluator) GetContext() *Context {
	return e.context
}

// Evaluate evaluates an expression in the current context
func (e *Evaluator) Evaluate(expr core.Expr) core.Expr {
	return e.evaluate(e.context, expr)
}

// evaluate is the main evaluation function
func (e *Evaluator) evaluate(_ *Context, expr core.Expr) core.Expr {
	ctx := e.context
	// Push current expression to stack for recursion tracking

	if err := ctx.stack.Push("evaluate", expr); err != nil {
		// Return recursion error with stack trace
		return core.NewErrorExprWithStack("RecursionError", err.Error(), []core.Expr{expr}, ctx.stack.GetFrames())
	}
	defer ctx.stack.Pop()
	return e.evaluateToFixedPoint(e.context, expr)
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
func (e *Evaluator) evaluateList(c *Context, list core.List) core.Expr {

	if len(list.Elements) == 0 {
		return list
	}

	// Get the head (function name)
	head := list.Elements[0]
	args := list.Elements[1:]

	// Evaluate the head to get the function name
	evaluatedHead := e.evaluate(c, head)

	// Check if head is an error - propagate it
	if core.IsError(evaluatedHead) {
		return evaluatedHead
	}

	// Check if head is a function expression (function application)
	if funcExpr, ok := evaluatedHead.(core.FunctionExpr); ok {

		//log.Printf("Converted Head to funcExpr: vars=%s, body=%s", funcExpr.Parameters, funcExpr.Body)
		return e.applyFunction(c, funcExpr, args)
	}
	/*
		if evaluatedHead.Head() == "Function" {
			//log.Printf("Got unparsed Function")
			fn := e.evaluate(ctx, evaluatedHead)
			// log.Printf("Reparsed Function: %v", fn)
			if funcExpr, ok := fn.(core.FunctionExpr); ok {
				return e.applyFunction(funcExpr, args, ctx)
			}
			log.Printf("FAIL123")
		}
	*/

	// Extract function name from evaluated head
	headName, ok := core.ExtractSymbol(evaluatedHead)
	if !ok {

		// Head is not a symbol, return unevaluated
		return list
	}

	// Apply attribute transformations before evaluation
	transformedList := e.applyAttributeTransformations(headName, list, c)

	if !transformedList.Equal(list) {
		// The list was transformed, re-evaluate it
		return e.evaluateList(c, transformedList)
	}
	// Handle OneIdentity attribute specially - it can return a non-List
	if c.symbolTable.HasAttribute(headName, OneIdentity) && len(list.Elements) == 2 {
		// OneIdentity: f(x) = x
		result := e.evaluate(c, list.Elements[1])
		return result
	}

	// Check for special forms first (these don't follow normal evaluation rules)
	if specialResult := e.evaluateSpecialForm(headName, args, c); specialResult != nil {
		return specialResult
	}
	// Try pattern-based function resolution
	return e.evaluatePatternFunction(headName, args, c)
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
	/*
		case "Function":
			return e.evaluateFunction(args, ctx)
	*/
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

// applyFunction applies a FunctionExpr to arguments
func (e *Evaluator) applyFunction(c *Context, funcExpr core.FunctionExpr, args []core.Expr) core.Expr {
	// Evaluate all arguments first
	evaluatedArgs := make([]core.Expr, len(args))
	for i, arg := range args {
		evaluatedArgs[i] = e.evaluate(c, arg)
		if core.IsError(evaluatedArgs[i]) {
			return evaluatedArgs[i]
		}
	}

	rules := make([]core.Expr, len(args))

	if funcExpr.Parameters == nil {
		// Anonymous
		for i := 0; i < len(args); i++ {
			name := core.NewSymbol(fmt.Sprintf("$%d", i+1))
			rules[i] = core.NewList("Rule", name, evaluatedArgs[i])
		}
		if len(args) > 0 {
			name := core.NewSymbol("$")
			rules = append(rules, core.NewList("Rule", name, evaluatedArgs[0]))
		}
	} else {
		// Named - Check argument count
		if len(args) != len(funcExpr.Parameters) {
			return core.NewErrorExpr("ArgumentError",
				fmt.Sprintf("Function expects %d arguments, got %d", len(funcExpr.Parameters), len(args)),
				args)
		}
		for i := 0; i < len(args); i++ {
			rules[i] = core.NewList("Rule", funcExpr.Parameters[i], evaluatedArgs[i])
		}
	}

	// create a rules list

	rlist := core.NewList("List", rules...)

	modified := functionReplaceAll(e, c, funcExpr.Body, rlist)

	result := e.Evaluate(modified)
	return result
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

func functionReplaceAll(e *Evaluator, c *Context, expr core.Expr, rule core.Expr) core.Expr {
	//log.Printf("rules %v , body %v", rule, expr)
	if fn, ok := expr.(core.FunctionExpr); ok {
		bodyOut := functionReplaceAll(e, c, fn.Body, rule)
		//log.Printf("Body in %q, body out %q", fn.Body, bodyOut)
		if bodyOut.Equal(fn.Body) {
			return expr
		}
		//log.Printf("Body changed 1:  should rewrite parameters")
		return core.NewFunction(fn.Parameters, bodyOut)
	}

	// First try to apply the rule to the current expression
	result := core.ReplaceAllWithRules(expr, rule)
	//log.Printf("Expr in: %q, Rule %q, Expr out %q", expr, rule, result)
	// If the rule matched at this level, we're done (don't recurse into replacement)
	if !result.Equal(expr) {
		return result
	}

	// If no match at this level, recursively apply to subexpressions
	if list, ok := expr.(core.List); ok && len(list.Elements) > 0 {
		// Create new list with transformed elements
		newElements := make([]core.Expr, len(list.Elements))
		changed := false

		for i, element := range list.Elements {
			newElement := functionReplaceAll(e, c, element, rule)
			newElements[i] = newElement
			if !newElement.Equal(element) {
				changed = true
			}
		}

		if changed {
			return core.NewListFromExprs(newElements...)
		}
	}

	// No changes made, return original expression
	return expr
}
