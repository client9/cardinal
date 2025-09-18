package big

import (
	"github.com/client9/bignum/mpq"
)

func (z *Rat) AddInt(x *Rat, y *Int) *Rat {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}

	// convert integer to rational
	// TODO can prob allocate on the stack
	num := new(Rat)
	mpq.SetZ(num.ptr, y.ptr)

	mpq.Add(z.ptr, x.ptr, num.ptr)

	return z
}
func (z *Rat) MulInt(x *Rat, y *Int) *Rat {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}

	// convert integer to rational
	// TODO can prob allocate on the stack
	num := new(Rat)
	mpq.SetZ(num.ptr, y.ptr)

	mpq.Mul(z.ptr, x.ptr, num.ptr)

	return z
}
