package stdlib

import (
	"fmt"
	"math"

	"github.com/client9/sexpr/core"
)

// Mixed mathematical operations that work with Expr types
// These functions handle type conversion and return appropriate types

// PlusNumbers adds a sequence of mixed numeric expressions
func PlusNumbers(args ...core.Expr) float64 {
	sum := 0.0
	for _, arg := range args {
		if val, ok := core.GetNumericValue(arg); ok {
			sum += val
		}
		// Skip non-numeric values - they'll be caught by pattern matching
	}
	return sum
}

// TimesNumbers multiplies a sequence of mixed numeric expressions
func TimesNumbers(args ...core.Expr) float64 {
	if len(args) == 0 {
		return 1.0
	}
	product := 1.0
	for _, arg := range args {
		if val, ok := core.GetNumericValue(arg); ok {
			product *= val
		}
		// Skip non-numeric values - they'll be caught by pattern matching
	}
	return product
}

// SubtractExprs performs mixed numeric subtraction on Expr arguments (returns float64)
// Pattern constraint ensures both arguments are numeric
func SubtractExprs(x, y core.Expr) float64 {
	val1, _ := core.GetNumericValue(x)
	val2, _ := core.GetNumericValue(y)
	return val1 - val2
}

// PowerExprs performs power operation on Expr numeric arguments
// Returns (float64, error) for clear type safety
func PowerExprs(base, exp core.Expr) (float64, error) {
	if !core.IsNumeric(base) || !core.IsNumeric(exp) {
		return 0, fmt.Errorf("MathematicalError")
	}

	baseVal, _ := core.GetNumericValue(base)
	expVal, _ := core.GetNumericValue(exp)
	result := math.Pow(baseVal, expVal)

	// Check for invalid results (NaN, Inf)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, fmt.Errorf("MathematicalError")
	}

	return result, nil
}

// DivideExprs performs division on Expr numeric arguments
// Returns (float64, error) for clear type safety
func DivideExprs(x, y core.Expr) (float64, error) {
	if !core.IsNumeric(x) || !core.IsNumeric(y) {
		return 0, fmt.Errorf("MathematicalError")
	}

	val1, _ := core.GetNumericValue(x)
	val2, _ := core.GetNumericValue(y)

	if val2 == 0 {
		return 0, fmt.Errorf("DivisionByZero")
	}

	return val1 / val2, nil
}
