package stdlib

import (
	"fmt"
	"math"

	"github.com/client9/sexpr/core"
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

// MinusInteger returns the negation of an integer
func MinusInteger(x int64) int64 {
	return -x
}

// MinusReal returns the negation of a real number
func MinusReal(x float64) float64 {
	return -x
}

// MinusExpr converts Minus(x) to Times(-1, x) as per Mathematica
func MinusExpr(x core.Expr) core.Expr {
	return core.NewList("Times", core.NewInteger(-1), x)
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

// PowerInteger - if exp >= 0, then it returns an integer, if exp < 0 returns the float value or error
func PowerInteger(base, exp int64) (core.Expr, error) {
	if exp == 0 {
		return core.NewInteger(1), nil
	}
	if exp < 0 {
		val, err := PowerNumbers(float64(base), float64(exp))
		if err != nil {
			return core.NewSymbolNull(), err
		}
		return core.NewReal(val), nil
	}
	val, err := PowerNumbers(float64(base), float64(exp))
	if err != nil {
		return core.NewSymbolNull(), err
	}
	return core.NewInteger(int64(val)), nil
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

// PlusExpr performs addition with light simplification
// Combines all numeric values and leaves symbolic terms unchanged
// Returns integers when all args are integers, float64 when mixed
func PlusExpr(args ...core.Expr) core.Expr {
	if len(args) == 0 {
		return core.NewInteger(0) // Plus() = 0
	}

	var intSum int64 = 0
	var realSum float64 = 0.0
	var hasIntegers bool = false
	var hasReals bool = false
	var nonNumeric []core.Expr

	// Separate numeric and non-numeric arguments
	for _, arg := range args {
		if intVal, ok := core.ExtractInt64(arg); ok {
			intSum += intVal
			hasIntegers = true
		} else if realVal, ok := core.ExtractFloat64(arg); ok {
			realSum += realVal
			hasReals = true
		} else {
			nonNumeric = append(nonNumeric, arg)
		}
	}

	// Build result elements
	var resultElements []core.Expr
	resultElements = append(resultElements, core.NewSymbol("Plus"))

	// Add combined numeric value
	if hasIntegers && hasReals {
		// Mixed: convert to float64
		totalSum := float64(intSum) + realSum
		if totalSum != 0.0 || len(nonNumeric) == 0 {
			resultElements = append(resultElements, core.NewReal(totalSum))
		}
	} else if hasIntegers {
		// All integers: keep as integer
		if intSum != 0 || len(nonNumeric) == 0 {
			resultElements = append(resultElements, core.NewInteger(intSum))
		}
	} else if hasReals {
		// All reals: keep as float64
		if realSum != 0.0 || len(nonNumeric) == 0 {
			resultElements = append(resultElements, core.NewReal(realSum))
		}
	}

	// Add non-numeric terms
	resultElements = append(resultElements, nonNumeric...)

	// Apply OneIdentity-like behavior: if only one element (plus head), return it
	if len(resultElements) == 2 {
		return resultElements[1]
	}

	// If no elements besides head, return 0
	if len(resultElements) == 1 {
		return core.NewInteger(0)
	}

	return core.List{Elements: resultElements}
}

// TimesExpr performs multiplication with light simplification
// Combines all numeric values and leaves symbolic terms unchanged
func TimesExpr(args ...core.Expr) core.Expr {
	if len(args) == 0 {
		return core.NewInteger(1) // Times() = 1
	}

	var intProduct int64 = 1
	var realProduct float64 = 1.0
	var hasIntegers bool = false
	var hasReals bool = false
	var nonNumeric []core.Expr

	// Separate numeric and non-numeric arguments
	for _, arg := range args {
		if intVal, ok := core.ExtractInt64(arg); ok {
			intProduct *= intVal
			hasIntegers = true
		} else if realVal, ok := core.ExtractFloat64(arg); ok {
			realProduct *= realVal
			hasReals = true
		} else {
			nonNumeric = append(nonNumeric, arg)
		}
	}

	// Check for zero (short-circuit)
	if (hasIntegers && intProduct == 0) || (hasReals && realProduct == 0.0) {
		return core.NewInteger(0)
	}

	// Build result elements
	var resultElements []core.Expr
	resultElements = append(resultElements, core.NewSymbol("Times"))

	// Add combined numeric value
	if hasIntegers && hasReals {
		// Mixed: convert to float64
		totalProduct := float64(intProduct) * realProduct
		if totalProduct != 1.0 || len(nonNumeric) == 0 {
			resultElements = append(resultElements, core.NewReal(totalProduct))
		}
	} else if hasIntegers {
		// All integers: keep as integer
		if intProduct != 1 || len(nonNumeric) == 0 {
			resultElements = append(resultElements, core.NewInteger(intProduct))
		}
	} else if hasReals {
		// All reals: keep as float64
		if realProduct != 1.0 || len(nonNumeric) == 0 {
			resultElements = append(resultElements, core.NewReal(realProduct))
		}
	}

	// Add non-numeric terms
	resultElements = append(resultElements, nonNumeric...)

	// Apply OneIdentity-like behavior: if only one element (plus head), return it
	if len(resultElements) == 2 {
		return resultElements[1]
	}

	// If no elements besides head, return 1
	if len(resultElements) == 1 {
		return core.NewInteger(1)
	}

	return core.List{Elements: resultElements}
}
