package sexpr

import (
	"strings"
	"testing"
)

func TestAssociationParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty association",
			input:    "{}",
			expected: "Association()",
		},
		{
			name:     "Single key-value pair",
			input:    "{name: \"Bob\"}",
			expected: "Association(Rule(name, \"Bob\"))",
		},
		{
			name:     "Multiple key-value pairs",
			input:    "{name: \"Bob\", age: 30}",
			expected: "Association(Rule(name, \"Bob\"), Rule(age, 30))",
		},
		{
			name:     "Mixed value types",
			input:    "{name: \"Bob\", age: 30, active: True}",
			expected: "Association(Rule(name, \"Bob\"), Rule(age, 30), Rule(active, True))",
		},
		{
			name:     "Trailing comma",
			input:    "{name: \"Bob\", age: 30,}",
			expected: "Association(Rule(name, \"Bob\"), Rule(age, 30))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			evaluator := NewEvaluator()
			result := evaluator.Evaluate(expr)

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestAssociationQ(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty association",
			input:    "AssociationQ({})",
			expected: "True",
		},
		{
			name:     "Non-empty association",
			input:    "AssociationQ({name: \"Bob\"})",
			expected: "True",
		},
		{
			name:     "List is not association",
			input:    "AssociationQ([1, 2, 3])",
			expected: "False",
		},
		{
			name:     "Integer is not association",
			input:    "AssociationQ(42)",
			expected: "False",
		},
		{
			name:     "String is not association",
			input:    "AssociationQ(\"test\")",
			expected: "False",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			evaluator := NewEvaluator()
			result := evaluator.Evaluate(expr)

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestAssociationKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty association",
			input:    "Keys({})",
			expected: "List()",
		},
		{
			name:     "Single key",
			input:    "Keys({name: \"Bob\"})",
			expected: "List(name)",
		},
		{
			name:     "Multiple keys",
			input:    "Keys({name: \"Bob\", age: 30})",
			expected: "List(name, age)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			evaluator := NewEvaluator()
			result := evaluator.Evaluate(expr)

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestAssociationValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty association",
			input:    "Values({})",
			expected: "List()",
		},
		{
			name:     "Single value",
			input:    "Values({name: \"Bob\"})",
			expected: "List(\"Bob\")",
		},
		{
			name:     "Multiple values",
			input:    "Values({name: \"Bob\", age: 30})",
			expected: "List(\"Bob\", 30)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			evaluator := NewEvaluator()
			result := evaluator.Evaluate(expr)

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestAssociationLength(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty association",
			input:    "Length({})",
			expected: "0",
		},
		{
			name:     "Single item",
			input:    "Length({name: \"Bob\"})",
			expected: "1",
		},
		{
			name:     "Multiple items",
			input:    "Length({name: \"Bob\", age: 30, active: True})",
			expected: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			evaluator := NewEvaluator()
			result := evaluator.Evaluate(expr)

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestAssociationPart(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Access existing key",
			input:    "Part({name: \"Bob\", age: 30}, name)",
			expected: "\"Bob\"",
		},
		{
			name:     "Access another existing key",
			input:    "Part({name: \"Bob\", age: 30}, age)",
			expected: "30",
		},
		{
			name:     "Access with string key",
			input:    "Part({\"key\": \"value\"}, \"key\")",
			expected: "\"value\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			evaluator := NewEvaluator()
			result := evaluator.Evaluate(expr)

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestAssociationPartErrors(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectsError bool
	}{
		{
			name:         "Access missing key",
			input:        "Part({name: \"Bob\"}, missing)",
			expectsError: true,
		},
		{
			name:         "Access missing string key",
			input:        "Part({\"key\": \"value\"}, \"missing\")",
			expectsError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			evaluator := NewEvaluator()
			result := evaluator.Evaluate(expr)

			isError := strings.HasPrefix(result.String(), "$Failed")
			if tt.expectsError && !isError {
				t.Errorf("Expected error, got %s", result.String())
			}
			if !tt.expectsError && isError {
				t.Errorf("Expected success, got error %s", result.String())
			}
		})
	}
}

func TestAssociationEquality(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Equal empty associations",
			input:    "SameQ({}, {})",
			expected: "True",
		},
		{
			name:     "Equal single-item associations",
			input:    "SameQ({name: \"Bob\"}, {name: \"Bob\"})",
			expected: "True",
		},
		{
			name:     "Different associations",
			input:    "SameQ({name: \"Bob\"}, {name: \"Alice\"})",
			expected: "False",
		},
		{
			name:     "Different key sets",
			input:    "SameQ({name: \"Bob\"}, {age: 30})",
			expected: "False",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			evaluator := NewEvaluator()
			result := evaluator.Evaluate(expr)

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestAssociationInsertionOrder(t *testing.T) {
	// Test that associations preserve insertion order
	expr, err := ParseString("{c: 3, a: 1, b: 2}")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	evaluator := NewEvaluator()
	evaluator.Evaluate(expr)

	// Keys should be in insertion order
	keysExpr, err := ParseString("Keys({c: 3, a: 1, b: 2})")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	keysResult := evaluator.Evaluate(keysExpr)
	expected := "List(c, a, b)"

	if keysResult.String() != expected {
		t.Errorf("Expected keys in insertion order %s, got %s", expected, keysResult.String())
	}
}

func TestAssociationPatternBehavior(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Keys with non-association",
			input:    "Keys(42)",
			expected: "Keys(42)", // Pattern doesn't match, returns unchanged
		},
		{
			name:     "Values with non-association",
			input:    "Values(\"test\")",
			expected: "Values(\"test\")", // Pattern doesn't match, returns unchanged
		},
		{
			name:     "Keys with multiple arguments",
			input:    "Keys({}, {})",
			expected: "Keys(Association(), Association())", // Pattern doesn't match, returns unchanged
		},
		{
			name:     "Values with multiple arguments",
			input:    "Values({}, {})",
			expected: "Values(Association(), Association())", // Pattern doesn't match, returns unchanged
		},
		{
			name:     "AssociationQ with multiple arguments",
			input:    "AssociationQ({}, {})",
			expected: "AssociationQ(Association(), Association())", // Pattern doesn't match, returns unchanged
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			evaluator := NewEvaluator()
			result := evaluator.Evaluate(expr)

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}
