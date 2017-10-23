package primative

import (
	"fmt"
)

// SetTest returns a 1-ary function that returns true if the argument is in the set S.
//	f_{S}(x) = { 1   if x ∈ S,
//	           { 0   otherwise,
func SetTest(S ...uint) Func {
	if len(S) == 0 {
		return Extend(Zero, 1)
	}

	isEqTo := IsEqual.Compose(
		Extend(Constant(S[0]), 1),
		Identity,
	)

	if len(S) == 1 {
		return isEqTo
	}

	return Addition.Compose(
		isEqTo,
		SetTest(S[1:]...),
	)
}

// endExtend returns an n-ary function that returns the value of the k-ary f function.
// 	g(x₁, …, x_n) = f(x_{n-k}, …, x_n), where k ≤ n
func endExtend(f Func, n uint) Func {
	k := f.Ary()
	if n < k {
		panic(fmt.Sprintf("%s is %d-ary, cannot extend to %d", f, k, n))
	}

	if n == k {
		return f
	}

	// Special case: extending a projection is as simple as increasing the n of the Projection.
	// endExtending requires also incrementing the i by the difference of the increase.
	if p, ok := f.(*projection); ok {
		// we built it ourselves to avoid having to make a 0-base to 1-base conversion.
		return &projection{
			n: n,
			i: p.i + (n - k),
		}
	}

	// such meta, very recursive, wow
	return Recurse(
		endExtend(f, n-1),
		Project(n+1, 1),
	)
}

// Extend returns an n-ary function that returns the value of the k-ary f function.
// 	g(x₁, …, x_n) = f(x₁, …, x_k), where k ≤ n
func Extend(f Func, n uint) Func {
	k := f.Ary()
	if n < k {
		panic(fmt.Sprintf("%s is %d-ary, cannot extend to %d", f, k, n))
	}

	if n == k {
		return f
	}

	type extender interface{
		extend(uint) Func
	}

	// Special case: some functions might have an extend function on them. In this case, use that instead.
	if f, ok := f.(extender); ok {
		return f.extend(n)
	}

	if k == 0 {
		// in the case of a 0-ary function, no reordering of parameters is necessary.
		// the underlying function won’t use any of them anyways.
		return endExtend(f, n)
	}

	// as endExtend trims arguments from the start,
	// we have to reverse the arguments to start, and then reverse them again before applying f.
	return Reverse(endExtend(Reverse(f), n))
}

// Reverse returns a k-ary function that composes f upon a reverse of the arguments given.
// 	g(x₁, …, x_n) = f(x_n, …, x₁)
func Reverse(f Func) Func {
	n := f.Ary()

	// special case: if f is 1-ary, or 0-ary, there there is no reversal necessary
	if n <= 1 {
		return f
	}

	args := make([]Func, n)

	for i := range args {
		args[i] = Project(n, n-uint(i))
	}

	return f.Compose(args...)
}

// Eponentiate(b:ℕ¹→ℕ¹, f:ℕ^k→ℕ¹) := g(x₀, …, x_k) = { b(x₁, …, x_k)                       when x₀ = 0,
//                                                   { f(g(x-1, x₁, …, x_n), x₁, …, x_n)   otherwise,
func Exponentiate(base, f Func) Func {
	k := f.Ary()

	if k == base.Ary() {
		k++
	}

	args := make([]Func, f.Ary())
	for i := range args {
		args[i] = Project(k+1, uint(i)+2)
	}
	args[0] = Project(k+1, 1)

	return Recurse(
		Extend(base, k-1),
		f.Compose(args...),
	)
}

func IfNotZero(cond Func, then Func) Func {
	return Multiplication.Compose(LogicalValue.Compose(cond), then)
}

func IfZero(cond Func, then Func) Func {
	return Multiplication.Compose(LogicalNot.Compose(cond), then)
}

func IfThenElse(cond Func, then Func, otherwise Func) Func {
	return Addition.Compose(
		IfNotZero(cond, then),
		IfZero(cond, otherwise),
	)
}

// ProductSeries returns a function that multiplies a series of applications of a function.
// 	g(x₀, …, x_k) = ∏_{n=0}^x₀ f(n, x₁, …, x_k) 
func ProductSeries(f Func) Func {
	k := f.Ary()

	return Recurse(
		Extend(Constant(1), k-1),
		Multiplication.Compose(
			Project(k+1, 1),
			endExtend(f, k+1),
		),
	)
}

// Summation returns a function that sums a series of applications of a function.
// 	g(x₀, …, x_k) = ∑_{n=0}^x₀ f(n, x₁, …, x_k) 
func Summation(f Func) Func {
	k := f.Ary()

	return Name(fmt.Sprintf("sum(%s)", f), Recurse(
		Extend(Zero, k-1),
		Addition.Compose(
			Project(k+1, 1),
			endExtend(f, k+1),
		),
	))
}

// Count returns a function that sums the signs of a series of applications of a function.
// 	g(x₀, …, x_k) = ∑_{n=0}^x₀ sgn(f(n, x₁, …, x_k))
func Count(f Func) Func {
	return Summation(LogicalValue.Compose(f))
}

// CountZeros returns a function that counts the number of zero values in a series of applications of a function.
// 	g(x₀, …, x_k) = ∑_{n=0}^x₀ |1 - sgn(f(n, x₁, …, x_k))|
func CountZeros(f Func) Func {
	return Summation(LogicalNot.Compose(f))
}
