package integration

import (
	"testing"
)

func TestNot(t *testing.T) {
	runTestCases(t, []TestCase{
		{name: "Not True", input: "Not(True)", expected: "False"},
		{name: "Not False", input: "Not(False)", expected: "True"},
		{name: "Not bool", input: "Not(1)", expected: "Not(1)"}, // returns unevaluated
	})
}
