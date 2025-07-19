package sexpr

import (
	"strings"
	"testing"
)

func TestEnhancedTypeSystem(t *testing.T) {
	// Test that typed patterns work with any symbol name, not just built-in types
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Custom type patterns
		{"Color pattern", "g(x_Color) := First(x); g(Color(1, 2, 3))", "1"},
		{"Point pattern", "g(x_Point) := \"got point\"; g(Point(0, 0))", "\"got point\""},
		{"Complex pattern", "g(x_Complex) := \"complex\"; g(Complex(1, 2))", "\"complex\""},

		// Mixed custom and built-in types
		{"Mixed types", "h(x_Integer) := \"int\"; h(x_Color) := \"color\"; h(42)", "\"int\""},
		{"Mixed types color", "h(x_Integer) := \"int\"; h(x_Color) := \"color\"; h(Color(1,2,3))", "\"color\""},

		// Custom types with sequence patterns
		{"Color sequence", "f(x__Color) := x; f(Color(1,2), Color(3,4))", "List(Color(1, 2), Color(3, 4))"},
		{"Color null sequence", "f(x___Color) := x; f()", "List()"},
		{"Color null sequence with args", "f(x___Color) := x; f(Color(1,2))", "List(Color(1, 2))"},

		// Type constraint failures
		{"Color constraint fail", "g(x_Color) := \"matched\"; g(42)", "g(42)"},
		{"Point constraint fail", "g(x_Point) := \"matched\"; g(\"hello\")", "g(\"hello\")"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringHelper(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestSpecificityBasedOrdering(t *testing.T) {
	// Test that pattern specificity works regardless of definition order
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// General first, specific second (should work due to auto-sorting)
		{"General then specific", "f(x_) := \"general\"; f(x_Integer) := \"specific\"; f(42)", "\"specific\""},
		{"General then custom", "f(x_) := \"general\"; f(x_Color) := \"color\"; f(Color(1,2,3))", "\"color\""},

		// Sequence patterns in different orders
		{"Null seq then seq", "g(x___) := \"null\"; g(x__) := \"seq\"; g(1, 2)", "\"seq\""},
		{"Seq then single", "g(x__) := \"seq\"; g(x_) := \"single\"; g(1)", "\"single\""},
		{"Random order", "h(x___) := \"null\"; h(x_) := \"single\"; h(x__) := \"seq\"; h(1)", "\"single\""},

		// Literal patterns (highest specificity)
		{"Literal wins", "f(x_Integer) := \"int\"; f(42) := \"forty-two\"; f(42)", "\"forty-two\""},
		{"Literal vs general", "f(x_) := \"general\"; f(0) := \"zero\"; f(0)", "\"zero\""},

		// Complex specificity hierarchy
		{"Full hierarchy", "test(x___) := \"null\"; test(x__) := \"seq\"; test(x_) := \"single\"; test(x_Integer) := \"int\"; test(x_Color) := \"color\"; test(42) := \"literal\"; test(42)", "\"literal\""},
		{"Full hierarchy int", "test(x___) := \"null\"; test(x__) := \"seq\"; test(x_) := \"single\"; test(x_Integer) := \"int\"; test(x_Color) := \"color\"; test(42) := \"literal\"; test(99)", "\"int\""},
		{"Full hierarchy color", "test(x___) := \"null\"; test(x__) := \"seq\"; test(x_) := \"single\"; test(x_Integer) := \"int\"; test(x_Color) := \"color\"; test(42) := \"literal\"; test(Color(1,2,3))", "\"color\""},
		{"Full hierarchy general", "test(x___) := \"null\"; test(x__) := \"seq\"; test(x_) := \"single\"; test(x_Integer) := \"int\"; test(x_Color) := \"color\"; test(42) := \"literal\"; test(\"hello\")", "\"single\""},
		{"Full hierarchy seq", "test(x___) := \"null\"; test(x__) := \"seq\"; test(x_) := \"single\"; test(x_Integer) := \"int\"; test(x_Color) := \"color\"; test(42) := \"literal\"; test(1, 2)", "\"seq\""},
		{"Full hierarchy null", "test(x___) := \"null\"; test(x__) := \"seq\"; test(x_) := \"single\"; test(x_Integer) := \"int\"; test(x_Color) := \"color\"; test(42) := \"literal\"; test()", "\"null\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringHelper(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestSpecificityRules(t *testing.T) {
	// Test the specificity calculation directly
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Built-in vs custom type specificity
		{"Builtin vs custom", "f(x_Integer) := \"builtin\"; f(x_Color) := \"custom\"; f(Color(1,2,3))", "\"custom\""},
		{"Custom vs builtin", "f(x_Color) := \"custom\"; f(x_Integer) := \"builtin\"; f(Color(1,2,3))", "\"custom\""},

		// Type vs non-type patterns
		{"Type vs general", "f(x_) := \"general\"; f(x_String) := \"string\"; f(\"hello\")", "\"string\""},

		// Sequence type constraints
		{"Typed sequence", "f(x__Integer) := \"int seq\"; f(x__) := \"seq\"; f(1, 2, 3)", "\"int seq\""},
		{"Typed null sequence", "f(x___String) := \"str null\"; f(x___) := \"null\"; f(\"a\", \"b\")", "\"str null\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringHelper(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestComplexPatternScenarios(t *testing.T) {
	// Test complex real-world scenarios
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Vector operations with different types
		{"Vector add", "add(Vector(x__), Vector(y__)) := Vector(Plus(x, y)); add(Vector(1, 2), Vector(3, 4))", "Vector(Plus(List(1, 2), List(3, 4)))"},

		// Graphics operations
		{"Graphics color", "draw(Shape(x_), Color(r_, g_, b_)) := [\"draw\", x, r, g, b]; draw(Shape(\"circle\"), Color(255, 0, 0))", "List(\"draw\", \"circle\", 255, 0, 0)"},

		// Recursive patterns
		{"Tree processing", "count(Tree(value_, children___)) := Plus(1, count(children)); count(Leaf(x_)) := 1; count(Leaf(42))", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringHelper(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Helper function to evaluate a string and return the result as a string
func evaluateStringHelper(t *testing.T, evaluator *Evaluator, input string) string {
	// Split by semicolon and evaluate each part
	parts := strings.Split(input, ";")
	var result Expr

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
