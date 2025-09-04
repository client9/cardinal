package core

// Generic sequence manipulation functions that work on any Sliceable expression

// Take extracts elements from a sliceable expression
// Take(expr, n) - takes first n elements if n > 0, last |n| elements if n < 0
func Take(expr Expr, n int64) Expr {
	sliceable := AsSliceable(expr)
	if sliceable == nil {
		return NewError("TypeError",
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
func Drop(expr Expr, n int64) Expr {
	sliceable := AsSliceable(expr)
	if sliceable == nil {
		return NewError("TypeError",
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
func Part(expr Expr, n int64) Expr {
	sliceable := AsSliceable(expr)
	if sliceable == nil {
		return NewError("TypeError",
			"Part requires a sliceable expression (List, String, or ByteArray)")
	}

	return sliceable.ElementAt(n)
}

// Rest removes the first element from a sliceable expression
// Rest(expr) - equivalent to Drop(expr, 1)
func Rest(expr Expr) Expr {
	return Drop(expr, 1)
}

// Most removes the last element from a sliceable expression
// Most(expr) - equivalent to Drop(expr, -1)
func Most(expr Expr) Expr {
	return Drop(expr, -1)
}

// First gets the first element from a sliceable expression
// First(expr) - equivalent to Part(expr, 1)
func First(expr Expr) Expr {
	return Part(expr, 1)
}

// Last gets the last element from a sliceable expression
// Last(expr) - equivalent to Part(expr, -1)
func Last(expr Expr) Expr {
	return Part(expr, -1)
}

// TakeRange extracts a range of elements from a sliceable expression
// TakeRange(expr, [start, stop]) - takes elements from start to stop (inclusive, 1-indexed)
func TakeRange(expr Expr, rangeList List) Expr {
	sliceable := AsSliceable(expr)
	if sliceable == nil {
		return NewError("TypeError",
			"Take requires a sliceable expression (List, String, or ByteArray)")
	}

	// Extract range indices
	if rangeList.Length() != 2 { // Head + two elements
		return NewError("ArgumentError",
			"Take with range requires exactly two indices")
	}

	args := rangeList.Tail()

	start, ok1 := ExtractInt64(args[0])
	stop, ok2 := ExtractInt64(args[1])
	if !ok1 || !ok2 {
		return NewError("ArgumentError",
			"Take indices must be integers")
	}

	return sliceable.Slice(start, stop)
}

// DropRange removes a range of elements from a sliceable expression
// DropRange(expr, [start, stop]) - removes elements from start to stop (inclusive, 1-indexed)
func DropRange(expr Expr, rangeList List) Expr {
	sliceable := AsSliceable(expr)
	if sliceable == nil {
		return NewError("TypeError",
			"Drop requires a sliceable expression (List, String, or ByteArray)")
	}

	// Extract range indices
	if rangeList.Length() != 2 {
		return NewError("ArgumentError",
			"Drop with range requires exactly two indices")
	}
	args := rangeList.Tail()

	start, ok1 := ExtractInt64(args[0])
	stop, ok2 := ExtractInt64(args[1])
	if !ok1 || !ok2 {
		return NewError("ArgumentError",
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
func createEmpty(expr Expr) Expr {
	switch e := expr.(type) {
	case List:
		return ListFrom(e.Head())
	case String:
		return NewString("")
	case ByteArray:
		return NewByteArray(nil)
	default:
		return NewError("TypeError", "Cannot create empty version of unknown type")
	}
}

// concatenateSliceable concatenates two sliceable expressions using the Concat method
func ConcatenateSliceable(left, right Expr) Expr {
	// Ensure both expressions are sliceable
	leftSliceable := AsSliceable(left)
	if leftSliceable == nil {
		return NewError("TypeError",
			"Left operand is not sliceable")
	}

	rightSliceable := AsSliceable(right)
	if rightSliceable == nil {
		return NewError("TypeError",
			"Right operand is not sliceable")
	}

	// Use the Join method - this handles type checking and implementation details
	return leftSliceable.Join(rightSliceable)
}

// joinSlices joins two slices for DropRange middle case
// This is a simplified implementation - a more sophisticated version would handle all expression types
func joinSlices(expr Expr, sliceable Sliceable, start1, stop1, start2, stop2 int64) Expr {
	// For now, return an error for complex joins - this would need type-specific implementation
	return NewError("NotImplemented",
		"Dropping middle ranges not yet implemented for this type")
}
