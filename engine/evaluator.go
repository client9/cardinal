package engine

import (
	"fmt"
	"sort"
	//	"log"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
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
	ctx := e.context
	if err := ctx.stack.Push("evaluate", expr); err != nil {
		return core.NewError("RecursionError", err.Error()).SetCaller(expr)
	}
	defer ctx.stack.Pop()
	result := e.evaluateToFixedPoint(e.context, expr)
	if err, ok := core.AsError(result); ok {
		return err.Wrap(expr)
	}
	return result
}

// evaluateToFixedPoint continues evaluating an expression until it reaches a fixed point
// (no more changes occur) or until a maximum number of iterations to prevent infinite loops
func (e *Evaluator) evaluateToFixedPoint(ctx *Context, expr core.Expr) core.Expr {
	next := e.evaluateExpr(ctx, expr)
	if core.IsError(next) {
		return next
	}

	// If the result is atomic, we can't evaluate further
	if next.IsAtom() {
		return next
	}
	// Check if we've reached a fixed point (no more changes)
	if next.Equal(expr) {
		return next
	}

	return e.Evaluate(next)
}

func (e *Evaluator) evaluateExpr(ctx *Context, expr core.Expr) core.Expr {
	switch ex := expr.(type) {
	case core.Symbol:
		// Check for variable binding first
		if value, ok := ctx.Get(ex); ok {
			return value
		}
		// Return the symbol itself if not bound
		return ex
	case core.List:
		result := e.evaluateList(ctx, ex)

		// downstream doesn't have access to the original
		// expression so fill it in here
		if err, ok := core.AsError(result); ok {
			if err.Arg == nil {
				err.Arg = expr
			}
			return err
		}
		return result
	default:
		// All other types (ByteArray, Association, ErrorExpr, etc.) evaluate to themselves
		return expr
	}
}

// evaluateList evaluates a list expression
func (e *Evaluator) evaluateList(c *Context, list core.List) core.Expr {
	// Get the head (function name)
	head := list.Head()
	args := list.Tail()

	// Evaluate the head to get the function name
	evaluatedHead := e.Evaluate(head)
	if _, ok := core.AsError(evaluatedHead); ok {
		return list
	}

	// Check if head is a function expression (function application)
	if funcExpr, ok := evaluatedHead.(core.FunctionExpr); ok {
		return e.applyFunction(c, funcExpr, args)
	}

	// Extract function name from evaluated head
	headName, ok := evaluatedHead.(core.Symbol)
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
	if c.symbolTable.HasAttribute(headName, OneIdentity) && list.Length() == 1 {
		// OneIdentity: f(x) = x
		args := list.Tail()
		result := e.Evaluate(args[0])
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
func (e *Evaluator) evaluatePatternFunction(headName core.Symbol, args []core.Expr, ctx *Context) core.Expr {

	// Evaluate arguments based on hold attributes
	evaluatedArgs := e.evaluateArguments(headName, args, ctx)

	// Check for errors in evaluated arguments
	for _, arg := range evaluatedArgs {
		if core.IsError(arg) {
			return arg
		}
	}

	// Create the function call expression for pattern matching
	callExpr := core.ListFrom(headName, evaluatedArgs...)

	// Try to find a matching pattern in the function registry
	if result, found := ctx.functionRegistry.CallFunction(callExpr, ctx, e); found {
		return result
	}

	// No pattern matched, return the unevaluated expression
	return callExpr
}

// evaluateArguments evaluates arguments based on hold attributes
func (e *Evaluator) evaluateArguments(headName core.Symbol, args []core.Expr, ctx *Context) []core.Expr {
	evaluatedArgs := make([]core.Expr, len(args))

	// TODO -- one lookup
	holdAll := ctx.symbolTable.HasAttribute(headName, HoldAll)
	holdFirst := ctx.symbolTable.HasAttribute(headName, HoldFirst)
	holdRest := ctx.symbolTable.HasAttribute(headName, HoldRest)

	for i, arg := range args {
		if holdAll || (holdFirst && i == 0) || (holdRest && i > 0) {
			evaluatedArgs[i] = arg // Don't evaluate
		} else {
			evaluatedArgs[i] = e.Evaluate(arg)
		}
	}

	return evaluatedArgs
}

// applyAttributeTransformations applies attribute-based transformations
func (e *Evaluator) applyAttributeTransformations(headName core.Symbol, list core.List, ctx *Context) core.List {
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
func (e *Evaluator) applyFlat(head core.Symbol, list core.List) core.List {
	if list.Length() == 0 {
		return list
	}

	listhead := list.Head()
	args := list.Tail()

	newArgs := []core.Expr{}

	for _, arg := range args {
		// If the argument is the same function, flatten it
		if argList, ok := arg.(core.List); ok {
			if argList.Head() == listhead {
				// Flatten: f(a, f(b, c), d) → f(a, b, c, d)
				newArgs = append(newArgs, argList.Tail()...)
				continue
			}
		}
		newArgs = append(newArgs, arg)
	}

	return core.ListFrom(head, newArgs...)
}

// applyOrderless implements the Orderless attribute (commutativity)
func (e *Evaluator) applyOrderless(list core.List) core.List {
	if list.Length() < 2 {
		// Not enough elements to sort (need at least head + 2 args)
		return list
	}

	head := list.Head()
	args := make([]core.Expr, list.Length())
	copy(args, list.Tail())

	// Sort arguments using canonical ordering
	sort.Slice(args, func(i, j int) bool {
		return core.CanonicalCompare(args[i], args[j])
	})

	// Reconstruct the list with sorted arguments
	resultElements := make([]core.Expr, list.Length()+1)
	resultElements[0] = head
	copy(resultElements[1:], args)

	return core.NewListFromExprs(resultElements...)
}

// applyOneIdentity implements the OneIdentity attribute
func (e *Evaluator) applyOneIdentity(list core.List) core.List {
	// OneIdentity is now handled specially in evaluateList
	// This function is kept for consistency but doesn't transform anything
	return list
}

// evaluateSpecialForm handles special forms that don't follow normal evaluation rules
// TODO
func (e *Evaluator) evaluateSpecialForm(headName core.Symbol, args []core.Expr, ctx *Context) core.Expr {
	switch headName.String() {
	// Special forms that are not yet moved to builtins (complex implementation)
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
		evaluatedArgs[i] = e.Evaluate(arg)
		if core.IsError(evaluatedArgs[i]) {
			return evaluatedArgs[i]
		}
	}

	rules := make([]core.Expr, len(args))

	if funcExpr.Parameters == nil {
		// Anonymous
		for i := 0; i < len(args); i++ {
			name := core.NewSymbol(fmt.Sprintf("$%d", i+1))
			rules[i] = core.ListFrom(symbol.Rule, name, evaluatedArgs[i])
		}
		if len(args) > 0 {
			name := core.NewSymbol("$")
			rules = append(rules, core.ListFrom(symbol.Rule, name, evaluatedArgs[0]))
		}
	} else {
		// Named - Check argument count
		if len(args) != len(funcExpr.Parameters) {
			return core.NewError(
				"ArgumentError",
				fmt.Sprintf("Function expects %d arguments, got %d",
					len(funcExpr.Parameters), len(args)))
		}
		for i := 0; i < len(args); i++ {
			rules[i] = core.ListFrom(symbol.Rule, funcExpr.Parameters[i], evaluatedArgs[i])
		}
	}

	// create a rules list

	rlist := core.NewList(symbol.List, rules...)

	modified := functionReplaceAll(e, c, funcExpr.Body, rlist)

	result := e.Evaluate(modified)
	return result
}

// evaluatePartSet implements slice assignment syntax: expr[index] = value
func (e *Evaluator) evaluatePartSet(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 3 {
		return core.NewError("ArgumentError",
			fmt.Sprintf("PartSet expects 3 arguments (expr, index, value), got %d", len(args)))
	}

	// Evaluate the expression being modified
	expr := e.Evaluate(args[0])
	if core.IsError(expr) {
		return expr
	}

	// Check if the expression is sliceable
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewError("TypeError",
			fmt.Sprintf("Expression of type %s is not sliceable", expr.String()))
	}

	// Evaluate index
	indexExpr := e.Evaluate(args[1])
	if core.IsError(indexExpr) {
		return indexExpr
	}

	// Extract integer value for index
	index, ok := core.ExtractInt64(indexExpr)
	if !ok {
		return core.NewError("TypeError",
			fmt.Sprintf("Part index must be an integer, got %s", indexExpr.String()))
	}

	// Evaluate value
	value := e.Evaluate(args[2])
	if core.IsError(value) {
		return value
	}

	// get modified list
	mlist := sliceable.SetElementAt(index, value)

	// is this variable?  i.e. a[2] = 100
	// then we need to update the variable
	if name, ok := args[0].(core.Symbol); ok {
		ctx.Set(name, mlist)

		// TODO right value
	}

	// list literal [1,2,3,4][2] = 100 just returns a new list
	return mlist
}

// evaluateSliceSet implements slice assignment syntax: expr[start:end] = value
func (e *Evaluator) evaluateSliceSet(args []core.Expr, ctx *Context) core.Expr {
	if len(args) != 4 {
		return core.NewError("ArgumentError",
			fmt.Sprintf("SliceSet expects 4 arguments (expr, start, end, value), got %d", len(args)))
	}

	// Evaluate the expression being modified
	expr := e.Evaluate(args[0])
	if core.IsError(expr) {
		return expr
	}

	// Check if the expression is sliceable
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewError("TypeError",
			fmt.Sprintf("Expression of type %s is not sliceable", expr.String()))
	}

	// Evaluate start index
	startExpr := e.Evaluate(args[1])
	if core.IsError(startExpr) {
		return startExpr
	}

	// Extract integer value for start
	start, ok := core.ExtractInt64(startExpr)
	if !ok {
		return core.NewError("TypeError",
			fmt.Sprintf("Slice start index must be an integer, got %s", startExpr.String()))
	}

	// Evaluate end index
	endExpr := e.Evaluate(args[2])
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
		return core.NewError("TypeError",
			fmt.Sprintf("Slice end index must be an integer, got %s", endExpr.String()))
	}

	// Evaluate value
	value := e.Evaluate(args[3])
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
	if list, ok := expr.(core.List); ok {
		// Create new list with transformed elements
		newElements := make([]core.Expr, list.Length()+1)
		changed := false

		for i, element := range list.AsSlice() {
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
