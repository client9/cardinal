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

// ParseString parses a string into an expression
func ParseString(input string) (core.Expr, error) {
	return core.ParseString(input)
}

// Parse is an alias for ParseString for convenience
func Parse(input string) (core.Expr, error) {
	return core.ParseString(input)
}

// EvaluateString is a convenience function that parses and evaluates a string
func EvaluateString(input string) (core.Expr, error) {
	expr, err := core.ParseString(input)
	if err != nil {
		return nil, err
	}
	e := NewEvaluator()
	return e.Evaluate(expr), nil
}
