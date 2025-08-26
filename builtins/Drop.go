package builtins

import (
	"fmt"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Drop
// TODO: isn't using generic version

// DropList drops the first or last n elements from a list and returns the remainder
// Drop(expr, n) drops first n elements; Drop(expr, -n) drops last n elements
//
// @ExprPattern (_List, _Integer)
func DropList(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	n, _ := core.ExtractInt64(args[1])
	return core.Drop(expr, n)
}

// DropListSingle drops the nth element from a list and returns the remainder
// Drop(expr, [n]) - removes the element at position n
//
// @ExprPattern (_, [_Integer])
func DropListSingle(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	list := args[1].(core.List)
	n := list.Tail()[0]
	//n,_ := core.ExtractInt64(list.Tail()[0])
	return core.DropRange(expr, core.NewList("List", n))
}

// DropListRange drops a range of elements from a list and returns the remainder
// Drop(expr, [n, m]) - removes elements from index n to m (inclusive)
// @ExprPattern (_, [_Integer, _Integer])
func DropListRange(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	list := args[1].(core.List)
	return core.DropRange(expr, list)
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
