// Package math mostly provides genericized wrappers around the standard math library.
// Some functions have different behaviors from the math library to be more compliant with the IEEE-754 standard.
//
// This package does not guarantee bit-identical results across architectures.
package math

import (
	"math"
)

// Float is a constraint that permits any builtin floating point type.
// If future releases of Go add new floating point types, this constraint will be modified to include them.
type Float interface{ float32 | float64 }

func cast2[OUT, IN1, IN2 Float](x IN1, y IN2) (xʹ, yʹ OUT) {
	return OUT(x), OUT(y)
}

func cast3[OUT, IN1, IN2, IN3 Float](x IN1, y IN2, z IN3) (xʹ, yʹ, zʹ OUT) {
	return OUT(x), OUT(y), OUT(z)
}

// NextUp returns the next representable float32 or float64 value after x towards positive infinity.
//
// Special cases are:
//
//	NextUp(NaN)  = NaN
//	NextUp(+Inf) = +Inf
func NextUp[FLOAT Float](x FLOAT) FLOAT {
	return NextAfter(x, Inf[FLOAT](0))
}

// NextDown returns the next representable float32 or float64 value after x towards negative infinity.
//
// Special cases are:
//
//	NextUp(NaN)  = NaN
//	NextUp(-Inf) = -Inf
func NextDown[FLOAT Float](x FLOAT) FLOAT {
	return NextAfter(x, Inf[FLOAT](-1))
}

// NextAfter returns the next representable float32 or float64 value after x towards y.
//
// Special cases are:
//
//	NextAfter(x, x)   = x
//	NextAfter(NaN, y) = NaN
//	NextAfter(x, NaN) = NaN
func NextAfter[FLOAT Float](x, y FLOAT) FLOAT {
	switch any(x).(type) {
	case float32:
		return FLOAT(math.Nextafter32(cast2[float32](x, y)))
	case float64:
		return FLOAT(math.Nextafter(cast2[float64](x, y)))
	default:
		panic("impossible type")
	}
}

// Max returns the larger of x or y.
//
// Special cases are:
//
//	Max(x, NaN)  = Max(NaN, x)  = x
//	Max(x, +Inf) = Max(+Inf, x) = +Inf
//	Max(+0, ±0)  = Max(±0, +0)  = +0
//	Max(-0, -0)  = -0
//
// Note that this differs from the built-in function max and the standard library math.Max when called with NaN.
// This change is to make Max compliant with the floating-point standard.
func Max[FLOAT Float](x, y FLOAT) FLOAT {
	switch {
	case IsNaN(x):
		return y
	case IsNaN(y):
		return x
	}

	return FLOAT(math.Max(cast2[float64](x, y)))
}

// Min returns the smaller of x or y.
//
// Special cases are:
//
//	Min(x, -Inf) = Min(-Inf, x) = -Inf
//	Min(x, NaN)  = Min(NaN, x)  = NaN
//	Min(-0, ±0)  = Min(±0, -0)  = -0
//
// Note that this differs from the built-in function min and the standard library math.Min when called with NaN.
// This change is to make Min compliant with the floating-point standard.
func Min[FLOAT Float](x, y FLOAT) FLOAT {
	switch {
	case IsNaN(x):
		return y
	case IsNaN(y):
		return x
	}

	return FLOAT(math.Min(cast2[float64](x, y)))
}

// Abs returns the absolute value of x.
//
// Special cases are:
//
//	Abs(±Inf) = +Inf
//	Abs(NaN)  = NaN
func Abs[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Abs(float64(x)))
}

// Sgn returns the sign function of x.
//
// Special cases are:
//
//	Sgn(NaN) = NaN
//	Sgn(±0)  = +0
func Sgn[FLOAT Float](x FLOAT) FLOAT {
	switch {
	case x != x:
		return x
	case x == 0:
		return 0
	case x < 0:
		return -1
	default:
		return 1
	}
}

// Dim returns the maximum of x-y or 0.
//
// Special cases are:
//
//	Dim(+Inf, +Inf) = NaN
//	Dim(-Inf, -Inf) = NaN
//	Dim(x, NaN)     = Dim(NaN, x) = NaN
func Dim[FLOAT Float](x, y FLOAT) FLOAT {
	// The special cases result in NaN after the subtraction:
	//      +Inf - +Inf = NaN
	//      -Inf - -Inf = NaN
	//       NaN - y    = NaN
	//         x - NaN  = NaN
	v := x - y
	if v <= 0 {
		// v is negative or 0
		return 0
	}
	// v is positive or NaN
	return v
}

// FMA returns x * y + z, computed with only one rounding.
// (That is, FMA returns the fused multiply-add of x, y, and z.)
func FMA[FLOAT Float](x, y, z FLOAT) FLOAT {
	return FLOAT(math.FMA(cast3[float64](x, y, z)))
}

// Mod returns the floating-point remainder of x/y.
// The magnitude of the result is less than y and its sign agrees with that of x.
//
// Special cases are:
//
//	Mod(±Inf, y) = NaN
//	Mod(NaN, y)  = NaN
//	Mod(x, 0)    = NaN
//	Mod(x, ±Inf) = x
//	Mod(x, NaN)  = NaN
func Mod[FLOAT Float](x, y FLOAT) FLOAT {
	return FLOAT(math.Mod(cast2[float64](x, y)))
}

// Remainder returns the IEEE 754 floating-point remainder of x/y.
//
// Special cases are:
//
//	Remainder(±Inf, y) = NaN
//	Remainder(NaN, y)  = NaN
//	Remainder(x, 0)    = NaN
//	Remainder(x, ±Inf) = x
//	Remainder(x, NaN)  = NaN
func Remainder[FLOAT Float](x, y FLOAT) FLOAT {
	return FLOAT(math.Remainder(cast2[float64](x, y)))
}

// Sqrt returns the square root of x.
//
// Special cases are:
//
//	Sqrt(+Inf)  = +Inf
//	Sqrt(±0)    = ±0
//	Sqrt(x < 0) = NaN
//	Sqrt(NaN)   = NaN
func Sqrt[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Sqrt(float64(x)))
}

// Cbrt returns the cube root of x.
//
// Special cases are:
//
//	Cbrt(±0)   = ±0
//	Cbrt(±Inf) = ±Inf
//	Cbrt(NaN)  = NaN
func Cbrt[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Cbrt(float64(x)))
}

// Exp returns e**2, the base-e exponential of x.
//
// Special cases are:
//
//	Exp(+Inf) = +Inf
//	Exp(NaN)  = NaN
//
// Very large values overflow to 0 or +Inf.
// Very small values underflow to 1.
func Exp[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Exp(float64(x)))
}

// ExpM1 returns e**x - 1, the base-e exponential of x minus 1.
// It is more accurate than [Exp](x) - 1 when x is near zero.
//
// Special cases are:
//
//	ExpM1(+Inf) = +Inf
//	ExpM1(-Inf) = -1
//	ExpM1(NaN)  = NaN
//
// Very large values overflow to -1 or +Inf.
func ExpM1[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Expm1(float64(x)))
}

// Exp2 returns 2**x, the base-2 exponential of x.
//
// Special cases are the same as [Exp].
func Exp2[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Exp2(float64(x)))
}

// Pow returns x**y, the base-x exponential of y.
//
// Special cases are (in order):
//
//	Pow(x, ±0)    = 1 for any x
//	Pow(1, y)     = 1 for any y
//	Pow(x, 1)     = x for any x
//	Pow(NaN, y)   = NaN
//	Pow(x, NaN)   = NaN
//	Pow(±0, y)    = ±Inf for y an odd integer < 0
//	Pow(±0, -Inf) = +Inf
//	Pow(±0, +Inf) = +0
//	Pow(±0, y)    = +Inf for finite y < 0 and not an odd integer
//	Pow(±0, y)    = ±0 for y an odd integer > 0
//	Pow(±0, y)    = +0 for finite y > 0 and not an odd integer
//	Pow(-1, ±Inf) = 1
//	Pow(x, +Inf)  = +Inf for |x| > 1
//	Pow(x, -Inf)  = +0 for |x| > 1
//	Pow(x, +Inf)  = +0 for |x| < 1
//	Pow(x, -Inf)  = +Inf for |x| < 1
//	Pow(+Inf, y)  = +Inf for y > 0
//	Pow(+Inf, y)  = +0 for y < 0
//	Pow(-Inf, y)  = Pow(-0, -y)
//	Pow(x, y)     = NaN for finite x < 0 and finite non-integer y
func Pow[FLOAT Float](x, y FLOAT) FLOAT {
	return FLOAT(math.Pow(cast2[float64](x, y)))
}

// Pow10 returns 10**n, the base-10 exponential of n.
//
// Special cases are:
//
//	Pow10(n) =    0 for n < -323
//	Pow10(n) = +Inf for n > 308
func Pow10[FLOAT Float](x int) FLOAT {
	return FLOAT(math.Pow10(x))
}

// ILogB returns the binary exponent of x as an integer.
//
// Special cases are:
//
//	ILogB(±Inf) = MaxInt32
//	ILogB(0)    = MinInt32
//	ILogB(NaN)  = MaxInt32
func ILogB[FLOAT Float](x FLOAT) int {
	return math.Ilogb(float64(x))
}

// LogB returns the binary exponent of x.
//
// Special cases are:
//
//	Logb(±Inf) = +Inf
//	Logb(0)    = -Inf
//	Logb(NaN)  = NaN
func LogB[FLOAT Float](x FLOAT) FLOAT {
	// special cases
	switch {
	case x == 0:
		return Inf[FLOAT](-1)
	case IsInf(x, 0):
		return Inf[FLOAT](1)
	case IsNaN(x):
		return x
	}
	return FLOAT(ILogB(x))
}

// Log returns the natural logarithm of x.
//
// Special cases are:
//
//	Log(+Inf)  = +Inf
//	Log(0)     = -Inf
//	Log(x < 0) = NaN
//	Log(NaN)   = NaN
func Log[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Log(float64(x)))
}

// Ln is an alias of [Log].
func Ln[FLOAT Float](x FLOAT) FLOAT {
	return Log(x)
}

// Log1p returns the natural logarithm of 1 plus its argument x.
// It is more accurate than [Log](1 + x) when x is near zero.
//
// Special cases are:
//
//	Log1p(+Inf)   = +Inf
//	Log1p(±0)     = ±0
//	Log1p(-1)     = -Inf
//	Log1p(x < -1) = NaN
//	Log1p(NaN)    = NaN
func Log1p[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Log1p(float64(x)))
}

// Ln1p is an alias of [Log1p].
func Ln1p[FLOAT Float](x FLOAT) FLOAT {
	return Log1p(x)
}

// Log2 returns the binary logarithm of x.
// The special cases are the same as for [Log].
func Log2[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Log2(float64(x)))
}

// Lb is an alias of [Log2].
func Lb[FLOAT Float](x FLOAT) FLOAT {
	return Log2(x)
}

// Log10 returns the decimal logarithm of x.
// The special cases are the same as for [Log].
func Log10[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Log10(float64(x)))
}
