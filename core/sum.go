package core

import (
	"math/big"

	"github.com/client9/sexpr/core/symbol"
)

func PlusList(args []Expr) Expr {
	intsum := PlusInteger{}
	ratsum := PlusRational{
		sum: rat64Zero,
	}
	realsum := PlusReal{}
	var nonnum []Expr

	for _, arg := range args {
		switch num := arg.(type) {
		case machineInt:
			intsum.Update(num)
		case bigInt:
			intsum.UpdateBig(num)
		case rat64:
			ratsum.UpdateRat64(num)
		case bigRat:
			ratsum.UpdateBigRat(num)
		case Real:
			realsum.Update(num)
		default:
			nonnum = append(nonnum, num)
		}
	}

	resultElements := make([]Expr, 0, 2+len(nonnum))
	resultElements = append(resultElements, symbol.Plus)

	// ints get turned into rationals
	if intsum.exists() && ratsum.exists() {
		ratsum.UpdateInt64(intsum.sum)
		if intsum.bigsum.Sign() != 0 {
			ratsum.UpdateBigInt(intsum.bigsum)
		}
	} else if intsum.exists() {
		total := intsum.Total()
		if realsum.exists() {
			realsum.Update(Real(intsum.Total().Float64()))
		} else {
			if len(nonnum) == 0 || total.Sign() != 0 {
				resultElements = append(resultElements, total)
			}
		}
	}
	if ratsum.exists() {
		total := ratsum.Total()
		if realsum.exists() {
			// All ints and rationals are added
			// Convert to float, and let next part deal with it
			realsum.Update(Real(total.Float64()))
		} else {
			if len(nonnum) == 0 || total.Sign() != 0 {
				if total.IsInt() {
					resultElements = append(resultElements, total.Numerator())
				} else {
					resultElements = append(resultElements, total)
				}
			}
		}
	}

	if realsum.exists() {
		total := realsum.Total()
		if len(nonnum) == 0 || total != 0.0 {
			resultElements = append(resultElements, total)
		}
	}

	// Add non-numeric terms
	resultElements = append(resultElements, nonnum...)

	// Apply OneIdentity-like behavior: if only one element (plus head), return it
	if len(resultElements) == 2 {
		return resultElements[1]
	}

	// If no elements besides head, return 0
	if len(resultElements) == 1 {
		return NewInteger(0)
	}

	return NewListFromExprs(resultElements...)
}

// Accumulators -- add similar types
//
// Plus and Times are unqiue in that they can take a list of many items.
//
// Other operators such as Subtract,Divide,Power are binary operators.
//
// The strategy is to add each type separately, and sort out the combined mess at the end.

// Adds a series of integers
type PlusInteger struct {
	sum    machineInt
	bigsum bigInt

	count int
}

func (a *PlusInteger) exists() bool {
	return a.count > 0
}

func (a *PlusInteger) Update(b machineInt) {
	if sumnext, ok := addInt64(a.sum.Int64(), b.Int64()); ok {
		a.count += 1
		a.sum = newMachineInt(sumnext)
		return
	}
	a.UpdateBig(a.sum.asBigInt())
	a.UpdateBig(b.asBigInt())
	a.sum = 0
}

func (a *PlusInteger) UpdateBig(b bigInt) {
	a.count += 1
	if a.bigsum.val == nil {
		a.bigsum = bigIntZero()
	}
	a.bigsum.add(b)
}

func (a *PlusInteger) Total() Integer {
	if a.bigsum.val == nil {
		return a.sum
	}

	if a.sum != 0 {
		a.bigsum.add(a.sum.asBigInt())
	}
	return a.bigsum
}

// Adds a series of integers
type PlusRational struct {
	sum    rat64
	bigsum bigRat

	count int
}

func (a *PlusRational) exists() bool {
	return a.count > 0
}
func (a *PlusRational) UpdateRat64(b rat64) {
	if sumnext, ok := addRat64(a.sum, b); ok {
		a.sum = sumnext
		a.count += 1
		return
	}
	a.UpdateBigRat(a.sum.asBigRat())
	a.UpdateBigRat(b.asBigRat())
	a.sum = rat64Zero
}

func (a *PlusRational) UpdateInt64(b machineInt) {
	if nextsum, ok := addRat64Int64(a.sum, b); ok {
		a.sum = nextsum
		a.count += 1
		return
	}
	a.UpdateBigRat(a.sum.asBigRat())
	a.sum = rat64Zero
	a.UpdateBigInt(b.asBigInt())
}

func (a *PlusRational) UpdateBigInt(b bigInt) {
	if a.bigsum.val == nil {
		a.bigsum.val = big.NewRat(0, 1)
	}
	a.count += 1
}
func (a *PlusRational) UpdateBigRat(b bigRat) {
	if a.bigsum.val == nil {
		a.bigsum.val = big.NewRat(0, 1)
	}
	a.count += 1
	a.bigsum.add(b)
}

func (a *PlusRational) Total() Rational {
	if a.bigsum.val == nil {
		return a.sum
	}
	a.bigsum.add(a.sum.asBigRat())
	return a.bigsum
}

// Adds a series of floats
type PlusReal struct {
	sum   Real
	count int
}

func (a *PlusReal) exists() bool {
	return a.count > 0
}
func (a *PlusReal) Update(b Real) {
	a.sum += b
	a.count += 1
}

func (a *PlusReal) Total() Real {
	return Real(a.sum)
}

func TimesList(args []Expr) Expr {
	intprod := TimesInteger{
		prod: newMachineInt(1),
	}
	ratprod := TimesRational{
		prod: rat64One,
	}
	realprod := TimesReal{
		prod: 1.0,
	}
	var nonnum []Expr

	for _, arg := range args {
		switch num := arg.(type) {
		case machineInt:
			intprod.Update(num)
		case bigInt:
			intprod.UpdateBig(num)
		case rat64:
			ratprod.UpdateRat64(num)
		case bigRat:
			ratprod.UpdateBig(num)
		case Real:
			realprod.Update(num)
		default:
			nonnum = append(nonnum, num)
		}
	}

	resultElements := make([]Expr, 0, 2+len(nonnum))
	resultElements = append(resultElements, symbol.Times)

	// ints get turned into rationals
	if intprod.exists() && ratprod.exists() {
		// * small
		ratprod.UpdateInt64(intprod.prod)
		if intprod.bigprod.Sign() != 0 {
			ratprod.UpdateBigInt(intprod.bigprod)
		}
	} else if intprod.exists() {
		total := intprod.Total()
		if total.Sign() == 0 {
			return NewInteger(0)
		}

		if realprod.exists() {
			realprod.Update(Real(intprod.Total().Float64()))
		} else {
			if !total.IsInt64() || total.Int64() != 1 {
				resultElements = append(resultElements, total)
			}
		}
	}
	if ratprod.exists() {
		total := ratprod.Total()
		if total.Sign() == 0 {
			return NewInteger(0)
		}
		if realprod.exists() {
			// All ints and rationals are added
			// Convert to float, and let next part deal with it
			realprod.Update(Real(total.Float64()))
		} else {
			if total.IsInt() {
				resultElements = append(resultElements, total.Numerator())
			} else {
				resultElements = append(resultElements, total)
			}
		}
	}

	if realprod.exists() {
		total := realprod.Total()
		if total == 0.0 {
			return NewInteger(0)
		}
		resultElements = append(resultElements, total)
	}

	// Add non-numeric terms
	resultElements = append(resultElements, nonnum...)

	// Apply OneIdentity-like behavior: if only one element (plus head), return it
	if len(resultElements) == 2 {
		return resultElements[1]
	}

	// If no elements besides head, return 0
	if len(resultElements) == 1 {
		return NewInteger(1)
	}

	return NewListFromExprs(resultElements...)

}

// Adds a series of integers
type TimesInteger struct {
	prod    machineInt
	bigprod bigInt
	count   int
}

func (a *TimesInteger) exists() bool {
	return a.count > 0
}

func (a *TimesInteger) Update(b machineInt) {
	if prodnext, ok := timesInt64(a.prod.Int64(), b.Int64()); ok {
		a.prod = newMachineInt(prodnext)
		a.count += 1
		return
	}
	a.UpdateBig(a.prod.asBigInt())
	a.UpdateBig(b.asBigInt())
	a.prod = 0
}

func (a *TimesInteger) UpdateBig(b bigInt) {
	if a.bigprod.val == nil {
		a.bigprod = bigIntOne()
	}
	a.bigprod.times(b)
	a.count += 1
}

func (a *TimesInteger) Total() Integer {
	if a.bigprod.val == nil {
		return a.prod
	}
	a.bigprod.times(a.prod.asBigInt())
	return a.bigprod
}

// Product of a series of rational numbers
// No attempt at normalization to integers is done until the end.
type TimesRational struct {
	prod    rat64
	bigprod bigRat
	count   int
}

func (a *TimesRational) exists() bool {
	return a.count > 0
}

func (a *TimesRational) UpdateInt64(b machineInt) {
	if nextprod, ok := timesRat64Int64(a.prod, b); ok {
		a.prod = nextprod
		a.count += 1
		return
	}
	panic("Not implemented")
}

func (a *TimesRational) UpdateRat64(b rat64) {
	if prodnext, ok := timesRat64(a.prod, b); ok {
		a.prod = prodnext
		a.count += 1
		return
	}
	a.UpdateBig(a.prod.asBigRat())
	a.UpdateBig(b.asBigRat())
	a.prod = rat64One
}

func (a *TimesRational) UpdateBig(b bigRat) {
	if a.bigprod.val == nil {
		a.bigprod.val = big.NewRat(1, 1)
	}
	a.count += 1
	a.bigprod.times(b)
}

func (a *TimesRational) UpdateBigInt(b bigInt) {
	if a.bigprod.val == nil {
		a.bigprod.val = big.NewRat(1, 1)
	}
	a.count += 1

	// turn bigInt into big.Rat
	r := newBigRat(new(big.Rat).SetFrac(b.val, big.NewInt(1)))
	a.bigprod.times(r)
}

func (a *TimesRational) Total() Rational {
	if a.bigprod.val == nil {
		return a.prod
	}
	a.bigprod.times(a.prod.asBigRat())
	return a.bigprod
}

// Adds a series of floats
type TimesReal struct {
	prod  Real
	count int
}

func (a *TimesReal) Update(b Real) {
	a.prod *= b
	a.count += 1
}

func (a *TimesReal) Total() Real {
	return a.prod
}
func (a *TimesReal) exists() bool {
	return a.count > 0
}
