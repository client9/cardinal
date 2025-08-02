package integration

import (
	"testing"
)

func TestSequenceSplicingBasic(t *testing.T) {
	tests := []TestCase{
		// The core issue that was reported - sequence splicing in Replace/ReplaceAll
		{
			name:     "Single sequence element",
			input:    `Replace(Zoo(1,3,x), Rule(Zoo(x_Integer, y_Integer, z__), Zoo(Plus(x,y), z)))`,
			expected: "Zoo(4, x)",
		},
		{
			name:     "Multiple sequence elements",
			input:    `Replace(Zoo(1,3,a,b,c), Rule(Zoo(x_Integer, y_Integer, z__), Zoo(Plus(x,y), z)))`,
			expected: "Zoo(4, a, b, c)",
		},
		{
			name:     "Empty sequence with BlankNullSequence",
			input:    `Replace(Zoo(1,3), Rule(Zoo(x_Integer, y_Integer, z___), Zoo(Plus(x,y), z)))`,
			expected: "Zoo(4)",
		},
		{
			name:     "ReplaceAll with sequence splicing",
			input:    `ReplaceAll(Zoo(1,3,x), Rule(Zoo(x_Integer, y_Integer, z__), Zoo(Plus(x,y), z)))`,
			expected: "Zoo(4, x)",
		},
		{
			name:     "ReplaceAll with nested sequences",
			input:    `ReplaceAll(List(Zoo(1,2,a), Zoo(3,4,b,c)), Rule(Zoo(x_Integer, y_Integer, z__), Zoo(Plus(x,y), z)))`,
			expected: "List(Zoo(3, a), Zoo(7, b, c))",
		},
		{
			name:     "Mixed regular and sequence variables",
			input:    `Replace(Zoo(42,hello,a,b,c), Rule(Zoo(n_Integer, word_, rest__), Zoo(word, n, rest)))`,
			expected: "Zoo(hello, 42, a, b, c)",
		},
	}

	runTestCases(t, tests)
}
