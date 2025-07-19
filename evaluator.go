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
	variables        map[string]Expr
	parent           *Context
	symbolTable      *SymbolTable
	builtins         map[string]BuiltinFunc // Legacy built-ins for compatibility
	functionRegistry *FunctionRegistry      // New unified pattern-based function system
	stack            *EvaluationStack
}

// NewContext creates a new evaluation context
func NewContext() *Context {
	ctx := &Context{
		variables:        make(map[string]Expr),
		parent:           nil,
		symbolTable:      NewSymbolTable(),
		builtins:         getBuiltinFunctions(), // Legacy built-ins for compatibility
		functionRegistry: NewFunctionRegistry(),
		stack:            NewEvaluationStack(1000), // Default max depth of 1000
	}
	
	// Register default built-in functions with patterns
	registerDefaultBuiltins(ctx.functionRegistry)
	
	return ctx
}

// NewChildContext creates a child context with a parent
func NewChildContext(parent *Context) *Context {
	return &Context{
		variables:        make(map[string]Expr),
		parent:           parent,
		symbolTable:      parent.symbolTable,      // Share symbol table with parent
		builtins:         parent.builtins,         // Share builtins with parent
		functionRegistry: parent.functionRegistry, // Share function registry with parent
		stack:            parent.stack,            // Share evaluation stack with parent
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
	case Atom:
		return e.evaluateAtom(ex, ctx)
	case List:
		return e.evaluateList(ex, ctx)
	default:
		return expr
	}
}

// evaluateAtom evaluates an atomic expression
func (e *Evaluator) evaluateAtom(atom Atom, ctx *Context) Expr {
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
func (e *Evaluator) evaluateList(list List, ctx *Context) Expr {
	if len(list.Elements) == 0 {
		return list
	}

	// Get the head of the expression
	head := list.Elements[0]
	
	// If head is not a symbol, evaluate it first
	if atom, ok := head.(Atom); !ok || atom.AtomType != SymbolAtom {
		evaluatedHead := e.evaluate(head, ctx)
		newList := List{Elements: make([]Expr, len(list.Elements))}
		newList.Elements[0] = evaluatedHead
		copy(newList.Elements[1:], list.Elements[1:])
		return e.evaluateList(newList, ctx)
	}

	headAtom := head.(Atom)
	headName := headAtom.Value.(string)
	args := list.Elements[1:]

	// Handle special forms that need special evaluation semantics
	if result := e.evaluateSpecialForm(headName, args, ctx); result != nil {
		return result
	}

	// Evaluate arguments based on hold attributes
	evaluatedArgs := e.evaluateArguments(headName, args, ctx)

	// Create new list with evaluated arguments
	newList := List{Elements: make([]Expr, len(evaluatedArgs)+1)}
	newList.Elements[0] = headAtom
	copy(newList.Elements[1:], evaluatedArgs)

	// Apply transformations based on attributes
	transformed := e.applyAttributeTransformations(headName, newList, ctx)

	// Check if OneIdentity reduced it to a single element
	// Only apply this if we originally had arguments (more than just the head)
	if len(args) > 0 && len(transformed.Elements) == 1 {
		return transformed.Elements[0]
	}

	// Check old built-in function overrides first (for backward compatibility)
	if builtin, exists := ctx.builtins[headName]; exists {
		args := transformed.Elements[1:]
		
		// Push function call to stack for error tracking
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
		
		// Check for errors in arguments and propagate them first
		for _, arg := range args {
			if IsError(arg) {
				// Add current function to stack trace if error doesn't have one
				if errorExpr, ok := arg.(*ErrorExpr); ok && len(errorExpr.StackTrace) == 0 {
					errorExpr.StackTrace = ctx.stack.GetFrames()
				}
				return arg
			}
		}
		
		// Execute the function
		result := builtin(args)
		
		// If the result is an error and doesn't have a stack trace, add one
		if errorExpr, ok := result.(*ErrorExpr); ok && len(errorExpr.StackTrace) == 0 {
			errorExpr.StackTrace = ctx.stack.GetFrames()
		}
		
		return result
	}

	// Try to evaluate using the unified pattern-based function system
	if result := e.evaluatePatternFunction(headName, transformed.Elements[1:], ctx); result != nil {
		return result
	}

	// Try to evaluate using user-defined functions (deprecated path for compatibility)
	if result := e.evaluateUserDefinedFunction(headName, transformed.Elements[1:], ctx); result != nil {
		return result
	}

	return transformed
}

// evaluatePatternFunction tries to evaluate a function using the unified pattern-based system
func (e *Evaluator) evaluatePatternFunction(headName string, args []Expr, ctx *Context) Expr {
	// Find the best matching function definition
	funcDef, bindings := ctx.functionRegistry.FindMatchingFunction(headName, args)
	if funcDef == nil {
		return nil
	}
	
	// Create a new child context with the pattern variable bindings
	funcCtx := NewChildContext(ctx)
	for varName, value := range bindings {
		funcCtx.Set(varName, value)
	}
	
	// Execute the function
	if funcDef.GoImpl != nil {
		// Built-in function with Go implementation
		// Push function call to stack for error tracking
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
		result := funcDef.GoImpl(args, funcCtx)
		
		// If the result is an error and doesn't have a stack trace, add one
		if errorExpr, ok := result.(*ErrorExpr); ok && len(errorExpr.StackTrace) == 0 {
			errorExpr.StackTrace = ctx.stack.GetFrames()
		}
		
		return result
	} else {
		// User-defined function with s-expression body
		return e.evaluate(funcDef.Body, funcCtx)
	}
}

// evaluateUserDefinedFunction tries to evaluate a user-defined function (DEPRECATED - kept for compatibility)
func (e *Evaluator) evaluateUserDefinedFunction(headName string, args []Expr, ctx *Context) Expr {
	// Look up the function definition
	if funcDef, exists := ctx.Get(headName); exists {
		// Check if it's a function list (multiple definitions)
		if funcList, ok := funcDef.(List); ok && len(funcList.Elements) > 0 {
			if head, ok := funcList.Elements[0].(Atom); ok && 
				head.AtomType == SymbolAtom && head.Value.(string) == "FunctionList" {
				
				// Try each function definition in order
				for _, def := range funcList.Elements[1:] {
					if result := e.tryFunctionDefinition(headName, def, args, ctx); result != nil {
						return result
					}
				}
				
				// No pattern matched, return unevaluated
				return List{Elements: append([]Expr{NewSymbolAtom(headName)}, args...)}
			}
		}
		
		// Single function definition
		if result := e.tryFunctionDefinition(headName, funcDef, args, ctx); result != nil {
			return result
		}
	}
	
	return nil
}

// tryFunctionDefinition tries to match and evaluate a single function definition
func (e *Evaluator) tryFunctionDefinition(headName string, funcDef Expr, args []Expr, ctx *Context) Expr {
	// Check if it's a function definition (Function[{params}, body])
	if funcList, ok := funcDef.(List); ok && len(funcList.Elements) == 3 {
		if head, ok := funcList.Elements[0].(Atom); ok && 
			head.AtomType == SymbolAtom && head.Value.(string) == "Function" {
			
			// Extract parameters and body
			params := funcList.Elements[1]
			body := funcList.Elements[2]
			
			// Create a new context for function evaluation
			funcCtx := NewChildContext(ctx)
			
			// Bind parameters to arguments using pattern matching
			if paramList, ok := params.(List); ok {
				// Check if any parameter is a sequence pattern
				hasSequencePattern := false
				for _, param := range paramList.Elements {
					if atom, ok := param.(Atom); ok && atom.AtomType == SymbolAtom {
						symbolName := atom.Value.(string)
						if isPatternVariable(symbolName) {
							patternInfo := parsePatternInfo(symbolName)
							if patternInfo.Type == BlankSequencePattern || patternInfo.Type == BlankNullSequencePattern {
								hasSequencePattern = true
								break
							}
						}
					}
				}
				
				if hasSequencePattern {
					// Use sequence pattern matching
					if !e.matchSequencePatterns(paramList.Elements, args, funcCtx) {
						return nil
					}
				} else {
					// Use traditional pattern matching (exact parameter count)
					if len(paramList.Elements) != len(args) {
						// Wrong number of arguments, pattern doesn't match
						return nil
					}
					
					for i, param := range paramList.Elements {
						if !e.matchPatternInternal(param, args[i], funcCtx, true) {
							// Pattern didn't match, try next definition
							return nil
						}
					}
				}
			}
			
			// All patterns matched, evaluate the function body
			return e.evaluate(body, funcCtx)
		}
	}
	
	return nil
}

// matchPattern matches a pattern against an expression and binds variables
func (e *Evaluator) matchPattern(pattern Expr, expr Expr, ctx *Context) bool {
	return e.matchPatternInternal(pattern, expr, ctx, false)
}

// matchPatternInternal matches a pattern with control over parameter binding behavior
func (e *Evaluator) matchPatternInternal(pattern Expr, expr Expr, ctx *Context, isParameterList bool) bool {
	switch pat := pattern.(type) {
	case Atom:
		if pat.AtomType == SymbolAtom {
			symbolName := pat.Value.(string)
			// Check if this is a pattern variable (ends with _)
			if isPatternVariable(symbolName) {
				// Parse the pattern variable to get complete pattern information
				patternInfo := parsePatternInfo(symbolName)
				
				// Check type constraint if present
				if !matchesType(expr, patternInfo.TypeName) {
					return false
				}
				
				if patternInfo.VarName == "" {
					// Anonymous pattern variable _ matches anything (if type constraint passed)
					return true
				}
				
				// Bind the variable to the expression
				ctx.Set(patternInfo.VarName, expr)
				return true
			} else {
				// Regular symbol behavior depends on context
				if isParameterList {
					// In parameter lists, regular symbols bind to values
					ctx.Set(symbolName, expr)
					return true
				} else {
					// In head patterns, regular symbols require exact matches
					if exprAtom, ok := expr.(Atom); ok && exprAtom.AtomType == SymbolAtom {
						return exprAtom.Value.(string) == symbolName
					}
					return false
				}
			}
		} else {
			// Literal atom - exact match required
			if exprAtom, ok := expr.(Atom); ok {
				return pat.AtomType == exprAtom.AtomType && pat.Value == exprAtom.Value
			}
			return false
		}
	case List:
		// Match structured expressions like Plus(x_, y_)
		if exprList, ok := expr.(List); ok {
			// For function parameters, we need to handle sequence patterns
			if isParameterList && len(pat.Elements) > 1 {
				// First, check that the heads match exactly (head is never a parameter)
				if !e.matchPatternInternal(pat.Elements[0], exprList.Elements[0], ctx, false) {
					return false
				}
				// Then use sequence pattern matching for the arguments
				return e.matchSequencePatterns(pat.Elements[1:], exprList.Elements[1:], ctx)
			}
			
			// For non-parameter lists, require exact length match
			if len(pat.Elements) != len(exprList.Elements) {
				return false
			}
			
			for i, patElem := range pat.Elements {
				// First element is the head (requires exact match), rest are parameters
				isParam := i > 0 && isParameterList
				if !e.matchPatternInternal(patElem, exprList.Elements[i], ctx, isParam) {
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

// matchSequencePatterns handles matching patterns that may contain sequence patterns
func (e *Evaluator) matchSequencePatterns(patterns []Expr, exprs []Expr, ctx *Context) bool {
	patternIndex := 0
	exprIndex := 0
	
	for patternIndex < len(patterns) {
		pattern := patterns[patternIndex]
		
		// Check if this pattern is a sequence pattern
		if atom, ok := pattern.(Atom); ok && atom.AtomType == SymbolAtom {
			symbolName := atom.Value.(string)
			if isPatternVariable(symbolName) {
				patternInfo := parsePatternInfo(symbolName)
				
				switch patternInfo.Type {
				case BlankPattern:
					// Single pattern: matches exactly one expression
					if exprIndex >= len(exprs) {
						return false // No more expressions to match
					}
					
					// Check type constraint
					if !matchesType(exprs[exprIndex], patternInfo.TypeName) {
						return false
					}
					
					// Bind variable if named
					if patternInfo.VarName != "" {
						ctx.Set(patternInfo.VarName, exprs[exprIndex])
					}
					
					exprIndex++
					patternIndex++
					
				case BlankSequencePattern:
					// Sequence pattern: matches one or more expressions
					if exprIndex >= len(exprs) {
						return false // Need at least one expression
					}
					
					// Calculate how many expressions this pattern should consume
					remainingPatterns := len(patterns) - patternIndex - 1
					minExprsNeeded := remainingPatterns // Minimum expressions needed for remaining patterns
					maxExprsAvailable := len(exprs) - exprIndex - minExprsNeeded
					
					if maxExprsAvailable < 1 {
						return false // Need at least one expression for BlankSequence
					}
					
					// Consume expressions for this sequence pattern
					sequenceExprs := make([]Expr, 0)
					for i := 0; i < maxExprsAvailable; i++ {
						if exprIndex >= len(exprs) {
							break
						}
						
						// Check type constraint
						if !matchesType(exprs[exprIndex], patternInfo.TypeName) {
							break
						}
						
						sequenceExprs = append(sequenceExprs, exprs[exprIndex])
						exprIndex++
					}
					
					if len(sequenceExprs) == 0 {
						return false // BlankSequence must match at least one expression
					}
					
					// Bind variable if named
					if patternInfo.VarName != "" {
						ctx.Set(patternInfo.VarName, NewList(append([]Expr{NewSymbolAtom("List")}, sequenceExprs...)...))
					}
					
					patternIndex++
					
				case BlankNullSequencePattern:
					// Null sequence pattern: matches zero or more expressions
					// Calculate how many expressions this pattern should consume
					remainingPatterns := len(patterns) - patternIndex - 1
					minExprsNeeded := remainingPatterns // Minimum expressions needed for remaining patterns
					maxExprsAvailable := len(exprs) - exprIndex - minExprsNeeded
					
					if maxExprsAvailable < 0 {
						maxExprsAvailable = 0 // Can consume zero expressions
					}
					
					// Consume expressions for this sequence pattern
					sequenceExprs := make([]Expr, 0)
					for i := 0; i < maxExprsAvailable; i++ {
						if exprIndex >= len(exprs) {
							break
						}
						
						// Check type constraint
						if !matchesType(exprs[exprIndex], patternInfo.TypeName) {
							break
						}
						
						sequenceExprs = append(sequenceExprs, exprs[exprIndex])
						exprIndex++
					}
					
					// Bind variable if named (can be empty list)
					if patternInfo.VarName != "" {
						ctx.Set(patternInfo.VarName, NewList(append([]Expr{NewSymbolAtom("List")}, sequenceExprs...)...))
					}
					
					patternIndex++
					
				default:
					return false // Unknown pattern type
				}
			} else {
				// Regular symbol pattern
				if exprIndex >= len(exprs) {
					return false
				}
				
				// Regular symbol in parameter list binds to value
				ctx.Set(symbolName, exprs[exprIndex])
				exprIndex++
				patternIndex++
			}
		} else {
			// Non-symbol pattern (literal)
			if exprIndex >= len(exprs) {
				return false
			}
			
			if !e.matchPatternInternal(pattern, exprs[exprIndex], ctx, false) {
				return false
			}
			
			exprIndex++
			patternIndex++
		}
	}
	
	// All patterns matched, check that all expressions were consumed
	return exprIndex == len(exprs)
}

// PatternType represents the type of pattern
type PatternType int

const (
	BlankPattern PatternType = iota     // _ - matches exactly one expression
	BlankSequencePattern                // __ - matches one or more expressions
	BlankNullSequencePattern            // ___ - matches zero or more expressions
)

// PatternInfo contains information about a parsed pattern
type PatternInfo struct {
	Type     PatternType
	VarName  string
	TypeName string
}

// PatternSpecificity represents how specific a pattern is (higher = more specific)
type PatternSpecificity int

const (
	SpecificityNullSequence PatternSpecificity = 1  // x___ (least specific)
	SpecificitySequence     PatternSpecificity = 2  // x__ 
	SpecificityGeneral      PatternSpecificity = 3  // x_
	SpecificityBuiltinType  PatternSpecificity = 4  // x_Integer, x_String, etc.
	SpecificityCustomType   PatternSpecificity = 5  // x_Color, x_Point, etc.
	SpecificityLiteral      PatternSpecificity = 6  // 42, "hello", exact values (most specific)
)

// getPatternSpecificity calculates the specificity of a single pattern
func getPatternSpecificity(pattern Expr) PatternSpecificity {
	switch pat := pattern.(type) {
	case Atom:
		if pat.AtomType == SymbolAtom {
			symbolName := pat.Value.(string)
			if isPatternVariable(symbolName) {
				patternInfo := parsePatternInfo(symbolName)
				return getPatternVariableSpecificity(patternInfo)
			}
			// Regular symbol (should not occur in patterns, but handle it)
			return SpecificityLiteral
		}
		// Literal atom (number, string, boolean)
		return SpecificityLiteral
	case List:
		// Structured pattern like Plus(x_, y_) - calculate based on elements
		if len(pat.Elements) == 0 {
			return SpecificityLiteral // Empty list
		}
		
		// Head must be literal for structured patterns
		headSpecificity := getPatternSpecificity(pat.Elements[0])
		if headSpecificity != SpecificityLiteral {
			// Head is not literal, this is not a valid structured pattern
			return SpecificityGeneral
		}
		
		// Specificity is based on the least specific parameter
		minSpecificity := SpecificityLiteral
		for i := 1; i < len(pat.Elements); i++ {
			paramSpecificity := getPatternSpecificity(pat.Elements[i])
			if paramSpecificity < minSpecificity {
				minSpecificity = paramSpecificity
			}
		}
		return minSpecificity
	default:
		return SpecificityGeneral
	}
}

// getPatternVariableSpecificity calculates specificity for pattern variables
func getPatternVariableSpecificity(info PatternInfo) PatternSpecificity {
	// Base specificity on pattern type
	var baseSpecificity PatternSpecificity
	switch info.Type {
	case BlankNullSequencePattern:
		baseSpecificity = SpecificityNullSequence
	case BlankSequencePattern:
		baseSpecificity = SpecificitySequence
	case BlankPattern:
		baseSpecificity = SpecificityGeneral
	default:
		baseSpecificity = SpecificityGeneral
	}
	
	// Increase specificity if there's a type constraint
	if info.TypeName != "" {
		if isBuiltinType(info.TypeName) {
			baseSpecificity = SpecificityBuiltinType
		} else {
			baseSpecificity = SpecificityCustomType
		}
	}
	
	return baseSpecificity
}

// isBuiltinType checks if a type name is a built-in type
func isBuiltinType(typeName string) bool {
	switch typeName {
	case "Integer", "Real", "Float", "Number", "Numeric", "String", "Boolean", "Bool", "Symbol", "Atom", "List":
		return true
	default:
		return false
	}
}

// getFunctionPatternSpecificity calculates overall specificity for a function pattern
func getFunctionPatternSpecificity(params Expr) PatternSpecificity {
	if paramList, ok := params.(List); ok {
		// Calculate total specificity as sum of individual parameter specificities
		// This ensures that patterns with more specific parameters rank higher
		totalSpecificity := PatternSpecificity(0)
		for _, param := range paramList.Elements {
			totalSpecificity += getPatternSpecificity(param)
		}
		return totalSpecificity
	}
	
	// Single parameter case
	return getPatternSpecificity(params)
}

// getFunctionDefinitionSpecificity extracts specificity from a Function definition
func getFunctionDefinitionSpecificity(funcDef Expr) PatternSpecificity {
	if funcList, ok := funcDef.(List); ok && len(funcList.Elements) == 3 {
		if head, ok := funcList.Elements[0].(Atom); ok && 
			head.AtomType == SymbolAtom && head.Value.(string) == "Function" {
			// Extract parameters (second element)
			params := funcList.Elements[1]
			return getFunctionPatternSpecificity(params)
		}
	}
	return SpecificityGeneral
}

// insertFunctionBySpecificity inserts a function definition in the correct position based on specificity
func insertFunctionBySpecificity(existingList List, newFuncDef Expr) List {
	newSpecificity := getFunctionDefinitionSpecificity(newFuncDef)
	
	// Create new elements slice with the new function inserted in the right position
	newElements := make([]Expr, 0, len(existingList.Elements)+1)
	newElements = append(newElements, existingList.Elements[0]) // Keep "FunctionList" head
	
	inserted := false
	for i := 1; i < len(existingList.Elements); i++ {
		existingSpecificity := getFunctionDefinitionSpecificity(existingList.Elements[i])
		
		// Insert before less specific patterns (higher specificity = more specific)
		if !inserted && newSpecificity > existingSpecificity {
			newElements = append(newElements, newFuncDef)
			inserted = true
		}
		newElements = append(newElements, existingList.Elements[i])
	}
	
	// If not inserted yet, append at end (least specific)
	if !inserted {
		newElements = append(newElements, newFuncDef)
	}
	
	return List{Elements: newElements}
}

// isPatternVariable checks if a symbol name represents a pattern variable
func isPatternVariable(name string) bool {
	// Pattern variables have the form: [varname][underscores][typename]
	// Examples: x_, x__, x___, _Integer, x_Integer, x__Integer, x___Integer
	// But NOT regular symbols with underscores in the middle: bool_test, my_function, etc.
	
	// Must start with underscore OR contain underscore followed by uppercase letter (type)
	if strings.HasPrefix(name, "_") {
		return true // _Integer, _, __, ___, etc.
	}
	
	// Look for pattern: letter(s) + underscore(s) + [optional type starting with uppercase]
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			// Found underscore - check if it's followed by more underscores or uppercase letter or end of string
			remaining := name[i:]
			if remaining == "_" || remaining == "__" || remaining == "___" {
				return true // x_, x__, x___
			}
			// Check if it's underscore(s) followed by type name (starts with uppercase)
			if len(remaining) > 1 {
				// Skip consecutive underscores
				j := 1
				for j < len(remaining) && remaining[j] == '_' {
					j++
				}
				if j < len(remaining) && remaining[j] >= 'A' && remaining[j] <= 'Z' {
					return true // x_Integer, x__Integer, x___Integer
				}
			}
			// If we found underscore but it doesn't match pattern variable format, it's not a pattern variable
			return false
		}
	}
	
	return false
}

// parsePatternVariable parses a pattern variable name and returns the variable name and type constraint
func parsePatternVariable(name string) (varName string, typeName string) {
	if !isPatternVariable(name) {
		return "", ""
	}
	
	// Split by underscore
	parts := strings.Split(name, "_")
	
	if len(parts) == 2 {
		if parts[0] == "" && parts[1] == "" {
			// Anonymous pattern: _ -> varName="", typeName=""
			return "", ""
		} else if parts[1] == "" {
			// Simple pattern: x_ -> varName="x", typeName=""
			return parts[0], ""
		} else if parts[0] != "" {
			// Named pattern: x_Integer -> varName="x", typeName="Integer"
			return parts[0], parts[1]
		}
	}
	
	// Invalid pattern
	return "", ""
}

// parsePatternInfo parses a pattern variable name and returns complete pattern information
func parsePatternInfo(name string) PatternInfo {
	if !isPatternVariable(name) {
		return PatternInfo{}
	}
	
	// Count consecutive underscores to determine pattern type
	underscoreCount := 0
	underscoreStart := -1
	
	// Find the first underscore
	for i, ch := range name {
		if ch == '_' {
			underscoreStart = i
			break
		}
	}
	
	if underscoreStart == -1 {
		return PatternInfo{} // No underscore found
	}
	
	// Count consecutive underscores
	for i := underscoreStart; i < len(name) && name[i] == '_'; i++ {
		underscoreCount++
	}
	
	// Determine pattern type based on underscore count
	var patternType PatternType
	switch underscoreCount {
	case 1:
		patternType = BlankPattern
	case 2:
		patternType = BlankSequencePattern
	case 3:
		patternType = BlankNullSequencePattern
	default:
		return PatternInfo{} // Invalid pattern
	}
	
	// Extract variable name (before underscores)
	varName := name[:underscoreStart]
	
	// Extract type name (after underscores)
	typeStart := underscoreStart + underscoreCount
	var typeName string
	if typeStart < len(name) {
		typeName = name[typeStart:]
	}
	
	return PatternInfo{
		Type:     patternType,
		VarName:  varName,
		TypeName: typeName,
	}
}

// matchesType checks if an expression matches a given type constraint
func matchesType(expr Expr, typeName string) bool {
	if typeName == "" {
		// No type constraint, matches anything
		return true
	}
	
	// Handle built-in types first
	switch typeName {
	case "Integer":
		if atom, ok := expr.(Atom); ok {
			return atom.AtomType == IntAtom
		}
	case "Real", "Float":
		if atom, ok := expr.(Atom); ok {
			return atom.AtomType == FloatAtom
		}
	case "Number", "Numeric":
		if atom, ok := expr.(Atom); ok {
			return atom.AtomType == IntAtom || atom.AtomType == FloatAtom
		}
	case "String":
		if atom, ok := expr.(Atom); ok {
			return atom.AtomType == StringAtom
		}
	case "Boolean", "Bool":
		if atom, ok := expr.(Atom); ok {
			return atom.AtomType == BoolAtom
		}
	case "Symbol":
		if atom, ok := expr.(Atom); ok {
			return atom.AtomType == SymbolAtom
		}
	case "Atom":
		_, ok := expr.(Atom)
		return ok
	case "List":
		_, ok := expr.(List)
		return ok
	default:
		// Handle arbitrary head symbols (e.g., Color, Point, etc.)
		if list, ok := expr.(List); ok && len(list.Elements) > 0 {
			if headAtom, ok := list.Elements[0].(Atom); ok && headAtom.AtomType == SymbolAtom {
				return headAtom.Value.(string) == typeName
			}
		}
		return false
	}
	
	return false
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

// applyFlat flattens nested applications of the same operator
func (e *Evaluator) applyFlat(headName string, list List) List {
	if len(list.Elements) <= 1 {
		return list
	}
	
	var newElements []Expr
	newElements = append(newElements, list.Elements[0]) // Keep the head
	
	for _, arg := range list.Elements[1:] {
		if argList, ok := arg.(List); ok && len(argList.Elements) > 0 {
			if argHead, ok := argList.Elements[0].(Atom); ok &&
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
	
	return List{Elements: newElements}
}

// applyOrderless sorts arguments for commutative operations
func (e *Evaluator) applyOrderless(list List) List {
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
	
	return List{Elements: newElements}
}

// applyOneIdentity handles f[x] -> x for functions with OneIdentity
func (e *Evaluator) applyOneIdentity(list List) List {
	if len(list.Elements) == 2 {
		// f[x] -> x when f has OneIdentity
		// Return the single argument directly (not wrapped in a list)
		return List{Elements: []Expr{list.Elements[1]}}
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
	}
	return nil, false
}

// Utility functions for numeric operations

// isNumeric checks if an expression is numeric
func isNumeric(expr Expr) bool {
	if atom, ok := expr.(Atom); ok {
		return atom.AtomType == IntAtom || atom.AtomType == FloatAtom
	}
	return false
}

// getNumericValue extracts numeric value from an expression
func getNumericValue(expr Expr) (float64, bool) {
	if atom, ok := expr.(Atom); ok {
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
	if atom, ok := expr.(Atom); ok {
		return atom.AtomType == BoolAtom
	}
	return false
}

// getBoolValue extracts boolean value from an expression
func getBoolValue(expr Expr) (bool, bool) {
	if atom, ok := expr.(Atom); ok && atom.AtomType == BoolAtom {
		return atom.Value.(bool), true
	}
	return false, false
}

// isSymbol checks if an expression is a symbol
func isSymbol(expr Expr) bool {
	if atom, ok := expr.(Atom); ok {
		return atom.AtomType == SymbolAtom
	}
	return false
}

// getSymbolName extracts symbol name from an expression
func getSymbolName(expr Expr) (string, bool) {
	if atom, ok := expr.(Atom); ok && atom.AtomType == SymbolAtom {
		return atom.Value.(string), true
	}
	return "", false
}

// patternsEqual checks if two patterns are structurally equivalent
// This ignores variable names and only compares pattern structure and types
func patternsEqual(pattern1, pattern2 Expr) bool {
	switch p1 := pattern1.(type) {
	case Atom:
		if p2, ok := pattern2.(Atom); ok {
			// For atoms, check if both are pattern variables or both are literals
			if p1.AtomType == SymbolAtom && p2.AtomType == SymbolAtom {
				name1 := p1.Value.(string)
				name2 := p2.Value.(string)
				
				isPattern1 := isPatternVariable(name1)
				isPattern2 := isPatternVariable(name2)
				
				if isPattern1 && isPattern2 {
					// Both are pattern variables - compare structure, not names
					info1 := parsePatternInfo(name1)
					info2 := parsePatternInfo(name2)
					return info1.Type == info2.Type && info1.TypeName == info2.TypeName
				} else if !isPattern1 && !isPattern2 {
					// Both are regular symbols - must match exactly
					return name1 == name2
				} else {
					// One is pattern, one is not - not equivalent
					return false
				}
			}
			// For non-symbol atoms, require exact match
			return p1.AtomType == p2.AtomType && p1.Value == p2.Value
		}
		return false
	case List:
		if p2, ok := pattern2.(List); ok {
			if len(p1.Elements) != len(p2.Elements) {
				return false
			}
			for i, elem1 := range p1.Elements {
				if !patternsEqual(elem1, p2.Elements[i]) {
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

// findFunctionWithPattern searches for a function definition with the given pattern
func findFunctionWithPattern(funcList List, targetPattern Expr) int {
	for i := 1; i < len(funcList.Elements); i++ {
		if funcDef, ok := funcList.Elements[i].(List); ok && len(funcDef.Elements) == 3 {
			if head, ok := funcDef.Elements[0].(Atom); ok && 
				head.AtomType == SymbolAtom && head.Value.(string) == "Function" {
				
				// Extract the pattern from this function definition
				if paramList, ok := funcDef.Elements[1].(List); ok {
					if patternsEqual(paramList, targetPattern) {
						return i
					}
				}
			}
		}
	}
	return -1
}

// replaceOrAddFunction replaces an existing function definition with the same pattern,
// or adds a new one if no matching pattern is found
func replaceOrAddFunction(existingList List, newFuncDef Expr) List {
	// Extract the pattern from the new function definition
	if funcDef, ok := newFuncDef.(List); ok && len(funcDef.Elements) == 3 {
		if head, ok := funcDef.Elements[0].(Atom); ok && 
			head.AtomType == SymbolAtom && head.Value.(string) == "Function" {
			
			targetPattern := funcDef.Elements[1]
			
			// Check if a function with this pattern already exists
			existingIndex := findFunctionWithPattern(existingList, targetPattern)
			
			if existingIndex != -1 {
				// Replace the existing function definition
				newElements := make([]Expr, len(existingList.Elements))
				copy(newElements, existingList.Elements)
				newElements[existingIndex] = newFuncDef
				return List{Elements: newElements}
			}
		}
	}
	
	// No existing pattern found, add the new function by specificity
	return insertFunctionBySpecificity(existingList, newFuncDef)
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

// registerDefaultBuiltins registers all built-in functions with their patterns
func registerDefaultBuiltins(registry *FunctionRegistry) {
	// Register built-in functions with pattern-based dispatch
	builtinPatterns := map[string]PatternFunc{
		// Arithmetic operations - support multiple arguments with sequence patterns
		"Plus(x___)":       wrapBuiltinFunc(EvaluatePlus),    // Zero or more arguments
		"Times(x___)":      wrapBuiltinFunc(EvaluateTimes),   // Zero or more arguments
		"Subtract(x_, y_)": wrapBuiltinFunc(EvaluateSubtract),
		"Divide(x_, y_)":   wrapBuiltinFunc(EvaluateDivide),
		"Power(x_, y_)":    wrapBuiltinFunc(EvaluatePower),
		
		// Comparison operations
		"Equal(x_, y_)":        wrapBuiltinFunc(EvaluateEqual),
		"Unequal(x_, y_)":      wrapBuiltinFunc(EvaluateUnequal),
		"Less(x_, y_)":         wrapBuiltinFunc(EvaluateLess),
		"Greater(x_, y_)":      wrapBuiltinFunc(EvaluateGreater),
		"LessEqual(x_, y_)":    wrapBuiltinFunc(EvaluateLessEqual),
		"GreaterEqual(x_, y_)": wrapBuiltinFunc(EvaluateGreaterEqual),
		"SameQ(x_, y_)":        wrapBuiltinFunc(EvaluateSameQ),
		"UnsameQ(x_, y_)":      wrapBuiltinFunc(EvaluateUnsameQ),
		
		// Logical operations (Not - And/Or are special forms)
		"Not(x_)": wrapBuiltinFunc(EvaluateNot),
		
		// List operations
		"Length(x_)":   wrapBuiltinFunc(EvaluateLength),
		"First(x_)":    wrapBuiltinFunc(EvaluateFirst),
		"Last(x_)":     wrapBuiltinFunc(EvaluateLast),
		"Rest(x_)":     wrapBuiltinFunc(EvaluateRest),
		"Most(x_)":     wrapBuiltinFunc(EvaluateMost),
		"Part(x_, i_)": wrapBuiltinFunc(EvaluatePart),
		
		// Type predicates
		"IntegerQ(x_)": wrapBuiltinFunc(EvaluateIntegerQ),
		"NumberQ(x_)":  wrapBuiltinFunc(EvaluateNumberQ),
		"StringQ(x_)":  wrapBuiltinFunc(EvaluateStringQ),
		"BooleanQ(x_)": wrapBuiltinFunc(EvaluateBooleanQ),
		"SymbolQ(x_)":  wrapBuiltinFunc(EvaluateSymbolQ),
		"ListQ(x_)":    wrapBuiltinFunc(EvaluateListQ),
		"AtomQ(x_)":    wrapBuiltinFunc(EvaluateAtomQ),
		"Head(x_)":     wrapBuiltinFuncNoErrorProp(EvaluateHead),
		
		// String functions
		"StringLength(x_)": wrapBuiltinFunc(EvaluateStringLength),
		"FullForm(x_)":     wrapBuiltinFunc(EvaluateFullForm),
	}
	
	// Register all patterns
	err := registry.RegisterPatternBuiltins(builtinPatterns)
	if err != nil {
		panic(fmt.Sprintf("Failed to register built-in patterns: %v", err))
	}
}

// wrapBuiltinFunc wraps an old-style BuiltinFunc to work with the new PatternFunc signature
func wrapBuiltinFunc(builtin BuiltinFunc) PatternFunc {
	return func(args []Expr, ctx *Context) Expr {
		// Check for errors in arguments and propagate them 
		// Note: Stack frame addition happens in the caller (evaluatePatternFunction)
		for _, arg := range args {
			if IsError(arg) {
				return arg
			}
		}
		
		return builtin(args)
	}
}

// wrapBuiltinFuncNoErrorProp wraps a BuiltinFunc that should NOT propagate errors
// (e.g., Head should analyze error expressions, not propagate them)
func wrapBuiltinFuncNoErrorProp(builtin BuiltinFunc) PatternFunc {
	return func(args []Expr, ctx *Context) Expr {
		// No error propagation - let the builtin handle errors as data
		return builtin(args)
	}
}