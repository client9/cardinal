package sexpr

import (
	"fmt"
)

// EvaluationStack represents the current evaluation call stack
type EvaluationStack struct {
	frames   []StackFrame
	depth    int
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
	functionRegistry *FunctionRegistry // Unified pattern-based function system
	stack            *EvaluationStack
}

// NewContext creates a new evaluation context
func NewContext() *Context {
	ctx := &Context{
		variables:        make(map[string]Expr),
		parent:           nil,
		symbolTable:      NewSymbolTable(),
		functionRegistry: NewFunctionRegistry(),
		stack:            NewEvaluationStack(1000), // Default max depth of 1000
	}

	// Set up built-in attributes
	setupBuiltinAttributes(ctx.symbolTable)

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

// SetStack sets the evaluation stack for the context
func (c *Context) SetStack(stack *EvaluationStack) {
	c.stack = stack
}

// GetFunctionRegistry returns the context's function registry
func (c *Context) GetFunctionRegistry() *FunctionRegistry {
	return c.functionRegistry
}

// GetSymbolTable returns the context's symbol table
func (c *Context) GetSymbolTable() *SymbolTable {
	return c.symbolTable
}
