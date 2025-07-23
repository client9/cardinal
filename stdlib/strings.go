package stdlib

import (
	"unicode/utf8"
)

// String manipulation functions

// StringLengthFunc returns the length of a string
func StringLengthFunc(s string) int64 {
	return int64(len(s))
}

// StringLengthStr returns the UTF-8 rune count of a string
func StringLengthStr(s string) int64 {
	return int64(utf8.RuneCountInString(s))
}