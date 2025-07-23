package stdlib

// Comparison operators for basic types - all return bool

// EqualInts checks if two integers are equal
func EqualInts(x, y int64) bool {
	return x == y
}

// EqualFloats checks if two floats are equal
func EqualFloats(x, y float64) bool {
	return x == y
}

// EqualStrings checks if two strings are equal
func EqualStrings(x, y string) bool {
	return x == y
}

// UnequalInts checks if two integers are not equal
func UnequalInts(x, y int64) bool {
	return x != y
}

// UnequalFloats checks if two floats are not equal
func UnequalFloats(x, y float64) bool {
	return x != y
}

// UnequalStrings checks if two strings are not equal
func UnequalStrings(x, y string) bool {
	return x != y
}

// LessExprs checks if x < y for numeric types
// func LessExprs(x, y Expr) bool {
// 	// Extract numeric values
// 	val1, ok1 := getNumericValue(x)
// 	val2, ok2 := getNumericValue(y)
// 	if ok1 && ok2 {
// 		return val1 < val2
// 	}
// 	return false // Fallback case - will be handled by wrapper
// }

// GreaterExprs checks if x > y for numeric types
// func GreaterExprs(x, y Expr) bool {
// 	// Extract numeric values
// 	val1, ok1 := getNumericValue(x)
// 	val2, ok2 := getNumericValue(y)
// 	if ok1 && ok2 {
// 		return val1 > val2
// 	}
// 	return false // Fallback case - will be handled by wrapper
// }

// LessEqualExprs checks if x <= y for numeric types
// func LessEqualExprs(x, y Expr) bool {
// 	// Extract numeric values
// 	val1, ok1 := getNumericValue(x)
// 	val2, ok2 := getNumericValue(y)
// 	if ok1 && ok2 {
// 		return val1 <= val2
// 	}
// 	return false // Fallback case - will be handled by wrapper
// }

// GreaterEqualExprs checks if x >= y for numeric types
// func GreaterEqualExprs(x, y Expr) bool {
// 	// Extract numeric values
// 	val1, ok1 := getNumericValue(x)
// 	val2, ok2 := getNumericValue(y)
// 	if ok1 && ok2 {
// 		return val1 >= val2
// 	}
// 	return false // Fallback case - will be handled by wrapper
// }

// SameQExprs checks if two expressions are structurally equal
// func SameQExprs(x, y Expr) bool {
// 	return x.Equal(y)
// }

// UnsameQExprs checks if two expressions are not structurally equal
// func UnsameQExprs(x, y Expr) bool {
// 	return !x.Equal(y)
// }