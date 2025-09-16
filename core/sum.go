package core

import (
	"github.com/client9/cardinal/core/symbol"
)

func PlusList(args []Expr) Expr {

	intsum := AccumulatorInteger{}

	ratsum := AccumulatorRational{
		sum: rat64Zero,
	}
	realsum := AccumulatorFloat64{}

	bigreal := AccumulatorBigFloat{}

	var nonnum []Expr

	for _, arg := range args {
		switch num := arg.(type) {
		case machineInt:
			intsum.Plus(num)
		case BigInt:
			intsum.PlusBig(num)
		case rat64:
			ratsum.PlusRat64(num)
		case BigRat:
			ratsum.PlusBigRat(num)
		case f64:
			realsum.PlusFloat64(num.Float64())
		case BigFloat:
			bigreal.Plus(num)
		default:
			nonnum = append(nonnum, num)
		}
	}

	resultElements := make([]Expr, 0, 2+len(nonnum))
	resultElements = append(resultElements, symbol.Plus)

	// it machine float64, then everything gets turninto to float64
	// if bigFloat, then everything gets turned into BigFloat
	// if int and rat, convert to one rat
	// if rat, add
	// if int, add
	// it not numeric, add

	// if we have machine precision float
	// then integer and rational stuff
	// compute and return

	var total Number

	// if machine float64 is found, we have just lower everything to float64
	if realsum.exists() {
		// if we have both types, add them into one rational
		// convert to float and update floatsum.
		if intsum.exists() && ratsum.exists() {
			ratsum.PlusInt64(intsum.sum)
			if intsum.bigcount {
				ratsum.PlusBigInt(intsum.bigsum)
			}
			realsum.PlusFloat64(ratsum.Total().Float64())
		} else if intsum.exists() {
			realsum.PlusFloat64(intsum.Total().Float64())
		} else if ratsum.exists() {
			realsum.PlusFloat64(ratsum.Total().Float64())
		}
		if bigreal.exists() {
			realsum.PlusFloat64(bigreal.Total().Float64())
		}

		total = realsum.Total()
	} else if bigreal.exists() {
		// don't have to worry about f64 since handled above
		// don't have to worry about glue together int and rat
		//  since they'll both be promoted to Big values anyways
		if intsum.exists() {
			bigreal.PlusInt(intsum.Total())
		}
		if ratsum.exists() {
			bigreal.PlusRat(ratsum.Total())
		}
		total = bigreal.Total()
	} else if intsum.exists() && ratsum.exists() {
		ratsum.PlusInt64(intsum.sum)
		if intsum.bigcount {
			ratsum.PlusBigInt(intsum.bigsum)
		}
		total = ratsum.Total()
	} else if ratsum.exists() {
		total = ratsum.Total()
	} else if intsum.exists() {
		total = intsum.Total()
	}
	if len(nonnum) == 0 || (total != nil && total.Sign() != 0) {
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
		return NewInteger(0)
	}

	return NewListFromExprs(resultElements...)
}

func TimesList(args []Expr) Expr {
	intsum := AccumulatorInteger{
		sum: newMachineInt(1),
	}
	ratsum := AccumulatorRational{
		sum: rat64One,
	}
	realsum := AccumulatorFloat64{
		sum: 1.0,
	}
	bigreal := AccumulatorBigFloat{
		// lazy initiziation
	}
	var nonnum []Expr

	for _, arg := range args {
		switch num := arg.(type) {
		case machineInt:
			intsum.TimesInt64(num)
		case BigInt:
			intsum.TimesBigInt(num)
		case rat64:
			ratsum.TimesRat64(num)
		case BigRat:
			ratsum.TimesBigRat(num)
		case f64:
			realsum.TimesFloat64(num.Float64())
		case BigFloat:
			bigreal.Times(num)
		default:
			nonnum = append(nonnum, num)
		}
	}

	resultElements := make([]Expr, 0, 2+len(nonnum))
	resultElements = append(resultElements, symbol.Times)

	var total Number

	// if machine float64 is found, we have just lower everything to float64
	if realsum.exists() {
		// if we have both types, add them into one rational
		// convert to float and update floatsum.
		if intsum.exists() && ratsum.exists() {
			ratsum.TimesInt64(intsum.sum)
			if intsum.bigcount {
				ratsum.TimesBigInt(intsum.bigsum)
			}
			realsum.TimesFloat64(ratsum.Total().Float64())
		} else if intsum.exists() {
			realsum.TimesFloat64(intsum.Total().Float64())
		} else if ratsum.exists() {
			realsum.TimesFloat64(ratsum.Total().Float64())
		}
		if bigreal.exists() {
			realsum.TimesFloat64(bigreal.Total().Float64())
		}

		total = realsum.Total()
	} else if bigreal.exists() {
		// don't have to worry about f64 since handled above
		// don't have to worry about glue together int and rat
		//  since they'll both be promoted to Big values anyways
		if intsum.exists() {
			bigreal.TimesInt(intsum.Total())
		}
		if ratsum.exists() {
			bigreal.TimesRat(ratsum.Total())
		}

		total = bigreal.Total()
	} else if intsum.exists() && ratsum.exists() {
		ratsum.TimesInt64(intsum.sum)
		if intsum.bigcount {
			ratsum.TimesBigInt(intsum.bigsum)
		}
		total = ratsum.Total()
	} else if ratsum.exists() {
		total = ratsum.Total()
	} else if intsum.exists() {
		total = intsum.Total()
	}
	if total != nil {
		if total.Sign() == 0 {
			return NewInteger(0)
		}

		if r, ok := total.(Rational); ok && r.IsInt() {
			total = r.Numerator()
		}

		if len(nonnum) == 0 {
			// if we don't have any non-numerical arguments, add it
			resultElements = append(resultElements, total)
		} else if val, ok := total.(Integer); ok {
			if !val.IsInt64() || val.Int64() != 1 {
				// if total is a integer, don't return Times(1, x)
				//  although this type of rule could be
				// done outside of here as part of general simplification
				resultElements = append(resultElements, total)
			}
		} else {
			// not the integer '1', add it
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
		return NewInteger(1)
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
type AccumulatorInteger struct {
	sum      machineInt
	bigsum   BigInt
	count    bool
	bigcount bool
}

func (a *AccumulatorInteger) exists() bool {
	return a.count || a.bigcount
}

func (a *AccumulatorInteger) Plus(b machineInt) {
	if sumnext, ok := addInt64(a.sum.Int64(), b.Int64()); ok {
		a.count = true
		a.sum = newMachineInt(sumnext)
		return
	}
	a.PlusBig(a.sum.asBigInt())
	a.PlusBig(b.asBigInt())
	a.sum = 0
}

func (a *AccumulatorInteger) PlusBig(b BigInt) {
	if !a.bigcount {
		a.bigsum.Set(&b)
		a.bigcount = true
		return
	}
	a.bigsum.add(b)
}

func (a *AccumulatorInteger) TimesInt64(b machineInt) {
	a.count = true
	if prodnext, ok := timesInt64(a.sum.Int64(), b.Int64()); ok {
		a.count = true
		a.sum = newMachineInt(prodnext)
		return
	}
	a.TimesBigInt(a.sum.asBigInt())
	a.TimesBigInt(b.asBigInt())
	a.sum = 0
}

func (a *AccumulatorInteger) TimesBigInt(b BigInt) {
	if !a.bigcount {
		a.bigsum.Set(&b)
		a.bigcount = true
		return
	}
	a.bigsum.times(b)
}

func (a *AccumulatorInteger) Total() Integer {
	if !a.bigcount {
		return a.sum
	}

	if a.sum != 0 {
		a.bigsum.add(a.sum.asBigInt())
	}
	return a.bigsum
}

// Adds a series of integers
type AccumulatorRational struct {
	sum    rat64
	bigsum BigRat

	count    bool
	bigcount bool
}

func (a *AccumulatorRational) exists() bool {
	return a.count || a.bigcount
}
func (a *AccumulatorRational) PlusRat64(b rat64) {
	if sumnext, ok := addRat64(a.sum, b); ok {
		a.sum = sumnext
		a.count = true
		return
	}
	a.PlusBigRat(a.sum.AsBigRat())
	a.PlusBigRat(b.AsBigRat())
	a.sum = rat64Zero
}

func (a *AccumulatorRational) PlusInt64(b machineInt) {
	if nextsum, ok := addRat64Int64(a.sum, b); ok {
		a.sum = nextsum
		a.count = true
		return
	}
	a.PlusBigRat(a.sum.AsBigRat())
	a.sum = rat64Zero
	a.PlusBigInt(b.AsBigInt())
}

func (a *AccumulatorRational) PlusBigInt(b BigInt) {
	if !a.bigcount {
		a.bigcount = true
		a.bigsum.SetInt(&b)
		return
	}
	a.bigsum.AddInt(&a.bigsum, &b)
}
func (a *AccumulatorRational) PlusBigRat(b BigRat) {
	if !a.count {
		a.bigcount = true
		a.bigsum.Set(&b)
	}
	a.bigsum.Add(&a.bigsum, &b)
}

func (a *AccumulatorRational) TimesInt64(b machineInt) {
	if nextprod, ok := timesRat64Int64(a.sum, b); ok {
		a.sum = nextprod
		a.count = true
		return
	}
	panic("Not implemented")
}

func (a *AccumulatorRational) TimesRat64(b rat64) {
	if prodnext, ok := timesRat64(a.sum, b); ok {
		a.sum = prodnext
		a.count = true
		return
	}
	a.TimesBigRat(a.sum.AsBigRat())
	a.TimesBigRat(b.AsBigRat())
	a.sum = rat64One
}

func (a *AccumulatorRational) TimesBigRat(b BigRat) {
	if !a.bigcount {
		a.bigsum.Set(&b)
		a.bigcount = true
		return
	}
	a.bigsum.Mul(&a.bigsum, &b)
}

func (a *AccumulatorRational) TimesBigInt(b BigInt) {
	if !a.bigcount {
		a.bigsum.SetInt(&b)
		a.bigcount = true
		return
	}

	a.bigsum.MulInt(&a.bigsum, &b)
}

func (a *AccumulatorRational) Total() Rational {
	if !a.bigcount {
		return a.sum
	}
	tmp := a.sum.AsBigRat()

	a.bigsum.Add(&a.bigsum, &tmp)
	return a.bigsum
}

// Adds a series of floats
type AccumulatorFloat64 struct {
	sum   float64
	count int
}

func (a *AccumulatorFloat64) exists() bool {
	return a.count > 0
}
func (a *AccumulatorFloat64) PlusFloat64(b float64) {
	a.sum += b
	a.count += 1
}
func (a *AccumulatorFloat64) TimesFloat64(b float64) {
	a.sum *= b
	a.count += 1
}

func (a *AccumulatorFloat64) Total() Real {
	return NewReal(a.sum)
}

type AccumulatorBigFloat struct {
	sum   BigFloat
	count int
}

func (a *AccumulatorBigFloat) exists() bool {
	return a.count > 0
}

func (a *AccumulatorBigFloat) Plus(b BigFloat) {
	if a.count == 0 {
		a.sum.Set(&b)
		a.count += 1
		return
	}
	a.sum.Add(&a.sum, &b)
}

func (a *AccumulatorBigFloat) Times(b BigFloat) {
	if a.count == 0 {
		a.sum.Set(&b)
		a.count += 1
		return
	}

	a.sum.Mul(&a.sum, &b)
}

func (a *AccumulatorBigFloat) PlusInt(n Integer) {
	ba := n.AsBigInt()
	a.sum.Add(&a.sum, new(BigFloat).SetInt(&ba))
}

func (a *AccumulatorBigFloat) TimesInt(n Integer) {
	ba := n.AsBigInt()
	a.sum.Mul(&a.sum, new(BigFloat).SetInt(&ba))
}

func (a *AccumulatorBigFloat) PlusRat(n Rational) {
	ba := n.AsBigRat()
	a.sum.Add(&a.sum, new(BigFloat).SetRat(&ba))
}

func (a *AccumulatorBigFloat) TimesRat(n Rational) {
	ba := n.AsBigRat()
	a.sum.Mul(&a.sum, new(BigFloat).SetRat(&ba))
}

func (a *AccumulatorBigFloat) Total() Real {
	return a.sum
}
