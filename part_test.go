package sexpr

import (
	"fmt"
	"testing"
)

func TestEvaluatePart(t *testing.T) {
	tests := []struct {
		name      string
		expr      Expr
		index     int
		expected  string
		hasError  bool
		errorType string
	}{
		{
			name: "Part 1 of simple list",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewIntAtom(10),
				NewIntAtom(20),
				NewIntAtom(30),
			}},
			index:    1,
			expected: "10",
			hasError: false,
		},
		{
			name: "Part 2 of simple list",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewIntAtom(10),
				NewIntAtom(20),
				NewIntAtom(30),
			}},
			index:    2,
			expected: "20",
			hasError: false,
		},
		{
			name: "Part 3 of simple list",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewIntAtom(10),
				NewIntAtom(20),
				NewIntAtom(30),
			}},
			index:    3,
			expected: "30",
			hasError: false,
		},
		{
			name: "Part 1 of function expression",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewSymbolAtom("a"),
				NewSymbolAtom("b"),
				NewSymbolAtom("c"),
			}},
			index:    1,
			expected: "a",
			hasError: false,
		},
		{
			name: "Part 2 of function expression",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewSymbolAtom("a"),
				NewSymbolAtom("b"),
				NewSymbolAtom("c"),
			}},
			index:    2,
			expected: "b",
			hasError: false,
		},
		{
			name: "Part with mixed types",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("Mixed"),
				NewIntAtom(42),
				NewStringAtom("hello"),
				NewBoolAtom(true),
				NewFloatAtom(3.14),
			}},
			index:    3,
			expected: "True",
			hasError: false,
		},
		{
			name: "Part with nested expression",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("Outer"),
				&List{Elements: []Expr{
					NewSymbolAtom("Inner"),
					NewIntAtom(1),
					NewIntAtom(2),
				}},
				NewIntAtom(3),
			}},
			index:    1,
			expected: "Inner(1, 2)",
			hasError: false,
		},
		{
			name: "Part -1 (last element)",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewIntAtom(10),
				NewIntAtom(20),
				NewIntAtom(30),
			}},
			index:    -1,
			expected: "30",
			hasError: false,
		},
		{
			name: "Part -2 (second to last)",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewIntAtom(10),
				NewIntAtom(20),
				NewIntAtom(30),
			}},
			index:    -2,
			expected: "20",
			hasError: false,
		},
		{
			name: "Part -3 (third to last)",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewIntAtom(10),
				NewIntAtom(20),
				NewIntAtom(30),
			}},
			index:    -3,
			expected: "10",
			hasError: false,
		},
		{
			name: "Part of single element list",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewStringAtom("only"),
			}},
			index:    1,
			expected: "\"only\"",
			hasError: false,
		},
		// Error cases
		{
			name: "Part index 0 - should error",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewIntAtom(10),
				NewIntAtom(20),
			}},
			index:     0,
			expected:  "",
			hasError:  true,
			errorType: "PartError",
		},
		{
			name: "Part index out of bounds (positive)",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewIntAtom(10),
				NewIntAtom(20),
			}},
			index:     5,
			expected:  "",
			hasError:  true,
			errorType: "PartError",
		},
		{
			name: "Part index out of bounds (negative)",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewIntAtom(10),
				NewIntAtom(20),
			}},
			index:     -5,
			expected:  "",
			hasError:  true,
			errorType: "PartError",
		},
		{
			name: "Part of empty list - should error",
			expr: &List{Elements: []Expr{}},
			index:     1,
			expected:  "",
			hasError:  true,
			errorType: "PartError",
		},
		{
			name: "Part of single element list (head only) - should error",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("OnlyHead"),
			}},
			index:     1,
			expected:  "",
			hasError:  true,
			errorType: "PartError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluatePart([]Expr{tt.expr, NewIntAtom(tt.index)})
			
			if tt.hasError {
				if !IsError(result) {
					t.Errorf("expected error for %s, got %s", tt.name, result.String())
					return
				}
				
				errorExpr := result.(*ErrorExpr)
				if errorExpr.ErrorType != tt.errorType {
					t.Errorf("expected error type %s, got %s", tt.errorType, errorExpr.ErrorType)
				}
			} else {
				if IsError(result) {
					t.Errorf("unexpected error: %s", result.String())
					return
				}
				
				if result.String() != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.String())
				}
			}
		})
	}
}

func TestEvaluatePart_NonIntegerIndex(t *testing.T) {
	testList := &List{Elements: []Expr{
		NewSymbolAtom("List"),
		NewIntAtom(1),
		NewIntAtom(2),
	}}

	tests := []struct {
		name     string
		index    Expr
		expected string
	}{
		{
			name:     "String index",
			index:    NewStringAtom("hello"),
			expected: "$Failed(PartError)",
		},
		{
			name:     "Float index",
			index:    NewFloatAtom(1.5),
			expected: "$Failed(PartError)",
		},
		{
			name:     "Boolean index",
			index:    NewBoolAtom(true),
			expected: "$Failed(PartError)",
		},
		{
			name:     "Symbol index",
			index:    NewSymbolAtom("x"),
			expected: "$Failed(PartError)",
		},
		{
			name: "List index",
			index: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewIntAtom(1),
			}},
			expected: "$Failed(PartError)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluatePart([]Expr{testList, tt.index})
			
			if !IsError(result) {
				t.Errorf("expected error for non-integer index, got %s", result.String())
				return
			}
			
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestEvaluatePart_NonListExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected string
	}{
		{
			name:     "Integer atom",
			expr:     NewIntAtom(42),
			expected: "$Failed(PartError)",
		},
		{
			name:     "String atom",
			expr:     NewStringAtom("hello"),
			expected: "$Failed(PartError)",
		},
		{
			name:     "Boolean atom",
			expr:     NewBoolAtom(true),
			expected: "$Failed(PartError)",
		},
		{
			name:     "Symbol atom",
			expr:     NewSymbolAtom("x"),
			expected: "$Failed(PartError)",
		},
		{
			name:     "Float atom",
			expr:     NewFloatAtom(3.14),
			expected: "$Failed(PartError)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluatePart([]Expr{tt.expr, NewIntAtom(1)})
			
			if !IsError(result) {
				t.Errorf("expected error for non-list expression, got %s", result.String())
				return
			}
			
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

// Test argument validation for Part
func TestPart_ArgumentValidation(t *testing.T) {
	testList := &List{Elements: []Expr{
		NewSymbolAtom("List"),
		NewIntAtom(1),
		NewIntAtom(2),
	}}

	t.Run("Part_no_args", func(t *testing.T) {
		result := EvaluatePart([]Expr{})
		if !IsError(result) {
			t.Errorf("expected error for no arguments, got %s", result.String())
		}
		
		errorExpr := result.(*ErrorExpr)
		if errorExpr.ErrorType != "ArgumentError" {
			t.Errorf("expected ArgumentError, got %s", errorExpr.ErrorType)
		}
	})
	
	t.Run("Part_one_arg", func(t *testing.T) {
		result := EvaluatePart([]Expr{testList})
		if !IsError(result) {
			t.Errorf("expected error for one argument, got %s", result.String())
		}
		
		errorExpr := result.(*ErrorExpr)
		if errorExpr.ErrorType != "ArgumentError" {
			t.Errorf("expected ArgumentError, got %s", errorExpr.ErrorType)
		}
	})
	
	t.Run("Part_too_many_args", func(t *testing.T) {
		result := EvaluatePart([]Expr{testList, NewIntAtom(1), NewIntAtom(2)})
		if !IsError(result) {
			t.Errorf("expected error for too many arguments, got %s", result.String())
		}
		
		errorExpr := result.(*ErrorExpr)
		if errorExpr.ErrorType != "ArgumentError" {
			t.Errorf("expected ArgumentError, got %s", errorExpr.ErrorType)
		}
	})
}

// Integration tests with the evaluator
func TestPart_Integration(t *testing.T) {
	eval := setupTestEvaluator()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic positive indexing
		{
			name:     "Part 1 of simple list",
			input:    "Part([a, b, c], 1)",
			expected: "a",
		},
		{
			name:     "Part 2 of simple list",
			input:    "Part([a, b, c], 2)",
			expected: "b",
		},
		{
			name:     "Part 3 of simple list",
			input:    "Part([a, b, c], 3)",
			expected: "c",
		},
		{
			name:     "Part of function call",
			input:    "Part(Plus(x, y, z), 2)",
			expected: "y",
		},
		{
			name:     "Part of numeric list",
			input:    "Part([10, 20, 30, 40], 3)",
			expected: "30",
		},
		
		// Negative indexing
		{
			name:     "Part -1 (last element)",
			input:    "Part([a, b, c], -1)",
			expected: "c",
		},
		{
			name:     "Part -2 (second to last)",
			input:    "Part([a, b, c], -2)",
			expected: "b",
		},
		{
			name:     "Part -3 (third to last)",
			input:    "Part([a, b, c], -3)",
			expected: "a",
		},
		
		// Single element lists
		{
			name:     "Part of single element list",
			input:    "Part([x], 1)",
			expected: "x",
		},
		{
			name:     "Part of single element with negative index",
			input:    "Part([x], -1)",
			expected: "x",
		},
		
		// Error cases
		{
			name:     "Part index 0 - should error",
			input:    "Part([a, b, c], 0)",
			expected: "$Failed(PartError)",
		},
		{
			name:     "Part index out of bounds",
			input:    "Part([a, b], 5)",
			expected: "$Failed(PartError)",
		},
		{
			name:     "Part negative index out of bounds",
			input:    "Part([a, b], -5)",
			expected: "$Failed(PartError)",
		},
		{
			name:     "Part of atom",
			input:    "Part(42, 1)",
			expected: "$Failed(PartError)",
		},
		{
			name:     "Part of empty list",
			input:    "Part([], 1)",
			expected: "$Failed(PartError)",
		},
		{
			name:     "Part with non-integer index",
			input:    "Part([a, b, c], \"hello\")",
			expected: "$Failed(PartError)",
		},
		
		// Combination with other functions
		{
			name:     "Part of Rest result",
			input:    "Part(Rest([a, b, c, d]), 2)",
			expected: "c",
		},
		{
			name:     "Part of Most result",
			input:    "Part(Most([a, b, c, d]), 1)",
			expected: "a",
		},
		{
			name:     "First equals Part[expr, 1]",
			input:    "Equal(First([x, y, z]), Part([x, y, z], 1))",
			expected: "True",
		},
		{
			name:     "Last equals Part[expr, -1]",
			input:    "Equal(Last([x, y, z]), Part([x, y, z], -1))",
			expected: "True",
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
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

// Test Part with complex expressions
func TestPart_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		index    int
		expected string
	}{
		{
			name: "Part of deeply nested expression",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("Outer"),
				&List{Elements: []Expr{
					NewSymbolAtom("Level1"),
					&List{Elements: []Expr{
						NewSymbolAtom("Level2"),
						NewIntAtom(1),
						NewIntAtom(2),
					}},
					NewIntAtom(3),
				}},
				NewIntAtom(4),
				NewIntAtom(5),
			}},
			index:    2,
			expected: "4",
		},
		{
			name: "Part of expression with string elements",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("StringList"),
				NewStringAtom("first"),
				NewStringAtom("second"),
				NewStringAtom("third"),
			}},
			index:    2,
			expected: "\"second\"",
		},
		{
			name: "Part of expression with mixed atom types",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("Mixed"),
				NewIntAtom(42),
				NewFloatAtom(3.14159),
				NewBoolAtom(false),
				NewStringAtom("test"),
				NewSymbolAtom("symbol"),
			}},
			index:    4,
			expected: "\"test\"",
		},
		{
			name: "Part with negative index on complex expression",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("Complex"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
				NewIntAtom(4),
				NewIntAtom(5),
			}},
			index:    -2,
			expected: "4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluatePart([]Expr{tt.expr, NewIntAtom(tt.index)})
			
			if IsError(result) {
				t.Errorf("unexpected error: %s", result.String())
				return
			}
			
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

// Test boundary conditions
func TestPart_BoundaryConditions(t *testing.T) {
	// Test various list sizes
	sizes := []int{1, 2, 3, 5, 10}
	
	for _, size := range sizes {
		t.Run(fmt.Sprintf("List_size_%d", size), func(t *testing.T) {
			// Create a list of the specified size
			elements := make([]Expr, size+1) // +1 for head
			elements[0] = NewSymbolAtom("List")
			for i := 1; i <= size; i++ {
				elements[i] = NewIntAtom(i * 10)
			}
			list := &List{Elements: elements}
			
			// Test first element
			result := EvaluatePart([]Expr{list, NewIntAtom(1)})
			if IsError(result) {
				t.Errorf("unexpected error accessing first element: %s", result.String())
			} else if result.String() != "10" {
				t.Errorf("expected 10, got %s", result.String())
			}
			
			// Test last element (positive index)
			result = EvaluatePart([]Expr{list, NewIntAtom(size)})
			if IsError(result) {
				t.Errorf("unexpected error accessing last element: %s", result.String())
			} else {
				expected := fmt.Sprintf("%d", size*10)
				if result.String() != expected {
					t.Errorf("expected %s, got %s", expected, result.String())
				}
			}
			
			// Test last element (negative index)
			result = EvaluatePart([]Expr{list, NewIntAtom(-1)})
			if IsError(result) {
				t.Errorf("unexpected error accessing last element with -1: %s", result.String())
			} else {
				expected := fmt.Sprintf("%d", size*10)
				if result.String() != expected {
					t.Errorf("expected %s, got %s", expected, result.String())
				}
			}
			
			// Test out of bounds (positive)
			result = EvaluatePart([]Expr{list, NewIntAtom(size + 1)})
			if !IsError(result) {
				t.Errorf("expected error for out of bounds access, got %s", result.String())
			}
			
			// Test out of bounds (negative)
			result = EvaluatePart([]Expr{list, NewIntAtom(-(size + 1))})
			if !IsError(result) {
				t.Errorf("expected error for out of bounds negative access, got %s", result.String())
			}
		})
	}
}
