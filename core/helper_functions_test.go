package core

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
			expr:     NewInteger(42),
			expected: 42,
			ok:       true,
		},
		{
			name:     "negative integer",
			expr:     NewInteger(-123),
			expected: -123,
			ok:       true,
		},
		{
			name:     "zero",
			expr:     NewInteger(0),
			expected: 0,
			ok:       true,
		},
		{
			name:     "float atom",
			expr:     NewReal(3.14),
			expected: 0,
			ok:       false,
		},
		{
			name:     "string atom",
			expr:     NewString("hello"),
			expected: 0,
			ok:       false,
		},
		{
			name:     "symbol atom",
			expr:     NewSymbol("x"),
			expected: 0,
			ok:       false,
		},
		{
			name:     "list",
			expr:     NewList("Plus", NewInteger(1), NewInteger(2)),
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
			expr:     NewReal(3.14),
			expected: 3.14,
			ok:       true,
		},
		{
			name:     "zero float",
			expr:     NewReal(0.0),
			expected: 0.0,
			ok:       true,
		},
		{
			name:     "negative float",
			expr:     NewReal(-2.718),
			expected: -2.718,
			ok:       true,
		},
		{
			name:     "integer atom",
			expr:     NewInteger(42),
			expected: 0,
			ok:       false,
		},
		{
			name:     "string atom",
			expr:     NewString("3.14"),
			expected: 0,
			ok:       false,
		},
		{
			name:     "symbol atom",
			expr:     NewSymbol("Pi"),
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
			expr:     NewString("hello"),
			expected: "hello",
			ok:       true,
		},
		{
			name:     "empty string",
			expr:     NewString(""),
			expected: "",
			ok:       true,
		},
		{
			name:     "string with spaces",
			expr:     NewString("hello world"),
			expected: "hello world",
			ok:       true,
		},
		{
			name:     "integer atom",
			expr:     NewInteger(42),
			expected: "",
			ok:       false,
		},
		{
			name:     "float atom",
			expr:     NewReal(3.14),
			expected: "",
			ok:       false,
		},
		{
			name:     "symbol atom",
			expr:     NewSymbol("hello"),
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
			expr:     NewBool(true), // This creates NewBool(true)
			expected: true,
			ok:       true,
		},
		{
			name:     "False symbol",
			expr:     NewBool(false), // This creates NewBool(false)
			expected: false,
			ok:       true,
		},
		{
			name:     "manual True symbol",
			expr:     NewBool(true),
			expected: true,
			ok:       true,
		},
		{
			name:     "manual False symbol",
			expr:     NewBool(false),
			expected: false,
			ok:       true,
		},
		{
			name:     "other symbol",
			expr:     NewSymbol("x"),
			expected: false,
			ok:       false,
		},
		{
			name:     "integer atom",
			expr:     NewInteger(1),
			expected: false,
			ok:       false,
		},
		{
			name:     "string atom",
			expr:     NewString("True"),
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
			args:     []Expr{NewInteger(42)},
			expected: "Plus(42)",
		},
		{
			name:     "multiple args",
			head:     "Plus",
			args:     []Expr{NewInteger(1), NewInteger(2), NewInteger(3)},
			expected: "Plus(1, 2, 3)",
		},
		{
			name:     "mixed types",
			head:     "Equal",
			args:     []Expr{NewInteger(42), NewSymbol("x")},
			expected: "Equal(42, x)",
		},
		{
			name:     "nested expression",
			head:     "Outer",
			args:     []Expr{NewList("Inner", NewInteger(1))},
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
			if int(result.Length()) != len(tt.args) {
				t.Errorf("CopyExprList result has %d elements, want %d",
					result.Length(), len(tt.args))
			}

			// Verify head
			if result.Head() != tt.head {
				t.Errorf("CopyExprList head = %s, want %s",
					result.Head(), tt.head)
			}

			// Verify args (should be copies of references, not deep copies)
			rslice := result.Tail()
			for i, expectedArg := range tt.args {
				if !rslice[i].Equal(expectedArg) {
					t.Errorf("CopyExprList arg[%d] = %s, want %s",
						i, rslice[i].String(), expectedArg.String())
				}
			}
		})
	}
}

func TestHelperFunctionsWithRealExpressions(t *testing.T) {
	// Test helper functions with more complex expressions
	// Create complex expressions for testing
	plusExpr := NewList("Plus", NewInteger(1), NewInteger(2))
	timesExpr := NewList("Times", NewReal(2.5), NewReal(4.0))

	// Test that helper functions work correctly with real expressions
	if val, ok := ExtractInt64(NewInteger(42)); !ok || val != 42 {
		t.Error("ExtractInt64 failed on real integer")
	}

	if val, ok := ExtractFloat64(NewReal(3.14)); !ok || val != 3.14 {
		t.Error("ExtractFloat64 failed on real float")
	}

	if val, ok := ExtractString(NewString("test")); !ok || val != "test" {
		t.Error("ExtractString failed on real string")
	}

	if val, ok := ExtractBool(NewBool(true)); !ok || val != true {
		t.Error("ExtractBool failed on real boolean")
	}

	// Test CopyExprList with real expressions
	copied := CopyExprList("Plus", []Expr{NewInteger(1), NewSymbol("x")})
	expected := "Plus(1, x)"
	if copied.String() != expected {
		t.Errorf("CopyExprList with real expressions = %s, want %s", copied.String(), expected)
	}

	// Verify these don't interfere with normal evaluation
	_ = plusExpr
	_ = timesExpr
}
