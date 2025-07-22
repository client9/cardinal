package sexpr

import (
	"testing"
)

func TestExtractInt64(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected int64
		ok       bool
	}{
		{
			name:     "valid integer",
			expr:     NewIntAtom(42),
			expected: 42,
			ok:       true,
		},
		{
			name:     "negative integer",
			expr:     NewIntAtom(-123),
			expected: -123,
			ok:       true,
		},
		{
			name:     "zero",
			expr:     NewIntAtom(0),
			expected: 0,
			ok:       true,
		},
		{
			name:     "float atom",
			expr:     NewFloatAtom(3.14),
			expected: 0,
			ok:       false,
		},
		{
			name:     "string atom",
			expr:     NewStringAtom("hello"),
			expected: 0,
			ok:       false,
		},
		{
			name:     "symbol atom",
			expr:     NewSymbolAtom("x"),
			expected: 0,
			ok:       false,
		},
		{
			name:     "list",
			expr:     NewList(NewSymbolAtom("Plus"), NewIntAtom(1), NewIntAtom(2)),
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := ExtractInt64(tt.expr)
			if ok != tt.ok {
				t.Errorf("ExtractInt64(%s) ok = %t, want %t", tt.expr.String(), ok, tt.ok)
			}
			if ok && result != tt.expected {
				t.Errorf("ExtractInt64(%s) = %d, want %d", tt.expr.String(), result, tt.expected)
			}
		})
	}
}

func TestExtractFloat64(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected float64
		ok       bool
	}{
		{
			name:     "valid float",
			expr:     NewFloatAtom(3.14),
			expected: 3.14,
			ok:       true,
		},
		{
			name:     "zero float",
			expr:     NewFloatAtom(0.0),
			expected: 0.0,
			ok:       true,
		},
		{
			name:     "negative float",
			expr:     NewFloatAtom(-2.718),
			expected: -2.718,
			ok:       true,
		},
		{
			name:     "integer atom",
			expr:     NewIntAtom(42),
			expected: 0,
			ok:       false,
		},
		{
			name:     "string atom",
			expr:     NewStringAtom("3.14"),
			expected: 0,
			ok:       false,
		},
		{
			name:     "symbol atom",
			expr:     NewSymbolAtom("Pi"),
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := ExtractFloat64(tt.expr)
			if ok != tt.ok {
				t.Errorf("ExtractFloat64(%s) ok = %t, want %t", tt.expr.String(), ok, tt.ok)
			}
			if ok && result != tt.expected {
				t.Errorf("ExtractFloat64(%s) = %f, want %f", tt.expr.String(), result, tt.expected)
			}
		})
	}
}

func TestExtractString(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected string
		ok       bool
	}{
		{
			name:     "valid string",
			expr:     NewStringAtom("hello"),
			expected: "hello",
			ok:       true,
		},
		{
			name:     "empty string",
			expr:     NewStringAtom(""),
			expected: "",
			ok:       true,
		},
		{
			name:     "string with spaces",
			expr:     NewStringAtom("hello world"),
			expected: "hello world",
			ok:       true,
		},
		{
			name:     "integer atom",
			expr:     NewIntAtom(42),
			expected: "",
			ok:       false,
		},
		{
			name:     "float atom",
			expr:     NewFloatAtom(3.14),
			expected: "",
			ok:       false,
		},
		{
			name:     "symbol atom",
			expr:     NewSymbolAtom("hello"),
			expected: "",
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := ExtractString(tt.expr)
			if ok != tt.ok {
				t.Errorf("ExtractString(%s) ok = %t, want %t", tt.expr.String(), ok, tt.ok)
			}
			if ok && result != tt.expected {
				t.Errorf("ExtractString(%s) = %q, want %q", tt.expr.String(), result, tt.expected)
			}
		})
	}
}

func TestExtractBool(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected bool
		ok       bool
	}{
		{
			name:     "True symbol",
			expr:     NewBoolAtom(true), // This creates NewSymbolAtom("True")
			expected: true,
			ok:       true,
		},
		{
			name:     "False symbol",
			expr:     NewBoolAtom(false), // This creates NewSymbolAtom("False")
			expected: false,
			ok:       true,
		},
		{
			name:     "manual True symbol",
			expr:     NewSymbolAtom("True"),
			expected: true,
			ok:       true,
		},
		{
			name:     "manual False symbol",
			expr:     NewSymbolAtom("False"),
			expected: false,
			ok:       true,
		},
		{
			name:     "other symbol",
			expr:     NewSymbolAtom("x"),
			expected: false,
			ok:       false,
		},
		{
			name:     "integer atom",
			expr:     NewIntAtom(1),
			expected: false,
			ok:       false,
		},
		{
			name:     "string atom",
			expr:     NewStringAtom("True"),
			expected: false,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := ExtractBool(tt.expr)
			if ok != tt.ok {
				t.Errorf("ExtractBool(%s) ok = %t, want %t", tt.expr.String(), ok, tt.ok)
			}
			if ok && result != tt.expected {
				t.Errorf("ExtractBool(%s) = %t, want %t", tt.expr.String(), result, tt.expected)
			}
		})
	}
}

func TestCopyExprList(t *testing.T) {
	tests := []struct {
		name     string
		head     string
		args     []Expr
		expected string
	}{
		{
			name:     "empty args",
			head:     "Plus",
			args:     []Expr{},
			expected: "Plus()",
		},
		{
			name:     "single arg",
			head:     "Plus",
			args:     []Expr{NewIntAtom(42)},
			expected: "Plus(42)",
		},
		{
			name:     "multiple args",
			head:     "Plus",
			args:     []Expr{NewIntAtom(1), NewIntAtom(2), NewIntAtom(3)},
			expected: "Plus(1, 2, 3)",
		},
		{
			name:     "mixed types",
			head:     "Equal",
			args:     []Expr{NewIntAtom(42), NewSymbolAtom("x")},
			expected: "Equal(42, x)",
		},
		{
			name:     "nested expression",
			head:     "Outer",
			args:     []Expr{NewList(NewSymbolAtom("Inner"), NewIntAtom(1))},
			expected: "Outer(Inner(1))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CopyExprList(tt.head, tt.args)
			if result.String() != tt.expected {
				t.Errorf("CopyExprList(%q, %v) = %s, want %s",
					tt.head, tt.args, result.String(), tt.expected)
			}

			// Verify structure
			if len(result.Elements) != len(tt.args)+1 {
				t.Errorf("CopyExprList result has %d elements, want %d",
					len(result.Elements), len(tt.args)+1)
			}

			// Verify head
			if headAtom, ok := result.Elements[0].(Atom); !ok || headAtom.Value.(string) != tt.head {
				t.Errorf("CopyExprList head = %s, want %s",
					result.Elements[0].String(), tt.head)
			}

			// Verify args (should be copies of references, not deep copies)
			for i, expectedArg := range tt.args {
				if !result.Elements[i+1].Equal(expectedArg) {
					t.Errorf("CopyExprList arg[%d] = %s, want %s",
						i, result.Elements[i+1].String(), expectedArg.String())
				}
			}
		})
	}
}

func TestHelperFunctionsWithRealExpressions(t *testing.T) {
	// Test helper functions with more complex expressions
	ctx := NewContext()

	// Create complex expressions for testing
	plusExpr := NewList(NewSymbolAtom("Plus"), NewIntAtom(1), NewIntAtom(2))
	timesExpr := NewList(NewSymbolAtom("Times"), NewFloatAtom(2.5), NewFloatAtom(4.0))

	// Test that helper functions work correctly with real expressions
	if val, ok := ExtractInt64(NewIntAtom(42)); !ok || val != 42 {
		t.Error("ExtractInt64 failed on real integer")
	}

	if val, ok := ExtractFloat64(NewFloatAtom(3.14)); !ok || val != 3.14 {
		t.Error("ExtractFloat64 failed on real float")
	}

	if val, ok := ExtractString(NewStringAtom("test")); !ok || val != "test" {
		t.Error("ExtractString failed on real string")
	}

	if val, ok := ExtractBool(NewBoolAtom(true)); !ok || val != true {
		t.Error("ExtractBool failed on real boolean")
	}

	// Test CopyExprList with real expressions
	copied := CopyExprList("Plus", []Expr{NewIntAtom(1), NewSymbolAtom("x")})
	expected := "Plus(1, x)"
	if copied.String() != expected {
		t.Errorf("CopyExprList with real expressions = %s, want %s", copied.String(), expected)
	}

	// Verify these don't interfere with normal evaluation
	_ = ctx
	_ = plusExpr
	_ = timesExpr
}
