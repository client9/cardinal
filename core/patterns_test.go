package core

import (
	"reflect"
	"testing"
)

// Test pattern analysis functions
func TestIsSymbolicBlank(t *testing.T) {
	tests := []struct {
		expr           Expr
		expectBlank    bool
		expectType     PatternType
		expectTypeExpr Expr
	}{
		{CreateBlankExpr(nil), true, BlankPattern, nil},
		{CreateBlankExpr(NewSymbol("Integer")), true, BlankPattern, NewSymbol("Integer")},
		{CreateBlankSequenceExpr(nil), true, BlankSequencePattern, nil},
		{CreateBlankNullSequenceExpr(NewSymbol("String")), true, BlankNullSequencePattern, NewSymbol("String")},
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

func TestIsSymbolicPattern(t *testing.T) {
	nameExpr := NewSymbol("x")
	blankExpr := CreateBlankExpr(nil)
	pattern := CreatePatternExpr(nameExpr, blankExpr)

	// Test valid pattern
	isPattern, gotName, gotBlank := IsSymbolicPattern(pattern)
	if !isPattern {
		t.Error("IsSymbolicPattern should return true for Pattern expression")
	}
	if !reflect.DeepEqual(gotName, nameExpr) {
		t.Errorf("IsSymbolicPattern name = %v, want %v", gotName, nameExpr)
	}
	if !reflect.DeepEqual(gotBlank, blankExpr) {
		t.Errorf("IsSymbolicPattern blank = %v, want %v", gotBlank, blankExpr)
	}

	// Test non-pattern
	isPattern, _, _ = IsSymbolicPattern(NewSymbol("x"))
	if isPattern {
		t.Error("IsSymbolicPattern should return false for non-Pattern expression")
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
		{NewList("List"), "List", true},
		{NewInteger(42), "", true}, // No constraint
		{NewObjectExpr("CustomType", NewInteger(1)), "CustomType", true},
		{NewObjectExpr("CustomType", NewInteger(1)), "OtherType", false},
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
		{CreateBlankExpr(nil), SpecificityGeneral*10 + 2},
		{CreateBlankExpr(NewSymbol("Integer")), SpecificityBuiltinType*10 + 2},
		{CreateBlankExpr(NewSymbol("CustomType")), SpecificityUserType*10 + 2},
		{CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(nil)), SpecificityGeneral*10 + 2},
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
		{CreateBlankExpr(nil), NewInteger(42), true},
		{CreateBlankExpr(nil), NewString("hello"), true},
		{CreateBlankExpr(NewSymbol("Integer")), NewInteger(42), true},
		{CreateBlankExpr(NewSymbol("Integer")), NewString("hello"), false},

		// Pattern expressions
		{CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(nil)), NewInteger(42), true},
		{CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(NewSymbol("String"))), NewString("hello"), true},
		{CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(NewSymbol("String"))), NewInteger(42), false},

		// List patterns
		{NewList("Plus", CreateBlankExpr(nil), CreateBlankExpr(nil)),
			NewList("Plus", NewInteger(1), NewInteger(2)), true},
		{NewList("Plus", CreateBlankExpr(nil), CreateBlankExpr(nil)),
			NewList("Times", NewInteger(1), NewInteger(2)), false},

		// Alternatives, single
		{NewList("Alternatives", CreateBlankExpr(NewSymbol("Integer")), CreateBlankExpr(NewSymbol("Real"))),
			NewInteger(2), true},
		{NewList("Alternatives", CreateBlankExpr(NewSymbol("Real")), CreateBlankExpr(NewSymbol("Integer"))),
			NewInteger(2), true},
		{NewList("Alternatives", CreateBlankExpr(NewSymbol("Real")), CreateBlankExpr(NewSymbol("Integer"))),
			NewString("2"), false},

		// Alternatives with Binding, single
		{NewList("Alternatives",
			CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(NewSymbol("Integer"))),
			CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(NewSymbol("Real")))),
			NewInteger(2), true},
		{NewList("Alternatives",
			CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(NewSymbol("Real"))),
			CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(NewSymbol("Integer")))),
			NewInteger(2), true},
		{NewList("Alternatives",
			CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(NewSymbol("Real"))),
			CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(NewSymbol("Integer")))),
			NewString("2"), false},

		// List
		{NewList("List",
			NewList("Alternatives",
				CreateBlankExpr(NewSymbol("Integer")),
				CreateBlankExpr(NewSymbol("Real"))),
			NewString("foo")),
			NewList("List", NewInteger(2), NewString("foo")), true},
		{NewList("List",
			NewList("Alternatives",
				CreateBlankExpr(NewSymbol("Integer")),
				CreateBlankExpr(NewSymbol("Real"))),
			NewString("junk")),
			NewList("List", NewInteger(2), NewString("foo")), false},
	}

	for _, test := range tests {
		result := matcher.TestMatch(test.pattern, test.expr)
		if result != test.expected {
			t.Errorf("TestMatch(%v, %v) = %v, want %v", test.pattern, test.expr, result, test.expected)
		}
	}
}
