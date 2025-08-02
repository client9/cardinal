package core

import (
	"fmt"
	"strings"
)

// Type extraction helper functions for builtin function wrappers

func AsError(expr Expr) (*ErrorExpr, bool) {
	if err, ok := expr.(*ErrorExpr); ok {
		return err, ok
	}
	return &ErrorExpr{}, false
}

// ExtractInt64 safely extracts an int64 value from an Expr
func ExtractInt64(expr Expr) (int64, bool) {
	// Check new Integer type first
	if i, ok := expr.(Integer); ok {
		return int64(i), true
	}
	return 0, false
}

// ExtractFloat64 safely extracts a float64 value from an Expr
func ExtractFloat64(expr Expr) (float64, bool) {
	// Check new Real type first
	if r, ok := expr.(Real); ok {
		return float64(r), true
	}
	return 0, false
}

// ExtractString safely extracts a string value from an Expr
func ExtractString(expr Expr) (string, bool) {
	// Check new String type first
	if s, ok := expr.(String); ok {
		return string(s), true
	}
	return "", false
}

// ExtractBool safely extracts a boolean value from an Expr
// Note: NewBool returns symbols "True"/"False", so we check for those
func ExtractBool(expr Expr) (bool, bool) {
	// Check new Symbol type first
	if s, ok := expr.(Symbol); ok {
		symbolName := string(s)
		if symbolName == "True" {
			return true, true
		} else if symbolName == "False" {
			return false, true
		}
	}
	return false, false
}

// ExtractByteArray safely extracts an ByteArray value from an Expr
func ExtractByteArray(expr Expr) (ByteArray, bool) {
	if ba, ok := expr.(ByteArray); ok {
		return ba, true
	}
	return ByteArray{}, false
}

// CopyExprList creates a new List expression from a head symbol and arguments
// This is useful for builtin functions that need to return unchanged expressions
func CopyExprList(head string, args []Expr) List {
	elements := make([]Expr, len(args)+1)
	elements[0] = NewSymbol(head) // Use new Symbol constructor
	copy(elements[1:], args)
	return List{Elements: elements}
}

// IsError checks if an expression is an error
func IsError(expr Expr) bool {
	_, ok := expr.(*ErrorExpr)
	return ok
}

// GetStackTrace returns a formatted string representation of the stack trace
func (e *ErrorExpr) GetStackTrace() string {
	if len(e.StackTrace) == 0 {
		return "No stack trace available"
	}

	var trace strings.Builder
	trace.WriteString("Stack trace:\n")
	for i, frame := range e.StackTrace {
		trace.WriteString(fmt.Sprintf("  %d. %s: %s", i+1, frame.Function, frame.Expression))
		if frame.Location != "" {
			trace.WriteString(fmt.Sprintf(" at %s", frame.Location))
		}
		trace.WriteString("\n")
	}
	return trace.String()
}

// GetDetailedMessage returns the error message with stack trace
func (e *ErrorExpr) GetDetailedMessage() string {
	return fmt.Sprintf("Error: %s\nMessage: %s\n%s", e.ErrorType, e.Message, e.GetStackTrace())
}

// GetNumericValue safely extracts a numeric value (int or float) as float64 from an Expr
func GetNumericValue(expr Expr) (float64, bool) {
	// Check new atomic types first
	if i, ok := expr.(Integer); ok {
		return float64(i), true
	}
	if r, ok := expr.(Real); ok {
		return float64(r), true
	}
	return 0, false
}

// IsNumeric checks if an expression represents a numeric value (int or float)
func IsNumeric(expr Expr) bool {
	_, ok := GetNumericValue(expr)
	return ok
}

// IsBool checks if an expression is a boolean value (True/False symbol)
func IsBool(expr Expr) bool {
	// Check new Symbol type first
	if s, ok := expr.(Symbol); ok {
		val := string(s)
		return val == "True" || val == "False"
	}
	return false
}

// IsSymbol checks if an expression is a symbol
func IsSymbol(expr Expr) bool {
	// Check new Symbol type first
	if _, ok := expr.(Symbol); ok {
		return true
	}
	return false
}

// ExtractSymbol safely extracts a symbol name from an Expr
func ExtractSymbol(expr Expr) (string, bool) {
	// Check new Symbol type first
	if s, ok := expr.(Symbol); ok {
		return string(s), true
	}
	return "", false
}

// CanonicalCompare provides a canonical comparison function for expressions
// Used for consistent ordering across mathematical functions and Orderless attribute
// Returns true if expr1 should come before expr2 in canonical ordering
func CanonicalCompare(expr1, expr2 Expr) bool {
	// Mathematical ordering: numbers first, then other expressions
	_, expr1IsNumber := GetNumericValue(expr1)
	_, expr2IsNumber := GetNumericValue(expr2)

	if expr1IsNumber && !expr2IsNumber {
		return true // Numbers come before non-numbers
	}
	if !expr1IsNumber && expr2IsNumber {
		return false // Non-numbers come after numbers
	}

	// Both are numbers or both are non-numbers, use standard ordering
	// First compare by length (complexity)
	cmp := expr1.Length() - expr2.Length()
	if cmp < 0 {
		return true
	}
	if cmp > 0 {
		return false
	}

	// If lengths are equal, compare by string representation for deterministic ordering
	return expr1.String() < expr2.String()
}
