package core

// Expr is the fundamental interface for all expressions in the system
type Expr interface {
	String() string
	InputForm() string
	Type() string
	Length() int64
	Equal(rhs Expr) bool
}