package core

// Expr is the fundamental interface for all expressions in the system
type Expr interface {
	String() string
	InputForm() string
	Type() string
	Length() int64
	Equal(rhs Expr) bool
	IsAtom() bool // Distinguishes atomic vs composite types
}

// New atomic types replacing the old Atom union type
type String string
type Integer int64
type Real float64
type Symbol string
