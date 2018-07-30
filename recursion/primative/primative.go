// Package primative provides functions to be able to build primative recursive functions.
package primative

// Func defines a primative recursive function.
type Func interface {
	// Ary returns the “Arity“ of the function, i.e. how many arguments it requires in Apply.
	Ary() uint

	// Apply evalutes the function given the arguments passed.
	Apply(x ...uint) uint

	// Apply returns a composition of this function over the slice of Func `g`.
	Compose(g ...Func) Func
}
