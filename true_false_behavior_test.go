package sexpr

import (
	"github.com/client9/sexpr/core"
	"testing"
)

// TestTrueFalseBehavior tests that True and False behave like Mathematica symbols
func TestTrueFalseBehavior(t *testing.T) {
	eval := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		// Head behavior - True and False should have Head = Symbol (not Boolean)
		{
			name:     "Head of True is Symbol",
			input:    "Head(True)",
			expected: "Symbol",
			desc:     "In Mathematica, Head[True] returns Symbol, not Boolean",
		},
		{
			name:     "Head of False is Symbol",
			input:    "Head(False)",
			expected: "Symbol",
			desc:     "In Mathematica, Head[False] returns Symbol, not Boolean",
		},

		// Identity behavior - True and False should remain as symbols
		{
			name:     "True remains as symbol",
			input:    "True",
			expected: "True",
			desc:     "True should not evaluate to a different form",
		},
		{
			name:     "False remains as symbol",
			input:    "False",
			expected: "False",
			desc:     "False should not evaluate to a different form",
		},

		// Attributes behavior
		{
			name:     "True has Constant and Protected attributes",
			input:    "Attributes(True)",
			expected: "List(Constant, Protected)",
			desc:     "True should have both Constant and Protected attributes",
		},
		{
			name:     "False has Constant and Protected attributes",
			input:    "Attributes(False)",
			expected: "List(Constant, Protected)",
			desc:     "False should have both Constant and Protected attributes",
		},

		// Type testing behavior
		{
			name:     "BooleanQ recognizes True",
			input:    "BooleanQ(True)",
			expected: "True",
			desc:     "BooleanQ should recognize True as a boolean even though it's a symbol",
		},
		{
			name:     "BooleanQ recognizes False",
			input:    "BooleanQ(False)",
			expected: "True",
			desc:     "BooleanQ should recognize False as a boolean even though it's a symbol",
		},
		{
			name:     "SymbolQ recognizes True",
			input:    "SymbolQ(True)",
			expected: "True",
			desc:     "True should be recognized as a symbol",
		},
		{
			name:     "SymbolQ recognizes False",
			input:    "SymbolQ(False)",
			expected: "True",
			desc:     "False should be recognized as a symbol",
		},

		// Logical operations behavior
		{
			name:     "And with True and False",
			input:    "And(True, False)",
			expected: "False",
			desc:     "Logical And should work with True/False symbols",
		},
		{
			name:     "Or with True and False",
			input:    "Or(True, False)",
			expected: "True",
			desc:     "Logical Or should work with True/False symbols",
		},
		{
			name:     "Not with True",
			input:    "Not(True)",
			expected: "False",
			desc:     "Logical Not should work with True symbol",
		},
		{
			name:     "Not with False",
			input:    "Not(False)",
			expected: "True",
			desc:     "Logical Not should work with False symbol",
		},

		// Complex logical operations
		{
			name:     "Nested logical operations",
			input:    "And(True, Or(False, True))",
			expected: "True",
			desc:     "Complex logical expressions should work correctly",
		},
		{
			name:     "Short-circuit evaluation with And",
			input:    "And(False, undefined_symbol)",
			expected: "False",
			desc:     "And should short-circuit and not evaluate undefined_symbol",
		},
		{
			name:     "Short-circuit evaluation with Or",
			input:    "Or(True, undefined_symbol)",
			expected: "True",
			desc:     "Or should short-circuit and not evaluate undefined_symbol",
		},

		// FullForm behavior
		{
			name:     "FullForm of True",
			input:    "FullForm(True)",
			expected: "\"True\"",
			desc:     "FullForm should return the string representation of the True symbol",
		},
		{
			name:     "FullForm of False",
			input:    "FullForm(False)",
			expected: "\"False\"",
			desc:     "FullForm should return the string representation of the False symbol",
		},

		// Pattern matching behavior
		{
			name:     "MatchQ True with True",
			input:    "MatchQ(True, True)",
			expected: "True",
			desc:     "True should match itself exactly",
		},
		{
			name:     "MatchQ False with False",
			input:    "MatchQ(False, False)",
			expected: "True",
			desc:     "False should match itself exactly",
		},
		{
			name:     "MatchQ True with symbol pattern",
			input:    "MatchQ(True, _Symbol)",
			expected: "True",
			desc:     "True should match symbol patterns",
		},
		{
			name:     "MatchQ False with symbol pattern",
			input:    "MatchQ(False, _Symbol)",
			expected: "True",
			desc:     "False should match symbol patterns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error for input '%s': %v", tt.input, err)
			}

			result := eval.Evaluate(expr)
			resultStr := result.String()

			if resultStr != tt.expected {
				t.Errorf("Input: %s\nExpected: %s\nGot: %s\nDescription: %s",
					tt.input, tt.expected, resultStr, tt.desc)
			}
		})
	}
}

// TestTrueFalseInternalConsistency tests that True/False behavior is internally consistent
func TestTrueFalseInternalConsistency(t *testing.T) {
	// Test that True and False symbols are recognized consistently by utility functions
	trueExpr := core.NewBool(true)
	falseExpr := core.NewBool(false)

	t.Run("isBool recognizes True symbol", func(t *testing.T) {
		if !core.IsBool(trueExpr) {
			t.Error("isBool should recognize True symbol as boolean")
		}
	})

	t.Run("isBool recognizes False symbol", func(t *testing.T) {
		if !core.IsBool(falseExpr) {
			t.Error("isBool should recognize False symbol as boolean")
		}
	})

	t.Run("getBoolValue extracts True", func(t *testing.T) {
		val, ok := getBoolValue(trueExpr)
		if !ok || !val {
			t.Error("getBoolValue should extract true from True symbol")
		}
	})

	t.Run("getBoolValue extracts False", func(t *testing.T) {
		val, ok := getBoolValue(falseExpr)
		if !ok || val {
			t.Error("getBoolValue should extract false from False symbol")
		}
	})

	t.Run("isSymbol recognizes True", func(t *testing.T) {
		if !isSymbol(trueExpr) {
			t.Error("isSymbol should recognize True as a symbol")
		}
	})

	t.Run("isSymbol recognizes False", func(t *testing.T) {
		if !isSymbol(falseExpr) {
			t.Error("isSymbol should recognize False as a symbol")
		}
	})
}

// TestTrueFalseLexerParserBehavior tests that True/False are lexed and parsed correctly
func TestTrueFalseLexerParserBehavior(t *testing.T) {
	// Test that True and False are lexed as SYMBOL tokens, not BOOLEAN tokens
	t.Run("True lexed as SYMBOL", func(t *testing.T) {
		expr, err := ParseString("True")
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		symbolName, ok := core.ExtractSymbol(expr)
		if !ok {
			t.Fatalf("Expected Symbol, got %T", expr)
		}

		if symbolName != "True" {
			t.Errorf("Expected 'True', got %v", symbolName)
		}
	})

	t.Run("False lexed as SYMBOL", func(t *testing.T) {
		expr, err := ParseString("False")
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		symbolName, ok := core.ExtractSymbol(expr)
		if !ok {
			t.Fatalf("Expected Symbol, got %T", expr)
		}

		if symbolName != "False" {
			t.Errorf("Expected 'False', got %v", symbolName)
		}
	})
}

// TestTrueFalseComparisonWithMathematica documents expected Mathematica equivalences
func TestTrueFalseComparisonWithMathematica(t *testing.T) {
	eval := NewEvaluator()

	// Document expected behavior compared to Mathematica
	mathematicaTests := []struct {
		name           string
		input          string
		ourResult      string
		mathematicaRef string
		desc           string
	}{
		{
			name:           "Head[True] in Mathematica",
			input:          "Head(True)",
			ourResult:      "Symbol",
			mathematicaRef: "Symbol",
			desc:           "Mathematica: Head[True] returns Symbol",
		},
		{
			name:           "Head[False] in Mathematica",
			input:          "Head(False)",
			ourResult:      "Symbol",
			mathematicaRef: "Symbol",
			desc:           "Mathematica: Head[False] returns Symbol",
		},
		{
			name:           "SymbolQ[True] in Mathematica",
			input:          "SymbolQ(True)",
			ourResult:      "True",
			mathematicaRef: "True",
			desc:           "Mathematica: SymbolQ[True] returns True",
		},
		{
			name:           "Attributes[True] in Mathematica",
			input:          "Attributes(True)",
			ourResult:      "List(Constant, Protected)",
			mathematicaRef: "{Constant, Protected}",
			desc:           "Mathematica: Attributes[True] returns {Constant, Protected}",
		},
	}

	for _, tt := range mathematicaTests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			resultStr := result.String()

			if resultStr != tt.ourResult {
				t.Errorf("Our result differs from expected\nInput: %s\nExpected: %s\nGot: %s\nMathematica reference: %s\nNote: %s",
					tt.input, tt.ourResult, resultStr, tt.mathematicaRef, tt.desc)
			}

			t.Logf("✓ %s: Our result '%s' matches expected (Mathematica: %s)",
				tt.desc, resultStr, tt.mathematicaRef)
		})
	}
}
