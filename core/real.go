package core

import (
	"strconv"

	"github.com/client9/cardinal/core/big"
)

type Real interface {
	Expr
	Number

	IsInt() bool
	Prec() uint
	IsFloat64() bool

	AsBigFloat() *big.Float
}

func ParseReal(s string) (Expr, error) {
	if len(s) < 17 {
		if num, err := strconv.ParseFloat(s, 64); err == nil {
			return f64(num), nil
		}
	}
	return new(big.Float).SetString(s)
}

// Convert any number-type to a big.Float
// Done as a function to prevent import cycles, etc,
// and to prevent Float from having to know everything
func ToBigFloat(z *big.Float, x Number) *big.Float {
	switch n := x.(type) {
	case *big.Float:
		// TODO adjust z's precision based on input
		z.Set(n)
	case *big.Int:
		z.SetInt(n)
	case *big.Rat:
		z.SetRat(n)
	case f64:
		// TODO adjust z to machine precision??
		z.SetPrec(53)
		z.SetFloat64(n.Float64())
	case machineInt:
		z.SetInt64(n.Int64())
	case rat64:
		// TODO, could directly create big float with numerator, then divide
		z.SetRat(n.AsBigRat())
	default:
		panic("ToBigFloat: unknown number type")
	}
	return z
}
