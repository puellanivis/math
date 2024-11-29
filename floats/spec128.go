package floats

import (
	"github.com/puellanivis/math/bits"
)

type binary128 struct {
	bits.Bits128
}

func (binary128) width() int {
	return 128
}

func (binary128) expWidth() int {
	return 15
}

func (binary128) mantWidth() int {
	return 128 - 15 - 1 // 112
}

var (
	exp2Overflow  = bits.Uint128{Hi: 0x400d000000000000, Lo: 0}
	exp2Underflow = bits.Uint128{Hi: 0x400d01b800000000, Lo: 0}
)

func (binary128) exp2OverUnder() (overflow, underflow bits.Uint128) {
	return exp2Overflow, exp2Underflow
}

var (
	expOverflow  = bits.Uint128{Hi: 0x400c62e42fefa39e, Lo: 0xf35793c7673007e6}
	expUnderflow = bits.Uint128{Hi: 0x400c654bb3b2c73e, Lo: 0xbb059fabb506ff34}
	expNearZero  = bits.Uint128{Hi: 0x3fe3000000000000, Lo: 0x0000000000000000}
)

func (binary128) expOverUnder() (overflow, underflow, nearZero bits.Uint128) {
	return expOverflow, expUnderflow, expNearZero
}

var (
	b128ln2hi = bits.Uint128{Hi: 0x3ffe62e42fee0000, Lo: 0x0000000000000000}
	b128ln2lo = bits.Uint128{Hi: 0x3fdea39ef35793c7, Lo: 0x6000000000000000}
)

func (binary128) ln2HiLoE() (hi, lo, ln2e bits.Uint128) {
	return b128ln2hi, b128ln2lo, Ln2E.bits
}

var binary128P5toP0 = []bits.Uint128{
	bits.Uint128{Hi: 0x3fe66376972bea4d, Lo: 0x0000000000000000},
	bits.Uint128{Hi: 0xbfebbbd41c5d26bf, Lo: 0x1000000000000000},
	bits.Uint128{Hi: 0x3ff11566aaf25de2, Lo: 0xc000000000000000},
	bits.Uint128{Hi: 0xbff66c16c16bebd9, Lo: 0x3000000000000000},
	bits.Uint128{Hi: 0x3ffc555555555555, Lo: 0x5000000000000000},
}

func (binary128) expPN() []bits.Uint128 {
	return binary128P5toP0
}
