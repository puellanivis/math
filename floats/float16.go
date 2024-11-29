package floats

import (
	"fmt"
	"math"
	"math/big"
)

// IEEE 754 floating-point limit values.
// MaxFloat16 is the largest finite value representable by the type.
// SmallestNonZeroFloat16 is the smallest positive non-zero value representable by the type.
var (
	MaxFloat16             = Float16{0b0_11110_1111111111} // 6.5504e+04
	SmallestNonzeroFloat16 = Float16{0b0_00000_0000000001} // 5.96046e-08
)

// Inf16 returns a IEEE 754 encoded positive infinity if sign >= 0,
// negative infinity if sign < 0.
func Inf16(sign bool) Float16 {
	return Float16{inf[binary16](sign)}
}

// NaN16 returns an IEEE 754 encoded “not-a-number” value.
func NaN16() Float16 {
	return Float16{nan[binary16]()}
}

// Float16WithRound is a IEEE 754 16-bit floating point number with specified rounding.
type Float16WithRound[RND RoundingMode] struct {
	bits uint16
}

// Float16 is an alias to a IEEE 754 16-bit floating point number with with rounding mode RoundTiesToEven.
type Float16 = Float16WithRound[RoundTiesToEven]

// Float16FromBits returns the IEEE 754 floating-point number corresponding to the binary representation of bits.
// Float16FromBits(x).Bits() == x
func Float16FromBits(bits uint16) Float16 {
	return Float16{bits}
}

// Float16WithRoundFromBits returns the IEEE 754 floating-point number with specified rounding corresponding to the binary representation of bits.
// Float16WithRoundFromBits[RoundingMode](x).Bits() == x
func Float16WithRoundFromBits[RND RoundingMode](bits uint16) Float16WithRound[RND] {
	return Float16WithRound[RND]{bits}
}

// Float16WithRoundFromFloat returns the IEEE 754 floating-point number closest in representation to the given floating-point argument,
// using the RoundTiesToEven rounding mode.
func Float16FromFloat[F ~float32 | ~float64 | *big.Float](val F) Float16 {
	return Float16WithRoundFromFloat[RoundTiesToEven](val)
}

// Float16WithRoundFromFloat returns the IEEE 754 floating-point number closest in representation to the given floating-point argument,
// using the specified rounding mode.
func Float16WithRoundFromFloat[RND RoundingMode, F ~float32 | ~float64 | *big.Float](val F) Float16WithRound[RND] {
	var rnd RND

	switch v := any(val).(type) {
	case float32:
		if math.IsNaN(float64(v)) {
			return Float16WithRound[RND](NaN16())
		}
		return Float16WithRound[RND]{convert[binary32, binary16](math.Float32bits(v), rnd)}
	case float64:
		if math.IsNaN(v) {
			return Float16WithRound[RND](NaN16())
		}
		return Float16WithRound[RND]{convert[binary64, binary16](math.Float64bits(v), rnd)}
	case *big.Float:
		return Float16WithRound[RND]{fromBigFloat[binary16, uint16](v, rnd)}
	default:
		panic(fmt.Sprintf("impossible type passed into Float16FromFloat: %T", v))
	}
}

// Format implements [fmt.Formatter].
func (x Float16WithRound[RND]) Format(f fmt.State, verb rune) {
	var rnd RND

	format[binary16](x.bits, f, verb, rnd)
}

// Float16 returns the number converted to an IEEE 754 16-bit floating-point number.
// There is no loss of precision.
func (x Float16WithRound[RND]) Float16() Float16WithRound[RND] {
	return x
}

// Float32 returns the number converted to an IEEE 754 32-bit floating-point number.
// There is no loss of precision.
func (x Float16WithRound[RND]) Float32() Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{convert[binary16, binary32](x.bits, rnd)}
}

// Float64 returns the number converted to an IEEE 754 64-bit floating-point number.
// There is no loss of precision.
func (x Float16WithRound[RND]) Float64() Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{convert[binary16, binary64](x.bits, rnd)}
}

// Float128 returns the number converted to an IEEE 754 128-bit floating-point number.
// There is no loss of precision.
func (x Float16WithRound[RND]) Float128() Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{convert[binary16, binary128](x.bits, rnd)}
}

// BFloat16 returns the number converted to a Google Brain floating-point number.
// This conversion loses 3 bits of precision.
func (x Float16WithRound[RND]) BFloat16() BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{convert[binary16, bfloat16](x.bits, rnd)}
}

// Bits returns the IEEE 754 floating-point encoded binary representation of the number.
// Float16WithRoundingFromBits[RoundingMode](x).Bits() == x
func (x Float16WithRound[RND]) Bits() uint16 {
	return x.bits
}

// IsInf reports whether the number is a an infinity, according to sign.
// If sign > 0, then IsInf reports whether the number is positive infinity.
// If sign < 0, then IsInf reports whether the number is negative infinity.
// If sign == 0, then IsInf reports whether the number is either infinity.
func (x Float16WithRound[RND]) IsInf(sign int) bool {
	if ok := isInf[binary16](x.bits); !ok {
		return false
	}

	if sign == 0 {
		return true
	}

	return sign == getSign[binary16](x.bits)
}

// IsNaN reports whether the number is an IEEE 754 “not-a-number” value.
func (x Float16WithRound[RND]) IsNaN() bool {
	return isNaN[binary16](x.bits)
}

// Sign returns the sign of the number.
// If x > 0, then it returns 1.
// If x < 0, then it returns -1.
// If x == 0, then it returns 0.
// If x is an IEEE 754 “not-a-number” value, it returns ±1 depending upon the sign bit of the number.
func (x Float16WithRound[RND]) Sign() int {
	return getSign[binary16](x.bits)
}

// SignBit reports whether x is negative or negative zero.
func (x Float16WithRound[RND]) SignBit() bool {
	return signBit[binary16](x.bits)
}

// Abs returns the absolute value of x.
//
// Special cases are:
//   NaN.Abs() = NaN
//   ±Inf.Abs() = +Inf
func (x Float16WithRound[RND]) Abs() Float16WithRound[RND] {
	return Float16WithRound[RND]{abs[binary16](x.bits)}
}

// Neg returns the negative value of x.
//
// Special cases are:
//   NaN.Neg() = NaN
//   ±Inf.Neg() = ∓Inf
func (x Float16WithRound[RND]) Neg() Float16WithRound[RND] {
	return Float16WithRound[RND]{neg[binary16](x.bits)}
}

// CopySign returns a value with the magnitude of the number and the sign based on the sign bit of sign.
func (x Float16WithRound[RND]) CopySign(sign Float16WithRound[RND]) Float16WithRound[RND] {
	return Float16WithRound[RND]{copySign[binary16](x.bits, sign.bits)}
}

// NextUp returns the smallest IEEE 754 floating-point value that is greater than the number.
//
// Special cases are:
//   NaN.NextUp() = NaN
//   +Inf.NextUp() = +Inf
//   -Inf.NextUp() = -MaxFloat
func (x Float16WithRound[RND]) NextUp() Float16WithRound[RND] {
	return Float16WithRound[RND]{nextUp[binary16](x.bits)}
}

// NextDown returns the largest IEEE 754 floating-point value that is less than the number.
//
// Special cases are:
//   NaN.NextDown() = NaN
//   +Inf.NextDown() = MaxFloat
//   -Inf.NextDown() = -Inf
func (x Float16WithRound[RND]) NextDown() Float16WithRound[RND] {
	return Float16WithRound[RND]{nextDown[binary16](x.bits)}
}

// Add returns the sum of x+y.
//
// Special cases are:
//   x + NaN = NaN + y = NaN
//   ±Inf + ∓Inf = NaN
//   ±Inf + ±Inf = ±Inf
//   ±Inf + y = x + ±Inf = ±Inf
func (x Float16WithRound[RND]) Add(y Float16WithRound[RND]) Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{add[binary16](x.bits, y.bits, rnd)}
}

// Sub returns the difference of x-y.
//
// Special cases are:
//   x - NaN = NaN - y = NaN
//   ±Inf - ±Inf = NaN
//   ±Inf - ∓Inf = ±Inf
//   ±Inf - y = ±Inf
//   x - ±Inf = ∓Inf
func (x Float16WithRound[RND]) Sub(y Float16WithRound[RND]) Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{sub[binary16](x.bits, y.bits, rnd)}
}

// Dim returns the maximum of x-y or 0.
//
// Special cases are:
//    x.Dim(NaN) = NaN.Dim(y) = NaN
//    ±Inf.Dim(±Inf) = NaN
//    +Inf.Dim(-Inf) = +Inf
//    -Inf.Dim(+Inf) = 0
func (x Float16WithRound[RND]) Dim(y Float16WithRound[RND]) Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{dim[binary16](x.bits, y.bits, rnd)}
}

// Mul returns the product of x*y.
//
// Special cases are:
//   x * NaN = NaN * y = NaN
//   Inf * 0 = 0 * Inf = NaN
//   ±x * ±Inf = ±Inf * ±y = +Inf
//   ±x * ∓Inf = ±Inf * ∓y = -Inf
func (x Float16WithRound[RND]) Mul(y Float16WithRound[RND]) Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{mul[binary16](x.bits, y.bits, rnd)}
}

// Div returns the quotient of x/y.
//
// Special cases are:
//   x / NaN = NaN / y = NaN
//   Inf / Inf = NaN
//   0 / 0 = NaN
//   ±x / ±0 = ±Inf / ±y = +Inf
//   ±x / ∓0 = ±Inf / ∓y = -Inf
//   ±x / ±Inf = +0
//   ±x / ∓Inf = -0
func (x Float16WithRound[RND]) Div(y Float16WithRound[RND]) Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{div[binary16](x.bits, y.bits, rnd)}
}

// Mod returns the floating-point remainer of x/y.
// The magnitude of the result is less than y, and its sign agrees with that of x.
//
// Special cases are:
//   x.Mod(NaN) = NaN.Mod(y) = NaN
//   ±Inf.Mod(y) = NaN
//   x.Mod(±Inf) = x
//   x.Mod(0) = NaN
func (x Float16WithRound[RND]) Mod(y Float16WithRound[RND]) Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{mod[binary16](x.bits, y.bits, rnd)}
}

// Modf returns integer and fractional floating-point numbers that sum to x.
// Both values have the same sign as x.
//
// Special cases are:
//   NaN.Modf() = NaN, NaN
//   ±Inf.Modf() = ±Inf, NaN
func (x Float16WithRound[RND]) ModF() (i, f Float16WithRound[RND]) {
	q, r := modf[binary16](x.bits)
	return Float16WithRound[RND]{q}, Float16WithRound[RND]{r}
}

// Less returns true if x < y.
//
// Special cases are:
//   NaN.Less(NaN) == false
//   NaN.Less(y) == true
//   x.Less(NaN) = false
func (x Float16WithRound[RND]) Less(y Float16WithRound[RND]) bool {
	return less[binary16](x.bits, y.bits)
}

// Compare returns
//   -1 if x is less than y,
//    0 if x equals y,
//   +1 if x is greater than y.
//
// A NaN is considered less than any non-NaN,
// a NaN is considered equal to a NaN,
// and -0 = 0.
func (x Float16WithRound[RND]) Compare(y Float16WithRound[RND]) int {
	return compare[binary16](x.bits, y.bits)
}

// Equal returns true if x == y.
//
// Special cases are:
//   x.Equal(NaN) = NaN.Equal(y) = false
//   0.Equal(0) = true
func (x Float16WithRound[RND]) Equal(y Float16WithRound[RND]) bool {
	order, ordered := fcmp[binary16](x.bits, y.bits)
	return order == 0 && ordered
}

// Cmp returns
//    0, false if either x or y is NaN,
//   -1, true if x is less than y,
//    0, true if x equals y,
//   +1, true if x is greater than y.
func (x Float16WithRound[RND]) Cmp(y Float16WithRound[RND]) (order int, ordered bool) {
	return fcmp[binary16](x.bits, y.bits)
}

// Min returns the smaller of x or y.
//
// Special cases are:
//   x.Min(NaN) = NaN.Min(x) = x
//   x.Min(-Inf) = -Inf.Min(x) = -Inf
//   -0.Min(±0) = ±0.Min(-0) = -0
//
// This differs from math.Min in that it returns the number rather than the NaN, if one of them is NaN.
// The IEEE 754 standard says it should return the canonical number.
func (x Float16WithRound[RND]) Min(y Float16WithRound[RND]) Float16WithRound[RND] {
	return Float16WithRound[RND]{fmin[binary16](x.bits, y.bits)}
}

// Max returns the larger of x or y.
//
// Special cases are:
//   x.Max(NaN) = NaN.Max(x) = x
//   x.Max(+Inf) = +Inf.Max(x) = +Inf
//   +0.Max(±0) = ±0.Max(+0) = +0
//   -0.Max(-0) = -0
//
// This differs from math.Max in that it returns the number rather than the NaN, if one of them is NaN.
// The IEEE 754 standard says it should return the canonical number.
func (x Float16WithRound[RND]) Max(y Float16WithRound[RND]) Float16WithRound[RND] {
	return Float16WithRound[RND]{fmax[binary16](x.bits, y.bits)}
}

// CmpMag returns
//    0, false if either x or y is NaN,
//   -1, true if |x| is less than |y|,
//    0, true if |x| equals |y|,
//   +1, true if |x| is greater than |y|.
func (x Float16WithRound[RND]) CmpMag(y Float16WithRound[RND]) (order int, ordered bool) {
	return fcmpMag[binary16](x.bits, y.bits)
}

// MinMag returns the smaller of magnitude of x or y.
//
// Special cases are:
//   x.MinMag(NaN) = NaN.MinMag(x) = x
//   x.MinMag(±Inf) = ±Inf.MinMag(x) = ±Inf
//   ±Inf.MinMag(±Inf) = ±Inf.MinMag(∓Inf) = ±Inf
//   ±0.MinMag(±0) = ±0.MinMag(∓0) = ±0
func (x Float16WithRound[RND]) MinMag(y Float16WithRound[RND]) Float16WithRound[RND] {
	return Float16WithRound[RND]{fminMag[binary16](x.bits, y.bits)}
}

// MaxMag returns the larger of magnitude of x or y.
//
// Special cases are:
//   x.MaxMag(NaN) = NaN.MaxMag(x) = x
//   x.MaxMag(±Inf) = ±Inf.MaxMag(x) = ±Inf
//   ±Inf.MaxMag(±Inf) = ±Inf.MaxMag(∓Inf) = ±Inf
//   ±0.MaxMag(±0) = ±0.MaxMag(∓0) = ±0
func (x Float16WithRound[RND]) MaxMag(y Float16WithRound[RND]) Float16WithRound[RND] {
	return Float16WithRound[RND]{fmaxMag[binary16](x.bits, y.bits)}
}

// Round returns the nearest integer, rounding ties away from zero.
//
// Special cases are:
//   NaN.Round() = NaN
//   ±Inf.Round() = ±Inf
//   ±0.Round() = ±0
func (x Float16WithRound[RND]) Round() Float16WithRound[RND] {
	return Float16WithRound[RND]{round[binary16](x.bits)}
}

// Round returns the nearest integer, rounding ties to even.
//
// Special cases are:
//   NaN.RoundToEven() = NaN
//   ±Inf.RoundToEven() = ±Inf
//   ±0.RoundToEven() = ±0
func (x Float16WithRound[RND]) RoundToEven() Float16WithRound[RND] {
	return Float16WithRound[RND]{roundToEven[binary16](x.bits)}
}

// Floor returns the greatest integer value less than or equal to x.
//
// Special cases are:
//   NaN.Floor() = NaN
//   ±Inf.Floor() = ±Inf
//   ±0.Floor() = ±0
func (x Float16WithRound[RND]) Floor() Float16WithRound[RND] {
	return Float16WithRound[RND]{floor[binary16](x.bits)}
}

// Trunc returns the integer value of x.
//
// Special cases are:
//   NaN.Trunc() = NaN
//   ±Inf.Trunc() = ±Inf
//   ±0.Trunc() = ±0
func (x Float16WithRound[RND]) Trunc() Float16WithRound[RND] {
	return Float16WithRound[RND]{trunc[binary16](x.bits)}
}

// Ceil returns the least integer value greater than or equal to x.
//
// Special cases are:
//   NaN.Ceil() = NaN
//   ±Inf.Ceil() = ±Inf
//   ±0.Ceil() = ±0
func (x Float16WithRound[RND]) Ceil() Float16WithRound[RND] {
	return Float16WithRound[RND]{ceil[binary16](x.bits)}
}

// Sqrt returns the square root of x.
//
// Special cases are:
//   NaN.Sqrt = NaN
//   +Inf.Sqrt() = +Inf
//   ±0.Sqrt() = ±0
//   -x.Sqrt() = NaN
func (x Float16WithRound[RND]) Sqrt() Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{sqrt[binary16](x.bits, rnd)}
}

// RSqrt returns the recipricol of the square root of x.
//
// Special cases are:
//   NaN.RSqrt() = NaN
//   +Inf.RSqrt() = 0
//   ±0.RSqrt() = ±Inf
//   -x.RSqrt = NaN
func (x Float16WithRound[RND]) RSqrt() Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{rsqrt[binary16](x.bits, rnd)}
}

// Hypot returns x.Mul(x).Add(y.Mul(y)).Sqrt(), taking care to avoid unnecessary overflow and underflow.
//
// Special cases are:
//   NaN.Hypot(y) = x.Hypot(NaN) = NaN
//   ±Inf.Hypot(y) = x.Hypot(±Inf) = +Inf
func (x Float16WithRound[RND]) Hypot(y Float16WithRound[RND]) Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{hypot[binary16](x.bits, y.bits, rnd)}
}

// Exp returns e**x, the base-e exponential of x.
//
// Special cases are:
//   NaN.Exp() = NaN
//   +Inf.Exp() = +Inf
//
// Very large values overflow to 0 or +Inf.
// Very small values underflow to 1.
func (x Float16WithRound[RND]) Exp() Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{exp[binary16](x.bits, rnd)}
}

// Exp2 returns 2**x, the base-2 exponential of x.
//
// Special cases are the same as [Exp].
func (x Float16WithRound[RND]) Exp2() Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{exp2[binary16](x.bits, rnd)}
}

// LogB returns the binary exponent of x.
//
// Special cases are:
//   NaN.LogB() = NaN
//   ±Inf.LogB() = +Inf
//   0.LogB() = -Inf
func (x Float16WithRound[RND]) LogB() Float16WithRound[RND] {
	return Float16WithRound[RND]{logb[binary16](x.bits)}
}

// ILogB returns the binary exponent of x as an integer and true,
// or false if the exponent cannot be represented as an int.
//
// Special cases are:
//   NaN.ILogB() = [math.MaxInt], false
//   ±Inf.ILogB() = [math.MaxInt], false
//   0.ILogB() = [math.MinInt], false
//
// N.B.: This returns MaxInt32 and MinInt32 regardless of the size of int.
func (x Float16WithRound[RND]) ILogB() (int, bool) {
	return ilogb[binary16](x.bits)
}
