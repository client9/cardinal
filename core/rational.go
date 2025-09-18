package core

import (
	"fmt"
	"github.com/client9/cardinal/core/big"
)

type Rational interface {
	Expr
	Number

	// returns true if denominator is 1
	IsInt() bool
	Sign() int

	AsDenom() Expr
	AsNum() Expr
	AsInv() Expr
	AsBigRat() *big.Rat
}

// this normalizes and may return an integer.
// prob wrong.
func NewRational(a, b int64) Expr {
	if b == 0 {
		panic("NewMachineRational got 0 denominator")
	}
	r := rat64{a, b}
	return r.StandardForm()
}

// The following functions return partially normalized Rationals.
// They do not convert to integers.
func addRat64(xi rat64, yi rat64) (rat64, bool) {

	newd, ok := lcm(xi.Denom().Int64(), yi.Denom().Int64())
	if !ok {
		return rat64Zero, false
	}

	xn, ok := timesInt64(xi.Num().Int64(), newd/xi.Denom().Int64())
	if !ok {
		return rat64Zero, false
	}
	yn, ok := timesInt64(yi.Num().Int64(), newd/yi.Denom().Int64())
	if !ok {
		return rat64Zero, false
	}
	numerator, ok := addInt64(xn, yn)
	if !ok {
		return rat64Zero, false
	}
	g := gcd(numerator, newd)

	return rat64{numerator / g, newd / g}, true
}

func timesRat64(xi, yi rat64) (rat64, bool) {
	fmt.Println("timesRat64", xi, yi)
	num, ok := timesInt64(xi.Num().Int64(), yi.Num().Int64())
	if !ok {
		return rat64Zero, false
	}

	den, ok := timesInt64(xi.Denom().Int64(), yi.Denom().Int64())
	if !ok {
		return rat64Zero, false
	}
	g := gcd(num, den)
	r := rat64{num / g, den / g}
	fmt.Println("timesRat64 Output:", num, den, g, r)
	return r, true
}
func timesRat64Int64(xi rat64, yi machineInt) (rat64, bool) {
	fmt.Println("Input", "rat=", xi, "int=", yi)
	num, ok := timesInt64(xi.Num().Int64(), yi.Int64())
	if !ok {
		return rat64Zero, false
	}
	den := xi.Denom().Int64()

	g := gcd(num, den)
	r := rat64{num / g, den / g}
	fmt.Println("Output", r)
	return r, true
}

// (a/b) * n = (a/b) * (n * b / b)
func addRat64Int64(xi rat64, yi machineInt) (rat64, bool) {
	xNum, ok := timesInt64(yi.Int64(), xi.Denom().Int64())
	if !ok {
		return rat64Zero, false
	}

	xNum, ok = addInt64(xNum, xi.Num().Int64())
	if !ok {
		return rat64Zero, false
	}
	g := gcd(xNum, xi.Denom().Int64())
	return rat64{xNum / g, xi.Denom().Int64() / g}, true
}

func addBigRatInt64(xi *big.Rat, yi machineInt) Rational {
	return addBigRatBigInt(xi, yi.AsBigInt())
}

func addRat64BigInt(xi rat64, yi *big.Int) Rational {
	return addBigRatBigInt(xi.AsBigRat(), yi)
}

func addBigRatBigInt(yi *big.Rat, xi *big.Int) Rational {
	numerator := new(big.Int).Mul(xi, yi.Denom())
	numerator.Add(numerator, yi.Num())

	denominator := new(big.Int).Set(yi.Denom())

	bigGCD := new(big.Int).GCD(nil, nil, numerator, denominator)
	denominator.Div(denominator, bigGCD)
	numerator.Div(numerator, bigGCD)
	return new(big.Rat).SetFrac(numerator, denominator)
}
