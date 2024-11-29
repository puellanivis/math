package bits

import (
	"cmp"
	"math/bits"
)

// Uint16 is an alias to uint16, to provide name consistency with Uint128.
type Uint16 = uint16

// Bits16 provides a genericable surface to abstract uint16 math.
type Bits16 struct{}

// Int converts to an int.
func (Bits16) Int(x uint16) int {
	return int(x)
}

// FromInt converts from an int.
func (Bits16) FromInt(i int) uint16 {
	return uint16(i)
}

// Add is bits.Add16.
func (Bits16) Add(x, y, carry uint16) (sum, carryOut uint16) {
	sum32 := uint32(x) + uint32(y) + uint32(carry)
	sum = uint16(sum32)
	carryOut = uint16(sum32 >> 16)
	return
}

// Sub is bits.Sub16.
func (Bits16) Sub(x, y, borrow uint16) (diff, borrowOut uint16) {
	diff = x - y - borrow
	// The difference will underflow if the top bit of x is not set and the top
	// bit of y is set (^x & y) or if they are the same (^(x ^ y)) and a borrow
	// from the lower place happens. If that borrow happens, the result will be
	// 1 - 1 - 1 = 0 - 0 - 1 = 1 (& diff).
	borrowOut = ((^x & y) | (^(x ^ y) & diff)) >> 15
	return
}

// Inc is a simplified increment.
func (Bits16) Inc(x uint16) uint16 {
	return x + 1
}

// Dec is a simplified decrement.
func (Bits16) Dec(x uint16) uint16 {
	return x - 1
}

// Mul is bits.Mul16.
func (Bits16) Mul(x, y uint16) (hi, lo uint16) {
	tmp := uint32(x) * uint32(y)
	hi, lo = uint16(tmp>>16), uint16(tmp)
	return
}

// Div is bits.Div16.
func (Bits16) Div(hi, lo, y uint16) (quo, rem uint16) {
	if y != 0 && y <= hi {
		panic("integer overflow")
	}
	z := uint32(hi)<<16 | uint32(lo)
	quo, rem = uint16(z/uint32(y)), uint16(z%uint32(y))
	return
}

// Not returns the bitwise inverse of all bits in the argument.
func (Bits16) Not(x uint16) uint16 {
	return ^x
}

// Or returns the bitwise OR of the arguments.
func (Bits16) Or(x, y uint16) uint16 {
	return x | y
}

// And returns the bitwise AND of the arguments.
func (Bits16) And(x, y uint16) uint16 {
	return x & y
}

// Mask masks out the mask bits from x.
func (Bits16) Mask(x, mask uint16) uint16 {
	return x &^ mask
}

// MaskInsert composes masking out the mask bits from x with ORing in the mask bits of y.
func (Bits16) MaskInsert(x, y, mask uint16) uint16 {
	return (x &^ mask) | (y & mask)
}

// Xor returns the bitwise XOR.
func (Bits16) Xor(x, y uint16) uint16 {
	return x ^ y
}

// Rotl is [bits.RotateLeft16].
func (Bits16) Rotl(x uint16, k int) uint16 {
	return bits.RotateLeft16(x, k)
}

// Pow2 returns the integer power of two.
//
// It is undefined behavior to use an x greater than or equal to the bit width.
func (Bits16) Pow2(x int) uint16 {
	return 1 << x
}

// Pow2m1 returns the integer power of two minus one.
//
// It is undefined behavior to use an x greater than the bit width.
func (Bits16) Pow2m1(x int) uint16 {
	return (1 << x) - 1
}

// Shl performs a left shift.
func (Bits16) Shl(x uint16, k int) uint16 {
	return x << k
}

// Shr performs a right shift.
func (Bits16) Shr(x uint16, k int) uint16 {
	return x >> k
}

// Lzcnt is [bits.LeadingZeros16].
func (Bits16) Lzcnt(x uint16) int {
	return bits.LeadingZeros16(x)
}

// Cmp is [cmp.Compare].
func (Bits16) Cmp(x, y uint16) int {
	return cmp.Compare(x, y)
}

// Eq returns true if x equals y.
func (Bits16) Eq(x, y uint16) bool {
	return x == y
}

// Neq returns true if x does not equal y.
func (Bits16) Neq(x, y uint16) bool {
	return x != y
}

// Lt returns true if x is less than y.
func (Bits16) Lt(x, y uint16) bool {
	return x < y
}

// Lte returns true if x is less than or equal to y.
func (Bits16) Lte(x, y uint16) bool {
	return x <= y
}

// Gt returns true if x is greater than y.
func (Bits16) Gt(x, y uint16) bool {
	return x > y
}

// Gte returns true if x is greater than or equal to y.
func (Bits16) Gte(x, y uint16) bool {
	return x >= y
}
