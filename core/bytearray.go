package core

import (
	"bytes"
	"fmt"
)

// ByteArray represents an array of bytes as a first-class expression type
type ByteArray struct {
	data []byte
}

// NewByteArray creates a new ByteArray from a byte slice
func NewByteArray(data []byte) ByteArray {
	// Make a copy to ensure immutability
	copied := make([]byte, len(data))
	copy(copied, data)
	return ByteArray{data: copied}
}

// NewByteArrayFromString creates a new ByteArray from a string
func NewByteArrayFromString(s string) ByteArray {
	return NewByteArray([]byte(s))
}

// Data returns a copy of the underlying byte data
func (b ByteArray) Data() []byte {
	copied := make([]byte, len(b.data))
	copy(copied, b.data)
	return copied
}

// Expr interface implementation

func (b ByteArray) Length() int64 {
	return int64(len(b.data))
}

func (b ByteArray) String() string {
	if len(b.data) <= 20 {
		return fmt.Sprintf("ByteArray[%v]", b.data)
	}
	return fmt.Sprintf("ByteArray[%v...]", b.data[:20])
}

func (b ByteArray) InputForm() string {
	return b.String()
}

func (b ByteArray) Type() string {
	return "ByteArray"
}

func (b ByteArray) Equal(rhs Expr) bool {
	rhsByteArray, ok := rhs.(ByteArray)
	if !ok {
		return false
	}
	return bytes.Equal(b.data, rhsByteArray.data)
}

// Sliceable interface implementation

// ElementAt returns the nth byte (1-indexed) as an integer atom
func (b ByteArray) ElementAt(n int64) Expr {
	length := int64(len(b.data))
	
	if length == 0 {
		return NewErrorExpr("PartError", "ByteArray is empty", []Expr{b})
	}
	
	// Handle negative indexing
	if n < 0 {
		n = length + n + 1
	}
	
	// Check bounds (1-indexed)
	if n <= 0 || n > length {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Part index %d is out of bounds for ByteArray with %d bytes", n, length),
			[]Expr{b})
	}
	
	// Convert to 0-based index and return byte as integer
	return NewIntAtom(int(b.data[n-1]))
}

// Slice returns a new ByteArray containing bytes from start to stop (inclusive, 1-indexed)
func (b ByteArray) Slice(start, stop int64) Expr {
	length := int64(len(b.data))
	
	if length == 0 {
		return NewByteArray(nil)
	}
	
	// Handle negative indexing
	if start < 0 {
		start = length + start + 1
	}
	if stop < 0 {
		stop = length + stop + 1
	}
	
	// Check bounds
	if start <= 0 || stop <= 0 || start > length || stop > length {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Slice indices [%d, %d] out of bounds for ByteArray with %d bytes",
				start, stop, length), []Expr{b})
	}
	
	if start > stop {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Start index %d is greater than stop index %d", start, stop),
			[]Expr{b})
	}
	
	// Convert to 0-based indices and create new ByteArray
	startIdx := start - 1
	stopIdx := stop // stop is inclusive, so we include it
	
	return NewByteArray(b.data[startIdx:stopIdx])
}

// Utility methods

// Append creates a new ByteArray with additional bytes appended
func (b ByteArray) Append(data ...byte) ByteArray {
	newData := make([]byte, len(b.data)+len(data))
	copy(newData, b.data)
	copy(newData[len(b.data):], data)
	return ByteArray{data: newData}
}

// ToStringAtom converts the ByteArray to a string atom (assuming UTF-8 encoding)
func (b ByteArray) ToStringAtom() Atom {
	return NewStringAtom(string(b.data))
}
