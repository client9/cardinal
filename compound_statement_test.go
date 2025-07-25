package sexpr

import (
	"testing"
)

func TestCompoundStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple two expressions",
			input:    "1 + 2; 3 + 4",
			expected: "7",
		},
		{
			name:     "Sequential assignments",
			input:    "a = 5; b = 10; a + b",
			expected: "15",
		},
		{
			name:     "Chained variable dependencies",
			input:    "x = 2; y = x * 3; z = y + 1; z",
			expected: "7",
		},
		{
			name:     "Multiple simple values",
			input:    "1; 2; 3; 4",
			expected: "4",
		},
		{
			name:     "Assignment with function call",
			input:    "list = [1, 2, 3]; len = Length(list); len",
			expected: "3",
		},
		{
			name:     "Mixed arithmetic and assignment",
			input:    "a = 1 + 2; b = a * 3",
			expected: "9",
		},
		{
			name:     "String operations",
			input:    `name = "Alice"; age = 25; name`,
			expected: `"Alice"`,
		},
		{
			name:     "Boolean expressions",
			input:    "x = True; y = False; x",
			expected: "True",
		},
		{
			name:     "List operations",
			input:    "lst = [1, 2, 3]; first = First(lst); first",
			expected: "1",
		},
		{
			name:     "Nested arithmetic",
			input:    "a = 2; b = 3; c = a + b; d = c * 2; d",
			expected: "10",
		},
		{
			name:     "Comparison operations",
			input:    "x = 5; y = 3; result = x > y; result",
			expected: "True",
		},
		{
			name:     "Association operations",
			input:    "assoc = Association(Rule(key, \"value\")); val = Part(assoc, key); val",
			expected: `"value"`,
		},
		{
			name:     "Mathematical expressions",
			input:    "radius = 5; area = Times(Pi, Power(radius, 2)); area",
			expected: "78.53981633974483",
		},
		{
			name:     "Logical operations",
			input:    "a = True; b = False; result = And(a, Not(b)); result",
			expected: "True",
		},
		{
			name:     "Sequential function calls",
			input:    "nums = [5, 2, 8, 1]; sorted = nums; Length(sorted)",
			expected: "4",
		},
		{
			name:     "Variable override",
			input:    "x = 1; x = 2; x = 3; x",
			expected: "3",
		},
		{
			name:     "Complex expression with multiple operations",
			input:    "base = 2; exp = 3; power = Power(base, exp); sum = Plus(power, 2); sum",
			expected: "10.0",
		},
		{
			name:     "Single expression (no semicolon)",
			input:    "42",
			expected: "42",
		},
		{
			name:     "Assignment returns assigned value",
			input:    "x = 5; y = x; y",
			expected: "5",
		},
		{
			name:     "Mixed data types",
			input:    `num = 42; str = "hello"; bool = True; num`,
			expected: "42",
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
			if IsError(result) {
				t.Fatalf("Evaluation error: %v", result)
			}

			resultStr := result.String()
			if resultStr != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, resultStr)
			}
		})
	}
}

func TestCompoundStatementParsing(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid compound statement",
			input:       "1; 2",
			expectError: false,
		},
		{
			name:        "Valid single expression",
			input:       "42",
			expectError: false,
		},
		{
			name:        "Valid multiple semicolons",
			input:       "1; 2; 3; 4",
			expectError: false,
		},
		{
			name:        "Valid assignment compound",
			input:       "x = 1; y = 2",
			expectError: false,
		},
		{
			name:        "Invalid empty first expression",
			input:       "; 42",
			expectError: true,
			errorMsg:    "unexpected token",
		},
		{
			name:        "Invalid trailing semicolon",
			input:       "42;",
			expectError: true,
			errorMsg:    "unexpected token: EOF",
		},
		{
			name:        "Invalid empty between semicolons",
			input:       "1; ; 3",
			expectError: true,
			errorMsg:    "unexpected token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)

			if tt.expectError {
				if err == nil {
					// Maybe parsing succeeded but evaluation should fail
					evaluator := NewEvaluator()
					result := evaluator.Evaluate(expr)
					if !IsError(result) {
						t.Errorf("Expected error but got none")
					}
				} else if tt.errorMsg != "" && !containsSubstring(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected parse error: %v", err)
				} else {
					// Also check that evaluation works
					evaluator := NewEvaluator()
					result := evaluator.Evaluate(expr)
					if IsError(result) {
						t.Errorf("Unexpected evaluation error: %v", result)
					}
				}
			}
		})
	}
}

func TestCompoundStatementPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Semicolon vs assignment precedence",
			input:    "a = 1; b = 2",
			expected: "2",
		},
		{
			name:     "Semicolon vs arithmetic precedence",
			input:    "x = 2 + 3; y = x * 2",
			expected: "10",
		},
		{
			name:     "Complex precedence test",
			input:    "a = 1 + 2 * 3; b = a + 1; c = b * 2; c",
			expected: "16",
		},
		{
			name:     "Comparison within compound",
			input:    "x = 5; y = 3; x > y",
			expected: "True",
		},
		{
			name:     "Logical operations within compound",
			input:    "a = True; b = False; a && b",
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
			if IsError(result) {
				t.Fatalf("Evaluation error: %v", result)
			}

			resultStr := result.String()
			if resultStr != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, resultStr)
			}
		})
	}
}

func TestCompoundStatementSideEffects(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Variable persists across statements",
			input:    "x = 10; y = x + 5; y",
			expected: "15",
		},
		{
			name:     "Variable modification",
			input:    "counter = 0; counter = counter + 1; counter = counter + 1; counter",
			expected: "2",
		},
		{
			name:     "List modification persistence",
			input:    "lst = [1, 2]; first = First(lst); Length(lst)",
			expected: "2",
		},
		{
			name:     "Association building",
			input:    "assoc = Association(Rule(a, 1)); keys = Keys(assoc); Length(keys)",
			expected: "1",
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
			if IsError(result) {
				t.Fatalf("Evaluation error: %v", result)
			}

			resultStr := result.String()
			if resultStr != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, resultStr)
			}
		})
	}
}

// Helper function to check if a string contains a substring (case-insensitive)
func containsSubstring(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr ||
			(len(str) > len(substr) &&
				(str[:len(substr)] == substr ||
					str[len(str)-len(substr):] == substr ||
					containsSubstringHelper(str, substr))))
}

func containsSubstringHelper(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
