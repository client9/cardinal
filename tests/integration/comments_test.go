package integration

import (
	"testing"
	
	"github.com/client9/sexpr"
)

func TestHashComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple inline comment",
			input:    "5 # this is a comment",
			expected: "5",
		},
		{
			name:     "Comment at start of line",
			input:    "# comment\n5",
			expected: "5",
		},
		{
			name:     "Comment only line",
			input:    "# just a comment",
			expected: "Null",
		},
		{
			name:     "Empty comment",
			input:    "#",
			expected: "Null",
		},
		{
			name:     "Assignment with comment",
			input:    "x = 42 # assign value",
			expected: "42",
		},
		{
			name:     "Expression with comment",
			input:    "Plus(1, 2) # addition",
			expected: "3",
		},
		{
			name:     "Multiple comments",
			input:    "# first comment\nx = 5 # inline comment\n# final comment\nx",
			expected: "5",
		},
		{
			name:     "Comment in list",
			input:    "[1, 2, 3] # list comment",
			expected: "List(1, 2, 3)",
		},
		{
			name:     "Comment in function call",
			input:    "Times(2, 3) # multiplication",
			expected: "6",
		},
		{
			name:     "Comment with special characters",
			input:    "10 # comment with !@#$%^&*()_+-={}[]|\\:;\"'<>?,./~`",
			expected: "10",
		},
	}

	runTestCases(t, tests)
}

func TestCommentsInComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Compound statement with comments",
			input:    "x = 1; # first assignment\ny = 2; # second assignment\nx + y",
			expected: "3",
		},
		{
			name:     "Rule with comment",
			input:    "rule = a:b; # create rule\nFullForm(rule)",
			expected: "\"Rule(a, b)\"",
		},
		{
			name:     "Nested expression with comments",
			input:    "Times(Plus(1, 2), # inner expression\nPlus(3, 4)) # outer expression",
			expected: "21",
		},
		{
			name:     "Pattern with comment",
			input:    "pattern = x_; # pattern variable\nFullForm(pattern)",
			expected: "\"Pattern(x, Blank())\"",
		},
	}

	runTestCases(t, tests)
}

func TestCommentsDoNotAffectTokenization(t *testing.T) {
	tests := []struct {
		name        string
		withComment string
		without     string
	}{
		{
			name:        "Number with comment",
			withComment: "42 # comment",
			without:     "42",
		},
		{
			name:        "String with comment",
			withComment: "\"hello\" # comment",
			without:     "\"hello\"",
		},
		{
			name:        "Symbol with comment",
			withComment: "myVar # comment",
			without:     "myVar",
		},
		{
			name:        "Operator with comment",
			withComment: "1 + 2 # comment",
			without:     "1 + 2",
		},
		{
			name:        "Function call with comment",
			withComment: "Plus(1, 2) # comment",
			without:     "Plus(1, 2)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Both expressions should produce the same result
			evaluateAndExpect(t, test.withComment, evaluateString(test.without))
		})
	}
}

// Helper function to evaluate a string and return the result
func evaluateString(input string) string {
	eval := sexpr.NewEvaluator()
	expr, err := sexpr.ParseString(input)
	if err != nil {
		return "ERROR"
	}
	result := eval.Evaluate(expr)
	return result.String()
}