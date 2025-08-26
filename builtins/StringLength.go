package builtins

import (
	"unicode/utf8"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol StringLength

// StringLengthRunes returns the UTF-8 rune count of a string
// @ExprPattern (_String)
func StringLengthRunes(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	s, _ := core.ExtractString(args[0])
	return core.NewInteger(int64(utf8.RuneCountInString(s)))
}
