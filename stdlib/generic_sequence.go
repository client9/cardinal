package stdlib

import (
	"github.com/client9/sexpr/core"
)

// Generic sequence manipulation functions that work on any Sliceable expression

// Take extracts elements from a sliceable expression
// Take(expr, n) - takes first n elements if n > 0, last |n| elements if n < 0
func Take(expr core.Expr, n int64) core.Expr {
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewError("TypeError",
			"Take requires a sliceable expression (List, String, or ByteArray)")
	}

	length := expr.Length()
	if length == 0 {
		return expr // Empty expression returns itself
	}

	if n == 0 {
		// Take 0 elements - return empty version of same type
		return createEmpty(expr)
	}

	if n > 0 {
		// Take first n elements: expr.Slice(1, n)
		if n > length {
			n = length // Don't take more than available
		}
		return sliceable.Slice(1, n)
	} else {
		// Take last |n| elements: expr.Slice(length + n + 1, length)
		absN := -n
		if absN > length {
			absN = length // Don't take more than available
		}
		start := length - absN + 1
		return sliceable.Slice(start, length)
	}
}

// Drop removes elements from a sliceable expression
// Drop(expr, n) - drops first n elements if n > 0, last |n| elements if n < 0
func Drop(expr core.Expr, n int64) core.Expr {
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewError("TypeError",
			"Drop requires a sliceable expression (List, String, or ByteArray)")
	}

	length := expr.Length()
	if length == 0 {
		return expr // Empty expression returns itself
	}

	if n == 0 {
		return expr // Drop 0 elements returns original
	}

	if n > 0 {
		// Drop first n elements: expr.Slice(n + 1, length)
		if n >= length {
			return createEmpty(expr) // Drop all elements
		}
		return sliceable.Slice(n+1, length)
	} else {
		// Drop last |n| elements: expr.Slice(1, length + n)
		absN := -n
		if absN >= length {
			return createEmpty(expr) // Drop all elements
		}
		return sliceable.Slice(1, length+n) // n is negative, so this is length - absN
	}
}

// Part extracts a single element from a sliceable expression
// Part(expr, n) - returns the nth element (1-indexed)
func Part(expr core.Expr, n int64) core.Expr {
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewError("TypeError",
			"Part requires a sliceable expression (List, String, or ByteArray)")
	}

	return sliceable.ElementAt(n)
}

// Rest removes the first element from a sliceable expression
// Rest(expr) - equivalent to Drop(expr, 1)
func Rest(expr core.Expr) core.Expr {
	return Drop(expr, 1)
}

// Most removes the last element from a sliceable expression
// Most(expr) - equivalent to Drop(expr, -1)
func Most(expr core.Expr) core.Expr {
	return Drop(expr, -1)
}

// First gets the first element from a sliceable expression
// First(expr) - equivalent to Part(expr, 1)
func First(expr core.Expr) core.Expr {
	return Part(expr, 1)
}

// Last gets the last element from a sliceable expression
// Last(expr) - equivalent to Part(expr, -1)
func Last(expr core.Expr) core.Expr {
	return Part(expr, -1)
}

// TakeRange extracts a range of elements from a sliceable expression
// TakeRange(expr, [start, stop]) - takes elements from start to stop (inclusive, 1-indexed)
func TakeRange(expr core.Expr, rangeList core.List) core.Expr {
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewError("TypeError",
			"Take requires a sliceable expression (List, String, or ByteArray)")
	}

	// Extract range indices
	if len(rangeList.Elements) != 3 { // Head + two elements
		return core.NewError("ArgumentError",
			"Take with range requires exactly two indices")
	}

	start, ok1 := core.ExtractInt64(rangeList.Elements[1])
	stop, ok2 := core.ExtractInt64(rangeList.Elements[2])
	if !ok1 || !ok2 {
		return core.NewError("ArgumentError",
			"Take indices must be integers")
	}

	return sliceable.Slice(start, stop)
}

// DropRange removes a range of elements from a sliceable expression
// DropRange(expr, [start, stop]) - removes elements from start to stop (inclusive, 1-indexed)
func DropRange(expr core.Expr, rangeList core.List) core.Expr {
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewError("TypeError",
			"Drop requires a sliceable expression (List, String, or ByteArray)")
	}

	// Extract range indices
	if len(rangeList.Elements) != 3 { // Head + two elements
		return core.NewError("ArgumentError",
			"Drop with range requires exactly two indices")
	}

	start, ok1 := core.ExtractInt64(rangeList.Elements[1])
	stop, ok2 := core.ExtractInt64(rangeList.Elements[2])
	if !ok1 || !ok2 {
		return core.NewError("ArgumentError",
			"Drop indices must be integers")
	}

	length := expr.Length()
	if length == 0 {
		return expr
	}

	// Drop range by combining two slices: [1, start-1] + [stop+1, length]
	if start <= 1 && stop >= length {
		// Dropping everything
		return createEmpty(expr)
	} else if start <= 1 {
		// Dropping from beginning: keep [stop+1, length]
		return sliceable.Slice(stop+1, length)
	} else if stop >= length {
		// Dropping to end: keep [1, start-1]
		return sliceable.Slice(1, start-1)
	} else {
		// Dropping middle: need to join [1, start-1] + [stop+1, length]
		return joinSlices(expr, sliceable, 1, start-1, stop+1, length)
	}
}

// Helper functions

// createEmpty creates an empty version of the given expression type
func createEmpty(expr core.Expr) core.Expr {
	switch e := expr.(type) {
	case core.List:
		// Return list with just the head
		if len(e.Elements) > 0 {
			return core.List{Elements: []core.Expr{e.Elements[0]}}
		}
		return core.List{Elements: []core.Expr{core.NewSymbol("List")}}
	case core.String:
		return core.NewString("")
	case core.ByteArray:
		return core.NewByteArray(nil)
	default:
		return core.NewError("TypeError", "Cannot create empty version of unknown type")
	}
}

// RotateLeft rotates elements of a sliceable expression to the left by n positions
// RotateLeft(expr, n) moves the first n elements to the end
func RotateLeft(expr core.Expr, n int64) core.Expr {
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewError("TypeError",
			"RotateLeft requires a sliceable expression (List, String, or ByteArray)")
	}

	length := expr.Length()
	if length == 0 {
		return expr // Empty expressions remain unchanged
	}

	// Normalize n to be within [0, length)
	n = n % length
	if n < 0 {
		n += length
	}

	if n == 0 {
		return expr // No rotation needed
	}

	// Rotate left by n: take elements [n+1:] + [1:n]
	// For 1-based indexing: slice from (n+1) to end, then from 1 to n
	rightPart := sliceable.Slice(n+1, length)
	leftPart := sliceable.Slice(1, n)

	// Concatenate the parts
	return concatenateSliceable(rightPart, leftPart)
}

// RotateRight rotates elements of a sliceable expression to the right by n positions
// RotateRight(expr, n) moves the last n elements to the beginning
func RotateRight(expr core.Expr, n int64) core.Expr {
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewError("TypeError",
			"RotateRight requires a sliceable expression (List, String, or ByteArray)")
	}

	length := expr.Length()
	if length == 0 {
		return expr // Empty expressions remain unchanged
	}

	// Normalize n to be within [0, length)
	n = n % length
	if n < 0 {
		n += length
	}

	if n == 0 {
		return expr // No rotation needed
	}

	// Rotate right by n: take elements [length-n+1:] + [1:length-n]
	// For 1-based indexing: slice from (length-n+1) to end, then from 1 to (length-n)
	rightPart := sliceable.Slice(length-n+1, length)
	leftPart := sliceable.Slice(1, length-n)

	// Concatenate the parts
	return concatenateSliceable(rightPart, leftPart)
}

// concatenateSliceable concatenates two sliceable expressions using the Concat method
func concatenateSliceable(left, right core.Expr) core.Expr {
	// Ensure both expressions are sliceable
	leftSliceable := core.AsSliceable(left)
	if leftSliceable == nil {
		return core.NewError("TypeError",
			"Left operand is not sliceable")
	}

	rightSliceable := core.AsSliceable(right)
	if rightSliceable == nil {
		return core.NewError("TypeError",
			"Right operand is not sliceable")
	}

	// Use the Join method - this handles type checking and implementation details
	return leftSliceable.Join(rightSliceable)
}

// joinSlices joins two slices for DropRange middle case
// This is a simplified implementation - a more sophisticated version would handle all expression types
func joinSlices(expr core.Expr, sliceable core.Sliceable, start1, stop1, start2, stop2 int64) core.Expr {
	// For now, return an error for complex joins - this would need type-specific implementation
	return core.NewError("NotImplemented",
		"Dropping middle ranges not yet implemented for this type")
}
