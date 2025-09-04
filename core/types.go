package core

// Expr is the fundamental interface for all expressions in the system
type Expr interface {
	String() string
	InputForm() string
	Head() Expr
	Equal(rhs Expr) bool

	// Length returns 0 is atomic, or the length of the list item.
	Length() int64

	// IsAtom returns true is not a compound element, and no futher reduction is possible.
	IsAtom() bool
}
