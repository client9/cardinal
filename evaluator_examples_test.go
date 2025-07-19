package sexpr

import (
	"fmt"
	"testing"
)

// Example demonstrates basic evaluator usage
func Example_evaluator() {
	// Set up built-in attributes
	eval := setupTestEvaluator()
	
	// Parse and evaluate simple arithmetic
	expr, _ := ParseString("1 + 2 * 3")
	result := eval.Evaluate(expr)
	fmt.Println("1 + 2 * 3 =", result.String())
	
	// Parse and evaluate with variables
	expr, _ = ParseString("x = 5")
	eval.Evaluate(expr)
	
	expr, _ = ParseString("x + 3")
	result = eval.Evaluate(expr)
	fmt.Println("x + 3 =", result.String())
	
	// Parse and evaluate comparison
	expr, _ = ParseString("5 > 3")
	result = eval.Evaluate(expr)
	fmt.Println("5 > 3 =", result.String())
	
	// Output:
	// 1 + 2 * 3 = 7
	// x + 3 = 8
	// 5 > 3 = True
}

// Example_evaluatorAttributes demonstrates attribute-aware evaluation
func Example_evaluatorAttributes() {
	// Set up built-in attributes
	eval := setupTestEvaluator()
	
	// Demonstrate Flat attribute (associativity)
	expr, _ := ParseString("Plus(1, Plus(2, 3))")
	result := eval.Evaluate(expr)
	fmt.Println("Plus(1, Plus(2, 3)) =", result.String())
	
	// Demonstrate Orderless attribute (commutativity)
	expr, _ = ParseString("Plus(c, a, b)")
	result = eval.Evaluate(expr)
	fmt.Println("Plus(c, a, b) =", result.String())
	
	// Demonstrate Hold attribute
	expr, _ = ParseString("Hold(1 + 2)")
	result = eval.Evaluate(expr)
	fmt.Println("Hold(1 + 2) =", result.String())
	
	// Demonstrate OneIdentity attribute
	expr, _ = ParseString("Plus(42)")
	result = eval.Evaluate(expr)
	fmt.Println("Plus(42) =", result.String())
	
	// Output:
	// Plus(1, Plus(2, 3)) = 6
	// Plus(c, a, b) = Plus(a, b, c)
	// Hold(1 + 2) = Hold(Plus(1, 2))
	// Plus(42) = 42
}

// Example_evaluatorControlStructures demonstrates control structures
func Example_evaluatorControlStructures() {
	// Set up built-in attributes
	eval := setupTestEvaluator()
	
	// Conditional evaluation
	expr, _ := ParseString("If(True, 1 + 2, 3 * 4)")
	result := eval.Evaluate(expr)
	fmt.Println("If(True, 1 + 2, 3 * 4) =", result.String())
	
	expr, _ = ParseString("If(False, 1 + 2, 3 * 4)")
	result = eval.Evaluate(expr)
	fmt.Println("If(False, 1 + 2, 3 * 4) =", result.String())
	
	// Assignment and delayed assignment
	expr, _ = ParseString("x = 2 + 3")
	result = eval.Evaluate(expr)
	fmt.Println("x = 2 + 3 returns", result.String())
	
	expr, _ = ParseString("x")
	result = eval.Evaluate(expr)
	fmt.Println("x evaluates to", result.String())
	
	expr, _ = ParseString("y := 2 + 3")
	eval.Evaluate(expr)
	
	expr, _ = ParseString("y")
	result = eval.Evaluate(expr)
	fmt.Println("y evaluates to", result.String())
	
	// Output:
	// If(True, 1 + 2, 3 * 4) = 3
	// If(False, 1 + 2, 3 * 4) = 12
	// x = 2 + 3 returns 5
	// x evaluates to 5
	// y evaluates to Plus(2, 3)
}

// Example_evaluatorConstants demonstrates built-in constants
func Example_evaluatorConstants() {
	// Set up built-in attributes
	eval := setupTestEvaluator()
	
	// Mathematical constants
	expr, _ := ParseString("Pi")
	result := eval.Evaluate(expr)
	fmt.Printf("Pi = %.6f\n", result.(Atom).Value.(float64))
	
	expr, _ = ParseString("E")
	result = eval.Evaluate(expr)
	fmt.Printf("E = %.6f\n", result.(Atom).Value.(float64))
	
	// Boolean constants
	expr, _ = ParseString("True")
	result = eval.Evaluate(expr)
	fmt.Println("True =", result.String())
	
	expr, _ = ParseString("False")
	result = eval.Evaluate(expr)
	fmt.Println("False =", result.String())
	
	// Output:
	// Pi = 3.141593
	// E = 2.718282
	// True = True
	// False = False
}

// TestEvaluatorIntegration demonstrates full integration with parsing and attributes
func TestEvaluatorIntegration(t *testing.T) {
	// Set up built-in attributes
	eval := setupTestEvaluator()
	
	// Test complex mathematical expression
	expr, err := ParseString("(1 + 2) * (3 + 4)")
	if err != nil {
		t.Fatal(err)
	}
	
	result := eval.Evaluate(expr)
	expected := "21"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
	
	// Test variable assignment and usage
	expr, _ = ParseString("x = 10")
	eval.Evaluate(expr)
	
	expr, _ = ParseString("y = 20")
	eval.Evaluate(expr)
	
	expr, _ = ParseString("x + y")
	result = eval.Evaluate(expr)
	expected = "30"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
	
	// Test complex conditional
	expr, _ = ParseString("If(x > y, x - y, y - x)")
	result = eval.Evaluate(expr)
	expected = "10"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
	
	// Test attribute transformations
	expr, _ = ParseString("Plus(1, Plus(2, Plus(3, 4)))")
	result = eval.Evaluate(expr)
	expected = "10"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
	
	// Test logical operations
	expr, _ = ParseString("And(True, Or(False, True))")
	result = eval.Evaluate(expr)
	expected = "True"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
	
	// Test comparison operations
	expr, _ = ParseString("And(3 > 2, 5 < 10)")
	result = eval.Evaluate(expr)
	expected = "True"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
	
	// Test SameQ/UnsameQ
	expr, _ = ParseString("SameQ(3, 3)")
	result = eval.Evaluate(expr)
	expected = "True"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
	
	expr, _ = ParseString("UnsameQ(3, 4)")
	result = eval.Evaluate(expr)
	expected = "True"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
}

// TestEvaluatorErrorHandling tests error handling and edge cases
func TestEvaluatorErrorHandling(t *testing.T) {
	// Set up built-in attributes
	eval := setupTestEvaluator()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "division by zero",
			input:    "1 / 0",
			expected: "$Failed(DivisionByZero)",
		},
		{
			name:     "undefined variable",
			input:    "undefinedVar",
			expected: "undefinedVar",
		},
		{
			name:     "invalid function",
			input:    "UnknownFunction(1, 2)",
			expected: "UnknownFunction(1, 2)",
		},
		{
			name:     "mixed numeric and symbolic",
			input:    "Plus(1, x, 3)",
			expected: "Plus(1, 3, x)",
		},
		{
			name:     "comparison with symbols",
			input:    "Less(x, y)",
			expected: "Less(x, y)",
		},
		{
			name:     "logical with non-boolean",
			input:    "And(True, x)",
			expected: "And(True, x)",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			
			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

// BenchmarkEvaluator benchmarks evaluator performance
func BenchmarkEvaluator(b *testing.B) {
	// Set up built-in attributes
	eval := setupTestEvaluator()
	
	// Parse expressions once
	expr1, _ := ParseString("1 + 2 * 3")
	expr2, _ := ParseString("Plus(1, 2, 3, 4, 5)")
	expr3, _ := ParseString("If(True, 1 + 2, 3 * 4)")
	
	b.Run("SimpleArithmetic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			eval.Evaluate(expr1)
		}
	})
	
	b.Run("MultipleArguments", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			eval.Evaluate(expr2)
		}
	})
	
	b.Run("ConditionalEvaluation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			eval.Evaluate(expr3)
		}
	})
	
	b.Run("WithAttributeTransformation", func(b *testing.B) {
		expr, _ := ParseString("Plus(1, Plus(2, Plus(3, 4)))")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			eval.Evaluate(expr)
		}
	})
}
