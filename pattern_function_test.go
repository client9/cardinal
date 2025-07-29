package sexpr

import (
	"github.com/client9/sexpr/core"
	"strings"
	"testing"
)

func TestPatternFunctionSystem(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Test basic arithmetic with new pattern system
		{"Plus two args", "Plus(1, 2)", "3"},
		{"Plus three args", "Plus(1, 2, 3)", "6"},
		{"Times two args", "Times(2, 3)", "6"},

		// Test user function override capability
		{"User function", "f(x_) := x + 1; f(5)", "6"},
		{"User function sequence", "g(x__) := Length(x); g(1, 2, 3)", "3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluatePatternTestHelper(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPatternFunctionRegistry(t *testing.T) {
	// Test the function registry directly
	registry := NewFunctionRegistry()

	// Register a simple pattern
	err := registry.RegisterPatternBuiltin("test(x_)", func(args []core.Expr, ctx *Context) core.Expr {
		if len(args) != 1 {
			return core.NewSymbol("error")
		}
		return args[0] // Return the argument unchanged
	})

	if err != nil {
		t.Fatalf("Failed to register pattern: %v", err)
	}

	// Test pattern matching
	funcDef, bindings := registry.FindMatchingFunction("test", []core.Expr{core.NewInteger(42)})

	if funcDef == nil {
		t.Fatalf("Expected to find matching function")
	}

	if len(bindings) != 1 {
		t.Errorf("Expected 1 binding, got %d", len(bindings))
	}

	if bindings["x"] == nil {
		t.Errorf("Expected binding for variable 'x'")
	}
}

func evaluatePatternTestHelper(t *testing.T, evaluator *Evaluator, input string) string {
	// Split by semicolon and evaluate each part
	parts := strings.Split(input, ";")

	var result core.Expr
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		expr, err := ParseString(part)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		result = evaluator.Evaluate(expr)
	}

	return result.String()
}
