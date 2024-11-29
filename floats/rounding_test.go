package floats

var (
	_ RoundingMode = RoundTowardZero{}
	_ RoundingMode = RoundTowardPositive{}
	_ RoundingMode = RoundTowardNegative{}
	_ RoundingMode = RoundTiesToAway{}
	_ RoundingMode = RoundTiesToEven{}
)
