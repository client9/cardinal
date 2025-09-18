package big

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/client9/bignum/mpz"
)

func newIntPtr(val int) mpz.IntPtr {
	return mpz.New(val)
	//return mpz.IntPtr(unsafe.Pointer(new(mpz.Int)))
}

type Int struct {
	ptr mpz.IntPtr
}

func NewInt(x int64) *Int {
	n := newIntPtr(int(x))
	z := &Int{
		ptr: n,
	}
	runtime.AddCleanup(z, mpz.Delete, n)
	return z
}

func (z *Int) init() {
	n := newIntPtr(0)
	z.ptr = n
	runtime.AddCleanup(z, mpz.Delete, n)
}

func NewIntTmp(x int64) *Int {
	n := newIntPtr(int(x))
	z := &Int{
		ptr: n,
	}
	return z
}
func (z *Int) Clear() {
	mpz.Delete(z.ptr)
}

func (z *Int) Abs(x *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	mpz.Abs(z.ptr, x.ptr)
	return z
}

func (z *Int) Add(x, y *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpz.Add(z.ptr, x.ptr, y.ptr)
	return z
}

// TODO ANDNOT -- write C so we don't make two round trips.
// TODO APPEND
// TODO APPENDTEXT

func (z *Int) Binomial(n, k int64) *Int {
	if z.ptr == nil {
		z.init()
	}

	// normalize
	if k < 0 {
		n = -n
		k = -k
	}
	if n > 0 {
		mpz.BinUiUi(z.ptr, uint(n), uint(k))
		return z
	}

	// TODO make sure N+k-1 is legit

	mpz.BinUiUi(z.ptr, uint(n+k-1), uint(k))
	// if k is even..
	if k&1 == 0 {
		mpz.Neg(z.ptr, z.ptr)
	}
	return z
}

// Bit returns the value of the i'th bit of z. That is, it
// returns (z>>i)&1. The bit index i must be >= 0.
func (z *Int) Bit(i int) uint {
	if i < 0 {
		panic("negative bit index")
	}
	if z.ptr == nil {
		return 0
	}
	return uint(mpz.Tstbit(z.ptr, uint(i)))
}

func (z *Int) BitLen() int {
	if z.ptr == nil {
		return 0
	}
	return mpz.Sizeinbase(z.ptr, 2)
}

// TODO BITS

func (z *Int) Bytes() []byte {
	if z.ptr == nil {
		return nil
	}
	return mpz.ExportGo(z.ptr)
}

func (z *Int) Cmp(y *Int) int {
	if z.ptr == nil {
		z.init()
	}
	if y.ptr == nil {
		y.init()
	}

	return mpz.Cmp(z.ptr, y.ptr)
}

func (z *Int) CmpAbs(y *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if y.ptr == nil {
		y.init()
	}

	mpz.Cmpabs(z.ptr, y.ptr)
	return z
}

func (z *Int) Div(x, y *Int) *Int {
	sgn := y.Sign()
	if sgn == 0 {
		panic("division by zero")
	}
	if z.ptr == nil {
		z.init()
	}

	// TODO: if x or m is unini we can skip steps
	if x.ptr == nil {
		x.init()
	}

	if sgn == 1 {
		mpz.FdivR(z.ptr, x.ptr, y.ptr)
	} else {
		mpz.CdivR(z.ptr, x.ptr, y.ptr)
	}
	return z
}

func (z *Int) DivMod(x, y, m *Int) (*Int, *Int) {
	// Sign handles uninitialized values so this is ok
	sgn := y.Sign()
	if sgn == 0 {
		panic("division by zero")
	}

	if z.ptr == nil {
		z.init()
	}

	// TODO: if x or m is unini we can skip steps
	if x.ptr == nil {
		x.init()
	}
	if m.ptr == nil {
		m.init()
	}

	if sgn == 1 {
		mpz.FdivQr(z.ptr, m.ptr, x.ptr, y.ptr)
	} else {
		mpz.CdivQr(z.ptr, m.ptr, x.ptr, y.ptr)
	}
	return z, m
}

func (z *Int) Exp(x, y, m *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	//
	// it's ok for m to be nil or uninitialized
	//
	if y.Sign() <= 0 {
		mpz.SetUi(z.ptr, 1)
		return z
	}

	if m == nil || m.Sign() == 0 {
		mpz.PowUi(z.ptr, x.ptr, mpz.GetUi(y.ptr))
	} else {
		mpz.Powm(z.ptr, x.ptr, y.ptr, m.ptr)
	}
	return z
}

// TODO FILLBYTES

func (z *Int) Float64() float64 {
	if z.ptr == nil {
		return 0.0
	}
	return mpz.GetD(z.ptr)
}

// TODO FORMAT

func (z *Int) GCD(x, y, a, b *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if a.ptr == nil {
		a.init()
	}
	if b.ptr == nil {
		b.init()
	}
	if a.Sign() <= 0 || b.Sign() <= 0 {
		z.SetInt64(0)
		if x != nil {
			x.SetInt64(0)
		}
		if y != nil {
			y.SetInt64(0)
		}
		return z
	}
	if x == nil && y == nil {
		mpz.Gcd(z.ptr, a.ptr, b.ptr)
		return z
	}

	if x != nil {
		if x.ptr == nil {
			x.init()
		}
	} else {
		x.SetInt64(0)
	}
	if y != nil {
		if y.ptr == nil {
			y.init()
		}
	} else {
		y.SetInt64(0)
	}

	mpz.Gcdext(z.ptr, x.ptr, y.ptr, a.ptr, b.ptr)
	return z
}

// TODO GOBDECODE
// TODO GOBENCODE

func (z *Int) Int64() int64 {
	if z.ptr == nil {
		return 0
	}
	return int64(mpz.GetSi(z.ptr))
}

func (z *Int) IsInt64() bool {
	if z.ptr == nil {
		return true
	}
	return mpz.FitsSlongP(z.ptr) == 1
}

func (z *Int) IsUint64() bool {
	if z.ptr == nil {
		return true
	}
	return mpz.FitsUlongP(z.ptr) == 1
}

// Lsh sets z = x << n and returns z.
func (z *Int) Lsh(x *Int, n uint) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	mpz.Mul2exp(z.ptr, x.ptr, n)
	return z
}

// MarshalJSON implements the json.Marshaler interface.
func (z *Int) MarshalJSON() ([]byte, error) {
	return []byte(z.String()), nil
}

// TODO MARSHALTEXT

func (z *Int) Mod(x, y *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpz.Mod(z.ptr, x.ptr, y.ptr)
	return z
}

func (z *Int) ModInverse(x, y *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpz.Invert(z.ptr, x.ptr, y.ptr)
	return z
}

// TODO MODSQRT

// TODO: if any are uninitialized, we can return 0
func (z *Int) Mul(x, y *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpz.Mul(z.ptr, x.ptr, y.ptr)
	return z
}

// MulRange sets z to the product of all integers
// in the range [a, b] inclusively and returns z.
// If a > b (empty range), the result is 1.
func (z *Int) MulRange(a, b int64) *Int {
	switch {
	case a > b:
		return z.SetInt64(1) // empty range
	case a <= 0 && b >= 0:
		return z.SetInt64(0) // range includes 0
	}
	// a <= b && (b < 0 || a > 0)

	// standard factorial
	if a == 1 && b >= 1 {
		mpz.FacUi(z.ptr, uint(b))
		return z
	}

	//if a > 1 && b >= 1 {
	// depends if factorial steps are cached or not
	// TBD  b! / (a-1)!
	//}

	// Slow
	z.SetInt64(a)
	for i := a + 1; i <= b; i++ {
		mpz.MulSi(z.ptr, z.ptr, int(i))
	}
	return z
}

// TODO: Whats Neg(0), if 0, then we can inline
func (z *Int) Neg(x *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	mpz.Neg(z.ptr, x.ptr)
	return z
}

func (z *Int) Not(x *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	mpz.Com(z.ptr, x.ptr)
	return z
}

func (z *Int) Or(x, y *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpz.Ior(z.ptr, x.ptr, y.ptr)
	return z
}

func (z *Int) ProbablyPrime(n int) bool {
	if z.ptr == nil {
		return false
	}
	return mpz.ProbabPrimeP(z.ptr, n) == 1
}

func (z *Int) Quo(x, y *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpz.TdivQ(z.ptr, x.ptr, y.ptr)
	return z
}

func (z *Int) QuoRem(x, y, r *Int) (*Int, *Int) {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	if r.ptr == nil {
		r.init()
	}
	mpz.TdivQr(z.ptr, r.ptr, x.ptr, y.ptr)
	return z, r
}

// TODO RAND

func (z *Int) Rem(x, y *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpz.TdivR(z.ptr, x.ptr, y.ptr)
	return z
}

// Rsh sets z = x >> n and returns z.
func (z *Int) Rsh(x *Int, n uint) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	// allows negatives (which is a left shift)
	mpz.TdivQ2exp(z.ptr, x.ptr, int(n))
	return z
}

// TODO SCAN

func (z *Int) Set(x *Int) *Int {
	if z.ptr == nil {
		n := newIntPtr(0)
		z.ptr = n
		runtime.AddCleanup(z, mpz.Delete, n)
	}
	mpz.Set(z.ptr, x.ptr)
	return z
}

// TODO SETBIT
// TODO SETBITS

// SetBit sets z to x, with x's i'th bit set to b (0 or 1).
// That is, if b is 1 SetBit sets z = x | (1 << i);
// if b is 0 SetBit sets z = x &^ (1 << i). If b is not 0 or 1,
// SetBit will panic.
func (z *Int) SetBit(x *Int, i int, b uint) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if z != x {
		z.Set(x)
	}
	if b == 0 {
		mpz.Clrbit(z.ptr, uint(i))
	} else {
		mpz.Setbit(z.ptr, uint(i))
	}
	return z
}

// / SetBytes interprets buf as the bytes of a big-endian unsigned
// integer, sets z to that value, and returns z.
func (z *Int) SetBytes(buf []byte) *Int {
	if z.ptr == nil {
		if len(buf) == 0 {
			// nothing to do
			return z
		}
		z.init()
	}
	if len(buf) == 0 {
		z.SetInt64(0)
		return z
	}
	mpz.ImportGo(z.ptr, buf)
	return z
}

func (z *Int) SetInt64(x int64) *Int {
	if z.ptr == nil {
		n := newIntPtr(int(x))
		mpz.SetSi(n, int(x))
		z.ptr = n
		runtime.AddCleanup(z, mpz.Delete, n)
	} else {
		mpz.SetSi(z.ptr, int(x))
	}
	return z
}

// SetString sets z to the value of s, interpreted in the given base,
// and returns z and a boolean indicating success. If SetString fails,
// the value of z is undefined but the returned value is nil.
//
// The base argument must be 0 or a value from 2 through MaxBase. If the base
// is 0, the string prefix determines the actual conversion base. A prefix of
// “0x” or “0X” selects base 16; the “0” prefix selects base 8, and a
// “0b” or “0B” prefix selects base 2. Otherwise the selected base is 10.

func (z *Int) setStringBase0(s string) (*Int, bool) {
	if len(s) == 0 {
		return nil, false
	}

	base := 0

	hasUnaryPlus := false
	hasUnaryMinus := false
	hasUnderscore := false
	octal := false
	hex := false

	idx := 0
	switch s[idx] {
	case '+':
		hasUnaryPlus = true
		idx += 1
	case '-':
		hasUnaryMinus = true
		idx += 1
	}

	// just a "+" or "-"
	if idx == len(s) {
		return nil, false
	}

	// next char must be a number
	ch := s[idx]
	switch {
	case '0' <= ch && ch <= '9':
	default:
		return nil, false
	}

	if s[idx] == '0' {
		idx += 1
		if idx == len(s) {
			z.SetInt64(0)
			return z, true
		}
		ch := s[idx]
		switch {
		case ch == 'x' || ch == 'X':
			idx += 1
			hex = true
		case ch == 'o' || ch == 'O':
			octal = true
			idx += 1
		case ch == 'b' || ch == 'B':
			idx += 1
		case '0' <= ch && ch <= '9':
			// leading zero followed by a number
		case ch == '_':
			hasUnderscore = true
		default:
			return nil, false
		}
		if idx == len(s) {
			return nil, false
		}
	}

	// validate all remaining characters
	for ; idx < len(s); idx += 1 {
		ch := s[idx]
		switch {
		case ch == '_':
			hasUnderscore = true
		case '0' <= ch && ch <= '9':
		case 'a' <= ch && ch <= 'f':
			if !hex {
				return nil, false
			}
		case 'A' <= ch && ch <= 'F':
			if !hex {
				return nil, false
			}
		default:
			return nil, false
		}
	}
	if hasUnderscore {
		// double underscore not allowed
		if strings.Contains(s, "__") {
			return nil, false
		}
		// trailing underscore not allowed
		if s[len(s)-1] == '_' {
			return nil, false
		}
		// remove the underscores
		s = strings.ReplaceAll(s, "_", "")
	}

	if hasUnaryPlus || hasUnaryMinus {
		s = s[1:]
	}

	// GMP doesn't support 0o for octal
	if octal {
		base = 8
		s = s[2:]
	}

	if z.ptr == nil {
		z.init()
	}
	if !mpz.SetStr(z.ptr, s, base) {
		return nil, false
	}
	if hasUnaryMinus {
		z.Neg(z)
	}
	return z, true
}
func (z *Int) SetString(s string, base int) (*Int, bool) {
	if base == 0 {
		return z.setStringBase0(s)
	}
	if base < 2 || base > 62 {
		return nil, false
	}
	if len(s) == 0 {
		return nil, false
	}

	hasUnaryPlus := false
	//hasUnaryMinus := false
	//leadingZero := false

	idx := 0
	b := s[idx]
	switch {
	case b == '+':
		hasUnaryPlus = true
		idx += 1
	case b == '-':
		//hasUnaryMinus = true
		idx += 1
	}

	// just a "+" or "-"
	if idx == len(s) {
		return nil, false
	}

	if s[idx] == '0' {
		//	hasLeadingZero = true
	}

	if base <= 36 {
		for ; idx < len(s); idx += 1 {
			ch := s[idx]
			switch {
			case '0' <= ch && ch <= '9':
			case 'a' <= ch && ch <= 'z':
			case 'A' <= ch && ch <= 'Z':
			default:
				return nil, false
			}
		}
	} else {
		// validate and swap case.
		bstr := []byte(s)
		for ; idx < len(s); idx += 1 {
			ch := bstr[idx]
			switch {
			case 'a' <= ch && ch <= 'z':
				bstr[idx] = ch - 'a' + 'A'
			case 'A' <= ch && ch <= 'Z':
				bstr[idx] = ch - 'A' + 'a'
			case '0' <= ch && ch <= '9':
				// NOP
			default:
				return nil, false
			}
			s = string(bstr)
		}
	}
	if hasUnaryPlus {
		s = s[1:]
	}

	if z.ptr == nil {
		z.init()
	}
	if !mpz.SetStr(z.ptr, s, base) {
		return nil, false
	}
	return z, true

}

func (z *Int) SetUint64(x uint64) *Int {
	if z.ptr == nil {
		n := newIntPtr(0)
		mpz.SetUi(n, uint(x))
		z.ptr = n
		runtime.AddCleanup(z, mpz.Delete, n)
	} else {
		mpz.SetUi(z.ptr, uint(x))
	}
	return z
}

func (z *Int) Sign() int {
	if z.ptr == nil {
		return 0
	}
	return mpz.Sgn(z.ptr)
}

func (z *Int) Sqrt(x *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	mpz.Sqrt(z.ptr, x.ptr)
	return z
}

func (z *Int) String() string {
	if z.ptr == nil {
		return ""
	}
	return mpz.GetStr(10, z.ptr)
}

func (z *Int) Sub(x, y *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpz.Sub(z.ptr, x.ptr, y.ptr)
	return z
}

// TODO TEXT
// TODO TRAILINGZEROBITS

func (z *Int) Uint64() uint64 {
	if z.ptr == nil {
		return 0
	}
	return uint64(mpz.GetUi(z.ptr))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (z *Int) UnmarshalJSON(x []byte) error {
	_, ok := z.SetString(string(x), 0)
	if !ok {
		return fmt.Errorf("math/big: cannot unmarshal %s into a *gmp.Int", x)
	}
	return nil
}

// TODO UNMARSHALTEXT

func (z *Int) Xor(x, y *Int) *Int {
	if z.ptr == nil {
		z.init()
	}
	if x.ptr == nil {
		x.init()
	}
	if y.ptr == nil {
		y.init()
	}
	mpz.Xor(z.ptr, x.ptr, y.ptr)
	return z
}
