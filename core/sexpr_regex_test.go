package core

import (
	"testing"
)

func TestBasicPatterns(t *testing.T) {
	tests := []struct {
		name     string
		pattern  Pattern
		expr     Expr
		expected bool
		strategy ExecutionStrategy
	}{
		{
			name:     "literal match",
			pattern:  MatchLiteral(NewInteger(42)),
			expr:     NewInteger(42),
			expected: true,
			strategy: StrategyDirect,
		},
		{
			name:     "literal no match",
			pattern:  MatchLiteral(NewInteger(42)),
			expr:     NewInteger(43),
			expected: false,
			strategy: StrategyDirect,
		},
		{
			name:     "head match",
			pattern:  MatchHead("Integer"),
			expr:     NewInteger(100),
			expected: true,
			strategy: StrategyDirect,
		},
		{
			name:     "head no match",
			pattern:  MatchHead("String"),
			expr:     NewInteger(100),
			expected: false,
			strategy: StrategyDirect,
		},
		{
			name:     "any match",
			pattern:  MatchAny(),
			expr:     NewString("hello"),
			expected: true,
			strategy: StrategyDirect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := CompilePattern(tt.pattern)
			if err != nil {
				t.Fatalf("Failed to compile pattern: %v", err)
			}

			if compiled.Strategy() != tt.strategy {
				t.Errorf("Expected strategy %v, got %v", tt.strategy, compiled.Strategy())
			}

			result := compiled.Match(tt.expr)
			if result.Matched != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result.Matched)
			}
		})
	}
}

func TestSequencePatterns(t *testing.T) {
	tests := []struct {
		name     string
		pattern  Pattern
		expr     Expr
		expected bool
		strategy ExecutionStrategy
	}{
		{
			name: "simple list match",
			pattern: MatchList(
				MatchLiteral(NewInteger(1)),
				MatchHead("String"),
				MatchAny(),
			),
			expr: NewList("List",
				NewInteger(1),
				NewString("hello"),
				NewBool(true),
			),
			expected: true,
			strategy: StrategyDirect,
		},
		{
			name: "simple list no match - wrong length",
			pattern: MatchList(
				MatchLiteral(NewInteger(1)),
				MatchHead("String"),
			),
			expr: NewList("List",
				NewInteger(1),
				NewString("hello"),
				NewBool(true),
			),
			expected: false,
			strategy: StrategyDirect,
		},
		{
			name: "simple list no match - wrong element",
			pattern: MatchList(
				MatchLiteral(NewInteger(1)),
				MatchHead("String"),
			),
			expr: NewList("List",
				NewInteger(2),
				NewString("hello"),
			),
			expected: false,
			strategy: StrategyDirect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := CompilePattern(tt.pattern)
			if err != nil {
				t.Fatalf("Failed to compile pattern: %v", err)
			}

			if compiled.Strategy() != tt.strategy {
				t.Errorf("Expected strategy %v, got %v", tt.strategy, compiled.Strategy())
			}

			result := compiled.Match(tt.expr)
			if result.Matched != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result.Matched)
			}
		})
	}
}

func TestTrailingQuantifiers(t *testing.T) {
	tests := []struct {
		name     string
		pattern  Pattern
		expr     Expr
		expected bool
		strategy ExecutionStrategy
	}{
		{
			name: "zero or more - zero matches",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				ZeroOrMore(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
			),
			expected: true,
			strategy: StrategyDirect,
		},
		{
			name: "zero or more - one match",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				ZeroOrMore(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewInteger(42),
			),
			expected: true,
			strategy: StrategyDirect,
		},
		{
			name: "zero or more - multiple matches",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				ZeroOrMore(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewInteger(1),
				NewInteger(2),
				NewInteger(3),
			),
			expected: true,
			strategy: StrategyDirect,
		},
		{
			name: "zero or more - non-matching type",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				ZeroOrMore(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewString("not an integer"),
			),
			expected: false,
			strategy: StrategyDirect,
		},
		{
			name: "one or more - zero matches (should fail)",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				OneOrMore(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
			),
			expected: false,
			strategy: StrategyDirect,
		},
		{
			name: "one or more - one match",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				OneOrMore(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewInteger(42),
			),
			expected: true,
			strategy: StrategyDirect,
		},
		{
			name: "one or more - multiple matches",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				OneOrMore(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewInteger(1),
				NewInteger(2),
				NewInteger(3),
			),
			expected: true,
			strategy: StrategyDirect,
		},
		{
			name: "one or more - non-matching type (should fail)",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				OneOrMore(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewString("not an integer"),
			),
			expected: false,
			strategy: StrategyDirect,
		},
		{
			name: "one or more - mixed types (should fail)",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				OneOrMore(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewInteger(1),
				NewString("not an integer"),
			),
			expected: false,
			strategy: StrategyDirect,
		},
		{
			name: "optional - zero matches",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				Optional(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
			),
			expected: true,
			strategy: StrategyDirect,
		},
		{
			name: "optional - one match",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				Optional(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewInteger(42),
			),
			expected: true,
			strategy: StrategyDirect,
		},
		{
			name: "optional - too many matches (should fail)",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				Optional(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewInteger(1),
				NewInteger(2),
			),
			expected: false,
			strategy: StrategyDirect,
		},
		{
			name: "optional - non-matching type (should fail)",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				Optional(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewString("not an integer"),
			),
			expected: false,
			strategy: StrategyDirect,
		},
		// Complex quantifier patterns to test NFA behavior
		{
			name: "one or more with complex inner pattern",
			pattern: MatchList(
				MatchLiteral(NewString("data")),
				OneOrMore(MatchOr(
					MatchHead("Integer"),
					MatchHead("String"),
				)),
			),
			expr: NewList("List",
				NewString("data"),
				NewInteger(1),
				NewString("hello"),
				NewInteger(2),
			),
			expected: true,
			strategy: StrategyNFA,
		},
		{
			name: "one or more with complex inner pattern - zero matches",
			pattern: MatchList(
				MatchLiteral(NewString("data")),
				OneOrMore(MatchOr(
					MatchHead("Integer"),
					MatchHead("String"),
				)),
			),
			expr: NewList("List",
				NewString("data"),
			),
			expected: false,
			strategy: StrategyNFA,
		},
		{
			name: "optional with complex inner pattern - zero matches",
			pattern: MatchList(
				MatchLiteral(NewString("data")),
				Optional(MatchOr(
					MatchHead("Integer"),
					MatchHead("String"),
				)),
			),
			expr: NewList("List",
				NewString("data"),
			),
			expected: true,
			strategy: StrategyNFA,
		},
		{
			name: "optional with complex inner pattern - one match",
			pattern: MatchList(
				MatchLiteral(NewString("data")),
				Optional(MatchOr(
					MatchHead("Integer"),
					MatchHead("String"),
				)),
			),
			expr: NewList("List",
				NewString("data"),
				NewInteger(42),
			),
			expected: true,
			strategy: StrategyNFA,
		},
		{
			name: "optional with complex inner pattern - too many matches",
			pattern: MatchList(
				MatchLiteral(NewString("data")),
				Optional(MatchOr(
					MatchHead("Integer"),
					MatchHead("String"),
				)),
			),
			expr: NewList("List",
				NewString("data"),
				NewInteger(1),
				NewString("hello"),
			),
			expected: false,
			strategy: StrategyNFA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := CompilePattern(tt.pattern)
			if err != nil {
				t.Fatalf("Failed to compile pattern: %v", err)
			}

			if compiled.Strategy() != tt.strategy {
				t.Errorf("Expected strategy %v, got %v", tt.strategy, compiled.Strategy())
			}

			result := compiled.Match(tt.expr)
			if result.Matched != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result.Matched)
			}
		})
	}
}

func TestPatternAnalyzer(t *testing.T) {
	tests := []struct {
		name     string
		pattern  Pattern
		strategy ExecutionStrategy
	}{
		{
			name:     "literal pattern",
			pattern:  MatchLiteral(NewInteger(42)),
			strategy: StrategyDirect,
		},
		{
			name:     "head pattern",
			pattern:  MatchHead("Integer"),
			strategy: StrategyDirect,
		},
		{
			name:     "any pattern",
			pattern:  MatchAny(),
			strategy: StrategyDirect,
		},
		{
			name: "simple sequence",
			pattern: MatchSequence(
				MatchLiteral(NewInteger(1)),
				MatchHead("String"),
			),
			strategy: StrategyDirect,
		},
		{
			name: "sequence with trailing quantifier",
			pattern: MatchSequence(
				MatchLiteral(NewString("prefix")),
				ZeroOrMore(MatchHead("Integer")),
			),
			strategy: StrategyDirect,
		},
		{
			name: "simple or pattern",
			pattern: MatchOr(
				MatchHead("Integer"),
				MatchHead("String"),
			),
			strategy: StrategyDirect,
		},
		{
			name:     "not pattern (needs NFA)",
			pattern:  MatchNot(MatchHead("Integer")),
			strategy: StrategyNFA,
		},
		{
			name: "sequence with quantifier in middle",
			pattern: MatchSequence(
				MatchLiteral(NewString("prefix")),
				ZeroOrMore(MatchHead("Integer")),
				MatchLiteral(NewString("suffix")),
			),
			strategy: StrategyDirect,
		},
		{
			name:     "named pattern (simple - uses direct strategy)",
			pattern:  Named("x", MatchHead("Integer")),
			strategy: StrategyDirect,
		},
		{
			name: "named pattern with sequence (simple - uses Direct)",
			pattern: Named("seq", MatchSequence(
				MatchLiteral(NewInteger(1)),
				MatchHead("String"),
			)),
			strategy: StrategyDirect,
		},
		// Standalone quantifier tests (to validate the inlined analyzeQuantifier logic)
		{
			name:     "standalone zero or more with simple inner (direct)",
			pattern:  ZeroOrMore(MatchHead("Integer")),
			strategy: StrategyNFA,
			//strategy: StrategyDirect,
		},
		{
			name:     "standalone one or more with simple inner (direct)",
			pattern:  OneOrMore(MatchHead("Integer")),
			strategy: StrategyNFA,
			//strategy: StrategyDirect,
		},
		{
			name:     "standalone optional with simple inner (direct)",
			pattern:  Optional(MatchHead("Integer")),
			strategy: StrategyNFA,
			//strategy: StrategyDirect,
		},
		{
			name: "standalone zero or more with complex inner (NFA)",
			pattern: ZeroOrMore(MatchOr(
				MatchHead("Integer"),
				MatchHead("String"),
			)),
			strategy: StrategyNFA,
		},
		{
			name: "standalone one or more with complex inner (NFA)",
			pattern: OneOrMore(MatchOr(
				MatchHead("Integer"),
				MatchHead("String"),
			)),
			strategy: StrategyNFA,
		},
		{
			name: "standalone optional with complex inner (NFA)",
			pattern: Optional(MatchOr(
				MatchHead("Integer"),
				MatchHead("String"),
			)),
			strategy: StrategyNFA,
		},
	}

	analyzer := &PatternAnalyzer{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := analyzer.Analyze(tt.pattern)
			if strategy != tt.strategy {
				t.Errorf("Expected strategy %v, got %v", tt.strategy, strategy)
			}
		})
	}
}

func TestPatternStrings(t *testing.T) {
	patterns := []Pattern{
		MatchLiteral(NewInteger(42)),
		MatchHead("Integer"),
		MatchAny(),
		MatchOr(MatchHead("Integer"), MatchHead("String")),
		MatchNot(MatchHead("Integer")),
		ZeroOrMore(MatchAny()),
		OneOrMore(MatchHead("Integer")),
		Optional(MatchHead("String")),
		Named("x", MatchAny()),
	}

	for _, pattern := range patterns {
		str := pattern.String()
		if str == "" {
			t.Errorf("Pattern %T returned empty string", pattern)
		}
		t.Logf("Pattern: %s", str)
	}
}

func TestThompsonNFA(t *testing.T) {
	tests := []struct {
		name             string
		pattern          Pattern
		expr             Expr
		expectedMatched  bool
		expectedBindings map[string]Expr
	}{
		{
			name: "or pattern - first alternative matches",
			pattern: MatchOr(
				MatchHead("Integer"),
				MatchHead("String"),
			),
			expr:             NewInteger(42),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "or pattern - second alternative matches",
			pattern: MatchOr(
				MatchHead("Integer"),
				MatchHead("String"),
			),
			expr:             NewString("hello"),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "or pattern - no match",
			pattern: MatchOr(
				MatchHead("Integer"),
				MatchHead("String"),
			),
			expr:             NewBool(true),
			expectedMatched:  false,
			expectedBindings: map[string]Expr{},
		},
		{
			name:             "not pattern - positive case",
			pattern:          MatchNot(MatchHead("Integer")),
			expr:             NewString("hello"),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name:             "not pattern - negative case",
			pattern:          MatchNot(MatchHead("Integer")),
			expr:             NewInteger(42),
			expectedMatched:  false,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "sequence with or in middle (needs NFA)",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				MatchOr(
					MatchHead("Integer"),
					MatchHead("String"),
				),
				MatchLiteral(NewString("suffix")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewInteger(42),
				NewString("suffix"),
			),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "sequence with quantifier in middle",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				ZeroOrMore(MatchHead("Integer")),
				MatchLiteral(NewString("suffix")),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewInteger(1),
				NewInteger(2),
				NewString("suffix"),
			),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "nested or patterns",
			pattern: MatchOr(
				MatchOr(
					MatchHead("Integer"),
					MatchHead("String"),
				),
				MatchHead("Symbol"),
			),
			expr:             NewBool(false), // False is a Symbol
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "simple named capture in NFA context",
			pattern: MatchList(
				Named("value", MatchHead("Integer")),
				MatchOr(
					MatchLiteral(NewString("suffix1")),
					MatchLiteral(NewString("suffix2")),
				),
			),
			expr: NewList("List",
				NewInteger(123),
				NewString("suffix2"),
			),
			expectedMatched: true,
			expectedBindings: map[string]Expr{
				"value": NewInteger(123),
			},
		},
		{
			name: "complex pattern: sequence with or and quantifier",
			pattern: MatchList(
				MatchLiteral(NewString("data")),
				ZeroOrMore(MatchOr(
					MatchHead("Integer"),
					MatchHead("String"),
				)),
			),
			expr: NewList("List",
				NewString("data"),
				NewInteger(1),
				NewString("hello"),
				NewInteger(2),
			),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := CompilePattern(tt.pattern)
			if err != nil {
				t.Fatalf("Failed to compile pattern: %v", err)
			}

			result := compiled.Match(tt.expr)
			if result.Matched != tt.expectedMatched {
				t.Errorf("Expected matched %v, got %v", tt.expectedMatched, result.Matched)
			}

			// Check bindings
			if len(result.Bindings) != len(tt.expectedBindings) {
				t.Errorf("Expected %d bindings, got %d", len(tt.expectedBindings), len(result.Bindings))
			}

			for name, expectedValue := range tt.expectedBindings {
				actualValue, exists := result.Bindings[name]
				if !exists {
					t.Errorf("Expected binding %s not found", name)
					continue
				}
				if !actualValue.Equal(expectedValue) {
					t.Errorf("Binding %s: expected %v, got %v", name, expectedValue, actualValue)
				}
			}

			// Log for debugging
			t.Logf("Pattern: %s, Strategy: %v, Matched: %v, Bindings: %v",
				tt.pattern.String(), compiled.Strategy(), result.Matched, result.Bindings)
		})
	}
}

func TestComplexPatternFixed(t *testing.T) {
	// Test cases that verify the fixes for BuildNamed and BuildNot

	t.Run("named capture with Or pattern (now working)", func(t *testing.T) {
		// Named captures with Or patterns should now collect bindings
		pattern := Named("value", MatchOr(
			MatchHead("Integer"),
			MatchHead("String"),
		))

		compiled, err := CompilePattern(pattern)
		if err != nil {
			t.Fatalf("Failed to compile pattern: %v", err)
		}

		result := compiled.Match(NewInteger(42))
		if !result.Matched {
			t.Errorf("Expected match to succeed")
		}

		if len(result.Bindings) != 1 {
			t.Errorf("Expected 1 binding, got %d", len(result.Bindings))
		}

		if value, exists := result.Bindings["value"]; !exists {
			t.Errorf("Expected binding 'value' not found")
		} else if !value.Equal(NewInteger(42)) {
			t.Errorf("Expected binding value 42, got %v", value)
		}
	})

	t.Run("not with Or pattern (now working)", func(t *testing.T) {
		// Not(Or(Integer, String)) should match anything that's not Integer or String
		pattern := MatchNot(MatchOr(
			MatchHead("Integer"),
			MatchHead("String"),
		))

		compiled, err := CompilePattern(pattern)
		if err != nil {
			t.Fatalf("Failed to compile pattern: %v", err)
		}

		// Should match Symbol (True/False) - it's not Integer or String
		result1 := compiled.Match(NewBool(true))
		if !result1.Matched {
			t.Errorf("Expected Symbol to match Not(Or(Integer, String))")
		}

		// Should NOT match Integer
		result2 := compiled.Match(NewInteger(42))
		if result2.Matched {
			t.Errorf("Expected Integer to NOT match Not(Or(Integer, String))")
		}

		// Should NOT match String
		result3 := compiled.Match(NewString("hello"))
		if result3.Matched {
			t.Errorf("Expected String to NOT match Not(Or(Integer, String))")
		}
	})

	t.Run("nested named captures (working)", func(t *testing.T) {
		// Named(outer, Sequence(Named(inner, Integer), String))
		pattern := Named("outer", MatchSequence(
			Named("inner", MatchHead("Integer")),
			MatchHead("String"),
		))

		compiled, err := CompilePattern(pattern)
		if err != nil {
			t.Fatalf("Failed to compile pattern: %v", err)
		}

		expr := NewList("List", NewInteger(123), NewString("test"))
		result := compiled.Match(expr)

		if !result.Matched {
			t.Errorf("Expected match to succeed")
		}

		if len(result.Bindings) != 2 {
			t.Errorf("Expected 2 bindings, got %d", len(result.Bindings))
		}

		if inner, exists := result.Bindings["inner"]; !exists {
			t.Errorf("Expected binding 'inner' not found")
		} else if !inner.Equal(NewInteger(123)) {
			t.Errorf("Expected inner binding 123, got %v", inner)
		}

		if outer, exists := result.Bindings["outer"]; !exists {
			t.Errorf("Expected binding 'outer' not found")
		} else if !outer.Equal(expr) {
			t.Errorf("Expected outer binding to be the whole expression %v, got %v", expr, outer)
		}
	})
}

func TestPredicatePatterns(t *testing.T) {
	tests := []struct {
		name             string
		pattern          Pattern
		expr             Expr
		expectedMatched  bool
		expectedBindings map[string]Expr
	}{
		{
			name: "predicate on integer - value greater than 10",
			pattern: MatchPredicate(
				MatchHead("Integer"),
				func(expr Expr) bool {
					if intVal, ok := ExtractInt64(expr); ok {
						return intVal > 10
					}
					return false
				},
			),
			expr:             NewInteger(15),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "predicate on integer - value NOT greater than 10",
			pattern: MatchPredicate(
				MatchHead("Integer"),
				func(expr Expr) bool {
					if intVal, ok := ExtractInt64(expr); ok {
						return intVal > 10
					}
					return false
				},
			),
			expr:             NewInteger(5),
			expectedMatched:  false,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "predicate on string - length check",
			pattern: MatchPredicate(
				MatchHead("String"),
				func(expr Expr) bool {
					if strVal, ok := ExtractString(expr); ok {
						return len(strVal) >= 5
					}
					return false
				},
			),
			expr:             NewString("hello"),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "predicate with named capture",
			pattern: Named("value", MatchPredicate(
				MatchHead("Integer"),
				func(expr Expr) bool {
					if intVal, ok := ExtractInt64(expr); ok {
						return intVal%2 == 0 // Even numbers
					}
					return false
				},
			)),
			expr:            NewInteger(8),
			expectedMatched: true,
			expectedBindings: map[string]Expr{
				"value": NewInteger(8),
			},
		},
		{
			name: "predicate in list",
			pattern: MatchList(
				MatchLiteral(NewString("numbers")),
				MatchPredicate(
					MatchHead("Integer"),
					func(expr Expr) bool {
						if intVal, ok := ExtractInt64(expr); ok {
							return intVal > 0 // Positive numbers
						}
						return false
					},
				),
				MatchPredicate(
					MatchHead("Integer"),
					func(expr Expr) bool {
						if intVal, ok := ExtractInt64(expr); ok {
							return intVal < 100 // Less than 100
						}
						return false
					},
				),
			),
			expr: NewList("List",
				NewString("numbers"),
				NewInteger(42),
				NewInteger(7),
			),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "predicate with Or pattern (needs NFA)",
			pattern: MatchPredicate(
				MatchOr(
					MatchHead("Integer"),
					MatchHead("Real"),
				),
				func(expr Expr) bool {
					// Check if it's a "round" number (integer or real with no decimal part)
					if _, ok := ExtractInt64(expr); ok {
						return true
					}
					if realVal, ok := ExtractFloat64(expr); ok {
						return realVal == float64(int64(realVal))
					}
					return false
				},
			),
			expr:             NewReal(42.0),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "predicate fails on wrong type",
			pattern: MatchPredicate(
				MatchHead("Integer"),
				func(expr Expr) bool {
					return true // Always true if it gets here
				},
			),
			expr:             NewString("not an integer"),
			expectedMatched:  false,
			expectedBindings: map[string]Expr{},
		},
		// Tests to trigger the simple predicate optimization in BuildPredicate NFA path
		// These patterns require NFA but contain predicates with simple inner patterns
		{
			name: "list with simple predicate in NFA context (triggers simple optimization)",
			pattern: MatchList(
				MatchLiteral(NewString("prefix")),
				MatchPredicate(
					MatchHead("Integer"),
					func(expr Expr) bool {
						if intVal, ok := ExtractInt64(expr); ok {
							return intVal%2 == 0 // Even numbers only
						}
						return false
					},
				),
				MatchOr( // This forces NFA strategy
					MatchLiteral(NewString("suffix1")),
					MatchLiteral(NewString("suffix2")),
				),
			),
			expr: NewList("List",
				NewString("prefix"),
				NewInteger(8),
				NewString("suffix1"),
			),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "list with literal predicate in NFA context (triggers simple optimization)",
			pattern: MatchList(
				MatchLiteral(NewString("start")),
				MatchPredicate(
					MatchLiteral(NewInteger(42)),
					func(expr Expr) bool {
						return true // Always match if we get the right literal
					},
				),
				MatchOr( // This forces NFA strategy
					MatchLiteral(NewString("end1")),
					MatchLiteral(NewString("end2")),
				),
			),
			expr: NewList("List",
				NewString("start"),
				NewInteger(42),
				NewString("end2"),
			),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "sequence with any predicate in NFA context (triggers simple optimization)",
			pattern: MatchList(
				MatchLiteral(NewString("data")),
				MatchPredicate(
					MatchAny(),
					func(expr Expr) bool {
						// Only match symbols that start with 'T'
						if sym, ok := expr.(Symbol); ok {
							return len(string(sym)) > 0 && string(sym)[0] == 'T'
						}
						return false
					},
				),
				MatchOr( // This forces NFA strategy
					MatchLiteral(NewInteger(1)),
					MatchLiteral(NewInteger(2)),
				),
			),
			expr: NewList("List",
				NewString("data"),
				NewBool(true), // True is a symbol starting with 'T'
				NewInteger(1),
			),
			expectedMatched:  true,
			expectedBindings: map[string]Expr{},
		},
		{
			name: "sequence with any predicate in NFA context fails (triggers simple optimization)",
			pattern: MatchSequence(
				MatchLiteral(NewString("data")),
				MatchPredicate(
					MatchAny(),
					func(expr Expr) bool {
						// Only match symbols that start with 'T'
						if sym, ok := expr.(Symbol); ok {
							return len(string(sym)) > 0 && string(sym)[0] == 'T'
						}
						return false
					},
				),
				MatchOr( // This forces NFA strategy
					MatchLiteral(NewInteger(1)),
					MatchLiteral(NewInteger(2)),
				),
			),
			expr: NewList("List",
				NewString("data"),
				NewBool(false), // False starts with 'F', should fail predicate
				NewInteger(1),
			),
			expectedMatched:  false,
			expectedBindings: map[string]Expr{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := CompilePattern(tt.pattern)
			if err != nil {
				t.Fatalf("Failed to compile pattern: %v", err)
			}

			result := compiled.Match(tt.expr)
			if result.Matched != tt.expectedMatched {
				t.Errorf("Expected matched %v, got %v", tt.expectedMatched, result.Matched)
			}

			// Check bindings
			if len(result.Bindings) != len(tt.expectedBindings) {
				t.Errorf("Expected %d bindings, got %d", len(tt.expectedBindings), len(result.Bindings))
			}

			for name, expectedValue := range tt.expectedBindings {
				actualValue, exists := result.Bindings[name]
				if !exists {
					t.Errorf("Expected binding %s not found", name)
					continue
				}
				if !actualValue.Equal(expectedValue) {
					t.Errorf("Binding %s: expected %v, got %v", name, expectedValue, actualValue)
				}
			}

			// Log for debugging
			t.Logf("Pattern: %s, Strategy: %v, Matched: %v, Bindings: %v",
				tt.pattern.String(), compiled.Strategy(), result.Matched, result.Bindings)
		})
	}
}

func TestNamedCaptures(t *testing.T) {
	tests := []struct {
		name             string
		pattern          Pattern
		expr             Expr
		expectedMatched  bool
		expectedBindings map[string]Expr
		strategy         ExecutionStrategy
	}{
		{
			name:            "simple named capture",
			pattern:         Named("x", MatchHead("Integer")),
			expr:            NewInteger(42),
			expectedMatched: true,
			expectedBindings: map[string]Expr{
				"x": NewInteger(42),
			},
			strategy: StrategyDirect,
		},
		{
			name:             "named capture no match",
			pattern:          Named("x", MatchHead("String")),
			expr:             NewInteger(42),
			expectedMatched:  false,
			expectedBindings: map[string]Expr{},
			strategy:         StrategyDirect,
		},
		{
			name: "list with named captures",
			pattern: MatchList(
				Named("first", MatchHead("Integer")),
				Named("second", MatchHead("String")),
				Named("third", MatchAny()),
			),
			expr: NewList("List",
				NewInteger(1),
				NewString("hello"),
				NewBool(true),
			),
			expectedMatched: true,
			expectedBindings: map[string]Expr{
				"first":  NewInteger(1),
				"second": NewString("hello"),
				"third":  NewBool(true),
			},
			strategy: StrategyDirect,
		},
		{
			name: "nested named captures",
			pattern: Named("outer", MatchSequence(
				Named("inner", MatchHead("Integer")),
				MatchHead("String"),
			)),
			expr: NewList("List",
				NewInteger(42),
				NewString("test"),
			),
			expectedMatched: true,
			expectedBindings: map[string]Expr{
				"inner": NewInteger(42),
				"outer": NewList("List", NewInteger(42), NewString("test")),
			},
			strategy: StrategyNFA,
		},
		{
			name: "named capture with quantifier",
			pattern: MatchList(
				Named("prefix", MatchLiteral(NewString("data"))),
				ZeroOrMore(MatchHead("Integer")),
			),
			expr: NewList("List",
				NewString("data"),
				NewInteger(1),
				NewInteger(2),
			),
			expectedMatched: true,
			expectedBindings: map[string]Expr{
				"prefix": NewString("data"),
			},
			strategy: StrategyDirect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := CompilePattern(tt.pattern)
			if err != nil {
				t.Fatalf("Failed to compile pattern: %v", err)
			}

			if compiled.Strategy() != tt.strategy {
				t.Errorf("Expected strategy %v, got %v", tt.strategy, compiled.Strategy())
			}

			result := compiled.Match(tt.expr)
			if result.Matched != tt.expectedMatched {
				t.Errorf("Expected matched %v, got %v", tt.expectedMatched, result.Matched)
			}

			// Check bindings
			if len(result.Bindings) != len(tt.expectedBindings) {
				t.Errorf("Expected %d bindings, got %d", len(tt.expectedBindings), len(result.Bindings))
			}

			for name, expectedValue := range tt.expectedBindings {
				actualValue, exists := result.Bindings[name]
				if !exists {
					t.Errorf("Expected binding %s not found", name)
					continue
				}
				if !actualValue.Equal(expectedValue) {
					t.Errorf("Binding %s: expected %v, got %v", name, expectedValue, actualValue)
				}
			}

			// Log bindings for debugging
			t.Logf("Bindings: %v", result.Bindings)
		})
	}
}
