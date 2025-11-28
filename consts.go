package math

import (
	"math"
)

// Mathematical constants.
const (
	E   = math.E                                                           // 2.71828182845904523536028747135266249775724709369995957496696763 https://oeis.org/A001113
	Pi  = math.Pi                                                          // 3.14159265358979323846264338327950288419716939937510582097494459 https://oeis.org/A000796
	Phi = math.Phi                                                         // 1.61803398874989484820458683436563811772030917980576286213544862 https://oeis.org/A001622
	Tau = 6.28318530717958647692528676655900576839433879875021164194988918 // https://oeis.org/A019692

	Sqrt2   = math.Sqrt2                                                       // 1.41421356237309504880168872420969807856967187537694807317667974 https://oeis.org/A002193
	SqrtE   = math.SqrtE                                                       // 1.64872127070012814684865078781416357165377610071014801157507931 https://oeis.org/A019774
	SqrtPi  = math.SqrtPi                                                      // 1.77245385090551602729816748334114518279754945612238712821380779 https://oeis.org/A002161
	SqrtPhi = math.SqrtPhi                                                     // 1.27201964951406896425242246173749149171560804184009624861664038 https://oeis.org/A139339
	SqrtTau = 2.50662827463100050241576528481104525300698674060993831662992358 // https://oeis.org/A019727

	Ln2    = math.Ln2    // 0.693147180559945309417232121458176568075500134360255254120680009 https://oeis.org/A002162
	Log2E  = math.Log2E  // 1 / Ln2
	Ln10   = math.Ln10   // 2.30258509299404568401799145468436420760110148862877297603332790 https://oeis.org/A002392
	Log10E = math.Log10E // 1 / Ln10
)

// Floating-point limit values.
// Max is the largest finite value representable by the type.
// SmallestNormal is the smallest positive, non-subnormal value representable by the type.
// SmallestNonzero is the smallest positive, non-zero value representable by the type.
const (
	MaxFloat32             = math.MaxFloat32             // 3.40282346638528859811704183484516925440e+38
	SmallestNormalFloat32  = 0x1p-126                    // 1.175494350822287507968736537222245677819e-38
	SmallestNonzeroFloat32 = math.SmallestNonzeroFloat32 // 1.401298464324817070923729583289916131280e-45

	MaxFloat64             = math.MaxFloat64             // 1.79769313486231570814527423731704356798070e+308
	SmallestNormalFloat64  = 0x1p-1022                   // 2.2250738585072013830902327173324040642192160e-308
	SmallestNonzeroFloat64 = math.SmallestNonzeroFloat64 // 4.9406564584124654417656879286822137236505980e-324
)

// Integer limit values.
const (
	MaxInt   = math.MaxInt   // MaxInt32 or MaxInt64 depending on intSize.
	MinInt   = math.MinInt   // MinInt32 or MinInt64 depending on intSize.
	MaxInt8  = math.MaxInt8  // 127
	MinInt8  = math.MinInt8  // 128
	MaxInt16 = math.MaxInt16 // 32767
	MinInt16 = math.MinInt16 // -32768
	MaxInt32 = math.MaxInt32 // 2147483647
	MinInt32 = math.MinInt32 // -2147483648
	MaxInt64 = math.MaxInt64 // 9223372036854775807
	MinInt64 = math.MinInt64 // -9223372036854775808

	MaxUint   = math.MaxUint   // MaxUint32 or MaxUint64 depending on intSize.
	MaxUint8  = math.MaxUint8  // 255
	MaxUint16 = math.MaxUint16 // 65535
	MaxUint32 = math.MaxUint32 // 4294967295
	MaxUint64 = math.MaxUint64 // 18446744073709551615
)
