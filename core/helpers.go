package core

import (
	"fmt"
	"strings"
)

// Type extraction helper functions for builtin function wrappers

// ExtractInt64 safely extracts an int64 value from an Expr
func ExtractInt64(expr Expr) (int64, bool) {
	if atom, ok := expr.(Atom); ok && atom.AtomType == IntAtom {
		return int64(atom.Value.(int)), true
	}
	return 0, false
}

// ExtractFloat64 safely extracts a float64 value from an Expr
func ExtractFloat64(expr Expr) (float64, bool) {
	if atom, ok := expr.(Atom); ok && atom.AtomType == FloatAtom {
		return atom.Value.(float64), true
	}
	return 0, false
}

// ExtractString safely extracts a string value from an Expr
func ExtractString(expr Expr) (string, bool) {
	if atom, ok := expr.(Atom); ok && atom.AtomType == StringAtom {
		return atom.Value.(string), true
	}
	return "", false
}

// ExtractBool safely extracts a boolean value from an Expr
// Note: NewBoolAtom returns symbols "True"/"False", so we check for those
func ExtractBool(expr Expr) (bool, bool) {
	if atom, ok := expr.(Atom); ok && atom.AtomType == SymbolAtom {
		if atom.Value.(string) == "True" {
			return true, true
		} else if atom.Value.(string) == "False" {
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
	elements[0] = NewSymbolAtom(head)
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
	if atom, ok := expr.(Atom); ok {
		switch atom.AtomType {
		case IntAtom:
			return float64(atom.Value.(int)), true
		case FloatAtom:
			return atom.Value.(float64), true
		}
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
	if atom, ok := expr.(Atom); ok && atom.AtomType == SymbolAtom {
		val := atom.Value.(string)
		return val == "True" || val == "False"
	}
	return false
}

// IsSymbol checks if an expression is a symbol
func IsSymbol(expr Expr) bool {
	if atom, ok := expr.(Atom); ok {
		return atom.AtomType == SymbolAtom
	}
	return false
}
