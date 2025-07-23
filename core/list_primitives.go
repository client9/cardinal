package core

// ElementAt returns the element at the given index (1-based, like Mathematica)
// Returns the element and true if successful, nil and false if index is out of bounds
func (l List) ElementAt(index int64) (Expr, bool) {
	if len(l.Elements) == 0 {
		return nil, false
	}
	
	// Convert to 0-based indexing and handle negative indices
	var idx int
	if index > 0 {
		idx = int(index) // 1-based to 0-based: elements[0] is head, elements[1] is first element
	} else if index < 0 {
		idx = len(l.Elements) + int(index) // Negative indexing from end
	} else {
		return nil, false // index 0 is invalid in 1-based system
	}
	
	// Check bounds (elements[0] is head, so valid indices are 1 to len-1)
	if idx < 1 || idx >= len(l.Elements) {
		return nil, false
	}
	
	return l.Elements[idx], true
}

// SliceStart returns a new list with the first n elements (excluding head)
func (l List) SliceStart(n int64) List {
	if len(l.Elements) == 0 || n <= 0 {
		return List{Elements: []Expr{l.Elements[0]}} // Return list with just head
	}
	
	endIdx := int(n) + 1 // +1 because elements[0] is head
	if endIdx > len(l.Elements) {
		endIdx = len(l.Elements)
	}
	
	newElements := make([]Expr, endIdx)
	copy(newElements, l.Elements[:endIdx])
	return List{Elements: newElements}
}

// SliceEnd returns a new list starting from index n to the end (1-based indexing)
func (l List) SliceEnd(n int64) List {
	if len(l.Elements) == 0 {
		return List{Elements: []Expr{l.Elements[0]}} // Return list with just head
	}
	
	startIdx := int(n) // n is 1-based, so element n is at index n in 0-based array (since elements[0] is head)
	if startIdx < 1 {
		startIdx = 1
	}
	if startIdx >= len(l.Elements) {
		return List{Elements: []Expr{l.Elements[0]}} // Return list with just head
	}
	
	newElements := make([]Expr, len(l.Elements)-startIdx+1)
	newElements[0] = l.Elements[0] // Copy head
	copy(newElements[1:], l.Elements[startIdx:])
	return List{Elements: newElements}
}

// SliceBetween returns a new list containing elements from start to end (inclusive, 1-based)
func (l List) SliceBetween(start, end int64) List {
	if len(l.Elements) == 0 {
		return List{Elements: []Expr{l.Elements[0]}} // Return list with just head
	}
	
	startIdx := int(start) // 1-based to 0-based with head offset
	endIdx := int(end) + 1   // +1 for inclusive end and head offset
	
	if startIdx < 1 {
		startIdx = 1
	}
	if endIdx > len(l.Elements) {
		endIdx = len(l.Elements)
	}
	if startIdx >= endIdx || startIdx >= len(l.Elements) {
		return List{Elements: []Expr{l.Elements[0]}} // Return list with just head
	}
	
	newElements := make([]Expr, endIdx-startIdx+1)
	newElements[0] = l.Elements[0] // Copy head
	copy(newElements[1:], l.Elements[startIdx:endIdx])
	return List{Elements: newElements}
}

// SliceExclude returns a new list with elements from start to end removed (1-based indexing)
func (l List) SliceExclude(start, end int64) List {
	if len(l.Elements) == 0 {
		return List{Elements: []Expr{l.Elements[0]}} // Return list with just head
	}
	
	startIdx := int(start) // 1-based to 0-based with head offset
	endIdx := int(end) + 1   // +1 for inclusive end and head offset
	
	if startIdx < 1 {
		startIdx = 1
	}
	if endIdx > len(l.Elements) {
		endIdx = len(l.Elements)
	}
	if startIdx >= len(l.Elements) || startIdx >= endIdx {
		// Nothing to exclude, return copy of original
		newElements := make([]Expr, len(l.Elements))
		copy(newElements, l.Elements)
		return List{Elements: newElements}
	}
	
	// Calculate new length: original - excluded range
	newLen := len(l.Elements) - (endIdx - startIdx)
	newElements := make([]Expr, newLen)
	
	// Copy head
	newElements[0] = l.Elements[0]
	
	// Copy elements before excluded range
	copy(newElements[1:], l.Elements[1:startIdx])
	
	// Copy elements after excluded range
	if endIdx < len(l.Elements) {
		copy(newElements[startIdx:], l.Elements[endIdx:])
	}
	
	return List{Elements: newElements}
}