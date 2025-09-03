package core

import (
	"math"
	"math/big"
)

func addInt64(x, y int64) (int64, bool) {
	if y > 0 {
		if x > math.MaxInt64-y {
			return 0, false
		}
	} else {
		if x < math.MinInt64-y {
			return 0, false
		}
	}
	return x + y, true
}

func timesInt64(x, y int64) (int64, bool) {
	if x == 0 || y == 0 {
		return 0, true
	}
	z := x * y

	// if z < 0 (true),   and x or y is negative
	// if z > 0 (false),  and both positive
	//
	// or another way is if you
	// z is positive, but x, y aren't positve, then overflow
	// z is negative, but x, y aren't mixed, then overflow
	if (z < 0) == ((x < 0) != (y < 0)) {
		if z/y == x {
			return z, true
		}
	}
	return 0, false
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) (int64, bool) {
	if a == 0 || b == 0 {
		return 0, true // LCM is 0 if either number is 0
	}
	// Ensure positive values for calculation
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if a == b {
		return a, true
	}

	prod, ok := timesInt64(a, b)
	if !ok {
		return 0, false
	}

	return prod / gcd(a, b), true
}

/*
func addIntRat(xi Integer, yi Rational) Rational {
	if mint, ok := xi.(machineInt); ok {
		if mrat, ok := yi.(rat64); ok {
			return addRat64Int64(mrat, mint)
		}
		return addBigRatInt64(yi.(bigRat), mint)
	}
	bint := xi.(bigInt)
	if mrat, ok := yi.(rat64); ok {
		return addRat64BigInt(mrat, bint)
	}
	return addBigRatBigInt(yi.(bigRat), bint)
}
*/

func addInteger(xi, yi Integer) Integer {
	// both are not bigints
	if xi.IsInt64() && yi.IsInt64() {
		if val, ok := addInt64(xi.Int64(), yi.Int64()); ok {
			return newMachineInt(val)
		}
	}
	// one is BigInt, or the addition would overflow.
	x := xi.asBigInt()
	y := yi.asBigInt()

	z := NewBigInt(0)
	z.add(x)
	z.add(y)
	return z
}

func timesInteger(xi, yi Integer) Integer {
	// both are not bigints
	if xi.IsInt64() && yi.IsInt64() {
		if val, ok := timesInt64(xi.Int64(), yi.Int64()); ok {
			return newMachineInt(val)
		}
	}
	// one is BigInt, or the addition would overflow.
	x := xi.asBigInt()
	y := yi.asBigInt()
	z := NewBigInt(1)
	z.times(x)
	z.times(y)
	return z
}

func PowerInteger(xi, yi Integer) Integer {
	if !yi.IsInt64() {
		panic("overflow")
	}
	y := yi.Int64()
	if xi.IsInt64() {
		return powerSmall(xi.Int64(), y)
	}
	return powerBig(xi.(bigInt), y)

}

func powerBig(xi bigInt, y int64) Integer {
	// make copy
	x := new(big.Int).Set(xi.val)
	r := big.NewInt(1)
	return powerBigLoop(x, y, r)
}

func powerBigLoop(x *big.Int, y int64, r *big.Int) Integer {
	for y > 0 {
		// if y is odd
		if y&1 == 1 {
			r.Mul(r, x)
		}
		x.Mul(x, x)
		y /= 2
	}
	return newBigInt(r)
}

func powerSmall(x, y int64) Integer {
	r := int64(1)
	for y > 0 {
		// if y is odd
		if y&1 == 1 {
			newr, ok1 := timesInt64(r, x)
			newx, ok2 := timesInt64(x, x)
			if !(ok1 && ok2) {
				return powerBigLoop(big.NewInt(x), y, big.NewInt(r))
			}
			r = newr
			x = newx
		} else {
			newx, ok := timesInt64(x, x)
			if !ok {
				return powerBigLoop(big.NewInt(x), y, big.NewInt(r))
			}
			x = newx
		}

		y /= 2
	}
	return newMachineInt(r)
}
