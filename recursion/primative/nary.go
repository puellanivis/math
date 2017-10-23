package primative

import (
	"fmt"
)

type funcN struct {
	ary uint
	f func(...uint) uint
	s string
}

func (f *funcN) Compose(g ...Func) Func {
	return Compose(f, g...)
}

func (f *funcN) Apply(x ...uint) uint {
	if f.ary != uint(len(x)) {
		panic(fmt.Sprintf("function is %d-ary, was given %d inputs", len(x)))
	}

	return f.f(x...)
}

func (f *funcN) Ary() uint {
	return f.ary
}

func (f *funcN) String() string {
	return f.s
}
