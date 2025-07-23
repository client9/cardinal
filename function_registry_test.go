package sexpr

import (
	"testing"
)

func TestFunctionRegistry_Basic(t *testing.T) {
	ctx := NewContext()
	args := []Expr{NewIntAtom(1), NewIntAtom(2)}

	// Test that builtin functions are registered
	funcDef, _ := ctx.functionRegistry.FindMatchingFunction("Plus", args)
	if funcDef == nil {
		t.Error("Plus should be registered as builtin")
	}

	funcDef, _ = ctx.functionRegistry.FindMatchingFunction("Equal", args)
	if funcDef == nil {
		t.Error("Equal should be registered as builtin")
	}

	funcDef, _ = ctx.functionRegistry.FindMatchingFunction("NonExistentFunction", args)
	if funcDef != nil {
		t.Error("NonExistentFunction should not be registered")
	}
}

func TestFunctionRegistry_CustomFunction(t *testing.T) {
	ctx := NewContext()

	// Define a custom function using the new pattern system
	customDouble := func(args []Expr, ctx *Context) Expr {
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

	// Register the custom function with pattern
	err := ctx.GetFunctionRegistry().RegisterPatternBuiltins(map[string]PatternFunc{
		"Double(x_)": customDouble,
	})
	if err != nil {
		t.Fatalf("Failed to register Double function: %v", err)
	}

	// Test that it's registered by finding a match
	args := []Expr{NewIntAtom(21)}
	funcDef, bindings := ctx.functionRegistry.FindMatchingFunction("Double", args)
	if funcDef == nil {
		t.Error("Double should be registered")
	}

	// Test the function works
	result := funcDef.GoImpl(args, ctx)
	expected := "42"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}

	// Check bindings
	if len(bindings) != 1 || bindings["x"].String() != "21" {
		t.Errorf("expected binding x=21, got %v", bindings)
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

	// Register the custom function using pattern system
	customMaxPattern := func(args []Expr, ctx *Context) Expr {
		// Wrap the old function to work with new signature
		return customMax(args)
	}
	err := eval.context.GetFunctionRegistry().RegisterPatternBuiltins(map[string]PatternFunc{
		"Max(x___)": customMaxPattern,
	})
	if err != nil {
		t.Fatalf("Failed to register Max function: %v", err)
	}

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

	// Register custom function using pattern system
	customSquarePattern := func(args []Expr, ctx *Context) Expr {
		return customSquare(args)
	}
	err := parentCtx.GetFunctionRegistry().RegisterPatternBuiltins(map[string]PatternFunc{
		"Square(x_)": customSquarePattern,
	})
	if err != nil {
		t.Fatalf("Failed to register Square function: %v", err)
	}

	// Create child context
	childCtx := NewChildContext(parentCtx)

	// Child should inherit the custom function
	args := []Expr{NewIntAtom(7)}
	funcDef, _ := childCtx.functionRegistry.FindMatchingFunction("Square", args)
	if funcDef == nil {
		t.Error("Child context should inherit Square function")
	}

	// Child should also inherit standard builtins
	args2 := []Expr{NewIntAtom(1), NewIntAtom(2)}
	funcDef2, _ := childCtx.functionRegistry.FindMatchingFunction("Plus", args2)
	if funcDef2 == nil {
		t.Error("Child context should inherit Plus function")
	}

	// Test that the inherited function works
	result := funcDef.GoImpl(args, childCtx)
	expected := "49"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
}

func TestFunctionRegistry_ErrorPropagation(t *testing.T) {
	eval := NewEvaluator()

	// Define a function that uses other builtins through the evaluator
	customAverage := func(args []Expr, ctx *Context) Expr {
		if len(args) == 0 {
			return NewErrorExpr("ArgumentError", "Average expects at least 1 argument", args)
		}

		// Check for errors in arguments first
		for _, arg := range args {
			if IsError(arg) {
				return arg
			}
		}

		// Sum all arguments using Plus through evaluator
		plusList := NewList(append([]Expr{NewSymbolAtom("Plus")}, args...)...)
		sumResult := eval.evaluate(plusList, ctx)

		// If sum failed, propagate the error
		if IsError(sumResult) {
			return sumResult
		}

		// Divide by count (use float to ensure real division)
		count := NewFloatAtom(float64(len(args)))
		divideList := NewList(NewSymbolAtom("Divide"), sumResult, count)
		avgResult := eval.evaluate(divideList, ctx)

		return avgResult
	}

	// Register custom function using pattern system
	err := eval.context.GetFunctionRegistry().RegisterPatternBuiltins(map[string]PatternFunc{
		"Average(x___)": customAverage, // customAverage now has the right signature
	})
	if err != nil {
		t.Fatalf("Failed to register Average function: %v", err)
	}

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
	} else {
		errorExpr := result2.(*ErrorExpr)
		if errorExpr.ErrorType != "DivisionByZero" {
			t.Errorf("expected DivisionByZero error, got %s", errorExpr.ErrorType)
		}
	}
}
