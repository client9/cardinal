package sexpr

import (
	"github.com/client9/sexpr/core"
	"testing"
)

func TestAtom_String(t *testing.T) {
	tests := []struct {
		name     string
		atom     Expr
		expected string
	}{
		{
			name:     "string atom",
			atom:     core.NewString("hello world"),
			expected: `"hello world"`,
		},
		{
			name:     "integer atom",
			atom:     core.NewInteger(42),
			expected: "42",
		},
		{
			name:     "float atom",
			atom:     core.NewReal(3.14159),
			expected: "3.14159",
		},
		{
			name:     "boolean true atom",
			atom:     core.NewBool(true),
			expected: "True",
		},
		{
			name:     "boolean false atom",
			atom:     core.NewBool(false),
			expected: "False",
		},
		{
			name:     "symbol atom",
			atom:     core.NewSymbol("mySymbol"),
			expected: "mySymbol",
		},
		{
			name:     "empty string atom",
			atom:     core.NewString(""),
			expected: `""`,
		},
		{
			name:     "zero integer atom",
			atom:     core.NewInteger(0),
			expected: "0",
		},
		{
			name:     "negative integer atom",
			atom:     core.NewInteger(-123),
			expected: "-123",
		},
		{
			name:     "zero float atom",
			atom:     core.NewReal(0.0),
			expected: "0.0",
		},
		{
			name:     "negative float atom",
			atom:     core.NewReal(-2.5),
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
			atom:     core.NewString("test"),
			expected: "String",
		},
		{
			name:     "integer atom type",
			atom:     core.NewInteger(42),
			expected: "Integer",
		},
		{
			name:     "float atom type",
			atom:     core.NewReal(3.14),
			expected: "Real",
		},
		{
			name:     "boolean atom type",
			atom:     core.NewBool(true),
			expected: "Symbol",
		},
		{
			name:     "symbol atom type",
			atom:     core.NewSymbol("x"),
			expected: "Symbol",
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
		list     List
		expected string
	}{
		{
			name:     "empty list",
			list:     NewList(),
			expected: "List()",
		},
		{
			name:     "single element list",
			list:     NewList(core.NewSymbol("Plus")),
			expected: "Plus()",
		},
		{
			name:     "simple function call",
			list:     NewList(core.NewSymbol("Plus"), core.NewInteger(1), core.NewInteger(2)),
			expected: "Plus(1, 2)",
		},
		{
			name:     "mixed types",
			list:     NewList(core.NewSymbol("List"), core.NewInteger(1), core.NewReal(2.5), core.NewString("hello"), core.NewBool(true)),
			expected: `List(1, 2.5, "hello", True)`,
		},
		{
			name:     "nested list",
			list:     NewList(core.NewSymbol("Plus"), core.NewInteger(1), NewList(core.NewSymbol("Times"), core.NewInteger(2), core.NewInteger(3))),
			expected: "Plus(1, Times(2, 3))",
		},
		{
			name:     "deeply nested",
			list:     NewList(core.NewSymbol("f"), NewList(core.NewSymbol("g"), NewList(core.NewSymbol("h"), core.NewInteger(1)))),
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
			list:     NewList(),
			expected: "List",
		},
		{
			name:     "non-empty list type",
			list:     NewList(core.NewSymbol("Plus"), core.NewInteger(1), core.NewInteger(2)),
			expected: "Plus",
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
		name          string
		constructor   func() Expr
		expectedType  string
		expectedValue interface{}
	}{
		{
			name:          "NewStringAtom",
			constructor:   func() Expr { return core.NewString("test") },
			expectedType:  "String",
			expectedValue: "test",
		},
		{
			name:          "NewIntAtom",
			constructor:   func() Expr { return core.NewInteger(42) },
			expectedType:  "Integer",
			expectedValue: 42,
		},
		{
			name:          "NewFloatAtom",
			constructor:   func() Expr { return core.NewReal(3.14) },
			expectedType:  "Real",
			expectedValue: 3.14,
		},
		{
			name:          "NewBoolAtom true",
			constructor:   func() Expr { return core.NewBool(true) },
			expectedType:  "Symbol",
			expectedValue: "True",
		},
		{
			name:          "NewBoolAtom false",
			constructor:   func() Expr { return core.NewBool(false) },
			expectedType:  "Symbol",
			expectedValue: "False",
		},
		{
			name:          "NewSymbolAtom",
			constructor:   func() Expr { return core.NewSymbol("mySymbol") },
			expectedType:  "Symbol",
			expectedValue: "mySymbol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := tt.constructor()

			if expr.Type() != tt.expectedType {
				t.Errorf("expected type %q, got %q", tt.expectedType, expr.Type())
			}

			// Check value based on the expected type
			switch tt.expectedType {
			case "String":
				if str, ok := expr.(core.String); ok {
					if string(str) != tt.expectedValue {
						t.Errorf("expected value %v, got %v", tt.expectedValue, string(str))
					}
				} else {
					t.Errorf("expected String, got %T", expr)
				}
			case "Integer":
				if integer, ok := expr.(core.Integer); ok {
					if int64(integer) != int64(tt.expectedValue.(int)) {
						t.Errorf("expected value %v, got %v", tt.expectedValue, int64(integer))
					}
				} else {
					t.Errorf("expected Integer, got %T", expr)
				}
			case "Real":
				if real, ok := expr.(core.Real); ok {
					if float64(real) != tt.expectedValue.(float64) {
						t.Errorf("expected value %v, got %v", tt.expectedValue, float64(real))
					}
				} else {
					t.Errorf("expected Real, got %T", expr)
				}
			case "Symbol":
				if symbolName, ok := core.ExtractSymbol(expr); ok {
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
		elements       []Expr
		expectedLength int
		expectedType   string
	}{
		{
			name:           "empty list",
			elements:       []Expr{},
			expectedLength: 0,
			expectedType:   "List",
		},
		/* this is NOT RIGHT!  Must have head.
		{
			name:           "single element",
			elements:       []Expr{core.NewInteger(1)},
			expectedLength: 1,
			expectedType:   "List",
		},
		*/
		{
			name:           "multiple elements",
			elements:       []Expr{core.NewSymbol("Plus"), core.NewInteger(1), core.NewInteger(2)},
			expectedLength: 3,
			expectedType:   "Plus",
		},
		{
			name:           "nested list",
			elements:       []Expr{core.NewSymbol("f"), NewList(core.NewSymbol("g"), core.NewInteger(1))},
			expectedLength: 2,
			expectedType:   "f",
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
				core.NewSymbol("Plus"),
				NewList(core.NewSymbol("Times"), core.NewInteger(2), core.NewSymbol("x")),
				core.NewInteger(5),
			),
			expected: "Plus(Times(2, x), 5)",
		},
		{
			name: "function definition",
			expr: NewList(
				core.NewSymbol("Function"),
				core.NewString("x"),
				NewList(core.NewSymbol("Power"), core.NewSymbol("x"), core.NewInteger(2)),
			),
			expected: `Function("x", Power(x, 2))`,
		},
		{
			name: "conditional expression",
			expr: NewList(
				core.NewSymbol("If"),
				NewList(core.NewSymbol("Greater"), core.NewSymbol("x"), core.NewInteger(0)),
				core.NewSymbol("x"),
				NewList(core.NewSymbol("Minus"), core.NewSymbol("x")),
			),
			expected: "If(Greater(x, 0), x, Minus(x))",
		},
		{
			name: "list with mixed types",
			expr: NewList(
				core.NewSymbol("List"),
				core.NewInteger(1),
				core.NewReal(2.5),
				core.NewString("hello"),
				core.NewBool(true),
				core.NewBool(false),
				core.NewSymbol("x"),
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
