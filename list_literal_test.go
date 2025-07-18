package sexpr

import (
	"testing"
)

func TestParseListLiteral_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty list",
			input:    "{}",
			expected: "List[]",
		},
		{
			name:     "Single element",
			input:    "{1}",
			expected: "List[1]",
		},
		{
			name:     "Multiple integers",
			input:    "{1, 2, 3}",
			expected: "List[1, 2, 3]",
		},
		{
			name:     "Mixed types",
			input:    "{1, \"hello\", True}",
			expected: "List[1, \"hello\", True]",
		},
		{
			name:     "Floating point numbers",
			input:    "{3.14, 2.71}",
			expected: "List[3.14, 2.71]",
		},
		{
			name:     "Symbols",
			input:    "{x, y, z}",
			expected: "List[x, y, z]",
		},
		{
			name:     "Nested expressions",
			input:    "{Plus[1, 2], Times[3, 4]}",
			expected: "List[Plus[1, 2], Times[3, 4]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			
			if expr.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, expr.String())
			}
		})
	}
}

func TestParseListLiteral_TrailingComma(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single element with trailing comma",
			input:    "{1,}",
			expected: "List[1]",
		},
		{
			name:     "Multiple elements with trailing comma",
			input:    "{1, 2, 3,}",
			expected: "List[1, 2, 3]",
		},
		{
			name:     "Mixed types with trailing comma",
			input:    "{\"hello\", True, 42,}",
			expected: "List[\"hello\", True, 42]",
		},
		{
			name:     "Nested expressions with trailing comma",
			input:    "{Plus[1, 2], x,}",
			expected: "List[Plus[1, 2], x]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			
			if expr.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, expr.String())
			}
		})
	}
}

func TestParseListLiteral_NestedLists(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "List containing empty list",
			input:    "{{}}",
			expected: "List[List[]]",
		},
		{
			name:     "List containing multiple lists",
			input:    "{{1, 2}, {3, 4}}",
			expected: "List[List[1, 2], List[3, 4]]",
		},
		{
			name:     "Deeply nested lists",
			input:    "{{{1}}}",
			expected: "List[List[List[1]]]",
		},
		{
			name:     "Mixed nested structures",
			input:    "{1, {2, 3}, 4}",
			expected: "List[1, List[2, 3], 4]",
		},
		{
			name:     "Complex nested with trailing commas",
			input:    "{{1, 2,}, {3,}, {}}",
			expected: "List[List[1, 2], List[3], List[]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			
			if expr.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, expr.String())
			}
		})
	}
}

func TestParseListLiteral_WithArithmetic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "List with arithmetic expressions",
			input:    "{1 + 2, 3 * 4}",
			expected: "List[Plus[1, 2], Times[3, 4]]",
		},
		{
			name:     "List with complex expressions",
			input:    "{x + y, z * 2, True && False}",
			expected: "List[Plus[x, y], Times[z, 2], And[True, False]]",
		},
		{
			name:     "List with comparisons",
			input:    "{1 == 2, 3 > 4, x <= y}",
			expected: "List[Equal[1, 2], Greater[3, 4], LessEqual[x, y]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			
			if expr.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, expr.String())
			}
		})
	}
}

func TestParseListLiteral_ErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Unclosed list",
			input: "{1, 2, 3",
		},
		{
			name:  "Missing comma",
			input: "{1 2 3}",
		},
		{
			name:  "Extra comma at start",
			input: "{, 1, 2}",
		},
		{
			name:  "Double comma",
			input: "{1,, 2}",
		},
		{
			name:  "Only comma",
			input: "{,}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseString(tt.input)
			if err == nil {
				t.Errorf("expected parse error for input: %s", tt.input)
			}
		})
	}
}

// Integration tests with the evaluator
func TestListLiteral_Integration(t *testing.T) {
	eval := setupTestEvaluator()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty list evaluation",
			input:    "{}",
			expected: "List[]",
		},
		{
			name:     "Simple list evaluation",
			input:    "{1, 2, 3}",
			expected: "List[1, 2, 3]",
		},
		{
			name:     "List with expressions",
			input:    "{1 + 2, 3 * 4}",
			expected: "List[3, 12]", // Expressions are evaluated
		},
		{
			name:     "Length of list literal",
			input:    "Length[{1, 2, 3, 4}]",
			expected: "4",
		},
		{
			name:     "ListQ on list literal",
			input:    "ListQ[{1, 2, 3}]",
			expected: "True",
		},
		{
			name:     "Head of list literal",
			input:    "Head[{1, 2, 3}]",
			expected: "List",
		},
		{
			name:     "Nested list evaluation",
			input:    "{{1, 2}, {3, 4}}",
			expected: "List[List[1, 2], List[3, 4]]",
		},
		{
			name:     "List with mixed types",
			input:    "{42, \"hello\", True, x}",
			expected: "List[42, \"hello\", True, x]",
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

// Test that list literals work correctly with Hold
func TestListLiteral_WithHold(t *testing.T) {
	eval := setupTestEvaluator()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Hold list literal",
			input:    "Hold[{1, 2, 3}]",
			expected: "Hold[List[1, 2, 3]]",
		},
		{
			name:     "Hold list with expressions",
			input:    "Hold[{1 + 2, 3 * 4}]",
			expected: "Hold[List[Plus[1, 2], Times[3, 4]]]",
		},
		{
			name:     "List containing Hold",
			input:    "{Hold[1 + 2], 3}",
			expected: "List[Hold[Plus[1, 2]], 3]",
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

// Test whitespace handling in list literals
func TestListLiteral_Whitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Extra spaces",
			input:    "{ 1 , 2 , 3 }",
			expected: "List[1, 2, 3]",
		},
		{
			name:     "Newlines in list",
			input:    "{\n1,\n2,\n3\n}",
			expected: "List[1, 2, 3]",
		},
		{
			name:     "Tabs and spaces",
			input:    "{\t1,\t\t2,   3\t}",
			expected: "List[1, 2, 3]",
		},
		{
			name:     "No spaces",
			input:    "{1,2,3}",
			expected: "List[1, 2, 3]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			
			if expr.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, expr.String())
			}
		})
	}
}