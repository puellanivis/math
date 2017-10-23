// Package primative provides functions to be able to build primative recursive functions.
package primative

import ()

type Func interface {
	Ary() uint
	Apply(x ...uint) uint
	Compose(g ...Func) Func
}
