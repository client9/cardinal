package core

import (
	"testing"

	"github.com/client9/cardinal/core/symbol"
)

func TestAtom_String(t *testing.T) {
	tests := []struct {
		name     string
		atom     Expr
		expected string
	}{
		{
			name:     "string atom",
			atom:     NewString("hello world"),
			expected: `"hello world"`,
		},
		{
			name:     "integer atom",
			atom:     NewInteger(42),
			expected: "42",
		},
		{
			name:     "float atom",
			atom:     NewReal(3.14159),
			expected: "3.14159",
		},
		{
			name:     "boolean true atom",
			atom:     NewBool(true),
			expected: "True",
		},
		{
			name:     "boolean false atom",
			atom:     NewBool(false),
			expected: "False",
		},
		{
			name:     "symbol. atom",
			atom:     NewSymbol("mySymbol"),
			expected: "mySymbol",
		},
		{
			name:     "empty string atom",
			atom:     NewString(""),
			expected: `""`,
		},
		{
			name:     "zero integer atom",
			atom:     NewInteger(0),
			expected: "0",
		},
		{
			name:     "negative integer atom",
			atom:     NewInteger(-123),
			expected: "-123",
		},
		{
			name:     "zero float atom",
			atom:     NewReal(0.0),
			expected: "0.0",
		},
		{
			name:     "negative float atom",
			atom:     NewReal(-2.5),
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
		atom     Expr
		expected string
	}{
		{
			name:     "string atom type",
			atom:     NewString("test"),
			expected: "String",
		},
		{
			name:     "integer atom type",
			atom:     NewInteger(42),
			expected: "Integer",
		},
		{
			name:     "float atom type",
			atom:     NewReal(3.14),
			expected: "Real",
		},
		{
			name:     "boolean atom type",
			atom:     NewBool(true),
			expected: "Symbol",
		},
		{
			name:     "symbol. atom type",
			atom:     NewSymbol("x"),
			expected: "Symbol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.atom.Head().String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestList_String(t *testing.T) {
	tests := []struct {
		name     string
		list     List
		expected string
	}{
		{
			name:     "empty list",
			list:     NewList(symbol.List),
			expected: "List()",
		},
		{
			name:     "single element list",
			list:     NewList(symbol.Plus),
			expected: "Plus()",
		},
		{
			name:     "simple function call",
			list:     NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
			expected: "Plus(1, 2)",
		},
		{
			name:     "mixed types",
			list:     NewList(symbol.List, NewInteger(1), NewReal(2.5), NewString("hello"), NewBool(true)),
			expected: `List(1, 2.5, "hello", True)`,
		},
		{
			name:     "nested list",
			list:     NewList(symbol.Plus, NewInteger(1), NewList(symbol.Times, NewInteger(2), NewInteger(3))),
			expected: "Plus(1, Times(2, 3))",
		},
		{
			name:     "deeply nested",
			list:     NewList(NewSymbol("f"), NewList(NewSymbol("g"), NewList(NewSymbol("h"), NewInteger(1)))),
			expected: "f(g(h(1)))",
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
		list     List
		expected string
	}{
		{
			name:     "empty list type",
			list:     NewList(symbol.List),
			expected: "List",
		},
		{
			name:     "non-empty list type",
			list:     NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
			expected: "Plus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.list.Head().String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestConstructorFunctions(t *testing.T) {
	tests := []struct {
		name          string
		constructor   func() Expr
		expectedType  string
		expectedValue interface{}
	}{
		{
			name:          "NewStringAtom",
			constructor:   func() Expr { return NewString("test") },
			expectedType:  "String",
			expectedValue: "test",
		},
		{
			name:          "NewIntAtom",
			constructor:   func() Expr { return NewInteger(42) },
			expectedType:  "Integer",
			expectedValue: 42,
		},
		{
			name:          "NewFloatAtom",
			constructor:   func() Expr { return NewReal(3.14) },
			expectedType:  "Real",
			expectedValue: 3.14,
		},
		{
			name:          "NewBoolAtom true",
			constructor:   func() Expr { return NewBool(true) },
			expectedType:  "Symbol",
			expectedValue: "True",
		},
		{
			name:          "NewBoolAtom false",
			constructor:   func() Expr { return NewBool(false) },
			expectedType:  "Symbol",
			expectedValue: "False",
		},
		{
			name:          "NewSymbolAtom",
			constructor:   func() Expr { return NewSymbol("mySymbol") },
			expectedType:  "Symbol",
			expectedValue: "mySymbol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := tt.constructor()

			if expr.Head().String() != tt.expectedType {
				t.Errorf("expected type %q, got %q", tt.expectedType, expr.Head().String())
			}

			// Check value based on the expected type
			switch tt.expectedType {
			case "String":
				if str, ok := expr.(String); ok {
					if string(str) != tt.expectedValue {
						t.Errorf("expected value %v, got %v", tt.expectedValue, string(str))
					}
				} else {
					t.Errorf("expected String, got %T", expr)
				}
			case "Integer":
				if integer, ok := expr.(Integer); ok {
					if integer.Int64() != int64(tt.expectedValue.(int)) {
						t.Errorf("expected value %v, got %v", tt.expectedValue, integer.Int64())
					}
				} else {
					t.Errorf("expected Integer, got %T", expr)
				}
			case "Real":
				if real, ok := expr.(Real); ok {
					if float64(real) != tt.expectedValue.(float64) {
						t.Errorf("expected value %v, got %v", tt.expectedValue, float64(real))
					}
				} else {
					t.Errorf("expected Real, got %T", expr)
				}
			case "Symbol":
				if symbolName, ok := ExtractSymbol(expr); ok {
					if symbolName != tt.expectedValue {
						t.Errorf("expected value %v, got %v", tt.expectedValue, symbolName)
					}
				} else {
					t.Errorf("expected Symbol, got %T", expr)
				}
			default:
				t.Errorf("unknown expected type: %s", tt.expectedType)
			}
		})
	}
}

func TestNewList(t *testing.T) {
	tests := []struct {
		name           string
		list           List
		expectedLength int
		expectedType   string
	}{
		{
			name:           "multiple elements",
			list:           NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
			expectedLength: 2,
			expectedType:   "Plus",
		},
		{
			name:           "nested list",
			list:           NewList(NewSymbol("f"), NewList(NewSymbol("g"), NewInteger(1))),
			expectedLength: 1,
			expectedType:   "f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.list.Head().String() != tt.expectedType {
				t.Errorf("expected type %q, got %q", tt.expectedType, tt.list.Head().String())
			}

			if len(tt.list.Tail()) != tt.expectedLength {
				t.Errorf("expected length %d, got %d", tt.expectedLength, len(tt.list.AsSlice()))
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
			expr: NewList(symbol.Plus,
				NewList(symbol.Times, NewInteger(2), NewSymbol("x")),
				NewInteger(5),
			),
			expected: "Plus(Times(2, x), 5)",
		},
		{
			name: "function definition",
			expr: NewList(
				symbol.Function,
				NewString("x"),
				NewList(symbol.Power, NewSymbol("x"), NewInteger(2)),
			),
			expected: `Function("x", Power(x, 2))`,
		},
		{
			name: "conditional expression",
			expr: NewList(
				symbol.If,
				NewList(symbol.Greater, NewSymbol("x"), NewInteger(0)),
				NewSymbol("x"),
				NewList(symbol.Minus, NewSymbol("x")),
			),
			expected: "If(Greater(x, 0), x, Minus(x))",
		},
		{
			name: "list with mixed types",
			expr: NewList(
				symbol.List,
				NewInteger(1),
				NewReal(2.5),
				NewString("hello"),
				NewBool(true),
				NewBool(false),
				NewSymbol("x"),
			),
			expected: `List(1, 2.5, "hello", True, False, x)`,
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
