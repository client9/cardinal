package sexpr

import (
	"strings"
	"testing"
)

func TestObjectExprBasics(t *testing.T) {
	// Create a simple user-defined type
	uint64Val := NewUint64Value(42)
	objExpr := NewObjectExpr("Uint64", uint64Val)

	// Test String()
	expected := "#2A" // 42 in hex
	if objExpr.String() != expected {
		t.Errorf("Expected String() %s, got %s", expected, objExpr.String())
	}

	// Test Type()
	if objExpr.Type() != "Uint64" {
		t.Errorf("Expected Type() Uint64, got %s", objExpr.Type())
	}

	// Test Equal() - same values
	other := NewObjectExpr("Uint64", NewUint64Value(42))
	if !objExpr.Equal(other) {
		t.Error("Expected equal ObjectExpr with same values to be equal")
	}

	// Test Equal() - different values
	different := NewObjectExpr("Uint64", NewUint64Value(43))
	if objExpr.Equal(different) {
		t.Error("Expected ObjectExpr with different values to not be equal")
	}

	// Test Equal() - different types
	differentType := NewObjectExpr("BigInt", NewUint64Value(42))
	if objExpr.Equal(differentType) {
		t.Error("Expected ObjectExpr with different TypeName to not be equal")
	}

	// Test Equal() - not ObjectExpr
	atom := NewIntAtom(42)
	if objExpr.Equal(atom) {
		t.Error("Expected ObjectExpr to not equal non-ObjectExpr")
	}
}

func TestUint64Constructor(t *testing.T) {
	eval := setupTestEvaluator()
	err := RegisterUint64(eval.context.functionRegistry)
	if err != nil {
		t.Fatalf("Failed to register Uint64: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
		isError  bool
	}{
		// Hex string constructor
		{
			name:     "valid hex string",
			input:    `Uint64("#FF")`,
			expected: "#FF",
			isError:  false,
		},
		{
			name:     "valid large hex",
			input:    `Uint64("#FFFFFFFFFFFFFFFF")`,
			expected: "#FFFFFFFFFFFFFFFF",
			isError:  false,
		},
		{
			name:     "valid zero hex",
			input:    `Uint64("#0")`,
			expected: "#0",
			isError:  false,
		},
		{
			name:    "missing # prefix",
			input:   `Uint64("FF")`,
			isError: true,
		},
		{
			name:    "invalid hex characters",
			input:   `Uint64("#GG")`,
			isError: true,
		},

		// Integer constructor
		{
			name:     "valid integer",
			input:    `Uint64(42)`,
			expected: "#2A",
			isError:  false,
		},
		{
			name:     "valid zero integer",
			input:    `Uint64(0)`,
			expected: "#0",
			isError:  false,
		},
		{
			name:     "valid large integer",
			input:    `Uint64(255)`,
			expected: "#FF",
			isError:  false,
		},
		{
			name:    "negative integer",
			input:   `Uint64(-1)`,
			isError: true,
		},

		// Error cases
		{
			name:    "no arguments",
			input:   `Uint64()`,
			isError: true,
		},
		{
			name:    "too many arguments",
			input:   `Uint64(1, 2)`,
			isError: true,
		},
		{
			name:    "float argument",
			input:   `Uint64(3.14)`,
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluateStringSimple(t, eval, tt.input)

			if tt.isError {
				if !strings.HasPrefix(result, "$Failed") {
					t.Errorf("Expected error, got %s", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

func TestUint64TypePredicate(t *testing.T) {
	eval := setupTestEvaluator()
	err := RegisterUint64(eval.context.functionRegistry)
	if err != nil {
		t.Fatalf("Failed to register Uint64: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Uint64 value",
			input:    `Uint64Q(Uint64("#FF"))`,
			expected: "True",
		},
		{
			name:     "integer",
			input:    `Uint64Q(42)`,
			expected: "False",
		},
		{
			name:     "string",
			input:    `Uint64Q("hello")`,
			expected: "False",
		},
		{
			name:     "symbol",
			input:    `Uint64Q(x)`,
			expected: "False",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluateStringSimple(t, eval, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestUint64Arithmetic(t *testing.T) {
	eval := setupTestEvaluator()
	err := RegisterUint64(eval.context.functionRegistry)
	if err != nil {
		t.Fatalf("Failed to register Uint64: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
		isError  bool
	}{
		// Pure Uint64 operations
		{
			name:     "Uint64 + Uint64",
			input:    `Plus(Uint64("#10"), Uint64("#20"))`,
			expected: "#30",
			isError:  false,
		},
		{
			name:     "Uint64 * Uint64",
			input:    `Times(Uint64("#10"), Uint64("#2"))`,
			expected: "#20",
			isError:  false,
		},

		// Mixed-type operations
		{
			name:     "Integer + Uint64",
			input:    `Plus(16, Uint64("#10"))`,
			expected: "#20",
			isError:  false,
		},
		{
			name:     "Uint64 + Integer",
			input:    `Plus(Uint64("#10"), 16)`,
			expected: "#20",
			isError:  false,
		},

		// Error cases
		{
			name:    "negative integer + Uint64",
			input:   `Plus(-5, Uint64("#10"))`,
			isError: true,
		},
		{
			name:    "Uint64 + negative integer",
			input:   `Plus(Uint64("#10"), -5)`,
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluateStringSimple(t, eval, tt.input)

			if tt.isError {
				if !strings.HasPrefix(result, "$Failed") {
					t.Errorf("Expected error, got %s", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

func TestHeadFunctionWithObjectExpr(t *testing.T) {
	eval := setupTestEvaluator()
	err := RegisterUint64(eval.context.functionRegistry)
	if err != nil {
		t.Fatalf("Failed to register Uint64: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Head of Uint64",
			input:    `Head(Uint64("#FF"))`,
			expected: "Uint64",
		},
		{
			name:     "Head of integer",
			input:    `Head(42)`,
			expected: "Integer",
		},
		{
			name:     "Head of string",
			input:    `Head("hello")`,
			expected: "String",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluateStringSimple(t, eval, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPatternMatchingWithObjectExpr(t *testing.T) {
	eval := setupTestEvaluator()
	err := RegisterUint64(eval.context.functionRegistry)
	if err != nil {
		t.Fatalf("Failed to register Uint64: %v", err)
	}

	// Register a test function that uses Uint64 pattern
	err = eval.context.functionRegistry.RegisterPatternBuiltin("testFunc(x_Uint64)", func(args []Expr, ctx *Context) Expr {
		return NewStringAtom("matched Uint64")
	})
	if err != nil {
		t.Fatalf("Failed to register test function: %v", err)
	}

	err = eval.context.functionRegistry.RegisterPatternBuiltin("testFunc(x_Integer)", func(args []Expr, ctx *Context) Expr {
		return NewStringAtom("matched Integer")
	})
	if err != nil {
		t.Fatalf("Failed to register test function: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Pattern matches Uint64",
			input:    `testFunc(Uint64("#FF"))`,
			expected: `"matched Uint64"`,
		},
		{
			name:     "Pattern matches Integer",
			input:    `testFunc(42)`,
			expected: `"matched Integer"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluateStringSimple(t, eval, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestMatchQWithObjectExpr(t *testing.T) {
	eval := setupTestEvaluator()
	err := RegisterUint64(eval.context.functionRegistry)
	if err != nil {
		t.Fatalf("Failed to register Uint64: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "MatchQ Uint64 with typed blank",
			input:    `MatchQ(Uint64("#FF"), x_Uint64)`,
			expected: "True",
		},
		{
			name:     "MatchQ Integer with Uint64 blank",
			input:    `MatchQ(42, x_Uint64)`,
			expected: "False",
		},
		{
			name:     "MatchQ Uint64 with generic blank",
			input:    `MatchQ(Uint64("#FF"), x_)`,
			expected: "True",
		},
		{
			name:     "MatchQ Uint64 with Integer blank",
			input:    `MatchQ(Uint64("#FF"), x_Integer)`,
			expected: "False",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluateStringSimple(t, eval, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestObjectExprStructuralEquality(t *testing.T) {
	uint64Val1 := NewUint64Value(42)
	uint64Val2 := NewUint64Value(42)
	uint64Val3 := NewUint64Value(43)

	obj1 := NewObjectExpr("Uint64", uint64Val1)
	obj2 := NewObjectExpr("Uint64", uint64Val2)
	obj3 := NewObjectExpr("Uint64", uint64Val3)
	obj4 := NewObjectExpr("BigInt", uint64Val1) // Different type name

	// Same values should be equal
	if !obj1.Equal(obj2) {
		t.Error("ObjectExpr with same TypeName and value should be equal")
	}

	// Different values should not be equal
	if obj1.Equal(obj3) {
		t.Error("ObjectExpr with different values should not be equal")
	}

	// Different type names should not be equal
	if obj1.Equal(obj4) {
		t.Error("ObjectExpr with different TypeName should not be equal")
	}

	// Should not equal non-ObjectExpr
	atom := NewIntAtom(42)
	if obj1.Equal(atom) {
		t.Error("ObjectExpr should not equal non-ObjectExpr")
	}
}
