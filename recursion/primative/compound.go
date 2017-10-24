package primative

import ()

var (
	// Identity(x) = x
	Identity = Name("id", Project(1, 1))

	// Increment(x) = x + 1
	Increment = Successor

	// Decrement(x) = { 0       if x == 0,
	//                { x - 1   otherwise,
	Decrement = Name("pred", Recurse(
		// f(0) = 0
		// f(x+1) = x
		Constant(0),
		endExtend(Identity, 2),
	))

	// Addition(x, y) = x + y
	Addition = Name("add", Exponentiate(Identity, Increment))

	// ReverseSubtraction(x, y) = { 0       if x > y,
	//                            { y - x   otherwise,
	ReverseSubtraction = Name("nsub", Exponentiate(Identity, Decrement))

	// Subtraction(x, y) = { 0       if y > x,
	//                     { x - y   otherwise,
	// NB: Since iteration is done upon the first parameter, we need swap/reverse the arguments.
	Subtraction = Name("sub", Reverse(ReverseSubtraction))

	// Difference(x, y) = |y - x|
	//                  = Addition(Subtraction(x, y), Subtraction(y, x))
	Difference = Name("diff", Addition.Compose(Subtraction, ReverseSubtraction))

	// Multiplication(x, y) = x * y
	Multiplication = Name("mul", Exponentiate(Zero, Addition))

	// Power(x, y) = y ^ x
	Power = Name("pow", Exponentiate(Constant(1), Multiplication))

	// Factorial(x) = x!
	Factorial = Name("fact", ProductSeries(SuccessorOf(Identity)))

	Sign = Name("sgn", Recurse(Zero, Extend(Constant(1), 2)))

	IsZero = Name("isz", Recurse(Constant(1), Extend(Zero, 2)))
	IsNotZero = Sign

	IsEqual = IsZero.Compose(Difference)
	IsNotEqual = IsNotZero.Compose(Difference)

	IsLessThanOrEqual = IsZero.Compose(Subtraction)
	IsGreaterThan = IsZero.Compose(Subtraction)

	IsLessThan = IsNotZero.Compose(ReverseSubtraction)
	IsGreaterThanOrEqual = IsZero.Compose(ReverseSubtraction)

	LogicalNot = Name("not", IsZero)
	LogicalValue = IsNotZero

	LogicalAnd = Multiplication.Compose(
			LogicalValue.Compose(Project(2, 1)),
			LogicalValue.Compose(Project(2, 2)),
	)
	LogicalOr  = LogicalValue.Compose( // result could be 2… so, clamp to {0,1} again
			Addition.Compose(
				LogicalValue.Compose(Project(2, 1)),
				LogicalValue.Compose(Project(2, 2)),
			),
	)
	LogicalXor = Difference.Compose(
			LogicalValue.Compose(Project(2, 1)),
			LogicalValue.Compose(Project(2, 2)),
	)

	Remainder = Name("rem", IfNotZero(Project(2, 2),
		// rem(0, y) = 0
		// rem(y, 0) = 0
		// rem(x+1, y) = sgn(|rem(x, y) + 1 - y|) * rem(x, y)
		Recurse(
			Extend(Zero, 1),
			IfNotZero(
				Difference.Compose(
					SuccessorOf(Project(3, 1)),
					Project(3, 3),
				),
				SuccessorOf(Project(3, 1)),
			),
		),
	))

	Quotient = Name("quo", IfNotZero(Project(2, 2),
		// q(0, y) = 0
		// q(x, 0) = 0
		// q(x, y) = ∑_{n=0}^x |1 - sgn(|rem(n, y) + 1 - y|)|
		CountZeros(
			Difference.Compose(
				SuccessorOf(
					Remainder.Compose(Project(2, 1), Project(2, 2)),
				),
				Project(2, 2),
			),
		),
	))
)
