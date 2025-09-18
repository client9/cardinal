package big

import (
	"fmt"
	"math"
	stdlib "math/big"
	"runtime"
	"strings"
	"unsafe"

	"github.com/client9/bignum/mpfr"
)

func newFloatPtr() mpfr.FloatPtr {
	return mpfr.FloatPtr(unsafe.Pointer(new(mpfr.Float)))
}

type Float struct {
	ptr mpfr.FloatPtr

	// TODO: storing prec and mode natively takes 16 bytes on 64-bit platforms
	//   * No need for 64-bit values
	//   * prec should be a uint32
	//   * mode could also be a uint32
	// Total: 8 bytes
	prec uint
	mode mpfr.RoundMode
}

func NewFloat(x float64) *Float {
	n := newFloatPtr()
	mpfr.InitSetD(n, x, mpfr.RNDN)
	z := &Float{
		ptr:  n,
		prec: 0,
		mode: 0,
	}
	runtime.AddCleanup(z, mpfr.Clear, n)
	return z
}

// TODO PARSEFLOAT

func (z *Float) init() {
	if z.ptr == nil {
		n := newFloatPtr()
		mpfr.InitSetD(n, 0, mpfr.RNDN)
		z.ptr = n
		runtime.AddCleanup(z, mpfr.Clear, n)
	}
}
func (z *Float) Abs(x *Float) *Float {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	mpfr.Abs(z.ptr, x.ptr, z.mode)
	return z
}

// TODO ACC

func (z *Float) Add(x, y *Float) *Float {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	if z.prec == 0 {
		z.prec = max(x.prec, y.prec)
	}
	mpfr.Add(z.ptr, x.ptr, y.ptr, z.mode)
	return z
}

// TODO APPEND
// TODO APPENDTEXT

// NONSTANDARD
// TODO make sure mpfr is ok with clearing Nul pointer.
func (z *Float) Clear() {
	if z.ptr != nil {
		mpfr.Clear(z.ptr)
		z.ptr = nil
		z.prec = 0
		z.mode = 0
	}
}

func (x *Float) Cmp(y *Float) int {
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	return mpfr.Cmp(x.ptr, y.ptr)
}

// TODO COPY

func (z *Float) Float32() float32 {
	if z.ptr == nil {
		return 0.0
	}
	return mpfr.GetFlt(z.ptr, z.mode)
}
func (z *Float) Float64() float64 {
	if z.ptr == nil {
		return 0.0
	}
	return mpfr.GetD(z.ptr, z.mode)
}

// TODO FORMAT
// TODO GOBDECODE
// TODO GOBENCODE

func (z *Float) Int() int {
	if z.ptr == nil {
		return 0
	}
	return int(mpfr.GetSi(z.ptr, z.mode))
}
func (z *Float) Int64() int64 {
	if z.ptr == nil {
		return 0
	}
	return int64(mpfr.GetSi(z.ptr, z.mode))
}

func (z *Float) IsInf() bool {
	if z.ptr == nil {
		return false
	}
	return mpfr.InfP(z.ptr) == 1
}

func (z *Float) IsInt() bool {
	if z.ptr == nil {
		return true
	}
	return mpfr.GetSi(z.ptr, z.mode) == 1
}

// TODO MANTEXP
// TODO MARSHALTEXT

func (z *Float) MinPrec() uint {
	if z.ptr == nil {
		return 0
	}
	// doesnt matter if z is initialized
	return mpfr.MinPrec(z.ptr)
}

func (z *Float) Mode() stdlib.RoundingMode {
	// doesnt matter if z is initialized
	return exportRoundingMode(z.mode)
}
func (z *Float) Mul(x, y *Float) *Float {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	if z.prec == 0 {
		z.prec = max(x.prec, y.prec)
	}
	mpfr.Mul(z.ptr, x.ptr, y.ptr, z.mode)
	return z
}

func (z *Float) Neg(x *Float) *Float {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	mpfr.Neg(z.ptr, x.ptr, z.mode)
	return z
}

// TODO PARSE

func (z *Float) Prec() uint {
	// doesn't matter if z is initialized or not
	return z.prec
}

func (z *Float) Quo(x, y *Float) *Float {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	if z.prec == 0 {
		z.prec = max(x.prec, y.prec)
	}
	mpfr.Div(z.ptr, x.ptr, y.ptr, z.mode)
	return z
}

// TODO RAT
// TODO SCAN

func (z *Float) Set(x *Float) *Float {
	if z.ptr == nil {
		n := newFloatPtr()
		mpfr.InitSet(n, x.ptr, mpfr.RNDN)
		z.ptr = n
		runtime.AddCleanup(z, mpfr.Clear, n)
		return z
	}
	z.prec = x.prec
	mpfr.Set(z.ptr, x.ptr, z.mode)
	return z
}

func (z *Float) SetFloat64(d float64) *Float {
	if z.ptr == nil {
		n := newFloatPtr()
		mpfr.InitSetD(n, d, mpfr.RNDN)
		z.ptr = n
		runtime.AddCleanup(z, mpfr.Clear, n)
		return z
	}
	mpfr.SetD(z.ptr, d, z.mode)
	return z
}

func (z *Float) SetInf(signbit bool) *Float {
	if z.ptr == nil {
		z.init()
	}
	if signbit {
		mpfr.SetInf(z.ptr, 1)
	} else {
		mpfr.SetInf(z.ptr, 0)
	}
	return z
}

func (z *Float) SetInt(x *Int) *Float {
	if z.ptr == nil {
		n := newFloatPtr()
		mpfr.SetZ(n, x.ptr, mpfr.RNDN)
		z.ptr = n
		runtime.AddCleanup(z, mpfr.Clear, n)
		return z
	}
	mpfr.SetZ(z.ptr, x.ptr, z.mode)
	return z
}

func (z *Float) SetInt64(d int64) *Float {
	if z.ptr == nil {
		n := newFloatPtr()
		mpfr.InitSetSi(n, d, mpfr.RNDN)
		z.ptr = n
		runtime.AddCleanup(z, mpfr.Clear, n)
		return z
	}

	mpfr.SetSi(z.ptr, d, z.mode)
	return z
}

// TODO SETMANTEXP

func (z *Float) SetMode(mode stdlib.RoundingMode) *Float {
	if z.ptr == nil {
		z.init()
	}
	z.mode = importRoundingMode(mode)
	return z
}

func (z *Float) SetPrec(prec uint) *Float {
	if z.ptr == nil {
		z.init()
	}
	mpfr.SetPrec(z.ptr, int(prec))
	z.prec = prec
	return z
}

func (z *Float) SetRat(x *Rat) *Float {
	if z.ptr == nil {
		n := newFloatPtr()
		z.ptr = n
		runtime.AddCleanup(z, mpfr.Clear, n)
	}
	mpfr.SetQ(z.ptr, x.ptr, z.mode)
	return z
}

// TODO BETTER PRECISION
func (z *Float) SetString(s string) (*Float, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("empty string")
	}
	integerDigitCount := strings.Index(s, ".")
	if integerDigitCount == -1 {
		return nil, fmt.Errorf("not a float")
	}
	digits := len(s) - 1
	precision := uint(math.Ceil(float64(digits) * math.Log2(10.0)))

	if z.ptr == nil {
		n := newFloatPtr()
		z.ptr = n
		runtime.AddCleanup(z, mpfr.Clear, n)
	}
	z.SetPrec(precision)
	if mpfr.SetStr(z.ptr, s, 10, z.mode) != 0 {
		return nil, fmt.Errorf("float conversion failed")
	}
	return z, nil
}

func (z *Float) SetUnt64(d uint64) {
	if z.ptr == nil {
		n := newFloatPtr()
		mpfr.InitSetUi(n, d, mpfr.RNDN)
		z.ptr = n
		runtime.AddCleanup(z, mpfr.Clear, n)
		return
	}

	mpfr.SetUi(z.ptr, d, z.mode)
}

func (x *Float) Sign() int {
	if x.ptr == nil {
		return 0
	}
	return mpfr.Sgn(x.ptr)
}

func (z *Float) Signbit() bool {
	if z.ptr == nil {
		return false
	}
	return 1 == mpfr.Signbit(z.ptr)
}

func (z *Float) Sqrt(x *Float) *Float {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}

	if z.prec == 0 {
		z.prec = x.prec
	}
	mpfr.Sqrt(z.ptr, x.ptr, z.mode)
	return z
}

func (z *Float) String() string {
	if z.ptr == nil {
		return ""
	}
	// matches Go
	//return mpfr.Sprintf3("%.10R*g", z.mode, z.ptr)
	// matches precision
	return mpfr.Sprintf3("%R*e", z.mode, z.ptr)

	// all digits
	s, _ := mpfr.GetStr(10, 0, z.ptr, z.mode)
	return s

}

func (z *Float) Sub(x, y *Float) *Float {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	if z.prec == 0 {
		z.prec = max(x.prec, y.prec)
	}
	mpfr.Sub(z.ptr, x.ptr, y.ptr, z.mode)
	return z
}

// TODO TEXT

func (z *Float) Unt64() uint64 {
	if z.ptr == nil {
		return 0
	}
	return uint64(mpfr.GetUi(z.ptr, z.mode))
}

// TODO UNMARSHALTEXT

func max(a, b uint) uint {
	if a >= b {
		return a
	}
	return b
}

func importRoundingMode(r stdlib.RoundingMode) mpfr.RoundMode {
	switch r {
	case stdlib.ToNearestEven:
		return mpfr.RNDN
	case stdlib.ToNearestAway:
		panic("ToNearestAway RoundingMode not supported")
	case stdlib.ToZero:
		return mpfr.RNDZ
	case stdlib.AwayFromZero:
		return mpfr.RNDA
	case stdlib.ToNegativeInf:
		return mpfr.RNDD
	case stdlib.ToPositiveInf:
		return mpfr.RNDU
	default:
		panic("unknown rounding mode")
	}
}

func exportRoundingMode(r mpfr.RoundMode) stdlib.RoundingMode {
	switch r {
	case mpfr.RNDN:
		return stdlib.ToNearestEven
	case mpfr.RNDZ:
		return stdlib.ToZero
	case mpfr.RNDA:
		return stdlib.AwayFromZero
	case mpfr.RNDD:
		return stdlib.ToNegativeInf
	case mpfr.RNDU:
		return stdlib.ToPositiveInf
	default:
		panic("unsupported MPFR rounding mode")
	}
}
