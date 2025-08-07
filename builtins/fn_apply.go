// Package builtins contains engine-dependent built-in functions that require evaluator access
package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// ApplyExpr applies a function to a list of arguments using EngineFunc signature
// Apply(f, {a, b, c}) -> f(a, b, c)
func ApplyExpr(e *engine.Evaluator, c *engine.Context, function core.Expr, list core.Expr) core.Expr {
	// Check if the second argument is a list
	listExpr, ok := list.(core.List)
	if !ok {
		return core.NewError("ArgumentError", "second argument must be a list")
	}

	// Extract elements (skip the head)
	elements := listExpr.Tail()

	// Create function application: function(element1, element2, ...)
	applicationElements := make([]core.Expr, len(elements)+1)
	applicationElements[0] = function
	copy(applicationElements[1:], elements)
	application := core.NewListFromExprs(applicationElements...)

	// Evaluate the function application using the evaluator
	return e.Evaluate(application)
}
