package core

// Expr is the fundamental interface for all expressions in the system
type Expr interface {
	String() string
	InputForm() string
	Head() string
	Length() int64
	Equal(rhs Expr) bool
	IsAtom() bool // Distinguishes atomic vs composite types
}

type Integer int64
type Real float64
type Symbol string
