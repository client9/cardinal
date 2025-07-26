package core

import (
	"bytes"
	"fmt"
	"strings"
)

// ByteArray represents an array of bytes as a first-class expression type
//
// Design Note: ByteArray is implemented as a struct with a private []byte field
// rather than as "type ByteArray []byte" to ensure immutability. This design
// prevents several issues:
//
//  1. Slice sharing: Go slices share underlying arrays, so "type ByteArray []byte"
//     would allow mutations to be visible across copies
//  2. Direct access: Users could directly modify bytes with b[i] = newValue
//  3. No encapsulation: Can't control access or enforce defensive copying
//
// The struct-based approach ensures immutability through:
// - Private data field (no direct access)
// - Defensive copying in NewByteArray() constructor
// - Defensive copying in Data() accessor
// - All operations return new ByteArray instances
//
// This follows functional programming principles and provides thread-safety,
// which is essential for s-expression evaluation where data structures should
// be immutable during pattern matching and evaluation.
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
	if len(b.data) == 0 {
		return "ByteArray()"
	}

	var parts []string
	limit := len(b.data)
	truncated := false

	if limit > 20 {
		limit = 20
		truncated = true
	}

	for i := 0; i < limit; i++ {
		parts = append(parts, fmt.Sprintf("%d", b.data[i]))
	}

	result := "ByteArray(" + strings.Join(parts, ", ") + ")"
	if truncated {
		result = result[:len(result)-1] + "...)"
	}

	return result
}

func (b ByteArray) InputForm() string {
	return b.String()
}

func (b ByteArray) Type() string {
	return "ByteArray"
}

func (b ByteArray) IsAtom() bool {
	return false
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
func (b ByteArray) ToStringAtom() Expr {
	return NewStringAtom(string(b.data))
}

// Join joins this ByteArray with another sliceable of the same type
func (b ByteArray) Join(other Sliceable) Expr {
	// Type check: ensure other is also a ByteArray
	otherByteArray, ok := other.(ByteArray)
	if !ok {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Cannot join %T with ByteArray", other),
			[]Expr{b, other.(Expr)})
	}

	// Get data from both byte arrays
	thisData := b.Data()
	otherData := otherByteArray.Data()

	// Create new byte array with combined data
	newData := make([]byte, len(thisData)+len(otherData))
	copy(newData, thisData)
	copy(newData[len(thisData):], otherData)

	return NewByteArray(newData)
}

// SetElementAt returns a new ByteArray with the nth byte replaced (1-indexed)
// Returns an error Expr if index is out of bounds or value is not a valid byte
func (b ByteArray) SetElementAt(n int64, value Expr) Expr {
	// Validate that value is an integer representing a valid byte (0-255)
	byteValue, ok := value.(Integer)
	if !ok {
		return NewErrorExpr("TypeError",
			"ByteArray assignment requires integer value (0-255)", []Expr{b, value})
	}

	if byteValue < 0 || byteValue > 255 {
		return NewErrorExpr("ValueError",
			fmt.Sprintf("Byte value %d is out of range (0-255)", byteValue),
			[]Expr{b, value})
	}

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

	// Create new ByteArray with byte replaced
	newData := make([]byte, length)
	copy(newData, b.data)
	newData[n-1] = byte(byteValue) // n is 1-indexed

	return NewByteArray(newData)
}

// SetSlice returns a new ByteArray with bytes from start to stop replaced by values (1-indexed)
// values can be a ByteArray, List of integers, or single integer
func (b ByteArray) SetSlice(start, stop int64, values Expr) Expr {
	length := int64(len(b.data))

	// Handle empty ByteArray
	if length == 0 {
		if start == 1 && stop == 0 {
			// Insert at beginning of empty ByteArray
			return b.convertToByteArray(values)
		}
		return NewErrorExpr("PartError", "ByteArray is empty", []Expr{b})
	}

	// Handle negative indexing
	if start < 0 {
		start = length + start + 1
	}
	if stop < 0 {
		stop = length + stop + 1
	}

	// Validate range
	if start <= 0 {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Start index %d must be positive", start),
			[]Expr{b})
	}

	if start > length+1 {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Start index %d is out of bounds for ByteArray with %d bytes", start, length),
			[]Expr{b})
	}

	// Handle special cases
	if stop < start-1 {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Stop index %d cannot be less than start index %d - 1", stop, start),
			[]Expr{b})
	}

	// Convert values to byte slice
	valueBytes, err := b.convertToByteSlice(values)
	if err != nil {
		return err
	}

	// Calculate range to replace
	replaceStart := start - 1 // Convert to 0-based
	replaceEnd := stop        // Convert to 0-based (exclusive end)
	if stop > length {
		replaceEnd = length
	}

	// Calculate size for new ByteArray
	oldRangeSize := replaceEnd - replaceStart
	newSize := length - oldRangeSize + int64(len(valueBytes))

	// Build new ByteArray
	newData := make([]byte, newSize)

	// Copy bytes before the replacement range
	if replaceStart > 0 {
		copy(newData[:replaceStart], b.data[:replaceStart])
	}

	// Copy replacement bytes
	copy(newData[replaceStart:replaceStart+int64(len(valueBytes))], valueBytes)

	// Copy bytes after the replacement range
	if replaceEnd < length {
		afterStart := replaceStart + int64(len(valueBytes))
		copy(newData[afterStart:], b.data[replaceEnd:])
	}

	return NewByteArray(newData)
}

// convertToByteArray converts an Expr to a ByteArray
func (b ByteArray) convertToByteArray(values Expr) Expr {
	if byteArray, ok := values.(ByteArray); ok {
		return byteArray
	}

	valueBytes, err := b.convertToByteSlice(values)
	if err != nil {
		return err
	}

	return NewByteArray(valueBytes)
}

// convertToByteSlice converts an Expr to a byte slice
func (b ByteArray) convertToByteSlice(values Expr) ([]byte, Expr) {
	// Handle ByteArray
	if byteArray, ok := values.(ByteArray); ok {
		return byteArray.Data(), nil
	}

	// Handle single integer (representing a byte)
	if byteValue, ok := values.(Integer); ok {
		if byteValue < 0 || byteValue > 255 {
			return nil, NewErrorExpr("ValueError",
				fmt.Sprintf("Byte value %d is out of range (0-255)", byteValue),
				[]Expr{b, values})
		}
		return []byte{byte(byteValue)}, nil
	}

	// Handle List of integers
	if list, ok := values.(List); ok && len(list.Elements) > 1 {
		// Extract bytes from List (excluding head)
		elements := list.Elements[1:]
		bytes := make([]byte, len(elements))

		for i, elem := range elements {
			if byteValue, ok := elem.(Integer); ok {
				if byteValue < 0 || byteValue > 255 {
					return nil, NewErrorExpr("ValueError",
						fmt.Sprintf("Byte value %d at position %d is out of range (0-255)", byteValue, i+1),
						[]Expr{b, values})
				}
				bytes[i] = byte(byteValue)
			} else {
				return nil, NewErrorExpr("TypeError",
					fmt.Sprintf("List element at position %d must be an integer (0-255)", i+1),
					[]Expr{b, values})
			}
		}
		return bytes, nil
	}

	return nil, NewErrorExpr("TypeError",
		"ByteArray slice assignment requires ByteArray, integer, or List of integers",
		[]Expr{b, values})
}
