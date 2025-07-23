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

// SameQExprs checks if two expressions are structurally equal
func SameQExprs(x, y core.Expr) bool {
	return x.Equal(y)
}

// UnsameQExprs checks if two expressions are not structurally equal
func UnsameQExprs(x, y core.Expr) bool {
	return !x.Equal(y)
}
