package core

type Number interface {
	Expr

	Sign() int
}
