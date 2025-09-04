package core

import (
	"fmt"
	"unicode/utf8"
)

type String string

func NewString(s string) String { return String(s) }

// String type implementation
func (s String) String() string {
	return fmt.Sprintf("\"%s\"", string(s))
}

func (s String) InputForm() string {
	return s.String()
}

func (s String) Head() Expr {
	return symbolString
}

func (s String) Length() int64 {
	return int64(utf8.RuneCountInString(string(s)))
}

func (s String) Equal(rhs Expr) bool {
	if other, ok := rhs.(String); ok {
		return s == other
	}
	return false
}

func (s String) IsAtom() bool {
	return true
}

// Sliceable interface implementation for String
func (s String) ElementAt(n int64) Expr {
	runes := []rune(s)
	r, err := ElementAt(runes, int(n))
	if err != nil {
		return NewError(err.Error(), "")
	}
	return String(string(r))
}

func (s String) Slice(start, stop int64) Expr {
	runes := []rune(s)
	r, err := Slice(runes, int(start), int(stop))
	if err != nil {
		return NewError(err.Error(), "")
	}
	return String(string(r))
}

func (s String) Join(other Sliceable) Expr {
	// Type check: ensure other is also a String
	otherStr, ok := other.(String)
	if !ok {
		return NewError("TypeError",
			fmt.Sprintf("Cannot join %T with String", other))
	}

	// Simple concatenation for strings .. no need for runes and generics
	return String(string(s) + string(otherStr))
}

func (s String) SetElementAt(n int64, value Expr) Expr {
	// Validate that value is a string
	valueStr, ok := value.(String)
	if !ok {
		return NewError("TypeError",
			"String assignment requires string value")
	}

	str := string(s)
	runes := []rune(str)
	length := int64(len(runes))

	if length == 0 {
		return NewError("PartError", "String is empty")
	}

	// Handle negative indexing
	if n < 0 {
		n = length + n + 1
	}

	// Check bounds (1-indexed)
	if n <= 0 || n > length {
		return NewError("PartError",
			fmt.Sprintf("Part index %d is out of bounds for string with %d characters", n, length))
	}

	valueRunes := []rune(string(valueStr))

	// For single character replacement, value should be a single character
	if len(valueRunes) != 1 {
		return NewError("ValueError", "Single character replacement requires exactly one character")
	}

	// Create new string with character replaced
	newRunes := make([]rune, length)
	copy(newRunes, runes)
	newRunes[n-1] = valueRunes[0] // n is 1-indexed

	return String(string(newRunes))
}

func (s String) SetSlice(start, stop int64, values Expr) Expr {
	// Validate that values is a string
	valueStr, ok := values.(String)
	if !ok {
		return NewError("TypeError",
			"String slice assignment requires string value")
	}

	str := string(s)
	runes := []rune(str)
	length := int64(len(runes))

	// Handle empty string
	if length == 0 {
		if start == 1 && stop == 0 {
			// Insert at beginning of empty string
			return values
		}
		return NewError("PartError", "String is empty")
	}

	// Handle negative indexing
	if start < 0 {
		start = length + start + 1
	}
	if stop < 0 {
		stop = length + stop + 1
	}

	// Validate range
	if start <= 0 {
		return NewError("PartError",
			fmt.Sprintf("Start index %d must be positive", start))
	}

	if start > length+1 {
		return NewError("PartError",
			fmt.Sprintf("Start index %d is out of bounds for string with %d characters", start, length))
	}

	// Handle special cases
	if stop < start-1 {
		return NewError("PartError",
			fmt.Sprintf("Stop index %d cannot be less than start index %d - 1", stop, start))
	}

	valueRunes := []rune(string(valueStr))

	// Calculate range to replace
	replaceStart := start - 1 // Convert to 0-based
	replaceEnd := stop        // Convert to 0-based (exclusive end)
	if stop > length {
		replaceEnd = length
	}

	// Build new string
	var newRunes []rune

	// Add characters before the replacement range
	if replaceStart > 0 {
		newRunes = append(newRunes, runes[:replaceStart]...)
	}

	// Add replacement characters
	newRunes = append(newRunes, valueRunes...)

	// Add characters after the replacement range
	if replaceEnd < length {
		newRunes = append(newRunes, runes[replaceEnd:]...)
	}

	return String(string(newRunes))
}
