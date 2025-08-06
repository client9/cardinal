package core

import (
	"strconv"
	"strings"
)

// Integer type implementation
func (i Integer) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Integer) InputForm() string {
	return i.String()
}

func (i Integer) Head() string {
	return "Integer"
}

func (i Integer) Length() int64 {
	return 0
}

func (i Integer) Equal(rhs Expr) bool {
	if other, ok := rhs.(Integer); ok {
		return i == other
	}
	return false
}

func (i Integer) IsAtom() bool {
	return true
}

// Real type implementation
func (r Real) String() string {
	str := strconv.FormatFloat(float64(r), 'f', -1, 64)
	if !strings.Contains(str, ".") {
		str += ".0"
	}
	return str
}

func (r Real) InputForm() string {
	return r.String()
}

func (r Real) Head() string {
	return "Real"
}

func (r Real) Length() int64 {
	return 0
}

func (r Real) Equal(rhs Expr) bool {
	if other, ok := rhs.(Real); ok {
		return r == other
	}
	return false
}

func (r Real) IsAtom() bool {
	return true
}

// Symbol type implementation
func (s Symbol) String() string {
	return string(s)
}

func (s Symbol) InputForm() string {
	return s.String()
}

func (s Symbol) Head() string {
	return "Symbol"
}

func (s Symbol) Length() int64 {
	return 0
}

func (s Symbol) Equal(rhs Expr) bool {
	if other, ok := rhs.(Symbol); ok {
		return s == other
	}
	return false
}

func (s Symbol) IsAtom() bool {
	return true
}
