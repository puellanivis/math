package primative

var (
	// Identity is the Identity function:
	// 	Identity(x) = x
	Identity = Name("id", Project(1, 1))

	// Increment is an alias of the Successor function:
	// 	Increment(x) = x + 1
	Increment = Successor

	// Decrement is the function that returns 0 or the value whose successor is the argument.
	// 	Decrement(x) = { 0       if x == 0,
	// 	               { x - 1   otherwise,
	Decrement = Name("pred", Recurse(
		// f(0) = 0
		// f(x+1) = x
		Constant(0),
		endExtend(Identity, 2),
	))

	// Addition returns the sum of the two arguments.
	//	Addition(x, y) = x + y
	Addition = Name("add", Exponentiate(Identity, Increment))

	// ReverseSubtraction returns the first argument subtracted from the second argument.
	// 	ReverseSubtraction(x, y) = { 0       if x > y,
	// 	                           { y - x   otherwise,
	ReverseSubtraction = Name("nsub", Exponentiate(Identity, Decrement))

	// Subtraction returns the second argument subtracted from the first argument.
	// 	Subtraction(x, y) = { 0       if y > x,
	// 	                    { x - y   otherwise,
	// NB: Since iteration is done upon the first parameter, we need swap/reverse the arguments.
	Subtraction = Name("sub", Reverse(ReverseSubtraction))

	// Difference returns the absolute value of the subtraction of the two arguments from each other.
	// 	Difference(x, y) = |y - x|
	// 	                 = Addition(Subtraction(x, y), Subtraction(y, x))
	Difference = Name("diff", Addition.Compose(Subtraction, ReverseSubtraction))

	// Multiplication returns the product of the two arguments.
	// 	Multiplication(x, y) = x * y
	Multiplication = Name("mul", Exponentiate(Zero, Addition))

	// Power returns the second argument raised to the power of the first argument.
	// 	Power(x, y) = y ^ x
	Power = Name("pow", Exponentiate(Constant(1), Multiplication))

	// Factorial returns the factorial of the argument.
	// 	Factorial(x) = x!
	Factorial = Name("fact", ProductSeries(Successor))

	// Sign returns the sign function for the argument (since all values are non-negative integers: 0 or 1)
	// 	Sign(x) = { 0   if x = 0,
	// 	          { 1   otherwise,
	Sign = Name("sgn", Recurse(Zero, Extend(Constant(1), 2)))

	// IsZero returns logically true (1) if x is zero (0), otherwise it returns logically false (0).
	// 	IsZero(x) = { 1   if x = 0,
	// 	            { 0   otherwise,
	IsZero = Name("isz", Recurse(Constant(1), Extend(Zero, 2)))

	// IsNotZero returns logically false (0) if x is zero (0), otherwise it returns logically true (1).
	// 	IsNotZero(x) = { 0   if x = 0,
	// 	               { 1   otherwise,
	IsNotZero = Sign

	// IsEqual returns logically true (1) if the arguments are the same unsigned integer,
	// otherwise it returns logically false (0).
	// 	IsEqual(x, y) = { 1   if x = y,
	// 	                { 0   otherwise,
	IsEqual = IsZero.Compose(Difference)

	// IsNotEqual returns logically false (0) if the arguments are the same unsigned integer,
	// otherwise it returns logically true (1).
	// 	IsNotEqual(x, y) = { 1   if x ≠ y,
	// 	                   { 0   otherwise,
	IsNotEqual = IsNotZero.Compose(Difference)

	// IsLessThanOrEqual returns logically true (1) if the first argument is less than or equal to the second argument.
	// otherwise it returns logically false (0).
	// 	IsLessThanOrEqual(x, y) = { 1   if x ≤ y,
	// 	                          { 0   otherwise,
	IsLessThanOrEqual = IsZero.Compose(Subtraction)

	// IsGreaterThan returns logically true (1) if the first argument is greater than the second argument.
	// otherwise it returns logically false (0).
	// 	IsGreaterThan(x, y) = { 1   if x > y,
	// 	                      { 0   otherwise,
	IsGreaterThan = IsNotZero.Compose(Subtraction)

	// IsLessThan returns logically true (1) if the first argument is less than the second argument.
	// otherwise it returns logically false (0).
	// 	IsLessThan(x, y) = { 1   if x < y,
	// 	                   { 0   otherwise,
	IsLessThan = IsNotZero.Compose(ReverseSubtraction)

	// IsGreaterThanOrEqual returns logically true (1) if the first argument is greater than or equal to the second argument.
	// otherwise it returns logically false (0).
	// 	IsGreatherThanOrEqual(x, y) = { 1   if x ≥ y,
	// 	                              { 0   otherwise,
	IsGreaterThanOrEqual = IsZero.Compose(ReverseSubtraction)

	// LogicalNot returns logically true (1) if the argument is logically false (0),
	// otherwise it returns logically false (0).
	// 	LogicalNot(x) = { 1   if x = 0,
	// 	                { 0   otherwise,
	LogicalNot = Name("not", IsZero)

	// LogicalValue returns logically false (0) if the argument logically false (0),
	// otherwise it returns logically true (1).
	// 	LogicalValue(x) = { 1   if x ≠ 0,
	// 	                  { 0   otherwise,
	LogicalValue = IsNotZero

	// LogicalAnd returns logically false (0) if either of the arguments is logically false (0),
	// otherwise it returns logically true (1).
	// 	LogicalAnd(x, y) = x ∧ y
	LogicalAnd = Name("and", PredicateArgs(Multiplication))

	// LogicalOr returns logically true (1) if either of the arguments is logically true (1),
	// otherwise it returns logically false (0).
	// 	LogicalOr(x, y) = x ∨ y
	// 	result could be 2… so, clamp to {0,1} again
	LogicalOr = Name("or", LogicalValue.Compose(PredicateArgs(Addition)))

	// LogicalXor returns logically true (1) if the arguments are different logical values,
	// otherwise it returns logically false (0).
	// 	LogicalXor(x, y) = x ⊻ y
	LogicalXor = Name("xor", PredicateArgs(Difference))

	// Remainder returns zero (0) if either argument is zero (0),
	// otherwise it returns the modulo of x and y.
	// 	Remainder(x, y) = { 0         if x = 0 ∨ y = 0,
	// 	                  { x mod y   otherwise,
	Remainder = Name("rem", IfNotZero(Project(2, 2),
		// rem(0, y) = 0
		// rem(y, 0) = 0
		// rem(x+1, y) = sgn(|rem(x, y) + 1 - y|) * (rem(x, y) + 1)
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

	// Quotient returns zero (0) if either argument is zero (0),
	// otherwise it returns the floor of the quotient of dividend x, and divisor y.
	// 	Quotient(x, y) = { 0         if x = 0 ∨ y = 0,
	// 	                 { ⌊x / y⌋   otherwise,
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
