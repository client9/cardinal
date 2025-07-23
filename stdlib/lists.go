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

// TakeList extracts the first or last n elements from a list
// Take(expr, n) takes first n elements; Take(expr, -n) takes last n elements
func TakeList(list core.List, n int64) core.Expr {
	if len(list.Elements) <= 1 {
		// Empty list - return empty list with same head
		return core.List{Elements: []core.Expr{list.Elements[0]}}
	}

	listLength := len(list.Elements) - 1 // Subtract head element
	
	if n == 0 {
		// Take 0 elements - return empty list with same head
		return core.List{Elements: []core.Expr{list.Elements[0]}}
	}

	var startIdx, endIdx int
	
	if n > 0 {
		// Take first n elements
		if n > int64(listLength) {
			n = int64(listLength) // Don't take more than available
		}
		startIdx = 1                // Start after head
		endIdx = 1 + int(n)        // End at head + n
	} else {
		// Take last |n| elements  
		absN := -n
		if absN > int64(listLength) {
			absN = int64(listLength) // Don't take more than available
		}
		startIdx = len(list.Elements) - int(absN) // Start at end - |n|
		endIdx = len(list.Elements)               // End at the very end
	}

	// Create new list with head + selected elements
	newElements := make([]core.Expr, 1+(endIdx-startIdx))
	newElements[0] = list.Elements[0] // Keep the head
	copy(newElements[1:], list.Elements[startIdx:endIdx])
	
	return core.List{Elements: newElements}
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
	
	// Get the element using PartList logic
	element := PartList(list, index)
	
	// If it's an error, return the error
	if element.String() == "Error" || (len(element.String()) > 7 && element.String()[:7] == "$Failed") {
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
	
	listLength := int64(len(list.Elements) - 1) // Subtract head element
	
	// Convert to 0-based indexing and handle negatives
	var startIdx, endIdx int64
	
	if start > 0 {
		startIdx = start - 1 // Convert from 1-based
	} else if start < 0 {
		startIdx = listLength + start // Negative indexing
	} else {
		// start == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError", 
			"Take index 0 is out of bounds (indices start at 1)", []core.Expr{list})
	}
	
	if end > 0 {
		endIdx = end - 1 // Convert from 1-based  
	} else if end < 0 {
		endIdx = listLength + end // Negative indexing
	} else {
		// end == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError", 
			"Take index 0 is out of bounds (indices start at 1)", []core.Expr{list})
	}
	
	// Bounds checking
	if startIdx < 0 || endIdx >= listLength || startIdx > endIdx {
		return core.NewErrorExpr("PartError", 
			fmt.Sprintf("Take range [%d, %d] is out of bounds for list with %d elements", 
				start, end, listLength), []core.Expr{list})
	}
	
	// Create result slice
	resultSize := endIdx - startIdx + 1
	newElements := make([]core.Expr, resultSize + 1) // +1 for head
	newElements[0] = list.Elements[0] // Keep head
	copy(newElements[1:], list.Elements[startIdx+1:endIdx+2]) // +1 to adjust for head
	
	return core.List{Elements: newElements}
}

// DropList drops the first or last n elements from a list and returns the remainder
// Drop(expr, n) drops first n elements; Drop(expr, -n) drops last n elements
func DropList(list core.List, n int64) core.Expr {
	if len(list.Elements) <= 1 {
		// Empty list - return empty list with same head
		return core.List{Elements: []core.Expr{list.Elements[0]}}
	}

	listLength := len(list.Elements) - 1 // Subtract head element
	
	if n == 0 {
		// Drop 0 elements - return original list
		return list
	}

	var startIdx, endIdx int
	
	if n > 0 {
		// Drop first n elements - keep the rest
		if n >= int64(listLength) {
			// Drop all elements - return empty list with same head
			return core.List{Elements: []core.Expr{list.Elements[0]}}
		}
		startIdx = 1 + int(n)       // Start after head + n dropped elements
		endIdx = len(list.Elements) // End at the very end
	} else {
		// Drop last |n| elements - keep the beginning
		absN := -n
		if absN >= int64(listLength) {
			// Drop all elements - return empty list with same head
			return core.List{Elements: []core.Expr{list.Elements[0]}}
		}
		startIdx = 1                                     // Start after head
		endIdx = len(list.Elements) - int(absN)         // End at total - |n|
	}

	// Create new list with head + remaining elements
	newElements := make([]core.Expr, 1+(endIdx-startIdx))
	newElements[0] = list.Elements[0] // Keep the head
	copy(newElements[1:], list.Elements[startIdx:endIdx])
	
	return core.List{Elements: newElements}
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
	
	listLength := int64(len(list.Elements) - 1) // Subtract head element
	
	// Convert to 0-based indexing and handle negatives
	var actualIndex int64
	
	if index > 0 {
		actualIndex = index - 1 // Convert from 1-based
	} else if index < 0 {
		actualIndex = listLength + index // Negative indexing
	} else {
		// index == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError", 
			"Drop index 0 is out of bounds (indices start at 1)", []core.Expr{list})
	}
	
	// Bounds checking
	if actualIndex < 0 || actualIndex >= listLength {
		return core.NewErrorExpr("PartError", 
			fmt.Sprintf("Drop index %d is out of bounds for list with %d elements", 
				index, listLength), []core.Expr{list})
	}
	
	// Create result without the specified element
	newElements := make([]core.Expr, listLength) // head + (listLength - 1) elements
	newElements[0] = list.Elements[0] // Keep head
	
	// Copy elements before the dropped index
	copy(newElements[1:], list.Elements[1:actualIndex+1])
	// Copy elements after the dropped index
	copy(newElements[actualIndex+1:], list.Elements[actualIndex+2:])
	
	return core.List{Elements: newElements}
}

// dropListRange is a helper function that drops a range of elements
func dropListRange(list core.List, start, end int64) core.Expr {
	if len(list.Elements) <= 1 {
		return core.List{Elements: []core.Expr{list.Elements[0]}} // Empty list with same head
	}
	
	listLength := int64(len(list.Elements) - 1) // Subtract head element
	
	// Convert to 0-based indexing and handle negatives
	var startIdx, endIdx int64
	
	if start > 0 {
		startIdx = start - 1 // Convert from 1-based
	} else if start < 0 {
		startIdx = listLength + start // Negative indexing
	} else {
		// start == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError", 
			"Drop index 0 is out of bounds (indices start at 1)", []core.Expr{list})
	}
	
	if end > 0 {
		endIdx = end - 1 // Convert from 1-based  
	} else if end < 0 {
		endIdx = listLength + end // Negative indexing
	} else {
		// end == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError", 
			"Drop index 0 is out of bounds (indices start at 1)", []core.Expr{list})
	}
	
	// Bounds checking
	if startIdx < 0 || endIdx >= listLength || startIdx > endIdx {
		return core.NewErrorExpr("PartError", 
			fmt.Sprintf("Drop range [%d, %d] is out of bounds for list with %d elements", 
				start, end, listLength), []core.Expr{list})
	}
	
	// Create result without the specified range
	remainingSize := listLength - (endIdx - startIdx + 1)
	newElements := make([]core.Expr, remainingSize + 1) // +1 for head
	newElements[0] = list.Elements[0] // Keep head
	
	// Copy elements before the dropped range
	copy(newElements[1:], list.Elements[1:startIdx+1])
	// Copy elements after the dropped range
	copy(newElements[startIdx+1:], list.Elements[endIdx+2:])
	
	return core.List{Elements: newElements}
}
