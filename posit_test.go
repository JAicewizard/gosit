package posit

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/cjdelisle/goposit"
)

func TestGetFloat(t *testing.T) {
	posit := posit{
		num: 0b0_0001_101_11011101_00000000_00000000,
		es:  3,
	}
	gopos := goposit.NewSlowPosit(16, 3)
	gopos.Bits = big.NewInt(int64(0b1_0001_101_11011101))
	fmt.Println(gopos.ToFloat())

	if math.Abs(Getfloat(posit)-3.553926944732666e-06) > 0.000001 {
		t.Fatal("got the wrong float:", Getfloat(posit), "expected:", 3.553926944732666e-06)
	}
}

func TestGetFloat2(t *testing.T) {
	posit := posit{
		num: 0b1_00001_10_11011101_00000000_00000000,
		es:  2,
	}
	gopos := goposit.NewSlowPosit(32, 2)
	gopos.Bits = big.NewInt(int64(posit.num))
	fmt.Println(gopos.ToFloat())

	if math.Abs(Getfloat(posit)+9312.000066757206) > 0.000001 {
		t.Fatal("got the wrong float:", Getfloat(posit), "expected:", -9312.000066757206)
	}
}

func TestGetFloat3(t *testing.T) {
	posit := posit{
		num: 0b0_00001_11_101110100000000000000000,
		es:  2,
	}
	gopos := goposit.NewSlowPosit(32, 2)
	gopos.Bits = big.NewInt(int64(posit.num))
	fmt.Println(gopos.ToFloat())

	if math.Abs(Getfloat(posit)-0.00021076202921221953) > 0.000001 {
		t.Fatal("got the wrong float:", Getfloat(posit), "expected:", 0.00021076202921221953)
	}
}

type opTest struct {
	a   posit
	b   posit
	exp posit
}

var addtests = map[string]opTest{
	"1": {
		a: posit{
			num: 0b0_00001_10_11011101_00000000_00000000,
			es:  2,
		},
		b: posit{
			num: 0b0_00001_10_11011101_00000000_00000000,
			es:  2,
		},
		exp: posit{
			num: 0b0_00001_11_11011101_00000000_00000000,
			es:  2,
		},
	},
	"2": {
		a: posit{
			num: 0b0_00001_10_11011101_00000000_00000000,
			es:  2,
		},
		b: posit{
			num: 0b0_00001_10_11111101_00000000_00000000,
			es:  2,
		},
		exp: posit{
			num: 0b0_00001_11_11101101_00000000_00000000,
			es:  2,
		},
	},
	"3": {
		a: posit{
			num: 0b1_11110_10_00100000_00000000_00000000,
			es:  2,
		},
		b: posit{
			num: 0b0_00001_10_11111101_00000000_00000000,
			es:  2,
		},
		exp: posit{
			num: 0b0_00001_10_00001101_00000000_00000000,
			es:  2,
		},
	},
	"4:add1-0": {
		a: posit{
			num: 0b01000101100011011010110001101011,
			es:  2,
		},
		b: posit{
			num: 0b10000100111000010000001001010100,
			es:  2,
		},
		exp: posit{
			num: 0b10000100111000010000010110110111,
			es:  2,
		},
	},
	"5:add1-1": {
		a: posit{
			num: 0b11001010100000010111011101000100,
			es:  2,
		},
		b: posit{
			num: 0b11100010010101101111000001110010,
			es:  2,
		},
		exp: posit{
			num: 0b11001001000101110011001101100000,
			es:  2,
		},
	},
	"6:add1-2": {
		a: posit{
			num: 0b10110001010101010011110101100011,
			es:  2,
		},
		b: posit{
			num: 0b01000000010010000100010011001001,
			es:  2,
		},
		exp: posit{
			num: 0b10110101011110010101111111001000,
			es:  2,
		},
	},
	"7:add1-3": {
		a: posit{
			num: 0b00001010100110011011111111110101,
			es:  2,
		},
		b: posit{
			num: 0b10011000110111111000101011101110,
			es:  2,
		},
		exp: posit{
			num: 0b10011000110111111001000000100001,
			es:  2,
		},
	},
	"8": {
		a: posit{
			num: 0b01011111100111100110110010101001,
			es:  2,
		},
		b: posit{
			num: 0b01010101011001000101111011111011,
			es:  2,
		},
		exp: posit{
			num: 0b01100001100101000010011100001010,
			es:  2,
		},
	},
	"9:add-4": {
		a: posit{
			num: 0b11001110011111001011111010111100,
			es:  2,
		},
		b: posit{
			num: 0b01110110011100110110110110010110,
			es:  2,
		},
		exp: posit{
			num: 0b01110110011100110101101010001111,
			es:  2,
		},
	},
	"10:negativeShift-0": {
		a: posit{
			num: 0b01111111111101010001101111000111,
			es:  2,
		},
		b: posit{
			num: 0b01100001110010010010111000000100,
			es:  2,
		},
		exp: posit{
			num: 0b01111111111101010001101111000111,
			es:  2,
		},
	},
	"11": {
		a: posit{
			num: 0b01111111111101010001101111000111,
			es:  2,
		},
		b: posit{
			num: 0b00000000000000000000000000000000,
			es:  2,
		},
		exp: posit{
			num: 0b01111111111101010001101111000111,
			es:  2,
		},
	},
	"12": {
		a: posit{
			num: 0b01111111111101010001101111000111,
			es:  2,
		},
		b: posit{
			num: 0b10000000000000000000000000000000,
			es:  2,
		},
		exp: posit{
			num: 0b10000000000000000000000000000000,
			es:  2,
		},
	},
	"13:restrictiveMask-0": {
		a: posit{
			num: 0b10000000000000000110111001011100,
			es:  2,
		},
		b: posit{
			num: 0b11111111111100101111000101010110,
			es:  2,
		},
		exp: posit{
			num: 0b10000000000000000110111001011100,
			es:  2,
		},
	},
}

func TestAdd(t *testing.T) {
	for name, test := range addtests {
		t.Run(name, func(t *testing.T) {
			res := AddPositSameES(test.a, test.b)
			t.Logf("exp:%#032b", test.exp.num)
			t.Logf("res:%#032b", res.num)
			if res != test.exp {
				t.Fatal("got the wrong number:", Getfloat(res), "expected:", Getfloat(test.exp))
			}
			res = AddPositSameES(test.b, test.a)
			t.Logf("res:%#032b", res.num)
			if res != test.exp {
				t.Fatal("got the wrong float:", Getfloat(res), "expected:", Getfloat(test.exp))
			}
		})
	}
}

var subtests = map[string]opTest{
	"1:zeroRes": {
		a: posit{
			num: 0b0_00001_10_11011101_00000000_00000000,
			es:  2,
		},
		b: posit{
			num: 0b0_00001_10_11011101_00000000_00000000,
			es:  2,
		},
		exp: posit{
			num: 0b0_0000000_00000000_00000000_00000000,
			es:  2,
		},
	},
	"2": {
		a: posit{
			num: 0b0_00001_10_11011101_00000000_00000000,
			es:  2,
		},
		b: posit{
			num: 0b0_00001_10_11111101_00000000_00000000,
			es:  2,
		},
		exp: posit{
			num: 0b1_111110_01_0000000_00000000_00000000,
			es:  2,
		},
	},
	"3": {
		a: posit{
			num: 0b1_11110_10_00100000_00000000_00000000,
			es:  2,
		},
		b: posit{
			num: 0b0_00001_10_11111101_00000000_00000000,
			es:  2,
		},
		exp: posit{
			num: 0b1_11110_00_10001001_10000000_00000000,
			es:  2,
		},
	},
}

func TestSub(t *testing.T) {
	for name, test := range subtests {
		t.Run(name, func(t *testing.T) {
			res := SubPositSameES(test.a, test.b)
			t.Logf("%#032b", res.num)
			if res != test.exp {
				t.Fatal("got the wrong number:", Getfloat(res), "expected:", Getfloat(test.exp))
			}
		})
	}
}

var multests = map[string]opTest{
	"1": {
		a: posit{
			num: 0b00011011100000110111111010010111,
			es:  2,
		},
		b: posit{
			num: 0b00000011100000010000011110100011,
			es:  2,
		},
		exp: posit{
			num: 0b00000001011110010010111110000101,
			es:  2,
		},
	},
	"2": {
		a: posit{
			num: 0b10100110111101110111011101100101,
			es:  2,
		},
		b: posit{
			num: 0b01100011011100100011011001101100,
			es:  2,
		},
		exp: posit{
			num: 0b10001111111001011110010001011011,
			es:  2,
		},
	},
	"3:mul0": {
		a: posit{
			num: 0b00000000000000000000000000000000,
			es:  2,
		},
		b: posit{
			num: 0b01100011011100100011011001101100,
			es:  2,
		},
		exp: posit{
			num: 0b00000000000000000000000000000000,
			es:  2,
		},
	},
	"4:mulInf": {
		a: posit{
			num: 0b00000000000000000000000000000000,
			es:  2,
		},
		b: posit{
			num: 0b01100011011100100011011001101100,
			es:  2,
		},
		exp: posit{
			num: 0b00000000000000000000000000000000,
			es:  2,
		},
	},
	"5": {
		a: posit{
			num: 0b00111000000000000000000000000000,
			es:  2,
		},
		b: posit{
			num: 0b01000000000000000000000000000000,
			es:  2,
		},
		exp: posit{
			num: 0b00111000000000000000000000000000,
			es:  2,
		},
	},
}

func TestMul(t *testing.T) {
	for name, test := range multests {
		t.Run(name, func(t *testing.T) {
			res := MulPositSameES(test.a, test.b)
			t.Logf("exp:%#032b", test.exp.num)
			t.Logf("res:%#032b", res.num)
			if res != test.exp {
				t.Fatal("got the wrong number:", Getfloat(res), "expected:", Getfloat(test.exp))
			}
			res = MulPositSameES(test.b, test.a)
			t.Logf("res:%#032b", res.num)
			if res != test.exp {
				t.Fatal("got the wrong float:", Getfloat(res), "expected:", Getfloat(test.exp))
			}
		})
	}
}

var divtests = map[string]opTest{
	"1": {
		a: posit{
			num: 0b00010001101100111010100101100001,
			es:  2,
		},
		b: posit{
			num: 0b01100100000010101011101101110000,
			es:  2,
		},
		exp: posit{
			num: 0b00000111011010010010000101101010,
			es:  2,
		},
	},
	"2:add1-0": {
		a: posit{
			num: 0b01111011001001100110000001010001,
			es:  2,
		},
		b: posit{
			num: 0b10000101010101010110010110101001,
			es:  2,
		},
		exp: posit{
			num: 0b10111100111101011001000011001110,
			es:  2,
		},
	},
	"3:div0": {
		a: posit{
			num: 0b00000000000000000000000000000000,
			es:  2,
		},
		b: posit{
			num: 0b01100011011100100011011001101100,
			es:  2,
		},
		exp: posit{
			num: 0b00000000000000000000000000000000,
			es:  2,
		},
	},
	"3:divbyInf": {
		a: posit{
			num: 0b01100011011100100011011001101100,
			es:  2,
		},
		b: posit{
			num: 0b10000000000000000000000000000000,
			es:  2,
		},
		exp: posit{
			num: 0b00000000000000000000000000000000,
			es:  2,
		},
	},
	"3:divby0": {
		a: posit{
			num: 0b01100011011100100011011001101100,
			es:  2,
		},
		b: posit{
			num: 0b00000000000000000000000000000000,
			es:  2,
		},
		exp: posit{
			num: 0b10000000000000000000000000000000,
			es:  2,
		},
	},
	"4:0sizedfrac": {
		a: posit{
			num: 0b00000000000000000000000001010011,
			es:  2,
		},
		b: posit{
			num: 0b01111010110110111000010110011101,
			es:  2,
		},
		exp: posit{
			num: 0b00000000000000000000000000000110,
			es:  2,
		},
	},
}

func TestDiv(t *testing.T) {
	for name, test := range divtests {
		t.Run(name, func(t *testing.T) {
			res := DivPositSameES(test.a, test.b)
			t.Logf("exp:%#032b", test.exp.num)
			t.Logf("res:%#032b", res.num)
			if res != test.exp {
				t.Fatal("got the wrong number:", Getfloat(res), "expected:", Getfloat(test.exp))
			}
		})
	}
}

var sqrttests = map[string]opTest{
	"1": {
		a: posit{
			num: 0b01001000000000000000000000000000,
			es:  2,
		},
		b: posit{},
		exp: posit{
			num: 0b01000011010100000100111100110011,
			es:  2,
		},
	},
	"2": {
		a: posit{
			num: 0b01010000000000000000000000000000,
			es:  2,
		},
		b: posit{},
		exp: posit{
			num: 0b01001000000000000000000000000000,
			es:  2,
		},
	},
	"3:add-0": {
		a: posit{
			num: 0b01110110100100110101001011011101,
			es:  2,
		},
		b: posit{},
		exp: posit{
			num: 0b01100110011010110101100011011110,
			es:  2,
		},
	},
	"4:negative": {
		a: posit{
			num: 0b11001011100010001011101001110010,
			es:  2,
		},
		b: posit{},
		exp: posit{
			num: 0b10000000000000000000000000000000,
			es:  2,
		},
	},
}

func TestSqrt(t *testing.T) {
	for name, test := range sqrttests {
		t.Run(name, func(t *testing.T) {
			res := SqrtPosit(test.a)
			t.Logf("exp:%#032b", test.exp.num)
			t.Logf("res:%#032b", res.num)
			if res != test.exp {
				t.Fatal("got the wrong number:", Getfloat(res), "expected:", Getfloat(test.exp))
			}
		})
	}
}

type opBenchCase struct {
	a  posit
	b  posit
	ag goposit.Posit32
	bg goposit.Posit32
}

var slowBenchcases = [...]opBenchCase{
	{
		a: posit{
			num: 0b0_00001_10_11011101_00000000_00000000,
			es:  2,
		},
		b: posit{
			num: 0b0_00001_10_11011101_00000000_00000000,
			es:  2,
		},
		ag: goposit.NewPosit32().SetBits(0b0_00001_10_11011101_00000000_00000000),
		bg: goposit.NewPosit32().SetBits(0b0_00001_10_11011101_00000000_00000000),
	},
	{
		a: posit{
			num: 0b0_00001_10_11011101_00000000_00000000,
			es:  2,
		},
		b: posit{
			num: 0b0_00001_10_11111101_00000000_00000000,
			es:  2,
		},
		ag: goposit.NewPosit32().SetBits(0b0_00001_10_11011101_00000000_00000000),
		bg: goposit.NewPosit32().SetBits(0b0_00001_10_11111101_00000000_00000000),
	},
	{
		a: posit{
			num: 0b1_11110_10_00100000_00000000_00000000,
			es:  2,
		},
		b: posit{
			num: 0b0_00001_10_11111101_00000000_00000000,
			es:  2,
		},
		ag: goposit.NewPosit32().SetBits(0b1_11110_10_00100000_00000000_00000000),
		bg: goposit.NewPosit32().SetBits(0b0_00001_10_11111101_00000000_00000000),
	},
	{
		a: posit{
			num: 0b01000101100011011010110001101011,
			es:  2,
		},
		b: posit{
			num: 0b10000100111000010000001001010100,
			es:  2,
		},
		ag: goposit.NewPosit32().SetBits(0b01000101100011011010110001101011),
		bg: goposit.NewPosit32().SetBits(0b10000100111000010000001001010100),
	},
	{
		a: posit{
			num: 0b11001010100000010111011101000100,
			es:  2,
		},
		b: posit{
			num: 0b11100010010101101111000001110010,
			es:  2,
		},
		ag: goposit.NewPosit32().SetBits(0b11001010100000010111011101000100),
		bg: goposit.NewPosit32().SetBits(0b11100010010101101111000001110010),
	},
	{
		a: posit{
			num: 0b10110001010101010011110101100011,
			es:  2,
		},
		b: posit{
			num: 0b01000000010010000100010011001001,
			es:  2,
		},
		ag: goposit.NewPosit32().SetBits(0b10110001010101010011110101100011),
		bg: goposit.NewPosit32().SetBits(0b01000000010010000100010011001001),
	},
	{
		a: posit{
			num: 0b00001010100110011011111111110101,
			es:  2,
		},
		b: posit{
			num: 0b10011000110111111000101011101110,
			es:  2,
		},
		ag: goposit.NewPosit32().SetBits(0b00001010100110011011111111110101),
		bg: goposit.NewPosit32().SetBits(0b10011000110111111000101011101110),
	},
	{
		a: posit{
			num: 0b01011111100111100110110010101001,
			es:  2,
		},
		b: posit{
			num: 0b01010101011001000101111011111011,
			es:  2,
		},
		ag: goposit.NewPosit32().SetBits(0b01011111100111100110110010101001),
		bg: goposit.NewPosit32().SetBits(0b01010101011001000101111011111011),
	},
	{
		a: posit{
			num: 0b11001110011111001011111010111100,
			es:  2,
		},
		b: posit{
			num: 0b01110110011100110110110110010110,
			es:  2,
		},
		ag: goposit.NewPosit32().SetBits(0b11001110011111001011111010111100),
		bg: goposit.NewPosit32().SetBits(0b01110110011100110110110110010110),
	},
	{
		a: posit{
			num: 0b01111111111101010001101111000111,
			es:  2,
		},
		b: posit{
			num: 0b01100001110010010010111000000100,
			es:  2,
		},
		ag: goposit.NewPosit32().SetBits(0b01111111111101010001101111000111),
		bg: goposit.NewPosit32().SetBits(0b01100001110010010010111000000100),
	},
}

func BenchmarkAddSlow(b *testing.B) {
	for n := 0; n < b.N; n++ {
		AddPositSameES(slowBenchcases[n%len(slowBenchcases)].a, slowBenchcases[n%len(slowBenchcases)].b)
	}
}
func BenchmarkAddSlowGoposit(b *testing.B) {
	for n := 0; n < b.N; n++ {
		slowBenchcases[n%len(slowBenchcases)].ag.Add(slowBenchcases[n%len(slowBenchcases)].bg)
	}
}
func BenchmarkMulSlow(b *testing.B) {
	for n := 0; n < b.N; n++ {
		MulPositSameES(slowBenchcases[n%len(slowBenchcases)].a, slowBenchcases[n%len(slowBenchcases)].b)
	}
}
func BenchmarkMulSlowGoposit(b *testing.B) {
	for n := 0; n < b.N; n++ {
		slowBenchcases[n%len(slowBenchcases)].ag.Mul(slowBenchcases[n%len(slowBenchcases)].bg)
	}
}
func BenchmarkDivSlow(b *testing.B) {
	for n := 0; n < b.N; n++ {
		DivPositSameES(slowBenchcases[n%len(slowBenchcases)].a, slowBenchcases[n%len(slowBenchcases)].b)
	}
}
func BenchmarkDivSlowGoposit(b *testing.B) {
	for n := 0; n < b.N; n++ {
		slowBenchcases[n%len(slowBenchcases)].ag.Div(slowBenchcases[n%len(slowBenchcases)].bg)
	}
}
func BenchmarkSqrtSlow(b *testing.B) {
	for n := 0; n < b.N; n++ {
		a := slowBenchcases[n%len(slowBenchcases)].a
		if a.num&(1<<31) != 0 {
			a.num = uint32(-int32(a.num))
		}
		SqrtPosit(a)
	}
}
