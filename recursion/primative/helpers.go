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

	if k == 0 {
		// If we’re extending a 0-ary function, then it doesn’t matter what order we work in.
		// So, just use the function we’ve already developed.
		return Extend(f, n)
	}

	if f, ok := f.(*named); ok {
		if p, ok := f.Func.(*projection); ok {
			return endExtend(p, n)
		}
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

	args := make([]Func, k)
	for i := range args {
		args[i] = Project(n, (n-k)+uint(i)+1)
	}

	return f.Compose(args...)
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

	type extender interface {
		extend(uint) Func
	}

	// Special case: some functions might have an extend function on them. In this case, use that instead.
	if f, ok := f.(extender); ok {
		return f.extend(n)
	}

	if k == 0 {
		// We cannot compose a 0-ary function to anything but itself.
		// However, we can recurse on a single value such that base case is the 0-ary function,
		// and each recursion returns the results of the last.
		// At this point, we now have a 1-ary function.
		// And THAT function, _can_ be composed via projections into an n-ary function.
		return Extend(Recurse(f, Project(2, 1)), n)
	}

	args := make([]Func, k)
	for i := range args {
		args[i] = Project(n, uint(i)+1)
	}

	return f.Compose(args...)
}

// PredicateArgs returns a k-ary function that composes f upon the predicate values of each argument.
// 	g(x₁, …, x_n) = f( sgn(x₁), …, sgn(x_n) )
func PredicateArgs(f Func) Func {
	n := f.Ary()

	args := make([]Func, n)

	for i := range args {
		args[i] = Sign.Compose(Project(n, uint(i+1)))
	}

	return f.Compose(args...)
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

// Exponentiate returns a an exponentiation of a function and a given base-case.
//
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

// IfNotZero returns a function that will be `then` when `cond` is logically true.
func IfNotZero(cond Func, then Func) Func {
	return Multiplication.Compose(LogicalValue.Compose(cond), then)
}

// IfZero returns a function that is `then` when `cond` is logicall not true.
func IfZero(cond Func, then Func) Func {
	return Multiplication.Compose(LogicalNot.Compose(cond), then)
}

// IfThenElse returns a function that will be `then` when `cond` is logically true,
// and `otherwise` when cond is logically false.
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

	return Recurse(
		Extend(Zero, k-1),
		Addition.Compose(
			Project(k+1, 1),
			endExtend(f, k+1),
		),
	)
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
