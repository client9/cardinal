// Package builtins contains engine-dependent built-in functions that require evaluator access
package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// MapExpr applies a function to each element of a list using EngineFunc signature
// Map(f, {a, b, c}) -> {f(a), f(b), f(c)}
func MapExpr(e *engine.Evaluator, c *engine.Context, function core.Expr, list core.Expr) core.Expr {

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
		result := e.Evaluate(application)
		resultElements[i+1] = result
	}

	return core.List{Elements: resultElements}
}
