package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"sort"
)

func main() {
	primes := []int{
		2, 3, 5, 7, 11,
		13, 17, 19, 23, 29,
		31, 37, 41, 43, 47,
		53, 59, 61, 67, 71,
		73, 79, 83, 89, 97,
	}

	primesBig := make([]*big.Int, 0, len(primes))
	for i := 0; i < len(primes); i++ {
		primesBig = append(primesBig, big.NewInt(int64(primes[i])))
	}

	n := len(primesBig)
	fmt.Printf("%d: %2d\n", n, primesBig)

	nBig := big.NewInt(int64(n))

	tmp := new(big.Int).Add(nBig, big.NewInt(0))
	acc := new(big.Int).Set(tmp)
	for _, pᵢ := range primesBig[1:] {
		pᵢ.Mul(pᵢ, acc)
		acc.Mul(acc, tmp)
	}

	fmt.Printf("%d: %2d\n", n, primesBig)
	// test assumptions
	for i, xᵢ := range primesBig {
		tmp := new(big.Int)
		for j, xⱼ := range primesBig[i+1:] {
			if tmp.Div(xⱼ, xᵢ).Cmp(nBig) < 0 {
				panic(fmt.Sprintf("assertion failed (x(%d) = %d) / (x(%d) = %d) < %d", j+i+1, xⱼ, i, xᵢ, len(primes)))
			}
		}
	}

	// randomize the input
	x := make([]*big.Int, 1, n+1)
	x[0] = new(big.Int)
	for _, i := range rand.Perm(n) {
		x = append(x, new(big.Int).Set(primesBig[i]))
	}

	fmt.Printf("%d: %2d\n", n, x)

	// x[0] = 0 so, this is equal to p = Σ_{i=1}^{n} ixᵢ and q = Σ_{i=1}^{n} xᵢ
	p, q := new(big.Int), new(big.Int)
	for i, xᵢ := range x {
		tmp.SetInt64(int64(i))
		p.Add(p, tmp.Mul(tmp, xᵢ))
		q.Add(q, xᵢ)
	}

	z := make([]*big.Int, len(x))
	for i := n; i >= 1; i-- {
		jBig := new(big.Int).Div(p, q)

		if !jBig.IsInt64() {
			panic(fmt.Sprintf("%d is not an int64", jBig))
		}

		j := jBig.Int64()

		fmt.Println(p, "/", q, "=", j)

		fmt.Println(j, ":", x[j+1], "?>", x[j])

		if x[j+1].Cmp(x[j]) > 0 {
			j++
			jBig.SetInt64(j)
		}

		z[i] = x[j]
		fmt.Printf("z[%d] ← x[%d] = %d\n", i, j, z[i])

		p.Sub(p, jBig.Mul(jBig, x[j]))
		q.Sub(q, x[j])

		x[j] = new(big.Int) // sets to zero
	}

	fmt.Printf("%d: %t: %2d\n", n, sort.SliceIsSorted(z[1:], func(i, j int) bool { return x[i].Cmp(x[j]) < 0 }), z[1:])
}
