// Package builtins contains engine-dependent built-in functions that require evaluator access
package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Map
// @ExprAttributes
//

// MapExpr applies a function to each element of a list using EngineFunc signature
// Map(f, {a, b, c}) -> {f(a), f(b), f(c)}
//
// @ExprPattern (_,_List)
func MapExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	fn := args[0]
	listExpr := args[1].(core.List)

	// If the list is empty or only has a head, return it unchanged
	if listExpr.Length() == 0 {
		return listExpr
	}

	// Extract head and elements
	head := listExpr.Head()
	elements := listExpr.Tail()

	// Apply the function to each element
	resultElements := make([]core.Expr, len(elements)+1)
	resultElements[0] = head // Keep the same head

	for i, element := range elements {
		// Create function application: function(element)
		applicationElements := []core.Expr{fn, element}
		application := core.NewListFromExprs(applicationElements...)

		// Evaluate the function application using the evaluator
		result := e.Evaluate(application)
		resultElements[i+1] = result
	}

	return core.NewListFromExprs(resultElements...)
}
