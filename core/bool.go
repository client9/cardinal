package core

import (
	"github.com/client9/sexpr/core/symbol"
)

func NewBool(value bool) Symbol {
	if value {
		return symbol.True
	} else {
		return symbol.False
	}
}

// IsBool checks if an expression is a boolean value (True/False symbol.)
func IsBool(expr Expr) bool {
	// Check new Symbol type first
	if s, ok := expr.(Symbol); ok {
		return s == symbol.True || s == symbol.False
	}
	return false
}

// ExtractBool safely extracts a boolean value from an Expr
// Note: NewBool returns symbol.s "True"/"False", so we check for those
func ExtractBool(expr Expr) (bool, bool) {
	// Check new Symbol type first
	if s, ok := expr.(Symbol); ok {
		if s == symbol.True {
			return true, true
		}
		if s == symbol.False {
			return false, true
		}
	}
	return false, false
}
