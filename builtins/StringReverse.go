package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol StringReverse

// StringReverse reverses a string
// @ExprPattern (_String)
func StringReverse(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	lhs, _ := core.ExtractString(args[0])
	runes := []rune(lhs)

	// Iterate with two pointers, one from the beginning and one from the end, swapping elements.
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Convert the reversed rune slice back to a string.
	return core.NewString(string(runes))
}
