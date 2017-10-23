package primative

import (
	"fmt"
)

type projection struct {
	n, i uint
}

func Project(n, i uint) Func {
	if n == 0 {
		panic("projection cannot be 0-ary")
	}

	if i < 1 || i > n {
		panic(fmt.Sprintf("projection requires 1 ≤ (i = %d) ≤ (n = %d)", i, n))
	}

	return &projection{
		n: n,
		i: i - 1, // precalculate base-1 to base-0 indexing
	}
}

func (p *projection) Apply(x ...uint) uint {
	if uint(len(x)) != p.n {
		panic(fmt.Sprintf("%s is %d-ary, was given %d inputs", p, p.n, len(x)))
	}

	return x[p.i]
}

func (p *projection) Ary() uint {
	return p.n
}

func (p *projection) String() string {
	// since we pre-calculated base-1 to base-0 indexing, we undo that here.
	i := p.i + 1

	return fmt.Sprintf("P%s%s", superscript.Transform(p.n), subscript.Transform(i))
}

func (p *projection) Format(s fmt.State, verb rune) {
	format(p, s, verb)
}

func (p *projection) Compose(g ...Func) Func {
	return Compose(p, g...)
}

func (p *projection) extend(n uint) Func {
	return &projection{
		n: n,
		i: p.i,
	}
}
