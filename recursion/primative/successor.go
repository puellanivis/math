package primative

import (
	"fmt"
)

type successor struct{}

func (S successor) Compose(g ...Func) Func {
	return Compose(S, g...)
}

func (S successor) Apply(x ...uint) uint {
	if len(x) != 1 {
		panic(fmt.Sprintf("Successor is 1-ary, was given %d inputs", len(x)))
	}

	if x[0] == ^uint(0) {
		panic("mathematical overflow")
	}

	return x[0] + 1
}

func (S successor) Ary() uint {
	return 1
}

func (S successor) String() string {
	return "S"
}

func (S successor) Format(s fmt.State, verb rune) {
	format(S, s, verb)
}

// Successor is the successor function.
var Successor successor

// SuccessorOf returns the successor of the given function.
// 	S(k) = k + 1
func SuccessorOf(f Func) Func {
	return Compose(Successor, f)
}
