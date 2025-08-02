package integration

import (
	"testing"
)

func TestValues(t *testing.T) {
	runTestCases(t, []TestCase{
		{name: "Empty association", input: "Values({})", expected: "List()"},
		{name: "Single value", input: "Values({name: \"Bob\"})", expected: "List(\"Bob\")"},
		{name: "Multiple values", input: "Values({name: \"Bob\", age: 30})", expected: "List(\"Bob\", 30)"},
	})
}
