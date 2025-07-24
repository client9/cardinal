package core

import (
	"fmt"
	"strings"
)

// List represents compound expressions
type List struct {
	Elements []Expr
}

func (l List) Length() int64 {
	// really should panic
	if len(l.Elements) == 0 {
		return 0
	}
	// element[0] is the head
	return int64(len(l.Elements)) - 1
}

func (l List) String() string {
	if len(l.Elements) == 0 {
		return "List()"
	}

	// Check if this is a List literal (head is "List")
	if len(l.Elements) > 0 {
		if headAtom, ok := l.Elements[0].(Atom); ok &&
			headAtom.AtomType == SymbolAtom && headAtom.Value.(string) == "List" {
			// This is a list literal: [element1, element2, ...]
			var elements []string
			for _, elem := range l.Elements[1:] {
				elements = append(elements, elem.String())
			}
			return fmt.Sprintf("List(%s)", strings.Join(elements, ", "))
		}
	}

	// This is a function call: head(arg1, arg2, ...)
	var elements []string
	for _, elem := range l.Elements {
		elements = append(elements, elem.String())
	}
	return fmt.Sprintf("%s(%s)", l.Elements[0].String(), strings.Join(elements[1:], ", "))
}

func (l List) InputForm() string {
	return l.inputFormWithPrecedence(PrecedenceLowest)
}

func (l List) Type() string {
	return "list"
}

func (l List) Equal(rhs Expr) bool {
	rhsList, ok := rhs.(List)
	if !ok {
		return false
	}

	// Lists must have same number of elements
	if len(l.Elements) != len(rhsList.Elements) {
		return false
	}

	// Recursively compare each element
	for i, elem := range l.Elements {
		if !elem.Equal(rhsList.Elements[i]) {
			return false
		}
	}

	return true
}

// Sliceable interface implementation

// ElementAt returns the nth element (1-indexed, excludes head)
// For a list [head, e1, e2, e3], ElementAt(1) returns e1
func (l List) ElementAt(n int64) Expr {
	if len(l.Elements) <= 1 {
		return NewErrorExpr("PartError", "List has no elements", []Expr{l})
	}
	
	length := l.Length() // Number of elements excluding head
	
	// Handle negative indexing
	if n < 0 {
		n = length + n + 1
	}
	
	// Check bounds (1-indexed)
	if n <= 0 || n > length {
		return NewErrorExpr("PartError", 
			fmt.Sprintf("Part index %d is out of bounds for list with %d elements", n, length), 
			[]Expr{l})
	}
	
	// Convert to 0-based index (adding 1 because Elements[0] is head)
	return l.Elements[n]
}

// Slice returns a new list containing elements from start to stop (inclusive, 1-indexed)
// For a list [head, e1, e2, e3, e4], Slice(2, 3) returns [head, e2, e3]
func (l List) Slice(start, stop int64) Expr {
	if len(l.Elements) <= 1 {
		// Empty list - return list with just head
		return List{Elements: []Expr{l.Elements[0]}}
	}
	
	length := l.Length()
	
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
			fmt.Sprintf("Slice indices [%d, %d] out of bounds for list with %d elements", 
				start, stop, length), []Expr{l})
	}
	
	if start > stop {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Start index %d is greater than stop index %d", start, stop),
			[]Expr{l})
	}
	
	// Create new list with head + sliced elements
	// Convert to 0-based indices (Elements[0] is head, Elements[1] is first element)
	startIdx := start      // Elements[start] is the start element
	stopIdx := stop + 1    // Elements[stop+1] is exclusive end for Go slice
	
	newElements := make([]Expr, stopIdx-startIdx+1)
	newElements[0] = l.Elements[0] // Copy head
	copy(newElements[1:], l.Elements[startIdx:stopIdx])
	
	return List{Elements: newElements}
}