package core

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// String type implementation
func (s String) String() string {
	return fmt.Sprintf("\"%s\"", string(s))
}

func (s String) InputForm() string {
	return s.String()
}

func (s String) Head() string {
	return "String"
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

// Integer type implementation
func (i Integer) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Integer) InputForm() string {
	return i.String()
}

func (i Integer) Head() string {
	return "Integer"
}

func (i Integer) Length() int64 {
	return 0
}

func (i Integer) Equal(rhs Expr) bool {
	if other, ok := rhs.(Integer); ok {
		return i == other
	}
	return false
}

func (i Integer) IsAtom() bool {
	return true
}

// Real type implementation
func (r Real) String() string {
	str := strconv.FormatFloat(float64(r), 'f', -1, 64)
	if !strings.Contains(str, ".") {
		str += ".0"
	}
	return str
}

func (r Real) InputForm() string {
	return r.String()
}

func (r Real) Head() string {
	return "Real"
}

func (r Real) Length() int64 {
	return 0
}

func (r Real) Equal(rhs Expr) bool {
	if other, ok := rhs.(Real); ok {
		return r == other
	}
	return false
}

func (r Real) IsAtom() bool {
	return true
}

// Symbol type implementation
func (s Symbol) String() string {
	return string(s)
}

func (s Symbol) InputForm() string {
	return s.String()
}

func (s Symbol) Head() string {
	return "Symbol"
}

func (s Symbol) Length() int64 {
	return 0
}

func (s Symbol) Equal(rhs Expr) bool {
	if other, ok := rhs.(Symbol); ok {
		return s == other
	}
	return false
}

func (s Symbol) IsAtom() bool {
	return true
}

// Sliceable interface implementation for String
func (s String) ElementAt(n int64) Expr {
	str := string(s)
	runes := []rune(str)
	length := int64(len(runes))

	if length == 0 {
		return NewErrorExpr("PartError", "String is empty", []Expr{s})
	}

	// Handle negative indexing
	if n < 0 {
		n = length + n + 1
	}

	// Check bounds (1-indexed)
	if n <= 0 || n > length {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Part index %d is out of bounds for string with %d characters", n, length),
			[]Expr{s})
	}

	// Convert to 0-based index and return character as string
	return String(string(runes[n-1]))
}

func (s String) Slice(start, stop int64) Expr {
	str := string(s)
	runes := []rune(str)
	length := int64(len(runes))

	if length == 0 {
		return s
	}

	// Handle negative indexing
	if start < 0 {
		start = length + start + 1
	}
	if stop < 0 {
		stop = length + stop + 1
	}

	// Check bounds
	if start <= 0 || stop <= 0 || start > length || stop > length {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Slice indices [%d, %d] out of bounds for string with %d characters",
				start, stop, length), []Expr{s})
	}

	if start > stop {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Start index %d is greater than stop index %d", start, stop),
			[]Expr{s})
	}

	// Convert to 0-based indices for Go slice
	startIdx := start - 1
	stopIdx := stop

	return String(string(runes[startIdx:stopIdx]))
}

func (s String) Join(other Sliceable) Expr {
	// Type check: ensure other is also a String
	otherStr, ok := other.(String)
	if !ok {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Cannot join %T with String", other),
			[]Expr{s, other.(Expr)})
	}

	// Simple concatenation for strings
	return String(string(s) + string(otherStr))
}

func (s String) SetElementAt(n int64, value Expr) Expr {
	// Validate that value is a string
	valueStr, ok := value.(String)
	if !ok {
		return NewErrorExpr("TypeError",
			"String assignment requires string value", []Expr{s, value})
	}

	str := string(s)
	runes := []rune(str)
	length := int64(len(runes))

	if length == 0 {
		return NewErrorExpr("PartError", "String is empty", []Expr{s})
	}

	// Handle negative indexing
	if n < 0 {
		n = length + n + 1
	}

	// Check bounds (1-indexed)
	if n <= 0 || n > length {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Part index %d is out of bounds for string with %d characters", n, length),
			[]Expr{s})
	}

	valueRunes := []rune(string(valueStr))

	// For single character replacement, value should be a single character
	if len(valueRunes) != 1 {
		return NewErrorExpr("ValueError",
			fmt.Sprintf("Single character replacement requires exactly one character, got %d", len(valueRunes)),
			[]Expr{s, value})
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
		return NewErrorExpr("TypeError",
			"String slice assignment requires string value", []Expr{s, values})
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
		return NewErrorExpr("PartError", "String is empty", []Expr{s})
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
		return NewErrorExpr("PartError",
			fmt.Sprintf("Start index %d must be positive", start),
			[]Expr{s})
	}

	if start > length+1 {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Start index %d is out of bounds for string with %d characters", start, length),
			[]Expr{s})
	}

	// Handle special cases
	if stop < start-1 {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Stop index %d cannot be less than start index %d - 1", stop, start),
			[]Expr{s})
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
