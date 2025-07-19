package sexpr

import (
	"testing"
)

func TestSimpleVariablePatterns(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple pattern variable",
			input:    "f(x_) := x + 1; f(5)",
			expected: "6",
		},
		{
			name:     "Pattern variable with expression",
			input:    "g(x_) := x * 2; g(3 + 4)",
			expected: "14",
		},
		{
			name:     "Anonymous pattern",
			input:    "h(_) := 42; h(anything)",
			expected: "42",
		},
		{
			name:     "Multiple pattern variables",
			input:    "add(x_, y_) := x + y; add(3, 4)",
			expected: "7",
		},
		{
			name:     "Mixed regular and pattern parameters",
			input:    "mixed(x, y_) := x * y + 1; mixed(5, 3)",
			expected: "16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result, err := evaluateExpression(evaluator, tt.input)
			if err != nil {
				t.Fatalf("Error evaluating expression: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestLiteralPatterns(t *testing.T) {
	tests := []struct {
		name     string
		setup    []string
		input    string
		expected string
	}{
		{
			name: "Factorial base case",
			setup: []string{
				"factorial(0) := 1",
				"factorial(n_) := n * factorial(n - 1)",
			},
			input:    "factorial(0)",
			expected: "1",
		},
		{
			name: "Factorial recursive case",
			setup: []string{
				"factorial(0) := 1",
				"factorial(n_) := n * factorial(n - 1)",
			},
			input:    "factorial(5)",
			expected: "120",
		},
		{
			name: "Number literal patterns",
			setup: []string{
				"sign(-1) := \"negative\"",
				"sign(0) := \"zero\"",
				"sign(1) := \"positive\"",
				"sign(n_) := \"unknown\"",
			},
			input:    "sign(-1)",
			expected: "\"negative\"",
		},
		{
			name: "Number literal fallback",
			setup: []string{
				"sign(-1) := \"negative\"",
				"sign(0) := \"zero\"",
				"sign(1) := \"positive\"",
				"sign(n_) := \"unknown\"",
			},
			input:    "sign(5)",
			expected: "\"unknown\"",
		},
		{
			name: "Boolean literal patterns",
			setup: []string{
				"bool_test(True) := \"it's true\"",
				"bool_test(False) := \"it's false\"",
				"bool_test(x_) := \"not boolean\"",
			},
			input:    "bool_test(True)",
			expected: "\"it's true\"",
		},
		{
			name: "Boolean literal fallback",
			setup: []string{
				"bool_test(True) := \"it's true\"",
				"bool_test(False) := \"it's false\"",
				"bool_test(x_) := \"not boolean\"",
			},
			input:    "bool_test(42)",
			expected: "\"not boolean\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			
			// Setup function definitions
			for _, setup := range tt.setup {
				_, err := evaluateExpression(evaluator, setup)
				if err != nil {
					t.Fatalf("Error setting up expression '%s': %v", setup, err)
				}
			}
			
			// Test the actual expression
			result, err := evaluateExpression(evaluator, tt.input)
			if err != nil {
				t.Fatalf("Error evaluating expression: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPatternPriority(t *testing.T) {
	tests := []struct {
		name     string
		setup    []string
		input    string
		expected string
	}{
		{
			name: "Pattern redefinition",
			setup: []string{
				"test(1) := \"first: one\"",
				"test(1) := \"second: one\"",
				"test(x_) := \"general case\"",
			},
			input:    "test(1)",
			expected: "\"second: one\"",
		},
		{
			name: "Fallback pattern",
			setup: []string{
				"test(1) := \"first: one\"",
				"test(1) := \"second: one\"",
				"test(x_) := \"general case\"",
			},
			input:    "test(2)",
			expected: "\"general case\"",
		},
		{
			name: "Specific before general",
			setup: []string{
				"fibonacci(0) := 0",
				"fibonacci(1) := 1",
				"fibonacci(n_) := fibonacci(n - 1) + fibonacci(n - 2)",
			},
			input:    "fibonacci(1)",
			expected: "1",
		},
		{
			name: "General case recursive",
			setup: []string{
				"fibonacci(0) := 0",
				"fibonacci(1) := 1",
				"fibonacci(n_) := fibonacci(n - 1) + fibonacci(n - 2)",
			},
			input:    "fibonacci(4)",
			expected: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			
			// Setup function definitions
			for _, setup := range tt.setup {
				_, err := evaluateExpression(evaluator, setup)
				if err != nil {
					t.Fatalf("Error setting up expression '%s': %v", setup, err)
				}
			}
			
			// Test the actual expression
			result, err := evaluateExpression(evaluator, tt.input)
			if err != nil {
				t.Fatalf("Error evaluating expression: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPatternMatching(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expr     string
		expected bool
	}{
		{
			name:     "Variable pattern matches number",
			pattern:  "x_",
			expr:     "5",
			expected: true,
		},
		{
			name:     "Variable pattern matches symbol",
			pattern:  "x_",
			expr:     "a",
			expected: true,
		},
		{
			name:     "Anonymous pattern matches anything",
			pattern:  "_",
			expr:     "anything",
			expected: true,
		},
		{
			name:     "Literal number matches",
			pattern:  "5",
			expr:     "5",
			expected: true,
		},
		{
			name:     "Literal number doesn't match different number",
			pattern:  "5",
			expr:     "6",
			expected: false,
		},
		{
			name:     "Literal boolean matches",
			pattern:  "True",
			expr:     "True",
			expected: true,
		},
		{
			name:     "Literal boolean doesn't match different boolean",
			pattern:  "True",
			expr:     "False",
			expected: false,
		},
		{
			name:     "Regular symbol matches same symbol",
			pattern:  "x",
			expr:     "x",
			expected: true,
		},
		{
			name:     "Regular symbol doesn't match different symbol",
			pattern:  "x",
			expr:     "y",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			ctx := NewContext()
			
			// Parse pattern and expression
			parser := NewParser(NewLexer(tt.pattern))
			pattern, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing pattern: %v", err)
			}
			
			parser = NewParser(NewLexer(tt.expr))
			expr, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing expression: %v", err)
			}
			
			// Test pattern matching
			result := evaluator.matchPattern(pattern, expr, ctx)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestPatternVariableBinding(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expr     string
		varName  string
		expected string
	}{
		{
			name:     "Variable pattern binds value",
			pattern:  "x_",
			expr:     "5",
			varName:  "x",
			expected: "5",
		},
		{
			name:     "Variable pattern binds expression",
			pattern:  "y_",
			expr:     "Plus(1, 2)",
			varName:  "y",
			expected: "Plus(1, 2)",
		},
		{
			name:     "Regular parameter binds value (in parameter context)",
			pattern:  "z",
			expr:     "42",
			varName:  "z",
			expected: "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			ctx := NewContext()
			
			// Parse pattern and expression
			parser := NewParser(NewLexer(tt.pattern))
			pattern, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing pattern: %v", err)
			}
			
			parser = NewParser(NewLexer(tt.expr))
			expr, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing expression: %v", err)
			}
			
			// Test pattern matching (use parameter context for regular symbols)
			matched := evaluator.matchPatternInternal(pattern, expr, ctx, true)
			if !matched {
				t.Fatalf("Pattern should have matched")
			}
			
			// Check if variable was bound correctly
			if value, exists := ctx.Get(tt.varName); exists {
				if value.String() != tt.expected {
					t.Errorf("Expected variable %s to be bound to %s, got %s", tt.varName, tt.expected, value.String())
				}
			} else {
				t.Errorf("Variable %s should have been bound", tt.varName)
			}
		})
	}
}

func TestMultipleFunctionDefinitions(t *testing.T) {
	tests := []struct {
		name     string
		setup    []string
		input    string
		expected string
	}{
		{
			name: "Multiple definitions stored correctly",
			setup: []string{
				"f(0) := \"zero\"",
				"f(1) := \"one\"",
				"f(x_) := \"other\"",
			},
			input:    "f(0)",
			expected: "\"zero\"",
		},
		{
			name: "Second definition accessible",
			setup: []string{
				"f(0) := \"zero\"",
				"f(1) := \"one\"",
				"f(x_) := \"other\"",
			},
			input:    "f(1)",
			expected: "\"one\"",
		},
		{
			name: "Fallback definition works",
			setup: []string{
				"f(0) := \"zero\"",
				"f(1) := \"one\"",
				"f(x_) := \"other\"",
			},
			input:    "f(99)",
			expected: "\"other\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			
			// Setup function definitions
			for _, setup := range tt.setup {
				_, err := evaluateExpression(evaluator, setup)
				if err != nil {
					t.Fatalf("Error setting up expression '%s': %v", setup, err)
				}
			}
			
			// Test the actual expression
			result, err := evaluateExpression(evaluator, tt.input)
			if err != nil {
				t.Fatalf("Error evaluating expression: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Helper function to evaluate expressions in tests
func evaluateExpression(evaluator *Evaluator, input string) (string, error) {
	// Split by semicolon and evaluate each part
	parts := splitBySemicolon(input)
	var result Expr
	
	for _, part := range parts {
		part = trimSpaces(part)
		if part == "" {
			continue
		}
		
		parser := NewParser(NewLexer(part))
		expr, err := parser.Parse()
		if err != nil {
			return "", err
		}
		
		result = evaluator.Evaluate(expr)
	}
	
	if result == nil {
		return "", nil
	}
	
	return result.String(), nil
}

// Helper function to split input by semicolon
func splitBySemicolon(input string) []string {
	parts := []string{}
	current := ""
	
	for _, char := range input {
		if char == ';' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		parts = append(parts, current)
	}
	
	return parts
}

// Helper function to trim spaces
func trimSpaces(s string) string {
	start := 0
	end := len(s)
	
	for start < end && s[start] == ' ' {
		start++
	}
	
	for end > start && s[end-1] == ' ' {
		end--
	}
	
	return s[start:end]
}

func TestNamedPatterns(t *testing.T) {
	tests := []struct {
		name     string
		setup    []string
		input    string
		expected string
	}{
		{
			name: "Integer pattern matches integer",
			setup: []string{
				"process(x_Integer) := \"got integer\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process(42)",
			expected: "\"got integer\"",
		},
		{
			name: "Integer pattern doesn't match float",
			setup: []string{
				"process(x_Integer) := \"got integer\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process(3.14)",
			expected: "\"got something else\"",
		},
		{
			name: "Real pattern matches float",
			setup: []string{
				"process(x_Real) := \"got real\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process(3.14)",
			expected: "\"got real\"",
		},
		{
			name: "Real pattern doesn't match integer",
			setup: []string{
				"process(x_Real) := \"got real\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process(42)",
			expected: "\"got something else\"",
		},
		{
			name: "Number pattern matches integer",
			setup: []string{
				"process(x_Number) := \"got number\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process(42)",
			expected: "\"got number\"",
		},
		{
			name: "Number pattern matches float",
			setup: []string{
				"process(x_Number) := \"got number\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process(3.14)",
			expected: "\"got number\"",
		},
		{
			name: "Symbol pattern matches symbol",
			setup: []string{
				"process(x_Symbol) := \"got symbol\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process(abc)",
			expected: "\"got symbol\"",
		},
		{
			name: "Symbol pattern doesn't match integer",
			setup: []string{
				"process(x_Symbol) := \"got symbol\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process(42)",
			expected: "\"got something else\"",
		},
		{
			name: "String pattern matches string",
			setup: []string{
				"process(x_String) := \"got string\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process(\"hello\")",
			expected: "\"got string\"",
		},
		{
			name: "Boolean pattern matches boolean",
			setup: []string{
				"process(x_Boolean) := \"got boolean\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process(True)",
			expected: "\"got boolean\"",
		},
		{
			name: "List pattern matches list",
			setup: []string{
				"process(x_List) := \"got list\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process([1, 2, 3])",
			expected: "\"got list\"",
		},
		{
			name: "Atom pattern matches atom",
			setup: []string{
				"process(x_Atom) := \"got atom\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process(42)",
			expected: "\"got atom\"",
		},
		{
			name: "Atom pattern doesn't match list",
			setup: []string{
				"process(x_Atom) := \"got atom\"",
				"process(x_) := \"got something else\"",
			},
			input:    "process([1, 2])",
			expected: "\"got something else\"",
		},
		{
			name: "Multiple named patterns",
			setup: []string{
				"combine(x_Integer, y_Integer) := x + y",
				"combine(x_String, y_String) := x",
				"combine(x_, y_) := \"mixed types\"",
			},
			input:    "combine(5, 10)",
			expected: "15",
		},
		{
			name: "Multiple named patterns - string case",
			setup: []string{
				"combine(x_Integer, y_Integer) := x + y",
				"combine(x_String, y_String) := x",
				"combine(x_, y_) := \"mixed types\"",
			},
			input:    "combine(\"hello\", \"world\")",
			expected: "\"hello\"",
		},
		{
			name: "Multiple named patterns - mixed case",
			setup: []string{
				"combine(x_Integer, y_Integer) := x + y",
				"combine(x_String, y_String) := x",
				"combine(x_, y_) := \"mixed types\"",
			},
			input:    "combine(5, \"hello\")",
			expected: "\"mixed types\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			
			// Setup function definitions
			for _, setup := range tt.setup {
				_, err := evaluateExpression(evaluator, setup)
				if err != nil {
					t.Fatalf("Error setting up expression '%s': %v", setup, err)
				}
			}
			
			// Test the actual expression
			result, err := evaluateExpression(evaluator, tt.input)
			if err != nil {
				t.Fatalf("Error evaluating expression: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestNamedPatternMatching(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expr     string
		expected bool
	}{
		{
			name:     "Integer pattern matches integer",
			pattern:  "x_Integer",
			expr:     "42",
			expected: true,
		},
		{
			name:     "Integer pattern doesn't match float",
			pattern:  "x_Integer",
			expr:     "3.14",
			expected: false,
		},
		{
			name:     "Real pattern matches float",
			pattern:  "x_Real",
			expr:     "3.14",
			expected: true,
		},
		{
			name:     "Real pattern doesn't match integer",
			pattern:  "x_Real",
			expr:     "42",
			expected: false,
		},
		{
			name:     "Number pattern matches integer",
			pattern:  "x_Number",
			expr:     "42",
			expected: true,
		},
		{
			name:     "Number pattern matches float",
			pattern:  "x_Number",
			expr:     "3.14",
			expected: true,
		},
		{
			name:     "Symbol pattern matches symbol",
			pattern:  "x_Symbol",
			expr:     "abc",
			expected: true,
		},
		{
			name:     "Symbol pattern doesn't match integer",
			pattern:  "x_Symbol",
			expr:     "42",
			expected: false,
		},
		{
			name:     "String pattern matches string",
			pattern:  "x_String",
			expr:     "\"hello\"",
			expected: true,
		},
		{
			name:     "Boolean pattern matches boolean",
			pattern:  "x_Boolean",
			expr:     "True",
			expected: true,
		},
		{
			name:     "List pattern matches list",
			pattern:  "x_List",
			expr:     "[1, 2, 3]",
			expected: true,
		},
		{
			name:     "Atom pattern matches atom",
			pattern:  "x_Atom",
			expr:     "42",
			expected: true,
		},
		{
			name:     "Atom pattern doesn't match list",
			pattern:  "x_Atom",
			expr:     "[1, 2]",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			ctx := NewContext()
			
			// Parse pattern and expression
			parser := NewParser(NewLexer(tt.pattern))
			pattern, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing pattern: %v", err)
			}
			
			parser = NewParser(NewLexer(tt.expr))
			expr, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing expression: %v", err)
			}
			
			// Test pattern matching (use parameter context for pattern variables)
			result := evaluator.matchPatternInternal(pattern, expr, ctx, true)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestNamedPatternVariableBinding(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expr     string
		varName  string
		expected string
	}{
		{
			name:     "Integer pattern binds value",
			pattern:  "x_Integer",
			expr:     "42",
			varName:  "x",
			expected: "42",
		},
		{
			name:     "Real pattern binds value",
			pattern:  "y_Real",
			expr:     "3.14",
			varName:  "y",
			expected: "3.14",
		},
		{
			name:     "Symbol pattern binds value",
			pattern:  "z_Symbol",
			expr:     "abc",
			varName:  "z",
			expected: "abc",
		},
		{
			name:     "String pattern binds value",
			pattern:  "s_String",
			expr:     "\"hello\"",
			varName:  "s",
			expected: "\"hello\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			ctx := NewContext()
			
			// Parse pattern and expression
			parser := NewParser(NewLexer(tt.pattern))
			pattern, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing pattern: %v", err)
			}
			
			parser = NewParser(NewLexer(tt.expr))
			expr, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing expression: %v", err)
			}
			
			// Test pattern matching (use parameter context for pattern variables)
			matched := evaluator.matchPatternInternal(pattern, expr, ctx, true)
			if !matched {
				t.Fatalf("Pattern should have matched")
			}
			
			// Check if variable was bound correctly
			if value, exists := ctx.Get(tt.varName); exists {
				if value.String() != tt.expected {
					t.Errorf("Expected variable %s to be bound to %s, got %s", tt.varName, tt.expected, value.String())
				}
			} else {
				t.Errorf("Variable %s should have been bound", tt.varName)
			}
		})
	}
}

func TestSameFunctionNameDifferentPatterns(t *testing.T) {
	tests := []struct {
		name     string
		setup    []string
		input    string
		expected string
		checkType func(string) bool
	}{
		{
			name: "double function - integer input",
			setup: []string{
				"double(x_Integer) := x * 2",
				"double(x_Real) := x * 2.0",
				"double(x_) := \"cannot double\"",
			},
			input:    "double(5)",
			expected: "10",
			checkType: func(result string) bool {
				// Should be integer (no decimal point)
				return result == "10"
			},
		},
		{
			name: "double function - real input",
			setup: []string{
				"double(x_Integer) := x * 2",
				"double(x_Real) := x * 2.0",
				"double(x_) := \"cannot double\"",
			},
			input:    "double(3.5)",
			expected: "7",
			checkType: func(result string) bool {
				// Should be real (has decimal point)
				return result == "7"
			},
		},
		{
			name: "double function - string input (fallback)",
			setup: []string{
				"double(x_Integer) := x * 2",
				"double(x_Real) := x * 2.0",
				"double(x_) := \"cannot double\"",
			},
			input:    "double(\"hello\")",
			expected: "\"cannot double\"",
			checkType: func(result string) bool {
				// Should be string
				return result == "\"cannot double\""
			},
		},
		{
			name: "process function - different types",
			setup: []string{
				"process(x_Integer) := x + 100",
				"process(x_String) := \"processed: \" + x",
				"process(x_Boolean) := \"bool: \" + x",
				"process(x_) := \"unknown type\"",
			},
			input:    "process(42)",
			expected: "142",
			checkType: func(result string) bool {
				return result == "142"
			},
		},
		{
			name: "process function - string type",
			setup: []string{
				"process(x_Integer) := x + 100",
				"process(x_String) := \"processed: \" + x",
				"process(x_Boolean) := \"bool: \" + x",
				"process(x_) := \"unknown type\"",
			},
			input:    "process(\"test\")",
			expected: "Plus(\"processed: \", \"test\")",
			checkType: func(result string) bool {
				return result == "Plus(\"processed: \", \"test\")"
			},
		},
		{
			name: "process function - boolean type",
			setup: []string{
				"process(x_Integer) := x + 100",
				"process(x_String) := \"processed: \" + x",
				"process(x_Boolean) := \"bool: \" + x",
				"process(x_) := \"unknown type\"",
			},
			input:    "process(True)",
			expected: "Plus(\"bool: \", True)",
			checkType: func(result string) bool {
				return result == "Plus(\"bool: \", True)"
			},
		},
		{
			name: "abs function - positive integer",
			setup: []string{
				"abs(x_Integer) := If(x < 0, Minus(x), x)",
				"abs(x_Real) := If(x < 0.0, Minus(x), x)",
				"abs(x_) := \"not a number\"",
			},
			input:    "abs(5)",
			expected: "5",
			checkType: func(result string) bool {
				return result == "5"
			},
		},
		{
			name: "abs function - positive real",
			setup: []string{
				"abs(x_Integer) := If(x < 0, Minus(x), x)",
				"abs(x_Real) := If(x < 0.0, Minus(x), x)",
				"abs(x_) := \"not a number\"",
			},
			input:    "abs(2.5)",
			expected: "2.5",
			checkType: func(result string) bool {
				return result == "2.5"
			},
		},
		{
			name: "abs function - non-numeric fallback",
			setup: []string{
				"abs(x_Integer) := If(x < 0, Minus(x), x)",
				"abs(x_Real) := If(x < 0.0, Minus(x), x)",
				"abs(x_) := \"not a number\"",
			},
			input:    "abs(\"hello\")",
			expected: "\"not a number\"",
			checkType: func(result string) bool {
				return result == "\"not a number\""
			},
		},
		{
			name: "negate function - integer vs real",
			setup: []string{
				"negate(x_Integer) := Minus(x)",
				"negate(x_Real) := Minus(x)",
				"negate(x_) := \"cannot negate\"",
			},
			input:    "negate(42)",
			expected: "Minus(42)",
			checkType: func(result string) bool {
				return result == "Minus(42)"
			},
		},
		{
			name: "negate function - real input",
			setup: []string{
				"negate(x_Integer) := Minus(x)",
				"negate(x_Real) := Minus(x)",
				"negate(x_) := \"cannot negate\"",
			},
			input:    "negate(3.14)",
			expected: "Minus(3.14)",
			checkType: func(result string) bool {
				return result == "Minus(3.14)"
			},
		},
		{
			name: "convert function - multiple numeric types",
			setup: []string{
				"convert(x_Integer) := \"int: \" + x",
				"convert(x_Real) := \"real: \" + x", 
				"convert(x_Number) := \"number: \" + x",
				"convert(x_) := \"other: \" + x",
			},
			input:    "convert(42)",
			expected: "Plus(\"int: \", 42)",
			checkType: func(result string) bool {
				// Should match most specific pattern (Integer before Number)
				return result == "Plus(\"int: \", 42)"
			},
		},
		{
			name: "convert function - real matches Real not Number",
			setup: []string{
				"convert(x_Integer) := \"int: \" + x",
				"convert(x_Real) := \"real: \" + x", 
				"convert(x_Number) := \"number: \" + x",
				"convert(x_) := \"other: \" + x",
			},
			input:    "convert(3.14)",
			expected: "Plus(\"real: \", 3.14)",
			checkType: func(result string) bool {
				// Should match most specific pattern (Real before Number)
				return result == "Plus(\"real: \", 3.14)"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			
			// Setup function definitions
			for _, setup := range tt.setup {
				_, err := evaluateExpression(evaluator, setup)
				if err != nil {
					t.Fatalf("Error setting up expression '%s': %v", setup, err)
				}
			}
			
			// Test the actual expression
			result, err := evaluateExpression(evaluator, tt.input)
			if err != nil {
				t.Fatalf("Error evaluating expression: %v", err)
			}
			
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
			
			// Additional type check if provided
			if tt.checkType != nil && !tt.checkType(result) {
				t.Errorf("Type check failed for result: %s", result)
			}
		})
	}
}

func TestNegativeNumberLiterals(t *testing.T) {
	tests := []struct {
		name     string
		setup    []string
		input    string
		expected string
	}{
		{
			name: "Negative integer literal matches Integer pattern",
			setup: []string{
				"test(x_Integer) := \"integer: \" + x",
				"test(x_) := \"not integer\"",
			},
			input:    "test(-5)",
			expected: "Plus(\"integer: \", -5)",
		},
		{
			name: "Negative real literal matches Real pattern",
			setup: []string{
				"test(x_Real) := \"real: \" + x",
				"test(x_) := \"not real\"",
			},
			input:    "test(-3.14)",
			expected: "Plus(\"real: \", -3.14)",
		},
		{
			name: "Negative integer vs positive integer - both match Integer",
			setup: []string{
				"sign(x_Integer) := If(x < 0, \"negative\", \"positive\")",
				"sign(x_) := \"not integer\"",
			},
			input:    "sign(-42)",
			expected: "\"negative\"",
		},
		{
			name: "Positive integer still matches Integer pattern",
			setup: []string{
				"sign(x_Integer) := If(x < 0, \"negative\", \"positive\")",
				"sign(x_) := \"not integer\"",
			},
			input:    "sign(42)",
			expected: "\"positive\"",
		},
		{
			name: "Minus expression with variable should not match Integer",
			setup: []string{
				"test(x_Integer) := \"integer: \" + x",
				"test(x_) := \"not integer: \" + x",
			},
			input:    "test(Minus(y))",
			expected: "Plus(\"not integer: \", Minus(y))",
		},
		{
			name: "Number pattern matches both positive and negative",
			setup: []string{
				"abs_val(x_Number) := If(x < 0, Minus(x), x)",
				"abs_val(x_) := \"not a number\"",
			},
			input:    "abs_val(-7)",
			expected: "Minus(-7)",
		},
		{
			name: "Number pattern matches negative real",
			setup: []string{
				"abs_val(x_Number) := If(x < 0, Minus(x), x)",
				"abs_val(x_) := \"not a number\"",
			},
			input:    "abs_val(-2.5)",
			expected: "Minus(-2.5)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			
			// Setup function definitions
			for _, setup := range tt.setup {
				_, err := evaluateExpression(evaluator, setup)
				if err != nil {
					t.Fatalf("Error setting up expression '%s': %v", setup, err)
				}
			}
			
			// Test the actual expression
			result, err := evaluateExpression(evaluator, tt.input)
			if err != nil {
				t.Fatalf("Error evaluating expression: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestHeadPatterns(t *testing.T) {
	tests := []struct {
		name     string
		setup    []string
		input    string
		expected string
	}{
		{
			name: "Plus head pattern with symbols",
			setup: []string{
				"simplify(Plus(x_, y_)) := x + y",
				"simplify(expr_) := expr",
			},
			input:    "simplify(Plus(a, b))",
			expected: "Plus(a, b)",
		},
		{
			name: "Times head pattern with symbols",
			setup: []string{
				"simplify(Plus(x_, y_)) := x + y",
				"simplify(Times(x_, y_)) := x * y",
				"simplify(expr_) := expr",
			},
			input:    "simplify(Times(a, b))",
			expected: "Times(a, b)",
		},
		{
			name: "Head pattern with literal values",
			setup: []string{
				"identity(Plus(0, x_)) := x",
				"identity(Times(1, x_)) := x",
				"identity(Times(0, x_)) := 0",
				"identity(expr_) := expr",
			},
			input:    "identity(Plus(0, a))",
			expected: "a",
		},
		{
			name: "Head pattern with different structure",
			setup: []string{
				"identity(Plus(0, x_)) := x",
				"identity(Times(1, x_)) := x",
				"identity(Times(0, x_)) := 0",
				"identity(expr_) := expr",
			},
			input:    "identity(Times(1, b))",
			expected: "b",
		},
		{
			name: "Head pattern with zero multiplication",
			setup: []string{
				"identity(Plus(0, x_)) := x",
				"identity(Times(1, x_)) := x",
				"identity(Times(0, x_)) := 0",
				"identity(expr_) := expr",
			},
			input:    "identity(Times(0, c))",
			expected: "0",
		},
		{
			name: "Fallback pattern for non-matching head",
			setup: []string{
				"identity(Plus(0, x_)) := x",
				"identity(Times(1, x_)) := x",
				"identity(Times(0, x_)) := 0",
				"identity(expr_) := expr",
			},
			input:    "identity(Power(2, 3))",
			expected: "8",
		},
		{
			name: "Fallback pattern for non-matching structure",
			setup: []string{
				"identity(Plus(0, x_)) := x",
				"identity(Times(1, x_)) := x",
				"identity(Times(0, x_)) := 0",
				"identity(expr_) := expr",
			},
			input:    "identity(Plus(1, 2))",
			expected: "3",
		},
		{
			name: "Nested head patterns",
			setup: []string{
				"expand(Plus(Times(x_, y_), z_)) := x * y + z",
				"expand(Times(Plus(x_, y_), z_)) := x * z + y * z",
				"expand(expr_) := expr",
			},
			input:    "expand(Plus(Times(a, b), c))",
			expected: "Plus(Times(a, b), c)",
		},
		{
			name: "Nested head patterns - distributive",
			setup: []string{
				"expand(Plus(Times(x_, y_), z_)) := x * y + z",
				"expand(Times(Plus(x_, y_), z_)) := x * z + y * z",
				"expand(expr_) := expr",
			},
			input:    "expand(Times(Plus(a, b), c))",
			expected: "Plus(Times(a, c), Times(b, c))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			
			// Setup function definitions
			for _, setup := range tt.setup {
				_, err := evaluateExpression(evaluator, setup)
				if err != nil {
					t.Fatalf("Error setting up expression '%s': %v", setup, err)
				}
			}
			
			// Test the actual expression
			result, err := evaluateExpression(evaluator, tt.input)
			if err != nil {
				t.Fatalf("Error evaluating expression: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestHeadPatternMatching(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expr     string
		expected bool
	}{
		{
			name:     "Plus head pattern matches Plus expression",
			pattern:  "Plus(x_, y_)",
			expr:     "Plus(a, b)",
			expected: true,
		},
		{
			name:     "Times head pattern matches Times expression",
			pattern:  "Times(x_, y_)",
			expr:     "Times(a, b)",
			expected: true,
		},
		{
			name:     "Plus head pattern doesn't match Times expression",
			pattern:  "Plus(x_, y_)",
			expr:     "Times(a, b)",
			expected: false,
		},
		{
			name:     "Head pattern with literal matches",
			pattern:  "Plus(0, x_)",
			expr:     "Plus(0, a)",
			expected: true,
		},
		{
			name:     "Head pattern with literal doesn't match different literal",
			pattern:  "Plus(0, x_)",
			expr:     "Plus(1, a)",
			expected: false,
		},
		{
			name:     "Head pattern with wrong arity doesn't match",
			pattern:  "Plus(x_, y_)",
			expr:     "Plus(a, b, c)",
			expected: false,
		},
		{
			name:     "Nested head pattern matches",
			pattern:  "Plus(Times(x_, y_), z_)",
			expr:     "Plus(Times(a, b), c)",
			expected: true,
		},
		{
			name:     "Nested head pattern doesn't match different structure",
			pattern:  "Plus(Times(x_, y_), z_)",
			expr:     "Plus(Power(a, b), c)",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			ctx := NewContext()
			
			// Parse pattern and expression
			parser := NewParser(NewLexer(tt.pattern))
			pattern, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing pattern: %v", err)
			}
			
			parser = NewParser(NewLexer(tt.expr))
			expr, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing expression: %v", err)
			}
			
			// Test pattern matching
			result := evaluator.matchPattern(pattern, expr, ctx)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestHeadPatternVariableBinding(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expr     string
		bindings map[string]string
	}{
		{
			name:    "Plus head pattern binds variables",
			pattern: "Plus(x_, y_)",
			expr:    "Plus(a, b)",
			bindings: map[string]string{
				"x": "a",
				"y": "b",
			},
		},
		{
			name:    "Head pattern with literal binds variable",
			pattern: "Plus(0, x_)",
			expr:    "Plus(0, a)",
			bindings: map[string]string{
				"x": "a",
			},
		},
		{
			name:    "Nested head pattern binds variables",
			pattern: "Plus(Times(x_, y_), z_)",
			expr:    "Plus(Times(a, b), c)",
			bindings: map[string]string{
				"x": "a",
				"y": "b",
				"z": "c",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			ctx := NewContext()
			
			// Parse pattern and expression
			parser := NewParser(NewLexer(tt.pattern))
			pattern, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing pattern: %v", err)
			}
			
			parser = NewParser(NewLexer(tt.expr))
			expr, err := parser.Parse()
			if err != nil {
				t.Fatalf("Error parsing expression: %v", err)
			}
			
			// Test pattern matching (use parameter context for regular symbols)
			matched := evaluator.matchPatternInternal(pattern, expr, ctx, true)
			if !matched {
				t.Fatalf("Pattern should have matched")
			}
			
			// Check if variables were bound correctly
			for varName, expectedValue := range tt.bindings {
				if value, exists := ctx.Get(varName); exists {
					if value.String() != expectedValue {
						t.Errorf("Expected variable %s to be bound to %s, got %s", varName, expectedValue, value.String())
					}
				} else {
					t.Errorf("Variable %s should have been bound", varName)
				}
			}
		})
	}
}