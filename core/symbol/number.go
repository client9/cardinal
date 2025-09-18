package symbol

type NumberExpr interface {
	Expr

	Sign() int
	AsNeg() Expr
	AsInv() Expr
	Float64() float64
}
