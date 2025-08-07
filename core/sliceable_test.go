package core

import (
	"testing"
)

func TestSliceableInterface(t *testing.T) {
	// Test List implements Sliceable
	list := NewList("List", NewInteger(1), NewInteger(2), NewInteger(3))

	if !IsSliceable(list) {
		t.Error("List should be sliceable")
	}

	sliceable := AsSliceable(list)
	if sliceable == nil {
		t.Error("List should cast to Sliceable")
	}

	// Test ElementAt
	elem := sliceable.ElementAt(2)
	if !elem.Equal(NewInteger(2)) {
		t.Errorf("ElementAt(2) expected 2, got %v", elem)
	}

	// Test Slice
	slice := sliceable.Slice(1, 2)
	expectedSlice := NewList("List", NewInteger(1), NewInteger(2))
	if !slice.Equal(expectedSlice) {
		t.Errorf("Slice(1,2) expected %v, got %v", expectedSlice, slice)
	}

	// Test String Atom implements Sliceable
	str := NewString("hello")

	if !IsSliceable(str) {
		t.Error("String Atom should be sliceable")
	}

	strSliceable := AsSliceable(str)
	if strSliceable == nil {
		t.Error("String Atom should cast to Sliceable")
	}

	// Test ElementAt on string
	char := strSliceable.ElementAt(2)
	if !char.Equal(NewString("e")) {
		t.Errorf("String ElementAt(2) expected 'e', got %v", char)
	}

	// Test Slice on string
	substr := strSliceable.Slice(2, 4)
	if !substr.Equal(NewString("ell")) {
		t.Errorf("String Slice(2,4) expected 'ell', got %v", substr)
	}

	// Test ByteArray implements Sliceable
	ba := NewByteArray([]byte{65, 66, 67, 68}) // "ABCD"

	if !IsSliceable(ba) {
		t.Error("ByteArray should be sliceable")
	}

	baSliceable := AsSliceable(ba)
	if baSliceable == nil {
		t.Error("ByteArray should cast to Sliceable")
	}

	// Test ElementAt on ByteArray
	byte1 := baSliceable.ElementAt(1)
	if !byte1.Equal(NewInteger(65)) {
		t.Errorf("ByteArray ElementAt(1) expected 65, got %v", byte1)
	}

	// Test Slice on ByteArray
	baSlice := baSliceable.Slice(2, 3)
	expectedBA := NewByteArray([]byte{66, 67})
	if !baSlice.Equal(expectedBA) {
		t.Errorf("ByteArray Slice(2,3) expected %v, got %v", expectedBA, baSlice)
	}

	// Test non-sliceable type
	intAtom := NewInteger(42)
	if IsSliceable(intAtom) {
		t.Error("Int Atom should not be sliceable")
	}

	if AsSliceable(intAtom) != nil {
		t.Error("Int Atom should not cast to Sliceable")
	}
}
