package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Times
// @ExprAttributes Flat Orderless OneIdentity NumericFunction

// TimesExpr performs multiplication with light simplification
// Combines all numeric values and leaves symbolic terms unchanged
// @ExprPattern (___)
func TimesExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) == 0 {
		return core.NewInteger(1) // Times() = 1
	}
	return core.TimesList(args)

	/*
	   intProduct := core.NewProductInteger()
	   var realProduct float64 = 1.0
	   var hasIntegers bool = false
	   var hasReals bool = false
	   var nonNumeric []core.Expr

	   // Separate numeric and non-numeric arguments

	   	for _, arg := range args {
	   		if intVal, ok := arg.(core.Integer); ok {
	   			intProduct.Update(intVal)
	   			hasIntegers = true
	   		} else if realVal, ok := core.ExtractFloat64(arg); ok {
	   			realProduct *= realVal
	   			hasReals = true
	   		} else {
	   			nonNumeric = append(nonNumeric, arg)
	   		}
	   	}

	   product := intProduct.Integer()
	   // Check for zero (short-circuit)

	   	if (hasIntegers && product.IsInt64() && product.Int64() == 0) || (hasReals && realProduct == 0.0) {
	   		return core.NewInteger(0)
	   	}

	   // Build result elements
	   var resultElements []core.Expr
	   resultElements = append(resultElements, symbol.Times)

	   // Add combined numeric value

	   	if hasIntegers && hasReals {
	   		// Mixed: convert to float64
	   		totalProduct := intProduct.Integer().Float64() * realProduct
	   		if totalProduct != 1.0 || len(nonNumeric) == 0 {
	   			resultElements = append(resultElements, core.NewReal(totalProduct))
	   		}
	   	} else if hasIntegers {

	   		// All integers: keep as integer
	   		if (product.IsInt64() && product.Int64() != 1) || len(nonNumeric) == 0 {
	   			resultElements = append(resultElements, product)
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

	   return core.NewListFromExprs(resultElements...)
	*/
}
