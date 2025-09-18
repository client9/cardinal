package big

import (
	"runtime"
	"unsafe"

	"github.com/client9/bignum/mpq"
	"github.com/client9/bignum/mpz"
)

type Rat struct {
	ptr mpq.RatPtr
}

// newRatPtr is where the magic happens.
// mpq_t is a struct of size 1, which is a pointer to the first element
//
//	this is exactly the same as "pointer to the struct"
//	same data, same layout.. but.. CGO gets cranky and thinks they are different.
//
// Fortuantely GMP doesn't actually use mpq_t but instead mpq_ptr
//
//	we'll use that, and isolate the pointer casting magic to this functiona
func newRatPtr() mpq.RatPtr {
	// new(mpg.Rat) is mpq_impl[1]
	// the pointer returns points at the first element, i.e. mpq_impl[0]
	// use unsafe to convert
	return mpq.RatPtr(unsafe.Pointer(new(mpq.Rat)))
}

func NewRat(a, b int64) *Rat {
	n := newRatPtr()
	mpq.Init(n)
	mpq.SetSi(n, int(a), uint(b))
	z := &Rat{
		ptr: n,
	}
	runtime.AddCleanup(z, mpq.Clear, n)
	return z
}

func NewRatTmp(a, b int64) *Rat {
	n := newRatPtr()
	mpq.Init(n)
	mpq.SetSi(n, int(a), uint(b))
	z := &Rat{
		ptr: n,
	}
	// NO CLEANUP
	//runtime.AddCleanup(z, mpq.Clear, &n)
	return z
}

func (z *Rat) init() {
	n := newRatPtr()
	mpq.Init(n)
	z.ptr = n
	runtime.AddCleanup(z, mpq.Clear, n)
}

func (z *Rat) Clear() {
	if z.ptr != nil {
		mpq.Clear(z.ptr)
		z.ptr = nil
	}
}

func (z *Rat) Abs(x *Rat) *Rat {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	mpq.Abs(z.ptr, x.ptr)
	return z
}

func (z *Rat) Add(x, y *Rat) *Rat {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpq.Add(z.ptr, x.ptr, y.ptr)
	return z
}

// TODO AppendText

func (x *Rat) Cmp(y *Rat) int {
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	return mpq.Cmp(x.ptr, y.ptr)
}

func (z *Rat) Denom() *Int {
	if z.ptr == nil {
		return NewInt(1)
	}
	// Do not set up Cleanup/Finalizer
	return &Int{
		ptr: mpq.Denref(z.ptr),
	}
}

func (z *Rat) Float32() float32 {
	// not f32 support in mpq
	return float32(z.Float64())
}

func (z *Rat) Float64() float64 {
	if z.ptr == nil {
		return 0.0
	}
	return mpq.GetD(z.ptr)
}

// TODO FloatPrec
// TODO FloatString
// TODO GobDecode
// TODO GobEncode

func (z *Rat) Inv(x *Rat) *Rat {
	if z.ptr == nil {
		z.init()
	}
	// TODO ERROR
	mpq.Inv(z.ptr, x.ptr)
	return z
}

func (z *Rat) IsInt() bool {
	if z.ptr == nil {
		return true
	}
	return mpz.CmpUi(mpq.Denref(z.ptr), 1) == 0
}

// TODO MarshalText

func (z *Rat) Mul(x, y *Rat) *Rat {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpq.Mul(z.ptr, x.ptr, y.ptr)
	return z
}

func (z *Rat) Neg(x *Rat) *Rat {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	mpq.Neg(z.ptr, x.ptr)
	return z
}

func (z *Rat) Num() *Int {
	if z.ptr == nil {
		return NewInt(0)
	}
	// Do not set up Cleanup/Finalizer
	return &Int{
		ptr: mpq.Numref(z.ptr),
	}
}

func (z *Rat) Quo(x, y *Rat) *Rat {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpq.Div(z.ptr, x.ptr, y.ptr)
	return z
}

// TOOD RatString
//    RETURNS PURE INT if denom == 1
// TODO SCAN

func (z *Rat) Set(x *Rat) *Rat {
	if z.ptr == nil {
		z.init()
	}
	mpq.Set(z.ptr, x.ptr)
	return z
}

func (z *Rat) SetFloat64(x float64) {
	if z.ptr == nil {
		n := newRatPtr()
		mpq.Init(n)
		mpq.SetD(n, x)
		z.ptr = n
		runtime.AddCleanup(z, mpq.Clear, n)
	} else {
		mpq.SetD(z.ptr, x)
	}
}

// TODO SETFRAC
func (z *Rat) SetFrac(a, b *Int) *Rat {
	return nil
}

func (z *Rat) SetFrac64(a, b int64) *Rat {
	if z.ptr == nil {
		z.init()
	}
	if b < 0 {
		a = -a
		b = -b
	}
	mpq.SetUi(z.ptr, uint(a), uint(b))
	return z
}

func (z *Rat) SetInt(x *Int) *Rat {
	if z.ptr == nil {
		n := newRatPtr()
		mpq.Init(n)
		z.ptr = n
		runtime.AddCleanup(z, mpq.Clear, n)
	}
	mpq.SetZ(z.ptr, x.ptr)
	return z
}

func (z *Rat) SetInt64(x int64) *Rat {
	if z.ptr == nil {
		n := newRatPtr()
		mpq.Init(n)
		mpq.SetSi(n, int(x), uint(1))
		z.ptr = n
		runtime.AddCleanup(z, mpq.Clear, n)
	} else {
		mpq.SetSi(z.ptr, int(x), uint(1))
	}
	return z
}

// TODO SETSTRING
func (z *Rat) SetString(s string) (*Rat, bool) {
	return nil, false
}

func (z *Rat) SetUint64(x uint64) {
	if z.ptr == nil {
		n := newRatPtr()
		mpq.Init(n)
		mpq.SetUi(n, uint(x), 1)
		z.ptr = n
		runtime.AddCleanup(z, mpq.Clear, n)
	} else {
		mpq.SetUi(z.ptr, uint(x), uint(1))
	}
}

func (z *Rat) Sign() int {
	if z.ptr == nil {
		return 0
	}
	return mpq.Sgn(z.ptr)
}

func (z *Rat) String() string {
	if z.ptr == nil {
		return ""
	}
	// Not in mpq -- calls custom C code
	return mpq.GetStr(10, z.ptr)
}

func (z *Rat) Sub(x, y *Rat) *Rat {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpq.Sub(z.ptr, x.ptr, y.ptr)
	return z
}

// TODO UNMARSHALTEXT
