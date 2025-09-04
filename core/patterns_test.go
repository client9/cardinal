package core

import (
	"reflect"
	"testing"

	"github.com/client9/sexpr/core/symbol"
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
		t.Errorf("Did not expected symbol.ic pattern")
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
		{ListFrom(symbol.Blank), true, BlankPattern, nil},
		{ListFrom(symbol.Blank, symbol.Integer), true, BlankPattern, symbol.Integer},
		{ListFrom(symbol.BlankSequence), true, BlankSequencePattern, nil},
		{ListFrom(symbol.BlankNullSequence, symbol.String), true, BlankNullSequencePattern, symbol.String},
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
		{NewList(symbol.List), "List", true},
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
		{ListFrom(symbol.Blank), SpecificityGeneral*10 + 2},
		{ListFrom(symbol.Blank, symbol.Integer), SpecificityBuiltinType*10 + 2},
		{ListFrom(symbol.Blank, NewSymbol("Foo")), SpecificityUserType*10 + 2},
		{ListFrom(symbol.Pattern, NewSymbol("x"), ListFrom(symbol.Blank)), SpecificityGeneral*10 + 2},
	}

	for _, test := range tests {
		result := GetPatternSpecificity(test.pattern)
		if result != test.expected {
			t.Errorf("GetPatternSpecificity(%v) = %v, want %v", test.pattern, result, test.expected)
		}
	}
}

// Test literal symbol. pattern distinction
func TestLiteralSymbolPatterns(t *testing.T) {
	matcher := NewPatternMatcher()

	// These should match exactly
	if !matcher.TestMatch(NewSymbol("s1"), NewSymbol("s1")) {
		t.Error("s1 pattern should match s1 symbol.")
	}

	// These should NOT match
	if matcher.TestMatch(NewSymbol("s1"), NewSymbol("s2")) {
		t.Error("s1 pattern should NOT match s2 symbol.")
	}

	// Symbol pattern should NOT match integer
	if matcher.TestMatch(NewSymbol("s1"), NewInteger(100)) {
		t.Error("s1 pattern should NOT match integer 100")
	}

	// Symbol pattern should NOT match different types
	if matcher.TestMatch(NewSymbol("True"), NewSymbol("False")) {
		t.Error("True pattern should NOT match False symbol.")
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
		{ListFrom(symbol.Blank), NewInteger(42), true},
		{ListFrom(symbol.Blank), NewString("hello"), true},
		{ListFrom(symbol.Blank, symbol.Integer), NewInteger(42), true},
		{ListFrom(symbol.Blank, symbol.Integer), NewString("hello"), false},

		// Pattern expressions
		{ListFrom(symbol.Pattern, NewSymbol("x"), ListFrom(symbol.Blank)), NewInteger(42), true},
		{ListFrom(symbol.Pattern, NewSymbol("x"), ListFrom(symbol.Blank, symbol.String)), NewString("hello"), true},
		{ListFrom(symbol.Pattern, NewSymbol("x"), ListFrom(symbol.Blank, symbol.String)), NewInteger(42), false},

		// Empty List patterns
		{ListFrom(symbol.List), ListFrom(symbol.List), true},
		{ListFrom(NewSymbol("Foo")), ListFrom(NewSymbol("Foo")), true},
		{ListFrom(symbol.List), ListFrom(NewSymbol("Foo")), false},

		// List patterns
		{ListFrom(symbol.Plus, ListFrom(symbol.Blank), ListFrom(symbol.Blank)),
			ListFrom(symbol.Plus, NewInteger(1), NewInteger(2)), true},
		{ListFrom(symbol.Plus, ListFrom(symbol.Blank), ListFrom(symbol.Blank)),
			ListFrom(symbol.Times, NewInteger(1), NewInteger(2)), false},

		// Alternatives, single
		{ListFrom(symbol.Alternatives, ListFrom(symbol.Blank, symbol.Integer), ListFrom(symbol.Blank, symbol.Real)),
			NewInteger(2), true},
		{ListFrom(symbol.Alternatives, ListFrom(symbol.Blank, symbol.Real), ListFrom(symbol.Blank, symbol.Integer)),
			NewInteger(2), true},
		{ListFrom(symbol.Alternatives, ListFrom(symbol.Blank, symbol.Real), ListFrom(symbol.Blank, symbol.Real)),
			NewString("2"), false},

		// Alternatives with Binding, single
		{ListFrom(symbol.Alternatives,
			ListFrom(symbol.Pattern, NewSymbol("x"), ListFrom(symbol.Blank, symbol.Integer)),
			ListFrom(symbol.Pattern, NewSymbol("x"), ListFrom(symbol.Blank, symbol.Real))),
			NewInteger(2), true},
		{ListFrom(symbol.Alternatives,
			ListFrom(symbol.Pattern, NewSymbol("x"), ListFrom(symbol.Blank, symbol.Real)),
			ListFrom(symbol.Pattern, NewSymbol("x"), ListFrom(symbol.Blank, symbol.Integer))),
			NewInteger(2), true},
		{ListFrom(symbol.Alternatives,
			ListFrom(symbol.Pattern, NewSymbol("x"), ListFrom(symbol.Blank, symbol.Real)),
			ListFrom(symbol.Pattern, NewSymbol("x"), ListFrom(symbol.Blank, symbol.Integer))),
			NewString("2"), false},

		// List
		{ListFrom(symbol.List,
			ListFrom(symbol.Alternatives,
				ListFrom(symbol.Blank, symbol.Integer),
				ListFrom(symbol.Blank, symbol.Real)),
			NewString("foo")),
			ListFrom(symbol.List, NewInteger(2), NewString("foo")), true},
		{ListFrom(symbol.List,
			ListFrom(symbol.Alternatives,
				ListFrom(symbol.Blank, symbol.Integer),
				ListFrom(symbol.Blank, symbol.Real)),
			NewString("junk")),
			ListFrom(symbol.List, NewInteger(2), NewString("foo")), false},
	}

	for _, test := range tests {
		result := matcher.TestMatch(test.expr, test.pattern)
		if result != test.expected {
			t.Errorf("TestMatch(%v, %v) = %v, want %v", test.expr, test.pattern, result, test.expected)
		}
	}
}
