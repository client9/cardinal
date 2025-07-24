package stdlib

import (
	"unicode/utf8"

	"github.com/client9/sexpr/core"
)

// String utility functions

// StringLengthRunes returns the UTF-8 rune count of a string
func StringLengthRunes(s string) int64 {
	return int64(utf8.RuneCountInString(s))
}

// NewByteArrayFromInts creates a new ByteArray from a slice of int64
//  values are cast to a byte.
func ByteArrayFromInts(src []int64) core.Expr {
	dest := make([]byte, len(src))
	for i,val := range src {
		dest[i] = byte(val)
	}
	// NewByteArray makes another un-needed copy
	return core.NewByteArray(dest)
}

// ByteArrayFromString creates a ByteArray from a string for byte-level operations
func ByteArrayFromString(s string) core.Expr {
	return core.NewByteArray([]byte(s))
}
// ByteArrayToString converts a ByteArray to a string (assuming UTF-8 encoding)
func ByteArrayToString(ba core.ByteArray) core.Expr {
	return ba.ToStringAtom()
}
