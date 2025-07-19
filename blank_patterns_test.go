package sexpr

import (
	"strings"
	"testing"
)

func TestBlankPatterns(t *testing.T) {
	// Test Blank pattern (_) - matches exactly one expression
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Blank pattern basic", "f(x_) := x; f(42)", "42"},
		{"Blank pattern with type", "f(x_Integer) := x; f(42)", "42"},
		{"Blank pattern with type fail", "f(x_Integer) := x; f(\"hello\")", "f(\"hello\")"},
		{"Anonymous blank pattern", "f(_) := \"matched\"; f(42)", "\"matched\""},
		{"Anonymous blank pattern with type", "f(_Integer) := \"matched\"; f(42)", "\"matched\""},
		{"Anonymous blank pattern with type fail", "f(_Integer) := \"matched\"; f(\"hello\")", "f(\"hello\")"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator() // Create fresh evaluator for each test
			result := evaluateString(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBlankSequencePatterns(t *testing.T) {
	// Test BlankSequence pattern (__) - matches one or more expressions
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"BlankSequence basic one arg", "f(x__) := x; f(42)", "List(42)"},
		{"BlankSequence basic two args", "f(x__) := x; f(42, 43)", "List(42, 43)"},
		{"BlankSequence basic three args", "f(x__) := x; f(42, 43, 44)", "List(42, 43, 44)"},
		{"BlankSequence with type", "f(x__Integer) := x; f(42, 43)", "List(42, 43)"},
		{"BlankSequence with type fail", "f(x__Integer) := x; f(42, \"hello\")", "f(42, \"hello\")"},
		{"BlankSequence no args fail", "f(x__) := x; f()", "f()"},
		{"Anonymous BlankSequence", "f(__) := \"matched\"; f(42)", "\"matched\""},
		{"Anonymous BlankSequence multi", "f(__) := \"matched\"; f(42, 43)", "\"matched\""},
		{"Anonymous BlankSequence with type", "f(__Integer) := \"matched\"; f(42, 43)", "\"matched\""},
		{"BlankSequence with following pattern", "f(x__, y_) := Subtract(x, y); f(1, 2, 3)", "Subtract(List(1, 2), 3)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator() // Create fresh evaluator for each test
			result := evaluateString(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBlankNullSequencePatterns(t *testing.T) {
	// Test BlankNullSequence pattern (___) - matches zero or more expressions
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"BlankNullSequence zero args", "f(x___) := x; f()", "List()"},
		{"BlankNullSequence one arg", "f(x___) := x; f(42)", "List(42)"},
		{"BlankNullSequence two args", "f(x___) := x; f(42, 43)", "List(42, 43)"},
		{"BlankNullSequence three args", "f(x___) := x; f(42, 43, 44)", "List(42, 43, 44)"},
		{"BlankNullSequence with type", "f(x___Integer) := x; f(42, 43)", "List(42, 43)"},
		{"BlankNullSequence with type fail", "f(x___Integer) := x; f(42, \"hello\")", "f(42, \"hello\")"},
		{"Anonymous BlankNullSequence", "f(___) := \"matched\"; f()", "\"matched\""},
		{"Anonymous BlankNullSequence multi", "f(___) := \"matched\"; f(42, 43)", "\"matched\""},
		{"Anonymous BlankNullSequence with type", "f(___Integer) := \"matched\"; f(42, 43)", "\"matched\""},
		{"BlankNullSequence with following pattern", "f(x___, y_) := Subtract(x, y); f(1, 2, 3)", "Subtract(List(1, 2), 3)"},
		{"BlankNullSequence empty with following pattern", "f(x___, y_) := Subtract(x, y); f(42)", "Subtract(List(), 42)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator() // Create fresh evaluator for each test
			result := evaluateString(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestMixedBlankPatterns(t *testing.T) {
	// Test combinations of different blank patterns
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Blank + BlankSequence", "f(x_, y__) := Subtract(x, y); f(1, 2, 3)", "Subtract(1, List(2, 3))"},
		{"Blank + BlankNullSequence", "f(x_, y___) := Subtract(x, y); f(1, 2, 3)", "Subtract(1, List(2, 3))"},
		{"Blank + BlankNullSequence empty", "f(x_, y___) := Subtract(x, y); f(1)", "Subtract(1, List())"},
		{"Multiple Blanks", "f(x_, y_, z_) := [x, y, z]; f(1, 2, 3)", "List(1, 2, 3)"},
		{"Blank + Blank + BlankSequence", "f(x_, y_, z__) := [x, y, z]; f(1, 2, 3, 4)", "List(1, 2, List(3, 4))"},
		{"Blank + Blank + BlankNullSequence", "f(x_, y_, z___) := [x, y, z]; f(1, 2, 3, 4)", "List(1, 2, List(3, 4))"},
		{"Blank + Blank + BlankNullSequence empty", "f(x_, y_, z___) := [x, y, z]; f(1, 2)", "List(1, 2, List())"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator() // Create fresh evaluator for each test
			result := evaluateString(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBlankPatternsWithTypes(t *testing.T) {
	// Test all blank patterns with various type constraints
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Integer types
		{"Blank Integer", "f(x_Integer) := x; f(42)", "42"},
		{"BlankSequence Integer", "f(x__Integer) := x; f(42, 43)", "List(42, 43)"},
		{"BlankNullSequence Integer", "f(x___Integer) := x; f(42, 43)", "List(42, 43)"},

		// Real types
		{"Blank Real", "f(x_Real) := x; f(3.14)", "3.14"},
		{"BlankSequence Real", "f(x__Real) := x; f(3.14, 2.71)", "List(3.14, 2.71)"},
		{"BlankNullSequence Real", "f(x___Real) := x; f(3.14, 2.71)", "List(3.14, 2.71)"},

		// Number types (Integer or Real)
		{"Blank Number", "f(x_Number) := x; f(42)", "42"},
		{"BlankSequence Number", "f(x__Number) := x; f(42, 3.14)", "List(42, 3.14)"},
		{"BlankNullSequence Number", "f(x___Number) := x; f(42, 3.14)", "List(42, 3.14)"},

		// String types
		{"Blank String", "f(x_String) := x; f(\"hello\")", "\"hello\""},
		{"BlankSequence String", "f(x__String) := x; f(\"hello\", \"world\")", "List(\"hello\", \"world\")"},
		{"BlankNullSequence String", "f(x___String) := x; f(\"hello\", \"world\")", "List(\"hello\", \"world\")"},

		// Boolean symbols (True/False are symbols in our Mathematica-compatible system)
		{"Blank Boolean as Symbol", "f(x_Symbol) := x; f(True)", "True"},
		{"BlankSequence Boolean as Symbol", "f(x__Symbol) := x; f(True, False)", "List(True, False)"},
		{"BlankNullSequence Boolean as Symbol", "f(x___Symbol) := x; f(True, False)", "List(True, False)"},

		// Symbol types
		{"Blank Symbol", "f(x_Symbol) := x; f(abc)", "abc"},
		{"BlankSequence Symbol", "f(x__Symbol) := x; f(abc, def)", "List(abc, def)"},
		{"BlankNullSequence Symbol", "f(x___Symbol) := x; f(abc, def)", "List(abc, def)"},

		// List types
		{"Blank List", "f(x_List) := x; f([1, 2, 3])", "List(1, 2, 3)"},
		{"BlankSequence List", "f(x__List) := x; f([1, 2], [3, 4])", "List(List(1, 2), List(3, 4))"},
		{"BlankNullSequence List", "f(x___List) := x; f([1, 2], [3, 4])", "List(List(1, 2), List(3, 4))"},

		// Atom types
		{"Blank Atom", "f(x_Atom) := x; f(42)", "42"},
		{"BlankSequence Atom", "f(x__Atom) := x; f(42, \"hello\")", "List(42, \"hello\")"},
		{"BlankNullSequence Atom", "f(x___Atom) := x; f(42, \"hello\")", "List(42, \"hello\")"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator() // Create fresh evaluator for each test
			result := evaluateString(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBlankPatternFailures(t *testing.T) {
	// Test cases that should fail to match
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Type mismatches
		{"Blank Integer fail", "f(x_Integer) := x; f(\"hello\")", "f(\"hello\")"},
		{"BlankSequence Integer fail", "f(x__Integer) := x; f(42, \"hello\")", "f(42, \"hello\")"},
		{"BlankNullSequence Integer fail", "f(x___Integer) := x; f(42, \"hello\")", "f(42, \"hello\")"},

		// BlankSequence with no arguments
		{"BlankSequence no args", "f(x__) := x; f()", "f()"},
		{"BlankSequence Integer no args", "f(x__Integer) := x; f()", "f()"},

		// Complex type constraints
		{"Mixed types fail", "f(x__Integer, y_String) := Plus(x, y); f(42, \"hello\", 3.14)", "f(42, \"hello\", 3.14)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator() // Create fresh evaluator for each test
			result := evaluateString(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBlankPatternParsing(t *testing.T) {
	// Test the parsePatternInfo function directly
	tests := []struct {
		name     string
		input    string
		expected PatternInfo
	}{
		{"Simple blank", "_", PatternInfo{Type: BlankPattern, VarName: "", TypeName: ""}},
		{"Named blank", "x_", PatternInfo{Type: BlankPattern, VarName: "x", TypeName: ""}},
		{"Typed blank", "_Integer", PatternInfo{Type: BlankPattern, VarName: "", TypeName: "Integer"}},
		{"Named typed blank", "x_Integer", PatternInfo{Type: BlankPattern, VarName: "x", TypeName: "Integer"}},

		{"Simple blank sequence", "__", PatternInfo{Type: BlankSequencePattern, VarName: "", TypeName: ""}},
		{"Named blank sequence", "x__", PatternInfo{Type: BlankSequencePattern, VarName: "x", TypeName: ""}},
		{"Typed blank sequence", "__Integer", PatternInfo{Type: BlankSequencePattern, VarName: "", TypeName: "Integer"}},
		{"Named typed blank sequence", "x__Integer", PatternInfo{Type: BlankSequencePattern, VarName: "x", TypeName: "Integer"}},

		{"Simple blank null sequence", "___", PatternInfo{Type: BlankNullSequencePattern, VarName: "", TypeName: ""}},
		{"Named blank null sequence", "x___", PatternInfo{Type: BlankNullSequencePattern, VarName: "x", TypeName: ""}},
		{"Typed blank null sequence", "___Integer", PatternInfo{Type: BlankNullSequencePattern, VarName: "", TypeName: "Integer"}},
		{"Named typed blank null sequence", "x___Integer", PatternInfo{Type: BlankNullSequencePattern, VarName: "x", TypeName: "Integer"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePatternInfo(tt.input)
			if result.Type != tt.expected.Type || result.VarName != tt.expected.VarName || result.TypeName != tt.expected.TypeName {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

// Helper function to evaluate a string and return the result as a string
func evaluateString(t *testing.T, evaluator *Evaluator, input string) string {
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
