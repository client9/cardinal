package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Flatten

// FlattenExpr flattens nested lists into a single one-dimensional list
// Flatten(List(1, 2, List(3, 4))) -> List(1, 2, 3, 4)
// For now, flattens all levels (no level specification)
//
// @ExprPattern (_List)
func FlattenExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	list, _ := args[0].(core.List)

	return flatten(list)

}

func flatten(list core.List) core.Expr {

	if list.Length() == 0 {
		// Empty list, return unchanged
		return list
	}

	// Extract the head (List, Zoo, etc.)
	head := list.HeadExpr()

	// Get the head name to determine if we should flatten sublists with the same head
	headName := head.String()

	// Flatten the elements recursively
	var flattenedElements []core.Expr
	flattenedElements = append(flattenedElements, head) // Keep the original head

	for _, element := range list.Tail() {

		// Check if this element is a list with the same head
		if elementList, ok := element.(core.List); ok {
			if elementHeadName, ok := core.ExtractSymbol(elementList.HeadExpr()); ok && elementHeadName == headName {
				// Same head - flatten this sublist's elements
				// First recursively flatten the sublist
				flattened := flatten(elementList)
				if flattenedList, ok := flattened.(core.List); ok && flattenedList.Length() > 0 {
					// Add all elements except the head
					flattenedElements = append(flattenedElements, flattenedList.Tail()...)
				}
			} else {
				// Different head - recursively flatten but keep as single element
				flattenedElements = append(flattenedElements, flatten(elementList))
			}
		} else {
			// Not a list - add as-is
			flattenedElements = append(flattenedElements, element)
		}
	}

	return core.NewListFromExprs(flattenedElements...)
}
