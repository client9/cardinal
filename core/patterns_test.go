package core

import (
	"reflect"
	"testing"
)

// Test pattern constructor functions
func TestCreateBlankExpr(t *testing.T) {
	// Test blank without type constraint
	blank := CreateBlankExpr(nil)
	expected := List{Elements: []Expr{NewSymbolAtom("Blank")}}
	if !reflect.DeepEqual(blank, expected) {
		t.Errorf("CreateBlankExpr(nil) = %v, want %v", blank, expected)
	}

	// Test blank with type constraint
	typeExpr := NewSymbolAtom("Integer")
	blankTyped := CreateBlankExpr(typeExpr)
	expectedTyped := List{Elements: []Expr{NewSymbolAtom("Blank"), typeExpr}}
	if !reflect.DeepEqual(blankTyped, expectedTyped) {
		t.Errorf("CreateBlankExpr(Integer) = %v, want %v", blankTyped, expectedTyped)
	}
}

func TestCreateBlankSequenceExpr(t *testing.T) {
	// Test sequence without type constraint
	seq := CreateBlankSequenceExpr(nil)
	expected := List{Elements: []Expr{NewSymbolAtom("BlankSequence")}}
	if !reflect.DeepEqual(seq, expected) {
		t.Errorf("CreateBlankSequenceExpr(nil) = %v, want %v", seq, expected)
	}

	// Test sequence with type constraint
	typeExpr := NewSymbolAtom("String")
	seqTyped := CreateBlankSequenceExpr(typeExpr)
	expectedTyped := List{Elements: []Expr{NewSymbolAtom("BlankSequence"), typeExpr}}
	if !reflect.DeepEqual(seqTyped, expectedTyped) {
		t.Errorf("CreateBlankSequenceExpr(String) = %v, want %v", seqTyped, expectedTyped)
	}
}

func TestCreatePatternExpr(t *testing.T) {
	nameExpr := NewSymbolAtom("x")
	blankExpr := CreateBlankExpr(nil)
	pattern := CreatePatternExpr(nameExpr, blankExpr)
	
	expected := List{Elements: []Expr{NewSymbolAtom("Pattern"), nameExpr, blankExpr}}
	if !reflect.DeepEqual(pattern, expected) {
		t.Errorf("CreatePatternExpr = %v, want %v", pattern, expected)
	}
}

// Test pattern analysis functions
func TestIsSymbolicBlank(t *testing.T) {
	tests := []struct {
		expr       Expr
		expectBlank bool
		expectType string
		expectTypeExpr Expr
	}{
		{CreateBlankExpr(nil), true, "Blank", nil},
		{CreateBlankExpr(NewSymbolAtom("Integer")), true, "Blank", NewSymbolAtom("Integer")},
		{CreateBlankSequenceExpr(nil), true, "BlankSequence", nil},
		{CreateBlankNullSequenceExpr(NewSymbolAtom("String")), true, "BlankNullSequence", NewSymbolAtom("String")},
		{NewSymbolAtom("x"), false, "", nil},
		{NewIntAtom(42), false, "", nil},
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
	nameExpr := NewSymbolAtom("x")
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
	isPattern, _, _ = IsSymbolicPattern(NewSymbolAtom("x"))
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
		name     string
		expectVar string
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
		{"x_", func() Expr { return CreatePatternExpr(NewSymbolAtom("x"), CreateBlankExpr(nil)) }},
		{"x_Integer", func() Expr { return CreatePatternExpr(NewSymbolAtom("x"), CreateBlankExpr(NewSymbolAtom("Integer"))) }},
		{"_", func() Expr { return CreateBlankExpr(nil) }},
		{"_Integer", func() Expr { return CreateBlankExpr(NewSymbolAtom("Integer")) }},
		{"x__", func() Expr { return CreatePatternExpr(NewSymbolAtom("x"), CreateBlankSequenceExpr(nil)) }},
		{"x___", func() Expr { return CreatePatternExpr(NewSymbolAtom("x"), CreateBlankNullSequenceExpr(nil)) }},
		{"regular", func() Expr { return NewSymbolAtom("regular") }},
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
		{NewIntAtom(42), "Integer", true},
		{NewIntAtom(42), "Number", true},
		{NewIntAtom(42), "String", false},
		{NewFloatAtom(3.14), "Real", true},
		{NewFloatAtom(3.14), "Number", true},
		{NewStringAtom("hello"), "String", true},
		{NewSymbolAtom("x"), "Symbol", true},
		{NewList(NewSymbolAtom("List")), "List", true},
		{NewIntAtom(42), "", true}, // No constraint
		{NewObjectExpr("CustomType", NewIntAtom(1)), "CustomType", true},
		{NewObjectExpr("CustomType", NewIntAtom(1)), "OtherType", false},
	}

	for _, test := range tests {
		result := MatchesType(test.expr, test.typeName)
		if result != test.expected {
			t.Errorf("MatchesType(%v, %q) = %v, want %v", test.expr, test.typeName, result, test.expected)
		}
	}
}

func TestIsBuiltinType(t *testing.T) {
	builtinTypes := []string{"Integer", "Real", "Number", "String", "Symbol", "List", "Atom"}
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
		{NewIntAtom(42), SpecificityLiteral},
		{NewSymbolAtom("x"), SpecificityLiteral},
		{CreateBlankExpr(nil), SpecificityGeneral},
		{CreateBlankExpr(NewSymbolAtom("Integer")), SpecificityBuiltinType},
		{CreateBlankExpr(NewSymbolAtom("CustomType")), SpecificityUserType},
		{CreatePatternExpr(NewSymbolAtom("x"), CreateBlankExpr(nil)), SpecificityGeneral},
	}

	for _, test := range tests {
		result := GetPatternSpecificity(test.pattern)
		if result != test.expected {
			t.Errorf("GetPatternSpecificity(%v) = %v, want %v", test.pattern, result, test.expected)
		}
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
		{NewIntAtom(42), NewIntAtom(42), true},
		{NewIntAtom(42), NewIntAtom(43), false},
		{NewSymbolAtom("x"), NewSymbolAtom("x"), true},
		{NewSymbolAtom("x"), NewSymbolAtom("y"), false},

		// Blank patterns
		{CreateBlankExpr(nil), NewIntAtom(42), true},
		{CreateBlankExpr(nil), NewStringAtom("hello"), true},
		{CreateBlankExpr(NewSymbolAtom("Integer")), NewIntAtom(42), true},
		{CreateBlankExpr(NewSymbolAtom("Integer")), NewStringAtom("hello"), false},

		// Pattern expressions
		{CreatePatternExpr(NewSymbolAtom("x"), CreateBlankExpr(nil)), NewIntAtom(42), true},
		{CreatePatternExpr(NewSymbolAtom("x"), CreateBlankExpr(NewSymbolAtom("String"))), NewStringAtom("hello"), true},
		{CreatePatternExpr(NewSymbolAtom("x"), CreateBlankExpr(NewSymbolAtom("String"))), NewIntAtom(42), false},

		// List patterns
		{NewList(NewSymbolAtom("Plus"), CreateBlankExpr(nil), CreateBlankExpr(nil)), 
		 NewList(NewSymbolAtom("Plus"), NewIntAtom(1), NewIntAtom(2)), true},
		{NewList(NewSymbolAtom("Plus"), CreateBlankExpr(nil), CreateBlankExpr(nil)), 
		 NewList(NewSymbolAtom("Times"), NewIntAtom(1), NewIntAtom(2)), false},
	}

	for _, test := range tests {
		result := matcher.TestMatch(test.pattern, test.expr)
		if result != test.expected {
			t.Errorf("TestMatch(%v, %v) = %v, want %v", test.pattern, test.expr, result, test.expected)
		}
	}
}