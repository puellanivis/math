package floats

import (
	"fmt"
	"math"
	"math/big"

	"github.com/puellanivis/math/bits"
)

var (
	MaxFloat128             = Float128{bits.Uint128{Hi: 0x7ffeffffffffffff, Lo: 0xffffffffffffffff}} // 1.1897314953572317650857593266280070e+4932
	SmallestNonzeroFloat128 = Float128{bits.Uint128{Hi: 0x0000000000000000, Lo: 0x0000000000000001}} // 6.4751751194380251109244389582276466e-4966
)

func Inf128(sign bool) Float128 {
	return Float128{inf[binary128](sign)}
}

func NaN128() Float128 {
	return Float128{nan[binary128]()}
}

type Float128WithRound[RND RoundingMode] struct {
	bits bits.Uint128
}

type Float128 = Float128WithRound[RoundTiesToEven]

func Float128FromBits[RND RoundingMode](bits bits.Uint128) Float128 {
	return Float128{bits}
}

func Float128WithRoundFromBits[RND RoundingMode](bits bits.Uint128) Float128WithRound[RND] {
	return Float128WithRound[RND]{bits}
}

func Float128FromFloat[F ~float32 | ~float64 | *big.Float](val F) Float128 {
	return Float128WithRoundFromFloat[RoundTiesToEven](val)
}

func Float128WithRoundFromFloat[RND RoundingMode, F ~float32 | ~float64 | *big.Float](val F) Float128WithRound[RND] {
	var rnd RND

	switch v := any(val).(type) {
	case float32:
		return Float128WithRound[RND]{convert[binary32, binary128](math.Float32bits(v), rnd)}
	case float64:
		return Float128WithRound[RND]{convert[binary64, binary128](math.Float64bits(v), rnd)}
	case *big.Float:
		return Float128WithRound[RND]{fromBigFloat[binary128, bits.Uint128](v, rnd)}
	default:
		panic("impossible type passed into Float128FromFloat")
	}
}

func (x Float128WithRound[RND]) Format(f fmt.State, verb rune) {
	var rnd RND

	format[binary128](x.bits, f, verb, rnd)
}

func (x Float128WithRound[RND]) Float16() Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{convert[binary128, binary16](x.bits, rnd)}
}

func (x Float128WithRound[RND]) Float32() Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{convert[binary128, binary32](x.bits, rnd)}
}

func (x Float128WithRound[RND]) Float64() Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{convert[binary128, binary64](x.bits, rnd)}
}

func (x Float128WithRound[RND]) Float128() Float128WithRound[RND] {
	return x
}

func (x Float128WithRound[RND]) BFloat16() BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{convert[binary128, bfloat16](x.bits, rnd)}
}

func (x Float128WithRound[RND]) Bits() bits.Uint128 {
	return x.bits
}

func (x Float128WithRound[RND]) IsInf(sign int) bool {
	if ok := isInf[binary128](x.bits); !ok {
		return false
	}

	if sign == 0 {
		return true
	}

	return sign == getSign[binary128](x.bits)
}

func (x Float128WithRound[RND]) IsNaN() bool {
	return isNaN[binary128](x.bits)
}

func (x Float128WithRound[RND]) Sign() int {
	return getSign[binary128](x.bits)
}

func (x Float128WithRound[RND]) SignBit() bool {
	return signBit[binary128](x.bits)
}

func (x Float128WithRound[RND]) Abs() Float128WithRound[RND] {
	return Float128WithRound[RND]{abs[binary128](x.bits)}
}

func (x Float128WithRound[RND]) Neg() Float128WithRound[RND] {
	return Float128WithRound[RND]{neg[binary128](x.bits)}
}

func (x Float128WithRound[RND]) CopySign(y Float128WithRound[RND]) Float128WithRound[RND] {
	return Float128WithRound[RND]{copySign[binary128](x.bits, y.bits)}
}

func (x Float128WithRound[RND]) NextUp() Float128WithRound[RND] {
	return Float128WithRound[RND]{nextUp[binary128](x.bits)}
}

func (x Float128WithRound[RND]) NextDown() Float128WithRound[RND] {
	return Float128WithRound[RND]{nextDown[binary128](x.bits)}
}

func (x Float128WithRound[RND]) Add(y Float128WithRound[RND]) Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{add[binary128](x.bits, y.bits, rnd)}
}

func (x Float128WithRound[RND]) Sub(y Float128WithRound[RND]) Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{sub[binary128](x.bits, y.bits, rnd)}
}

func (x Float128WithRound[RND]) Dim(y Float128WithRound[RND]) Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{dim[binary128](x.bits, y.bits, rnd)}
}

func (x Float128WithRound[RND]) Mul(y Float128WithRound[RND]) Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{mul[binary128](x.bits, y.bits, rnd)}
}

func (x Float128WithRound[RND]) Div(y Float128WithRound[RND]) Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{div[binary128](x.bits, y.bits, rnd)}
}

func (x Float128WithRound[RND]) Mod(y Float128WithRound[RND]) Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{mod[binary128](x.bits, y.bits, rnd)}
}

func (x Float128WithRound[RND]) ModF() (i, f Float128WithRound[RND]) {
	q, r := modf[binary128](x.bits)
	return Float128WithRound[RND]{q}, Float128WithRound[RND]{r}
}

func (x Float128WithRound[RND]) Less(y Float128WithRound[RND]) bool {
	return less[binary128](x.bits, y.bits)
}

func (x Float128WithRound[RND]) Compare(y Float128WithRound[RND]) int {
	return compare[binary128](x.bits, y.bits)
}

func (x Float128WithRound[RND]) Equal(y Float128WithRound[RND]) bool {
	order, _ := fcmp[binary128](x.bits, y.bits)
	return order == 0
}

func (x Float128WithRound[RND]) Cmp(y Float128WithRound[RND]) (order int, ordered bool) {
	return fcmp[binary128](x.bits, y.bits)
}

func (x Float128WithRound[RND]) Min(y Float128WithRound[RND]) Float128WithRound[RND] {
	return Float128WithRound[RND]{fmin[binary128](x.bits, y.bits)}
}

func (x Float128WithRound[RND]) Max(y Float128WithRound[RND]) Float128WithRound[RND] {
	return Float128WithRound[RND]{fmax[binary128](x.bits, y.bits)}
}

func (x Float128WithRound[RND]) CmpMag(y Float128WithRound[RND]) (order int, ordered bool) {
	return fcmpMag[binary128](x.bits, y.bits)
}

func (x Float128WithRound[RND]) MinMag(y Float128WithRound[RND]) Float128WithRound[RND] {
	return Float128WithRound[RND]{fminMag[binary128](x.bits, y.bits)}
}

func (x Float128WithRound[RND]) MaxMag(y Float128WithRound[RND]) Float128WithRound[RND] {
	return Float128WithRound[RND]{fmaxMag[binary128](x.bits, y.bits)}
}

func (x Float128WithRound[RND]) Round() Float128WithRound[RND] {
	return Float128WithRound[RND]{round[binary128](x.bits)}
}

func (x Float128WithRound[RND]) RoundToEven() Float128WithRound[RND] {
	return Float128WithRound[RND]{roundToEven[binary128](x.bits)}
}

func (x Float128WithRound[RND]) Floor() Float128WithRound[RND] {
	return Float128WithRound[RND]{floor[binary128](x.bits)}
}

func (x Float128WithRound[RND]) Trunc() Float128WithRound[RND] {
	return Float128WithRound[RND]{trunc[binary128](x.bits)}
}

func (x Float128WithRound[RND]) Ceil() Float128WithRound[RND] {
	return Float128WithRound[RND]{ceil[binary128](x.bits)}
}

func (x Float128WithRound[RND]) Sqrt() Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{sqrt[binary128](x.bits, rnd)}
}

func (x Float128WithRound[RND]) RSqrt() Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{rsqrt[binary128](x.bits, rnd)}
}

func (x Float128WithRound[RND]) Hypot(y Float128WithRound[RND]) Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{hypot[binary128](x.bits, y.bits, rnd)}
}

func (x Float128WithRound[RND]) Exp() Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{exp[binary128](x.bits, rnd)}
}

func (x Float128WithRound[RND]) Exp2() Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{exp2[binary128](x.bits, rnd)}
}

func (x Float128WithRound[RND]) LogB() Float128WithRound[RND] {
	return Float128WithRound[RND]{logb[binary128](x.bits)}
}

func (x Float128WithRound[RND]) ILogB() (int, bool) {
	return ilogb[binary128](x.bits)
}
