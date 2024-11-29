package floats

import (
	"github.com/puellanivis/math/bits"
)

type binary16 struct {
	bits.Bits16
}

func (binary16) width() int {
	return 16
}

func (binary16) expWidth() int {
	return 5
}

func (binary16) mantWidth() int {
	return 16 - 5 - 1 // 10
}

func (binary16) exp2OverUnder() (overflow, underflow uint16) {
	return 0x4c00, 0x4e00
}

func (binary16) expOverUnder() (overflow, underflow, nearZero uint16) {
	return 0x498c, 0x4c55, 0x0000
}

func (binary16) ln2HiLoE() (hi, lo, ln2e uint16) {
	return 0x398c, 0x0001, 0x3dc5
}

var binary16P5toP0 = []uint16{
	0x0001,
	0x801c,
	0x0456,
	0x99b0,
	0x3155,
}

func (binary16) expPN() []uint16 {
	return binary16P5toP0
}

type bfloat16 struct {
	bits.Bits16
}

func (bfloat16) width() int {
	return 16
}

func (bfloat16) expWidth() int {
	return 8
}

func (bfloat16) mantWidth() int {
	return 16 - 8 - 1 // 7
}

func (bfloat16) exp2OverUnder() (overflow, underflow uint16) {
	return 0x4300, 0x4305
}

func (bfloat16) expOverUnder() (overflow, underflow, nearZero uint16) {
	return 0x558c, 0x55c3, 0x0000
}

func (bfloat16) ln2HiLoE() (hi, lo, ln2e uint16) {
	return 0x3f31, 0x2f52, 0x3fb9
}

var bfloat16P5toP0 = []uint16{
	0x3332,
	0xb5de,
	0x388b,
	0xbb36,
	0x3e2b,
}

func (bfloat16) expPN() []uint16 {
	return bfloat16P5toP0
}
