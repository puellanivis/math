package bits

// Uint defines a constraint on available unsigned ints supported.
type Uint interface {
	Uint16 | uint32 | Uint64 | Uint128
}

// Bits defines a generic interface of operations upon a Uint.
type Bits[U Uint] interface {
	Int(x U) int
	FromInt(int) U

	Add(x, y, carry U) (sum, carryOut U)
	Sub(x, y, borrow U) (diff, borrowOut U)
	Inc(x U) U
	Dec(x U) U

	Mul(x, y U) (hi, lo U)
	Div(hi, lo, y U) (quo, rem U)

	Not(U) U
	Or(x, y U) U
	And(x, y U) U
	Mask(x, mask U) U
	MaskInsert(x, y, mask U) U
	Xor(x, y U) U

	Rotl(x U, k int) U
	Pow2(x int) U
	Pow2m1(x int) U
	Shl(x U, k int) U
	Shr(x U, k int) U
	Lzcnt(x U) int

	Cmp(x, y U) int
	Eq(x, y U) bool
	Neq(x, y U) bool
	Lt(x, y U) bool
	Lte(x, y U) bool
	Gt(x, y U) bool
	Gte(x, y U) bool
}
