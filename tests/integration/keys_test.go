package integration

import (
	"testing"
)

func TestKeys(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Empty association",
			input:    "Keys({})",
			expected: "List()",
		},
		{
			name:     "Single key",
			input:    "Keys({name: \"Bob\"})",
			expected: "List(name)",
		},
		{
			name:     "Multiple keys",
			input:    "Keys({name: \"Bob\", age: 30})",
			expected: "List(name, age)",
		},
	}

	runTestCases(t, tests)
}
