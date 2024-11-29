package bits

import (
	"cmp"
	"fmt"
	"math/big"
	"math/bits"
	"strconv"
)

const wordsPerUint128 = 128 / strconv.IntSize

// Uint128 defines a 128 bit unsigned integer.
type Uint128 struct {
	Hi, Lo uint64
}

func (u *Uint128) set(words []big.Word) {
	if strconv.IntSize == 32 {
		var w [4]big.Word
		copy(w[:], words)

		u.Hi = (uint64(w[3]) << 32) | uint64(w[2])
		u.Lo = (uint64(w[1]) << 32) | uint64(w[0])
	}

	var w [2]big.Word
	copy(w[:], words)

	u.Hi, u.Lo = uint64(w[1]), uint64(w[0])
}

func (u Uint128) words() []big.Word {
	if strconv.IntSize == 32 {
		const mask = (1 << 32) - 1
		return []big.Word{
			big.Word(u.Lo & mask),
			big.Word(u.Lo >> 32),
			big.Word(u.Hi & mask),
			big.Word(u.Hi >> 32),
		}
	}

	return []big.Word{
		big.Word(u.Lo),
		big.Word(u.Hi),
	}
}

// Format implements [fmt.Formatter].
func (u Uint128) Format(f fmt.State, verb rune) {
	new(big.Int).SetBits(u.words()).Format(f, verb)
}

// Bits128 provides a genericable surface to abstract uint64 math.
type Bits128 struct{}

// Int converts to an int.
//
// This is a simple cast, and will discard bits above strconv.IntSize.
func (Bits128) Int(x Uint128) int {
	return int(x.Lo)
}

// FromInt converts from an int.
func (Bits128) FromInt(i int) (u Uint128) {
	return Uint128{Lo: uint64(i)}
}

// Add is bits.Add128.
func (Bits128) Add(x, y, carry Uint128) (sum, carryOut Uint128) {
	sum.Lo, carryOut.Lo = bits.Add64(x.Lo, y.Lo, carry.Lo)
	sum.Hi, carryOut.Lo = bits.Add64(x.Hi, y.Hi, carryOut.Lo)
	return
}

// Sub is bits.Sub128.
func (Bits128) Sub(x, y, borrow Uint128) (diff, borrowOut Uint128) {
	diff.Lo, borrowOut.Lo = bits.Sub64(x.Lo, y.Lo, borrow.Lo)
	diff.Hi, borrowOut.Lo = bits.Sub64(x.Hi, y.Hi, borrowOut.Lo)
	return
}

// Inc is a simplified increment.
func (b Bits128) Inc(x Uint128) Uint128 {
	lo, carry := bits.Add64(x.Lo, 1, 0)
	hi, _ := bits.Add64(x.Hi, 0, carry)
	return Uint128{Lo: lo, Hi: hi}
}

// Dec is a simplified decrement.
func (b Bits128) Dec(x Uint128) Uint128 {
	lo, borrow := bits.Sub64(x.Lo, 1, 0)
	hi, _ := bits.Sub64(x.Hi, 0, borrow)
	return Uint128{Lo: lo, Hi: hi}
}

// Mul is bits.Mul128.
func (b Bits128) Mul(x, y Uint128) (hi, lo Uint128) {
	// TODO: implement without using math/big.

	bx := new(big.Int).SetBits(x.words())
	by := new(big.Int).SetBits(y.words())

	bprod := new(big.Int).Mul(bx, by)

	prod := bprod.Bits()

	sep := wordsPerUint128
	if sep > len(prod) {
		sep = len(prod)
	}

	lo.set(prod[:sep])
	hi.set(prod[sep:])

	return
}

// Div is bits.Div128.
func (Bits128) Div(hi, lo, y Uint128) (quo, rem Uint128) {
	// TODO: implement without using math/big.

	bz := new(big.Int).SetBits(hi.words())
	bz.Lsh(bz, 128)
	bz.Or(bz, new(big.Int).SetBits(lo.words()))

	by := new(big.Int).SetBits(y.words())

	bquo, brem := new(big.Int).QuoRem(bz, by, new(big.Int))

	quo.set(bquo.Bits())
	rem.set(brem.Bits())

	return
}

// Not returns the bitwise inverse of all bits in the argument.
func (Bits128) Not(x Uint128) Uint128 {
	return Uint128{Lo: ^x.Lo, Hi: ^x.Hi}
}

// Or returns the bitwise OR of the arguments.
func (Bits128) Or(x, y Uint128) Uint128 {
	return Uint128{Lo: x.Lo | y.Lo, Hi: x.Hi | y.Hi}
}

// And returns the bitwise AND of the arguments.
func (Bits128) And(x, y Uint128) Uint128 {
	return Uint128{Lo: x.Lo & y.Lo, Hi: x.Hi & y.Hi}
}

// Mask masks out the mask bits from x.
func (Bits128) Mask(x, mask Uint128) Uint128 {
	return Uint128{Lo: x.Lo &^ mask.Lo, Hi: x.Hi &^ mask.Hi}
}

// MaskInsert composes masking out the mask bits from x with ORing in the mask bits of y.
func (Bits128) MaskInsert(x, y, mask Uint128) Uint128 {
	return Uint128{
		Lo: (x.Lo &^ mask.Lo) | (y.Lo & mask.Lo),
		Hi: (x.Hi &^ mask.Hi) | (y.Hi & mask.Hi),
	}
}

// Xor returns the bitwise XOR of the arguments.
func (Bits128) Xor(x, y Uint128) Uint128 {
	return Uint128{Lo: x.Lo ^ y.Lo, Hi: x.Hi ^ y.Hi}
}

// Rotl is bits.RotateLeft128.
func (Bits128) Rotl(x Uint128, k int) (r Uint128) {
	c := uint(k) % 128

	if c <= 64 {
		r.Hi = (x.Hi << c) | (x.Lo >> (64 - c))
		r.Lo = (x.Lo << c) | (x.Hi >> (64 - c))
		return
	}

	c -= 64

	r.Hi = (x.Lo << c) | (x.Hi >> (64 - c))
	r.Lo = (x.Hi << c) | (x.Lo >> (64 - c))

	return
}

// Pow2 returns the integer power of two.
//
// It is undefined behavior to use an x greater than or equal to the bit width.
func (Bits128) Pow2(x int) Uint128 {
	if x >= 64 {
		return Uint128{Hi: 1 << (x - 64)}
	}
	return Uint128{Lo: 1 << x}
}

// Pow2m1 returns the integer power of two minus one.
//
// It is undefined behavior to use an x greater than the bit width.
func (Bits128) Pow2m1(x int) Uint128 {
	if x >= 64 {
		return Uint128{Hi: (1 << (x - 64)) - 1, Lo: ^uint64(0)}
	}
	return Uint128{Lo: (1 << x) - 1}
}

// Shl performs a left shift.
func (Bits128) Shl(x Uint128, k int) (r Uint128) {
	if k <= 64 {
		return Uint128{
			Hi: (x.Hi << k) | (x.Lo >> (64 - k)),
			Lo: (x.Lo << k),
		}
	}

	return Uint128{Hi: x.Lo << (k - 64)}
}

// Shr performs a right shift.
func (Bits128) Shr(x Uint128, k int) Uint128 {
	if k <= 64 {
		return Uint128{
			Hi: (x.Hi >> k),
			Lo: (x.Lo >> k) | (x.Hi << (64 - k)),
		}
	}

	return Uint128{Lo: x.Hi >> (k - 64)}
}

// Lzcnt is bits.LeadingZeros128.
func (Bits128) Lzcnt(x Uint128) int {
	lz := bits.LeadingZeros64(x.Hi)
	if lz == 64 {
		lz += bits.LeadingZeros64(x.Lo)
	}
	return lz
}

// Cmp is cmp.Compare.
func (Bits128) Cmp(x, y Uint128) int {
	if i := cmp.Compare(x.Hi, y.Hi); i != 0 {
		return i
	}
	return cmp.Compare(x.Lo, y.Lo)
}

// Eq returns true if x equals y.
func (b Bits128) Eq(x, y Uint128) bool {
	return b.Cmp(x, y) == 0
}

// Neq returns true if x does not equal y.
func (b Bits128) Neq(x, y Uint128) bool {
	return b.Cmp(x, y) != 0
}

// Lt returns true if x is less than y.
func (b Bits128) Lt(x, y Uint128) bool {
	return b.Cmp(x, y) < 0
}

// Lte returns true if x is less than or equal to y.
func (b Bits128) Lte(x, y Uint128) bool {
	return b.Cmp(x, y) <= 0
}

// Gt returns true if x is greater than y.
func (b Bits128) Gt(x, y Uint128) bool {
	return b.Cmp(x, y) > 0
}

// Gte returns true if x is greater than or equal to y.
func (b Bits128) Gte(x, y Uint128) bool {
	return b.Cmp(x, y) >= 0
}
