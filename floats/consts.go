package floats

import (
	"github.com/puellanivis/math/bits"
)

// Constants as per from Wolfram|Alpha:
// e         = 2¹  × 1.5bf0a8b14576_95355fb8ac404e7a_7
// φ         = 2⁰  × 1.9e3779b97f4a_7c15f39cc0605ced_e
// π         = 2¹  × 1.921fb54442d1_8469898cc51701b8_3
// τ         = 2²  × 1.921fb54442d1_8469898cc51701b8_3
// √2        = 2⁰  × 1.6a09e667f3bc_c908b2fb1366ea95_7d
// √e        = 2⁰  × 1.a61298e1e069_bc972dfefab6df33_f
// √φ        = 2⁰  × 1.45a3146a8845_5e92554501121ec4_d
// √π        = 2⁰  × 1.c5bf891b4ef6_aa79c3b0520d5db9_3
// √τ        = 2¹  × 1.40d931ff6270_59657ca41fae722c_e
// ln 2      = 2⁻¹ × 1.62e42fefa39e_f35793c7673007e5_e
// (ln 2)⁻¹  = 2⁰  × 1.71547652b82f_e1777d0ffda0d23a_7d
// ln 10     = 2¹  × 1.26bb1bbb5551_582dd4adac5705a6_1
// (ln 10)⁻¹ = 2⁻² × 1.bcb7b1526e50_e32a6ab7555f5a67_b

// Various mathematical constants, encoded in Float128 for maximum precision.
var (
	E   = Float128{bits.Uint128{Hi: 0x40005bf0a8b14576, Lo: 0x95355fb8ac404e7a}}
	Phi = Float128{bits.Uint128{Hi: 0x3fff9e3779b97f4a, Lo: 0x7c15f39cc0605cee}}
	Pi  = Float128{bits.Uint128{Hi: 0x4000921fb54442d1, Lo: 0x8469898cc51701b8}}
	Tau = Float128{bits.Uint128{Hi: 0x4001921fb54442d1, Lo: 0x8469898cc51701b8}}

	Sqrt2   = Float128{bits.Uint128{Hi: 0x3fff6a09e667f3bc, Lo: 0xc908b2fb1366ea95}}
	SqrtE   = Float128{bits.Uint128{Hi: 0x3fffa61298e1e069, Lo: 0xbc972dfefab6df34}}
	SqrtPhi = Float128{bits.Uint128{Hi: 0x3fff45a3146a8845, Lo: 0x5e92554501121ec5}}
	SqrtPi  = Float128{bits.Uint128{Hi: 0x3fffc5bf891b4ef6, Lo: 0xaa79c3b0520d5db9}}
	SqrtTau = Float128{bits.Uint128{Hi: 0x400040d931ff6270, Lo: 0x59657ca41fae722d}}

	Ln2   = Float128{bits.Uint128{Hi: 0x3ffe62e42fefa39e, Lo: 0xf35793c7673007e6}}
	Ln2E  = Float128{bits.Uint128{Hi: 0x3fff71547652b82f, Lo: 0xe1777d0ffda0d23a}}
	Ln10  = Float128{bits.Uint128{Hi: 0x400026bb1bbb5551, Lo: 0x582dd4adac5705a6}}
	Ln10E = Float128{bits.Uint128{Hi: 0x3ffdbcb7b1526e50, Lo: 0xe32a6ab7555f5a68}}
)
