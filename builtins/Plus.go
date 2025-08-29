package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Plus
// @ExprAttributes Flat Listable NumericFunction OneIdentity Orderless Protected

// PlusExpr performs addition with light simplification
// Combines all numeric values and leaves symbolic terms unchanged
// Returns integers when all args are integers, float64 when mixed
// @ExprPattern (___)
func PlusExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
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

	var resultElements = make([]core.Expr, 0, 4+len(nonNumeric))
	resultElements = append(resultElements, symbol.Plus)

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

	return core.NewListFromExprs(resultElements...)
}
