package stdlib

import (
	"github.com/client9/cardinal/core"
)

// ByteArrayReverse reverses a byte array
func ByteArrayReverse(lhs []byte) []byte {
	end := len(lhs) - 1
	ba := make([]byte, len(lhs))
	for i, b := range lhs {
		ba[end-i] = b
	}
	return ba
}

// ByteArrayAppend adds one byte to the end of a ByteArray
func ByteArrayAppend(lhs core.ByteArray, b int64) core.Expr {
	return lhs.Append(byte(b))
}

// NewByteArrayFromInts creates a new ByteArray from a slice of int64
//
//	values are cast to a byte.
func ByteArrayFromInts(src []int64) core.Expr {
	dest := make([]byte, len(src))
	for i, val := range src {
		dest[i] = byte(val)
	}
	// NewByteArray makes another un-needed copy
	return core.NewByteArray(dest)
}

// ByteArrayToString converts a ByteArray to a string (assuming UTF-8 encoding)
func ByteArrayToString(ba core.ByteArray) core.Expr {
	return ba.ToStringAtom()
}
