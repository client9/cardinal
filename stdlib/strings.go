package stdlib

import (
	"unicode/utf8"

//	"github.com/client9/sexpr/core"
)

// String manipulation functions

// StringLengthStr returns the UTF-8 rune count of a string
func StringLengthRunes(s string) int64 {
	return int64(utf8.RuneCountInString(s))
}
