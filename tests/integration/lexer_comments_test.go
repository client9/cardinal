package integration

import (
	"testing"
)

func TestCommentsParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Comment should not affect parsing result",
			input:    "Plus(1, 2) # this is a comment",
			expected: "3",
		},
		{
			name:     "Comment at start of expression",
			input:    "# comment\nTimes(3, 4)",
			expected: "12",
		},
		{
			name:     "Empty line with comment",
			input:    "# just a comment",
			expected: "Null",
		},
		{
			name:     "Multiple line comments",
			input:    "# first comment\n# second comment\n5",
			expected: "5",
		},
		{
			name:     "Interleaved comments and code",
			input:    "x = 1; # assign\n# comment line\ny = 2; # another assign\nx + y",
			expected: "3",
		},
	}

	runTestCases(t, tests)
}

// Test that comments work the same as if they weren't there
func TestCommentsEquivalence(t *testing.T) {
	testPairs := []struct {
		name            string
		withComments    string
		withoutComments string
	}{
		{
			name:            "Simple arithmetic",
			withComments:    "1 + 2 # addition",
			withoutComments: "1 + 2",
		},
		{
			name:            "Function call",
			withComments:    "Power(2, 3) # exponentiation",
			withoutComments: "Power(2, 3)",
		},
		{
			name:            "Assignment",
			withComments:    "x = 42 # the answer",
			withoutComments: "x = 42",
		},
		{
			name:            "List creation",
			withComments:    "[1, 2, 3] # numbers",
			withoutComments: "[1, 2, 3]",
		},
	}

	for _, pair := range testPairs {
		t.Run(pair.name, func(t *testing.T) {
			result1 := evaluateString(pair.withComments)
			result2 := evaluateString(pair.withoutComments)

			if result1 != result2 {
				t.Errorf("Results should be identical:\nWith comments: %s\nWithout comments: %s",
					result1, result2)
			}
		})
	}
}
