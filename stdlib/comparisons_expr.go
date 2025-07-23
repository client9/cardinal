package stdlib

import (
	"github.com/client9/sexpr/core"
)

// Comparison operations that work with Expr types
// These functions handle type conversion and return bool values

// EqualExprs checks if two expressions are equal
func EqualExprs(x, y core.Expr) bool {
	return x.Equal(y)
}

// UnequalExprs checks if two expressions are not equal
func UnequalExprs(x, y core.Expr) bool {
	return !x.Equal(y)
}

// LessExprs checks if x < y for numeric types
func LessExprs(x, y core.Expr) bool {
	// Extract numeric values
	val1, ok1 := core.GetNumericValue(x)
	val2, ok2 := core.GetNumericValue(y)
	if ok1 && ok2 {
		return val1 < val2
	}
	return false // Fallback case - will be handled by wrapper
}

// GreaterExprs checks if x > y for numeric types
func GreaterExprs(x, y core.Expr) bool {
	// Extract numeric values
	val1, ok1 := core.GetNumericValue(x)
	val2, ok2 := core.GetNumericValue(y)
	if ok1 && ok2 {
		return val1 > val2
	}
	return false // Fallback case - will be handled by wrapper
}

// LessEqualExprs checks if x <= y for numeric types
func LessEqualExprs(x, y core.Expr) bool {
	// Extract numeric values
	val1, ok1 := core.GetNumericValue(x)
	val2, ok2 := core.GetNumericValue(y)
	if ok1 && ok2 {
		return val1 <= val2
	}
	return false // Fallback case - will be handled by wrapper
}

// GreaterEqualExprs checks if x >= y for numeric types
func GreaterEqualExprs(x, y core.Expr) bool {
	// Extract numeric values
	val1, ok1 := core.GetNumericValue(x)
	val2, ok2 := core.GetNumericValue(y)
	if ok1 && ok2 {
		return val1 >= val2
	}
	return false // Fallback case - will be handled by wrapper
}

// SameQExprs checks if two expressions are structurally equal
func SameQExprs(x, y core.Expr) bool {
	return x.Equal(y)
}

// UnsameQExprs checks if two expressions are not structurally equal
func UnsameQExprs(x, y core.Expr) bool {
	return !x.Equal(y)
}