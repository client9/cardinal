package sexpr

import (
	"strings"
	"testing"

	"github.com/client9/sexpr/core"
)

func TestMap_BasicFunctionality(t *testing.T) {
	evaluator := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Map with Plus function on single elements",
			input:    `Map(Plus($1, 10) &, [1, 2, 3])`,
			expected: `List(11, 12, 13)`,
		},
		{
			name:     "Map with Times function on single elements",
			input:    `Map(Times($1, 2) &, [2, 3, 4])`,
			expected: `List(4, 6, 8)`,
		},
		{
			name:     "Map with Length function",
			input:    `Map(Length, [[1, 2, 3], [4, 5], [6]])`,
			expected: `List(3, 2, 1)`,
		},
		{
			name:     "Map with Head function",
			input:    `Map(Head, [Plus(a, b), Times(x, y), List(1, 2)])`,
			expected: `List(Plus, Times, List)`,
		},
		{
			name:     "Map with empty list",
			input:    `Map(Plus, [])`,
			expected: `List()`,
		},
		{
			name:     "Map with single element",
			input:    `Map(Length, [[1, 2, 3]])`,
			expected: `List(3)`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result.String())
			}
		})
	}
}

func TestMap_WithAmpersandSyntax(t *testing.T) {
	evaluator := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Map with & syntax - double elements",
			input:    `Map(Plus($1, $1) &, [1, 2, 3, 4])`,
			expected: `List(2, 4, 6, 8)`,
		},
		{
			name:     "Map with & syntax - square elements",
			input:    `Map(Times($1, $1) &, [2, 3, 4])`,
			expected: `List(4, 9, 16)`,
		},
		{
			name:     "Map with & syntax - increment",
			input:    `Map(Plus($1, 1) &, [10, 20, 30])`,
			expected: `List(11, 21, 31)`,
		},
		{
			name:     "Map with & syntax - complex expression",
			input:    `Map(Plus(Times($1, 2), 1) &, [1, 2, 3])`,
			expected: `List(3, 5, 7)`,
		},
		{
			name:     "Map with & syntax - string operations",
			input:    `Map(StringLength($1) &, ["hello", "world", "test"])`,
			expected: `List(5, 5, 4)`,
		},
		{
			name:     "Map with & syntax - nested lists",
			input:    `Map(Length($1) &, [[1, 2], [3, 4, 5], [6]])`,
			expected: `List(2, 3, 1)`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result.String())
			}
		})
	}
}

func TestMap_EdgeCases(t *testing.T) {
	evaluator := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Map with non-list should error",
			input:    `Map(Plus, 42)`,
			expected: `ErrorExpr(ArgumentError)`,
		},
		{
			name:     "Map with wrong number of arguments returns unevaluated",
			input:    `Map(Plus)`,
			expected: `Map(Plus)`,
		},
		{
			name:     "Map with too many arguments returns unevaluated",
			input:    `Map(Plus, [1, 2], [3, 4])`,
			expected: `Map(Plus, List(1, 2), List(3, 4))`,
		},
		{
			name:     "Map preserves head of input list",
			input:    `Map(Plus($1, 1) &, MyList(1, 2, 3))`,
			expected: `MyList(2, 3, 4)`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if core.IsError(result) {
				if !strings.Contains(test.expected, "ErrorExpr") {
					t.Errorf("Expected %s, got error: %s", test.expected, result.String())
				}
			} else if result.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result.String())
			}
		})
	}
}

func TestApply_BasicFunctionality(t *testing.T) {
	evaluator := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Apply with Plus function",
			input:    `Apply(Plus, [1, 2, 3, 4])`,
			expected: `10`,
		},
		{
			name:     "Apply with Times function",
			input:    `Apply(Times, [2, 3, 4])`,
			expected: `24`,
		},
		{
			name:     "Apply with List function",
			input:    `Apply(List, [1, 2, 3])`,
			expected: `List(1, 2, 3)`,
		},
		{
			name:     "Apply with single argument",
			input:    `Apply(Length, [[1, 2, 3, 4]])`,
			expected: `4`,
		},
		{
			name:     "Apply with empty list",
			input:    `Apply(Plus, [])`,
			expected: `0`,
		},
		{
			name:     "Apply with two arguments",
			input:    `Apply(Power, [2, 8])`,
			expected: `256`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result.String())
			}
		})
	}
}

func TestApply_WithAmpersandSyntax(t *testing.T) {
	evaluator := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Apply with & syntax - two arguments",
			input:    `Apply(Plus($1, $2) &, [10, 20])`,
			expected: `30`,
		},
		{
			name:     "Apply with & syntax - three arguments",
			input:    `Apply($1 + $2 + $3 &, [10, 20, 30])`,
			expected: `60`,
		},
		{
			name:     "Apply with & syntax - multiplication",
			input:    `Apply(Times($1, $2, $3) &, [2, 3, 4])`,
			expected: `24`,
		},
		{
			name:     "Apply with & syntax - complex expression",
			input:    `Apply(Plus(Times($1, $2), $3) &, [3, 4, 5])`,
			expected: `17`,
		},
		{
			name:     "Apply with & syntax - single argument",
			input:    `Apply(Times($1, $1) &, [7])`,
			expected: `49`,
		},
		{
			name:     "Apply with & syntax - higher slots",
			input:    `Apply($1 + $2 + $3 + $4 &, [1, 10, 100, 1000])`,
			expected: `1111`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result.String())
			}
		})
	}
}

func TestApply_WithRegularFunction(t *testing.T) {
	evaluator := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Apply with Function syntax - two parameters",
			input:    `Apply(Function([x, y], x + y), [15, 25])`,
			expected: `40`,
		},
		{
			name:     "Apply with Function syntax - three parameters",
			input:    `Apply(Function([a, b, c], a * b + c), [2, 5, 3])`,
			expected: `13`,
		},
		{
			name:     "Apply with Function syntax - single parameter",
			input:    `Apply(Function([x], x * x), [6])`,
			expected: `36`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result.String())
			}
		})
	}
}

func TestApply_EdgeCases(t *testing.T) {
	evaluator := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Apply with non-list should error",
			input:    `Apply(Plus, 42)`,
			expected: `ErrorExpr(ArgumentError)`,
		},
		{
			name:     "Apply with wrong number of arguments returns unevaluated",
			input:    `Apply(Plus)`,
			expected: `Apply(Plus)`,
		},
		{
			name:     "Apply with too many arguments returns unevaluated",
			input:    `Apply(Plus, [1, 2], [3, 4])`,
			expected: `Apply(Plus, List(1, 2), List(3, 4))`,
		},
		{
			name:     "Apply ignores list head",
			input:    `Apply(Plus, MyList(5, 10, 15))`,
			expected: `30`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if core.IsError(result) {
				if !strings.Contains(test.expected, "ErrorExpr") {
					t.Errorf("Expected %s, got error: %s", test.expected, result.String())
				}
			} else if result.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result.String())
			}
		})
	}
}

func TestMapApply_Integration(t *testing.T) {
	evaluator := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Map then Apply - sum of squares",
			input:    `Apply(Plus, Map(Times($1, $1) &, [1, 2, 3, 4]))`,
			expected: `30`,
		},
		{
			name:     "Apply then Map - distribute and square",
			input:    `Map(Times($1, $1) &, Apply(List, [2, 3, 4]))`,
			expected: `List(4, 9, 16)`,
		},
		{
			name:     "Nested Map with Apply",
			input:    `Map(Apply(Plus, $1) &, [[1, 2], [3, 4], [5, 6]])`,
			expected: `List(3, 7, 11)`,
		},
		{
			name:     "Apply with Map as function",
			input:    `Apply(Function([x], Map(Plus($1, 1) &, x)), [[1, 2, 3]])`,
			expected: `List(2, 3, 4)`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result.String())
			}
		})
	}
}

func TestMapApply_WithBuiltinFunctions(t *testing.T) {
	evaluator := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Map with IntegerQ predicate",
			input:    `Map(IntegerQ, [1, 2.5, 3, "hello"])`,
			expected: `List(True, False, True, False)`,
		},
		{
			name:     "Map with StringQ predicate",
			input:    `Map(StringQ, [1, "hello", 3.14, "world"])`,
			expected: `List(False, True, False, True)`,
		},
		{
			name:     "Apply with Equal comparison",
			input:    `Apply(Equal, [5, 5])`,
			expected: `True`,
		},
		{
			name:     "Apply with Less comparison",
			input:    `Apply(Less, [3, 7])`,
			expected: `True`,
		},
		{
			name:     "Map with Not function",
			input:    `Map(Not, [True, False, True])`,
			expected: `List(False, True, False)`,
		},
		{
			name:     "Apply with MatchQ",
			input:    `Apply(MatchQ, [42, _Integer])`,
			expected: `True`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result.String())
			}
		})
	}
}
