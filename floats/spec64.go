package floats

import (
	"github.com/puellanivis/math/bits"
)

type binary64 struct {
	bits.Bits64
}

func (binary64) width() int {
	return 64
}

func (binary64) expWidth() int {
	return 11
}

func (binary64) mantWidth() int {
	return 64 - 11 - 1 // 52
}

func (binary64) exp2OverUnder() (overflow, underflow uint64) {
	return 0x4090000000000000, 0x4090c80000000000
}

func (binary64) expOverUnder() (overflow, underflow, nearZero uint64) {
	return 0x40862e42fefa39ef, 0xc0874910d52d3052, 0x3e30000000000000
}

func (binary64) ln2HiLoE() (hi, lo, ln2e uint64) {
	return 0x3fe62e42fee00000, 0x3dea39ef35793c76, 0x3ff71547652b82fe
}

var binary64P5toP0 = []uint64{
	0x3E663769_72BEA4D0, // 4.13813679705723846039e-08
	0xBEBBBD41_C5D26BF1, // -1.65339022054652515390e-06
	0x3F11566A_AF25DE2C, // 6.61375632143793436117e-05
	0xBF66C16C_16BEBD93, // -2.77777777770155933842e-03
	0x3FC55555_55555555, // 1.66666666666666657415e-01
}

func (binary64) expPN() []uint64 {
	return binary64P5toP0
}
