// Package builtins contains engine-dependent built-in functions that require evaluator access
package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// Evaluator interface for functions that need recursive evaluation and context access
type Evaluator interface {
	Evaluate(expr core.Expr) core.Expr
	GetContext() *engine.Context
}

// MapExpr applies a function to each element of a list using EngineFunc signature
// Map(f, {a, b, c}) -> {f(a), f(b), f(c)}
func MapExpr(evaluator Evaluator, function core.Expr, list core.Expr) core.Expr {

	// Check if the second argument is a list
	listExpr, ok := list.(core.List)
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"Map expects a list as the second argument", []core.Expr{list})
	}

	// If the list is empty or only has a head, return it unchanged
	if len(listExpr.Elements) <= 1 {
		return listExpr
	}

	// Extract head and elements
	head := listExpr.Elements[0]
	elements := listExpr.Elements[1:]

	// Apply the function to each element
	resultElements := make([]core.Expr, len(elements)+1)
	resultElements[0] = head // Keep the same head

	for i, element := range elements {
		// Create function application: function(element)
		applicationElements := []core.Expr{function, element}
		application := core.List{Elements: applicationElements}

		// Evaluate the function application using the evaluator
		result := evaluator.Evaluate(application)
		resultElements[i+1] = result
	}

	return core.List{Elements: resultElements}
}

// ApplyExpr applies a function to a list of arguments using EngineFunc signature
// Apply(f, {a, b, c}) -> f(a, b, c)
func ApplyExpr(evaluator Evaluator, function core.Expr, list core.Expr) core.Expr {
	// Check if the second argument is a list
	listExpr, ok := list.(core.List)
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"Apply expects a list as the second argument", []core.Expr{list})
	}

	// If the list is empty, just return the function
	if len(listExpr.Elements) <= 1 {
		return function
	}

	// Extract elements (skip the head)
	elements := listExpr.Elements[1:]

	// Create function application: function(element1, element2, ...)
	applicationElements := make([]core.Expr, len(elements)+1)
	applicationElements[0] = function
	copy(applicationElements[1:], elements)
	application := core.List{Elements: applicationElements}

	// Evaluate the function application using the evaluator
	return evaluator.Evaluate(application)
}
