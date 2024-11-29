package floats

import (
	"fmt"
	"math"
	"math/big"
)

var (
	MaxFloat64             = Float64{0b0_11111111110_1111111111111111111111111111111111111111111111111111} // 1.79769313486231570814527423731704356798070e+308
	SmallestNonzeroFloat64 = Float64{0b0_00000000000_0000000000000000000000000000000000000000000000000001} // 4.9406564584124654417656879286822137236505980e-324
)

func Inf64(sign bool) Float64 {
	return Float64{inf[binary64](sign)}
}

func NaN64() Float64 {
	return Float64{nan[binary64]()}
}

type Float64WithRound[RND RoundingMode] struct {
	bits uint64
}

type Float64 = Float64WithRound[RoundTiesToEven]

func Float64FromBits(bits uint64) Float64 {
	return Float64{bits}
}

func Float64WithRoundFromBits[RND RoundingMode](bits uint64) Float64WithRound[RND] {
	return Float64WithRound[RND]{bits}
}

func Float64FromFloat[F ~float32 | ~float64 | *big.Float](val F) Float64 {
	return Float64WithRoundFromFloat[RoundTiesToEven](val)
}

func Float64WithRoundFromFloat[RND RoundingMode, F ~float32 | ~float64 | *big.Float](val F) Float64WithRound[RND] {
	var rnd RND

	switch v := any(val).(type) {
	case float32:
		return Float64WithRound[RND]{convert[binary32, binary64](math.Float32bits(v), rnd)}
	case float64:
		return Float64WithRound[RND]{convert[binary64, binary64](math.Float64bits(v), rnd)}
	case *big.Float:
		return Float64WithRound[RND]{fromBigFloat[binary64, uint64](v, rnd)}
	default:
		panic("impossible type passed into Float64FromFloat")
	}
}

func (x Float64WithRound[RND]) Native() float64 {
	return math.Float64frombits(x.bits)
}

func (x Float64WithRound[RND]) Format(f fmt.State, verb rune) {
	var rnd RND

	format[binary64](x.bits, f, verb, rnd)
}

func (x Float64WithRound[RND]) Float16() Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{convert[binary64, binary16](x.bits, rnd)}
}

func (x Float64WithRound[RND]) Float32() Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{convert[binary64, binary32](x.bits, rnd)}
}

func (x Float64WithRound[RND]) Float64() Float64WithRound[RND] {
	return x
}

func (x Float64WithRound[RND]) Float128() Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{convert[binary64, binary128](x.bits, rnd)}
}

func (x Float64WithRound[RND]) BFloat16() BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{convert[binary64, bfloat16](x.bits, rnd)}
}

func (x Float64WithRound[RND]) Bits() uint64 {
	return x.bits
}

func (x Float64WithRound[RND]) IsInf(sign int) bool {
	if ok := isInf[binary64](x.bits); !ok {
		return false
	}

	if sign == 0 {
		return true
	}

	return sign == getSign[binary64](x.bits)
}

func (x Float64WithRound[RND]) IsNaN() bool {
	return isNaN[binary64](x.bits)
}

func (x Float64WithRound[RND]) Sign() int {
	return getSign[binary64](x.bits)
}

func (x Float64WithRound[RND]) SignBit() bool {
	return signBit[binary64](x.bits)
}

func (x Float64WithRound[RND]) Abs() Float64WithRound[RND] {
	return Float64WithRound[RND]{abs[binary64](x.bits)}
}

func (x Float64WithRound[RND]) Neg() Float64WithRound[RND] {
	return Float64WithRound[RND]{neg[binary64](x.bits)}
}

func (x Float64WithRound[RND]) CopySign(y Float64WithRound[RND]) Float64WithRound[RND] {
	return Float64WithRound[RND]{copySign[binary64](x.bits, y.bits)}
}

func (x Float64WithRound[RND]) NextUp() Float64WithRound[RND] {
	return Float64WithRound[RND]{nextUp[binary64](x.bits)}
}

func (x Float64WithRound[RND]) NextDown() Float64WithRound[RND] {
	return Float64WithRound[RND]{nextDown[binary64](x.bits)}
}

func (x Float64WithRound[RND]) Add(y Float64WithRound[RND]) Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{add[binary64](x.bits, y.bits, rnd)}
}

func (x Float64WithRound[RND]) Sub(y Float64WithRound[RND]) Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{sub[binary64](x.bits, y.bits, rnd)}
}

func (x Float64WithRound[RND]) Dim(y Float64WithRound[RND]) Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{dim[binary64](x.bits, y.bits, rnd)}
}

func (x Float64WithRound[RND]) Mul(y Float64WithRound[RND]) Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{mul[binary64](x.bits, y.bits, rnd)}
}

func (x Float64WithRound[RND]) Div(y Float64WithRound[RND]) Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{div[binary64](x.bits, y.bits, rnd)}
}

func (x Float64WithRound[RND]) Mod(y Float64WithRound[RND]) Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{mod[binary64](x.bits, y.bits, rnd)}
}

func (x Float64WithRound[RND]) ModF() (i, f Float64WithRound[RND]) {
	q, r := modf[binary64](x.bits)
	return Float64WithRound[RND]{q}, Float64WithRound[RND]{r}
}

func (x Float64WithRound[RND]) Less(y Float64WithRound[RND]) bool {
	return less[binary64](x.bits, y.bits)
}

func (x Float64WithRound[RND]) Compare(y Float64WithRound[RND]) int {
	return compare[binary64](x.bits, y.bits)
}

func (x Float64WithRound[RND]) Equal(y Float64WithRound[RND]) bool {
	order, _ := fcmp[binary64](x.bits, y.bits)
	return order == 0
}

func (x Float64WithRound[RND]) Cmp(y Float64WithRound[RND]) (order int, ordered bool) {
	return fcmp[binary64](x.bits, y.bits)
}

func (x Float64WithRound[RND]) Min(y Float64WithRound[RND]) Float64WithRound[RND] {
	return Float64WithRound[RND]{fmin[binary64](x.bits, y.bits)}
}

func (x Float64WithRound[RND]) Max(y Float64WithRound[RND]) Float64WithRound[RND] {
	return Float64WithRound[RND]{fmax[binary64](x.bits, y.bits)}
}

func (x Float64WithRound[RND]) CmpMag(y Float64WithRound[RND]) (order int, ordered bool) {
	return fcmpMag[binary64](x.bits, y.bits)
}

func (x Float64WithRound[RND]) MinMag(y Float64WithRound[RND]) Float64WithRound[RND] {
	return Float64WithRound[RND]{fminMag[binary64](x.bits, y.bits)}
}

func (x Float64WithRound[RND]) MaxMag(y Float64WithRound[RND]) Float64WithRound[RND] {
	return Float64WithRound[RND]{fmaxMag[binary64](x.bits, y.bits)}
}

func (x Float64WithRound[RND]) Round() Float64WithRound[RND] {
	return Float64WithRound[RND]{round[binary64](x.bits)}
}

func (x Float64WithRound[RND]) RoundToEven() Float64WithRound[RND] {
	return Float64WithRound[RND]{roundToEven[binary64](x.bits)}
}

func (x Float64WithRound[RND]) Floor() Float64WithRound[RND] {
	return Float64WithRound[RND]{floor[binary64](x.bits)}
}

func (x Float64WithRound[RND]) Trunc() Float64WithRound[RND] {
	return Float64WithRound[RND]{trunc[binary64](x.bits)}
}

func (x Float64WithRound[RND]) Ceil() Float64WithRound[RND] {
	return Float64WithRound[RND]{ceil[binary64](x.bits)}
}

func (x Float64WithRound[RND]) Sqrt() Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{sqrt[binary64](x.bits, rnd)}
}

func (x Float64WithRound[RND]) RSqrt() Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{rsqrt[binary64](x.bits, rnd)}
}

func (x Float64WithRound[RND]) Hypot(y Float64WithRound[RND]) Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{hypot[binary64](x.bits, y.bits, rnd)}
}

func (x Float64WithRound[RND]) Exp() Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{exp[binary64](x.bits, rnd)}
}

func (x Float64WithRound[RND]) Exp2() Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{exp2[binary64](x.bits, rnd)}
}

func (x Float64WithRound[RND]) LogB() Float64WithRound[RND] {
	return Float64WithRound[RND]{logb[binary64](x.bits)}
}

func (x Float64WithRound[RND]) ILogB() (int, bool) {
	return ilogb[binary64](x.bits)
}
