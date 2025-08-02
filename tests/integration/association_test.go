package integration

import (
	"testing"
)

func TestAssociation(t *testing.T) {
	runTestCases(t, []TestCase{
		{
			name:     "Empty association",
			input:    "{}",
			expected: "Association()",
		},
		{
			name:     "Single key-value pair",
			input:    "{ name: \"Bob\"}",
			expected: "Association(Rule(name, \"Bob\"))",
		},
		{
			name:     "Multiple key-value pairs",
			input:    "{ name: \"Bob\", age: 30}",
			expected: "Association(Rule(name, \"Bob\"), Rule(age, 30))",
		},
		{
			name:     "Mixed value types",
			input:    "{ name: \"Bob\", age: 30, active: True}",
			expected: "Association(Rule(name, \"Bob\"), Rule(age, 30), Rule(active, True))",
		},
		{
			name:     "Trailing comma",
			input:    "{ name: \"Bob\", age: 30,}",
			expected: "Association(Rule(name, \"Bob\"), Rule(age, 30))",
		},
		{
			name:     "Empty association",
			input:    "Length({})",
			expected: "0",
		},
		{
			name:     "Single item",
			input:    "Length({ name: \"Bob\"})",
			expected: "1",
		},
		{
			name:     "Multiple items",
			input:    "Length({ name: \"Bob\", age: 30, active: True})",
			expected: "3",
		},
		{
			name:     "Access existing key",
			input:    "Part({ name: \"Bob\", age: 30}, name)",
			expected: "\"Bob\"",
		},
		{
			name:     "Access another existing key",
			input:    "Part({ name: \"Bob\", age: 30}, age)",
			expected: "30",
		},
		{
			name:     "Access with string key",
			input:    "Part({\"key\": \"value\"}, \"key\")",
			expected: "\"value\"",
		},
		{
			name:     "Equal empty associations",
			input:    "SameQ({}, {})",
			expected: "True",
		},
		{
			name:     "Equal single-item associations",
			input:    "SameQ({ name: \"Bob\"}, { name: \"Bob\"})",
			expected: "True",
		},
		{
			name:     "Different associations",
			input:    "SameQ({ name: \"Bob\"}, { name: \"Alice\"})",
			expected: "False",
		},
		{
			name:     "Different key sets",
			input:    "SameQ({ name: \"Bob\"}, {age: 30})",
			expected: "False",
		},
		{
			name:     "Insertion order preserved - keys",
			input:    "Keys({c:3, b:2, a:1})",
			expected: "List(c, b, a)",
		},
		{
			name:     "Insertion order preserved - values",
			input:    "Values({c:3, b:2, a:1})",
			expected: "List(3, 2, 1)",
		},
		{
			name:     "InputForm - empty",
			input:    "InputForm(Association())",
			expected: "\"{}\"",
		},
		{
			name:     "InputForm - full",
			input:    "InputForm(Association(Rule(a,1), Rule(b,2)))",
			expected: "\"{a: 1, b: 2}\"",
		},
		{
			name:     "Add with Part syntax",
			input:    "m = Association(Rule(a,1), Rule(b,2)); Part(m, c) = 3",
			expected: "List(a, b, c)",
			skip:     true,
		},
		{
			name:     "Strings and Symbols are different keys",
			input:    "Assert(Length(Keys({a:1,\"a\":2})) == 2)",
			expected: "True",
		},
		{
			name:     "Add with Part syntax",
			input:    "m = Association(Rule(a,1), Rule(b,2)); Part(m, c) = 3",
			expected: "List(a, b, c)",
			skip:     true,
		},
		{
			name:     "Add with Part syntax",
			input:    "m = {a:1,b:2}; Part(m, c) = 3",
			expected: "List(a, b, c)",
			skip:     true,
		},
		/*
			{
				name:		 "Add with slice syntax",
				input:				 "map = {a:1,b:2}; map[c] = 3; Keys(map)",
				expected:		 "List(a, b)",
							},
		*/
	})
}
