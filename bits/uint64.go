package bits

import (
	"cmp"
	"math/bits"
)

// Uint64 is an alias to uint64, to provide name consistency with Uint128.
type Uint64 = uint64

// Bits64 provides a genericable surface to abstract uint64 math.
type Bits64 struct{}

// Int converts to an int.
func (Bits64) Int(x uint64) int {
	return int(x)
}

// FromInt converts from an int.
func (Bits64) FromInt(i int) uint64 {
	return uint64(i)
}

// Add is [bits.Add64].
func (Bits64) Add(x, y, carry uint64) (sum, carryOut uint64) {
	return bits.Add64(x, y, carry)
}

// Sub is [bits.Sub64].
func (Bits64) Sub(x, y, borrow uint64) (diff, borrowOut uint64) {
	return bits.Sub64(x, y, borrow)
}

// Inc is a simplified increment.
func (Bits64) Inc(x uint64) uint64 {
	return x + 1
}

// Dec is a simplified decrement.
func (Bits64) Dec(x uint64) uint64 {
	return x - 1
}

// Mul is [bits.Mul64].
func (Bits64) Mul(x, y uint64) (hi, lo uint64) {
	return bits.Mul64(x, y)
}

// Div is [bits.Div64].
func (Bits64) Div(hi, lo, y uint64) (quo, rem uint64) {
	return bits.Div64(hi, lo, y)
}

// Not returns the bitwise inverse of all bits in the argument.
func (Bits64) Not(x uint64) uint64 {
	return ^x
}

// Or returns the bitwise OR of the arguments.
func (Bits64) Or(x, y uint64) uint64 {
	return x | y
}

// And returns the bitwise AND of the arguments.
func (Bits64) And(x, y uint64) uint64 {
	return x & y
}

// Mask masks out the mask bits from x.
func (Bits64) Mask(x, mask uint64) uint64 {
	return x &^ mask
}

// MaskInsert composes masking out the mask bits from x with ORing in the mask bits of y.
func (Bits64) MaskInsert(x, y, mask uint64) uint64 {
	return (x &^ mask) | (y & mask)
}

// Xor returns the bitwise XOR of the arguments.
func (Bits64) Xor(x, y uint64) uint64 {
	return x ^ y
}

// Rotl is [bits.RotateLeft64].
func (Bits64) Rotl(x uint64, k int) uint64 {
	return bits.RotateLeft64(x, k)
}

// Pow2 returns the integer power of two.
//
// It is undefined behavior to use an x greater than or equal to the bit width.
func (Bits64) Pow2(x int) uint64 {
	return 1 << x
}

// Pow2m1 returns the integer power of two minus one.
//
// It is undefined behavior to use an x greater than the bit width.
func (Bits64) Pow2m1(x int) uint64 {
	return (1 << x) - 1
}

// Shl performs a left shift.
func (Bits64) Shl(x uint64, k int) uint64 {
	return x << k
}

// Shr performs a right shift.
func (Bits64) Shr(x uint64, k int) uint64 {
	return x >> k
}

// Lzcnt is [bits.LeadingZeros64].
func (Bits64) Lzcnt(x uint64) int {
	return bits.LeadingZeros64(x)
}

// Cmp is [cmp.Compare].
func (Bits64) Cmp(x, y uint64) int {
	return cmp.Compare(x, y)
}

// Eq returns true if x equals y.
func (Bits64) Eq(x, y uint64) bool {
	return x == y
}

// Neq returns true if x does not equal y.
func (Bits64) Neq(x, y uint64) bool {
	return x != y
}

// Lt returns true if x is less than y.
func (Bits64) Lt(x, y uint64) bool {
	return x < y
}

// Lte returns true if x is less than or equal to y.
func (Bits64) Lte(x, y uint64) bool {
	return x <= y
}

// Gt returns true if x is greater than y.
func (Bits64) Gt(x, y uint64) bool {
	return x > y
}

// Gte returns true if x is greater than or equal to y.
func (Bits64) Gte(x, y uint64) bool {
	return x >= y
}
