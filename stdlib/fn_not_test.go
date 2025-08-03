package stdlib

import (
	"testing"

	"github.com/client9/sexpr/core"
)

func TestNotExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    core.Expr
		expected core.Expr
	}{
		{
			name:     "Not True",
			input:    core.NewBool(true),
			expected: core.NewBool(false),
		},
		{
			name:     "Not False",
			input:    core.NewBool(false),
			expected: core.NewBool(true),
		},
		{
			name:     "Not symbol (symbolic behavior)",
			input:    core.NewSymbol("x"),
			expected: core.NewList("Not", core.NewSymbol("x")),
		},
		{
			name:     "Not number (symbolic behavior)",
			input:    core.NewInteger(42),
			expected: core.NewList("Not", core.NewInteger(42)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := NotExpr(test.input)
			if !result.Equal(test.expected) {
				t.Errorf("NotExpr(%s) = %s, expected %s",
					test.input.String(), result.String(), test.expected.String())
			}
		})
	}
}
