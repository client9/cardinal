package stdlib

import (
	"fmt"
	"unicode/utf8"

	"github.com/client9/sexpr/core"
)

// String manipulation functions

// StringLengthStr returns the UTF-8 rune count of a string
func StringLengthRunes(s string) int64 {
	return int64(utf8.RuneCountInString(s))
}

// StringTakeByte extracts the first or last n bytes from a string
// StringTakeByte(str, n) takes first n bytes; StringTakeByte(str, -n) takes last n bytes
func StringTakeByte(str string, n int64) core.Expr {
	strLength := len(str) // byte length

	if n == 0 {
		// Take 0 bytes - return empty string
		return core.NewStringAtom("")
	}

	var startIdx, endIdx int

	if n > 0 {
		// Take first n bytes
		if n > int64(strLength) {
			n = int64(strLength) // Don't take more than available
		}
		startIdx = 0
		endIdx = int(n)
	} else {
		// Take last |n| bytes
		absN := -n
		if absN > int64(strLength) {
			absN = int64(strLength) // Don't take more than available
		}
		startIdx = strLength - int(absN)
		endIdx = strLength
	}

	return core.NewStringAtom(str[startIdx:endIdx])
}

// StringTakeByteAt takes the nth byte from a string and returns it as a single-byte string
// StringTakeByteAt(str, [n]) - returns string containing the nth byte
func StringTakeByteAt(str string, indexList core.List) core.Expr {
	// Extract single integer from List(n_Integer)
	if len(indexList.Elements) != 2 { // Head + one element
		return core.NewErrorExpr("ArgumentError",
			"StringTakeByteAt with list spec requires exactly one index", []core.Expr{indexList})
	}

	index, ok := core.ExtractInt64(indexList.Elements[1])
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"StringTakeByteAt index must be an integer", []core.Expr{indexList.Elements[1]})
	}

	return stringTakeByteAtSingle(str, index)
}

// StringTakeByteRange takes a range of bytes from a string
// StringTakeByteRange(str, [n, m]) - takes bytes from index n to m (inclusive, 1-based)
func StringTakeByteRange(str string, indexList core.List) core.Expr {
	// Extract two integers from List(n_Integer, m_Integer)
	if len(indexList.Elements) != 3 { // Head + two elements
		return core.NewErrorExpr("ArgumentError",
			"StringTakeByteRange with range spec requires exactly two indices", []core.Expr{indexList})
	}

	start, ok1 := core.ExtractInt64(indexList.Elements[1])
	end, ok2 := core.ExtractInt64(indexList.Elements[2])
	if !ok1 || !ok2 {
		return core.NewErrorExpr("ArgumentError",
			"StringTakeByteRange indices must be integers", indexList.Elements[1:])
	}

	return stringTakeByteRange(str, start, end)
}

// StringDropByte drops the first or last n bytes from a string and returns the remainder
// StringDropByte(str, n) drops first n bytes; StringDropByte(str, -n) drops last n bytes
func StringDropByte(str string, n int64) core.Expr {
	strLength := len(str) // byte length

	if n == 0 {
		// Drop 0 bytes - return original string
		return core.NewStringAtom(str)
	}

	var startIdx, endIdx int

	if n > 0 {
		// Drop first n bytes - keep the rest
		if n >= int64(strLength) {
			// Drop all bytes - return empty string
			return core.NewStringAtom("")
		}
		startIdx = int(n)
		endIdx = strLength
	} else {
		// Drop last |n| bytes - keep the beginning
		absN := -n
		if absN >= int64(strLength) {
			// Drop all bytes - return empty string
			return core.NewStringAtom("")
		}
		startIdx = 0
		endIdx = strLength - int(absN)
	}

	return core.NewStringAtom(str[startIdx:endIdx])
}

// StringDropByteAt drops the nth byte from a string and returns the remainder
// StringDropByteAt(str, [n]) - removes the byte at position n
func StringDropByteAt(str string, indexList core.List) core.Expr {
	// Extract single integer from List(n_Integer)
	if len(indexList.Elements) != 2 { // Head + one element
		return core.NewErrorExpr("ArgumentError",
			"StringDropByteAt with list spec requires exactly one index", []core.Expr{indexList})
	}

	index, ok := core.ExtractInt64(indexList.Elements[1])
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"StringDropByteAt index must be an integer", []core.Expr{indexList.Elements[1]})
	}

	return stringDropByteAtSingle(str, index)
}

// StringDropByteRange drops a range of bytes from a string and returns the remainder
// StringDropByteRange(str, [n, m]) - removes bytes from index n to m (inclusive, 1-based)
func StringDropByteRange(str string, indexList core.List) core.Expr {
	// Extract two integers from List(n_Integer, m_Integer)
	if len(indexList.Elements) != 3 { // Head + two elements
		return core.NewErrorExpr("ArgumentError",
			"StringDropByteRange with range spec requires exactly two indices", []core.Expr{indexList})
	}

	start, ok1 := core.ExtractInt64(indexList.Elements[1])
	end, ok2 := core.ExtractInt64(indexList.Elements[2])
	if !ok1 || !ok2 {
		return core.NewErrorExpr("ArgumentError",
			"StringDropByteRange indices must be integers", indexList.Elements[1:])
	}

	return stringDropByteRange(str, start, end)
}

// Helper functions

// stringTakeByteAtSingle is a helper function that takes a single byte
func stringTakeByteAtSingle(str string, index int64) core.Expr {
	strLength := int64(len(str))

	if strLength == 0 {
		return core.NewErrorExpr("PartError",
			"StringTakeByteAt: string is empty", []core.Expr{core.NewStringAtom(str)})
	}

	// Convert to 0-based indexing and handle negatives
	var actualIndex int64

	if index > 0 {
		actualIndex = index - 1 // Convert from 1-based
	} else if index < 0 {
		actualIndex = strLength + index // Negative indexing
	} else {
		// index == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringTakeByteAt index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if actualIndex < 0 || actualIndex >= strLength {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringTakeByteAt index %d is out of bounds for string with %d bytes",
				index, strLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Return single byte as string
	return core.NewStringAtom(string(str[actualIndex]))
}

// stringTakeByteRange is a helper function that implements the range logic
func stringTakeByteRange(str string, start, end int64) core.Expr {
	strLength := int64(len(str))

	if strLength == 0 {
		return core.NewStringAtom("")
	}

	// Convert to 0-based indexing and handle negatives
	var startIdx, endIdx int64

	if start > 0 {
		startIdx = start - 1 // Convert from 1-based
	} else if start < 0 {
		startIdx = strLength + start // Negative indexing
	} else {
		// start == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringTakeByteRange index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	if end > 0 {
		endIdx = end - 1 // Convert from 1-based
	} else if end < 0 {
		endIdx = strLength + end // Negative indexing
	} else {
		// end == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringTakeByteRange index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if startIdx < 0 || endIdx >= strLength || startIdx > endIdx {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringTakeByteRange range [%d, %d] is out of bounds for string with %d bytes",
				start, end, strLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Return substring (endIdx+1 because Go slicing is exclusive on end)
	return core.NewStringAtom(str[startIdx : endIdx+1])
}

// stringDropByteAtSingle is a helper function that drops a single byte
func stringDropByteAtSingle(str string, index int64) core.Expr {
	strLength := int64(len(str))

	if strLength == 0 {
		return core.NewStringAtom("")
	}

	// Convert to 0-based indexing and handle negatives
	var actualIndex int64

	if index > 0 {
		actualIndex = index - 1 // Convert from 1-based
	} else if index < 0 {
		actualIndex = strLength + index // Negative indexing
	} else {
		// index == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringDropByteAt index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if actualIndex < 0 || actualIndex >= strLength {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringDropByteAt index %d is out of bounds for string with %d bytes",
				index, strLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Create result without the specified byte
	result := str[:actualIndex] + str[actualIndex+1:]
	return core.NewStringAtom(result)
}

// stringDropByteRange is a helper function that drops a range of bytes
func stringDropByteRange(str string, start, end int64) core.Expr {
	strLength := int64(len(str))

	if strLength == 0 {
		return core.NewStringAtom("")
	}

	// Convert to 0-based indexing and handle negatives
	var startIdx, endIdx int64

	if start > 0 {
		startIdx = start - 1 // Convert from 1-based
	} else if start < 0 {
		startIdx = strLength + start // Negative indexing
	} else {
		// start == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringDropByteRange index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	if end > 0 {
		endIdx = end - 1 // Convert from 1-based
	} else if end < 0 {
		endIdx = strLength + end // Negative indexing
	} else {
		// end == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringDropByteRange index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if startIdx < 0 || endIdx >= strLength || startIdx > endIdx {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringDropByteRange range [%d, %d] is out of bounds for string with %d bytes",
				start, end, strLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Create result without the specified range
	result := str[:startIdx] + str[endIdx+1:]
	return core.NewStringAtom(result)
}

// StringTake extracts the first or last n runes from a string
// StringTake(str, n) takes first n runes; StringTake(str, -n) takes last n runes
func StringTake(str string, n int64) core.Expr {
	runes := []rune(str)
	runeLength := len(runes)

	if n == 0 {
		// Take 0 runes - return empty string
		return core.NewStringAtom("")
	}

	var startIdx, endIdx int

	if n > 0 {
		// Take first n runes
		if n > int64(runeLength) {
			n = int64(runeLength) // Don't take more than available
		}
		startIdx = 0
		endIdx = int(n)
	} else {
		// Take last |n| runes
		absN := -n
		if absN > int64(runeLength) {
			absN = int64(runeLength) // Don't take more than available
		}
		startIdx = runeLength - int(absN)
		endIdx = runeLength
	}

	return core.NewStringAtom(string(runes[startIdx:endIdx]))
}

// StringTakeAt takes the nth rune from a string and returns it as a single-rune string
// StringTakeAt(str, [n]) - returns string containing the nth rune
func StringTakeAt(str string, indexList core.List) core.Expr {
	// Extract single integer from List(n_Integer)
	if len(indexList.Elements) != 2 { // Head + one element
		return core.NewErrorExpr("ArgumentError",
			"StringTakeAt with list spec requires exactly one index", []core.Expr{indexList})
	}

	index, ok := core.ExtractInt64(indexList.Elements[1])
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"StringTakeAt index must be an integer", []core.Expr{indexList.Elements[1]})
	}

	return stringTakeAtSingle(str, index)
}

// StringTakeRange takes a range of runes from a string
// StringTakeRange(str, [n, m]) - takes runes from index n to m (inclusive, 1-based)
func StringTakeRange(str string, indexList core.List) core.Expr {
	// Extract two integers from List(n_Integer, m_Integer)
	if len(indexList.Elements) != 3 { // Head + two elements
		return core.NewErrorExpr("ArgumentError",
			"StringTake with range spec requires exactly two indices", []core.Expr{indexList})
	}

	start, ok1 := core.ExtractInt64(indexList.Elements[1])
	end, ok2 := core.ExtractInt64(indexList.Elements[2])
	if !ok1 || !ok2 {
		return core.NewErrorExpr("ArgumentError",
			"StringTake indices must be integers", indexList.Elements[1:])
	}

	return stringTakeRange(str, start, end)
}

// StringDrop drops the first or last n runes from a string and returns the remainder
// StringDrop(str, n) drops first n runes; StringDrop(str, -n) drops last n runes
func StringDrop(str string, n int64) core.Expr {
	runes := []rune(str)
	runeLength := len(runes)

	if n == 0 {
		// Drop 0 runes - return original string
		return core.NewStringAtom(str)
	}

	var startIdx, endIdx int

	if n > 0 {
		// Drop first n runes - keep the rest
		if n >= int64(runeLength) {
			// Drop all runes - return empty string
			return core.NewStringAtom("")
		}
		startIdx = int(n)
		endIdx = runeLength
	} else {
		// Drop last |n| runes - keep the beginning
		absN := -n
		if absN >= int64(runeLength) {
			// Drop all runes - return empty string
			return core.NewStringAtom("")
		}
		startIdx = 0
		endIdx = runeLength - int(absN)
	}

	return core.NewStringAtom(string(runes[startIdx:endIdx]))
}

// StringDropAt drops the nth rune from a string and returns the remainder
// StringDropAt(str, [n]) - removes the rune at position n
func StringDropAt(str string, indexList core.List) core.Expr {
	// Extract single integer from List(n_Integer)
	if len(indexList.Elements) != 2 { // Head + one element
		return core.NewErrorExpr("ArgumentError",
			"StringDropAt with list spec requires exactly one index", []core.Expr{indexList})
	}

	index, ok := core.ExtractInt64(indexList.Elements[1])
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"StringDropAt index must be an integer", []core.Expr{indexList.Elements[1]})
	}

	return stringDropAtSingle(str, index)
}

// StringDropRange drops a range of runes from a string and returns the remainder
// StringDropRange(str, [n, m]) - removes runes from index n to m (inclusive, 1-based)
func StringDropRange(str string, indexList core.List) core.Expr {
	// Extract two integers from List(n_Integer, m_Integer)
	if len(indexList.Elements) != 3 { // Head + two elements
		return core.NewErrorExpr("ArgumentError",
			"StringDrop with range spec requires exactly two indices", []core.Expr{indexList})
	}

	start, ok1 := core.ExtractInt64(indexList.Elements[1])
	end, ok2 := core.ExtractInt64(indexList.Elements[2])
	if !ok1 || !ok2 {
		return core.NewErrorExpr("ArgumentError",
			"StringDrop indices must be integers", indexList.Elements[1:])
	}

	return stringDropRange(str, start, end)
}

// Helper functions for rune-based operations

// stringTakeAtSingle is a helper function that takes a single rune
func stringTakeAtSingle(str string, index int64) core.Expr {
	runes := []rune(str)
	runeLength := int64(len(runes))

	if runeLength == 0 {
		return core.NewErrorExpr("PartError",
			"StringTakeAt: string is empty", []core.Expr{core.NewStringAtom(str)})
	}

	// Convert to 0-based indexing and handle negatives
	var actualIndex int64

	if index > 0 {
		actualIndex = index - 1 // Convert from 1-based
	} else if index < 0 {
		actualIndex = runeLength + index // Negative indexing
	} else {
		// index == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringTakeAt index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if actualIndex < 0 || actualIndex >= runeLength {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringTakeAt index %d is out of bounds for string with %d characters",
				index, runeLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Return single rune as string
	return core.NewStringAtom(string(runes[actualIndex]))
}

// stringTakeRange is a helper function that implements the range logic
func stringTakeRange(str string, start, end int64) core.Expr {
	runes := []rune(str)
	runeLength := int64(len(runes))

	if runeLength == 0 {
		return core.NewStringAtom("")
	}

	// Convert to 0-based indexing and handle negatives
	var startIdx, endIdx int64

	if start > 0 {
		startIdx = start - 1 // Convert from 1-based
	} else if start < 0 {
		startIdx = runeLength + start // Negative indexing
	} else {
		// start == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringTakeRange index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	if end > 0 {
		endIdx = end - 1 // Convert from 1-based
	} else if end < 0 {
		endIdx = runeLength + end // Negative indexing
	} else {
		// end == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringTakeRange index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if startIdx < 0 || endIdx >= runeLength || startIdx > endIdx {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringTakeRange range [%d, %d] is out of bounds for string with %d characters",
				start, end, runeLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Return substring (endIdx+1 because Go slicing is exclusive on end)
	return core.NewStringAtom(string(runes[startIdx : endIdx+1]))
}

// stringDropAtSingle is a helper function that drops a single rune
func stringDropAtSingle(str string, index int64) core.Expr {
	runes := []rune(str)
	runeLength := int64(len(runes))

	if runeLength == 0 {
		return core.NewStringAtom("")
	}

	// Convert to 0-based indexing and handle negatives
	var actualIndex int64

	if index > 0 {
		actualIndex = index - 1 // Convert from 1-based
	} else if index < 0 {
		actualIndex = runeLength + index // Negative indexing
	} else {
		// index == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringDropAt index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if actualIndex < 0 || actualIndex >= runeLength {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringDropAt index %d is out of bounds for string with %d characters",
				index, runeLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Create result without the specified rune
	result := append(runes[:actualIndex], runes[actualIndex+1:]...)
	return core.NewStringAtom(string(result))
}

// stringDropRange is a helper function that drops a range of runes
func stringDropRange(str string, start, end int64) core.Expr {
	runes := []rune(str)
	runeLength := int64(len(runes))

	if runeLength == 0 {
		return core.NewStringAtom("")
	}

	// Convert to 0-based indexing and handle negatives
	var startIdx, endIdx int64

	if start > 0 {
		startIdx = start - 1 // Convert from 1-based
	} else if start < 0 {
		startIdx = runeLength + start // Negative indexing
	} else {
		// start == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringDropRange index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	if end > 0 {
		endIdx = end - 1 // Convert from 1-based
	} else if end < 0 {
		endIdx = runeLength + end // Negative indexing
	} else {
		// end == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringDropRange index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if startIdx < 0 || endIdx >= runeLength || startIdx > endIdx {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringDropRange range [%d, %d] is out of bounds for string with %d characters",
				start, end, runeLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Create result without the specified range
	result := append(runes[:startIdx], runes[endIdx+1:]...)
	return core.NewStringAtom(string(result))
}
