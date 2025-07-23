package core

import (
	"fmt"
	"strconv"
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
		return int64(len(a.Value.(string)))
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
		return strconv.FormatFloat(a.Value.(float64), 'f', -1, 64)
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