package sexpr

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// BuiltinFunc represents a builtin function signature
type BuiltinFunc func([]Expr) Expr

// EvaluationStack represents the current evaluation call stack
type EvaluationStack struct {
	frames []StackFrame
	depth  int
	maxDepth int
}

// NewEvaluationStack creates a new evaluation stack with the given maximum depth
func NewEvaluationStack(maxDepth int) *EvaluationStack {
	return &EvaluationStack{
		frames:   make([]StackFrame, 0, maxDepth),
		depth:    0,
		maxDepth: maxDepth,
	}
}

// Push adds a new frame to the stack and checks for recursion limits
func (s *EvaluationStack) Push(function, expression string) error {
	if s.depth >= s.maxDepth {
		return fmt.Errorf("maximum recursion depth exceeded: %d", s.maxDepth)
	}
	
	frame := StackFrame{
		Function:   function,
		Expression: expression,
		Location:   "", // Can be set later if needed
	}
	
	s.frames = append(s.frames, frame)
	s.depth++
	return nil
}

// Pop removes the top frame from the stack
func (s *EvaluationStack) Pop() {
	if s.depth > 0 {
		s.frames = s.frames[:len(s.frames)-1]
		s.depth--
	}
}

// GetFrames returns a copy of the current stack frames
func (s *EvaluationStack) GetFrames() []StackFrame {
	frames := make([]StackFrame, len(s.frames))
	copy(frames, s.frames)
	return frames
}

// Depth returns the current stack depth
func (s *EvaluationStack) Depth() int {
	return s.depth
}

// Context represents the evaluation context with variable bindings and symbol attributes
type Context struct {
	variables   map[string]Expr
	parent      *Context
	symbolTable *SymbolTable
	builtins    map[string]BuiltinFunc
	stack       *EvaluationStack
}

// NewContext creates a new evaluation context
func NewContext() *Context {
	return &Context{
		variables:   make(map[string]Expr),
		parent:      nil,
		symbolTable: NewSymbolTable(),
		builtins:    getBuiltinFunctions(),
		stack:       NewEvaluationStack(1000), // Default max depth of 1000
	}
}

// NewChildContext creates a child context with a parent
func NewChildContext(parent *Context) *Context {
	return &Context{
		variables:   make(map[string]Expr),
		parent:      parent,
		symbolTable: parent.symbolTable, // Share symbol table with parent
		builtins:    parent.builtins,    // Share builtins with parent
		stack:       parent.stack,       // Share evaluation stack with parent
	}
}

// Set sets a variable in the context
func (c *Context) Set(name string, value Expr) {
	c.variables[name] = value
}

// Get retrieves a variable from the context (searches up the parent chain)
func (c *Context) Get(name string) (Expr, bool) {
	if value, ok := c.variables[name]; ok {
		return value, true
	}
	if c.parent != nil {
		return c.parent.Get(name)
	}
	return nil, false
}

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

	switch ex := expr.(type) {
	case *Atom:
		return e.evaluateAtom(ex, ctx)
	case *List:
		return e.evaluateList(ex, ctx)
	default:
		return expr
	}
}

// evaluateAtom evaluates an atomic expression
func (e *Evaluator) evaluateAtom(atom *Atom, ctx *Context) Expr {
	switch atom.AtomType {
	case SymbolAtom:
		symbolName := atom.Value.(string)
		
		// Check for built-in constants
		if value, ok := e.getBuiltinConstant(symbolName); ok {
			return value
		}
		
		// Look up variable in context
		if value, ok := ctx.Get(symbolName); ok {
			return value
		}
		
		// Return the symbol unchanged if not found
		return atom
	default:
		// Other atoms (numbers, strings, booleans) evaluate to themselves
		return atom
	}
}

// evaluateList evaluates a list expression
func (e *Evaluator) evaluateList(list *List, ctx *Context) Expr {
	if len(list.Elements) == 0 {
		return list
	}

	// Get the head of the expression
	head := list.Elements[0]
	
	// If head is not a symbol, evaluate it first
	if atom, ok := head.(*Atom); !ok || atom.AtomType != SymbolAtom {
		evaluatedHead := e.evaluate(head, ctx)
		newList := &List{Elements: make([]Expr, len(list.Elements))}
		newList.Elements[0] = evaluatedHead
		copy(newList.Elements[1:], list.Elements[1:])
		return e.evaluateList(newList, ctx)
	}

	headAtom := head.(*Atom)
	headName := headAtom.Value.(string)
	args := list.Elements[1:]

	// Handle special forms that need special evaluation semantics
	if result := e.evaluateSpecialForm(headName, args, ctx); result != nil {
		return result
	}

	// Evaluate arguments based on hold attributes
	evaluatedArgs := e.evaluateArguments(headName, args, ctx)

	// Create new list with evaluated arguments
	newList := &List{Elements: make([]Expr, len(evaluatedArgs)+1)}
	newList.Elements[0] = headAtom
	copy(newList.Elements[1:], evaluatedArgs)

	// Apply transformations based on attributes
	transformed := e.applyAttributeTransformations(headName, newList, ctx)

	// Check if OneIdentity reduced it to a single element
	// Only apply this if we originally had arguments (more than just the head)
	if len(args) > 0 && len(transformed.Elements) == 1 {
		return transformed.Elements[0]
	}

	// Try to evaluate as a built-in function
	if result := e.evaluateBuiltin(headName, transformed.Elements[1:], ctx); result != nil {
		return result
	}

	return transformed
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
func (e *Evaluator) applyAttributeTransformations(headName string, list *List, ctx *Context) *List {
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

// applyFlat flattens nested applications of the same operator
func (e *Evaluator) applyFlat(headName string, list *List) *List {
	if len(list.Elements) <= 1 {
		return list
	}
	
	var newElements []Expr
	newElements = append(newElements, list.Elements[0]) // Keep the head
	
	for _, arg := range list.Elements[1:] {
		if argList, ok := arg.(*List); ok && len(argList.Elements) > 0 {
			if argHead, ok := argList.Elements[0].(*Atom); ok &&
				argHead.AtomType == SymbolAtom &&
				argHead.Value.(string) == headName {
				// Flatten: Plus[a, Plus[b, c]] -> Plus[a, b, c]
				newElements = append(newElements, argList.Elements[1:]...)
			} else {
				newElements = append(newElements, arg)
			}
		} else {
			newElements = append(newElements, arg)
		}
	}
	
	return &List{Elements: newElements}
}

// applyOrderless sorts arguments for commutative operations
func (e *Evaluator) applyOrderless(list *List) *List {
	if len(list.Elements) <= 2 {
		return list
	}
	
	// Sort arguments by their string representation for consistent ordering
	args := make([]Expr, len(list.Elements)-1)
	copy(args, list.Elements[1:])
	
	sort.Slice(args, func(i, j int) bool {
		return args[i].String() < args[j].String()
	})
	
	newElements := make([]Expr, len(list.Elements))
	newElements[0] = list.Elements[0]
	copy(newElements[1:], args)
	
	return &List{Elements: newElements}
}

// applyOneIdentity handles f[x] -> x for functions with OneIdentity
func (e *Evaluator) applyOneIdentity(list *List) *List {
	if len(list.Elements) == 2 {
		// f[x] -> x when f has OneIdentity
		// Return the single argument directly (not wrapped in a list)
		return &List{Elements: []Expr{list.Elements[1]}}
	}
	return list
}

// evaluateSpecialForm handles special forms that need custom evaluation
func (e *Evaluator) evaluateSpecialForm(headName string, args []Expr, ctx *Context) Expr {
	switch headName {
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
	case "If":
		return e.evaluateIf(args, ctx)
	case "And":
		return e.evaluateAnd(args, ctx)
	case "Or":
		return e.evaluateOr(args, ctx)
	}
	return nil
}

// evaluateBuiltin evaluates built-in functions
func (e *Evaluator) evaluateBuiltin(headName string, args []Expr, ctx *Context) Expr {
	// Propagate errors from arguments first
	for _, arg := range args {
		if IsError(arg) {
			return arg
		}
	}
	
	// Look up function in context's builtin registry
	if fn, exists := ctx.builtins[headName]; exists {
		// Push function call to stack
		argsStr := ""
		if len(args) > 0 {
			argStrs := make([]string, len(args))
			for i, arg := range args {
				argStrs[i] = arg.String()
			}
			argsStr = fmt.Sprintf("[%s]", strings.Join(argStrs, ", "))
		} else {
			argsStr = "[]"
		}
		
		funcCallStr := headName + argsStr
		if err := ctx.stack.Push(headName, funcCallStr); err != nil {
			return NewErrorExprWithStack("RecursionError", err.Error(), args, ctx.stack.GetFrames())
		}
		defer ctx.stack.Pop()
		
		// Execute the function
		result := fn(args)
		
		// If the result is an error and doesn't have a stack trace, add one
		if errorExpr, ok := result.(*ErrorExpr); ok && len(errorExpr.StackTrace) == 0 {
			errorExpr.StackTrace = ctx.stack.GetFrames()
		}
		
		return result
	}
	
	return nil
}

// getBuiltinConstant returns built-in mathematical constants
func (e *Evaluator) getBuiltinConstant(name string) (Expr, bool) {
	switch name {
	case "Pi":
		return NewFloatAtom(math.Pi), true
	case "E":
		return NewFloatAtom(math.E), true
	case "True":
		return NewBoolAtom(true), true
	case "False":
		return NewBoolAtom(false), true
	}
	return nil, false
}

// Utility functions for numeric operations

// isNumeric checks if an expression is numeric
func isNumeric(expr Expr) bool {
	if atom, ok := expr.(*Atom); ok {
		return atom.AtomType == IntAtom || atom.AtomType == FloatAtom
	}
	return false
}

// getNumericValue extracts numeric value from an expression
func getNumericValue(expr Expr) (float64, bool) {
	if atom, ok := expr.(*Atom); ok {
		switch atom.AtomType {
		case IntAtom:
			return float64(atom.Value.(int)), true
		case FloatAtom:
			return atom.Value.(float64), true
		}
	}
	return 0, false
}

// createNumericResult creates the appropriate numeric atom
func createNumericResult(value float64) Expr {
	if value == float64(int(value)) {
		return NewIntAtom(int(value))
	}
	return NewFloatAtom(value)
}

// isBool checks if an expression is a boolean
func isBool(expr Expr) bool {
	if atom, ok := expr.(*Atom); ok {
		return atom.AtomType == BoolAtom
	}
	return false
}

// getBoolValue extracts boolean value from an expression
func getBoolValue(expr Expr) (bool, bool) {
	if atom, ok := expr.(*Atom); ok && atom.AtomType == BoolAtom {
		return atom.Value.(bool), true
	}
	return false, false
}

// isSymbol checks if an expression is a symbol
func isSymbol(expr Expr) bool {
	if atom, ok := expr.(*Atom); ok {
		return atom.AtomType == SymbolAtom
	}
	return false
}

// getSymbolName extracts symbol name from an expression
func getSymbolName(expr Expr) (string, bool) {
	if atom, ok := expr.(*Atom); ok && atom.AtomType == SymbolAtom {
		return atom.Value.(string), true
	}
	return "", false
}

// RegisterBuiltin registers a new builtin function in the context
func (c *Context) RegisterBuiltin(name string, fn BuiltinFunc) {
	c.builtins[name] = fn
}

// HasBuiltin checks if a builtin function is registered
func (c *Context) HasBuiltin(name string) bool {
	_, exists := c.builtins[name]
	return exists
}

// GetBuiltin retrieves a builtin function by name
func (c *Context) GetBuiltin(name string) (BuiltinFunc, bool) {
	fn, exists := c.builtins[name]
	return fn, exists
}

// SetStack sets the evaluation stack for the context
func (c *Context) SetStack(stack *EvaluationStack) {
	c.stack = stack
}

// GetContext returns the evaluator's context
func (e *Evaluator) GetContext() *Context {
	return e.context
}

// getBuiltinFunctions returns the standard builtin function registry
func getBuiltinFunctions() map[string]BuiltinFunc {
	return map[string]BuiltinFunc{
		// Arithmetic operations
		"Plus":     EvaluatePlus,
		"Times":    EvaluateTimes,
		"Subtract": EvaluateSubtract,
		"Divide":   EvaluateDivide,
		"Power":    EvaluatePower,
		
		// Comparison operations
		"Equal":        EvaluateEqual,
		"Unequal":      EvaluateUnequal,
		"Less":         EvaluateLess,
		"Greater":      EvaluateGreater,
		"LessEqual":    EvaluateLessEqual,
		"GreaterEqual": EvaluateGreaterEqual,
		
		// Logical operations
		"Not":     EvaluateNot,
		"SameQ":   EvaluateSameQ,
		"UnsameQ": EvaluateUnsameQ,
		
		// Introspection operations
		"Head":      EvaluateHead,
		"Length":    EvaluateLength,
		
		// Predicate operations
		"ListQ":     EvaluateListQ,
		"NumberQ":   EvaluateNumberQ,
		"BooleanQ":  EvaluateBooleanQ,
		"IntegerQ":  EvaluateIntegerQ,
		"AtomQ":     EvaluateAtomQ,
		"SymbolQ":   EvaluateSymbolQ,
		"StringQ":   EvaluateStringQ,
		
		// String operations
		"StringLength": EvaluateStringLength,
		"FullForm":     EvaluateFullForm,
		
		// List access operations
		"First": EvaluateFirst,
		"Last":  EvaluateLast,
		"Rest":  EvaluateRest,
		"Most":  EvaluateMost,
		"Part":  EvaluatePart,
	}
}

// setupBuiltinAttributes sets up standard attributes for built-in functions
func setupBuiltinAttributes(symbolTable *SymbolTable) {
	// Reset attributes
	symbolTable.Reset()
	
	// Arithmetic operations
	symbolTable.SetAttributes("Plus", []Attribute{Flat, Orderless, OneIdentity})
	symbolTable.SetAttributes("Times", []Attribute{Flat, Orderless, OneIdentity})
	symbolTable.SetAttributes("Power", []Attribute{OneIdentity})
	
	// Control structures
	symbolTable.SetAttributes("Hold", []Attribute{HoldAll})
	symbolTable.SetAttributes("If", []Attribute{HoldRest})
	symbolTable.SetAttributes("While", []Attribute{HoldAll})
	symbolTable.SetAttributes("CompoundExpression", []Attribute{HoldAll})
	symbolTable.SetAttributes("Module", []Attribute{HoldAll})
	symbolTable.SetAttributes("Block", []Attribute{HoldAll})
	
	// Assignment operations
	symbolTable.SetAttributes("Set", []Attribute{HoldFirst})
	symbolTable.SetAttributes("SetDelayed", []Attribute{HoldAll})
	symbolTable.SetAttributes("Unset", []Attribute{HoldFirst})
	
	// Logical operations
	symbolTable.SetAttributes("And", []Attribute{Flat, Orderless, HoldAll})
	symbolTable.SetAttributes("Or", []Attribute{Flat, Orderless, HoldAll})
	
	// Constants
	symbolTable.SetAttributes("Pi", []Attribute{Constant, Protected})
	symbolTable.SetAttributes("E", []Attribute{Constant, Protected})
	symbolTable.SetAttributes("True", []Attribute{Constant, Protected})
	symbolTable.SetAttributes("False", []Attribute{Constant, Protected})
}