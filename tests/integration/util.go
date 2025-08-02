package integration

import (
	"testing"

	"github.com/client9/sexpr"
	"github.com/client9/sexpr/core"
)

type TestCase struct {
	name     string
	input    string
	expected string

	// TODO use Enum
	errorType string

	skip      bool
}

// evaluateAndExpect is a test helper that parses input, evaluates it, and checks the result
func evaluateAndExpect(t *testing.T, tt TestCase) {
	t.Helper()

	result, err := sexpr.EvaluateString(tt.input)
	if err != nil {
		t.Errorf("%s: parse error %s", L4(tt.name), tt.input)
		return
	}
	if tt.errorType == "" {
		if core.IsError(result) {
			t.Errorf("%s: got error for %q, got: %q", L4(tt.name), tt.input, result.String())
			return
		}
		if result.String() != tt.expected {
			t.Errorf("%s: Input: %q\nExpected: %q\nGot: %q", L4(tt.name), tt.input, tt.expected, result.String())
			return
		}
		return
	}

	// expected error case

	errorExpr, ok := result.(*core.ErrorExpr)
	if !ok {
		t.Errorf("%s: expected error %s, but got ordinary result %s", L4(tt.name), tt.errorType, result.String())
		return
	}
	if errorExpr.ErrorType != tt.errorType {
		t.Errorf("%s: expected error type %q for input %q, got %q", L4(tt.name), tt.errorType, tt.input, errorExpr.ErrorType)
		return
	}
	return
}

// evaluateAndExpectError is a test helper that expects an error of a specific type
func evaluateAndExpectError(t *testing.T, input, errorType string) {
	t.Helper()
	eval := sexpr.NewEvaluator()
	ctx := eval.GetContext()
	expr, err := sexpr.ParseString(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	result := eval.Evaluate(ctx, expr)
	if !core.IsError(result) {
		t.Errorf("Expected error for %q, got: %q", input, result.String())
		return
	}
	errorExpr := result.(*core.ErrorExpr)
	if errorExpr.ErrorType != errorType {
		t.Errorf("Expected error type %q for input %q, got %q", errorType, input, errorExpr.ErrorType)
	}
}

// runTestCases runs a slice of test cases using evaluateAndExpected
func runTestCases(t *testing.T, tests []TestCase) {
	t.Helper()
	for _, tt := range tests {
		if tt.skip {
			continue	
		}
		//t.Run(tt.name, func(t *testing.T) {
			//t.Helper()
			evaluateAndExpect(t, tt)
		//})
	}
}

// runErrorTestCases runs a slice of test cases that expect errors
func runErrorTestCases(t *testing.T, tests []struct {
	name      string
	input     string
	errorType string
}) {
	t.Helper()
	for _, tt := range tests {
				evaluateAndExpectError(t, tt.input, tt.errorType)
	}
}

// Helper function to evaluate a string and return the result
func evaluateString(input string) string {
	expr, err := sexpr.EvaluateString(input)
	if err != nil {
		return "ERROR"
	}
	return expr.String()
}
