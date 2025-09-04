package core

import (
	"reflect"
	"testing"
)

func TestIsSymbolicPattern(t *testing.T) {
	p := MustParse("Pattern(x, 1)")
	ok, _, _ := IsSymbolicPattern(p)
	if !ok {
		t.Errorf("expected smbolic pattern")
	}

	p = MustParse("Just(x, 1)")
	ok, _, _ = IsSymbolicPattern(p)
	if ok {
		t.Errorf("Did not expected symbolic pattern")
	}

}

// Test pattern analysis functions
func TestIsSymbolicBlank(t *testing.T) {
	tests := []struct {
		expr           Expr
		expectBlank    bool
		expectType     PatternType
		expectTypeExpr Expr
	}{
		{ListFrom(symbolBlank), true, BlankPattern, nil},
		{ListFrom(symbolBlank, symbolInteger), true, BlankPattern, symbolInteger},
		{ListFrom(symbolBlankSequence), true, BlankSequencePattern, nil},
		{ListFrom(symbolBlankNullSequence, symbolString), true, BlankNullSequencePattern, symbolString},
		{NewSymbol("x"), false, PatternUnknown, nil},
		{NewInteger(42), false, PatternUnknown, nil},
	}

	for _, test := range tests {
		isBlank, blankType, typeExpr := IsSymbolicBlank(test.expr)
		if isBlank != test.expectBlank {
			t.Errorf("IsSymbolicBlank(%v) blank = %v, want %v", test.expr, isBlank, test.expectBlank)
		}
		if blankType != test.expectType {
			t.Errorf("IsSymbolicBlank(%v) type = %v, want %v", test.expr, blankType, test.expectType)
		}
		if !reflect.DeepEqual(typeExpr, test.expectTypeExpr) {
			t.Errorf("IsSymbolicBlank(%v) typeExpr = %v, want %v", test.expr, typeExpr, test.expectTypeExpr)
		}
	}
}

// Test type matching functions
func TestMatchesType(t *testing.T) {
	tests := []struct {
		expr     Expr
		typeName string
		expected bool
	}{
		{NewInteger(42), "Integer", true},
		{NewInteger(42), "Number", true},
		{NewInteger(42), "String", false},
		{NewReal(3.14), "Real", true},
		{NewReal(3.14), "Number", true},
		{NewString("hello"), "String", true},
		{NewSymbol("x"), "Symbol", true},
		{NewList(symbolList), "List", true},
		{NewInteger(42), "", true}, // No constraint
		{NewObjectExpr(NewSymbol("CustomType"), NewInteger(1)), "CustomType", true},
		{NewObjectExpr(NewSymbol("CustomType"), NewInteger(1)), "OtherType", false},
	}

	for _, test := range tests {
		result := MatchesType(test.expr, test.typeName)
		if result != test.expected {
			t.Errorf("MatchesType(%v, %q) = %v, want %v", test.expr, test.typeName, result, test.expected)
		}
	}
}

func TestIsBuiltinType(t *testing.T) {
	builtinTypes := []string{"Integer", "Real", "Number", "String", "Symbol", "List", "Rule", "ByteArray", "Association"}
	for _, typeName := range builtinTypes {
		if !IsBuiltinType(typeName) {
			t.Errorf("IsBuiltinType(%q) should return true", typeName)
		}
	}

	userTypes := []string{"CustomType", "MyClass", ""}
	for _, typeName := range userTypes {
		if IsBuiltinType(typeName) {
			t.Errorf("IsBuiltinType(%q) should return false", typeName)
		}
	}
}

// Test pattern specificity functions
func TestGetPatternSpecificity(t *testing.T) {
	tests := []struct {
		pattern  Expr
		expected PatternSpecificity
	}{
		{NewInteger(42), SpecificityLiteral * 100},
		{NewSymbol("x"), SpecificityLiteral * 100},
		{ListFrom(symbolBlank), SpecificityGeneral*10 + 2},
		{ListFrom(symbolBlank, symbolInteger), SpecificityBuiltinType*10 + 2},
		{ListFrom(symbolBlank, NewSymbol("Foo")), SpecificityUserType*10 + 2},
		{ListFrom(symbolPattern, NewSymbol("x"), ListFrom(symbolBlank)), SpecificityGeneral*10 + 2},
	}

	for _, test := range tests {
		result := GetPatternSpecificity(test.pattern)
		if result != test.expected {
			t.Errorf("GetPatternSpecificity(%v) = %v, want %v", test.pattern, result, test.expected)
		}
	}
}

// Test literal symbol pattern distinction
func TestLiteralSymbolPatterns(t *testing.T) {
	matcher := NewPatternMatcher()

	// These should match exactly
	if !matcher.TestMatch(NewSymbol("s1"), NewSymbol("s1")) {
		t.Error("s1 pattern should match s1 symbol")
	}

	// These should NOT match
	if matcher.TestMatch(NewSymbol("s1"), NewSymbol("s2")) {
		t.Error("s1 pattern should NOT match s2 symbol")
	}

	// Symbol pattern should NOT match integer
	if matcher.TestMatch(NewSymbol("s1"), NewInteger(100)) {
		t.Error("s1 pattern should NOT match integer 100")
	}

	// Symbol pattern should NOT match different types
	if matcher.TestMatch(NewSymbol("True"), NewSymbol("False")) {
		t.Error("True pattern should NOT match False symbol")
	}
}

// Test pure pattern matcher
func TestPatternMatcher(t *testing.T) {
	matcher := NewPatternMatcher()

	tests := []struct {
		pattern  Expr
		expr     Expr
		expected bool
	}{
		// Literal matches
		{NewInteger(42), NewInteger(42), true},
		{NewInteger(42), NewInteger(43), false},
		{NewSymbol("x"), NewSymbol("x"), true},
		{NewSymbol("x"), NewSymbol("y"), false},

		// Blank patterns
		{ListFrom(symbolBlank), NewInteger(42), true},
		{ListFrom(symbolBlank), NewString("hello"), true},
		{ListFrom(symbolBlank, symbolInteger), NewInteger(42), true},
		{ListFrom(symbolBlank, symbolInteger), NewString("hello"), false},

		// Pattern expressions
		{ListFrom(symbolPattern, NewSymbol("x"), ListFrom(symbolBlank)), NewInteger(42), true},
		{ListFrom(symbolPattern, NewSymbol("x"), ListFrom(symbolBlank, symbolString)), NewString("hello"), true},
		{ListFrom(symbolPattern, NewSymbol("x"), ListFrom(symbolBlank, symbolString)), NewInteger(42), false},

		// Empty List patterns
		{ListFrom(symbolList), ListFrom(symbolList), true},
		{ListFrom(NewSymbol("Foo")), ListFrom(NewSymbol("Foo")), true},
		{ListFrom(symbolList), ListFrom(NewSymbol("Foo")), false},

		// List patterns
		{ListFrom(symbolPlus, ListFrom(symbolBlank), ListFrom(symbolBlank)),
			ListFrom(symbolPlus, NewInteger(1), NewInteger(2)), true},
		{ListFrom(symbolPlus, ListFrom(symbolBlank), ListFrom(symbolBlank)),
			ListFrom(symbolTimes, NewInteger(1), NewInteger(2)), false},

		// Alternatives, single
		{ListFrom(symbolAlternatives, ListFrom(symbolBlank, symbolInteger), ListFrom(symbolBlank, symbolReal)),
			NewInteger(2), true},
		{ListFrom(symbolAlternatives, ListFrom(symbolBlank, symbolReal), ListFrom(symbolBlank, symbolInteger)),
			NewInteger(2), true},
		{ListFrom(symbolAlternatives, ListFrom(symbolBlank, symbolReal), ListFrom(symbolBlank, symbolReal)),
			NewString("2"), false},

		// Alternatives with Binding, single
		{ListFrom(symbolAlternatives,
			ListFrom(symbolPattern, NewSymbol("x"), ListFrom(symbolBlank, symbolInteger)),
			ListFrom(symbolPattern, NewSymbol("x"), ListFrom(symbolBlank, symbolReal))),
			NewInteger(2), true},
		{ListFrom(symbolAlternatives,
			ListFrom(symbolPattern, NewSymbol("x"), ListFrom(symbolBlank, symbolReal)),
			ListFrom(symbolPattern, NewSymbol("x"), ListFrom(symbolBlank, symbolInteger))),
			NewInteger(2), true},
		{ListFrom(symbolAlternatives,
			ListFrom(symbolPattern, NewSymbol("x"), ListFrom(symbolBlank, symbolReal)),
			ListFrom(symbolPattern, NewSymbol("x"), ListFrom(symbolBlank, symbolInteger))),
			NewString("2"), false},

		// List
		{ListFrom(symbolList,
			ListFrom(symbolAlternatives,
				ListFrom(symbolBlank, symbolInteger),
				ListFrom(symbolBlank, symbolReal)),
			NewString("foo")),
			ListFrom(symbolList, NewInteger(2), NewString("foo")), true},
		{ListFrom(symbolList,
			ListFrom(symbolAlternatives,
				ListFrom(symbolBlank, symbolInteger),
				ListFrom(symbolBlank, symbolReal)),
			NewString("junk")),
			ListFrom(symbolList, NewInteger(2), NewString("foo")), false},
	}

	for _, test := range tests {
		result := matcher.TestMatch(test.expr, test.pattern)
		if result != test.expected {
			t.Errorf("TestMatch(%v, %v) = %v, want %v", test.expr, test.pattern, result, test.expected)
		}
	}
}
