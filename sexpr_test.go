package sexpr

import (
	"testing"
)

func TestAtom_String(t *testing.T) {
	tests := []struct {
		name     string
		atom     *Atom
		expected string
	}{
		{
			name:     "string atom",
			atom:     NewStringAtom("hello world"),
			expected: `"hello world"`,
		},
		{
			name:     "integer atom",
			atom:     NewIntAtom(42),
			expected: "42",
		},
		{
			name:     "float atom",
			atom:     NewFloatAtom(3.14159),
			expected: "3.14159",
		},
		{
			name:     "boolean true atom",
			atom:     NewBoolAtom(true),
			expected: "True",
		},
		{
			name:     "boolean false atom",
			atom:     NewBoolAtom(false),
			expected: "False",
		},
		{
			name:     "symbol atom",
			atom:     NewSymbolAtom("mySymbol"),
			expected: "mySymbol",
		},
		{
			name:     "empty string atom",
			atom:     NewStringAtom(""),
			expected: `""`,
		},
		{
			name:     "zero integer atom",
			atom:     NewIntAtom(0),
			expected: "0",
		},
		{
			name:     "negative integer atom",
			atom:     NewIntAtom(-123),
			expected: "-123",
		},
		{
			name:     "zero float atom",
			atom:     NewFloatAtom(0.0),
			expected: "0",
		},
		{
			name:     "negative float atom",
			atom:     NewFloatAtom(-2.5),
			expected: "-2.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.atom.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestAtom_Type(t *testing.T) {
	tests := []struct {
		name     string
		atom     *Atom
		expected string
	}{
		{
			name:     "string atom type",
			atom:     NewStringAtom("test"),
			expected: "string",
		},
		{
			name:     "integer atom type",
			atom:     NewIntAtom(42),
			expected: "int",
		},
		{
			name:     "float atom type",
			atom:     NewFloatAtom(3.14),
			expected: "float64",
		},
		{
			name:     "boolean atom type",
			atom:     NewBoolAtom(true),
			expected: "bool",
		},
		{
			name:     "symbol atom type",
			atom:     NewSymbolAtom("x"),
			expected: "symbol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.atom.Type()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestList_String(t *testing.T) {
	tests := []struct {
		name     string
		list     *List
		expected string
	}{
		{
			name:     "empty list",
			list:     NewList(),
			expected: "{}",
		},
		{
			name:     "single element list",
			list:     NewList(NewSymbolAtom("Plus")),
			expected: "Plus[]",
		},
		{
			name:     "simple function call",
			list:     NewList(NewSymbolAtom("Plus"), NewIntAtom(1), NewIntAtom(2)),
			expected: "Plus[1, 2]",
		},
		{
			name:     "mixed types",
			list:     NewList(NewSymbolAtom("List"), NewIntAtom(1), NewFloatAtom(2.5), NewStringAtom("hello"), NewBoolAtom(true)),
			expected: `List[1, 2.5, "hello", True]`,
		},
		{
			name:     "nested list",
			list:     NewList(NewSymbolAtom("Plus"), NewIntAtom(1), NewList(NewSymbolAtom("Times"), NewIntAtom(2), NewIntAtom(3))),
			expected: "Plus[1, Times[2, 3]]",
		},
		{
			name:     "deeply nested",
			list:     NewList(NewSymbolAtom("f"), NewList(NewSymbolAtom("g"), NewList(NewSymbolAtom("h"), NewIntAtom(1)))),
			expected: "f[g[h[1]]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.list.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestList_Type(t *testing.T) {
	tests := []struct {
		name     string
		list     *List
		expected string
	}{
		{
			name:     "empty list type",
			list:     NewList(),
			expected: "list",
		},
		{
			name:     "non-empty list type",
			list:     NewList(NewSymbolAtom("Plus"), NewIntAtom(1), NewIntAtom(2)),
			expected: "list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.list.Type()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestConstructorFunctions(t *testing.T) {
	tests := []struct {
		name        string
		constructor func() Expr
		expectedType string
		expectedValue interface{}
	}{
		{
			name:        "NewStringAtom",
			constructor: func() Expr { return NewStringAtom("test") },
			expectedType: "string",
			expectedValue: "test",
		},
		{
			name:        "NewIntAtom",
			constructor: func() Expr { return NewIntAtom(42) },
			expectedType: "int",
			expectedValue: 42,
		},
		{
			name:        "NewFloatAtom",
			constructor: func() Expr { return NewFloatAtom(3.14) },
			expectedType: "float64",
			expectedValue: 3.14,
		},
		{
			name:        "NewBoolAtom true",
			constructor: func() Expr { return NewBoolAtom(true) },
			expectedType: "bool",
			expectedValue: true,
		},
		{
			name:        "NewBoolAtom false",
			constructor: func() Expr { return NewBoolAtom(false) },
			expectedType: "bool",
			expectedValue: false,
		},
		{
			name:        "NewSymbolAtom",
			constructor: func() Expr { return NewSymbolAtom("mySymbol") },
			expectedType: "symbol",
			expectedValue: "mySymbol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := tt.constructor()
			
			if expr.Type() != tt.expectedType {
				t.Errorf("expected type %q, got %q", tt.expectedType, expr.Type())
			}
			
			atom, ok := expr.(*Atom)
			if !ok {
				t.Errorf("expected Atom, got %T", expr)
				return
			}
			
			if atom.Value != tt.expectedValue {
				t.Errorf("expected value %v, got %v", tt.expectedValue, atom.Value)
			}
		})
	}
}

func TestNewList(t *testing.T) {
	tests := []struct {
		name            string
		elements        []Expr
		expectedLength  int
		expectedType    string
	}{
		{
			name:            "empty list",
			elements:        []Expr{},
			expectedLength:  0,
			expectedType:    "list",
		},
		{
			name:            "single element",
			elements:        []Expr{NewIntAtom(1)},
			expectedLength:  1,
			expectedType:    "list",
		},
		{
			name:            "multiple elements",
			elements:        []Expr{NewSymbolAtom("Plus"), NewIntAtom(1), NewIntAtom(2)},
			expectedLength:  3,
			expectedType:    "list",
		},
		{
			name:            "nested list",
			elements:        []Expr{NewSymbolAtom("f"), NewList(NewSymbolAtom("g"), NewIntAtom(1))},
			expectedLength:  2,
			expectedType:    "list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewList(tt.elements...)
			
			if list.Type() != tt.expectedType {
				t.Errorf("expected type %q, got %q", tt.expectedType, list.Type())
			}
			
			if len(list.Elements) != tt.expectedLength {
				t.Errorf("expected length %d, got %d", tt.expectedLength, len(list.Elements))
			}
		})
	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected string
	}{
		{
			name: "mathematical expression",
			expr: NewList(
				NewSymbolAtom("Plus"),
				NewList(NewSymbolAtom("Times"), NewIntAtom(2), NewSymbolAtom("x")),
				NewIntAtom(5),
			),
			expected: "Plus[Times[2, x], 5]",
		},
		{
			name: "function definition",
			expr: NewList(
				NewSymbolAtom("Function"),
				NewStringAtom("x"),
				NewList(NewSymbolAtom("Power"), NewSymbolAtom("x"), NewIntAtom(2)),
			),
			expected: `Function["x", Power[x, 2]]`,
		},
		{
			name: "conditional expression",
			expr: NewList(
				NewSymbolAtom("If"),
				NewList(NewSymbolAtom("Greater"), NewSymbolAtom("x"), NewIntAtom(0)),
				NewSymbolAtom("x"),
				NewList(NewSymbolAtom("Minus"), NewSymbolAtom("x")),
			),
			expected: "If[Greater[x, 0], x, Minus[x]]",
		},
		{
			name: "list with mixed types",
			expr: NewList(
				NewSymbolAtom("List"),
				NewIntAtom(1),
				NewFloatAtom(2.5),
				NewStringAtom("hello"),
				NewBoolAtom(true),
				NewBoolAtom(false),
				NewSymbolAtom("x"),
			),
			expected: `List[1, 2.5, "hello", True, False, x]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.expr.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}