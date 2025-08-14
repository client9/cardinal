package core

import (
	"fmt"
	"testing"
)

// TestExampleUsage demonstrates the s-expression regex system with practical examples
func TestExampleUsage(t *testing.T) {
	t.Skip()
	// Example 1: Match a list with specific structure
	// Pattern: List(1, MatchHead("String"), MatchAny())
	// Should match: [1, "hello", anything]
	pattern1 := MatchSequence(
		MatchLiteral(NewInteger(1)),
		MatchHead("String"),
		MatchAny(),
	)

	expr1 := NewList("List",
		NewInteger(1),
		NewString("hello"),
		NewBool(true),
	)

	compiled1, _ := CompilePattern(pattern1)
	result1 := compiled1.Match(expr1)

	fmt.Printf("Example 1: %v (Strategy: %v)\n", result1.Matched, compiled1.Strategy())

	// Example 2: Match a function call with variable arguments
	// Pattern: List("Plus", MatchHead("Integer")*)
	// Should match: Plus(1, 2, 3, 4, ...)
	pattern2 := MatchSequence(
		MatchLiteral(NewString("prefix")),
		ZeroOrMore(MatchHead("Integer")),
	)

	expr2 := NewList("List",
		NewString("prefix"),
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	)

	compiled2, _ := CompilePattern(pattern2)
	result2 := compiled2.Match(expr2)

	fmt.Printf("Example 2: %v (Strategy: %v)\n", result2.Matched, compiled2.Strategy())

	// Example 3: Match optional elements
	// Pattern: List("config", MatchHead("String"), MatchHead("Integer")?)
	// Should match: List("config", "name", 42) or List("config", "name")
	pattern3 := MatchSequence(
		MatchLiteral(NewString("config")),
		MatchHead("String"),
		Optional(MatchHead("Integer")),
	)

	expr3a := NewList("List",
		NewString("config"),
		NewString("name"),
		NewInteger(42),
	)

	expr3b := NewList("List",
		NewString("config"),
		NewString("name"),
	)

	compiled3, _ := CompilePattern(pattern3)
	result3a := compiled3.Match(expr3a)
	result3b := compiled3.Match(expr3b)

	fmt.Printf("Example 3a: %v (Strategy: %v)\n", result3a.Matched, compiled3.Strategy())
	fmt.Printf("Example 3b: %v\n", result3b.Matched)

	// All examples should match and use direct strategy
	if !result1.Matched || !result2.Matched || !result3a.Matched || !result3b.Matched {
		t.Error("Expected all examples to match")
	}

	// These simple patterns should use Direct strategy
	if compiled1.Strategy() != StrategyDirect || compiled2.Strategy() != StrategyDirect || compiled3.Strategy() != StrategyDirect {
		t.Error("Expected all examples to use Direct strategy (simple patterns use direct optimization)")
	}
}

// TestPredicateExamples demonstrates practical predicate usage scenarios
func TestPredicateExamples(t *testing.T) {
	t.Skip()
	// Example 1: Validate configuration entries
	// Pattern: List("config", String, Integer > 0, Optional(Boolean))
	configPattern := MatchSequence(
		MatchLiteral(NewString("config")),
		MatchHead("String"), // config name
		MatchPredicate(
			MatchHead("Integer"),
			func(expr Expr) bool {
				if val, ok := ExtractInt64(expr); ok {
					return val > 0 // positive values only
				}
				return false
			},
		),
		Optional(MatchHead("Symbol")), // optional boolean flag
	)

	validConfig := NewList("List",
		NewString("config"),
		NewString("timeout"),
		NewInteger(30),
		NewBool(true),
	)

	invalidConfig := NewList("List",
		NewString("config"),
		NewString("retries"),
		NewInteger(-1), // Invalid: negative number
	)

	compiledConfig, _ := CompilePattern(configPattern)

	result1 := compiledConfig.Match(validConfig)
	result2 := compiledConfig.Match(invalidConfig)

	fmt.Printf("Valid config matches: %v\n", result1.Matched)
	fmt.Printf("Invalid config matches: %v\n", result2.Matched)

	if !result1.Matched {
		t.Error("Expected valid config to match")
	}
	if result2.Matched {
		t.Error("Expected invalid config to not match")
	}

	// Example 2: Filter numeric data with business rules
	// Pattern: Named capture of numbers in range [1, 100] (uses NFA due to Or pattern)
	numberInRangePattern := Named("validNumber", MatchPredicate(
		MatchOr(MatchHead("Integer"), MatchHead("Real")),
		func(expr Expr) bool {
			if intVal, ok := ExtractInt64(expr); ok {
				return intVal >= 1 && intVal <= 100
			}
			if realVal, ok := ExtractFloat64(expr); ok {
				return realVal >= 1.0 && realVal <= 100.0
			}
			return false
		},
	))

	compiledRange, _ := CompilePattern(numberInRangePattern)

	testValues := []Expr{
		NewInteger(50),  // Valid
		NewReal(75.5),   // Valid
		NewInteger(0),   // Invalid: too small
		NewInteger(150), // Invalid: too large
		NewString("50"), // Invalid: wrong type
	}

	for i, val := range testValues {
		result := compiledRange.Match(val)
		fmt.Printf("Value %d (%v): matches=%v", i+1, val, result.Matched)
		if result.Matched {
			fmt.Printf(", captured=%v", result.Bindings["validNumber"])
		}
		fmt.Printf("\n")
	}

	// Example 3: Complex validation with multiple predicates
	// Pattern: Function call with valid function name and arguments
	functionCallPattern := MatchSequence(
		MatchPredicateNamed(
			MatchHead("Symbol"),
			func(expr Expr) bool {
				if symVal, ok := expr.(Symbol); ok {
					name := string(symVal)
					// Valid function names: start with letter, contain only alphanumeric
					if len(name) == 0 {
						return false
					}
					if !((name[0] >= 'A' && name[0] <= 'Z') || (name[0] >= 'a' && name[0] <= 'z')) {
						return false
					}
					for _, c := range name[1:] {
						if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
							return false
						}
					}
					return true
				}
				return false
			},
			"ValidFunctionName",
		),
		ZeroOrMore(MatchAny()), // Arguments can be anything
	)

	validCall := NewList("List",
		NewSymbol("Plus"),
		NewInteger(1),
		NewInteger(2),
	)

	invalidCall := NewList("List",
		NewSymbol("123Invalid"), // Invalid function name
		NewInteger(1),
	)

	compiledFunc, _ := CompilePattern(functionCallPattern)

	result3 := compiledFunc.Match(validCall)
	result4 := compiledFunc.Match(invalidCall)

	fmt.Printf("Valid function call matches: %v\n", result3.Matched)
	fmt.Printf("Invalid function call matches: %v\n", result4.Matched)

	if !result3.Matched {
		t.Error("Expected valid function call to match")
	}
	if result4.Matched {
		t.Error("Expected invalid function call to not match")
	}
}

// TestPerformanceComparison shows the difference between strategies
func TestPerformanceComparison(t *testing.T) {
	t.Skip()
	// Fast path: simple patterns
	fastPattern := MatchSequence(
		MatchLiteral(NewString("data")),
		ZeroOrMore(MatchHead("Integer")),
	)

	// Complex pattern that would need NFA (when implemented)
	complexPattern := MatchOr(
		MatchHead("Integer"),
		MatchHead("String"),
	)

	compiledFast, _ := CompilePattern(fastPattern)
	compiledComplex, _ := CompilePattern(complexPattern)

	fmt.Printf("Fast pattern strategy: %v\n", compiledFast.Strategy())
	fmt.Printf("Complex pattern strategy: %v\n", compiledComplex.Strategy())

	if compiledFast.Strategy() != StrategyDirect {
		t.Error("Expected fast pattern to use direct strategy")
	}

	if compiledComplex.Strategy() != StrategyNFA {
		t.Error("Expected complex pattern to use NFA strategy")
	}
}

// TestNamedCaptureExamples demonstrates named captures with practical examples
func TestNamedCaptureExamples(t *testing.T) {
	t.Skip()
	// Example 1: Extract function name and arguments
	// Pattern: Simple Named patterns with MatchHead and MatchAny (uses Direct strategy)
	pattern1 := MatchSequence(
		Named("funcName", MatchHead("Symbol")),
		Named("firstArg", MatchAny()),
		ZeroOrMore(MatchAny()),
	)

	expr1 := NewList("List",
		NewSymbol("Plus"),
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	)

	compiled1, _ := CompilePattern(pattern1)
	result1 := compiled1.Match(expr1)

	fmt.Printf("Function extraction: %v\n", result1.Matched)
	fmt.Printf("Function name: %v\n", result1.Bindings["funcName"])
	fmt.Printf("First argument: %v\n", result1.Bindings["firstArg"])

	// Example 2: Configuration parsing
	// Pattern: Simple Named patterns (uses Direct strategy)
	pattern2 := MatchSequence(
		MatchLiteral(NewString("config")),
		Named("name", MatchHead("String")),
		Named("value", MatchAny()),
	)

	expr2 := NewList("List",
		NewString("config"),
		NewString("timeout"),
		NewInteger(30),
	)

	compiled2, _ := CompilePattern(pattern2)
	result2 := compiled2.Match(expr2)

	fmt.Printf("Config parsing: %v\n", result2.Matched)
	fmt.Printf("Config name: %v\n", result2.Bindings["name"])
	fmt.Printf("Config value: %v\n", result2.Bindings["value"])

	// Example 3: Nested structure extraction
	// Extract data from nested structures (complex Named with nested sequence uses NFA)
	pattern3 := Named("response", MatchSequence(
		MatchLiteral(NewString("response")),
		Named("status", MatchHead("Integer")),
		Named("data", MatchAny()),
	))

	expr3 := NewList("List",
		NewString("response"),
		NewInteger(200),
		NewList("List",
			NewString("user"),
			NewString("john"),
		),
	)

	compiled3, _ := CompilePattern(pattern3)
	result3 := compiled3.Match(expr3)

	fmt.Printf("Response parsing: %v\n", result3.Matched)
	fmt.Printf("Status: %v\n", result3.Bindings["status"])
	fmt.Printf("Data: %v\n", result3.Bindings["data"])
	fmt.Printf("Full response: %v\n", result3.Bindings["response"])

	// All examples should match
	if !result1.Matched || !result2.Matched || !result3.Matched {
		t.Error("Expected all examples to match")
	}

	// Check strategies: simple Named patterns use Direct, complex ones use NFA
	if compiled1.Strategy() != StrategyDirect {
		t.Errorf("Expected pattern1 to use Direct strategy (simple Named), got %v", compiled1.Strategy())
	}
	if compiled2.Strategy() != StrategyDirect {
		t.Errorf("Expected pattern2 to use Direct strategy (simple Named), got %v", compiled2.Strategy())
	}
	if compiled3.Strategy() != StrategyNFA {
		t.Errorf("Expected pattern3 to use NFA strategy (complex Named), got %v", compiled3.Strategy())
	}
}
