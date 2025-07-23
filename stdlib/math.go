package stdlib

import (
	"fmt"
	"math"
)

// Arithmetic functions for mathematical operations
// These are pure Go functions that are automatically wrapped

// PlusIdentity returns the additive identity (0), which is returned by Plus()
func PlusIdentity() int64 {
	return 0
}

// TimesIdentity returns the multiplicative identity (1), which is returned by Times()
func TimesIdentity() int64 {
	return 1
}

// PlusIntegers adds a sequence of integers
func PlusIntegers(args ...int64) int64 {
	sum := int64(0)
	for _, v := range args {
		sum += v
	}
	return sum
}

// TimesIntegers multiplies a sequence of integers
func TimesIntegers(args ...int64) int64 {
	if len(args) == 0 {
		return 1
	}
	product := int64(1)
	for _, v := range args {
		product *= v
	}
	return product
}

// PlusReals adds a sequence of real numbers
func PlusReals(args ...float64) float64 {
	sum := 0.0
	for _, v := range args {
		sum += v
	}
	return sum
}

// TimesReals multiplies a sequence of real numbers
func TimesReals(args ...float64) float64 {
	if len(args) == 0 {
		return 1.0
	}
	product := 1.0
	for _, v := range args {
		product *= v
	}
	return product
}

// PowerReal computes real base to integer exponent
func PowerReal(base float64, exp int64) float64 {
	if exp == 0 {
		return 1.0
	}
	if exp < 0 {
		return 1.0 / PowerReal(base, -exp)
	}

	result := 1.0
	for i := int64(0); i < exp; i++ {
		result *= base
	}
	return result
}

// SubtractIntegers performs integer subtraction
func SubtractIntegers(x, y int64) int64 {
	return x - y
}

// DivideIntegers performs integer division on int64 arguments
// Returns (int64, error) for clear type safety using Go's integer division
func DivideIntegers(x, y int64) (int64, error) {
	if y == 0 {
		return 0, fmt.Errorf("DivisionByZero")
	}

	return x / y, nil
}

// PowerNumbers performs power operation on numeric arguments
// Returns (float64, error) for clear type safety
func PowerNumbers(base, exp float64) (float64, error) {
	result := math.Pow(base, exp)

	// Check for invalid results (NaN, Inf)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, fmt.Errorf("MathematicalError")
	}

	return result, nil
}

// DivideNumbers performs division on numeric arguments
// Returns (float64, error) for clear type safety
func DivideNumbers(x, y float64) (float64, error) {
	if y == 0 {
		return 0, fmt.Errorf("DivisionByZero")
	}

	return x / y, nil
}

// SubtractNumbers performs mixed numeric subtraction (returns float64)
func SubtractNumbers(x, y float64) float64 {
	return x - y
}
