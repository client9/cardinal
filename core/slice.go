package core

import (
	"fmt"
	"slices"
)

func ElementAt[T any](s []T, n int) (T, error) {
	var zero T
	length := len(s)
	// Handle negative indexing
	if n < 0 {
		n = length + n + 1
	}

	// Check bounds (1-indexed)
	if n <= 0 || n > length {
		return zero, fmt.Errorf("Bounds error")
	}
	return s[n-1], nil
}

func Slice[S ~[]E, E any](s S, start, stop int) (S, error) {
	length := len(s)
	if length == 0 {
		return s, nil
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
		return nil, fmt.Errorf("indexes out of bounds")
	}

	if start > stop {
		return nil, fmt.Errorf("start index is greater than stop index")
	}

	// Convert to 0-based indices and create new ByteArray
	startIdx := start - 1
	stopIdx := stop // stop is inclusive, so we include it

	return s[startIdx:stopIdx], nil
}

// Join make a new slice of the two inputs
func Join[S ~[]E, E any](a, b S) (S, error) {
	// Create new byte array with combined data
	newData := make(S, len(a)+len(b))
	copy(newData, a)
	copy(newData[len(a):], b)
	return newData, nil
}

func SetElementAt[S ~[]E, E any](s S, n int, val E) (S, error) {
	length := len(s)

	if length == 0 {
		return nil, fmt.Errorf("Part Error")
	}

	// Handle negative indexing
	if n < 0 {
		n = length + n + 1
	}

	// Check bounds (1-indexed)
	if n <= 0 || n > length {
		return nil, fmt.Errorf("Part Error")
	}

	// Create new slice with element replaced
	newData := make([]E, length)
	copy(newData, s)
	newData[n-1] = val // n is 1-indexed
	return newData, nil
}

func Replace[S ~[]E, E any](s S, i, j int, v S) (S, error) {

	length := len(s)
	// Handle negative indexing
	if i < 0 {
		i = length + i + 1
	}
	if j < 0 {
		j = length + j + 1
	}
	if i <= 0 || i > length+1 {
		return nil, fmt.Errorf("PartError")
	}
	if j < i-1 {
		return nil, fmt.Errorf("PartError")
	}

	out := make(S, length)
	copy(out, s)
	slices.Replace(s, i, j, v...)
	return out, nil
}
