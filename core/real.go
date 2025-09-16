package core

import (
	"fmt"
	"strconv"
)

type Real interface {
	Expr

	Neg() Real
	Sign() int
	Float64() float64
	IsInt() bool
	Prec() uint
	IsFloat64() bool
	AsBigFloat() BigFloat
}

func ParseReal(s string) (Expr, error) {
	if len(s) < 17 {
		if num, err := strconv.ParseFloat(s, 64); err == nil {
			return f64(num), nil
		}
	}
	z, ok := new(BigFloat).SetString(s)
	if !ok {
		return nil, fmt.Errorf("invalid floating point string")
	}
	return z, nil
}
