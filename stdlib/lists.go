package stdlib

import (
	"fmt"

	"github.com/client9/sexpr/core"
)

// List manipulation functions that work with core.List and core.Expr types

// LengthExpr returns the length of an expression
func LengthExpr(expr core.Expr) int64 {
	switch ex := expr.(type) {
	case core.List:
		// For lists, return the number of elements (excluding the head)
		if len(ex.Elements) == 0 {
			return 0 // Empty list has length 0
		}
		return int64(len(ex.Elements) - 1) // Subtract 1 for the head
	case core.ObjectExpr:
		// Handle Association - need to access AssociationValue but it's in main package
		// For now, just return 0 for ObjectExpr types
		return 0
	default:
		// For atoms and other expressions, length is 0
		return 0
	}
}

// FirstExpr returns the first element of a list (after the head)
func FirstExpr(list core.List) core.Expr {
	// For lists, return the first element after the head
	if len(list.Elements) <= 1 {
		// Return error for empty lists
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("First: expression %s has no elements", list.String()), []core.Expr{list})
	}
	return list.Elements[1] // Index 1 is first element after head (index 0)
}

// LastExpr returns the last element of a list
func LastExpr(list core.List) core.Expr {
	// For lists, return the last element
	if len(list.Elements) <= 1 {
		// Return error for empty lists
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Last: expression %s has no elements", list.String()), []core.Expr{list})
	}
	return list.Elements[len(list.Elements)-1] // Last element
}

// RestExpr returns a new list with the first element after head removed
func RestExpr(list core.List) core.Expr {
	// For lists, return a new list with the first element after head removed
	if len(list.Elements) <= 1 {
		// Return error for empty lists
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Rest: expression %s has no elements", list.String()), []core.Expr{list})
	}

	// Create new list: head + elements[2:] (skip first element after head)
	if len(list.Elements) == 2 {
		// Special case: if only head and one element, return just the head
		return core.List{Elements: []core.Expr{list.Elements[0]}}
	}

	newElements := make([]core.Expr, len(list.Elements)-1)
	newElements[0] = list.Elements[0]        // Keep the head
	copy(newElements[1:], list.Elements[2:]) // Copy everything after the first element
	return core.List{Elements: newElements}
}

// MostExpr returns a new list with the last element removed
func MostExpr(list core.List) core.Expr {
	// For lists, return a new list with the last element removed
	if len(list.Elements) <= 1 {
		// Return error for empty lists
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Most: expression %s has no elements", list.String()), []core.Expr{list})
	}

	// Create new list with all elements except the last one
	if len(list.Elements) == 2 {
		// Special case: if only head and one element, return just the head
		return core.List{Elements: []core.Expr{list.Elements[0]}}
	}

	newElements := make([]core.Expr, len(list.Elements)-1)
	copy(newElements, list.Elements[:len(list.Elements)-1])
	return core.List{Elements: newElements}
}

// PartList extracts an element from a list by integer index (1-based)
func PartList(list core.List, index int64) core.Expr {
	// For lists, return the element at the specified index (1-based)
	if len(list.Elements) <= 1 {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Part: expression %s has no elements", list.String()), []core.Expr{list})
	}

	// Handle negative indexing: -1 is last element, -2 is second to last, etc.
	var actualIndex int
	if index < 0 {
		// Negative indexing: -1 = last, -2 = second to last, etc.
		actualIndex = len(list.Elements) + int(index)
	} else if index > 0 {
		// Positive 1-based indexing: convert to 0-based for internal use
		actualIndex = int(index)
	} else {
		// index == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Part index %d is out of bounds (indices start at 1)", index), []core.Expr{list})
	}

	// Check bounds (remember: list.Elements[0] is the head, actual elements start at index 1)
	if actualIndex < 1 || actualIndex >= len(list.Elements) {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Part index %d is out of bounds for expression with %d elements",
				index, len(list.Elements)-1), []core.Expr{list})
	}

	return list.Elements[actualIndex]
}