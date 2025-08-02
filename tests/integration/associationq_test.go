package integration

import (
	"testing"
)

func TestAssociationQ(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Empty association",
			input:    "AssociationQ({})",
			expected: "True",
		},
		{
			name:     "Non-empty association",
			input:    "AssociationQ({name: \"Bob\"})",
			expected: "True",
		},
		{
			name:     "List is not association",
			input:    "AssociationQ([1, 2, 3])",
			expected: "False",
		},
		{
			name:     "Integer is not association",
			input:    "AssociationQ(42)",
			expected: "False",
		},
		{
			name:     "String is not association",
			input:    "AssociationQ(\"test\")",
			expected: "False",
		},
	}

	runTestCases(t, tests)
}
