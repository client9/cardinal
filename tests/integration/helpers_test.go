package integration

import (
	"testing"

	"github.com/client9/sexpr"
	"github.com/client9/sexpr/core"
)

// evaluateAndExpect is a test helper that parses input, evaluates it, and checks the result
func evaluateAndExpect(t *testing.T, input, expected string) {
	t.Helper()
	eval := sexpr.NewEvaluator()
	expr, err := sexpr.ParseString(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	result := eval.Evaluate(expr)
	if result.String() != expected {
		t.Errorf("Input: %q\nExpected: %q\nGot: %q", input, expected, result.String())
	}
}

// evaluateAndExpectError is a test helper that expects an error of a specific type
func evaluateAndExpectError(t *testing.T, input, errorType string) {
	t.Helper()
	eval := sexpr.NewEvaluator()
	expr, err := sexpr.ParseString(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	result := eval.Evaluate(expr)
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
func runTestCases(t *testing.T, tests []struct {
	name     string
	input    string
	expected string
}) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluateAndExpect(t, tt.input, tt.expected)
		})
	}
}

// runErrorTestCases runs a slice of test cases that expect errors
func runErrorTestCases(t *testing.T, tests []struct {
	name      string
	input     string
	errorType string
}) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluateAndExpectError(t, tt.input, tt.errorType)
		})
	}
}