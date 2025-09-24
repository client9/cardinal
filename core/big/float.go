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
		prec: 53,
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

func (z *Float) Int() *Int {
	i := NewInt(0)
	if z.ptr == nil {
		return i
	}
	// round to zero
	mpfr.GetZ(i.ptr, z.ptr, mpfr.RNDZ)
	return i
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
		z.prec = x.prec
		runtime.AddCleanup(z, mpfr.Clear, n)
		return z
	}
	mpfr.Set(z.ptr, x.ptr, z.mode)
	z.prec = x.prec
	return z
}

func (z *Float) SetFloat64(d float64) *Float {
	if z.ptr == nil {
		n := newFloatPtr()
		mpfr.InitSetD(n, d, mpfr.RNDN)
		z.ptr = n
		z.prec = 53
		runtime.AddCleanup(z, mpfr.Clear, n)
		return z
	}
	mpfr.SetD(z.ptr, d, z.mode)
	z.prec = 53
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
		mpfr.InitSetZ(n, x.ptr, mpfr.RNDN)
		z.ptr = n
		z.prec = uint(x.BitLen())
		if z.prec < 64 {
			z.prec = 64
		}
		runtime.AddCleanup(z, mpfr.Clear, n)
		return z
	}
	mpfr.SetZ(z.ptr, x.ptr, z.mode)
	if z.prec == 0 {
		z.prec = uint(x.BitLen())
		if z.prec < 64 {
			z.prec = 64
		}
	}
	return z
}

func (z *Float) SetInt64(d int64) *Float {
	if z.ptr == nil {
		n := newFloatPtr()
		mpfr.InitSetSi(n, d, mpfr.RNDN)
		z.ptr = n
		z.prec = 64
		runtime.AddCleanup(z, mpfr.Clear, n)
		return z
	}
	mpfr.SetSi(z.ptr, d, z.mode)
	if z.prec == 0 {
		z.prec = 64
	}
	return z
}

func (z *Float) SetMantExp(mant *Float, exp int) *Float {
	z.Set(mant)
	mpfr.Mul2si(z.ptr, z.ptr, exp, mpfr.RNDN)
	return z
}

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
	mpfr.PrecRound(z.ptr, int(prec), z.mode)
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

func (z *Float) SetString(s string) (*Float, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("empty string")
	}
	if z.ptr == nil {
		z.init()
	}

	// HACK
	// Assume base 10 for now

	s2 := s
	if s2[0] == '-' || s2[0] == '+' {
		s2 = s2[1:]
	}
	idx := strings.IndexAny(s2, "eE")
	if idx != -1 {
		s2 = s2[:idx]
	}
	prec := len(s2)
	if strings.IndexByte(s2, '.') != -1 {
		prec -= 1
	}
	p := uint(math.Floor(float64(prec) / math.Log10(2.0)))
	// do before setting
	z.SetPrec(p)
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

	d := z.Float64()
	// if bigger than 10^6, or smaller than 10^-6 use exponential form
	if (d <= 0.000001) || (d >= 1000000) {
		return mpfr.Sprintf3("%R*e", mpfr.RNDN, z.ptr)
	}

	// MPFR in "g" mode switches to "e" if d < 0.0001
	// Want 10^-6
	if d <= 0.0001 {
		s, e := mpfr.GetStr(10, 0, z.ptr, mpfr.RNDN)
		return "0." + strings.Repeat("0", -e) + s
	}
	// stdlib uses bits for precision, convert to decimal digits
	prec := int(math.Ceil(float64(z.Prec()) * math.Log10(2)))
	// make the format string
	fstr := fmt.Sprintf("%%.%dR*G", prec)

	// and make the output
	return mpfr.Sprintf3(fstr, mpfr.RNDN, z.ptr)

	// matches Go
	//return mpfr.Sprintf3("%.10R*g", z.mode, z.ptr)
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
