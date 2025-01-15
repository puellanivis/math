package bits

import (
	"cmp"
	"math/bits"
)

// Uint32 is an alias to uint32, to provide name consistency with Uint128.
type Uint32 = uint32

// Bits32 provides a genericable surface to abstract uint32 math.
type Bits32 struct{}

// Int converts to an int.
func (Bits32) Int(x uint32) int {
	return int(x)
}

// FromInt converts from an int.
func (Bits32) FromInt(i int) uint32 {
	return uint32(i)
}

// Add is [bits.Add32].
func (Bits32) Add(x, y, carry uint32) (sum, carryOut uint32) {
	return bits.Add32(x, y, carry)
}

// Sub is [bits.Sub32].
func (Bits32) Sub(x, y, borrow uint32) (diff, borrowOut uint32) {
	return bits.Sub32(x, y, borrow)
}

// Inc is a simplified increment.
func (Bits32) Inc(x uint32) uint32 {
	return x + 1
}

// Dec is a simplified decrement.
func (Bits32) Dec(x uint32) uint32 {
	return x - 1
}

// Mul is [bits.Mul32].
func (Bits32) Mul(x, y uint32) (hi, lo uint32) {
	return bits.Mul32(x, y)
}

// Div is [bits.Div32].
func (Bits32) Div(hi, lo, y uint32) (quo, rem uint32) {
	return bits.Div32(hi, lo, y)
}

// Not returns the bitwise inverse of all bits in the argument.
func (Bits32) Not(x uint32) uint32 {
	return ^x
}

// Or returns the bitwise OR of the arguments.
func (Bits32) Or(x, y uint32) uint32 {
	return x | y
}

// And returns the bitwise AND of the arguments.
func (Bits32) And(x, y uint32) uint32 {
	return x & y
}

// Mask masks out the mask bits from x.
func (Bits32) Mask(x, mask uint32) uint32 {
	return x &^ mask
}

// MaskInsert composes masking out the mask bits from x with ORing in the mask bits of y.
func (Bits32) MaskInsert(x, y, mask uint32) uint32 {
	return (x &^ mask) | (y & mask)
}

// Xor returns the bitwise XOR.
func (Bits32) Xor(x, y uint32) uint32 {
	return x ^ y
}

// Rotl is [bits.RotateLeft32].
func (Bits32) Rotl(x uint32, k int) uint32 {
	return bits.RotateLeft32(x, k)
}

// Pow2 returns the integer power of two.
//
// It is undefined behavior to use an x greater than or equal to the bit width.
func (Bits32) Pow2(x int) uint32 {
	return 1 << x
}

// Pow2m1 returns the integer power of two minus one.
//
// It is undefined behavior to use an x greater than the bit width.
func (Bits32) Pow2m1(x int) uint32 {
	return (1 << x) - 1
}

// Shl performs a left shift.
func (Bits32) Shl(x uint32, k int) uint32 {
	return x << k
}

// Shr performs a right shift.
func (Bits32) Shr(x uint32, k int) uint32 {
	return x >> k
}

// Lzcnt is [bits.LeadingZeros32].
func (Bits32) Lzcnt(x uint32) int {
	return bits.LeadingZeros32(x)
}

// Cmp is [cmp.Compare].
func (Bits32) Cmp(x, y uint32) int {
	return cmp.Compare(x, y)
}

// Eq returns true if x equals y.
func (Bits32) Eq(x, y uint32) bool {
	return x == y
}

// IsZero returns true if x is zero.
func (b Bits32) IsZero(x uint32) bool {
	return b.Eq(x, 0)
}

// Neq returns true if x does not equal y.
func (Bits32) Neq(x, y uint32) bool {
	return x != y
}

// Lt returns true if x is less than y.
func (Bits32) Lt(x, y uint32) bool {
	return x < y
}

// Lte returns true if x is less than or equal to y.
func (Bits32) Lte(x, y uint32) bool {
	return x <= y
}

// Gt returns true if x is greater than y.
func (Bits32) Gt(x, y uint32) bool {
	return x > y
}

// Gte returns true if x is greater than or equal to y.
func (Bits32) Gte(x, y uint32) bool {
	return x >= y
}
