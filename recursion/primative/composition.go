package primative

import (
	"fmt"
	"strings"
)

type composition struct {
	k, m uint

	f Func
	g []Func
}

func (c *composition) Compose(g ...Func) Func {
	return Compose(c, g...)
}

func (c *composition) Apply(x ...uint) uint {
	if uint(len(x)) != c.m {
		panic(fmt.Sprintf("%s is %d-ary, was given %d inputs", c, c.m, len(x)))
	}

	xPrime := make([]uint, c.k)
	for i := range c.g {
		xPrime[i] = c.g[i].Apply(x...)
	}

	return c.f.Apply(xPrime...)
}

func (c *composition) Ary() uint {
	return c.m
}

func (c *composition) String() string {
	var args []string

	for i := range c.g {
		args = append(args, fmt.Sprint(c.g[i]))
	}

	return fmt.Sprintf("%s(%s)", c.f, strings.Join(args, ", "))
}

func (c *composition) Verbose(prec int) string {
	if prec == 0 {
		return c.String()
	}

	if prec < 0 {
		var args []string
		for _, g := range c.g {
			args = append(args, fmt.Sprintf("%+v", g))
		}

		return fmt.Sprintf("%+v(%s)", c.f, strings.Join(args, ", "))
	}

	var args []string
	for _, g := range c.g {
		args = append(args, fmt.Sprintf("%.*v", prec-1, g))
	}

	return fmt.Sprintf("%.*v(%s)", prec-1, c.f, strings.Join(args, ", "))
}

func (c *composition) Format(s fmt.State, verb rune) {
	format(c, s, verb)
}

// Compose returns the composition of f and a slice of Func g.
//
// Given a k-ary primitive recursive function f,
// and k m-ary primitive recursive functions g₁,…,gₖ,
// the composition of f with g₁,…,gₖ is the m-ary function h(x₁,…,xₘ)=f(g₁(x₁,…,xₘ),…,gₖ(x₁,…,xₘ))
//
// This function will panic if:
// * `len(g) != f.Ary()`
// * There exists an i such that `g[i].Ary() != g[0].Ary()`
func Compose(f Func, g ...Func) Func {
	k := f.Ary()

	if uint(len(g)) != k {
		panic(fmt.Sprintf("Composition of %d-ary f needs %d g functions, got %d", k, k, len(g)))
	}

	if k < 1 {
		// Identity composition
		return f
	}

	m := g[0].Ary()
	for i := range g[1:] {
		if g[i].Ary() != m {
			panic(fmt.Sprintf("Composition of g functions must all be %d-ary functions", m))
		}
	}

	return &composition{
		k: k,
		m: m,
		f: f,
		g: g,
	}
}
