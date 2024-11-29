package floats

import (
	"fmt"
	"math"
	"math/big"
)

var (
	MaxFloat32             = Float32{0b0_11111110_11111111111111111111111} // 3.40282346638528859811704183484516925440e+38
	SmallestNonzeroFloat32 = Float32{0b0_00000000_00000000000000000000001} // 1.401298464324817070923729583289916131280e-45
)

func Inf32(sign bool) Float32 {
	return Float32{inf[binary32](sign)}
}

func NaN32() Float32 {
	return Float32{nan[binary32]()}
}

type Float32WithRound[RND RoundingMode] struct {
	bits uint32
}

type Float32 = Float32WithRound[RoundTiesToEven]

func Float32FromBits(bits uint32) Float32 {
	return Float32{bits}
}

func Float32WithRoundFromBits[RND RoundingMode](bits uint32) Float32WithRound[RND] {
	return Float32WithRound[RND]{bits}
}

func Float32FromFloat[F ~float32 | ~float64 | *big.Float](val F) Float32 {
	return Float32WithRoundFromFloat[RoundTiesToEven](val)
}

func Float32WithRoundFromFloat[RND RoundingMode, F ~float32 | ~float64 | *big.Float](val F) Float32WithRound[RND] {
	var rnd RND

	switch v := any(val).(type) {
	case float32:
		return Float32WithRound[RND]{convert[binary32, binary32](math.Float32bits(v), rnd)}
	case float64:
		return Float32WithRound[RND]{convert[binary64, binary32](math.Float64bits(v), rnd)}
	case *big.Float:
		return Float32WithRound[RND]{fromBigFloat[binary32, uint32](v, rnd)}
	default:
		panic("impossible type passed into Float32FromFloat")
	}
}

func (x Float32WithRound[RND]) Native() float32 {
	return math.Float32frombits(x.bits)
}

func (x Float32WithRound[RND]) Format(f fmt.State, verb rune) {
	var rnd RND

	format[binary32](x.bits, f, verb, rnd)
}

func (x Float32WithRound[RND]) Float16() Float16WithRound[RND] {
	var rnd RND

	return Float16WithRound[RND]{convert[binary32, binary16](x.bits, rnd)}
}

func (x Float32WithRound[RND]) Float32() Float32WithRound[RND] {
	return x
}

func (x Float32WithRound[RND]) Float64() Float64WithRound[RND] {
	var rnd RND

	return Float64WithRound[RND]{convert[binary32, binary64](x.bits, rnd)}
}

func (x Float32WithRound[RND]) Float128() Float128WithRound[RND] {
	var rnd RND

	return Float128WithRound[RND]{convert[binary32, binary128](x.bits, rnd)}
}

func (x Float32WithRound[RND]) BFloat16() BFloat16WithRound[RND] {
	var rnd RND

	return BFloat16WithRound[RND]{convert[binary32, bfloat16](x.bits, rnd)}
}

func (x Float32WithRound[RND]) Bits() uint32 {
	return x.bits
}

func (x Float32WithRound[RND]) IsInf(sign int) bool {
	if ok := isInf[binary32](x.bits); !ok {
		return false
	}

	if sign == 0 {
		return true
	}

	return sign == getSign[binary32](x.bits)
}

func (x Float32WithRound[RND]) IsNaN() bool {
	return isNaN[binary32](x.bits)
}

func (x Float32WithRound[RND]) Sign() int {
	return getSign[binary32](x.bits)
}

func (x Float32WithRound[RND]) SignBit() bool {
	return signBit[binary32](x.bits)
}

func (x Float32WithRound[RND]) Abs() Float32WithRound[RND] {
	return Float32WithRound[RND]{abs[binary32](x.bits)}
}

func (x Float32WithRound[RND]) Neg() Float32WithRound[RND] {
	return Float32WithRound[RND]{neg[binary32](x.bits)}
}

func (x Float32WithRound[RND]) CopySign(y Float32WithRound[RND]) Float32WithRound[RND] {
	return Float32WithRound[RND]{copySign[binary32](x.bits, y.bits)}
}

func (x Float32WithRound[RND]) NextUp() Float32WithRound[RND] {
	return Float32WithRound[RND]{nextUp[binary32](x.bits)}
}

func (x Float32WithRound[RND]) NextDown() Float32WithRound[RND] {
	return Float32WithRound[RND]{nextDown[binary32](x.bits)}
}

func (x Float32WithRound[RND]) Add(y Float32WithRound[RND]) Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{add[binary32](x.bits, y.bits, rnd)}
}

func (x Float32WithRound[RND]) Sub(y Float32WithRound[RND]) Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{sub[binary32](x.bits, y.bits, rnd)}
}

func (x Float32WithRound[RND]) Dim(y Float32WithRound[RND]) Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{dim[binary32](x.bits, y.bits, rnd)}
}

func (x Float32WithRound[RND]) Mul(y Float32WithRound[RND]) Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{mul[binary32](x.bits, y.bits, rnd)}
}

func (x Float32WithRound[RND]) Div(y Float32WithRound[RND]) Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{div[binary32](x.bits, y.bits, rnd)}
}

func (x Float32WithRound[RND]) Mod(y Float32WithRound[RND]) Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{mod[binary32](x.bits, y.bits, rnd)}
}

func (x Float32WithRound[RND]) ModF() (i, f Float32WithRound[RND]) {
	q, r := modf[binary32](x.bits)
	return Float32WithRound[RND]{q}, Float32WithRound[RND]{r}
}

func (x Float32WithRound[RND]) Less(y Float32WithRound[RND]) bool {
	return less[binary32](x.bits, y.bits)
}

func (x Float32WithRound[RND]) Compare(y Float32WithRound[RND]) int {
	return compare[binary32](x.bits, y.bits)
}

func (x Float32WithRound[RND]) Equal(y Float32WithRound[RND]) bool {
	order, _ := fcmp[binary32](x.bits, y.bits)
	return order == 0
}

func (x Float32WithRound[RND]) Cmp(y Float32WithRound[RND]) (order int, ordered bool) {
	return fcmp[binary32](x.bits, y.bits)
}

func (x Float32WithRound[RND]) Min(y Float32WithRound[RND]) Float32WithRound[RND] {
	return Float32WithRound[RND]{fmin[binary32](x.bits, y.bits)}
}

func (x Float32WithRound[RND]) Max(y Float32WithRound[RND]) Float32WithRound[RND] {
	return Float32WithRound[RND]{fmax[binary32](x.bits, y.bits)}
}

func (x Float32WithRound[RND]) CmpMag(y Float32WithRound[RND]) (order int, ordered bool) {
	return fcmpMag[binary32](x.bits, y.bits)
}

func (x Float32WithRound[RND]) MinMag(y Float32WithRound[RND]) Float32WithRound[RND] {
	return Float32WithRound[RND]{fminMag[binary32](x.bits, y.bits)}
}

func (x Float32WithRound[RND]) MaxMag(y Float32WithRound[RND]) Float32WithRound[RND] {
	return Float32WithRound[RND]{fmaxMag[binary32](x.bits, y.bits)}
}

func (x Float32WithRound[RND]) Round() Float32WithRound[RND] {
	return Float32WithRound[RND]{round[binary32](x.bits)}
}

func (x Float32WithRound[RND]) RoundToEven() Float32WithRound[RND] {
	return Float32WithRound[RND]{roundToEven[binary32](x.bits)}
}

func (x Float32WithRound[RND]) Floor() Float32WithRound[RND] {
	return Float32WithRound[RND]{floor[binary32](x.bits)}
}

func (x Float32WithRound[RND]) Trunc() Float32WithRound[RND] {
	return Float32WithRound[RND]{trunc[binary32](x.bits)}
}

func (x Float32WithRound[RND]) Ceil() Float32WithRound[RND] {
	return Float32WithRound[RND]{ceil[binary32](x.bits)}
}

func (x Float32WithRound[RND]) Sqrt() Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{sqrt[binary32](x.bits, rnd)}
}

func (x Float32WithRound[RND]) RSqrt() Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{rsqrt[binary32](x.bits, rnd)}
}

func (x Float32WithRound[RND]) Hypot(y Float32WithRound[RND]) Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{hypot[binary32](x.bits, y.bits, rnd)}
}

func (x Float32WithRound[RND]) Exp() Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{exp[binary32](x.bits, rnd)}
}

func (x Float32WithRound[RND]) Exp2() Float32WithRound[RND] {
	var rnd RND

	return Float32WithRound[RND]{exp2[binary32](x.bits, rnd)}
}

func (x Float32WithRound[RND]) LogB() Float32WithRound[RND] {
	return Float32WithRound[RND]{logb[binary32](x.bits)}
}

func (x Float32WithRound[RND]) ILogB() (int, bool) {
	return ilogb[binary32](x.bits)
}
