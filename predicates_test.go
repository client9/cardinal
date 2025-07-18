package sexpr

import (
	"testing"
)

func TestEvaluateLength(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected int
	}{
		{
			name:     "Empty list",
			arg:      &List{Elements: []Expr{}},
			expected: 0,
		},
		{
			name: "Single element list",
			arg: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(42),
			}},
			expected: 1,
		},
		{
			name: "Multiple element list",
			arg: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
			}},
			expected: 3,
		},
		{
			name:     "Integer atom",
			arg:      NewIntAtom(42),
			expected: 0,
		},
		{
			name:     "String atom",
			arg:      NewStringAtom("hello"),
			expected: 0,
		},
		{
			name:     "Symbol atom",
			arg:      NewSymbolAtom("x"),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateLength([]Expr{tt.arg})
			
			if !isNumeric(result) {
				t.Errorf("expected numeric result, got %T", result)
				return
			}
			
			val, _ := getNumericValue(result)
			if int(val) != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, int(val))
			}
		})
	}
}

func TestEvaluateLength_ArgumentValidation(t *testing.T) {
	tests := []struct {
		name string
		args []Expr
	}{
		{
			name: "No arguments",
			args: []Expr{},
		},
		{
			name: "Too many arguments",
			args: []Expr{NewIntAtom(1), NewIntAtom(2)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateLength(tt.args)
			
			if !IsError(result) {
				t.Errorf("expected error for %s, got %s", tt.name, result.String())
			}
		})
	}
}

func TestEvaluateListQ(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected bool
	}{
		{
			name:     "Empty list",
			arg:      &List{Elements: []Expr{}},
			expected: true,
		},
		{
			name: "Non-empty list",
			arg: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				NewIntAtom(2),
			}},
			expected: true,
		},
		{
			name:     "Integer atom",
			arg:      NewIntAtom(42),
			expected: false,
		},
		{
			name:     "String atom",
			arg:      NewStringAtom("hello"),
			expected: false,
		},
		{
			name:     "Symbol atom",
			arg:      NewSymbolAtom("x"),
			expected: false,
		},
		{
			name:     "Boolean atom",
			arg:      NewBoolAtom(true),
			expected: false,
		},
		{
			name:     "Float atom",
			arg:      NewFloatAtom(3.14),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateListQ([]Expr{tt.arg})
			
			if !isBool(result) {
				t.Errorf("expected boolean result, got %T", result)
				return
			}
			
			val, _ := getBoolValue(result)
			if val != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, val)
			}
		})
	}
}

func TestEvaluateNumberQ(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected bool
	}{
		{
			name:     "Integer atom",
			arg:      NewIntAtom(42),
			expected: true,
		},
		{
			name:     "Float atom",
			arg:      NewFloatAtom(3.14),
			expected: true,
		},
		{
			name:     "Negative integer",
			arg:      NewIntAtom(-5),
			expected: true,
		},
		{
			name:     "Zero",
			arg:      NewIntAtom(0),
			expected: true,
		},
		{
			name:     "String atom",
			arg:      NewStringAtom("hello"),
			expected: false,
		},
		{
			name:     "Symbol atom",
			arg:      NewSymbolAtom("x"),
			expected: false,
		},
		{
			name:     "Boolean atom",
			arg:      NewBoolAtom(true),
			expected: false,
		},
		{
			name: "List",
			arg: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
			}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateNumberQ([]Expr{tt.arg})
			
			if !isBool(result) {
				t.Errorf("expected boolean result, got %T", result)
				return
			}
			
			val, _ := getBoolValue(result)
			if val != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, val)
			}
		})
	}
}

func TestEvaluateBooleanQ(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected bool
	}{
		{
			name:     "Boolean true",
			arg:      NewBoolAtom(true),
			expected: true,
		},
		{
			name:     "Boolean false",
			arg:      NewBoolAtom(false),
			expected: true,
		},
		{
			name:     "Integer atom",
			arg:      NewIntAtom(42),
			expected: false,
		},
		{
			name:     "Float atom",
			arg:      NewFloatAtom(3.14),
			expected: false,
		},
		{
			name:     "String atom",
			arg:      NewStringAtom("hello"),
			expected: false,
		},
		{
			name:     "Symbol atom",
			arg:      NewSymbolAtom("x"),
			expected: false,
		},
		{
			name: "List",
			arg: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
			}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateBooleanQ([]Expr{tt.arg})
			
			if !isBool(result) {
				t.Errorf("expected boolean result, got %T", result)
				return
			}
			
			val, _ := getBoolValue(result)
			if val != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, val)
			}
		})
	}
}

func TestEvaluateIntegerQ(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected bool
	}{
		{
			name:     "Integer atom",
			arg:      NewIntAtom(42),
			expected: true,
		},
		{
			name:     "Negative integer",
			arg:      NewIntAtom(-5),
			expected: true,
		},
		{
			name:     "Zero",
			arg:      NewIntAtom(0),
			expected: true,
		},
		{
			name:     "Float atom",
			arg:      NewFloatAtom(3.14),
			expected: false,
		},
		{
			name:     "Float that could be integer",
			arg:      NewFloatAtom(42.0),
			expected: false,
		},
		{
			name:     "String atom",
			arg:      NewStringAtom("hello"),
			expected: false,
		},
		{
			name:     "Symbol atom",
			arg:      NewSymbolAtom("x"),
			expected: false,
		},
		{
			name:     "Boolean atom",
			arg:      NewBoolAtom(true),
			expected: false,
		},
		{
			name: "List",
			arg: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
			}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateIntegerQ([]Expr{tt.arg})
			
			if !isBool(result) {
				t.Errorf("expected boolean result, got %T", result)
				return
			}
			
			val, _ := getBoolValue(result)
			if val != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, val)
			}
		})
	}
}

func TestEvaluateAtomQ(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected bool
	}{
		{
			name:     "Integer atom",
			arg:      NewIntAtom(42),
			expected: true,
		},
		{
			name:     "Float atom",
			arg:      NewFloatAtom(3.14),
			expected: true,
		},
		{
			name:     "String atom",
			arg:      NewStringAtom("hello"),
			expected: true,
		},
		{
			name:     "Symbol atom",
			arg:      NewSymbolAtom("x"),
			expected: true,
		},
		{
			name:     "Boolean atom",
			arg:      NewBoolAtom(true),
			expected: true,
		},
		{
			name: "Empty list",
			arg:  &List{Elements: []Expr{}},
			expected: false,
		},
		{
			name: "Non-empty list",
			arg: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
			}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateAtomQ([]Expr{tt.arg})
			
			if !isBool(result) {
				t.Errorf("expected boolean result, got %T", result)
				return
			}
			
			val, _ := getBoolValue(result)
			if val != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, val)
			}
		})
	}
}

func TestEvaluateSymbolQ(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected bool
	}{
		{
			name:     "Symbol atom",
			arg:      NewSymbolAtom("x"),
			expected: true,
		},
		{
			name:     "Symbol with special chars",
			arg:      NewSymbolAtom("$Variable"),
			expected: true,
		},
		{
			name:     "Integer atom",
			arg:      NewIntAtom(42),
			expected: false,
		},
		{
			name:     "Float atom",
			arg:      NewFloatAtom(3.14),
			expected: false,
		},
		{
			name:     "String atom",
			arg:      NewStringAtom("hello"),
			expected: false,
		},
		{
			name:     "Boolean atom",
			arg:      NewBoolAtom(true),
			expected: false,
		},
		{
			name: "List",
			arg: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
			}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateSymbolQ([]Expr{tt.arg})
			
			if !isBool(result) {
				t.Errorf("expected boolean result, got %T", result)
				return
			}
			
			val, _ := getBoolValue(result)
			if val != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, val)
			}
		})
	}
}

// Integration tests with the evaluator
func TestPredicates_Integration(t *testing.T) {
	eval := setupTestEvaluator()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Length tests
		{
			name:     "Length of empty list",
			input:    "Length([])",
			expected: "0",
		},
		{
			name:     "Length of evaluated function call",
			input:    "Length(Plus(1, 2, 3))",
			expected: "0", // Plus(1,2,3) evaluates to 6, Length(6) = 0
		},
		{
			name:     "Length of atom",
			input:    "Length(42)",
			expected: "0",
		},
		
		// ListQ tests
		{
			name:     "ListQ on empty list",
			input:    "ListQ([])",
			expected: "True",
		},
		{
			name:     "ListQ on evaluated function call",
			input:    "ListQ(Plus(1, 2))",
			expected: "False", // Plus(1,2) evaluates to 3, ListQ(3) = False
		},
		{
			name:     "ListQ on atom",
			input:    "ListQ(42)",
			expected: "False",
		},
		
		// NumberQ tests
		{
			name:     "NumberQ on integer",
			input:    "NumberQ(42)",
			expected: "True",
		},
		{
			name:     "NumberQ on float",
			input:    "NumberQ(3.14)",
			expected: "True",
		},
		{
			name:     "NumberQ on string",
			input:    "NumberQ(\"hello\")",
			expected: "False",
		},
		
		// BooleanQ tests
		{
			name:     "BooleanQ on True",
			input:    "BooleanQ(True)",
			expected: "True",
		},
		{
			name:     "BooleanQ on False",
			input:    "BooleanQ(False)",
			expected: "True",
		},
		{
			name:     "BooleanQ on integer",
			input:    "BooleanQ(1)",
			expected: "False",
		},
		
		// IntegerQ tests
		{
			name:     "IntegerQ on integer",
			input:    "IntegerQ(42)",
			expected: "True",
		},
		{
			name:     "IntegerQ on float",
			input:    "IntegerQ(3.14)",
			expected: "False",
		},
		{
			name:     "IntegerQ on symbol",
			input:    "IntegerQ(x)",
			expected: "False",
		},
		
		// AtomQ tests
		{
			name:     "AtomQ on integer",
			input:    "AtomQ(42)",
			expected: "True",
		},
		{
			name:     "AtomQ on symbol",
			input:    "AtomQ(x)",
			expected: "True",
		},
		{
			name:     "AtomQ on evaluated function call",
			input:    "AtomQ(Plus(1, 2))",
			expected: "True", // Plus(1,2) evaluates to 3, AtomQ(3) = True
		},
		{
			name:     "AtomQ on held expression",
			input:    "AtomQ(Hold(Plus(1, 2)))",
			expected: "False", // Hold[Plus(1,2)) is a list, AtomQ = False
		},
		{
			name:     "ListQ on held expression",
			input:    "ListQ(Hold(Plus(1, 2)))",
			expected: "True", // Hold[Plus(1,2)) is a list, ListQ = True
		},
		{
			name:     "Length of held expression",
			input:    "Length(Hold(Plus(1, 2)))",
			expected: "1", // Hold(Plus(1,2)) has 1 argument
		},
		
		// SymbolQ tests
		{
			name:     "SymbolQ on symbol",
			input:    "SymbolQ(x)",
			expected: "True",
		},
		{
			name:     "SymbolQ on integer",
			input:    "SymbolQ(42)",
			expected: "False",
		},
		{
			name:     "SymbolQ on string",
			input:    "SymbolQ(\"hello\")",
			expected: "False",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			
			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

// Test argument validation for all predicate functions
func TestPredicates_ArgumentValidation(t *testing.T) {
	functions := []struct {
		name string
		fn   func([]Expr) Expr
	}{
		{"ListQ", EvaluateListQ},
		{"NumberQ", EvaluateNumberQ},
		{"BooleanQ", EvaluateBooleanQ},
		{"IntegerQ", EvaluateIntegerQ},
		{"AtomQ", EvaluateAtomQ},
		{"SymbolQ", EvaluateSymbolQ},
	}
	
	for _, fn := range functions {
		t.Run(fn.name+"_no_args", func(t *testing.T) {
			result := fn.fn([]Expr{})
			if !IsError(result) {
				t.Errorf("expected error for no arguments, got %s", result.String())
			}
		})
		
		t.Run(fn.name+"_too_many_args", func(t *testing.T) {
			result := fn.fn([]Expr{NewIntAtom(1), NewIntAtom(2)})
			if !IsError(result) {
				t.Errorf("expected error for too many arguments, got %s", result.String())
			}
		})
	}
}
