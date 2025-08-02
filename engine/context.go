package engine

import (
	"fmt"
	"github.com/client9/sexpr/core"
)

// EvaluationStack represents the current evaluation call stack
type EvaluationStack struct {
	frames   []core.StackFrame
	depth    int
	maxDepth int
}

// NewEvaluationStack creates a new evaluation stack with the given maximum depth
func NewEvaluationStack(maxDepth int) *EvaluationStack {
	return &EvaluationStack{
		frames:   make([]core.StackFrame, 0, maxDepth),
		depth:    0,
		maxDepth: maxDepth,
	}
}

// Push adds a new frame to the stack and checks for recursion limits
func (s *EvaluationStack) Push(function string, expression core.Expr) error {
	if s.depth >= s.maxDepth {
		return fmt.Errorf("maximum recursion depth exceeded: %d", s.maxDepth)
	}

	frame := core.StackFrame{
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
func (s *EvaluationStack) GetFrames() []core.StackFrame {
	frames := make([]core.StackFrame, len(s.frames))
	copy(frames, s.frames)
	return frames
}

// Depth returns the current stack depth
func (s *EvaluationStack) Depth() int {
	return s.depth
}

// Context represents the evaluation context with variable bindings and symbol attributes
type Context struct {
	variables        map[string]core.Expr
	parent           *Context
	symbolTable      *SymbolTable
	functionRegistry *FunctionRegistry // Unified pattern-based function system
	stack            *EvaluationStack
	scopedVars       map[string]bool // Variables that are locally scoped (for Block)
}

// NewContext creates a new evaluation context
func NewContext() *Context {
	ctx := &Context{
		variables:        make(map[string]core.Expr),
		parent:           nil,
		symbolTable:      NewSymbolTable(),
		functionRegistry: NewFunctionRegistry(),
		stack:            NewEvaluationStack(1000), // Default max depth of 1000
		scopedVars:       make(map[string]bool),
	}

	// Note: Builtin attributes and functions are now registered by the top-level API
	// This allows breaking the circular import between engine and wrapped packages

	return ctx
}

// NewChildContext creates a child context with a parent
func NewChildContext(parent *Context) *Context {
	return &Context{
		variables:        make(map[string]core.Expr),
		parent:           parent,
		symbolTable:      parent.symbolTable,      // Share symbol table with parent
		functionRegistry: parent.functionRegistry, // Share function registry with parent
		stack:            parent.stack,            // Share evaluation stack with parent
		scopedVars:       make(map[string]bool),
	}
}

// GetFunctionDefinitions returns a list of patterns registered to the given symbol
// exposed for debugging.
func (c *Context) GetFunctionDefinitions(name string) []FunctionDef {
	return c.functionRegistry.GetFunctionDefinitions(name)
}

// Set sets a variable in the context
// If this is a child context and the variable is not in scopedVars, set it in the parent
// Returns an error if the symbol is Protected
func (c *Context) Set(name string, value core.Expr) error {
	// Check if symbol is protected
	if c.symbolTable.HasAttribute(name, Protected) {
		return fmt.Errorf("symbol %s is Protected", name)
	}

	// If this variable is explicitly scoped to this context, set it here
	if c.scopedVars[name] {
		c.variables[name] = value
		return nil
	}

	// If this is a child context and variable is not scoped here, set in parent
	if c.parent != nil {
		return c.parent.Set(name, value)
	}

	// Otherwise set in current context (root context or explicitly local)
	c.variables[name] = value
	return nil
}

// Get retrieves a variable from the context (searches up the parent chain)
func (c *Context) Get(name string) (core.Expr, bool) {
	if value, ok := c.variables[name]; ok {
		return value, true
	}
	if c.parent != nil {
		return c.parent.Get(name)
	}
	return nil, false
}

// Delete removes a variable from the context
func (c *Context) Delete(name string) error {
	if c.symbolTable.HasAttribute(name, Protected) {
		return fmt.Errorf("symbol %s is Protected", name)
	}
	delete(c.variables, name)
	if c.parent != nil {
		return c.parent.Delete(name)
	}
	return nil
}

// AddScopedVar marks a variable as locally scoped to this context
func (c *Context) AddScopedVar(name string) {
	c.scopedVars[name] = true
}

// NewBlockContext creates a child context for Block evaluation with specified scoped variables
func NewBlockContext(parent *Context, scopedVarNames []string) *Context {
	ctx := NewChildContext(parent)
	for _, varName := range scopedVarNames {
		ctx.AddScopedVar(varName)
	}
	return ctx
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
