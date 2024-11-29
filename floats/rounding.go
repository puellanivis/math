package floats

import (
	"fmt"

	"github.com/puellanivis/math/bits"
)

var _ = fmt.Println

// RoundingMode defines a rounding mode.
type RoundingMode interface {
	finiteOverflow(sign bool) bool

	round16(f *binary[binary16, uint16])
	round32(f *binary[binary32, uint32])
	round64(f *binary[binary64, uint64])
	round128(f *binary[binary128, bits.Uint128])
	roundBF16(f *binary[bfloat16, uint16])
}

func incAway[SPEC spec[D], D datum]() D {
	var spec SPEC

	return spec.Pow2m1(spec.expWidth())
}

func incNear[SPEC spec[D], D datum]() D {
	var spec SPEC

	return spec.Pow2(spec.expWidth() - 1)
}

func incNearEven[SPEC spec[D], D datum](x D) (m, r D) {
	var spec SPEC
	var z D

	x = spec.And(spec.Shr(x, spec.expWidth()), spec.Pow2(0))

	r, carry := spec.Add(spec.Not(z), x, z)
	m, _ = spec.Add(incNear[SPEC](), spec.Not(z), carry)
	return m, r
}

func incNearOdd[SPEC spec[D], D datum](x D) (m, r D) {
	var spec SPEC
	var z D

	x = spec.And(spec.Not(spec.Shr(x, spec.expWidth())), spec.Pow2(0))

	r, carry := spec.Add(spec.Not(z), x, z)
	m, _ = spec.Add(incNear[SPEC](), spec.Not(z), carry)
	return m, r
}

func overflow[SPEC spec[D], D datum](sign bool, rounding RoundingMode) D {
	x := inf[SPEC](sign)

	if rounding.finiteOverflow(sign) {
		var spec SPEC
		return spec.Dec(x)
	}

	return x
}

func applyRounding[B binary[SPEC, D], SPEC spec[D], D datum](f *B, rounding RoundingMode) {
	switch f := any(f).(type) {
	case *binary[binary16, uint16]:
		rounding.round16(f)
	case *binary[binary32, uint32]:
		rounding.round32(f)
	case *binary[binary64, uint64]:
		rounding.round64(f)
	case *binary[binary128, bits.Uint128]:
		rounding.round128(f)
	case *binary[bfloat16, uint16]:
		rounding.roundBF16(f)
	default:
		panic(fmt.Errorf("unsupported type in closed type-switch: %T", f))
	}
}

// RoundTowardZero rounds infinitely precise results to the floating-point numbers
// closest to and no greater in magnitude than the infinitely precise result.
//
// IEEE-754 roundTowardZero
type RoundTowardZero struct{}

func (RoundTowardZero) finiteOverflow(_ bool) bool {
	return true
}

func (RoundTowardZero) round16(f *binary[binary16, uint16]) {
	f.trunc()
}

func (RoundTowardZero) round32(f *binary[binary32, uint32]) {
	f.trunc()
}

func (RoundTowardZero) round64(f *binary[binary64, uint64]) {
	f.trunc()
}

func (RoundTowardZero) round128(f *binary[binary128, bits.Uint128]) {
	f.trunc()
}

func (RoundTowardZero) roundBF16(f *binary[bfloat16, uint16]) {
	f.trunc()
}

// RoundTowardPositive rounds infinitely precise results to the floating-point numbers
// (possibly +∞) closest to and no lesser than the infinitely precise result.
//
// IEEE-754 roundTowardPositive
type RoundTowardPositive struct{}

func (RoundTowardPositive) finiteOverflow(sign bool) bool {
	return sign
}

func (RoundTowardPositive) round16(f *binary[binary16, uint16]) {
	inf, nan := f.classify()
	switch {
	case inf, nan:
	case !f.s:
		// for positive numbers, this is a round away from zero.
		f.add(incAway[binary16](), 0)
		fallthrough
	default:
		// for negative numbers, this is a truncation
		f.trunc()
	}
}

func (RoundTowardPositive) round32(f *binary[binary32, uint32]) {
	inf, nan := f.classify()
	switch {
	case inf, nan:
	case !f.s:
		// for positive numbers, this is a round away from zero.
		f.add(incAway[binary32](), 0)
		fallthrough
	default:
		// for negative numbers, this is a truncation
		f.trunc()
	}
}

func (RoundTowardPositive) round64(f *binary[binary64, uint64]) {
	inf, nan := f.classify()
	switch {
	case inf, nan:
	case !f.s:
		// for positive numbers, this is a round away from zero.
		f.add(incAway[binary64](), 0)
		fallthrough
	default:
		// for negative numbers, this is a truncation
		f.trunc()
	}
}

func (RoundTowardPositive) round128(f *binary[binary128, bits.Uint128]) {
	inf, nan := f.classify()
	switch {
	case inf, nan:
	case !f.s:
		// for positive numbers, this is a round away from zero.
		f.add(incAway[binary128](), bits.Uint128{})
		fallthrough
	default:
		// for negative numbers, this is a truncation
		f.trunc()
	}
}

func (RoundTowardPositive) roundBF16(f *binary[bfloat16, uint16]) {
	inf, nan := f.classify()
	switch {
	case inf, nan:
	case !f.s:
		// for positive numbers, this is a round away from zero.
		f.add(incAway[binary16](), 0)
		fallthrough
	default:
		// for negative numbers, this is a truncation
		f.trunc()
	}
}

// RoundTowardNegative rounds infinitely precise results to the floating-point numbers
// (possibly -∞) closest to and no greater than the infinitely precise result.
//
// IEEE-754 roundTowardNegative
type RoundTowardNegative struct{}

func (RoundTowardNegative) finiteOverflow(sign bool) bool {
	return !sign
}

func (RoundTowardNegative) round16(f *binary[binary16, uint16]) {
	inf, nan := f.classify()
	switch {
	case inf, nan:
	case f.s:
		// for negative numbers, this is a round away from zero.
		f.add(incAway[binary16](), 0)
		fallthrough
	default:
		// for positive numbers, this is a truncation.
		f.trunc()
	}
}

func (RoundTowardNegative) round32(f *binary[binary32, uint32]) {
	inf, nan := f.classify()
	switch {
	case inf, nan:
	case f.s:
		// for negative numbers, this is a round away from zero.
		f.add(incAway[binary32](), 0)
		fallthrough
	default:
		// for positive numbers, this is a truncation.
		f.trunc()
	}
}

func (RoundTowardNegative) round64(f *binary[binary64, uint64]) {
	inf, nan := f.classify()
	switch {
	case inf, nan:
	case f.s:
		// for negative numbers, this is a round away from zero.
		f.add(incAway[binary64](), 0)
		fallthrough
	default:
		// for positive numbers, this is a truncation.
		f.trunc()
	}
}

func (RoundTowardNegative) round128(f *binary[binary128, bits.Uint128]) {
	inf, nan := f.classify()
	switch {
	case inf, nan:
	case f.s:
		// for negative numbers, this is a round away from zero.
		f.add(incAway[binary128](), bits.Uint128{})
		fallthrough
	default:
		// for positive numbers, this is a truncation.
		f.trunc()
	}
}

func (RoundTowardNegative) roundBF16(f *binary[bfloat16, uint16]) {
	inf, nan := f.classify()
	switch {
	case inf, nan:
	case f.s:
		// for negative numbers, this is a round away from zero.
		f.add(incAway[binary16](), 0)
		fallthrough
	default:
		// for positive numbers, this is a truncation.
		f.trunc()
	}
}

// RoundTiesToAway rounds infinitely precise results to the floating-point numbers
// (possibly ±∞) nearest to the infinitely precise result;
// if the two nearest floating-point numbers bracketing an unrepresentable infinitely precise result are equally near,
// it will return the one with larger magnitude.
//
// IEEE-754 roundTiesToAway
type RoundTiesToAway struct{}

func (RoundTiesToAway) finiteOverflow(_ bool) bool {
	return false
}

func (RoundTiesToAway) round16(f *binary[binary16, uint16]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNear[binary16](), 0)
	f.trunc()
}

func (RoundTiesToAway) round32(f *binary[binary32, uint32]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNear[binary32](), 0)
	f.trunc()
}

func (RoundTiesToAway) round64(f *binary[binary64, uint64]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNear[binary64](), 0)
	f.trunc()
}

func (RoundTiesToAway) round128(f *binary[binary128, bits.Uint128]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNear[binary128](), bits.Uint128{})
	f.trunc()
}

func (RoundTiesToAway) roundBF16(f *binary[bfloat16, uint16]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNear[binary16](), 0)
	f.trunc()
}

// RoundTiesToEven rounds infinitely precise results to the floating-point numbers
// (possibly ±∞) nearest to the infinitely precise result;
// if the two nearest floating-point numbers bracketing an unrepresentable infinitely precise result are equally near,
// it will return the one with an even least significant digit.
//
// IEEE-754 roundTiesToEven
type RoundTiesToEven struct{}

func (RoundTiesToEven) finiteOverflow(_ bool) bool {
	return false
}

func (RoundTiesToEven) round16(f *binary[binary16, uint16]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNearEven[binary16](f.m))
	f.trunc()
}

func (RoundTiesToEven) round32(f *binary[binary32, uint32]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNearEven[binary32](f.m))
	f.trunc()
}

func (RoundTiesToEven) round64(f *binary[binary64, uint64]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNearEven[binary64](f.m))
	f.trunc()
}

func (RoundTiesToEven) round128(f *binary[binary128, bits.Uint128]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNearEven[binary128](f.m))
	f.trunc()
}

func (RoundTiesToEven) roundBF16(f *binary[bfloat16, uint16]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNearEven[bfloat16](f.m))
	f.trunc()
}

// RoundTiesToOdd rounds infinitely precise results to the floating-point numbers
// (possibly ±∞) nearest to the infinitely precise result;
// if the two nearest floating-point numbers bracketing an unrepresentable infinitely precise result are equally near,
// it will return the one with an odd least significant digit.
type RoundTiesToOdd struct{}

func (RoundTiesToOdd) finiteOverflow(_ bool) bool {
	return false
}

func (RoundTiesToOdd) round16(f *binary[binary16, uint16]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNearOdd[binary16](f.m))
	f.trunc()
}

func (RoundTiesToOdd) round32(f *binary[binary32, uint32]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNearOdd[binary32](f.m))
	f.trunc()
}

func (RoundTiesToOdd) round64(f *binary[binary64, uint64]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNearOdd[binary64](f.m))
	f.trunc()
}

func (RoundTiesToOdd) round128(f *binary[binary128, bits.Uint128]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNearOdd[binary128](f.m))
	f.trunc()
}

func (RoundTiesToOdd) roundBF16(f *binary[bfloat16, uint16]) {
	inf, nan := f.classify()
	if inf || nan {
		return
	}
	f.add(incNearOdd[bfloat16](f.m))
	f.trunc()
}
