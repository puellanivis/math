package primative

import (
	"fmt"
)

type successor struct{}

func (s successor) Compose(g ...Func) Func {
	return Compose(s, g...)
}

func (s successor) Apply(x ...uint) uint {
	if len(x) != 1 {
		panic(fmt.Sprintf("Successor is 1-ary, was given %d inputs", len(x)))
	}

	if x[0] == ^uint(0) {
		panic("mathematical overflow")
	}

	return x[0] + 1
}

func (s successor) Ary() uint {
	return 1
}

func (s successor) String() string {
	return "S"
}

func (S successor) Format(s fmt.State, verb rune) {
	format(S, s, verb)
}

var Successor successor

func SuccessorOf(f Func) Func {
	return Compose(Successor, f)
}
