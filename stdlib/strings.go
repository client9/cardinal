package stdlib

import (
	"unicode/utf8"
)

// String utility functions

// StringLengthRunes returns the UTF-8 rune count of a string
func StringLengthRunes(s string) int64 {
	return int64(utf8.RuneCountInString(s))
}

// StringAppend appends a string to another string.
// Note: Append normally adds an element of list to a list
//
//	We don't have "character" or "rune" type, and characters
//	are single character *strings*.
//
// So while the intention of Append is Append("foo", "d")
// it can be used for string joining Append("foo", "bar")
//
// This is added for the principle of least surprise.
func StringAppend(lhs, rhs string) string {
	return lhs + rhs
}

// StringReverse reverses a string
func StringReverse(lhs string) string {
	runes := []rune(lhs)

	// Iterate with two pointers, one from the beginning and one from the end, swapping elements.
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Convert the reversed rune slice back to a string.
	return string(runes)
}
