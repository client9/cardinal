package core

import (
	"fmt"
	"strings"
)

// List represents compound expressions
type List struct {
	Elements *[]Expr
}

func NewList(head string, args ...Expr) List {
	elements := make([]Expr, len(args)+1)
	elements[0] = NewSymbol(head)
	copy(elements[1:], args)
	return List{Elements: &elements}
}

// NewListFromExprs creates a List directly from expressions (for special cases)
// Use NewList instead when possible, as it enforces the Symbol-head convention
func NewListFromExprs(elements ...Expr) List {
	return List{Elements: &elements}
}

// Copy does a shallow clone of the List
// TBD if this should return List or Expr
func (l List) Copy() List {
	newElements := make([]Expr, len(*l.Elements))
	copy(newElements, *(l.Elements))
	return List{Elements: &newElements}
}

func (l List) Length() int64 {
	// really should panic
	if len(*l.Elements) == 0 {
		return 0
	}
	// element[0] is the head
	return int64(len(*l.Elements)) - 1
}

func (l List) String() string {
	if len(*l.Elements) == 0 {
		return "List()"
	}

	// Check if this is a List literal (head is "List")
	if len(*l.Elements) > 0 {
		isListLiteral := false

		// Check new Symbol type first
		if l.Head() == "List" {
			isListLiteral = true
		}

		if isListLiteral {
			// This is a list literal: [element1, element2, ...]
			var elements []string
			for _, elem := range l.Tail() {
				elements = append(elements, elem.String())
			}
			return fmt.Sprintf("List(%s)", strings.Join(elements, ", "))
		}
	}

	// This is a function call: head(arg1, arg2, ...)
	var elements []string
	for _, elem := range l.AsSlice() {
		elements = append(elements, elem.String())
	}
	return fmt.Sprintf("%s(%s)", l.Head(), strings.Join(elements[1:], ", "))
}
func (l List) InputForm() string {
	return l.inputFormWithPrecedence(PrecedenceLowest)
}

func (l List) HeadExpr() Expr {
	return (*l.Elements)[0]
}

func (l List) Head() string {
	if name, ok := ExtractSymbol(l.HeadExpr()); ok {
		return name
	}
	panic("Head of List is not a symbol")
}

func (l List) Tail() []Expr {
	return (*l.Elements)[1:]
}

func (l List) AsSlice() []Expr {
	return *(l.Elements)
}

func (l List) Equal(rhs Expr) bool {
	rhsList, ok := rhs.(List)
	if !ok {
		return false
	}

	lhsSlice := l.AsSlice()
	rhsSlice := rhsList.AsSlice()

	if len(lhsSlice) != len(rhsSlice) {
		return false
	}

	// Recursively compare each element
	for i, elem := range lhsSlice {
		if !elem.Equal(rhsSlice[i]) {
			return false
		}
	}

	return true
}

func (l List) IsAtom() bool {
	return false
}

// Sliceable interface implementation

// ElementAt returns the nth element (1-indexed, excludes head)
// For a list [head, e1, e2, e3], ElementAt(1) returns e1
func (l List) ElementAt(n int64) Expr {
	e, err := ElementAt(l.Tail(), int(n))
	if err != nil {
		return NewError(err.Error(), "")
	}
	return e
}

// Slice returns a new list containing elements from start to stop (inclusive, 1-indexed)
// For a list [head, e1, e2, e3, e4], Slice(2, 3) returns [head, e2, e3]
func (l List) Slice(start, stop int64) Expr {
	e, err := Slice(l.Tail(), int(start), int(stop))
	if err != nil {
		return NewError(err.Error(), "")
	}
	newElements := make([]Expr, len(e)+1)
	newElements[0] = l.HeadExpr()
	copy(newElements[1:], e)
	return List{Elements: &newElements}
}

// Join joins this list with another sliceable of the same type
// Both lists must have the same head to be joined
func (l List) Join(other Sliceable) Expr {
	// Type check: ensure other is also a List
	otherList, ok := other.(List)
	if !ok {
		return NewError("TypeError",
			fmt.Sprintf("Cannot join %T with List", other))
	}

	// Handle empty lists
	if l.Length() == 0 {
		return otherList // Return the other list if this one is empty
	}
	if otherList.Length() == 0 {
		return l // Return this list if the other one is empty
	}

	if l.Head() != otherList.Head() {
		return NewError("TypeError",
			fmt.Sprintf("Cannot join lists with different heads: %s and %s",
				l.Head(), otherList.Head()))
	}

	// Create new list with combined elements
	// newElements = [head, l.elements[1:], otherList.elements[1:]]
	newElements := make([]Expr, 1+l.Length()+otherList.Length())
	newElements[0] = l.HeadExpr()

	// Copy elements from first list (excluding head)
	copy(newElements[1:], l.Tail())

	// Copy elements from second list (excluding head)
	copy(newElements[1+l.Length():], otherList.Tail())

	return List{Elements: &newElements}
}

// Appends an expression to the end of a List
func (l List) Append(e Expr) List {
	dest := make([]Expr, l.Length()+2)
	copy(dest, l.AsSlice())
	dest[len(dest)-1] = e
	return List{Elements: &dest}
}

// SetElementAt returns a new List with the nth element replaced (1-indexed)
// Returns an error Expr if index is out of bounds
func (l List) SetElementAt(n int64, value Expr) Expr {
	length := l.Length() // Number of elements excluding head

	if length == 0 {
		return NewError("PartError", "List has no elements")
	}

	// Handle negative indexing
	if n < 0 {
		n = length + n + 1
	}

	// Check bounds (1-indexed)
	if n <= 0 || n > length {
		return NewError("PartError",
			fmt.Sprintf("Part index %d is out of bounds for list with %d elements", n, length))
	}

	// Create new list with element replaced
	newElements := make([]Expr, len(*l.Elements))
	copy(newElements, *(l.Elements))
	newElements[n] = value // n is 1-indexed, but array is 0-indexed after head

	l.Elements = &newElements
	//return value
	return List{Elements: &newElements}
}

// SetSlice returns a new List with elements from start to stop replaced by values (1-indexed)
// values can be a single Expr, List, or other Sliceable
func (l List) SetSlice(start, stop int64, values Expr) Expr {
	length := l.Length()
	if length == 0 {
		// Empty list - can only insert at position 1
		if start == 1 && stop == 0 {
			return l.insertValues(1, values)
		}
		return NewError("PartError", "List has no elements")
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
		return NewError("PartError",
			fmt.Sprintf("Start index %d must be positive", start))
	}

	if start > length+1 {
		return NewError("PartError",
			fmt.Sprintf("Start index %d is out of bounds for list with %d elements", start, length))
	}

	// Handle special cases
	if stop < start-1 {
		return NewError("PartError",
			fmt.Sprintf("Stop index %d cannot be less than start index %d - 1", stop, start))
	}

	// Convert values to slice
	var valueSlice []Expr
	if valuesList, ok := values.(List); ok && valuesList.Length() > 0 {
		// Extract elements from List (excluding head)
		valueSlice = valuesList.Tail()
	} else if sliceable := AsSliceable(values); sliceable != nil {
		// Handle other sliceable types by converting to list
		if valuesList, ok := values.(List); ok {
			valueSlice = valuesList.Tail()
		} else {
			// For non-List sliceables, treat as single element
			valueSlice = []Expr{values}
		}
	} else {
		// Single value
		valueSlice = []Expr{values}
	}

	// Calculate new list size
	oldRangeSize := int64(0)
	if stop >= start {
		oldRangeSize = stop - start + 1
	}
	newSize := 1 + int64(length) - oldRangeSize + int64(len(valueSlice))

	// Create new list
	newElements := make([]Expr, newSize)

	// Copy head
	newElements[0] = l.HeadExpr()

	// Copy elements before the range
	if start > 1 {
		copy(newElements[1:start], (*l.Elements)[1:start])
	}

	// Insert new values
	if len(valueSlice) > 0 {
		copy(newElements[start:start+int64(len(valueSlice))], valueSlice)
	}

	// Copy elements after the range
	if stop < length {
		afterStart := start + int64(len(valueSlice))
		copy(newElements[afterStart:], (*l.Elements)[stop+1:])
	}

	return List{Elements: &newElements}
}

// insertValues is a helper method for inserting values at a specific position
func (l List) insertValues(pos int64, values Expr) Expr {
	// Convert values to slice
	var valueSlice []Expr
	if valuesList, ok := values.(List); ok && valuesList.Length() > 0 {
		valueSlice = valuesList.Tail()
	} else {
		valueSlice = []Expr{values}
	}

	// Create new list with values inserted
	newSize := 1 + int(l.Length()) + len(valueSlice)
	newElements := make([]Expr, newSize)

	// Copy head
	newElements[0] = l.HeadExpr()

	// Copy elements before insertion point
	if pos > 1 {
		copy(newElements[1:pos], l.Tail()[:pos])
	}

	// Insert new values
	copy(newElements[pos:pos+int64(len(valueSlice))], valueSlice)

	// Copy remaining elements
	if pos <= int64(l.Length())+1 {
		copy(newElements[pos+int64(len(valueSlice)):], l.Tail()[pos:])
	}

	return List{Elements: &newElements}
}
