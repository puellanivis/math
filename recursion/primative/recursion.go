package primative

import (
	"fmt"
)

type recursion struct {
	ary uint
	f   Func
	g   Func
}

func (r *recursion) Compose(g ...Func) Func {
	return Compose(r, g...)
}

func (r *recursion) Apply(x ...uint) uint {
	if uint(len(x)) != r.ary {
		panic(fmt.Sprintf("%s is %d-ary, was given %d inputs", r, r.ary, len(x)))
	}

	if x[0] == 0 {
		// h(0, x₁, …, x_k) = f(x₁, …, x_k)
		return r.f.Apply(x[1:]...)
	}

	// build "stack" frame…
	x_ := make([]uint, len(x)+1)

	// OPTIMIZATION
	// if the recursive function is a constant, then short-circuit any recursion, and argument assignment.
	if g, ok := r.g.(*extendedConst); ok {
		return g.Apply(x_...)
	}

	// x′ = { _, x₀-1, x₁, …, x_k }
	x_[1] = x[0] - 1
	copy(x_[2:], x[1:])

	// OPTIMIZATION
	// if the recursive function is a projection with i > 1, then we can short-circuit any recursion.
	if p, ok := r.g.(*projection); ok && p.i > 0 {
		return p.Apply(x_...)
	}

	// x′₀ = h(x₀-1, x₁, …, x_k)
	x_[0] = r.Apply(x_[1:]...)

	// h(x₀, x₁, …, x_k) = g(h(x₀-1, x₁, …, x_k), x₀-1, x₁, …, x_k)
	return r.g.Apply(x_...)
}

func (r *recursion) Ary() uint {
	return r.ary
}

func (r *recursion) String() string {
	return fmt.Sprintf("R(%s, %s)", r.f, r.g)
}

func (r *recursion) Verbose(prec int) string {
	if prec == 0 {
		return r.String()
	}

	if prec < 1 {
		return fmt.Sprintf("R(%+v, %+v)", r.f, r.g)
	}

	return fmt.Sprintf("R(%.*v, %.*v)", prec, r.f, prec, r.g)
}

func (r *recursion) Format(s fmt.State, verb rune) {
	format(r, s, verb)
}

func Recurse(f, g Func) Func {
	k := f.Ary()

	if g.Ary() != k+2 {
		panic(fmt.Sprintf("Primative Recursion of %d-ary f, must have %d-ary g, got %d-ary", k, k+2, g.Ary()))
	}

	return &recursion{
		ary: k + 1,
		f:   f,
		g:   g,
	}
}
