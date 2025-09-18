package core

import (
	"strconv"

	"github.com/client9/cardinal/core/big"
)

type Integer interface {
	Expr
	Number

	AsBigInt() *big.Int

	//
	// These are satified by big.Int
	//
	IsInt64() bool
	Int64() int64
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
	return z, true
}
