package floats

import (
	"math"
	"math/big"
	"testing"

	"github.com/puellanivis/math/bits"
)

func TestFloat128ConvertConstants(t *testing.T) {
	type test struct {
		name  string
		value any
	}

	tests := []test{}

	for _, tt := range tests {
		x, err := fromAny[binary128](t, tt.value)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		t.Errorf("%s: bits.Uint128{ Hi: %#016x, Lo: %#016x }", tt.name, x.Hi, x.Lo)
	}
}

func TestFloat128NumConstants(t *testing.T) {
	type test struct {
		name  string
		value string
		bits  bits.Uint128
	}

	p2e := bits.Uint128{Hi: 0x4001a52d3c6f8818, Lo: 0xe8aede80fd627e32}
	p2π := bits.Uint128{Hi: 0x40021a6637e666f8, Lo: 0x34db52a7e461c2d8}
	p2φ := bits.Uint128{Hi: 0x400088e77d62aa7b, Lo: 0x3b8037ce5204704f}

	// Constants retrieved from Wolfram|Alpha:
	tests := []test{
		{"zero", "0", bits.Uint128{Hi: 0x0000_0000_0000_0000, Lo: 0}},
		{"one", "1", bits.Uint128{Hi: 0x3fff_0000_0000_0000, Lo: 0}},
		{"e", "2.71828182845904523536028747135266249775724709369995957496696763", E.bits},
		{"π", "0x1.921fb54442d18469898cc51701b8p+01", Pi.bits},
		{"π", "3.14159265358979323846264338327950288419716939937510582097494459", Pi.bits},
		{"φ", "1.61803398874989484820458683436563811772030917980576286213544862", Phi.bits},
		{"√2", "1.41421356237309504880168872420969807856967187537694807317667974", Sqrt2.bits},
		{"√e", "1.64872127070012814684865078781416357165377610071014801157507931", SqrtE.bits},
		{"√π", "1.77245385090551602729816748334114518279754945612238712821380779", SqrtPi.bits},
		{"√φ", "1.27201964951406896425242246173749149171560804184009624861664038", SqrtPhi.bits},
		{"ln2", "0.693147180559945309417232121458176568075500134360255254120680009", Ln2.bits},
		{"1/ln2", "1.44269504088896340735992468100189213742664595415298593413544940", Ln2E.bits},
		{"ln10", "2.30258509299404568401799145468436420760110148862877297603332790", Ln10.bits},
		{"1/ln10", "0.4342944819032518276511289189166050822943970058036665661144537", Ln10E.bits},
		{"maxMag", "1.1897314953572317650857593266280070e+4932", MaxFloat128.bits},
		{"minMag", "6.4751751194380251109244389582276466e-4966", SmallestNonzeroFloat128.bits},
		{"2^e", "6.5808859910179209708515424038864864915730774383480740051215126610", p2e},
		{"2^π", "8.8249778270762876238564296042080015817044108152714849266689598650", p2π},
		{"2^φ", "3.0695645076529788214628616541515824382196579505278255446333356097", p2φ},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := new(big.Float).SetPrec(256)
			if _, _, err := f.Parse(tt.value, 0); err != nil {
				t.Fatal("unexpected error:", err)
			}

			f = f.SetMode(big.ToNearestEven).SetPrec(128)

			b128 := Float128FromFloat(f)

			t.Logf("%s: %.34e %.34e", tt.name, f, b128)
			t.Logf("%s: %032x", tt.name, b128.bits)

			if b128.bits != tt.bits {
				expected := Float128{tt.bits}
				t.Logf("expected: %.34e, %b", expected, expected)
				t.Errorf("Parse(%q):\n  actual: %032x\nexpected: %032x", tt.value, b128.bits, tt.bits)
			}
		})
	}
}

func TestFloat128Numbers(t *testing.T) {
	type test struct {
		name string
		in   float64
		bits bits.Uint128
	}

	var spec binary128
	exp2of := expBias[binary128]() + 1
	exp2uf := expBias[binary128]() - 1 + spec.mantWidth()

	exp2ofBits, exp2ufBits := spec.exp2OverUnder()

	tests := []test{
		{"zero", 0, bits.Uint128{Hi: 0, Lo: 0}},
		{"one", 1, bits.Uint128{Hi: 0x3fff_0000_0000_0000, Lo: 0}},
		{"one+tiny", 0x1.ffep0, bits.Uint128{Hi: 0x3fff_ffe0_0000_0000, Lo: 0}}, // rounding required
		{"pi", math.Pi, bits.Uint128{Hi: 0x4000_921f_b544_42d1, Lo: 0x8469_898c_c517_01b8}},
		{"of", float64(exp2of), exp2ofBits}, // abs(overflow)
		{"uf", float64(exp2uf), exp2ufBits}, // abs(underflow)
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := Float128FromFloat(tt.in)

			if tt.name == "pi" {
				f.bits = Pi.bits
			}

			b64 := f.Float64()
			f64m1, f64p1 := f.NextDown().Float64().Native(), f.NextUp().Float64().Native()

			if f64m1 > tt.in {
				t.Errorf("Roundtrip(%x) = %x, expected %x <= %x", tt.in, b64, f64m1, tt.in)
			}

			if tt.in > f64p1 {
				t.Errorf("Roundtrip(%x) = %x, expected %x <= %x", tt.in, b64, tt.in, f64p1)
			}

			bits := f.Bits()

			t.Logf("e:\n%.34e\n%.34e", f, b64)
			t.Logf("x:\n0x%.28x\n0x%.28x", f, b64)
			t.Logf("b:\n0b%b\n0b%b", f, b64)
			t.Logf("x:\n%#032x\n%#016x", f.bits, b64.bits)

			if bits != tt.bits {
				t.Errorf("Float128FromFloat(%x)\n  actual: %032x\nexpected: %032x", tt.in, bits, tt.bits)
			}
		})
	}
}

func TestFloat128OpAdd(t *testing.T) {
	type test struct {
		name string
		x, y Float128
		bits bits.Uint128
	}

	tests := []test{
		{"tau", Pi, Pi, bits.Uint128{Hi: 0x4001921fb54442d1, Lo: 0x8469898cc51701b8}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, g := tt.x, tt.y

			t.Logf("x:\n%%e: %.34e\n%%b: 0b%b", f, f)
			t.Logf("x:\n%%e: %.34e\n%%b: 0b%b", g, g)

			res := f.Add(g)
			bits := res.Bits()

			t.Logf("r:\n%%e: %.34e\n%%b: 0b%b", res, res)

			if bits != tt.bits {
				t.Errorf("Float128(%x) + Float128(%x):\n  actual: %032x\nexpected: %032x", tt.x, tt.y, bits, tt.bits)
			}
		})
	}
}

func TestFloat128OpMul(t *testing.T) {
	type test struct {
		name string
		x, y Float128
		bits bits.Uint128
	}

	tests := []test{
		//                                    0000000000000000
		{"pi_squared", Pi, Pi, bits.Uint128{Hi: 0x40023bd3cc9be45d, Lo: 0xe5a4adc4d9b30118}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, g := tt.x, tt.y

			t.Logf("x:\n%%e: %.34e\n%%b: 0b%b", f, f)
			t.Logf("x:\n%%e: %.34e\n%%b: 0b%b", g, g)

			res := f.Mul(g)
			bits := res.Bits()

			t.Logf("r:\n%%e: %.34e\n%%b: 0b%b", res, res)

			if bits != tt.bits {
				t.Errorf("Float128(%x) × Float128(%x):\n  actual: %032x\nexpected: %032x", tt.x, tt.y, bits, tt.bits)
			}
		})
	}
}

func TestFloat128OpDiv(t *testing.T) {
	type test struct {
		name string
		x, y Float128
		bits bits.Uint128
	}

	tests := []test{
		{"log2e", Float128FromFloat(1.0), Ln2, bits.Uint128{Hi: 0x3fff71547652b82f, Lo: 0xe1777d0ffda0d23a}},
		{"log10e", Float128FromFloat(1.0), Ln10, bits.Uint128{Hi: 0x3ffdbcb7b1526e50, Lo: 0xe32a6ab7555f5a68}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, g := tt.x, tt.y

			t.Logf("x:\n%%e: %.34e\n%%b: 0b%b", f, f)
			t.Logf("x:\n%%e: %.34e\n%%b: 0b%b", g, g)

			res := f.Div(g)
			bits := res.Bits()

			t.Logf("r:\n%%e: %.34e\n%%b: 0b%b", res, res)

			if bits != tt.bits {
				t.Errorf("Float128(%x) ÷ Float128(%x)\n  actual: %032x\nexpected: %032x", tt.x, tt.y, bits, tt.bits)
			}
		})
	}
}

func TestFloat128OpExp2(t *testing.T) {
	type test struct {
		name string
		x    Float128
		bits bits.Uint128
	}

	tests := []test{
		{"2^e", E, bits.Uint128{Hi: 0x4001a52d3c6f8818, Lo: 0xe8aede80fd627e32}},
		{"2^π", Pi, bits.Uint128{Hi: 0x40021a6637e666f8, Lo: 0x34db52a7e461c2d8}},
		{"2^φ", Phi, bits.Uint128{Hi: 0x400088e77d62aa7b, Lo: 0x3b8037ce5204704f}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			res := f.Exp2()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.34e\nactual: %.34e\nexpect: %.34e", f, res, Float128{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, Float128{tt.bits})

			if bits.Hi != tt.bits.Hi {
				t.Errorf("2 ** Float128(%x)\n  actual: %032x\nexpected: %032x", tt.x, bits, tt.bits)
			}
		})
	}
}

func TestFloat128OpExp(t *testing.T) {
	type test struct {
		name string
		x    Float128
		bits bits.Uint128
	}

	tests := []test{
		{"e^e", E, bits.Uint128{Hi: 0x4002e4efb75e4527, Lo: 0xaf5a730011f27cf7}},
		{"e^π", Pi, bits.Uint128{Hi: 0x4003724046eb0933, Lo: 0x99ecda7489f9ab77}},
		{"e^φ", Phi, bits.Uint128{Hi: 0x400142c339d4a2b2, Lo: 0x23092cc388b50f19}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			res := f.Exp()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.34e\nactual: %.34e\nexpect: %.34e", f, res, Float128{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, Float128{tt.bits})

			if bits.Hi != tt.bits.Hi {
				t.Errorf("e ** Float128(%x)\n  actual: %032x\nexpected: %032x", tt.x, bits, tt.bits)
			}
		})
	}
}

func TestFloat128OpSqrt(t *testing.T) {
	type test struct {
		name string
		x    Float128
		bits bits.Uint128
	}

	tests := []test{
		{"√2", Float128{two[binary128]()}, Sqrt2.bits},
		{"√e", E, SqrtE.bits},
		{"√π", Pi, SqrtPi.bits},
		{"√φ", Phi, SqrtPhi.bits},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			res := f.Sqrt()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.34e\nactual: %.34e\nexpect: %.34e", f, res, Float128{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, Float128{tt.bits})

			if bits != tt.bits {
				t.Errorf("√Float128(%x)\n  actual: %032x\nexpected: %032x", tt.x, bits, tt.bits)
			}
		})
	}
}
