package integration

import (
	"testing"
)

func TestUnset(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Set Unset",
			input:    "Set(x, 42); Unset(x); x",
			expected: "x",
		},
		{
			name:     "Unset nonexistant",
			input:    "Unset(x)",
			expected: "Null",
		},
		{
			name:      "Unset non-symbol",
			input:     "Unset(100)",
			expected:  "",
			errorType: "ArgumentError",
		},

		// TODO: need better error handling
		{
			name:      "Unset empty",
			input:     "Unset()",
			expected:  "",
			errorType: "???",
			skip:      true,
		},
		{
			name:      "Unset Protected",
			input:     "Unset(Plus)",
			expected:  "",
			errorType: "Protected",
		},
	}

	runTestCases(t, tests)
}
