package core

import (
	"math/big"
	"strconv"

	"github.com/client9/cardinal/core/symbol"
)

// MachineInteger
type machineInt int64

func MustInt64(e Expr) int64 {
	return e.(Integer).Int64()
}

func NewInteger(n int64) Integer {
	return newMachineInt(n)
}

func newMachineInt(i int64) machineInt {
	return machineInt(i)
}

// Integer type implementation
func (i machineInt) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i machineInt) InputForm() string {
	return i.String()
}

func (i machineInt) Head() Expr {
	return symbol.Integer
}

func (i machineInt) Length() int64 {
	return 0
}

func (i machineInt) IsAtom() bool {
	return true
}

func (i machineInt) Equal(rhs Expr) bool {
	switch intval := rhs.(type) {
	case machineInt:
		return i == intval
	case BigInt:
		return intval.IsInt64() && intval.Int64() == i.Int64()
	default:
		return false
	}
}

// Integer Interface
func (i machineInt) Float64() float64 {
	return float64(i)
}
func (i machineInt) IsInt64() bool {
	return true
}
func (i machineInt) Int64() int64 {
	return int64(i)
}

func (i machineInt) Inv() Expr {
	// TODO ZERO
	return rat64{1, i.Int64()}
}

func (i machineInt) Sign() int {
	if i < 0 {
		return -1
	}
	if i > 0 {
		return 1
	}
	return 0
}

func (i machineInt) AsBigInt() BigInt {
	return BigInt{
		val: big.NewInt(int64(i)),
	}
}
func (i machineInt) asBigInt() BigInt {
	return newBigInt(big.NewInt(int64(i)))
}

// TBD ... unclear if needed or if it's just an optimization of "Times(-1, x)"
func (i machineInt) Neg() Integer {
	return machineInt(-i)
}
