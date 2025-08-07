package stdlib

import (
	"fmt"
	"sort"

	"github.com/client9/sexpr/core"
)

// List manipulation functions that work with core.List and core.Expr types

// LengthExpr returns the length of an expression
func LengthExpr(expr core.Expr) int64 {
	return expr.Length()
}

// ListAppends adds an expression to the end of List type
func ListAppend(list core.List, e core.Expr) core.Expr {
	return list.Append(e)
}

// FirstExpr returns the first element of a list (after the head)
func FirstExpr(list core.List) core.Expr {
	// Use the ElementAt primitive method
	return list.ElementAt(1)
}

// LastExpr returns the last element of a list
func LastExpr(list core.List) core.Expr {
	// Use the ElementAt primitive method with negative indexing
	return list.ElementAt(-1)
}

// RestExpr returns a new list with the first element after head removed
func RestExpr(list core.List) core.Expr {
	// For lists, return a new list with the first element after head removed
	if list.Length() == 0 {
		// Return error for empty lists
		return core.NewError("PartError",
			fmt.Sprintf("Rest: expression %s has no elements", list.String()))
	}

	// Use the modern Slice method to get elements from index 2 onwards
	listLength := list.Length()
	if listLength == 1 {
		// Only head, return empty list with same head
		return core.NewList(list.Head())
	}
	return list.Slice(2, listLength)
}

// MostExpr returns a new list with the last element removed
func MostExpr(list core.List) core.Expr {
	listLength := list.Length()
	// For lists, return a new list with the last element removed
	if listLength == 0 {
		// Return error for empty lists
		return core.NewError("PartError",
			fmt.Sprintf("Most: expression %s has no elements", list.String()))
	}

	// Use the modern Slice method to get elements from 1 to length-1
	if listLength == 1 {
		// Special case: if only one element, return just the head
		return core.NewList(list.Head())
	}
	return list.Slice(1, listLength-1)
}

// PartList extracts an element from a list by integer index (1-based)
func PartList(list core.List, index int64) core.Expr {
	// Use the ElementAt primitive method
	return list.ElementAt(index)
}

// TakeList extracts the first or last n elements from a list
// Take(expr, n) takes first n elements; Take(expr, -n) takes last n elements
func TakeList(list core.List, n int64) core.Expr {
	listLength := list.Length()
	if listLength == 0 || n == 0 {
		return core.NewList(list.Head())
	}

	if n > 0 {
		// Take first n elements using modern Slice method
		if n > listLength {
			n = listLength // Don't take more than available
		}
		return list.Slice(1, n)
	} else {
		// Take last |n| elements using modern Slice method
		absN := -n
		if absN > listLength {
			absN = listLength // Don't take more than available
		}
		// Calculate starting index: total elements - |n| + 1
		startIdx := listLength - absN + 1
		return list.Slice(startIdx, listLength)
	}
}

// TakeListSingle takes the nth element from a list and returns it as a single-element list
// Take(expr, [n]) - returns List(element_n)
func TakeListSingle(list core.List, indexList core.List) core.Expr {
	// Extract single integer from List(n_Integer)
	if indexList.Length() != 1 { // Head + one element
		return core.NewError("ArgumentError",
			"Take with list spec requires exactly one index")
	}
	args := indexList.Tail()
	index, ok := core.ExtractInt64(args[0])
	if !ok {
		return core.NewError("ArgumentError",
			"Take index must be an integer")
	}

	// Use the ElementAt primitive method
	element := list.ElementAt(index)
	if core.IsError(element) {
		return element
	}

	// Wrap the element in a list with the same head as the original list
	return core.NewList(list.Head(), element)
}

// TakeListRange takes a range of elements from a list
// Take(expr, [n, m]) - takes elements from index n to m (inclusive)
func TakeListRange(list core.List, indexList core.List) core.Expr {
	// Extract two integers from List(n_Integer, m_Integer)
	if indexList.Length() != 2 { // Head + two elements
		return core.NewError("ArgumentError",
			"Take with range spec requires exactly two indices")
	}
	args := indexList.Tail()
	start, ok1 := core.ExtractInt64(args[0])
	end, ok2 := core.ExtractInt64(args[1])
	if !ok1 || !ok2 {
		return core.NewError("ArgumentError",
			"Take indices must be integers")
	}

	return takeListRange(list, start, end)
}

// takeListRange is a helper function that implements the range logic
func takeListRange(list core.List, start, end int64) core.Expr {
	listLength := list.Length()
	if listLength == 0 {
		return core.NewList(list.Head())
	}

	// Validate indices
	if start == 0 || end == 0 {
		return core.NewError("PartError",
			"Take index 0 is out of bounds (indices start at 1)")
	}

	// Convert negative indices to positive
	actualStart := start
	actualEnd := end
	if start < 0 {
		actualStart = listLength + start + 1
	}
	if end < 0 {
		actualEnd = listLength + end + 1
	}

	// Bounds checking
	if actualStart < 1 || actualEnd > listLength || actualStart > actualEnd {
		return core.NewError("PartError",
			fmt.Sprintf("Take range [%d, %d] is out of bounds for list with %d elements",
				start, end, listLength))
	}

	// Use the modern Slice method
	return list.Slice(actualStart, actualEnd)
}

// DropList drops the first or last n elements from a list and returns the remainder
// Drop(expr, n) drops first n elements; Drop(expr, -n) drops last n elements
func DropList(list core.List, n int64) core.Expr {
	listLength := list.Length()
	if listLength == 0 {
		return core.NewList(list.Head())
	}

	if n == 0 {
		// Drop 0 elements - return original list
		return list
	}

	if n > 0 {
		// Drop first n elements using modern Slice method
		if n >= listLength {
			// Drop all elements - return empty list with same head
			return core.NewList(list.Head())
		}
		return list.Slice(n+1, listLength) // Start from n+1 to end
	} else {
		// Drop last |n| elements using modern Slice method
		absN := -n
		if absN >= listLength {
			return core.NewList(list.Head())
		}
		return list.Slice(1, listLength-absN) // Keep first (length - |n|) elements
	}
}

// DropListSingle drops the nth element from a list and returns the remainder
// Drop(expr, [n]) - removes the element at position n
func DropListSingle(list core.List, indexList core.List) core.Expr {
	// Extract single integer from List(n_Integer)
	if indexList.Length() != 1 { // Head + one element
		return core.NewError("ArgumentError",
			"Drop with list spec requires exactly one index")
	}
	args := indexList.Tail()
	index, ok := core.ExtractInt64(args[0])
	if !ok {
		return core.NewError("ArgumentError",
			"Drop index must be an integer")
	}

	return dropListSingle(list, index)
}

// DropListRange drops a range of elements from a list and returns the remainder
// Drop(expr, [n, m]) - removes elements from index n to m (inclusive)
func DropListRange(list core.List, indexList core.List) core.Expr {
	// Extract two integers from List(n_Integer, m_Integer)
	if indexList.Length() != 2 {
		return core.NewError("ArgumentError",
			"Drop with range spec requires exactly two indices")
	}

	args := indexList.Tail()
	start, ok1 := core.ExtractInt64(args[0])
	end, ok2 := core.ExtractInt64(args[1])

	if !ok1 || !ok2 {
		return core.NewError("ArgumentError",
			"Drop indices must be integers")
	}

	return dropListRange(list, start, end)
}

// dropListSingle is a helper function that drops a single element
func dropListSingle(list core.List, index int64) core.Expr {
	listLength := list.Length()

	if listLength == 0 {
		return core.NewList(list.Head())
	}

	// Validate index
	if index == 0 {
		return core.NewError("PartError",
			"Drop index 0 is out of bounds (indices start at 1)")
	}

	// Convert negative indices to positive
	actualIndex := index
	if index < 0 {
		actualIndex = listLength + index + 1
	}

	// Bounds checking
	if actualIndex < 1 || actualIndex > listLength {
		return core.NewError("PartError", "Index out of bounds")
	}

	// Use modern Slice and Join methods to exclude single element
	if actualIndex == 1 {
		// Dropping first element
		if listLength == 1 {
			return core.NewList(list.Head())
		}
		return list.Slice(2, listLength)
	} else if actualIndex == listLength {
		// Dropping last element
		return list.Slice(1, listLength-1)
	} else {
		// Dropping middle element - join before and after
		before := list.Slice(1, actualIndex-1)
		after := list.Slice(actualIndex+1, listLength)
		if sliceable, ok := before.(core.Sliceable); ok {
			return sliceable.Join(after.(core.Sliceable))
		}
		return core.NewError("InternalError", "Failed to join slices")
	}
}

// dropListRange is a helper function that drops a range of elements
func dropListRange(list core.List, start, end int64) core.Expr {
	listLength := list.Length()
	if listLength == 0 {
		return core.NewList(list.Head())
	}

	// Validate indices
	if start == 0 || end == 0 {
		return core.NewError("PartError",
			"Drop index 0 is out of bounds (indices start at 1)")
	}

	// Convert negative indices to positive
	actualStart := start
	actualEnd := end
	if start < 0 {
		actualStart = listLength + start + 1
	}
	if end < 0 {
		actualEnd = listLength + end + 1
	}

	// Bounds checking
	if actualStart < 1 || actualEnd > listLength || actualStart > actualEnd {
		return core.NewError("PartError",
			fmt.Sprintf("Drop range [%d, %d] is out of bounds for list with %d elements",
				start, end, listLength))
	}

	// Use modern Slice and Join methods to exclude the range
	if actualStart == 1 && actualEnd == listLength {
		// Dropping everything
		return core.NewList(list.Head())
	} else if actualStart == 1 {
		// Dropping from beginning
		return list.Slice(actualEnd+1, listLength)
	} else if actualEnd == listLength {
		// Dropping to end
		return list.Slice(1, actualStart-1)
	} else {
		// Dropping middle range - join before and after
		before := list.Slice(1, actualStart-1)
		after := list.Slice(actualEnd+1, listLength)
		if sliceable, ok := before.(core.Sliceable); ok {
			return sliceable.Join(after.(core.Sliceable))
		}
		return core.NewError("InternalError", "Failed to join slices")
	}
}

// FlattenExpr flattens nested lists into a single one-dimensional list
// Flatten(List(1, 2, List(3, 4))) -> List(1, 2, 3, 4)
// For now, flattens all levels (no level specification)
func FlattenExpr(expr core.Expr) core.Expr {
	list, ok := expr.(core.List)
	if !ok {
		// If it's not a list, return it unchanged
		return expr
	}

	if list.Length() == 0 {
		// Empty list, return unchanged
		return expr
	}

	// Extract the head (List, Zoo, etc.)
	head := list.HeadExpr()

	// Get the head name to determine if we should flatten sublists with the same head
	headName, isSymbol := core.ExtractSymbol(head)
	if !isSymbol {
		// If head is not a symbol, return unchanged
		return expr
	}

	// Flatten the elements recursively
	var flattenedElements []core.Expr
	flattenedElements = append(flattenedElements, head) // Keep the original head

	for _, element := range list.Tail() {

		// Check if this element is a list with the same head
		if elementList, ok := element.(core.List); ok {
			if elementHeadName, ok := core.ExtractSymbol(elementList.HeadExpr()); ok && elementHeadName == headName {
				// Same head - flatten this sublist's elements
				// First recursively flatten the sublist
				flattened := FlattenExpr(element)
				if flattenedList, ok := flattened.(core.List); ok && flattenedList.Length() > 0 {
					// Add all elements except the head
					flattenedElements = append(flattenedElements, flattenedList.Tail()...)
				}
			} else {
				// Different head - recursively flatten but keep as single element
				flattenedElements = append(flattenedElements, FlattenExpr(element))
			}
		} else {
			// Not a list - add as-is
			flattenedElements = append(flattenedElements, element)
		}
	}

	return core.NewListFromExprs(flattenedElements...)
}

// Sort sorts the elements of a list using canonical ordering
// Uses the same ordering as the Orderless attribute and mathematical functions
func Sort(expr core.Expr) core.Expr {
	list, ok := expr.(core.List)
	if !ok || list.Length() < 2 {
		// Not a list or too few elements to sort
		return expr
	}

	head := list.HeadExpr()
	args := make([]core.Expr, list.Length())
	copy(args, list.Tail())

	// Sort arguments using canonical ordering
	sort.Slice(args, func(i, j int) bool {
		return core.CanonicalCompare(args[i], args[j])
	})

	// Reconstruct the list with sorted arguments
	resultElements := make([]core.Expr, list.Length()+1)
	resultElements[0] = head
	copy(resultElements[1:], args)

	return core.NewListFromExprs(resultElements...)
}
