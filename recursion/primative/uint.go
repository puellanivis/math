package primative

import (
	"math"
)

type uintType interface {
	Clone() uintType
	Increment() uintType
	Decrement() uintType
}

type bigUint []uint

func (u bigUint) Clone() uintType {
	e := 0
	if u.allOnes() {
		// if we are all ones, then extend the clone by one
		// this gives better chances that we need not extend
		// such a clone later in an Increment.
		e = 1
	}

	n := make(bigUint, len(u)+e)
	copy(n, u)
	return n
}

func (u bigUint) allOnes() bool {
	test := ^uint(0)

	for i := range u {
		test &= u[i]
	}

	return test == ^uint(0)
}

func (u bigUint) allZeros() bool {
	test := uint(0)

	for i := range u {
		test |= u[i]
	}

	return test == 0
}

func (u bigUint) Increment() uintType {
	for i := range u {
		u[i]++

		if u[i] != 0 {
			return u
		}
	}

	// if we get here, all of the previous digits carried.
	return append(u, 1)
}

func (u bigUint) Decrement() uintType {
	if u.allZeros() {
		return simpleUint(0)
	}

	for i := range u {
		u[i]--

		if u[i] != ^uint(0) {
			return u
		}
	}

	panic("decremented an all zero natural number")
}

type simpleUint uint

func (u simpleUint) Clone() uintType {
	return u
}

func (u simpleUint) Increment() uintType {
	u++

	if u == 0 {
		return bigUint{0, 1}
	}

	return u
}

func (u simpleUint) Decrement() uintType {
	if u == 0 {
		return u
	}

	return u - 1
}

func uintFromFloat(k float64) uintType {
	if k == 0 {
		return simpleUint(0)
	}

	l := math.Ceil(math.Log2(k) / 64)
	if l <= 1 {
		return uintValue(uint64(k))
	}

	panic("float is too large")
}

func uintValue(k interface{}) uintType {
	switch k := k.(type) {
	case int8:
		if k >= 0 {
			return simpleUint(k)
		}
	case int16:
		if k >= 0 {
			return simpleUint(k)
		}
	case int32:
		if k >= 0 {
			return simpleUint(k)
		}
	case int64:
		if k >= 0 {
			return uintValue(uint64(k))
		}

	case uint8:
		return simpleUint(k)
	case uint16:
		return simpleUint(k)
	case uint32:
		return simpleUint(k)
	case uint64:
		w := simpleUint(k)
		if uint64(w) == k {
			return w
		}
		return bigUint{uint(k), uint(k >> 32)}

	case float32:
		w := math.Trunc(float64(k))
		if w == float64(k) {
			return uintFromFloat(w)
		}
	case float64:
		w := math.Trunc(k)
		if w == k {
			return uintFromFloat(w)
		}

	}

	panic("invalid natural number")
}
