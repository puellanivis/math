package floats

import (
	"math"
	"math/big"
	"testing"
)

func TestFloat32ConvertConstants(t *testing.T) {
	type test struct {
		name  string
		value any
	}

	tests := []test{}

	for _, tt := range tests {
		x, err := fromAny[binary32](t, tt.value)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		t.Errorf("%s: uint32( %#08x )", tt.name, x)
	}
}

func TestFloat32NumConstants(t *testing.T) {
	type test struct {
		name  string
		value string
		bits  uint32
	}

	tests := []test{
		{"zero", "0", 0x0000_0000},
		{"one", "1", 0x3f80_0000},
		{"e", "2.71828182845904523536028747135266249775724709369995957496696763", E.Float32().bits},
		{"π", "0x1.921fb54442d18469898cc51701b8p+01", Pi.Float32().bits},
		{"π", "3.14159265358979323846264338327950288419716939937510582097494459", Pi.Float32().bits},
		{"φ", "1.61803398874989484820458683436563811772030917980576286213544862", Phi.Float32().bits},
		{"√2", "1.41421356237309504880168872420969807856967187537694807317667974", Sqrt2.Float32().bits},
		{"√e", "1.64872127070012814684865078781416357165377610071014801157507931", SqrtE.Float32().bits},
		{"√π", "1.77245385090551602729816748334114518279754945612238712821380779", SqrtPi.Float32().bits},
		{"√φ", "1.27201964951406896425242246173749149171560804184009624861664038", SqrtPhi.Float32().bits},
		{"ln2", "0.693147180559945309417232121458176568075500134360255254120680009", Ln2.Float32().bits},
		{"1/ln2", "1.44269504088896340735992468100189213742664595415298593413544940", Ln2E.Float32().bits},
		{"ln10", "2.30258509299404568401799145468436420760110148862877297603332790", Ln10.Float32().bits},
		{"1/ln10", "0.4342944819032518276511289189166050822943970058036665661144537", Ln10E.Float32().bits},
		{"maxMag", "3.40282346638528859811704183484516925440e+38", MaxFloat32.bits},
		{"minMag", "1.401298464324817070923729583289916131280e-45", SmallestNonzeroFloat32.bits},
		{"2^e", "6.5808859910179209708515424038864864915730774383480740051215126610", 0x40d2969e},
		{"2^π", "8.8249778270762876238564296042080015817044108152714849266689598650", 0x410d331c},
		{"2^φ", "3.0695645076529788214628616541515824382196579505278255446333356097", 0x404473bf},
	}

	for _, tt := range tests {
		f := new(big.Float).SetPrec(32)
		if _, _, err := f.Parse(tt.value, 0); err != nil {
			t.Fatal("unexpected error:", err)
		}

		b32 := Float32FromFloat(f)

		t.Logf("%s: %.7e %.7e", tt.name, f, b32)
		t.Logf("%s: %08x", tt.name, b32.bits)

		if b32.bits != tt.bits {
			expected := Float32{tt.bits}
			t.Logf("expected: %.7e, %b", expected, expected)
			t.Errorf("Parse(%q):\n  actual: %08x\nexpected: %08x", tt.value, b32.bits, tt.bits)
		}
	}
}

func TestFloat32Numbers(t *testing.T) {
	type test struct {
		name string
		in   float32
		bits uint32
	}

	var spec binary32
	exp2of := expBias[binary32]() + 1
	exp2uf := expBias[binary32]() - 1 + spec.mantWidth()

	exp2ofBits, exp2ufBits := spec.exp2OverUnder()

	tests := []test{
		{"zero", 0, 0},
		{"one", 1, 0x3f80_0000},
		{"pi", math.Pi, 0x4049_0fdb},
		{"of", float32(exp2of), exp2ofBits}, // abs(overflow)
		{"uf", float32(exp2uf), exp2ufBits}, // abs(underflow)
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := Float32FromFloat(tt.in)

			b32 := f.Float32()
			f32m1, f32p1 := f.NextDown().Native(), f.NextUp().Native()

			if f32m1 > tt.in {
				t.Errorf("Roundtrip(%x) = %x, expected %x <= %x", tt.in, b32, f32m1, tt.in)
			}

			if tt.in > f32p1 {
				t.Errorf("Roundtrip(%x) = %x, expected %x <= %x", tt.in, b32, tt.in, f32p1)
			}

			bits := f.Bits()

			t.Logf("e:\n%.7e\n%.7e", f, b32)
			t.Logf("x:\n0x%.6x\n0x%.6x", f, b32)
			t.Logf("b:\n0b%b\n0b%b", f, b32)
			t.Logf("x:\n%#08x\n%#08x", f.bits, b32.bits)

			if bits != tt.bits {
				t.Errorf("Float32FromFloat(%x)\n  actual: %08x\nexpected: %08x", tt.in, bits, tt.bits)
			}
		})
	}
}

func TestFloat32OpExp2(t *testing.T) {
	type test struct {
		name string
		x    Float32
		bits uint32
	}

	tests := []test{
		{"2^e", E.Float32(), 0x40d2969e},
		{"2^π", Pi.Float32(), 0x410d331c},
		{"2^φ", Phi.Float32(), 0x404473bf},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			res := f.Exp2()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.7e\nactual: %.7e\nexpect: %.7e", f, res, Float32{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, Float32{tt.bits})

			if bits&^0xFFFF != tt.bits&^0xFFFF {
				t.Errorf("2 ** Float32(%x)\n  actual: %08x\nexpected: %08x", tt.x, bits, tt.bits)
			}
		})
	}
}

func TestFloat32OpExp(t *testing.T) {
	type test struct {
		name string
		x    Float32
		bits uint32
	}

	tests := []test{
		{"e^e", E.Float32(), 0x417277dc},
		{"e^π", Pi.Float32(), 0x41b92023},
		{"e^φ", Phi.Float32(), 0x40a1619d},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			res := f.Exp()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.7e\nactual: %.7e\nexpect: %.7e", f, res, Float32{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, Float32{tt.bits})

			if bits&^0xFFFF != tt.bits&^0xFFFF {
				t.Errorf("e ** Float32(%x)\n  actual: %08x\nexpected: %08x", tt.x, bits, tt.bits)
			}
		})
	}
}
