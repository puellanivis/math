package floats

import (
	"github.com/puellanivis/math/bits"
)

type binary32 struct {
	bits.Bits32
}

func (binary32) width() int {
	return 32
}

func (binary32) expWidth() int {
	return 8
}

func (binary32) mantWidth() int {
	return 32 - 8 - 1 // 23
}

func (binary32) exp2OverUnder() (overflow, underflow uint32) {
	return 0x43000000, 0x43150000
}

func (binary32) expOverUnder() (overflow, underflow, nearZero uint32) {
	return 0x42b17218, 0x42cff1b4, 0x31800000
}

func (binary32) ln2HiLoE() (hi, lo, ln2e uint32) {
	return 0x3f317218, 0x2f51cf7a, 0x3fb8aa3b
}

var binary32P5toP0 = []uint32{
	0x3331bb4c,
	0xb5ddea0e,
	0x388ab355,
	0xbb360b61,
	0x3e2aaaab,
}

func (binary32) expPN() []uint32 {
	return binary32P5toP0
}
