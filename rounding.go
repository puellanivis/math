package math

import (
	"math"
)

// Floor returns the greatest integer value less than or equal to x.
//
// Special cases are:
//
//	Floor(±0)   = ±0
//	Floor(±Inf) = ±Inf
//	Floor(NaN)  = NaN
func Floor[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Floor(float64(x)))
}

// Ceil returns the least integer value greater than or equal to x.
//
// Special cases are:
//
//	Ceil(±0)   = ±0
//	Ceil(±Inf) = ±Inf
//	Ceil(NaN)  = NaN
func Ceil[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Ceil(float64(x)))
}

// Trunc returns the integer value of x.
//
// Special cases are:
//
//	Trunc(±0)   = ±0
//	Trunc(±Inf) = ±Inf
//	Trunc(NaN)  = NaN
func Trunc[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Trunc(float64(x)))
}

// Modf returns integer and fractional floating-point numbers that sum to f.
// Both values have the same sign as f.
//
// Special cases are:
//
//	Modf(±Inf) = ±Inf, NaN
//	Modf(NaN)  = NaN, NaN
func Modf[FLOAT Float](x FLOAT) (int, frac FLOAT) {
	return cast2[FLOAT](math.Modf(float64(x)))
}

// Round returns the nearest integer, rounding half away from zero.
//
// Special cases are:
//
//	Round(±0)   = ±0
//	Round(±Inf) = ±Inf
//	Round(NaN)  = NaN
func Round[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.Round(float64(x)))
}

// RoundToEven returns the nearest integer, rounding ties to even.
//
// Special cases are:
//
//	RoundToEven(±0)   = ±0
//	RoundToEven(±Inf) = ±Inf
//	RoundToEven(NaN)  = NaN
func RoundToEven[FLOAT Float](x FLOAT) FLOAT {
	return FLOAT(math.RoundToEven(float64(x)))
}
