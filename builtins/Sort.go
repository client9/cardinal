package builtins

import (
	"sort"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Sort

// Sort sorts the elements of a list using canonical ordering
// Uses the same ordering as the Orderless attribute and mathematical functions
//
// @ExprPattern (_)
func Sort(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	list, ok := expr.(core.List)
	if !ok || list.Length() < 2 {
		// Not a list or too few elements to sort
		return expr
	}

	head := list.Head()
	elements := make([]core.Expr, list.Length())
	copy(elements, list.Tail())

	// Sort arguments using canonical ordering
	sort.Slice(elements, func(i, j int) bool {
		return core.CanonicalCompare(elements[i], elements[j])
	})

	// Reconstruct the list with sorted arguments
	resultElements := make([]core.Expr, list.Length()+1)
	resultElements[0] = head
	copy(resultElements[1:], elements)

	return core.NewListFromExprs(resultElements...)
}
