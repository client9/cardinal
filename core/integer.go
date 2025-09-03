package core

import (
	"math/big"
	"strconv"
)

type Integer interface {
	Expr

	IsInt64() bool
	Int64() int64
	asBigInt() bigInt

	Sign() int

	// Inv returns the reciprocal or Power(x,-1) of the value
	// // TODO zero
	Inv() Expr

	// TBD if actually needed
	Neg() Integer

	// TBD
	Float64() float64
}

func NewIntegerFromString(s string) (Integer, bool) {
	if len(s) < 19 {
		value, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return newMachineInt(0), false
		}
		return newMachineInt(value), true
	}
	z, ok := new(big.Int).SetString(s, 0)
	if !ok {
		return newMachineInt(0), false
	}
	return bigInt{val: z}, true
}
