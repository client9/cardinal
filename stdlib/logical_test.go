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

func TestMatchQExprs(t *testing.T) {
	tests := []struct {
		name     string
		expr     core.Expr
		pattern  core.Expr
		expected bool
	}{
		// Basic wildcard patterns
		{
			name:     "Simple wildcard match",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2)),
			pattern:  createList("Zoo", createPattern("x", createBlank()), createPattern("y", createBlank())),
			expected: true,
		},
		{
			name:     "Wildcard mismatch - different head",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2)),
			pattern:  createList("Bar", createPattern("x", createBlank()), createPattern("y", createBlank())),
			expected: false,
		},

		// Type-constrained patterns
		{
			name:     "Type constraint match",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2)),
			pattern:  createList("Zoo", createPattern("x", createBlankWithType("Integer")), createPattern("y", createBlankWithType("Integer"))),
			expected: true,
		},
		{
			name:     "Type constraint mismatch",
			expr:     createList("Zoo", core.NewInteger(1), core.NewSymbol("a")),
			pattern:  createList("Zoo", createPattern("x", createBlankWithType("Integer")), createPattern("y", createBlankWithType("Integer"))),
			expected: false,
		},

		// Sequence patterns (BlankNullSequence - z___)
		{
			name:     "BlankNullSequence with zero elements",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2)),
			pattern:  createList("Zoo", createPattern("x", createBlankWithType("Integer")), createPattern("y", createBlankWithType("Integer")), createPattern("z", createBlankNullSequence())),
			expected: true,
		},
		{
			name:     "BlankNullSequence with one element",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2), core.NewSymbol("a")),
			pattern:  createList("Zoo", createPattern("x", createBlankWithType("Integer")), createPattern("y", createBlankWithType("Integer")), createPattern("z", createBlankNullSequence())),
			expected: true,
		},
		{
			name:     "BlankNullSequence with multiple elements",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2), core.NewSymbol("a"), core.NewSymbol("b")),
			pattern:  createList("Zoo", createPattern("x", createBlankWithType("Integer")), createPattern("y", createBlankWithType("Integer")), createPattern("z", createBlankNullSequence())),
			expected: true,
		},

		// Sequence patterns with type constraints
		{
			name:     "BlankNullSequence with type constraint - success",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2), core.NewInteger(3), core.NewInteger(4)),
			pattern:  createList("Zoo", createPattern("x", createBlankWithType("Integer")), createPattern("y", createBlankWithType("Integer")), createPattern("z", createBlankNullSequenceWithType("Integer"))),
			expected: true,
		},
		{
			name:     "BlankNullSequence with type constraint - failure",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2), core.NewSymbol("a"), core.NewSymbol("b")),
			pattern:  createList("Zoo", createPattern("x", createBlankWithType("Integer")), createPattern("y", createBlankWithType("Integer")), createPattern("z", createBlankNullSequenceWithType("Integer"))),
			expected: false,
		},

		// BlankSequence patterns (z__) - must match at least one
		{
			name:     "BlankSequence with zero elements - should fail",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2)),
			pattern:  createList("Zoo", createPattern("x", createBlankWithType("Integer")), createPattern("y", createBlankWithType("Integer")), createPattern("z", createBlankSequence())),
			expected: false,
		},
		{
			name:     "BlankSequence with one element",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2), core.NewSymbol("a")),
			pattern:  createList("Zoo", createPattern("x", createBlankWithType("Integer")), createPattern("y", createBlankWithType("Integer")), createPattern("z", createBlankSequence())),
			expected: true,
		},

		// Literal patterns
		{
			name:     "Exact literal match",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2)),
			pattern:  createList("Zoo", core.NewInteger(1), core.NewInteger(2)),
			expected: true,
		},
		{
			name:     "Literal mismatch",
			expr:     createList("Zoo", core.NewInteger(1), core.NewInteger(2)),
			pattern:  createList("Zoo", core.NewInteger(1), core.NewInteger(3)),
			expected: false,
		},

		// Mixed patterns
		{
			name:     "Mixed literal and wildcard",
			expr:     createList("Zoo", core.NewInteger(1), core.NewSymbol("a")),
			pattern:  createList("Zoo", core.NewInteger(1), createPattern("x", createBlank())),
			expected: true,
		},

		// Edge cases
		{
			name:     "Empty lists",
			expr:     createList("Zoo"),
			pattern:  createList("Zoo"),
			expected: true,
		},
		{
			name:     "Pattern longer than expression",
			expr:     createList("Zoo", core.NewInteger(1)),
			pattern:  createList("Zoo", createPattern("x", createBlankWithType("Integer")), createPattern("y", createBlankWithType("Integer"))),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MatchQExprs(test.expr, test.pattern)
			if result != test.expected {
				t.Errorf("MatchQExprs(%s, %s) = %v, expected %v", 
					test.expr.String(), test.pattern.String(), result, test.expected)
			}
		})
	}
}

// Helper functions to create pattern expressions

func createList(head string, elements ...core.Expr) core.Expr {
	allElements := append([]core.Expr{core.NewSymbol(head)}, elements...)
	return core.List{Elements: allElements}
}

func createPattern(varName string, blankExpr core.Expr) core.Expr {
	return core.CreatePatternExpr(core.NewSymbol(varName), blankExpr)
}

func createBlank() core.Expr {
	return core.CreateBlankExpr(nil)
}

func createBlankWithType(typeName string) core.Expr {
	return core.CreateBlankExpr(core.NewSymbol(typeName))
}

func createBlankSequence() core.Expr {
	return core.CreateBlankSequenceExpr(nil)
}

func createBlankSequenceWithType(typeName string) core.Expr {
	return core.CreateBlankSequenceExpr(core.NewSymbol(typeName))
}

func createBlankNullSequence() core.Expr {
	return core.CreateBlankNullSequenceExpr(nil)
}

func createBlankNullSequenceWithType(typeName string) core.Expr {
	return core.CreateBlankNullSequenceExpr(core.NewSymbol(typeName))
}