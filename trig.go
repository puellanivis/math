package math

import (
	"math"
)

// Acos returns the arccosine, in radians, of x.
//
// Special case is:
//
//	Acos(x) = NaN if x < -1 or x > 1
func Acos[FLOAT Float](x FLOAT) (θ FLOAT) {
	return FLOAT(math.Acos(float64(x)))
}

// Acosh returns the inverse hyperbolic cosine of x.
//
// Special cases are:
//
//	Acosh(+Inf) = +Inf
//	Acosh(x)    = NaN if x < 1
//	Acosh(NaN)  = NaN
func Acosh[FLOAT Float](x FLOAT) (θ FLOAT) {
	return FLOAT(math.Acosh(float64(x)))
}

// Asin returns the arcsine, in radians, of x.
//
// Special cases are:
//
//	Asin(±0) = ±0
//	Asin(x) = NaN if x < -1 or x > 1
func Asin[FLOAT Float](x FLOAT) (θ FLOAT) {
	return FLOAT(math.Asin(float64(x)))
}

// Asinh returns the inverse hyperbolic sine of x.
//
// Special cases are:
//
//	Asinh(±0)   = ±0
//	Asinh(±Inf) = ±Inf
//	Asinh(NaN)  = NaN
func Asinh[FLOAT Float](x FLOAT) (θ FLOAT) {
	return FLOAT(math.Asinh(float64(x)))
}

// Atan returns the arctangent, in radians, of x.
//
// Special cases are:
//
//	Atan(±0) = ±0
//	Atan(±Inf) = ±Pi/2
func Atan[FLOAT Float](x FLOAT) (θ FLOAT) {
	return FLOAT(math.Atan(float64(x)))
}

// Atanh returns the inverse hyperbolic tangent of x.
//
// Special cases are:
//
//	Atanh(1)  = +Inf
//	Atanh(±0) = ±0
//	Atanh(-1) = -Inf
//	Atanh(x)  = NaN if x < -1 or x > 1
func Atanh[FLOAT Float](x FLOAT) (θ FLOAT) {
	return FLOAT(math.Atanh(float64(x)))
}

// Cos returns the cosine of the radian argument θ.
//
// Special cases are:
//
//	Cos(±Inf) = NaN
//	Cos(NaN)  = NaN
func Cos[FLOAT Float](θ FLOAT) (x FLOAT) {
	return FLOAT(math.Cos(float64(θ)))
}

// Cosh returns the hyperbolic cosine of x.
//
// Special cases are:
//
//	Cosh(±0)   = 1
//	Cosh(±Inf) = +Inf
//	Cosh(NaN)  = NaN
func Cosh[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Cosh(float64(x)))
}

// Sin returns the sine of the radian argument θ.
//
// Special cases are:
//
//	Sin(±0)   = ±0
//	Sin(±Inf) = NaN
//	Sin(NaN)  = NaN
func Sin[FLOAT Float](θ FLOAT) (x FLOAT) {
	return FLOAT(math.Sin(float64(θ)))
}

// Sinh returns the hyperbolic sine of x.
//
// Special cases are:
//
//	Sinh(±0)   = ±0
//	Sinh(±Inf) = ±Inf
//	Sinh(NaN)  = NaN
func Sinh[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Sinh(float64(x)))
}

// Tan returns the tangent of the radian argument θ.
//
// Special cases are:
//
//	Tan(±0)   = ±0
//	Tan(±Inf) = NaN
//	Tan(NaN)  = NaN
func Tan[FLOAT Float](θ FLOAT) (x FLOAT) {
	return FLOAT(math.Tan(float64(θ)))
}

// Tanh returns the hyperbolic tangent of x.
//
// Special cases are:
//
//	Tanh(±0) = ±0
//	Tanh(±Inf) = ±1
//	Tanh(NaN) = NaN
func Tanh[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Tanh(float64(x)))
}

// Atan2 returns the arc tangent of y/x, using
// the signs of the two to determine the quadrant
// of the return value.
//
// Special cases are (in order):
//
//	Atan2(y, NaN)     = NaN
//	Atan2(NaN, x)     = NaN
//	Atan2(+0, x>=0)   = +0
//	Atan2(-0, x>=0)   = -0
//	Atan2(+0, x<=-0)  = +τ/2
//	Atan2(-0, x<=-0)  = -τ/2
//	Atan2(y>0, 0)     = +τ/4
//	Atan2(y<0, 0)     = -τ/4
//	Atan2(+Inf, +Inf) = +τ/8
//	Atan2(-Inf, +Inf) = -τ/8
//	Atan2(+Inf, -Inf) = 3τ/8
//	Atan2(-Inf, -Inf) = -3τ/8
//	Atan2(y, +Inf)    = 0
//	Atan2(y>0, -Inf)  = +τ/2
//	Atan2(y<0, -Inf)  = -τ/2
//	Atan2(+Inf, x)    = +τ/4
//	Atan2(-Inf, x)    = -τ/4
func Atan2[FLOAT Float](y, x FLOAT) (θ FLOAT) {
	return FLOAT(math.Atan2(cast2[float64](y, x)))
}

// CosSin returns [Cos](x), [Sin](x).
//
// Special cases are:
//
//	CosSin(±0)   = 1, ±0
//	CosSin(±Inf) = NaN, NaN
//	CosSin(NaN)  = NaN, NaN
func CosSin[FLOAT Float](θ FLOAT) (cos, sin FLOAT) {
	sin, cos = cast2[FLOAT](math.Sincos(float64(θ)))
	return
}

// SinCos returns [Sin](x), [Cos](x).
//
// Special cases are:
//
//	SinCos(±0)   = ±0, 1
//	SinCos(±Inf) = NaN, NaN
//	SinCos(NaN)  = NaN, NaN
func SinCos[FLOAT Float](θ FLOAT) (sin, cos FLOAT) {
	sin, cos = cast2[FLOAT](math.Sincos(float64(θ)))
	return
}

// Hypot returns [Sqrt](p*p + q*q),
// taking care to avoid unnecessary overflow and underflow.
//
// Special cases are:
//
//	Hypot(±Inf, q) = +Inf
//	Hypot(p, ±Inf) = +Inf
//	Hypot(NaN, q)  = NaN
//	Hypot(p, NaN)  = NaN
func Hypot[FLOAT Float](p, q FLOAT) FLOAT {
	return FLOAT(math.Hypot(cast2[float64](p, q)))
}
