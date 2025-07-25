package core

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// AtomType represents the type of an atomic value
type AtomType int

const (
	StringAtom AtomType = iota
	IntAtom
	FloatAtom
	SymbolAtom
)

// Atom represents atomic values (strings, integers, floats, symbols)
type Atom struct {
	AtomType AtomType
	Value    interface{}
}

func (a Atom) Length() int64 {
	if a.AtomType == StringAtom {
		return int64(utf8.RuneCountInString(a.Value.(string)))
	}
	return 0
}

func (a Atom) String() string {
	switch a.AtomType {
	case StringAtom:
		return fmt.Sprintf("\"%s\"", a.Value.(string))
	case IntAtom:
		return strconv.Itoa(a.Value.(int))
	case FloatAtom:
		// Always show decimal point to distinguish from integers
		str := strconv.FormatFloat(a.Value.(float64), 'f', -1, 64)
		if !strings.Contains(str, ".") {
			str += ".0"
		}
		return str
	case SymbolAtom:
		return a.Value.(string)
	default:
		return ""
	}
}

func (a Atom) InputForm() string {
	// For atoms, InputForm is the same as String()
	return a.String()
}

func (a Atom) Type() string {
	switch a.AtomType {
	case StringAtom:
		return "string"
	case IntAtom:
		return "int"
	case FloatAtom:
		return "float64"
	case SymbolAtom:
		return "symbol"
	default:
		return "unknown"
	}
}

func (a Atom) Equal(rhs Expr) bool {
	rhsAtom, ok := rhs.(Atom)
	if !ok {
		return false
	}

	// Must have same atom type
	if a.AtomType != rhsAtom.AtomType {
		return false
	}

	// Compare values based on type
	switch a.AtomType {
	case StringAtom:
		return a.Value.(string) == rhsAtom.Value.(string)
	case IntAtom:
		return a.Value.(int) == rhsAtom.Value.(int)
	case FloatAtom:
		return a.Value.(float64) == rhsAtom.Value.(float64)
	case SymbolAtom:
		return a.Value.(string) == rhsAtom.Value.(string)
	default:
		return false
	}
}

// Sliceable interface implementation (only for StringAtom)

// ElementAt returns the nth character (1-indexed) for string atoms
func (a Atom) ElementAt(n int64) Expr {
	if a.AtomType != StringAtom {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("ElementAt not supported for %s", a.Type()), []Expr{a})
	}

	str := a.Value.(string)
	runes := []rune(str)
	length := int64(len(runes))

	if length == 0 {
		return NewErrorExpr("PartError", "String is empty", []Expr{a})
	}

	// Handle negative indexing
	if n < 0 {
		n = length + n + 1
	}

	// Check bounds (1-indexed)
	if n <= 0 || n > length {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Part index %d is out of bounds for string with %d characters", n, length),
			[]Expr{a})
	}

	// Convert to 0-based index and return character as string
	return NewStringAtom(string(runes[n-1]))
}

// Slice returns a substring from start to stop (inclusive, 1-indexed) for string atoms
func (a Atom) Slice(start, stop int64) Expr {
	if a.AtomType != StringAtom {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Slice not supported for %s", a.Type()), []Expr{a})
	}

	str := a.Value.(string)
	runes := []rune(str)
	length := int64(len(runes))

	if length == 0 {
		return NewStringAtom("")
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
			fmt.Sprintf("Slice indices [%d, %d] out of bounds for string with %d characters",
				start, stop, length), []Expr{a})
	}

	if start > stop {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Start index %d is greater than stop index %d", start, stop),
			[]Expr{a})
	}

	// Convert to 0-based indices and return substring
	return NewStringAtom(string(runes[start-1 : stop]))
}

// Join joins this string atom with another sliceable of the same type
func (a Atom) Join(other Sliceable) Expr {
	// Only string atoms support joining
	if a.AtomType != StringAtom {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Atom of type %d does not support join", a.AtomType),
			[]Expr{a})
	}

	// Type check: ensure other is also a string atom
	otherAtom, ok := other.(Atom)
	if !ok {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Cannot join %T with string", other),
			[]Expr{a, other.(Expr)})
	}

	if otherAtom.AtomType != StringAtom {
		return NewErrorExpr("TypeError",
			fmt.Sprintf("Cannot join atom %d with string", otherAtom.AtomType),
			[]Expr{a, otherAtom})
	}

	// Join the string values
	thisStr := a.Value.(string)
	otherStr := otherAtom.Value.(string)

	return NewStringAtom(thisStr + otherStr)
}
