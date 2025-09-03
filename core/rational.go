package core

import (
	"math/big"
)

type Rational interface {
	Expr

	Denominator() Integer
	Numerator() Integer

	// returns true if denominator is 1
	IsInt() bool

	Sign() int

	Float64() float64

	// TBD if needed
	Neg() Rational

	// Inv returns the reciprocal
	Inv() Rational

	asBigRat() bigRat
}

// this normalizes and may return an integer.
// prob wrong.
func NewRational(a, b int64) Expr {
	if b == 0 {
		panic("NewMachineRational got 0 denominator")
	}
	if a == 0 {
		return newMachineInt(0)
	}
	if a == b {
		return newMachineInt(1)
	}
	// if denom is negative, normalize
	if b < 0 {
		a = -a
		b = -b
	}

	if b == 1 {
		return newMachineInt(a)
	}
	if a == 1 || a == -1 {
		return rat64{a, b}
	}
	g := gcd(a, b)
	if g == 1 {
		return rat64{a, b}
	}

	return rat64{a / g, b / g}
}

// The following functions return partially normalized Rationals.
// They do not convert to integers.
func addRat64(xi rat64, yi rat64) (rat64, bool) {

	newd, ok := lcm(xi.Denominator().Int64(), yi.Denominator().Int64())
	if !ok {
		return rat64Zero, false
	}

	xn, ok := timesInt64(xi.Numerator().Int64(), newd/xi.Denominator().Int64())
	if !ok {
		return rat64Zero, false
	}
	yn, ok := timesInt64(yi.Numerator().Int64(), newd/yi.Denominator().Int64())
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
	num, ok := timesInt64(xi.Numerator().Int64(), yi.Numerator().Int64())
	if !ok {
		return rat64Zero, false
	}

	den, ok := timesInt64(xi.Denominator().Int64(), yi.Denominator().Int64())
	if !ok {
		return rat64Zero, false
	}
	g := gcd(num, den)
	return rat64{num / g, den / g}, true
}
func timesRat64Int64(xi rat64, yi machineInt) (rat64, bool) {
	num, ok := timesInt64(xi.Numerator().Int64(), yi.Int64())
	if !ok {
		return rat64Zero, false
	}
	den := xi.Denominator().Int64()
	g := gcd(num, den)
	return rat64{num / g, den / g}, true
}

func addRat64Int64(xi rat64, yi machineInt) (rat64, bool) {
	xNum, ok := timesInt64(yi.Int64(), xi.Denominator().Int64())
	if !ok {
		return rat64Zero, false
	}

	xNum, ok = addInt64(xNum, xi.Numerator().Int64())
	if !ok {
		return rat64Zero, false
	}
	g := gcd(xNum, xi.Denominator().Int64())
	return rat64{xNum / g, xi.Denominator().Int64() / g}, true
}

func addBigRatInt64(xi bigRat, yi machineInt) Rational {
	return addBigRatBigInt(xi, yi.asBigInt())
}

func addRat64BigInt(xi rat64, yi bigInt) Rational {
	return addBigRatBigInt(xi.asBigRat(), yi)
}

func addBigRatBigInt(yi bigRat, xi bigInt) Rational {
	numerator := new(big.Int).Mul(xi.val, yi.val.Denom())
	numerator.Add(numerator, yi.val.Num())

	denominator := new(big.Int).Set(yi.val.Denom())

	bigGCD := new(big.Int).GCD(nil, nil, numerator, denominator)
	denominator.Div(denominator, bigGCD)
	numerator.Div(numerator, bigGCD)
	return newBigRat(new(big.Rat).SetFrac(numerator, denominator))
}
