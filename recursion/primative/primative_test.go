package primative

import (
	"testing"
)

func TestBasics(t *testing.T) {
	C := Constant(1)

	if k := C.Apply(); k != 1 {
		t.Errorf("%.1v(): expected 1, got %d", C, k)
	}

	S := Successor

	if k := S.Apply(1); k != 2 {
		t.Errorf("%.1v(1): expected 2, got %d", S, k)
	}

	h := Compose(S, Compose(S, C))
	if k := h.Apply(); k != 3 {
		t.Errorf("%.1v(): expected 3, got %d", h, k)
	}
}

func TestAddition(t *testing.T) {
	add := Addition

	if k := add.Apply(0, 42); k != 42 {
		t.Errorf("%.1v(0, 42): expected 42, got %d", add, k)
	}

	if k := add.Apply(1, 1); k != 2 {
		t.Errorf("%.1v(1, 1): expected 2, got %d", add, k)
	}

	if k := add.Apply(2, 1); k != 3 {
		t.Errorf("%.1v(2, 1): expected 3, got %d", add, k)
	}

	if k := add.Apply(1, 2); k != 3 {
		t.Errorf("%.1v(1, 2): expected 3, got %d", add, k)
	}

	if k := add.Apply(2, 2); k != 4 {
		t.Errorf("%.1v(2, 2): expected 4, got %d", add, k)
	}
}

func TestMultiplication(t *testing.T) {
	mul := Multiplication

	if k := mul.Apply(0, 42); k != 0 {
		t.Errorf("%.1v(0, 42): expected 0, got %d", mul, k)
	}

	if k := mul.Apply(1, 42); k != 42 {
		t.Errorf("%.1v(1, 42): expected 42, got %d", mul, k)
	}

	if k := mul.Apply(2, 42); k != 84 {
		t.Errorf("%.1v(2, 42): expected 84, got %d", mul, k)
	}

	if k := mul.Apply(6, 7); k != 42 {
		t.Errorf("%.1v(6, 7): expected 42, got %d", mul, k)
	}
}

func TestPower(t *testing.T) {
	pow := Power

	if k := pow.Apply(0, 42); k != 1 {
		t.Errorf("%.1v(0, 42): expected 1, got %d", pow, k)
	}

	if k := pow.Apply(1, 42); k != 42 {
		t.Errorf("%.1v(1, 42): expected 42, got %d", pow, k)
	}

	if k := pow.Apply(2, 42); k != 42*42 {
		t.Errorf("%.1v(2, 42): expected %d, got %d", pow, 42*42, k)
	}
}

func TestFactorial(t *testing.T) {
	fact := Factorial

	if k := fact.Apply(0); k != 1 {
		t.Errorf("%.1v(0): expected 1, got %d", fact, k)
	}

	if k := fact.Apply(1); k != 1 {
		t.Errorf("%.1v(1): expected 1, got %d", fact, k)
	}

	if k := fact.Apply(2); k != 2 {
		t.Errorf("%.1v(2): expected 2, got %d", fact, k)
	}

	if k := fact.Apply(3); k != 6 {
		t.Errorf("%.1v(3): expected 6, got %d", fact, k)
	}

	if k := fact.Apply(4); k != 24 {
		t.Errorf("%.1v(4): expected 24, got %d", fact, k)
	}
}

func TestEqualities(t *testing.T) {
	is0 := IsZero

	if k := is0.Apply(0); k != 1 {
		t.Errorf("%.1v(0): expected 1, got %d", is0, k)
	}

	if k := is0.Apply(1); k != 0 {
		t.Errorf("%.1v(1): expected 0, got %d", is0, k)
	}

	if k := is0.Apply(42); k != 0 {
		t.Errorf("%.1v(42): expected 0, got %d", is0, k)
	}

	isEq := IsEqual

	if k := isEq.Apply(0, 0); k != 1 {
		t.Errorf("%.1v(0, 0): expected 1, got %d", isEq, k)
	}

	if k := isEq.Apply(1, 1); k != 1 {
		t.Errorf("%.1v(1, 1): expected 1, got %d", isEq, k)
	}

	if k := isEq.Apply(42, 42); k != 1 {
		t.Errorf("%.1v(42, 42): expected 1, got %d", isEq, k)
	}

	if k := isEq.Apply(42, 1); k != 0 {
		t.Errorf("%.1v(42, 1): expected 0, got %d", isEq, k)
	}

	if k := isEq.Apply(1, 42); k != 0 {
		t.Errorf("%.1v(1, 42): expected 0, got %d", isEq, k)
	}

	if k := isEq.Apply(42, 1); k != 0 {
		t.Errorf("%.1v(42, 1): expected 0, got %d", isEq, k)
	}

	if k := isEq.Apply(1, 42); k != 0 {
		t.Errorf("%.1v(1, 42): expected 0, got %d", isEq, k)
	}

	isLt := IsLessThan

	if k := isLt.Apply(1, 42); k != 1 {
		t.Errorf("%.1v(1, 42): expected 1, got %d", isLt, k)
	}
	if k := isLt.Apply(42, 42); k != 0 {
		t.Errorf("%.1v(42, 42): expected 0, got %d", isLt, k)
	}

	inSet := SetTest(2, 3, 5, 7, 11, 13)

	if k := inSet.Apply(1); k != 0 {
		t.Errorf("%.1v(1): expected 0, got %d", inSet, k)
	}

	if k := inSet.Apply(2); k != 1 {
		t.Errorf("%.1v(2): expected 1, got %d", inSet, k)
	}

	if k := inSet.Apply(13); k != 1 {
		t.Errorf("%.1v(13): expected 1, got %d", inSet, k)
	}

	if k := inSet.Apply(42); k != 0 {
		t.Errorf("%.1v(42): expected 0, got %d", inSet, k)
	}
}

func TestDivision(t *testing.T) {
	rem := Remainder

	if k := rem.Apply(6, 7); k != 6 {
		t.Errorf("%.1v(6, 7): expected 6, got %d", rem, k)
	}

	if k := rem.Apply(8, 7); k != 1 {
		t.Errorf("%.1v(8, 7): expected 1, got %d", rem, k)
	}

	if k := rem.Apply(42, 7); k != 0 {
		t.Errorf("%.1v(42, 7): expected 0, got %d", rem, k)
	}

	if k := rem.Apply(0, 7); k != 0 {
		t.Errorf("%.1v(0, 7): expected 0, got %d", rem, k)
	}

	if k := rem.Apply(7, 0); k != 0 {
		t.Errorf("%.1v(7, 0): expected 0, got %d", rem, k)
	}

	div := Quotient

	if k := div.Apply(0, 42); k != 0 {
		t.Errorf("%.1v(0, 42): expected 0, got %d", div, k)
	}

	if k := div.Apply(2, 42); k != 0 { // got 2
		t.Errorf("%.1v(2, 42): expected 0, got %d", div, k)
	}

	if k := div.Apply(43, 42); k != 1 { // got 41?
		t.Errorf("%.1v(43, 42): expected 1, got %d", div, k)
	}

	if k := div.Apply(42, 7); k != 6 {
		t.Errorf("%.1v(42, 7): expected 6, got %d", div, k)
	}
}

func TestSubtraction(t *testing.T) {
	pred := Decrement

	if k := pred.Apply(0); k != 0 {
		t.Errorf("%.1v(0): expected 0, got %d", pred, k)
	}

	if k := pred.Apply(1); k != 0 {
		t.Errorf("%.1v(1): expected 0, got %d", pred, k)
	}

	if k := pred.Apply(2); k != 1 {
		t.Errorf("%.1v(2): expected 1, got %d", pred, k)
	}

	if k := pred.Apply(42); k != 41 {
		t.Errorf("%.1v(42): expected 41, got %d", pred, k)
	}

	sub := Subtraction

	if k := sub.Apply(2, 2); k != 0 {
		t.Errorf("%.1v(2, 2): expected 0, got %d", sub, k)
	}

	if k := sub.Apply(2, 1); k != 1 {
		t.Errorf("%.1v(2, 1): expected 1, got %d", sub, k)
	}

	if k := sub.Apply(2, 0); k != 2 {
		t.Errorf("%.1v(2, 0): expected 2, got %d", sub, k)
	}

	if k := sub.Apply(42, 12); k != 30 {
		t.Errorf("%.1v(42, 12): expected 30, got %d", sub, k)
	}

	diff := Difference

	if k := diff.Apply(12, 42); k != 30 {
		t.Errorf("%.1v(12, 42): expected 30, got %d", diff, k)
	}

	if k := diff.Apply(42, 12); k != 30 {
		t.Errorf("%.1v(42, 12): expected 30, got %d", diff, k)
	}
}
