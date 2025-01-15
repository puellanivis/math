package floats

import (
	"fmt"
	"math"
	"math/big"

	"github.com/puellanivis/math/bits"
)

var _ = fmt.Println

type datum = bits.Uint

func set[D2 datum, D1 datum](dst *D2, src D1) {
	switch v := any(dst).(type) {
	case *bits.Uint128:
		switch src := any(src).(type) {
		case bits.Uint128:
			*v = src
		case uint64:
			v.Lo = uint64(src)
		case uint32:
			v.Lo = uint64(src)
		case uint16:
			v.Lo = uint64(src)
		}

	case *uint64:
		switch src := any(src).(type) {
		case bits.Uint128:
			*v = uint64(src.Lo)
		case uint64:
			*v = uint64(src)
		case uint32:
			*v = uint64(src)
		case uint16:
			*v = uint64(src)
		}

	case *uint32:
		switch src := any(src).(type) {
		case bits.Uint128:
			*v = uint32(src.Lo)
		case uint64:
			*v = uint32(src)
		case uint32:
			*v = uint32(src)
		case uint16:
			*v = uint32(src)
		}

	case *uint16:
		switch src := any(src).(type) {
		case bits.Uint128:
			*v = uint16(src.Lo)
		case uint64:
			*v = uint16(src)
		case uint32:
			*v = uint16(src)
		case uint16:
			*v = uint16(src)
		}
	}
}

type spec[D datum] interface {
	width() int
	expWidth() int
	mantWidth() int

	exp2OverUnder() (over, under D)
	expOverUnder() (over, under, nearZero D)
	ln2HiLoE() (hi, lo, e D)
	expPN() []D

	bits.Bits[D]
}

type binary[SPEC spec[D], D datum] struct {
	s bool
	e int
	m D
	r D
}

func fromBigFloat[SPEC spec[D], D datum](v *big.Float, rounding RoundingMode) D {
	if v.IsInf() {
		return inf[SPEC](v.Signbit())
	}

	var spec SPEC
	var z D

	mant := new(big.Float).SetPrec(uint(spec.width()))
	exp := v.MantExp(mant)

	tmp := new(big.Float).SetPrec(uint(spec.width()))

	mant.Mul(mant, tmp.SetFloat64(2))
	exp = exp + expBias[SPEC]() - 1

	g := binary[SPEC, D]{
		s: v.Signbit(),
		e: exp,
	}

	switch any(z).(type) {
	case bits.Uint128:
		var t D

		h := new(big.Float).Mul(mant, tmp.SetFloat64(math.Ldexp(1.0, 63)))
		hi, _ := h.Uint64()
		set(&t, hi)

		g.m = spec.Shl(t, 64)

		l := new(big.Float).Sub(h, tmp.SetUint64(hi))

		l.Mul(l, tmp.SetFloat64(math.Ldexp(1.0, 64)))
		lo, _ := l.Uint64()
		set(&t, lo)

		g.m = spec.Or(g.m, t)

	default:
		h := new(big.Float).Mul(mant, tmp.SetFloat64(math.Ldexp(1.0, spec.width()-1)))
		hi, _ := h.Uint64()
		set(&g.m, hi)
	}

	switch {
	case g.e >= expMax[SPEC]():
		// EXCEPTION: overflow
		return overflow[SPEC](g.s, rounding)

	case g.e <= 1:
		// forced into sub-norm
		g.shr(1 - g.e)
		g.e = 1
	}

	if spec.IsZero(g.m) {
		// EXCEPTION: underflow
		return z
	}

	g.renorm()

	applyRounding(&g, rounding)

	return g.encode()
}

func signMask[SPEC spec[D], D datum]() D {
	var spec SPEC
	return spec.Pow2(spec.width() - 1)
}

func magMask[SPEC spec[D], D datum]() D {
	var spec SPEC
	return spec.Pow2m1(spec.width() - 1)
}

func expMask[SPEC spec[D], D datum]() D {
	var spec SPEC
	return spec.Shl(spec.Pow2m1(spec.expWidth()), spec.mantWidth())
}

func mantMask[SPEC spec[D], D datum]() D {
	var spec SPEC
	return spec.Pow2m1(spec.mantWidth())
}

func quietMask[SPEC spec[D], D datum]() D {
	var spec SPEC
	return spec.Pow2(spec.mantWidth() - 1)
}

func expMax[SPEC spec[D], D datum]() int {
	var spec SPEC
	return (1 << spec.expWidth()) - 1
}

func expBias[SPEC spec[D], D datum]() int {
	var spec SPEC
	return (1 << (spec.expWidth() - 1)) - 1
}

func half[SPEC spec[D], D datum]() D {
	var spec SPEC
	return spec.Shl(spec.FromInt(expBias[SPEC]()-1), spec.mantWidth())
}

func one[SPEC spec[D], D datum]() D {
	var spec SPEC
	return spec.Shl(spec.FromInt(expBias[SPEC]()), spec.mantWidth())
}

func two[SPEC spec[D], D datum]() D {
	var spec SPEC
	return spec.Shr(signMask[SPEC](), 1)
}

func inf[SPEC spec[D], D datum](sign bool) D {
	var spec SPEC

	bits := magInf[SPEC]()
	if sign {
		return spec.Or(signMask[SPEC](), bits)
	}
	return bits
}

func nan[SPEC spec[D], D datum]() D {
	var spec SPEC

	return spec.Or(expMask[SPEC](), quietMask[SPEC]())
}

func mag[SPEC spec[D], D datum](bits D) (sign, mag D) {
	var spec SPEC

	sign = spec.And(bits, signMask[SPEC]())
	mag = spec.And(bits, magMask[SPEC]())
	return
}

func magInf[SPEC spec[D], D datum]() D {
	return expMask[SPEC]()
}

// decomp returns the decomposed sign, exponent, and mantissa of the encoded floating point.
// It does not align either the sign or mantissa, and leaves them simply in place.
//
// For reasons of utility, sign is _not_ right aligned into the LSB position.
// Turns out, one needs the plain masked sign bit in-place more often than bit-shifted.
// Plus, most of the time, we internally covert it to a bool anyways.
//
// The exponent and mantissa are aligned to the LSB position.
func decomp[SPEC spec[D], D datum](bits D) (sign D, exp int, mant D) {
	var spec SPEC

	sign = spec.And(bits, signMask[SPEC]())
	exp = spec.Int(spec.Shr(spec.And(bits, expMask[SPEC]()), spec.mantWidth()))
	mant = spec.And(bits, mantMask[SPEC]())
	return
}

func decode[SPEC spec[D], D datum](bits D) binary[SPEC, D] {
	var spec SPEC

	sign, exp, mant := decomp[SPEC](bits)

	f := binary[SPEC, D]{
		s: !spec.IsZero(sign),
		e: int(exp),
		m: spec.Shl(mant, spec.expWidth()),
	}

	switch f.e {
	case 0: // sub-normal: no implied leading bit,
		f.e = 1 // and exponent is one larger than.
	case expMax[SPEC](): // NaN and Inf
	default:
		f.m = spec.Or(f.m, signMask[SPEC]()) // add implied leading bit
	}

	return f
}

func (f *binary[SPEC, D]) encode() D {
	var spec SPEC

	bits := spec.Shr(f.m, spec.expWidth())

	if f.e != 1 {
		bits = spec.And(bits, mantMask[SPEC]()) // clear the top bit of mantissa as it is implied
		bits = spec.Or(bits, spec.Shl(spec.FromInt(f.e), spec.mantWidth()))
	}

	// if f.e == 1, and the top-bit of the mantissa is set,
	// then we have actually un-denormalized, and the encoded exponent needs to be 1,
	// which is right where it’s already at.

	if f.s {
		bits = spec.Or(bits, signMask[SPEC]())
	}

	return bits
}

func convert[SPEC1 spec[D1], SPEC2 spec[D2], D1, D2 datum](bits D1, rounding RoundingMode) D2 {
	var spec1 SPEC1
	var spec2 SPEC2

	Δw := spec2.width() - spec1.width()
	if Δw == 0 {
		var x D2
		set(&x, bits)
		return x
	}

	f := decode[SPEC1](bits)

	fInf, fNaN := f.classify()
	if fInf {
		return inf[SPEC2](f.s)
	}

	g := binary[SPEC2, D2]{
		s: f.s,

		// we don’t need to worry about overflow here.
		// Inf and NaN cases have already been handled,
		// so nothing else should be even close to overfloating MaxInt.
		e: f.ilogb() + expBias[SPEC2](),
	}

	if Δw < 0 {
		// we’re scaling down, cast after the shift or we will clip out the part we need.
		set(&g.r, spec1.Shr(f.m, -Δw-spec2.width()))
		set(&g.m, spec1.Shr(f.m, -Δw))
	} else {
		// we’re scaling up, cast before the shift or we will shift the whole thing into the bitbucket.
		set(&g.m, f.m)
		g.m = spec2.Shl(g.m, Δw)
	}

	if fNaN {
		// converting up to a larger data-type, and back down needs to preserve the NaN payload.
		bits := spec2.Or(expMask[SPEC2](), spec2.Shr(g.m, spec2.expWidth()))

		if f.s {
			return spec2.Or(bits, signMask[SPEC2]())
		}

		return bits
	}

	switch {
	case g.e >= expMax[SPEC2]():
		// EXCEPTION: overflow
		return overflow[SPEC2](g.s, rounding)

	case g.e <= 1:
		// forced into sub-norm
		g.shr(1 - g.e)
		g.e = 1
	}

	var spec SPEC2

	if spec.IsZero(g.m) {
		// EXCEPTION: underflow
	}

	g.renorm()

	applyRounding(&g, rounding)

	return g.encode()
}

func (f *binary[SPEC, D]) trunc() {
	var spec SPEC
	var z D

	f.m = spec.Mask(f.m, incAway[SPEC]())
	f.r = z
}

func (f *binary[SPEC, D]) shr(shift int) {
	var spec SPEC

	f.r = spec.Shr(f.r, shift)
	if lshift := spec.width() - shift; lshift >= 0 {
		f.r = spec.Or(f.r, spec.Shl(f.m, lshift))
	}
	f.m = spec.Shr(f.m, shift)
}

func (f *binary[SPEC, D]) shl(shift int) {
	var spec SPEC

	f.m = spec.Shl(f.m, shift)
	if rshift := spec.width() - shift; rshift >= 0 {
		f.m = spec.Or(f.m, spec.Shr(f.r, rshift))
	}
	f.r = spec.Shl(f.r, shift)
}

func (f *binary[SPEC, D]) add(inc, rem D) {
	var spec SPEC
	var z D

	var carry D
	f.r, carry = spec.Add(f.r, rem, z)
	f.m, carry = spec.Add(f.m, inc, carry)
	f.e += spec.Int(carry)

	if f.e == expMax[SPEC]() {
		f.m = z
		return
	}

	f.m = spec.Or(f.m, carry)
	f.m = spec.Rotl(f.m, -spec.Int(carry))
}

func (f *binary[SPEC, D]) sub(dec, rem D) {
	var spec SPEC
	var z D

	var borrow D
	f.r, borrow = spec.Sub(f.r, rem, z)
	f.m, _ = spec.Sub(f.m, dec, borrow)
	f.renorm()

	if spec.IsZero(f.m) {
		f.s = false
	}
}

func (f *binary[SPEC, D]) mulPrim(g *binary[SPEC, D]) {
	var spec SPEC
	var z D

	f.s = f.s != g.s
	f.e += g.e - expBias[SPEC]()
	f.e++

	if f.e < 1 {
		// underflow
		f.e = 1
		f.m = z
		return
	}

	f.m, f.r = spec.Mul(f.m, g.m)

	f.renorm()

	if f.e >= expMax[SPEC]() {
		f.e = expMax[SPEC]()
		f.m = z
	}
}

func (f *binary[SPEC, D]) divPrim(g *binary[SPEC, D]) {
	var spec SPEC
	var z D

	f.s = f.s != g.s
	f.e -= g.e - expBias[SPEC]()
	f.e--

	for spec.Gte(f.m, g.m) {
		f.e++
		f.shr(1)
	}

	if f.e < 1 {
		// underflow
		f.e = 1
		f.m = z
		return
	}

	f.m, f.r = spec.Div(f.m, f.r, g.m)

	f.renorm()

	if f.e >= expMax[SPEC]() {
		f.e = expMax[SPEC]()
		f.m = z
	}
}

func (f *binary[SPEC, D]) renorm() {
	var spec SPEC

	lz := spec.Lzcnt(f.m)
	if lz == 0 {
		return
	}

	if lz == spec.width() {
		f.e = 1
		return
	}

	lz = min(lz, f.e-1)

	f.shl(lz)
	f.e -= lz
}

func (f binary[SPEC, D]) classify() (inf, nan bool) {
	if f.e != expMax[SPEC]() {
		return false, false
	}

	var spec SPEC

	isInf := spec.IsZero(f.m)

	return isInf, !isInf
}

func (f binary[SPEC, D]) isNaN() bool {
	var spec SPEC

	return f.e == expMax[SPEC]() && !spec.IsZero(f.m)
}

func (f binary[SPEC, D]) isInf() bool {
	var spec SPEC
	var z D

	return f.e == expMax[SPEC]() && spec.IsZero(f.m)
}

func (f binary[SPEC, D]) isZero() bool {
	var spec SPEC

	return f.e == 1 && spec.IsZero(f.m)
}

func (f *binary[SPEC, D]) ilogb() int {
	if f.isZero() {
		return math.MinInt32
	}
	if f.e == expMax[SPEC]() {
		return math.MaxInt32
	}
	return f.e - expBias[SPEC]()
}

func format[SPEC spec[D], D datum](x D, f fmt.State, verb rune, rounding RoundingMode) {

	switch verb {
	case 'b':
		var spec SPEC

		s, e, m := decomp[SPEC](x)
		if !spec.IsZero(s) {
			fmt.Fprintf(f, "1_%0*b_%0*b", spec.expWidth(), e, spec.mantWidth(), m)
		} else {
			fmt.Fprintf(f, "0_%0*b_%0*b", spec.expWidth(), e, spec.mantWidth(), m)
		}

	default:
		switch x := any(x).(type) {
		case bits.Uint128:
			xf := decode[binary128](x)
			if xf.isZero() {
				fmt.Fprintf(f, fmt.FormatString(f, verb), 0.0)
				return
			}

			one := 1.0
			if xf.s {
				one = -one
			}

			g := new(big.Float).SetPrec(128).SetUint64(xf.m.Hi)
			g.Mul(g, new(big.Float).SetFloat64(0x1p64))
			g.Add(g, new(big.Float).SetUint64(xf.m.Lo))
			g.Mul(g, new(big.Float).SetMantExp(new(big.Float).SetFloat64(one), xf.ilogb()-127))

			fmt.Fprintf(f, fmt.FormatString(f, verb), g)
			return
		}

		g := math.Float64frombits(convert[SPEC, binary64](D(x), rounding))
		fmt.Fprintf(f, fmt.FormatString(f, verb), g)
	}
}

func isNaN[SPEC spec[D], D datum](x D) bool {
	var spec SPEC

	_, m := mag[SPEC](x)
	return spec.Gt(m, magInf[SPEC]())
}

func isInf[SPEC spec[D], D datum](x D) bool {
	var spec SPEC

	_, m := mag[SPEC](x)
	return spec.Eq(m, magInf[SPEC]())
}

func signBit[SPEC spec[D], D datum](x D) bool {
	var spec SPEC

	s, _ := mag[SPEC](x)

	return !spec.IsZero(s)
}

func getSign[SPEC spec[D], D datum](x D) int {
	var spec SPEC

	s, m := mag[SPEC](x)

	if spec.IsZero(m) {
		return 0
	}

	if !spec.IsZero(s) {
		return -1
	}

	return 1
}

func abs[SPEC spec[D], D datum](x D) D {
	_, m := mag[SPEC](x)
	return m
}

func neg[SPEC spec[D], D datum](x D) D {
	var spec SPEC

	return spec.Xor(x, signMask[SPEC]())
}

func copySign[SPEC spec[D], D datum](x, y D) D {
	var spec SPEC

	_, m := mag[SPEC](x)
	s, _ := mag[SPEC](y)
	return spec.Or(s, m)
}

func nextUp[SPEC spec[D], D datum](x D) D {
	var spec SPEC
	var z D

	s, m := mag[SPEC](x)

	if spec.Gte(m, magInf[SPEC]()) {
		if spec.IsZero(s) || spec.Gt(m, magInf[SPEC]()) {
			// +∞ and NaN do not change
			return x
		}
	}

	if spec.IsZero(m) {
		// ±zero is smallest positive sub-normal
		return spec.Inc(z)
	}

	if !spec.IsZero(s) {
		return spec.Dec(x)
	}

	return spec.Inc(x)
}

func nextDown[SPEC spec[D], D datum](x D) D {
	var spec SPEC
	var z D

	s, m := mag[SPEC](x)

	if spec.Gte(m, magInf[SPEC]()) {
		if !spec.IsZero(s) || spec.Gt(m, magInf[SPEC]()) {
			// -∞ and NaN do not change
			return x
		}
	}

	if spec.IsZero(m) {
		// ±zero is smallest negative sub-normal
		return spec.Or(signMask[SPEC](), spec.Inc(z))
	}

	if !spec.IsZero(s) {
		return spec.Inc(x)
	}

	return spec.Dec(x)
}

func add[SPEC spec[D], D datum](x, y D, rounding RoundingMode) D {
	f, g := decode[SPEC](x), decode[SPEC](y)

	fInf, fNaN := f.classify()
	gInf, gNaN := g.classify()

	switch {
	case fNaN:
		return x
	case gNaN:
		return y

	case fInf && gInf:
		if f.s == g.s {
			return x
		}

		// EXCEPTION: illegal operation: adding: ±∞ + ∓∞
		return nan[SPEC]()

	case fInf:
		return x
	case gInf:
		return y
	}

	if lessMagPrim[SPEC](x, y) {
		f, g = g, f
	}

	g.shr(f.e - g.e)

	if f.s != g.s {
		f.sub(g.m, g.r)
	} else {
		f.add(g.m, g.r)
	}

	applyRounding(&f, rounding)

	return f.encode()
}

func sub[SPEC spec[D], D datum](x, y D, rounding RoundingMode) D {
	return add[SPEC](x, neg[SPEC](y), rounding)
}

func dim[SPEC spec[D], D datum](x, y D, rounding RoundingMode) D {
	diff := add[SPEC](x, neg[SPEC](y), rounding)

	s, m := mag[SPEC](diff)

	var spec SPEC
	var z D

	if spec.Gt(m, magInf[SPEC]()) {
		return diff // NaN → NaN
	}

	if !spec.IsZero(s) {
		// diff < 0, return zero.
		return z
	}

	return diff
}

// compare is not standards compliant, but is written this way to be compliant with the Go standard cmp library.
//
// Soruce:
//
//	“Every NaN shall compare unordered with everything, including itself.”
//	— IEEE Std 754-2008, § 5.11 pg. 29
func compare[SPEC spec[D], D datum](x, y D) int {
	order, ordered := fcmp[SPEC](x, y)
	if ordered {
		return order
	}

	xNaN := isNaN[SPEC](x)
	yNaN := isNaN[SPEC](y)

	switch {
	case xNaN && yNaN:
		return 0
	case xNaN:
		return -1
	default:
		return 1
	}
}

// As less depends upon compare, which is not standards compliant, this call is also not standards compliant.
func less[SPEC spec[D], D datum](x, y D) bool {
	return compare[SPEC](x, y) < 0
}

func fcmp[SPEC spec[D], D datum](x, y D) (order int, ordered bool) {
	xs, xm := mag[SPEC](x)
	ys, ym := mag[SPEC](y)

	var spec SPEC

	if spec.Gt(xm, magInf[SPEC]()) {
		return 0, false
	}

	if spec.Gt(ym, magInf[SPEC]()) {
		return 0, false
	}

	if spec.Neq(xs, ys) {
		// signs differ

		if spec.IsZero(spec.Or(xm, ym)) {
			// special case: ±zero == ∓zero
			return 0, true
		}

		if !spec.IsZero(xs) {
			// left is negative, and right is opposite
			return -1, true
		}

		return 1, true
	}

	return spec.Cmp(xm, ym), true
}

func fmin[SPEC spec[D], D datum](x, y D) D {
	order, ordered := fcmp[SPEC](x, y)
	if !ordered {
		// The numbers are unordered, so, at least one of x or y must be a NaN.

		if isNaN[SPEC](y) {
			// If y is NaN, then return x; whatever it is.
			// If x is also NaN, then we are preferencially returning x over y.
			return x
		}

		// Otherwise, x must be NaN,
		// So, we return y instead.
		return y
	}

	if order > 0 {
		// If x > y, then return y.
		return y
	}

	if order == 0 {
		var spec SPEC

		// The only case where the bits of x and y are different here,
		// is if either x or y is -0, and we want to return -0 if either is.
		return spec.Or(x, y)
	}

	// Preferentially, return x over y.
	return x
}

func fmax[SPEC spec[D], D datum](x, y D) D {
	order, ordered := fcmp[SPEC](x, y)
	if !ordered {
		// The numbers are unordered, so, at least one of x or y must be a NaN.

		if isNaN[SPEC](y) {
			// If y is NaN, then return x; whatever it is.
			// If x is also NaN, then we are preferencially returning x over y.
			return x
		}

		// Otherwise, x must be NaN,
		// So, we return y instead.
		return y
	}

	if order < 0 {
		// If x < y, then return y.
		return y
	}

	if order == 0 {
		var spec SPEC

		// The only case where the bits of x and y are different here,
		// is if either x or y is -0, and we want to return 0 if only one is.
		return spec.And(x, y)
	}

	// Preferentially, return x over y.
	return x
}

func lessMagPrim[SPEC spec[D], D datum](x, y D) bool {
	_, xm := mag[SPEC](x)
	_, ym := mag[SPEC](y)

	var spec SPEC

	return spec.Lt(xm, ym)
}

func fcmpMag[SPEC spec[D], D datum](x, y D) (order int, ordered bool) {
	_, xm := mag[SPEC](x)
	_, ym := mag[SPEC](y)

	var spec SPEC

	if spec.Gt(xm, magInf[SPEC]()) {
		return 0, false
	}

	if spec.Gt(ym, magInf[SPEC]()) {
		return 0, false
	}

	return spec.Cmp(xm, ym), true
}

func fminMag[SPEC spec[D], D datum](x, y D) D {
	order, ordered := fcmpMag[SPEC](x, y)
	if !ordered {
		// The numbers are unordered, so, at least one of x or y must be a NaN.

		if isNaN[SPEC](y) {
			// If y is NaN, then return x; whatever it is.
			// If x is also NaN, then we are preferencially returning x over y.
			return x
		}

		// Otherwise, x must be NaN,
		// So, we return y instead.
		return y
	}

	if order > 0 {
		// If x > y, then return y.
		return y
	}

	// Preferentially, return x over y.
	return x
}

func fmaxMag[SPEC spec[D], D datum](x, y D) D {
	order, ordered := fcmpMag[SPEC](x, y)
	if !ordered {
		// The numbers are unordered, so, at least one of x or y must be a NaN.

		if isNaN[SPEC](y) {
			// If y is NaN, then return x; whatever it is.
			// If x is also NaN, then we are preferencially returning x over y.
			return x
		}

		// Otherwise, x must be NaN,
		// So, we return y instead.
		return y
	}

	if order < 0 {
		// If x < y, then return y.
		return y
	}

	// Preferentially, return x over y.
	return x
}

// normalize returns a normal number y and exponent exp
// satisfying mag == y × 2**exp. It assumes x is positive, finite, and non-zero.
func normalize[SPEC spec[D], D datum](x D) (y D, exp int) {
	s, e, m := decomp[SPEC](x)

	if e != 0 {
		return x, 0
	}

	var spec SPEC

	// offset necessary to put the top bit of the mantissa into the exponent.
	exp = (spec.Lzcnt(m) - expBias[SPEC]()) + 1

	// There is no need to mask the implicit top bit out, as it fills in the exponent field for us.
	return spec.Or(s, spec.Shl(m, exp)), exp
}

func frexp[SPEC spec[D], D datum](x D) (frac D, exp int) {
	_, m := mag[SPEC](x)

	var spec SPEC

	switch {
	case spec.IsZero(m):
		return x, 0 // correctly return -0
	case spec.Gte(m, magInf[SPEC]()):
		return x, 0
	}

	_, exp, _ = decomp[SPEC](x)

	x, e := normalize[SPEC](x)
	exp += e - expBias[SPEC]() + 1
	x = spec.MaskInsert(x, half[SPEC](), expMask[SPEC]())
	return x, exp
}

func ldexp[SPEC spec[D], D datum](frac D, exp int) D {
	s, m := mag[SPEC](frac)

	var spec SPEC
	var z D

	switch {
	case spec.IsZero(m):
		return frac // ±0
	case spec.Gte(m, magInf[SPEC]()):
		return frac // ±∞ and NaN
	}

	frac, e := normalize[SPEC](frac)
	exp += e

	_, e, _ = decomp[SPEC](frac)
	exp += e - expBias[SPEC]()

	if exp < -(expBias[SPEC]() + spec.mantWidth()) {
		// underflow
		return s // ±0
	}

	if exp >= expBias[SPEC]() {
		// overflow
		return spec.Or(s, magInf[SPEC]()) // ±∞
	}

	if exp < -(expBias[SPEC]() - 1) {
		// subnormal
		panic("subnormals not handled")
		return z // TODO correct this
	}

	ne := spec.Shl(spec.FromInt(exp+expBias[SPEC]()), spec.mantWidth())
	return spec.MaskInsert(frac, ne, expMask[SPEC]())
}

func modf[SPEC spec[D], D datum](x D) (i, f D) {
	i = trunc[SPEC](x)
	f = sub[SPEC](x, i, RoundTowardZero{})
	return
}

func round[SPEC spec[D], D datum](x D) D {
	s, e, _ := decomp[SPEC](x)

	var spec SPEC
	var z D

	e -= expBias[SPEC]()
	if e < 0 {
		// Round abs(x) < 1 including denormals.
		// start with s, which is ±0

		if e == -1 {
			// e == bias-1 is one half, all of these numbers round away
			return spec.Or(s, one[SPEC]()) // ±1
		}

		return s // ±0
	}

	if e >= spec.mantWidth() {
		// This exponent is large enough to not have any fractional part.
		return x
	}

	// Round any abs(x) >= 1 containing a fractional component [0,1)

	x, _ = spec.Add(x, spec.Shr(quietMask[SPEC](), e), z) // add one half, scaled into our fractional range.
	return spec.Mask(x, spec.Shr(mantMask[SPEC](), e))
}

func roundToEven[SPEC spec[D], D datum](x D) D {
	s, e, m := decomp[SPEC](x)

	var spec SPEC
	var z D

	e -= expBias[SPEC]()
	if e < 0 {
		// Round abs(x) < 1 including denormals.
		if e == -1 && !spec.IsZero(m) {
			// e == bias-1 and m == 0 is one half, which rounds down to zero (even).
			// Every other mantissa is greater than one-half, so round away from zero.
			return spec.Or(s, one[SPEC]()) // ±1
		}

		return s // ±0
	}

	if e >= spec.mantWidth() {
		// This exponent is large enough to not have any fractional part.
		return x
	}

	// Round abs(x) >= 1
	// - Add 0.499.. or 0.5 before truncating,
	//   depending on whether the truncated number is even or odd (respectively).

	// 1’s bit into lowest bit
	even := spec.And(spec.Shr(x, spec.mantWidth()-e), spec.Pow2(0))

	// 4.999… + even
	halfm1, _ := spec.Add(spec.Dec(quietMask[SPEC]()), even, z)

	x, _ = spec.Add(x, spec.Shr(halfm1, e), z)

	// truncate whatever fraction is left.
	return spec.Mask(x, spec.Shr(mantMask[SPEC](), e))
}

func floor[SPEC spec[D], D datum](x D) D {
	s, m := mag[SPEC](x)

	var spec SPEC

	if spec.IsZero(m) {
		// zeros return themselves.
		return x
	}

	if spec.Gte(m, magInf[SPEC]()) {
		// NaN and Inf return themselves.
		return x
	}

	if spec.IsZero(s) {
		// positive numbers round toward zero.
		return trunc[SPEC](x)
	}

	d, fract := modf[SPEC](spec.Xor(x, s))
	if !spec.IsZero(fract) {
		d = add[SPEC](d, one[SPEC](), RoundTowardZero{})
	}
	return spec.Or(d, s)
}

func truncToInt[SPEC spec[D], D datum](x D) int {
	return int(math.Float64frombits(convert[SPEC, binary64](x, RoundTowardZero{})))
}

func trunc[SPEC spec[D], D datum](x D) D {
	s, e, _ := decomp[SPEC](x)

	if e == expMax[SPEC]() {
		// NaN and Inf return themselves.
		return x
	}

	var spec SPEC

	e -= expBias[SPEC]()
	if e < 0 {
		// abs(x) < 1 all round down to ±0
		return s
	}

	shift := spec.mantWidth()

	if e >= shift {
		// This exponent is large enough to not have any fractional part.
		return x
	}

	// Drop the lowest mantissa-e bits, the fractional part, leaving only the integer part.
	return spec.Mask(x, spec.Pow2m1(shift-e))
}

func ceil[SPEC spec[D], D datum](x D) D {
	return neg[SPEC](floor[SPEC](neg[SPEC](x)))
}

// msub returns x * y - z, computed with only one rounding.
func msub[SPEC spec[D], D datum](x, y, z D, rounding RoundingMode) D {
	return madd[SPEC](x, y, neg[SPEC](z), rounding)
}

// mnsub returns -(x * y) + z, computed with only one rounding.
func mnsub[SPEC spec[D], D datum](x, y, z D, rounding RoundingMode) D {
	return neg[SPEC](madd[SPEC](x, y, neg[SPEC](z), rounding))
}

// madd returns x * y + z, computed with only one rounding.
func madd[SPEC spec[D], D datum](x, y, z D, rounding RoundingMode) D {
	f, g, h := decode[SPEC](x), decode[SPEC](y), decode[SPEC](z)

	fInf, fNaN := f.classify()
	gInf, gNaN := g.classify()
	hInf, hNaN := h.classify()

	switch {
	case fNaN:
		return x
	case gNaN:
		return y
	case hNaN:
		return z

	case fInf:
		if g.isZero() {
			// EXCEPTION: illegal operation: adding: ±∞ × ±0
			return nan[SPEC]()
		}

		newSign := f.s != g.s // sign of product x*y

		if newSign != h.s {
			// EXCEPTION: illegal operation: adding: ±∞ + ∓∞
			return nan[SPEC]()
		}

		// assuming fInf
		return inf[SPEC](newSign)

	case gInf:
		if f.isZero() {
			// EXCEPTION: illegal operation: adding: ±0 × ±∞
			return nan[SPEC]()
		}

		newSign := f.s != g.s // sign of product x*y

		if newSign != h.s {
			// EXCEPTION: illegal operation: adding: ±∞ + ∓∞
			return nan[SPEC]()
		}

		return inf[SPEC](newSign)
	}

	f.mulPrim(&g)

	if f.e == expMax[SPEC]() {
		// arithmetic overflow to ±∞

		if !hInf {
			// EXCEPTION: overflow
			return overflow[SPEC](f.s, rounding)
		}

		// recheck for sign mismatch infinite addition
		if f.s != h.s {
			// EXCEPTION: illegal operation: adding: ±∞ + ∓∞
			return nan[SPEC]()
		}

		return z
	}

	if hInf {
		return z
	}

	var spec SPEC

	if f.e < h.e || (f.e == h.e && spec.Lt(f.m, h.m)) {
		f, h = h, f
	}

	h.shr(f.e - h.e)

	if f.s != h.s {
		f.sub(h.m, h.r)
	} else {
		f.add(h.m, h.r)
	}

	applyRounding(&f, rounding)

	return f.encode()
}

func mul[SPEC spec[D], D datum](x, y D, rounding RoundingMode) D {
	f, g := decode[SPEC](x), decode[SPEC](y)

	fInf, fNaN := f.classify()
	gInf, gNaN := g.classify()

	switch {
	case fNaN:
		return x
	case gNaN:
		return y

	case fInf:
		if g.isZero() {
			// EXCEPTION: illegal operation: adding: ±∞ × ±0
			return nan[SPEC]()
		}

		return inf[SPEC](f.s != g.s)

	case gInf:
		if f.isZero() {
			// EXCEPTION: illegal operation: adding: ±0 × ±∞
			return nan[SPEC]()
		}

		return inf[SPEC](f.s != g.s)
	}

	f.mulPrim(&g)

	applyRounding(&f, rounding)

	return f.encode()
}

func rcp[SPEC spec[D], D datum](x D, rounding RoundingMode) D {
	return div[SPEC](one[SPEC](), x, rounding)
}

func div[SPEC spec[D], D datum](x, y D, rounding RoundingMode) D {
	f, g := decode[SPEC](x), decode[SPEC](y)

	fInf, fNaN := f.classify()
	gInf, gNaN := g.classify()

	switch {
	case fNaN:
		return x
	case gNaN:
		return y

	case fInf:
		if gInf {
			// EXCEPTION: illegal operation: dividing: ±∞ ÷ ±∞
			return nan[SPEC]()
		}

		return inf[SPEC](f.s != g.s)

	case gInf:
		var z D

		if f.s != g.s {
			return signMask[SPEC]()
		}

		return z

	case g.isZero():
		if f.isZero() {
			// EXCEPTION: illegal operation: dividing: ±0 ÷ ±0
			return nan[SPEC]()
		}

		// EXCEPTION: divide by zero
		return inf[SPEC](f.s != g.s)
	}

	f.divPrim(&g)

	applyRounding(&f, rounding)

	return f.encode()
}

func mod[SPEC spec[D], D datum](x, y D, rounding RoundingMode) D {
	s, xm := mag[SPEC](x)
	_, ym := mag[SPEC](y)

	var spec SPEC

	switch {
	case spec.Gt(xm, magInf[SPEC]()):
		return x
	case spec.Gt(ym, magInf[SPEC]()):
		return y

	case spec.Eq(xm, magInf[SPEC]()):
		// EXCEPTION: invalid operation: mod(±∞, y)
		return nan[SPEC]()

	case spec.IsZero(ym):
		// EXCEPTION: invalid operation: mod(x, 0)
		return nan[SPEC]()

	case spec.Eq(ym, magInf[SPEC]()):
		return x // mod(x, ±Inf) → x
	}

	yfr, yexp := frexp[SPEC](ym)

	for spec.Gte(xm, ym) {
		rfr, rexp := frexp[SPEC](xm)
		if spec.Lt(rfr, yfr) {
			rexp--
		}
		xm = sub[SPEC](xm, ldexp[SPEC](y, rexp-yexp), rounding)
	}

	return spec.Or(s, xm)
}

func sqrt[SPEC spec[D], D datum](x D, rounding RoundingMode) D {
	sign, m := mag[SPEC](x)

	var spec SPEC
	var z D

	switch {
	case spec.IsZero(m) || spec.Gte(m, magInf[SPEC]()):
		if spec.Neq(m, magInf[SPEC]()) || spec.IsZero(sign) {
			// either: not Infinity, or positive sign
			return x
		}

		// √(-∞), fallthrough to exception
		fallthrough

	case !spec.IsZero(sign):
		// EXCEPTION: illegal operation: √(-f)
		return nan[SPEC]()
	}

	_, exp, _ := decomp[SPEC](x)

	if exp == 0 {
		// subnorm
		var e int
		x, e = normalize[SPEC](x)
		exp -= e
		exp++
	}

	shift := spec.mantWidth()

	exp -= expBias[SPEC]()
	x = spec.MaskInsert(x, spec.Pow2(shift), expMask[SPEC]())

	if exp&1 == 1 { // odd exp, double x to make it even
		x = spec.Shl(x, 1)
	}
	exp >>= 1 // exp = exp/2, exponent of square root
	// generate sqrt(x) bit by bit
	x = spec.Shl(x, 1)
	var q, s D                // q = sqrt(x)
	r := spec.Pow2(shift + 1) // r = moving bit from MSB to LSB
	for !spec.IsZero(r) {
		t, _ := spec.Add(s, r, z)
		if spec.Lte(t, x) {
			s, _ = spec.Add(t, r, z)
			x, _ = spec.Sub(x, t, z)
			q, _ = spec.Add(q, r, z)
		}
		x = spec.Shl(x, 1)
		r = spec.Shr(r, 1)
	}

	// final rounding
	if !spec.IsZero(x) {
		q, _ = spec.Add(q, spec.And(q, spec.Pow2(0)), z)
	}

	exp += expBias[SPEC]() - 1

	x, _ = spec.Add(spec.Shr(q, 1), spec.Shl(spec.FromInt(exp), shift), z)
	return x
}

func rsqrt[SPEC spec[D], D datum](x D, rounding RoundingMode) D {
	s, m := mag[SPEC](x)

	var spec SPEC
	var z D

	switch {
	case spec.Eq(x, magInf[SPEC]()):
		return z // +0 with no exception
	case spec.IsZero(m):
		// EXCEPTION: divide by zero
		return spec.Or(s, magInf[SPEC]())
	}

	return rcp[SPEC](sqrt[SPEC](x, rounding), rounding)
}

func hypot[SPEC spec[D], D datum](x, y D, rounding RoundingMode) D {
	_, x = mag[SPEC](x)
	_, y = mag[SPEC](y)

	var spec SPEC
	var z D

	magInf := magInf[SPEC]()

	switch {
	case spec.Eq(x, magInf) || spec.Eq(y, magInf):
		return magInf
	case spec.Gt(x, magInf):
		return x
	case spec.Gt(y, magInf):
		return y
	}

	if spec.Lt(x, y) {
		x, y = y, x
	}

	if spec.IsZero(x) {
		return z
	}

	y = div[SPEC](y, x, rounding)
	return mul[SPEC](x, sqrt[SPEC](madd[SPEC](y, y, one[SPEC](), rounding), rounding), rounding)

}

func exp[SPEC spec[D], D datum](x D, rounding RoundingMode) D {
	s, m := mag[SPEC](x)

	var spec SPEC
	var z D

	overflowVal, underflowVal, nearZero := spec.expOverUnder()

	switch {
	case spec.Gte(m, magInf[SPEC]()):
		if spec.Eq(x, inf[SPEC](true)) {
			// 2**-∞ == 0
			return z
		}

		// 2**∞ = ∞, NaN → NaN
		return x

	case spec.IsZero(s):
		if spec.Gte(m, overflowVal) {
			// EXCEPTION: overflow
			return overflow[SPEC](false, rounding)
		}

		// argument reduction; x = r×lg(e) + k with |r| ≤ ln(2)/2.

	case !spec.IsZero(s):
		if spec.Gt(m, underflowVal) {
			// We’re dealing with abs(x) here, so it must be _greater than_ Underflow.

			// EXCEPTION: underflow
			return z
		}

	case spec.Lt(m, nearZero):
		return add[SPEC](one[SPEC](), x, rounding)
	}

	Ln2Hi, Ln2Lo, Ln2E := spec.ln2HiLoE()

	k := round[SPEC](mul[SPEC](Ln2E, x, rounding))

	hi := mnsub[SPEC](k, Ln2Hi, x, rounding)
	lo := mul[SPEC](k, Ln2Lo, rounding)

	return expmulti[SPEC](hi, lo, truncToInt[SPEC](k), rounding)
}

func exp2[SPEC spec[D], D datum](x D, rounding RoundingMode) D {
	s, m := mag[SPEC](x)

	var spec SPEC
	var z D

	overflowVal, underflowVal := spec.exp2OverUnder()

	switch {
	case spec.Gte(m, magInf[SPEC]()):
		if spec.Eq(x, inf[SPEC](true)) {
			// 2**-∞ == 0
			return z
		}

		// 2**∞ = ∞, NaN → NaN
		return x

	case spec.IsZero(s):
		if spec.Gte(m, overflowVal) {
			// EXCEPTION: overflow
			return overflow[SPEC](false, rounding)
		}

	case !spec.IsZero(s):
		// We’re dealing with abs(x) here, so it must be _greater than_ Underflow.
		if spec.Gt(m, underflowVal) {
			// EXCEPTION: underflow
			return z
		}
	}

	// argument reduction; x = r×lg(e) + k with |r| ≤ ln(2)/2.
	k := round[SPEC](x)
	t := sub[SPEC](x, k, rounding)

	Ln2Hi, Ln2Lo, _ := spec.ln2HiLoE()

	hi := mul[SPEC](t, Ln2Hi, rounding)
	lo := mul[SPEC](neg[SPEC](t), Ln2Lo, rounding)

	return expmulti[SPEC](hi, lo, truncToInt[SPEC](k), rounding)
}

// expmulti returns e**r × 2**k where r = hi - lo and |r| ≤ ln(2)/2
func expmulti[SPEC spec[D], D datum](hi, lo D, k int, rounding RoundingMode) D {
	r := sub[SPEC](hi, lo, rounding)
	t := mul[SPEC](r, r, rounding)

	var spec SPEC

	PP := spec.expPN()

	// pp = (P1 + t(P2 + t(P3 + t(P4 + t·P5))))
	pp := PP[0]
	for _, pn := range PP[1:] {
		pp = madd[SPEC](pp, t, pn, rounding)
	}

	// c = r - t·pp
	c := mnsub[SPEC](pp, t, r, rounding)

	// y = 1 - ((lo - (r·c)/(2-c)) - hi)

	rc := mul[SPEC](r, c, rounding)
	tmc := rcp[SPEC](sub[SPEC](two[SPEC](), c, rounding), rounding)

	lo = mnsub[SPEC](rc, tmc, lo, rounding)
	hi = sub[SPEC](lo, hi, rounding)

	y := sub[SPEC](one[SPEC](), hi, rounding)

	return ldexp[SPEC](y, k)
}

func ilogb[SPEC spec[D], D datum](x D) (int, bool) {
	_, m := mag[SPEC](x)

	var spec SPEC

	switch {
	case spec.IsZero(m):
		// EXCEPTION: invalid operation: ilogb(0) = ⟨implementation defined⟩
		// But the result must be outside the range ±2×(emax+p-1)
		return math.MinInt, false
	case spec.Gte(m, magInf[SPEC]()):
		// EXCEPTION: invalid operation: ilogb(NaN or ±∞) = ⟨implementation defined⟩
		// But the result must be outside the range ±2×(emax+p-1)
		return math.MaxInt, false
	}

	_, e, _ := decomp[SPEC](x)

	return e - expBias[SPEC](), true
}

func logb[SPEC spec[D], D datum](x D) D {
	_, m := mag[SPEC](x)

	var spec SPEC

	switch {
	case spec.IsZero(m):
		// EXCEPTION: divide by zero: logB(0) = -∞
		return inf[SPEC](true)

	case spec.Gte(m, magInf[SPEC]()):
		if spec.Eq(m, magInf[SPEC]()) {
			return magInf[SPEC]() // logB(±∞) = +∞
		}

		return nan[SPEC]()
	}

	_, e, _ := decomp[SPEC](x)

	return spec.FromInt(e - expBias[SPEC]())
}
