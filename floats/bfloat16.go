package floats

import (
	"fmt"
	"math"
	"math/big"
)

// Google Brain floating-point limit values.
// Max is the largest finite value representable by the type.
// SmallestNonZero is the smallest positive non-zero value representable by the type.
var (
	MaxBFloat16             = BFloat16{0b0_11111110_1111111} // 3.3895e+38
	SmallestNonzeroBFloat16 = BFloat16{0b0_00000000_0000001} // 9.1835e-41
)

// InfBFloat16 returns a Google Brain floating-point encoded positive infinity if sign >= 0, negative infinity if sign < 0.
func InfBFloat16(sign bool) BFloat16 {
	return BFloat16{inf[bfloat16](sign)}
}

// NaNBFloat16 returns a Google Brain floating-point encoded "not-a-number" value.
func NaNBFloat16() BFloat16 {
	return BFloat16{nan[bfloat16]()}
}

// BFloat16WithRound is a Google Brain floating-point number with specified rounding.
type BFloat16WithRound[RND RoundingMode] struct {
	bits uint16
}

// BFloat16 is an alias to a Google Brain floating-point number with rounding toward nearest, with ties to even.
type BFloat16 = BFloat16WithRound[RoundTiesToEven]

// BFloat16FromBits returns the Google Brain floating-point number corresponding to the binary representation of bits.
// BFloat16FromBits(x).Bits() == x
func BFloat16FromBits(bits uint16) BFloat16 {
	return BFloat16{bits}
}

// BFloat16WithRoundFromBits returns the Google Brain floating-point number with specified rounding corresponding to the binary representation of bits.
// BFloat16WithRoundingFromBits[RoundingMode](x).Bits() == x
func BFloat16WithRoundFromBits[RND RoundingMode](bits uint16) BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{bits}
}

// BFloat16FromFloat returns the Google Brain floating-point number closest in representation to the given floating point argument using round toward nearest, with ties to even.
func BFloat16FromFloat[F ~float32 | ~float64 | *big.Float](val F) BFloat16 {
	return BFloat16WithRoundFromFloat[RoundTiesToEven](val)
}

// BFloat16WithRoundFromFloat returns the Google Brain floating-point number closest in representation to the given floating point argument using the specified rounding mode.
func BFloat16WithRoundFromFloat[RND RoundingMode, F ~float32 | ~float64 | *big.Float](val F) BFloat16WithRound[RND] {
	var rnd RND

	switch v := any(val).(type) {
	case float32:
		if math.IsNaN(float64(v)) {
			return BFloat16WithRound[RND](NaNBFloat16())
		}
		return BFloat16WithRound[RND]{convert[binary32, bfloat16](math.Float32bits(v), rnd)}
	case float64:
		if math.IsNaN(v) {
			return BFloat16WithRound[RND](NaNBFloat16())
		}
		return BFloat16WithRound[RND]{convert[binary64, bfloat16](math.Float64bits(v), rnd)}
	case *big.Float:
		return BFloat16WithRound[RND]{fromBigFloat[bfloat16, uint16](v, rnd)}
	default:
		panic(fmt.Sprintf("impossible type passed into BFloat16FromFloat: %T", v))
	}
}

// Format implements [fmt.Formatter].
func (x BFloat16WithRound[RND]) Format(f fmt.State, verb rune) {
	var rnd RND

	format[bfloat16](x.bits, f, verb, rnd)
}

// Float16 returns the number converted to an IEEE 754 16-bit floating-point number.
func (x BFloat16WithRound[RND]) Float16() Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{convert[bfloat16, binary16](x.bits, rnd)}
}

// Float32 returns the number converted to an IEEE 754 32-bit floating-point number.
func (x BFloat16WithRound[RND]) Float32() Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{convert[bfloat16, binary32](x.bits, rnd)}
}

// Bin64 returns the number converted to an IEEE 754 64-bit floating-point number.
func (x BFloat16WithRound[RND]) Float64() Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{convert[bfloat16, binary64](x.bits, rnd)}
}

// Float128 returns the number converted to an IEEE 754 128-bit floating-point number.
func (x BFloat16WithRound[RND]) Float128() Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{convert[bfloat16, binary128](x.bits, rnd)}
}

// BFloat16 returns the number converted to a Google Brain floating-point number.
func (x BFloat16WithRound[RND]) BFloat16() BFloat16WithRound[RND] {
	return x
}

// Bits returns the Google Brain floating-point encoded binary representation of the number.
// BFloat16WithRoundingFromBits[RoundingMode](x).Bits() == x
func (x BFloat16WithRound[RND]) Bits() uint16 {
	return x.bits
}

func (x BFloat16WithRound[RND]) IsInf(sign int) bool {
	if ok := isInf[bfloat16](x.bits); !ok {
		return false
	}

	if sign == 0 {
		return true
	}

	return sign == getSign[bfloat16](x.bits)
}

func (x BFloat16WithRound[RND]) IsNaN() bool {
	return isNaN[bfloat16](x.bits)
}

func (x BFloat16WithRound[RND]) Sign() int {
	return getSign[bfloat16](x.bits)
}

func (x BFloat16WithRound[RND]) SignBit() bool {
	return signBit[bfloat16](x.bits)
}

func (x BFloat16WithRound[RND]) Abs() BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{abs[bfloat16](x.bits)}
}

func (x BFloat16WithRound[RND]) Neg() BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{neg[bfloat16](x.bits)}
}

func (x BFloat16WithRound[RND]) CopySign(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{copySign[bfloat16](x.bits, y.bits)}
}

func (x BFloat16WithRound[RND]) NextUp() BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{nextUp[bfloat16](x.bits)}
}

func (x BFloat16WithRound[RND]) NextDown() BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{nextDown[bfloat16](x.bits)}
}

func (x BFloat16WithRound[RND]) Add(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{add[bfloat16](x.bits, y.bits, rnd)}
}

func (x BFloat16WithRound[RND]) Sub(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{sub[bfloat16](x.bits, y.bits, rnd)}
}

func (x BFloat16WithRound[RND]) Dim(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{dim[bfloat16](x.bits, y.bits, rnd)}
}

func (x BFloat16WithRound[RND]) Mul(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{mul[bfloat16](x.bits, y.bits, rnd)}
}

func (x BFloat16WithRound[RND]) Div(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{div[bfloat16](x.bits, y.bits, rnd)}
}

func (x BFloat16WithRound[RND]) Mod(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{mod[bfloat16](x.bits, y.bits, rnd)}
}

func (x BFloat16WithRound[RND]) ModF() (i, f BFloat16WithRound[RND]) {
	q, r := modf[bfloat16](x.bits)
	return BFloat16WithRound[RND]{q}, BFloat16WithRound[RND]{r}
}

func (x BFloat16WithRound[RND]) Less(y BFloat16WithRound[RND]) bool {
	return less[bfloat16](x.bits, y.bits)
}

func (x BFloat16WithRound[RND]) Compare(y BFloat16WithRound[RND]) int {
	return compare[bfloat16](x.bits, y.bits)
}

func (x BFloat16WithRound[RND]) Equal(y BFloat16WithRound[RND]) bool {
	order, _ := fcmp[bfloat16](x.bits, y.bits)
	return order == 0
}

func (x BFloat16WithRound[RND]) Cmp(y BFloat16WithRound[RND]) (order int, ordered bool) {
	return fcmp[bfloat16](x.bits, y.bits)
}

func (x BFloat16WithRound[RND]) Min(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{fmin[bfloat16](x.bits, y.bits)}
}

func (x BFloat16WithRound[RND]) Max(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{fmax[bfloat16](x.bits, y.bits)}
}

func (x BFloat16WithRound[RND]) CmpMag(y BFloat16WithRound[RND]) (order int, ordered bool) {
	return fcmpMag[bfloat16](x.bits, y.bits)
}

func (x BFloat16WithRound[RND]) MinMag(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{fminMag[bfloat16](x.bits, y.bits)}
}

func (x BFloat16WithRound[RND]) MaxMag(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{fmaxMag[bfloat16](x.bits, y.bits)}
}

func (x BFloat16WithRound[RND]) Round() BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{round[bfloat16](x.bits)}
}

func (x BFloat16WithRound[RND]) RoundToEven() BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{roundToEven[bfloat16](x.bits)}
}

func (x BFloat16WithRound[RND]) Floor() BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{floor[bfloat16](x.bits)}
}

func (x BFloat16WithRound[RND]) Trunc() BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{trunc[bfloat16](x.bits)}
}

func (x BFloat16WithRound[RND]) Ceil() BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{ceil[bfloat16](x.bits)}
}

func (x BFloat16WithRound[RND]) Sqrt() BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{sqrt[bfloat16](x.bits, rnd)}
}

func (x BFloat16WithRound[RND]) RSqrt() BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{rsqrt[bfloat16](x.bits, rnd)}
}

func (x BFloat16WithRound[RND]) Hypot(y BFloat16WithRound[RND]) BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{hypot[bfloat16](x.bits, y.bits, rnd)}
}

func (x BFloat16WithRound[RND]) Exp() BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{exp[bfloat16](x.bits, rnd)}
}

func (x BFloat16WithRound[RND]) Exp2() BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{exp2[bfloat16](x.bits, rnd)}
}

func (x BFloat16WithRound[RND]) LogB() BFloat16WithRound[RND] {
	return BFloat16WithRound[RND]{logb[bfloat16](x.bits)}
}

func (x BFloat16WithRound[RND]) ILogB() (int, bool) {
	return ilogb[bfloat16](x.bits)
}
