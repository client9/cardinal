package sexpr

import (
	"testing"
)

func TestFunctionRegistry_Basic(t *testing.T) {
	ctx := NewContext()

	// Test that builtin functions are registered
	if !ctx.HasBuiltin("Plus") {
		t.Error("Plus should be registered as builtin")
	}

	if !ctx.HasBuiltin("Equal") {
		t.Error("Equal should be registered as builtin")
	}

	if ctx.HasBuiltin("NonExistentFunction") {
		t.Error("NonExistentFunction should not be registered")
	}
}

func TestFunctionRegistry_CustomFunction(t *testing.T) {
	ctx := NewContext()

	// Define a custom function
	customDouble := func(args []Expr) Expr {
		if len(args) != 1 {
			return NewErrorExpr("ArgumentError", "Double expects 1 argument", args)
		}

		if isNumeric(args[0]) {
			val, _ := getNumericValue(args[0])
			return createNumericResult(val * 2)
		}

		// Return symbolic form if not numeric
		return List{Elements: []Expr{NewSymbolAtom("Double"), args[0]}}
	}

	// Register the custom function
	ctx.RegisterBuiltin("Double", customDouble)

	// Test that it's registered
	if !ctx.HasBuiltin("Double") {
		t.Error("Double should be registered")
	}

	// Test that we can retrieve it
	fn, exists := ctx.GetBuiltin("Double")
	if !exists {
		t.Error("Should be able to retrieve Double function")
	}

	// Test the function works
	result := fn([]Expr{NewIntAtom(21)})
	expected := "42"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
}

func TestFunctionRegistry_Evaluator(t *testing.T) {
	eval := NewEvaluator()

	// Define a custom Max function
	customMax := func(args []Expr) Expr {
		if len(args) == 0 {
			return NewErrorExpr("ArgumentError", "Max expects at least 1 argument", args)
		}

		// Check if all arguments are numeric
		for _, arg := range args {
			if !isNumeric(arg) {
				// Return symbolic form if any argument is not numeric
				elements := make([]Expr, len(args)+1)
				elements[0] = NewSymbolAtom("Max")
				copy(elements[1:], args)
				return List{Elements: elements}
			}
		}

		// Find the maximum value
		maxVal, _ := getNumericValue(args[0])
		for _, arg := range args[1:] {
			val, _ := getNumericValue(arg)
			if val > maxVal {
				maxVal = val
			}
		}

		return createNumericResult(maxVal)
	}

	// Register the custom function
	eval.context.RegisterBuiltin("Max", customMax)

	// Test evaluation with the custom function
	expr, err := ParseString("Max(1, 5, 3, 9, 2)")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := eval.Evaluate(expr)
	expected := "9"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
}

func TestFunctionRegistry_ChildContextInheritance(t *testing.T) {
	parentCtx := NewContext()

	// Add a custom function to the parent
	customSquare := func(args []Expr) Expr {
		if len(args) != 1 {
			return NewErrorExpr("ArgumentError", "Square expects 1 argument", args)
		}

		if isNumeric(args[0]) {
			val, _ := getNumericValue(args[0])
			return createNumericResult(val * val)
		}

		return List{Elements: []Expr{NewSymbolAtom("Square"), args[0]}}
	}

	parentCtx.RegisterBuiltin("Square", customSquare)

	// Create child context
	childCtx := NewChildContext(parentCtx)

	// Child should inherit the custom function
	if !childCtx.HasBuiltin("Square") {
		t.Error("Child context should inherit Square function")
	}

	// Child should also inherit standard builtins
	if !childCtx.HasBuiltin("Plus") {
		t.Error("Child context should inherit Plus function")
	}

	// Test that the inherited function works
	fn, exists := childCtx.GetBuiltin("Square")
	if !exists {
		t.Error("Should be able to retrieve Square function from child")
	}

	result := fn([]Expr{NewIntAtom(7)})
	expected := "49"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
}

func TestFunctionRegistry_ErrorPropagation(t *testing.T) {
	eval := NewEvaluator()

	// Define a function that uses other builtins
	customAverage := func(args []Expr) Expr {
		if len(args) == 0 {
			return NewErrorExpr("ArgumentError", "Average expects at least 1 argument", args)
		}

		// Sum all arguments using Plus
		sumResult := EvaluatePlus(args)

		// If sum failed, propagate the error
		if IsError(sumResult) {
			return sumResult
		}

		// Divide by count
		count := NewIntAtom(len(args))
		avgResult := EvaluateDivide([]Expr{sumResult, count})

		return avgResult
	}

	eval.context.RegisterBuiltin("Average", customAverage)

	// Test with valid arguments
	expr, err := ParseString("Average(1, 2, 3, 4)")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := eval.Evaluate(expr)
	expected := "2.5"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}

	// Test with an error case (average includes division by zero)
	expr2, err2 := ParseString("Average(10, Divide(1, 0), 30)")
	if err2 != nil {
		t.Fatalf("Parse error: %v", err2)
	}

	result2 := eval.Evaluate(expr2)
	if !IsError(result2) {
		t.Errorf("expected error propagation, got %s", result2.String())
	}

	errorExpr := result2.(*ErrorExpr)
	if errorExpr.ErrorType != "DivisionByZero" {
		t.Errorf("expected DivisionByZero error, got %s", errorExpr.ErrorType)
	}
}

func TestFunctionRegistry_OverrideBuiltin(t *testing.T) {
	ctx := NewContext()

	// Define a custom Plus that only adds the first two arguments
	customPlus := func(args []Expr) Expr {
		if len(args) < 2 {
			return NewErrorExpr("ArgumentError", "CustomPlus expects at least 2 arguments", args)
		}

		// Only add first two arguments
		return EvaluatePlus(args[:2])
	}

	// Override the builtin Plus
	ctx.RegisterBuiltin("Plus", customPlus)

	eval := NewEvaluatorWithContext(ctx)

	// Test that our custom Plus is used
	expr, err := ParseString("Plus(1, 2, 3, 4)")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := eval.Evaluate(expr)
	expected := "3" // Only 1 + 2, ignoring 3 and 4
	if result.String() != expected {
		t.Errorf("expected %s (custom Plus), got %s", expected, result.String())
	}
}
