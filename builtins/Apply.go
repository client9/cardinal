// Package builtins contains engine-dependent built-in functions that require evaluator access
package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Apply
// TODO: is arg0 a Symbol

// ApplyExpr applies a function to a list of arguments using EngineFunc signature
// Apply(f, {a, b, c}) -> f(a, b, c)
//
// @ExprPattern (_,_List)
func ApplyExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	//function core.Expr, list core.Expr) core.Expr {

	fn := args[0]
	elements := args[1].(core.List).Tail()

	// Create function application: function(element1, element2, ...)
	applicationElements := make([]core.Expr, len(elements)+1)
	applicationElements[0] = fn
	copy(applicationElements[1:], elements)
	application := core.NewListFromExprs(applicationElements...)

	// Evaluate the function application using the evaluator
	return e.Evaluate(application)
}
