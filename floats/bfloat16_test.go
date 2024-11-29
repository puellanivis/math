package floats

import (
	"math"
	"math/big"
	"testing"
)

func TestBFloat16ConvertConstants(t *testing.T) {
	type test struct {
		name  string
		value any
	}

	tests := []test{}

	for _, tt := range tests {
		x, err := fromAny[binary16](t, tt.value)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		t.Errorf("%s: uint16( %#04x )", tt.name, x)
	}
}

func TestBFloat16NumConstants(t *testing.T) {
	type test struct {
		name  string
		value string
		bits  uint16
	}

	tests := []test{
		{"zero", "0", 0x0000},
		{"one", "1", 0x3f80},
		{"e", "2.71828182845904523536028747135266249775724709369995957496696763", E.BFloat16().bits},
		{"π", "0x1.921fb54442d18469898cc51701b8p+01", Pi.BFloat16().bits},
		{"π", "3.14159265358979323846264338327950288419716939937510582097494459", Pi.BFloat16().bits},
		{"φ", "1.61803398874989484820458683436563811772030917980576286213544862", Phi.BFloat16().bits},
		{"√2", "1.41421356237309504880168872420969807856967187537694807317667974", Sqrt2.BFloat16().bits},
		{"√e", "1.64872127070012814684865078781416357165377610071014801157507931", SqrtE.BFloat16().bits},
		{"√π", "1.77245385090551602729816748334114518279754945612238712821380779", SqrtPi.BFloat16().bits},
		{"√φ", "1.27201964951406896425242246173749149171560804184009624861664038", SqrtPhi.BFloat16().bits},
		{"ln2", "0.693147180559945309417232121458176568075500134360255254120680009", Ln2.BFloat16().bits},
		{"1/ln2", "1.44269504088896340735992468100189213742664595415298593413544940", Ln2E.BFloat16().bits},
		{"ln10", "2.30258509299404568401799145468436420760110148862877297603332790", Ln10.BFloat16().bits},
		{"1/ln10", "0.4342944819032518276511289189166050822943970058036665661144537", Ln10E.BFloat16().bits},
		{"maxMag", "3.3895e+38", MaxBFloat16.bits},
		{"minMag", "9.1835e-41", SmallestNonzeroBFloat16.bits},
		{"2^e", "6.5808859910179209708515424038864864915730774383480740051215126610", 0x40d3},
		{"2^π", "8.8249778270762876238564296042080015817044108152714849266689598650", 0x410d},
		{"2^φ", "3.0695645076529788214628616541515824382196579505278255446333356097", 0x4044},
	}

	for _, tt := range tests {
		f := new(big.Float).SetPrec(16)
		if _, _, err := f.Parse(tt.value, 0); err != nil {
			t.Fatal("unexpected error:", err)
		}

		b16 := BFloat16FromFloat(f)

		t.Logf("%s: %.4e %.4e", tt.name, f, b16)
		t.Logf("%s: %04x", tt.name, b16.bits)

		if b16.bits != tt.bits {
			expected := BFloat16{tt.bits}
			t.Logf("expected: %.4e, %b", expected, expected)
			t.Errorf("Parse(%q) = %04x, expected: %04x", tt.value, b16.bits, tt.bits)
		}
	}
}

func TestBFloat16Numbers(t *testing.T) {
	type test struct {
		name string
		in   float32
		bits uint16
	}

	var spec bfloat16
	exp2of := expBias[bfloat16]() + 1
	exp2uf := expBias[bfloat16]() - 1 + spec.mantWidth()

	exp2ofBits, exp2ufBits := spec.exp2OverUnder()

	tests := []test{
		{"zero", 0, 0},
		{"one", 1, 0x3f80},
		{"one+tiny", 0x1.ffep0, 0x4000}, // rounding required
		{"pi", math.Pi, 0x4049},
		{"of", float32(exp2of), exp2ofBits}, // abs(overflow)
		{"uf", float32(exp2uf), exp2ufBits}, // abs(underflow)
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := BFloat16FromFloat(tt.in)

			b32 := f.Float32()
			f32m1, f32p1 := f.NextDown().Float32().Native(), f.NextUp().Float32().Native()

			if f32m1 > tt.in {
				t.Errorf("Roundtrip(%x) = %x, expected %x <= %x", tt.in, b32, f32m1, tt.in)
			}

			if tt.in > f32p1 {
				t.Errorf("Roundtrip(%x) = %x, expected %x <= %x", tt.in, b32, tt.in, f32p1)
			}

			bits := f.Bits()

			t.Logf("e:\n%.3e\n%.3e", f, b32)
			t.Logf("x:\n0x%.2x\n0x%.2x", f, b32)
			t.Logf("b:\n0b%b\n0b%b", f, b32)
			t.Logf("x:\n%#04x\n%#04x", f.bits, b32.bits)

			if bits != tt.bits {
				t.Errorf("BFloat16FromFloat(%x)\n  actual: %04x\nexpected: %04x", tt.in, bits, tt.bits)
			}
		})
	}
}

func TestBFloat16OpAdd(t *testing.T) {
	type test struct {
		name string
		x, y float32
		bits uint16
	}

	tests := []test{
		{"tau", math.Pi, math.Pi, 0x40c9},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, g := BFloat16FromFloat(tt.x), BFloat16FromFloat(tt.y)

			t.Logf("x:\n%%e: %.3e\n%%b: 0b%b", f, f)
			t.Logf("x:\n%%e: %.3e\n%%b: 0b%b", g, g)

			res := f.Add(g)
			bits := res.Bits()

			t.Logf("r:\n%%e: %.3e\n%%b: 0b%b", res, res)

			if bits != tt.bits {
				t.Errorf("BFloat16(%x) + BFloat16(%x):\n  actual: %04x\nexpected: %04x", tt.x, tt.y, bits, tt.bits)
			}
		})
	}
}

func TestBFloat16OpMul(t *testing.T) {
	type test struct {
		name string
		x, y float32
		bits uint16
	}

	tests := []test{
		//                                 0000
		{"pi_squared", math.Pi, math.Pi, 0x411e},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, g := BFloat16FromFloat(tt.x), BFloat16FromFloat(tt.y)

			t.Logf("x:\n%%e: %.3e\n%%b: 0b%b", f, f)
			t.Logf("x:\n%%e: %.3e\n%%b: 0b%b", g, g)

			res := f.Mul(g)
			bits := res.Bits()

			t.Logf("r:\n%%e: %.3e\n%%b: 0b%b", res, res)

			if bits != tt.bits {
				t.Errorf("BFloat16(%x) × BFloat16(%x):\n  actual: %04x\nexpected: %04x", tt.x, tt.y, bits, tt.bits)
			}
		})
	}
}

func TestBFloat16OpDiv(t *testing.T) {
	type test struct {
		name string
		x, y BFloat16
		bits uint16
	}

	tests := []test{
		{"log2e", BFloat16FromFloat(1.0), BFloat16(Ln2.BFloat16()), 0x3fb9},
		{"log10e", BFloat16FromFloat(1.0), BFloat16(Ln10.BFloat16()), 0x3edf},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, g := tt.x, tt.y

			t.Logf("x:\n%%e: %.3e\n%%b: 0b%b", f, f)
			t.Logf("x:\n%%e: %.3e\n%%b: 0b%b", g, g)

			res := f.Div(g)
			bits := res.Bits()

			t.Logf("r:\n%%e: %.3e\n%%b: 0b%b", res, res)

			if bits != tt.bits {
				t.Errorf("BFloat16(%x) ÷ BFloat16(%x)\n  actual: %04x\nexpected: %04x", tt.x, tt.y, bits, tt.bits)
			}
		})
	}
}

func TestBFloat16OpExp2(t *testing.T) {
	type test struct {
		name string
		x    BFloat16
		bits uint16
	}

	tests := []test{
		{"2^e", BFloat16(E.BFloat16()), 0x40d3},
		{"2^π", BFloat16(Pi.BFloat16()), 0x410d},
		{"2^φ", BFloat16(Phi.BFloat16()), 0x4044},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			res := f.Exp2()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.3e\nactual: %.3e\nexpect: %.3e", f, res, BFloat16{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, BFloat16{tt.bits})

			if bits&^0x0f != tt.bits&^0x0f {
				t.Errorf("2 ** BFloat16(%x)\n  actual: %04x\nexpected: %04x", tt.x, bits, tt.bits)
			}
		})
	}
}

func TestBFloat16OpExp(t *testing.T) {
	type test struct {
		name string
		x    BFloat16
		bits uint16
	}

	tests := []test{
		{"e^e", BFloat16(E.BFloat16()), 0x4172},
		{"e^π", BFloat16(Pi.BFloat16()), 0x41b9},
		{"e^φ", BFloat16(Phi.BFloat16()), 0x40a1},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			res := f.Exp()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.3e\nactual: %.3e\nexpect: %.3e", f, res, BFloat16{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, BFloat16{tt.bits})

			if bits&^0x0f != tt.bits&^0x0f {
				t.Errorf("e ** BFloat16(%x)\n  actual: %04x\nexpected: %04x", tt.x, bits, tt.bits)
			}
		})
	}
}

func TestBFloat16OpSqrt(t *testing.T) {
	type test struct {
		name string
		x    BFloat16
		bits uint16
	}

	tests := []test{
		{"√2", BFloat16{two[bfloat16]()}, Sqrt2.BFloat16().bits},
		{"√e", BFloat16(E.BFloat16()), SqrtE.BFloat16().bits},
		{"√π", BFloat16(Pi.BFloat16()), SqrtPi.BFloat16().bits},
		{"√φ", BFloat16(Phi.BFloat16()), SqrtPhi.BFloat16().bits},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			res := f.Sqrt()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.3e\nactual: %.3e\nexpect: %.3e", f, res, BFloat16{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, BFloat16{tt.bits})

			if bits != tt.bits {
				t.Errorf("√BFloat16(%x)\n  actual: %04x\nexpected: %04x", tt.x, bits, tt.bits)
			}
		})
	}
}
