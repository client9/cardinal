package core

import (
	"reflect"
	"testing"
)

// Test pattern constructor functions
func TestCreateBlankExpr(t *testing.T) {
	// Test blank without type constraint
	blank := CreateBlankExpr(nil)
	expected := List{Elements: []Expr{NewSymbol("Blank")}}
	if !reflect.DeepEqual(blank, expected) {
		t.Errorf("CreateBlankExpr(nil) = %v, want %v", blank, expected)
	}

	// Test blank with type constraint
	typeExpr := NewSymbol("Integer")
	blankTyped := CreateBlankExpr(typeExpr)
	expectedTyped := List{Elements: []Expr{NewSymbol("Blank"), typeExpr}}
	if !reflect.DeepEqual(blankTyped, expectedTyped) {
		t.Errorf("CreateBlankExpr(Integer) = %v, want %v", blankTyped, expectedTyped)
	}
}

func TestCreateBlankSequenceExpr(t *testing.T) {
	// Test sequence without type constraint
	seq := CreateBlankSequenceExpr(nil)
	expected := List{Elements: []Expr{NewSymbol("BlankSequence")}}
	if !reflect.DeepEqual(seq, expected) {
		t.Errorf("CreateBlankSequenceExpr(nil) = %v, want %v", seq, expected)
	}

	// Test sequence with type constraint
	typeExpr := NewSymbol("String")
	seqTyped := CreateBlankSequenceExpr(typeExpr)
	expectedTyped := List{Elements: []Expr{NewSymbol("BlankSequence"), typeExpr}}
	if !reflect.DeepEqual(seqTyped, expectedTyped) {
		t.Errorf("CreateBlankSequenceExpr(String) = %v, want %v", seqTyped, expectedTyped)
	}
}

func TestCreatePatternExpr(t *testing.T) {
	nameExpr := NewSymbol("x")
	blankExpr := CreateBlankExpr(nil)
	pattern := CreatePatternExpr(nameExpr, blankExpr)

	expected := List{Elements: []Expr{NewSymbol("Pattern"), nameExpr, blankExpr}}
	if !reflect.DeepEqual(pattern, expected) {
		t.Errorf("CreatePatternExpr = %v, want %v", pattern, expected)
	}
}

// Test pattern analysis functions
func TestIsSymbolicBlank(t *testing.T) {
	tests := []struct {
		expr           Expr
		expectBlank    bool
		expectType     string
		expectTypeExpr Expr
	}{
		{CreateBlankExpr(nil), true, "Blank", nil},
		{CreateBlankExpr(NewSymbol("Integer")), true, "Blank", NewSymbol("Integer")},
		{CreateBlankSequenceExpr(nil), true, "BlankSequence", nil},
		{CreateBlankNullSequenceExpr(NewSymbol("String")), true, "BlankNullSequence", NewSymbol("String")},
		{NewSymbol("x"), false, "", nil},
		{NewInteger(42), false, "", nil},
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

// Test pattern parsing functions
func TestIsPatternVariable(t *testing.T) {
	tests := []struct {
		name   string
		expect bool
	}{
		{"x_", true},
		{"x_Integer", true},
		{"x__", true},
		{"x___", true},
		{"_", true},
		{"__", true},
		{"___", true},
		{"x", false},
		{"abc", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsPatternVariable(test.name)
		if result != test.expect {
			t.Errorf("IsPatternVariable(%q) = %v, want %v", test.name, result, test.expect)
		}
	}
}

func TestParsePatternVariable(t *testing.T) {
	tests := []struct {
		name       string
		expectVar  string
		expectType string
	}{
		{"x_Integer", "x", "Integer"},
		{"var_String", "var", "String"},
		{"x_", "x", ""},
		{"_Integer", "", "Integer"},
		{"_", "", ""},
		{"abc", "", ""},
	}

	for _, test := range tests {
		varName, typeName := ParsePatternVariable(test.name)
		if varName != test.expectVar {
			t.Errorf("ParsePatternVariable(%q) var = %q, want %q", test.name, varName, test.expectVar)
		}
		if typeName != test.expectType {
			t.Errorf("ParsePatternVariable(%q) type = %q, want %q", test.name, typeName, test.expectType)
		}
	}
}

func TestConvertPatternStringToSymbolic(t *testing.T) {
	tests := []struct {
		name     string
		expected func() Expr
	}{
		{"x_", func() Expr { return CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(nil)) }},
		{"x_Integer", func() Expr { return CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(NewSymbol("Integer"))) }},
		{"_", func() Expr { return CreateBlankExpr(nil) }},
		{"_Integer", func() Expr { return CreateBlankExpr(NewSymbol("Integer")) }},
		{"x__", func() Expr { return CreatePatternExpr(NewSymbol("x"), CreateBlankSequenceExpr(nil)) }},
		{"x___", func() Expr { return CreatePatternExpr(NewSymbol("x"), CreateBlankNullSequenceExpr(nil)) }},
		{"regular", func() Expr { return NewSymbol("regular") }},
	}

	for _, test := range tests {
		result := ConvertPatternStringToSymbolic(test.name)
		expected := test.expected()
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("ConvertPatternStringToSymbolic(%q) = %v, want %v", test.name, result, expected)
		}
	}
}

func TestParsePatternInfo(t *testing.T) {
	tests := []struct {
		name     string
		expected PatternInfo
	}{
		{"x_", PatternInfo{Type: BlankPattern, VarName: "x", TypeName: ""}},
		{"x_Integer", PatternInfo{Type: BlankPattern, VarName: "x", TypeName: "Integer"}},
		{"x__", PatternInfo{Type: BlankSequencePattern, VarName: "x", TypeName: ""}},
		{"x___", PatternInfo{Type: BlankNullSequencePattern, VarName: "x", TypeName: ""}},
		{"_", PatternInfo{Type: BlankPattern, VarName: "", TypeName: ""}},
		{"regular", PatternInfo{}},
	}

	for _, test := range tests {
		result := ParsePatternInfo(test.name)
		if result != test.expected {
			t.Errorf("ParsePatternInfo(%q) = %v, want %v", test.name, result, test.expected)
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
	}

	for _, test := range tests {
		result := matcher.TestMatch(test.pattern, test.expr)
		if result != test.expected {
			t.Errorf("TestMatch(%v, %v) = %v, want %v", test.pattern, test.expr, result, test.expected)
		}
	}
}
