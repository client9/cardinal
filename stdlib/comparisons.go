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

// EqualNumbers has input casted to float64 before
func EqualNumbers(x, y Number) bool {
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

// UnequalFloats checks if two numbers are not equal.
//
//	Automatically casted to float64.
func UnequalNumbers(x, y Number) bool {
	return x != y
}

// UnequalStrings checks if two strings are not equal
func UnequalStrings(x, y string) bool {
	return x != y
}

// LessNumber checks if x < y for numeric types
func LessNumber(x, y Number) bool {
	return x < y
}

// LessEqualNumber checks if x <= y for numeric types
func LessEqualNumber(x, y Number) bool {
	return x <= y
}

// GreaterNumber checks if x < y for numeric types
func GreaterNumber(x, y Number) bool {
	return x > y
}

// GreaterEqualNumber checks if x <= y for numeric types
func GreaterEqualNumber(x, y Number) bool {
	return x >= y
}
