package floats

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestFloat16ConvertConstants(t *testing.T) {
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

func TestFloat16NumConstants(t *testing.T) {
	type test struct {
		name  string
		value string
		bits  uint16
	}

	tests := []test{
		{"zero", "0", 0x0000},
		{"one", "1", 0x3c00},
		{"e", "2.71828182845904523536028747135266249775724709369995957496696763", E.Float16().bits},
		{"π", "0x1.921fb54442d18469898cc51701b8p+01", Pi.Float16().bits},
		{"π", "3.14159265358979323846264338327950288419716939937510582097494459", Pi.Float16().bits},
		{"φ", "1.61803398874989484820458683436563811772030917980576286213544862", Phi.Float16().bits},
		{"√2", "1.41421356237309504880168872420969807856967187537694807317667974", Sqrt2.Float16().bits},
		{"√e", "1.64872127070012814684865078781416357165377610071014801157507931", SqrtE.Float16().bits},
		{"√π", "1.77245385090551602729816748334114518279754945612238712821380779", SqrtPi.Float16().bits},
		{"√φ", "1.27201964951406896425242246173749149171560804184009624861664038", SqrtPhi.Float16().bits},
		{"ln2", "0.693147180559945309417232121458176568075500134360255254120680009", Ln2.Float16().bits},
		{"1/ln2", "1.44269504088896340735992468100189213742664595415298593413544940", Ln2E.Float16().bits},
		{"ln10", "2.30258509299404568401799145468436420760110148862877297603332790", Ln10.Float16().bits},
		{"1/ln10", "0.4342944819032518276511289189166050822943970058036665661144537", Ln10E.Float16().bits},
		{"maxMag", "6.5504e+04", MaxFloat16.bits},
		{"minMag", "5.96046e-08", SmallestNonzeroFloat16.bits},
		{"2^e", "6.5808859910179209708515424038864864915730774383480740051215126610", 0x4695},
		{"2^π", "8.8249778270762876238564296042080015817044108152714849266689598650", 0x486a},
		{"2^φ", "3.0695645076529788214628616541515824382196579505278255446333356097", 0x4224},
	}

	for _, tt := range tests {
		f := new(big.Float).SetPrec(16)
		if _, _, err := f.Parse(tt.value, 0); err != nil {
			t.Fatal("unexpected error:", err)
		}

		b16 := Float16FromFloat(f)

		t.Logf("%s: %.4e %.4e", tt.name, f, b16)
		t.Logf("%s: %04x", tt.name, b16.bits)

		if b16.bits != tt.bits {
			expected := Float16{tt.bits}
			t.Logf("expected: %.4e, %b", expected, expected)
			t.Errorf("Parse(%q) = %04x, expected: %04x", tt.value, b16.bits, tt.bits)
		}
	}
}

func TestFloat16Numbers(t *testing.T) {
	type test struct {
		name string
		in   float32
		bits uint16
	}

	var spec binary16
	exp2of := expBias[binary16]() + 1
	exp2uf := expBias[binary16]() - 1 + spec.mantWidth()

	exp2ofBits, exp2ufBits := spec.exp2OverUnder()

	tests := []test{
		{"zero", 0, 0},
		{"one", 1, 0x3c00},
		{"one+tiny", 0x1.ffep0, 0x4000}, // rounding required
		{"pi", math.Pi, 0x4248},
		{"of", float32(exp2of), exp2ofBits}, // abs(overflow)
		{"uf", float32(exp2uf), exp2ufBits}, // abs(underflow)
		{"epsilon12", 0x1p-12, 0x0c00},
		{"epsilon13", 0x1p-13, 0x0800},
		{"epsilon14", 0x1p-14, 0x0400},
		{"epsilon15", 0x1p-15, 0x0200},
		{"epsilon16", 0x1p-16, 0x0100},
		{"epsilon17", 0x1p-17, 0x0080},
		{"epsilon18", 0x1p-18, 0x0040},
		{"epsilon19", 0x1p-19, 0x0020},
		{"epsilon20", 0x1p-20, 0x0010},
		{"epsilon21", 0x1p-21, 0x0008},
		{"epsilon22", 0x1p-22, 0x0004},
		{"epsilon23", 0x1p-23, 0x0002},
		{"epsilon24", 0x1p-24, 0x0001},
		{"epsilon25", 0x1p-25, 0x0000}, // round towards even
		{"epsilon26", 0x1p-26, 0x0000},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := Float16FromFloat(tt.in)
			b32 := f.Float32()
			f32m1, f32p1 := f.NextDown().Float32().Native(), f.NextUp().Float32().Native()

			if f32m1 > tt.in {
				t.Errorf("Roundtrip(%x) = %x, expected %x <= %x", tt.in, b32, f32m1, tt.in)
			}

			if tt.in > f32p1 {
				t.Errorf("Roundtrip(%x) = %x, expected %x <= %x", tt.in, b32, tt.in, f32p1)
			}

			bits := f.Bits()

			t.Logf("e:\n%.4e\n%.4e", f, b32)
			t.Logf("x:\n0x%.3x\n0x%.3x", f, b32)
			t.Logf("b:\n0b%b\n0b%b", f, b32)
			t.Logf("x:\n%#04x\n%#04x", f.bits, b32.bits)

			if bits != tt.bits {
				t.Errorf("Float16FromFloat(%x) = %04x, but expected %04x", tt.in, bits, tt.bits)
			}
		})
	}
}

func TestFloat16OpAdd(t *testing.T) {
	type test struct {
		name string
		x, y float32
		bits uint16
	}

	negZ := math.Float32frombits(1 << 31)
	negInf := math.Float32frombits(Inf32(true).Bits())
	posInf := math.Float32frombits(Inf32(false).Bits())

	tests := []test{
		{"zero add zero", 0, 0, 0x0000},
		{"zero add -zero", 0, negZ, 0x0000},
		{"-zero add zero", negZ, 0, 0x0000},
		{"-zero add -zero", negZ, negZ, 0x8000},
		{"zero add one", 0, 1, 0x3c00},
		{"one add zero", 1, 0, 0x3c00},
		{"-zero add one", negZ, 1, 0x3c00},
		{"one add -zero", 1, negZ, 0x3c00},
		{"zero add -one", 0, -1, 0xbc00},
		{"-one add zero", -1, 0, 0xbc00},
		{"-zero add -one", negZ, -1, 0xbc00},
		{"-one add -zero", -1, negZ, 0xbc00},
		{"tau", math.Pi, math.Pi, 0x4648},
		{"one add -one", 1, -1, 0x0000},
		{"two add -two", 2, -2, 0x0000},
		{"tiny add -almost(tiny)", 0x1p-14, -0x1.ff8p-15, 0x0001},
		{"nextDown(inf) add epsilon", 0x1.ffcp15, 0x1p4, 0x7c00}, // round towards even
		{"+inf add +inf", posInf, posInf, 0x7c00},
		{"-inf add -inf", negInf, negInf, 0xfc00},
		{"+inf add -inf", posInf, negInf, 0x7e00},
		{"-inf add +inf", negInf, posInf, 0x7e00},
		{"+inf add 0", posInf, 0, 0x7c00},
		{"-inf add 0", negInf, 0, 0xfc00},
		{"+inf add -0", posInf, negZ, 0x7c00},
		{"-inf add -0", negInf, negZ, 0xfc00},
		{"0 add +inf", 0, posInf, 0x7c00},
		{"0 add -inf", 0, negInf, 0xfc00},
		{"-0 add +inf", negZ, posInf, 0x7c00},
		{"-0 add -inf", negZ, negInf, 0xfc00},
		{"8.125 add -8.125", 8.125, -8.125, 0x0000},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, g := Float16FromFloat(tt.x), Float16FromFloat(tt.y)

			// f = f.NextDown().NextUp()

			t.Logf("x: %x: 0b%b", f, f)
			t.Logf("y: %x: 0b%b", g, g)

			res := f.Add(g)
			bits := res.Bits()

			t.Logf("r: %x: 0b%b", res, res)

			if bits != tt.bits {
				t.Errorf("Float16(%x) + Float16(%x) = %04x, but expected %04x", tt.x, tt.y, bits, tt.bits)
			}
		})
	}
}

func TestFloat16OpMul(t *testing.T) {
	type test struct {
		name string
		x, y float32
		bits uint16
	}

	negZ := math.Float32frombits(1 << 31)
	negInf := math.Float32frombits(Inf32(true).Bits())
	posInf := math.Float32frombits(Inf32(false).Bits())

	tests := []test{
		{"zero mul zero", 0, 0, 0x0000},
		{"zero mul -zero", 0, negZ, 0x8000},
		{"-zero mul zero", negZ, 0, 0x8000},
		{"-zero mul -zero", negZ, negZ, 0x0000},
		{"zero mul one", 0, 1, 0x0000},
		{"one mul zero", 1, 0, 0x0000},
		{"-zero mul one", negZ, 1, 0x8000},
		{"one mul -zero", 1, negZ, 0x8000},
		{"zero mul -one", 0, -1, 0x8000},
		{"-one mul zero", -1, 0, 0x8000},
		{"-zero mul -one", negZ, -1, 0x0000},
		{"-one mul -zero", -1, negZ, 0x0000},
		{"pi_squared", math.Pi, math.Pi, 0x48ef},
		{"one mul -one", 1, -1, 0xbc00},
		{"two mul -two", 2, -2, 0xc400},
		{"tiny mul -almost(tiny)", 0x1p-14, -0x1.ff8p-15, 0x8000},
		{"nextDown(inf) mul 1", 0x1.ffcp15, 1, 0x7bff},
		{"nextDown(inf) mul 1+epsilon_10_001", 0x1.ffcp15, 0x1.0022p0, 0x7c00}, // round towards even
		{"+inf mul +inf", posInf, posInf, 0x7c00},
		{"-inf mul -inf", negInf, negInf, 0x7c00},
		{"+inf mul -inf", posInf, negInf, 0xfc00},
		{"-inf mul +inf", negInf, posInf, 0xfc00},
		{"+inf mul 0", posInf, 0, 0x7e00},
		{"-inf mul 0", negInf, 0, 0x7e00},
		{"+inf mul -0", posInf, negZ, 0x7e00},
		{"-inf mul -0", negInf, negZ, 0x7e00},
		{"0 mul +inf", 0, posInf, 0x7e00},
		{"0 mul -inf", 0, negInf, 0x7e00},
		{"-0 mul +inf", negZ, posInf, 0x7e00},
		{"-0 mul -inf", negZ, negInf, 0x7e00},
		{"8.125 mul -8.125", 8.125, -8.125, 0xd420}, // fraction falls off
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, g := Float16FromFloat(tt.x), Float16FromFloat(tt.y)

			t.Logf("x: %x: 0b%b", f, f)
			t.Logf("y: %x: 0b%b", g, g)

			r32 := f.Float32().Native() * g.Float32().Native()
			t.Logf("r32: %.7f, 0b%b", r32, Float16FromFloat(r32))

			res := f.Mul(g)
			bits := res.Bits()

			t.Logf("  r: %.7f: 0b%b", res, res)

			if bits != tt.bits {
				t.Errorf("Float16(%x) × Float16(%x) = %04x, but expected %04x", tt.x, tt.y, bits, tt.bits)
			}
		})
	}
}

func TestFloat16OpDiv(t *testing.T) {
	type test struct {
		name string
		x, y float32
		bits uint16
	}

	negZ := math.Float32frombits(1 << 31)
	negInf := math.Float32frombits(Inf32(true).Bits())
	posInf := math.Float32frombits(Inf32(false).Bits())

	tests := []test{
		{"zero div zero", 0, 0, 0x7e00},
		{"zero div -zero", 0, negZ, 0x7e00},
		{"-zero div zero", negZ, 0, 0x7e00},
		{"-zero div -zero", negZ, negZ, 0x7e00},
		{"zero div one", 0, 1, 0x0000},
		{"one div zero", 1, 0, 0x7c00},
		{"-zero div one", negZ, 1, 0x8000},
		{"one div -zero", 1, negZ, 0xfc00},
		{"zero div -one", 0, -1, 0x8000},
		{"-one div zero", -1, 0, 0xfc00},
		{"-zero div -one", negZ, -1, 0x0000},
		{"-one div -zero", -1, negZ, 0x7c00},
		{"pi div pi", math.Pi, math.Pi, 0x3c00},
		{"one div -one", 1, -1, 0xbc00},
		{"two div -two", 2, -2, 0xbc00},
		{"tiny div -almost(tiny)", 0x1p-14, -0x1.ff8p-15, 0xbc01},
		{"nextDown(inf) div 1", 0x1.ffcp15, 1, 0x7bff},
		{"nextDown(inf) div 1-epsilon_10_001", 0x1.ffcp15, 0x0.9978p0, 0x7c00}, // round towards even
		{"+inf div +inf", posInf, posInf, 0x7e00},
		{"-inf div -inf", negInf, negInf, 0x7e00},
		{"+inf div -inf", posInf, negInf, 0x7e00},
		{"-inf div +inf", negInf, posInf, 0x7e00},
		{"+inf div 0", posInf, 0, 0x7c00},
		{"-inf div 0", negInf, 0, 0xfc00},
		{"+inf div -0", posInf, negZ, 0xfc00},
		{"-inf div -0", negInf, negZ, 0x7c00},
		{"0 div +inf", 0, posInf, 0x0000},
		{"0 div -inf", 0, negInf, 0x8000},
		{"-0 div +inf", negZ, posInf, 0x8000},
		{"-0 div -inf", negZ, negInf, 0x0000},
		{"8.125 div -8.125", 8.125, -8.125, 0xbc00}, // fraction falls off
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, g := Float16FromFloat(tt.x), Float16FromFloat(tt.y)

			t.Logf("x: %.4f: 0b%b", f, f)
			t.Logf("y: %.4f: 0b%b", g, g)

			r32 := f.Float32().Native() / g.Float32().Native()
			t.Logf("r32: %.7f, 0b%b", r32, Float16FromFloat(r32))

			res := f.Div(g)
			bits := res.Bits()

			t.Logf("  r: %.7f: 0b%b", res, res)

			if bits != tt.bits {
				t.Errorf("Float16(%x) ÷ Float16(%x) = %04x, but expected %04x", tt.x, tt.y, bits, tt.bits)
			}
		})
	}
}

func testRounding(t *testing.T, in uint16, nearTieAway, nearTieEven, toNeg, toZero, toPos uint16) {
	t.Helper()

	f := Float16FromBits(in)

	var g Float16
	var bits uint16

	g = f.Round()
	bits = g.Bits()
	if expected := nearTieAway; bits != expected {
		t.Logf("⸤0b%b⸣ -> 0b%b", f, g)
		t.Logf("⸤%f⸣ -> %f", f, g)
		t.Errorf("Round(%04x) = %04x, but expected %04x", f.Bits(), bits, expected)
	}

	g = f.RoundToEven()
	bits = g.Bits()
	if expected := nearTieEven; bits != expected {
		t.Logf("⸤0b%b⸣ₑ -> 0b%b", f, g)
		t.Logf("⸤%f⸣ₑ -> %f", f, g)
		t.Errorf("RoundToEven(%04x) = %04x, but expected %04x", f.Bits(), bits, expected)
	}

	g = f.Floor()
	bits = g.Bits()
	if expected := toNeg; bits != expected {
		t.Logf("⸤0b%b⸥ -> 0b%b", f, g)
		t.Logf("⸤%f⸥ -> %f", f, g)
		t.Errorf("Floor(%04x) = %04x, but expected %04x", f.Bits(), bits, expected)
	}

	g = f.Trunc()
	bits = g.Bits()
	if expected := toZero; bits != expected {
		t.Logf("int(0b%b) -> 0b%b", f, g)
		t.Logf("int(%f) -> %f", f, g)
		t.Errorf("Trunc(%04x) = %04x, but expected %04x", f.Bits(), bits, expected)
	}

	g = f.Ceil()
	bits = g.Bits()
	if expected := toPos; bits != expected {
		t.Logf("⸢0b%b⸣ -> 0b%b", f, g)
		t.Logf("⸢%f⸣ -> %f", f, g)
		t.Errorf("Ceil(%04x) = %04x, but expected %04x", f.Bits(), bits, expected)
	}
}

func TestFloat16Rounding(t *testing.T) {
	type test struct {
		name string
		in   float32

		nearTieAway uint16
		nearTieEven uint16
		toNeg       uint16
		toZero      uint16
		toPos       uint16
	}

	negZ := math.Float32frombits(1 << 31)
	negInf := math.Float32frombits(Inf32(true).Bits())
	posInf := math.Float32frombits(Inf32(false).Bits())

	tests := []test{
		{"-inf", negInf, 0xfc00, 0xfc00, 0xfc00, 0xfc00, 0xfc00},
		{"-half", -0.5, 0xbc00, 0x8000, 0xbc00, 0x8000, 0x8000},
		{"-zero", negZ, 0x8000, 0x8000, 0x8000, 0x8000, 0x8000},
		{"zero", 0, 0, 0, 0, 0, 0},
		{"+half", 0.5, 0x3c00, 0x0000, 0x0000, 0x0000, 0x3c00},
		{"+inf", posInf, 0x7c00, 0x7c00, 0x7c00, 0x7c00, 0x7c00},
		{"epsilon", 0x1p-25, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000}, // underflow to zero
		{"epsilon", 0x1p-26, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000}, // underflow to zero
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := Float16FromFloat(tt.in)

			testRounding(t, f.bits, tt.nearTieAway, tt.nearTieEven, tt.toNeg, tt.toZero, tt.toPos)
		})
	}

	// test all mantissa
	for m := uint16(0); false && m < 1<<10; m++ {
		m := m
		n := m | 0x8000

		t.Run(fmt.Sprintf("%b", m), func(t *testing.T) {
			// test sub-normals
			if m != 0 {
				testRounding(t, m, 0x0000, 0x0000, 0x0000, 0x0000, 0x3c00)
				testRounding(t, n, 0x8000, 0x8000, 0xbc00, 0x8000, 0x8000)
			}

			for e := uint16(1); e < 0x0E; e++ { // test all exponents below one-half
				f := m | e<<10
				g := n | e<<10

				testRounding(t, f, 0x0000, 0x0000, 0x0000, 0x0000, 0x3c00)
				testRounding(t, g, 0x8000, 0x8000, 0xbc00, 0x8000, 0x8000)
			}

			if m != 0 { // test all exponents at one-half (but not one-half itself)
				f := m | 0xE<<10
				g := n | 0xE<<10

				testRounding(t, f, 0x3c00, 0x3c00, 0x0000, 0x0000, 0x3c00)
				testRounding(t, g, 0xbc00, 0xbc00, 0xbc00, 0x8000, 0x8000)
			}

			// test all infinities and NaNs
			mMax := m | 0x1F<<10
			nMax := n | 0x1F<<10

			testRounding(t, mMax, mMax, mMax, mMax, mMax, mMax)
			testRounding(t, nMax, nMax, nMax, nMax, nMax, nMax)
		})
	}

	for i := uint16(1); i < 1<<9; i++ { // test all exponents > abs(1), where nextUp(f) != f+0.5
		f := float32(i)

		m := Float16FromFloat(f).bits
		n := m | 0x8000

		mu := Float16FromFloat(f + 1).bits
		nd := mu | 0x8000

		md := Float16FromFloat(f - 1).bits
		nu := md | 0x8000

		t.Run(fmt.Sprintf("%f", f), func(t *testing.T) {
			testRounding(t, m, m, m, m, m, m)
			testRounding(t, n, n, n, n, n, n)

			testRounding(t, nextUp[binary16](m), m, m, m, m, mu)
			testRounding(t, nextDown[binary16](n), n, n, nd, n, n)

			testRounding(t, nextDown[binary16](m), m, m, md, md, m)
			testRounding(t, nextUp[binary16](n), n, n, n, nu, nu)
		})

		f += 0.5

		even := m
		if i&1 != 0 {
			// we are odd, so round to even will round up.
			even = mu
		}

		mh := Float16FromFloat(f).bits
		nh := mh | 0x8000

		t.Run(fmt.Sprintf("%f", f), func(t *testing.T) {
			testRounding(t, mh, mu, even, m, m, mu)
			testRounding(t, nh, nd, even|0x8000, nd, n, n)

			testRounding(t, nextUp[binary16](mh), mu, mu, m, m, mu)
			testRounding(t, nextDown[binary16](nh), nd, nd, nd, n, n)

			testRounding(t, nextDown[binary16](mh), m, m, m, m, mu)
			testRounding(t, nextUp[binary16](nh), n, n, nd, n, n)
		})
	}

	for i := uint16(1) << 9; i < 1<<10; i++ { // test the numbers where nextUp(f) == f+0.5
		f := float32(i)

		m := Float16FromFloat(f).bits
		n := m | 0x8000

		mu := Float16FromFloat(f + 1).bits
		nd := mu | 0x8000

		md := Float16FromFloat(f - 1).bits
		nu := md | 0x8000

		even := m
		evend := m
		if i&1 != 0 {
			// we are odd, so round to even will round up.
			even = mu
			evend = md
		}

		t.Run(fmt.Sprintf("%f", f), func(t *testing.T) {
			testRounding(t, m, m, m, m, m, m)
			testRounding(t, n, n, n, n, n, n)

			testRounding(t, nextUp[binary16](m), mu, even, m, m, mu)
			testRounding(t, nextDown[binary16](n), nd, even|0x8000, nd, n, n)

			testRounding(t, nextDown[binary16](m), m, evend, md, md, m)
			testRounding(t, nextUp[binary16](n), n, evend|0x8000, n, nu, nu)
		})
	}

	for i := uint16(1) << 10; i != 0; i++ { // test all the large integers now
		f := Float16FromFloat(float32(i))

		if uint16(f.Float32().Native()) != i {
			// we’ve lost the integer precision for this number, skip it.
			continue
		}

		m := f.bits
		n := m | 0x8000

		testRounding(t, m, m, m, m, m, m)
		testRounding(t, n, n, n, n, n, n)
	}
}

func TestFloat16OpExp2(t *testing.T) {
	type test struct {
		name string
		x    Float16
		bits uint16
	}

	tests := []test{
		{"2^e", E.Float16(), 0x4695},
		{"2^π", Pi.Float16(), 0x486a},
		{"2^φ", Phi.Float16(), 0x4224},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			res := f.Exp2()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.4e\nactual: %.4e\nexpect: %.4e", f, res, Float16{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, Float16{tt.bits})

			if bits&^0xFF != tt.bits&^0xFF {
				t.Logf("expect:\n%%e: %.4e\n%%b: 0b%b", Float16{tt.bits}, Float16{tt.bits})
				t.Errorf("2 ** Float16(%x)\n  actual: %04x\nexpected: %04x", tt.x, bits, tt.bits)
			}
		})
	}
}

func TestFloat16OpExp(t *testing.T) {
	type test struct {
		name string
		x    Float16
		bits uint16
	}

	tests := []test{
		{"e^e", E.Float16(), 0x4b94},
		{"e^π", Pi.Float16(), 0x4dc9},
		{"e^φ", Phi.Float16(), 0x450b},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := tt.x

			res := f.Exp()
			bits := res.Bits()

			t.Logf("%%e:\n     x: %.4e\nactual: %.4e\nexpect: %.4e", f, res, Float16{tt.bits})
			t.Logf("%%b:\n     x: %b\nactual: %b\nexpect: %b", f, res, Float16{tt.bits})

			if bits&^0xFF != tt.bits&^0xFF {
				t.Errorf("e ** Float16(%x)\n  actual: %04x\nexpected: %04x", tt.x, bits, tt.bits)
			}
		})
	}
}
