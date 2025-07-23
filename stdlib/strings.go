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

// TakeString extracts the first or last n bytes from a string
// Take(str, n) takes first n bytes; Take(str, -n) takes last n bytes
func TakeString(str string, n int64) core.Expr {
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

// TakeStringSingle takes the nth byte from a string and returns it as a single-byte string
// Take(str, [n]) - returns string containing the nth byte
func TakeStringSingle(str string, indexList core.List) core.Expr {
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

	return takeStringSingle(str, index)
}

// TakeStringRange takes a range of bytes from a string
// Take(str, [n, m]) - takes bytes from index n to m (inclusive, 1-based)
func TakeStringRange(str string, indexList core.List) core.Expr {
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

	return takeStringRange(str, start, end)
}

// DropString drops the first or last n bytes from a string and returns the remainder
// Drop(str, n) drops first n bytes; Drop(str, -n) drops last n bytes
func DropString(str string, n int64) core.Expr {
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

// DropStringSingle drops the nth byte from a string and returns the remainder
// Drop(str, [n]) - removes the byte at position n
func DropStringSingle(str string, indexList core.List) core.Expr {
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

	return dropStringSingle(str, index)
}

// DropStringRange drops a range of bytes from a string and returns the remainder
// Drop(str, [n, m]) - removes bytes from index n to m (inclusive, 1-based)
func DropStringRange(str string, indexList core.List) core.Expr {
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

	return dropStringRange(str, start, end)
}

// Helper functions

// takeStringSingle is a helper function that takes a single byte
func takeStringSingle(str string, index int64) core.Expr {
	strLength := int64(len(str))

	if strLength == 0 {
		return core.NewErrorExpr("PartError",
			"Take: string is empty", []core.Expr{core.NewStringAtom(str)})
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
			"Take index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if actualIndex < 0 || actualIndex >= strLength {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Take index %d is out of bounds for string with %d bytes",
				index, strLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Return single byte as string
	return core.NewStringAtom(string(str[actualIndex]))
}

// takeStringRange is a helper function that implements the range logic
func takeStringRange(str string, start, end int64) core.Expr {
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
			"Take index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	if end > 0 {
		endIdx = end - 1 // Convert from 1-based
	} else if end < 0 {
		endIdx = strLength + end // Negative indexing
	} else {
		// end == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"Take index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if startIdx < 0 || endIdx >= strLength || startIdx > endIdx {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Take range [%d, %d] is out of bounds for string with %d bytes",
				start, end, strLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Return substring (endIdx+1 because Go slicing is exclusive on end)
	return core.NewStringAtom(str[startIdx : endIdx+1])
}

// dropStringSingle is a helper function that drops a single byte
func dropStringSingle(str string, index int64) core.Expr {
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
			"Drop index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if actualIndex < 0 || actualIndex >= strLength {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Drop index %d is out of bounds for string with %d bytes",
				index, strLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Create result without the specified byte
	result := str[:actualIndex] + str[actualIndex+1:]
	return core.NewStringAtom(result)
}

// dropStringRange is a helper function that drops a range of bytes
func dropStringRange(str string, start, end int64) core.Expr {
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
			"Drop index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	if end > 0 {
		endIdx = end - 1 // Convert from 1-based
	} else if end < 0 {
		endIdx = strLength + end // Negative indexing
	} else {
		// end == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"Drop index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if startIdx < 0 || endIdx >= strLength || startIdx > endIdx {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("Drop range [%d, %d] is out of bounds for string with %d bytes",
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

// StringTakeSingle takes the nth rune from a string and returns it as a single-rune string
// StringTake(str, [n]) - returns string containing the nth rune
func StringTakeSingle(str string, indexList core.List) core.Expr {
	// Extract single integer from List(n_Integer)
	if len(indexList.Elements) != 2 { // Head + one element
		return core.NewErrorExpr("ArgumentError",
			"StringTake with list spec requires exactly one index", []core.Expr{indexList})
	}

	index, ok := core.ExtractInt64(indexList.Elements[1])
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"StringTake index must be an integer", []core.Expr{indexList.Elements[1]})
	}

	return stringTakeSingle(str, index)
}

// StringTakeRange takes a range of runes from a string
// StringTake(str, [n, m]) - takes runes from index n to m (inclusive, 1-based)
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

// StringDropSingle drops the nth rune from a string and returns the remainder
// StringDrop(str, [n]) - removes the rune at position n
func StringDropSingle(str string, indexList core.List) core.Expr {
	// Extract single integer from List(n_Integer)
	if len(indexList.Elements) != 2 { // Head + one element
		return core.NewErrorExpr("ArgumentError",
			"StringDrop with list spec requires exactly one index", []core.Expr{indexList})
	}

	index, ok := core.ExtractInt64(indexList.Elements[1])
	if !ok {
		return core.NewErrorExpr("ArgumentError",
			"StringDrop index must be an integer", []core.Expr{indexList.Elements[1]})
	}

	return stringDropSingle(str, index)
}

// StringDropRange drops a range of runes from a string and returns the remainder
// StringDrop(str, [n, m]) - removes runes from index n to m (inclusive, 1-based)
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

// stringTakeSingle is a helper function that takes a single rune
func stringTakeSingle(str string, index int64) core.Expr {
	runes := []rune(str)
	runeLength := int64(len(runes))

	if runeLength == 0 {
		return core.NewErrorExpr("PartError",
			"StringTake: string is empty", []core.Expr{core.NewStringAtom(str)})
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
			"StringTake index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if actualIndex < 0 || actualIndex >= runeLength {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringTake index %d is out of bounds for string with %d characters",
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
			"StringTake index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	if end > 0 {
		endIdx = end - 1 // Convert from 1-based
	} else if end < 0 {
		endIdx = runeLength + end // Negative indexing
	} else {
		// end == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringTake index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if startIdx < 0 || endIdx >= runeLength || startIdx > endIdx {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringTake range [%d, %d] is out of bounds for string with %d characters",
				start, end, runeLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Return substring (endIdx+1 because Go slicing is exclusive on end)
	return core.NewStringAtom(string(runes[startIdx : endIdx+1]))
}

// stringDropSingle is a helper function that drops a single rune
func stringDropSingle(str string, index int64) core.Expr {
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
			"StringDrop index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if actualIndex < 0 || actualIndex >= runeLength {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringDrop index %d is out of bounds for string with %d characters",
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
			"StringDrop index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	if end > 0 {
		endIdx = end - 1 // Convert from 1-based
	} else if end < 0 {
		endIdx = runeLength + end // Negative indexing
	} else {
		// end == 0 is invalid in 1-based indexing
		return core.NewErrorExpr("PartError",
			"StringDrop index 0 is out of bounds (indices start at 1)", []core.Expr{core.NewStringAtom(str)})
	}

	// Bounds checking
	if startIdx < 0 || endIdx >= runeLength || startIdx > endIdx {
		return core.NewErrorExpr("PartError",
			fmt.Sprintf("StringDrop range [%d, %d] is out of bounds for string with %d characters",
				start, end, runeLength), []core.Expr{core.NewStringAtom(str)})
	}

	// Create result without the specified range
	result := append(runes[:startIdx], runes[endIdx+1:]...)
	return core.NewStringAtom(string(result))
}
