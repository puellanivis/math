package floats

import (
	"math"
	"math/big"
	"testing"
)

func TestFloat64ConvertConstants(t *testing.T) {
	type test struct {
		name  string
		value any
	}

	tests := []test{}

	for _, tt := range tests {
		x, err := fromAny[binary64](t, tt.value)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		t.Errorf("%s: uint64( %#016x )", tt.name, x)
	}
}

func TestFloat64NumConstants(t *testing.T) {
	type test struct {
		name  string
		value string
		bits  uint64
	}

	tests := []test{
		{"zero", "0", 0x0000_0000_0000_0000},
		{"one", "1", 0x3ff0_0000_0000_0000},
		{"e", "2.71828182845904523536028747135266249775724709369995957496696763", E.Float64().bits},
		{"π", "0x1.921fb54442d18469898cc51701b8p+01", Pi.Float64().bits},
		{"π", "3.14159265358979323846264338327950288419716939937510582097494459", Pi.Float64().bits},
		{"φ", "1.61803398874989484820458683436563811772030917980576286213544862", Phi.Float64().bits},
		{"√2", "1.41421356237309504880168872420969807856967187537694807317667974", Sqrt2.Float64().bits},
		{"√e", "1.64872127070012814684865078781416357165377610071014801157507931", SqrtE.Float64().bits},
		{"√π", "1.77245385090551602729816748334114518279754945612238712821380779", SqrtPi.Float64().bits},
		{"√φ", "1.27201964951406896425242246173749149171560804184009624861664038", SqrtPhi.Float64().bits},
		{"ln2", "0.693147180559945309417232121458176568075500134360255254120680009", Ln2.Float64().bits},
		{"1/ln2", "1.44269504088896340735992468100189213742664595415298593413544940", Ln2E.Float64().bits},
		{"ln10", "2.30258509299404568401799145468436420760110148862877297603332790", Ln10.Float64().bits},
		{"1/ln10", "0.4342944819032518276511289189166050822943970058036665661144537", Ln10E.Float64().bits},
		{"maxMag", "1.79769313486231570814527423731704356798070e+308", MaxFloat64.bits},
		{"minMag", "4.9406564584124654417656879286822137236505980e-324", SmallestNonzeroFloat64.bits},
		{"2^e", "6.5808859910179209708515424038864864915730774383480740051215126610", 0x401a52d3c6f8818f},
		{"2^π", "8.8249778270762876238564296042080015817044108152714849266689598650", 0x4021a6637e666f83},
		{"2^φ", "3.0695645076529788214628616541515824382196579505278255446333356097", 0x40088e77d62aa7b4},
	}

	for _, tt := range tests {
		f := new(big.Float).SetPrec(64)
		if _, _, err := f.Parse(tt.value, 0); err != nil {
			t.Fatal("unexpected error:", err)
		}

		b64 := Float64FromFloat(f)

		t.Logf("%s: %.16e %.16e", tt.name, f, b64)
		t.Logf("%s: %016x", tt.name, b64.bits)

		if b64.bits != tt.bits {
			expected := Float64{tt.bits}
			t.Logf("expected: %.16e, %b", expected, expected)
			t.Errorf("Parse(%q):\n  actual: %016x\nexpected: %016x", tt.value, b64.bits, tt.bits)
		}
	}
}

func TestFloat64Numbers(t *testing.T) {
	type test struct {
		name string
		in   float64
		bits uint64
	}

	var spec binary64
	exp2of := expBias[binary64]() + 1
	exp2uf := expBias[binary64]() - 1 + spec.mantWidth()

	exp2ofBits, exp2ufBits := spec.exp2OverUnder()

	tests := []test{
		{"zero", 0, 0},
		{"one", 1, 0x3ff0_0000_0000_0000},
		{"pi", math.Pi, 0x4009_21fb_5444_2d18},
		{"of", float64(exp2of), exp2ofBits}, // abs(overflow)
		{"uf", float64(exp2uf), exp2ufBits}, // abs(underflow)
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := Float64FromFloat(tt.in)

			b64 := f.Float64()
			f64m1, f64p1 := f.NextDown().Native(), f.NextUp().Native()

			if f64m1 > tt.in {
				t.Errorf("Roundtrip(%x) = %x, expected %x <= %x", tt.in, b64, f64m1, tt.in)
			}

			if tt.in > f64p1 {
				t.Errorf("Roundtrip(%x) = %x, expected %x <= %x", tt.in, b64, tt.in, f64p1)
			}

			bits := f.Bits()

			t.Logf("e:\n%.16e\n%.16e", f, b64)
			t.Logf("x:\n0x%.13x\n0x%.13x", f, b64)
			t.Logf("b:\n0b%b\n0b%b", f, b64)
			t.Logf("x:\n%#016x\n%#016x", f.bits, b64.bits)

			if bits != tt.bits {
				t.Errorf("Float64FromFloat(%x)\n  actual: %016x\nexpected: %016x", tt.in, bits, tt.bits)
			}
		})
	}
}

func TestFloat64OpExp2(t *testing.T) {
	type test struct {
		name string
		x    Float64
		bits uint64
	}

	tests := []test{
		{"2^e", E.Float64(), 0x401a52d3c6f8818f},
		{"2^π", Pi.Float64(), 0x4021a6637e666f83},
		{"2^φ", Phi.Float64(), 0x40088e77d62aa7b4},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			t.Logf("x:\n%%e: %.16e\n%%b: 0b%b", f, f)

			res := f.Exp2()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.16e\nactual: %.16e\nexpect: %.16e", f, res, Float64{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, Float64{tt.bits})

			if bits&^0xFFFF_FFFF != tt.bits&^0xFFFF_FFFF {
				t.Logf("expect:\n%%e: %.16e\n%%b: 0b%b", Float64{tt.bits}, Float64{tt.bits})
				t.Errorf("2 ** Float64(%x)\n  actual: %016x\nexpected: %016x", tt.x, bits, tt.bits)
			}
		})
	}
}

func TestFloat64OpExp(t *testing.T) {
	type test struct {
		name string
		x    Float64
		bits uint64
	}

	tests := []test{
		{"e^e", E.Float64(), 0x402e4efb75e4527b},
		{"e^π", Pi.Float64(), 0x403724046eb0933a},
		{"e^φ", Phi.Float64(), 0x40142c339d4a2b22},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			res := f.Exp()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.16e\nactual: %.16e\nexpect: %.16e", f, res, Float64{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, Float64{tt.bits})

			if bits&^0xFFFF_FFFF != tt.bits&^0xFFFF_FFFF {
				t.Errorf("e ** Float64(%x)\n  actual: %016x\nexpected: %016x", tt.x, bits, tt.bits)
			}
		})
	}
}
