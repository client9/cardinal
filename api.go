package sexpr

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// Public API - these are the only functions external code should use

// NewEvaluator creates a new evaluator with a fresh context and registers all builtins
func NewEvaluator() *engine.Evaluator {
	evaluator := engine.NewEvaluator()
	// Set up attributes and register functions using generated code
	SetupBuiltinAttributes(evaluator.GetContext().GetSymbolTable())
	RegisterDefaultBuiltins(evaluator.GetContext().GetFunctionRegistry())
	return evaluator
}

// NewEvaluatorWithContext creates an evaluator with a specific context
func NewEvaluatorWithContext(ctx *engine.Context) *engine.Evaluator {
	return engine.NewEvaluatorWithContext(ctx)
}

// ParseString parses a string into an expression
func ParseString(input string) (core.Expr, error) {
	return engine.ParseString(input)
}

// Parse is an alias for ParseString for convenience
func Parse(input string) (core.Expr, error) {
	return engine.ParseString(input)
}

// EvaluateString is a convenience function that parses and evaluates a string
func EvaluateString(input string) (core.Expr, error) {
	expr, err := ParseString(input)
	if err != nil {
		return nil, err
	}
	e := NewEvaluator()
	c := e.GetContext()
	return e.Evaluate(c, expr), nil
}

// NewContext creates a new evaluation context
func NewContext() *engine.Context {
	return engine.NewContext()
}

// SetupBuiltinAttributes and RegisterDefaultBuiltins are now provided by generated builtin_setup.go

// AttributesToString converts attributes to a string representation
func AttributesToString(attrs []engine.Attribute) string {
	return engine.AttributesToString(attrs)
}
