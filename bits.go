package math

import (
	"math"
)

// Bits returns the IEEE 754 binary representation of f,
// with the sign bit of f and the result in the same bit position,
// and Bits[UINT](FromBits[FLOAT](x)) == x.
func Bits[UINT uint32 | uint64, FLOAT Float](f FLOAT) UINT {
	switch any(f).(type) {
	case float32:
		return UINT(Bits32(float32(f)))
	case float64:
		return UINT(Bits64(float64(f)))
	default:
		panic("impossible type")
	}
}

// FromBits returns the floating-point number corresponding
// to the IEEE 754 binary representation b, with the sign bit of b
// and the result in the same bit position.
// FromBits[FLOAT](Bits[UINT](x)) == x.
func FromBits[FLOAT Float, UINT uint32 | uint64](bits UINT) FLOAT {
	var f FLOAT
	switch any(f).(type) {
	case float32:
		return FLOAT(FromBits32(uint32(bits)))
	case float64:
		return FLOAT(FromBits64(uint64(bits)))
	default:
		panic("impossible type")
	}
}

// Bits32 returns the IEEE 754 binary representation of f,
// with the sign bit of f and the result in the same bit position,
// and Bits32(FromBits32(x)) == x.
func Bits32(x float32) uint32 {
	return math.Float32bits(x)
}

// FromBits32 returns the floating-point number corresponding
// to the IEEE 754 binary representation b, with the sign bit of b
// and the result in the same bit position.
// FromBits32(Bits32(x)) == x.
func FromBits32(x uint32) float32 {
	return math.Float32frombits(x)
}

// Bits64 returns the IEEE 754 binary representation of f,
// with the sign bit of f and the result in the same bit position,
// and Bits64(FromBits64(x)) == x.
func Bits64(x float64) uint64 {
	return math.Float64bits(x)
}

// FromBits64 returns the floating-point number corresponding
// to the IEEE 754 binary representation b, with the sign bit of b
// and the result in the same bit position.
// FromBits64(Bits64(x)) == x.
func FromBits64(x uint64) float64 {
	return math.Float64frombits(x)
}

const (
	signBit32 = 1 << 31
)

// SignBit reports whether x is negative or negative zero.
func SignBit[FLOAT Float](x FLOAT) bool {
	switch x := any(x).(type) {
	case float32:
		return Bits32(x)&signBit32 != 0
	case float64:
		return math.Signbit(x)
	default:
		panic("impossible type")
	}
}

// CopySign returns a value with the magnitude of f and the sign of sign.
func CopySign[FLOAT Float](f, sign FLOAT) FLOAT {
	switch any(f).(type) {
	case float32:
		f, sign := cast2[float32](f, sign)
		return FLOAT(FromBits32(Bits32(f)&^signBit32 | Bits32(sign)&signBit32))
	case float64:
		return FLOAT(math.Copysign(cast2[float64](f, sign)))
	default:
		panic("impossible type")
	}
}

// Inf returns positive infinity if sign >= 0, negative infinity if sign < 0.
func Inf[FLOAT Float](sign int) FLOAT {
	return FLOAT(math.Inf(sign))
}

// Inf32 returns [Inf] as a float32.
func Inf32(sign int) float32 {
	return float32(math.Inf(sign))
}

// Inf64 returns [Inf] as a float64.
func Inf64(sign int) float64 {
	return float64(math.Inf(sign))
}

// IsInf reports whether f is an infinity, according to sign.
// If sign > 0, IsInf reports whether f is positive infinity.
// If sign < 0, IsInf reports whether f is negative infinity.
// If sign == 0, IsInf reports whether f is either infinity.
func IsInf[FLOAT Float](f FLOAT, sign int) (is bool) {
	return math.IsInf(float64(f), sign)
}

// NaN returns an IEEE 754 “not-a-number” value.
func NaN[FLOAT Float]() FLOAT {
	return FLOAT(math.NaN())
}

// NaN32 returns an IEEE 754 “not-a-number” value as a float32.
func NaN32() float32 {
	return float32(math.NaN())
}

// NaN64 returns an IEEE 754 “not-a-number” value as a float64.
func NaN64() float64 {
	return math.NaN()
}

// IsNaN reports whether f is an IEEE 754 “not-a-number” value.
func IsNaN[FLOAT Float](f FLOAT) (is bool) {
	// IEEE 754 says that only NaNs satisfy f != f.
	return f != f
}

// FrExp breaks f into a normalized fraction and an integral power of two.
// It returns frac and exp satisfying f == frac × 2**exp,
// with the absolute value of frac in the interval [½, 1).
//
// Special cases are:
//
//	FrExp(±0)   = ±0, 0
//	FrExp(±Inf) = ±Inf, 0
//	FrExp(NaN)  = NaN, 0
func FrExp[FLOAT Float](f FLOAT) (frac FLOAT, exp int) {
	fracʹ, exp := math.Frexp(float64(f))
	return FLOAT(fracʹ), exp
}

// LdExp is the inverse of [FrExp].
// It returns frac × 2**exp.
//
// Special cases are:
//
//	LdExp(±0, exp)   = ±0
//	LdExp(±Inf, exp) = ±Inf
//	LdExp(NaN, exp)  = NaN
func LdExp[FLOAT Float](frac FLOAT, exp int) FLOAT {
	return FLOAT(math.Ldexp(float64(frac), exp))
}
