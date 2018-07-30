package primative

import (
	"fmt"
	"strings"
)

type constant uint

// Zero is the 0-ary function that returns 0.
const Zero constant = 0

// Constant returns a 0-ary function that returns the given value c.
func Constant(c uint) Func {
	return Name(fmt.Sprintf("%d", c), constant(c))
}

func (c constant) Apply(x ...uint) uint {
	if len(x) != 0 {
		panic(fmt.Sprintf("Constant is 0-ary, was given %d inputs", len(x)))
	}

	return uint(c)
}

func (c constant) Ary() uint {
	return 0
}

// String returns a constant in the form of calls of successor upon the Zero 0-ary function.
func (c constant) String() string {
	return fmt.Sprintf("%s0%s", strings.Repeat("S(", int(c)), strings.Repeat(")", int(c)))
}

func (c constant) Format(s fmt.State, verb rune) {
	format(c, s, verb)
}

func (c constant) Compose(g ...Func) Func {
	return Compose(c, g...)
}

func (c constant) extend(n uint) Func {
	return Name(fmt.Sprintf("C%s%s", superscript.Transform(n), subscript.Transform(uint(c))), &extendedConst{
		n: n,
		c: c,
	})
}

type extendedConst struct {
	n uint
	c constant
}

func (c *extendedConst) Apply(x ...uint) uint {
	if uint(len(x)) != c.n {
		panic(fmt.Sprintf("function is %d-ary, was given %d inputs", c.n, len(x)))
	}

	return uint(c.c)
}

func (c *extendedConst) Ary() uint {
	return c.n
}

// String returns the extended constant as a form of successive recursions selecting the first value,
// with the end-base-case as the successor construction of the constant.
func (c *extendedConst) String() string {
	s := c.c.String()

	if c.n == 1 {
		return fmt.Sprintf("R(%s, P²₁)", s)
	}

	return fmt.Sprintf("R(%s, P²₁)(P%s₁)", s, superscript.Transform(c.n))
}

func (c *extendedConst) Format(s fmt.State, verb rune) {
	format(c, s, verb)
}

func (c *extendedConst) Compose(g ...Func) Func {
	return Compose(c, g...)
}

func (c *extendedConst) extend(n uint) Func {
	return Name(fmt.Sprintf("C%s%s", superscript.Transform(n), subscript.Transform(uint(c.c))), &extendedConst{
		n: n,
		c: c.c,
	})
}
