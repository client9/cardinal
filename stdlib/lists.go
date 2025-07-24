package stdlib

import (
	"fmt"

	"github.com/client9/sexpr/core"
)

// List manipulation functions that work with core.List and core.Expr types

// LengthExpr returns the length of an expression
func LengthExpr(expr core.Expr) int64 {
	return expr.Length()
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
	if len(list.Elements) <= 1 {
		// Return error for empty lists
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Rest: expression %s has no elements", list.String()), []core.Expr{list})
	}

	// Use the SliceEnd primitive method to get elements from index 2 onwards
	return list.SliceEnd(2)
}

// MostExpr returns a new list with the last element removed
func MostExpr(list core.List) core.Expr {
	// For lists, return a new list with the last element removed
	if len(list.Elements) <= 1 {
		// Return error for empty lists
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Most: expression %s has no elements", list.String()), []core.Expr{list})
	}

	// Use the SliceBetween primitive method to get elements from 1 to length-1
	listLength := list.Length()
	if listLength == 1 {
		// Special case: if only one element, return just the head
		return core.List{Elements: []core.Expr{list.Elements[0]}}
	}
	return list.SliceBetween(1, listLength-1)
}

// PartList extracts an element from a list by integer index (1-based)
func PartList(list core.List, index int64) core.Expr {
	// Use the ElementAt primitive method
	return list.ElementAt(index)
}

// TakeList extracts the first or last n elements from a list
// Take(expr, n) takes first n elements; Take(expr, -n) takes last n elements
func TakeList(list core.List, n int64) core.Expr {
	if len(list.Elements) <= 1 {
		// Empty list - return empty list with same head
		return core.List{Elements: []core.Expr{list.Elements[0]}}
	}

	listLength := list.Length()

	if n == 0 {
		// Take 0 elements - return empty list with same head
		return core.List{Elements: []core.Expr{list.Elements[0]}}
	}

	if n > 0 {
		// Take first n elements using SliceStart primitive
		if n > listLength {
			n = listLength // Don't take more than available
		}
		return list.SliceStart(n)
	} else {
		// Take last |n| elements using SliceEnd primitive
		absN := -n
		if absN > listLength {
			absN = listLength // Don't take more than available
		}
		// Calculate starting index: total elements - |n| + 1
		startIdx := listLength - absN + 1
		return list.SliceEnd(startIdx)
	}
}

// TakeListSingle takes the nth element from a list and returns it as a single-element list
// Take(expr, [n]) - returns List(element_n)
func TakeListSingle(list core.List, indexList core.List) core.Expr {
	// Extract single integer from List(n_Integer)
	if len(indexList.Elements) != 2 { // Head + one element
		return core.NewErrorExpr("ArgumentError",
			"Take with list spec requires exactly one index", []core.Expr{indexList})
	}

	index, ok := core.ExtractInt64(indexList.Elements[1])
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"Take index must be an integer", []core.Expr{indexList.Elements[1]})
	}

	// Use the ElementAt primitive method
	element := list.ElementAt(index)
	if core.IsError(element) {
		return element
	}

	// Wrap the element in a list with the same head as the original list
	return core.List{Elements: []core.Expr{list.Elements[0], element}}
}

// TakeListRange takes a range of elements from a list
// Take(expr, [n, m]) - takes elements from index n to m (inclusive)
func TakeListRange(list core.List, indexList core.List) core.Expr {
	// Extract two integers from List(n_Integer, m_Integer)
	if len(indexList.Elements) != 3 { // Head + two elements
		return core.NewErrorExpr("ArgumentError",
			"Take with range spec requires exactly two indices", []core.Expr{indexList})
	}

	start, ok1 := core.ExtractInt64(indexList.Elements[1])
	end, ok2 := core.ExtractInt64(indexList.Elements[2])
	if !ok1 || !ok2 {
		return core.NewErrorExpr("ArgumentError",
			"Take indices must be integers", indexList.Elements[1:])
	}

	return takeListRange(list, start, end)
}

// takeListRange is a helper function that implements the range logic
func takeListRange(list core.List, start, end int64) core.Expr {
	if len(list.Elements) <= 1 {
		return core.List{Elements: []core.Expr{list.Elements[0]}} // Empty list with same head
	}

	listLength := list.Length()

	// Validate indices
	if start == 0 || end == 0 {
		return core.NewErrorExpr("PartError",
			"Take index 0 is out of bounds (indices start at 1)", []core.Expr{list})
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
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Take range [%d, %d] is out of bounds for list with %d elements",
				start, end, listLength), []core.Expr{list})
	}

	// Use the SliceBetween primitive method
	return list.SliceBetween(actualStart, actualEnd)
}

// DropList drops the first or last n elements from a list and returns the remainder
// Drop(expr, n) drops first n elements; Drop(expr, -n) drops last n elements
func DropList(list core.List, n int64) core.Expr {
	if len(list.Elements) <= 1 {
		// Empty list - return empty list with same head
		return core.List{Elements: []core.Expr{list.Elements[0]}}
	}

	listLength := list.Length()

	if n == 0 {
		// Drop 0 elements - return original list
		return list
	}

	if n > 0 {
		// Drop first n elements - use SliceEnd to keep the rest
		if n >= listLength {
			// Drop all elements - return empty list with same head
			return core.List{Elements: []core.Expr{list.Elements[0]}}
		}
		return list.SliceEnd(n + 1) // Start from n+1 to end
	} else {
		// Drop last |n| elements - use SliceStart to keep the beginning
		absN := -n
		if absN >= listLength {
			// Drop all elements - return empty list with same head
			return core.List{Elements: []core.Expr{list.Elements[0]}}
		}
		return list.SliceStart(listLength - absN) // Keep first (length - |n|) elements
	}
}

// DropListSingle drops the nth element from a list and returns the remainder
// Drop(expr, [n]) - removes the element at position n
func DropListSingle(list core.List, indexList core.List) core.Expr {
	// Extract single integer from List(n_Integer)
	if len(indexList.Elements) != 2 { // Head + one element
		return core.NewErrorExpr("ArgumentError",
			"Drop with list spec requires exactly one index", []core.Expr{indexList})
	}

	index, ok := core.ExtractInt64(indexList.Elements[1])
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"Drop index must be an integer", []core.Expr{indexList.Elements[1]})
	}

	return dropListSingle(list, index)
}

// DropListRange drops a range of elements from a list and returns the remainder
// Drop(expr, [n, m]) - removes elements from index n to m (inclusive)
func DropListRange(list core.List, indexList core.List) core.Expr {
	// Extract two integers from List(n_Integer, m_Integer)
	if len(indexList.Elements) != 3 { // Head + two elements
		return core.NewErrorExpr("ArgumentError",
			"Drop with range spec requires exactly two indices", []core.Expr{indexList})
	}

	start, ok1 := core.ExtractInt64(indexList.Elements[1])
	end, ok2 := core.ExtractInt64(indexList.Elements[2])
	if !ok1 || !ok2 {
		return core.NewErrorExpr("ArgumentError",
			"Drop indices must be integers", indexList.Elements[1:])
	}

	return dropListRange(list, start, end)
}

// dropListSingle is a helper function that drops a single element
func dropListSingle(list core.List, index int64) core.Expr {
	if len(list.Elements) <= 1 {
		return core.List{Elements: []core.Expr{list.Elements[0]}} // Empty list with same head
	}

	listLength := list.Length()

	// Validate index
	if index == 0 {
		return core.NewErrorExpr("PartError",
			"Drop index 0 is out of bounds (indices start at 1)", []core.Expr{list})
	}

	// Convert negative indices to positive
	actualIndex := index
	if index < 0 {
		actualIndex = listLength + index + 1
	}

	// Bounds checking
	if actualIndex < 1 || actualIndex > listLength {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Drop index %d is out of bounds for list with %d elements",
				index, listLength), []core.Expr{list})
	}

	// Use the SliceExclude primitive method to remove single element
	return list.SliceExclude(actualIndex, actualIndex)
}

// dropListRange is a helper function that drops a range of elements
func dropListRange(list core.List, start, end int64) core.Expr {
	if len(list.Elements) <= 1 {
		return core.List{Elements: []core.Expr{list.Elements[0]}} // Empty list with same head
	}

	listLength := list.Length()

	// Validate indices
	if start == 0 || end == 0 {
		return core.NewErrorExpr("PartError",
			"Drop index 0 is out of bounds (indices start at 1)", []core.Expr{list})
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
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Drop range [%d, %d] is out of bounds for list with %d elements",
				start, end, listLength), []core.Expr{list})
	}

	// Use the SliceExclude primitive method to remove the range
	return list.SliceExclude(actualStart, actualEnd)
}
