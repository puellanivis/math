package math

import (
	"math"
)

// Erf returns the error function of x.
//
// Special cases are:
//
//	Erf(+Inf) = 1
//	Erf(-Inf) = -1
//	Erf(NaN)  = NaN
func Erf[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Erf(float64(x)))
}

// ErfInv returns the inverse error function of x.
//
// Special cases are:
//
//	ErfInv(1)   = +Inf
//	ErfInv(-1)  = -Inf
//	ErfInv(x)   = NaN if x < -1 or x > 1
//	ErfInv(NaN) = NaN
func ErfInv[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Erfinv(float64(x)))
}

// Erfc returns the complementary error function of x.
//
// Special cases are:
//
//	Erfc(+Inf) = 0
//	Erfc(-Inf) = 2
//	Erfc(NaN)  = NaN
func Erfc[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Erfc(float64(x)))
}

// ErfcInv returns the inverse of [Erfc](x).
//
// Special cases are:
//
//	ErfcInv(0)   = +Inf
//	ErfcInv(2)   = -Inf
//	ErfcInv(x)   = NaN if x < 0 or x > 2
//	ErfcInv(NaN) = NaN
func ErfcInv[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Erfcinv(float64(x)))
}

// Gamma is an alias of [Γ].
func Gamma[FLOAT Float](x FLOAT) FLOAT {
	return Γ(x)
}

// Γ returns the Gamma function of x.
//
// Special cases are:
//
//	Γ(+Inf) = +Inf
//	Γ(+0)   = +Inf
//	Γ(-0)   = -Inf
//	Γ(x)    = NaN for integer x < 0
//	Γ(-Inf) = NaN
//	Γ(NaN)  = NaN
func Γ[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Gamma(float64(x)))
}

// LnGamma is an alias of [LnΓ].
func LnGamma[FLOAT Float](x FLOAT) (lngamma FLOAT, sign int) {
	return LnΓ(x)
}

// LnΓ returns the natural logarithm and sign (-1 or +1) of [Γ](x).
//
// Special cases are:
//
//	LnΓ(+Inf)     = +Inf
//	LnΓ(0)        = +Inf
//	LnΓ(-integer) = +Inf
//	LnΓ(-Inf)     = -Inf
//	LnΓ(NaN)      = NaN
func LnΓ[FLOAT Float](x FLOAT) (lnΓ FLOAT, sign int) {
	lnΓʹ, sign := math.Lgamma(float64(x))
	return FLOAT(lnΓʹ), sign
}

// J0 returns the order-zero Bessel function of the first kind.
//
// Special cases are:
//
//	J0(±Inf) = 0
//	J0(0)    = 1
//	J0(NaN)  = NaN
func J0[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.J0(float64(x)))
}

// J1 returns the order-one Bessel function of the first kind.
//
// Special cases are:
//
//	J1(±Inf) = 0
//	J1(NaN)  = NaN
func J1[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.J1(float64(x)))
}

// Jn returns the order-n Bessel function of the first kind.
//
// Special cases are:
//
//	Jn(n, ±Inf) = 0
//	Jn(n, NaN)  = NaN
func Jn[FLOAT Float](n int, x FLOAT) FLOAT {
	return FLOAT(math.Jn(n, float64(x)))
}

// Y0 returns the order-zero Bessel function of the second kind.
//
// Special cases are:
//
//	Y0(+Inf)  = 0
//	Y0(0)     = -Inf
//	Y0(x < 0) = NaN
//	Y0(NaN)   = NaN
func Y0[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Y0(float64(x)))
}

// Y1 returns the order-one Bessel function of the second kind.
//
// Special cases are:
//
//	Y1(+Inf)  = 0
//	Y1(0)     = -Inf
//	Y1(x < 0) = NaN
//	Y1(NaN)   = NaN
func Y1[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Y1(float64(x)))
}

// Yn returns the order-n Bessel function of the second kind.
//
// Special cases are:
//
//	Yn(n, +Inf)  = 0
//	Yn(n ≥ 0, 0) = -Inf
//	Yn(n < 0, 0) = +Inf if n is odd, -Inf if n is even
//	Yn(n, x < 0) = NaN
//	Yn(n, NaN)   = NaN
func Yn[FLOAT Float](n int, x FLOAT) FLOAT {
	return FLOAT(math.Yn(n, float64(x)))
}
