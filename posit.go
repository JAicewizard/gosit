package posit

import (
	"math"
	"math/bits"
)

type posit struct {
	num uint32
	es  uint8
}

func Newposit32FromBits(bits uint32, es uint8) posit {
	if es > 31 {
		panic("ES cannot be more then 31, that doesnt even make sense")
	}
	return posit{num: bits, es: es}
}

func Getfloat(posit posit) float64 {
	if posit.num == 0 {
		return 0.0
	}
	neg := posit.num>>31 != 0
	if neg {
		posit.num = uint32(-int32(posit.num))
	}
	eseed := 1 << (1 << uint16(posit.es))
	var n int8
	var m int8
	n = 1
	if posit.num&(1<<30) == 0 {
		for posit.num&(1<<(30-n)) == 0 {
			n++
		}
		m = -n
	} else {
		for posit.num&(1<<(30-n)) != 0 {
			n++
		}
		m = n - 1
	}
	regime := math.Pow(float64(eseed), float64(m))
	fracBits := (31 - posit.es - 1 - uint8(n)) //we need to bitshift by the remaining bits. this is 31 - 1 - n - es
	exp := (uint32(posit.num) & (0b00111111111111111111111111111111 >> n)) >> fracBits
	frac_2 := (posit.num & ((0b00111111111111111111111111111111) >> (n + int8(posit.es))))

	frac := 1 + float64(frac_2)/float64(uint32(0b1<<(fracBits))-1)

	if !neg {
		return regime * math.Pow(2.0, float64(exp)) * frac
	} else {
		return -(regime * float64(int32(1)<<exp) * frac)
	}
}
func Negate(p posit) posit {
	p.num = uint32(-int32(p.num))
	return p
}

func AddPositSameES(a, b posit) posit {
	if a.es != b.es {
		panic("es must be the same for AddPositFast")
	}
	if a.num == 0 || b.num == 0 {
		a.num |= b.num
		return a
	}
	if a.num == 1<<31 || b.num == 1<<31 {
		return Newposit32FromBits(1<<31, a.es)
	}
	aneg := a.num>>31 != 0
	if aneg {
		a.num = uint32(-int32(a.num))
	}

	xneg := aneg
	bneg := b.num>>31 != 0
	if bneg {
		xneg = !xneg
		b.num = uint32(-int32(b.num))
	}

	//This allows less branching later on.
	//This potentially leaves bneg in the wrong position, dont use bneg!!
	if a.num < b.num {
		a.num, b.num = b.num, a.num
		aneg = bneg
		// aneg, bneg = bneg, aneg
	} else if a.num == -b.num {
		a.num = 0
		return a
	}
	//A
	var an uint8
	var m int8
	if a.num&(1<<30) == 0 {
		an = uint8(bits.LeadingZeros32(a.num << 1))
		m = int8(-an)
	} else {
		an = uint8(bits.LeadingZeros32(-(a.num << 1)))
		m = int8(an - 1)
	}
	an++
	aexp := (a.num << ((an + 1) & 0x1f)) >> ((32 - a.es) & 0x1f)
	afrac_2 := (a.num << ((a.es + an) & 0x1f)) | (0b1 << 31)
	ascale := (int32(m) << (b.es & 0x1f)) + int32(aexp)

	//B
	var bn uint8
	if b.num&(1<<30) == 0 {
		bn = uint8(bits.LeadingZeros32(b.num << 1))
		m = int8(-bn)
	} else {
		bn = uint8(bits.LeadingZeros32(-(b.num << 1)))
		m = int8(bn - 1)
	}
	bn++
	bexp := (b.num << ((bn + 1) & 0x1f)) >> ((32 - a.es) & 0x1f)
	bfrac_2 := (b.num << ((a.es + bn) & 0x1f)) | (0b1 << 31)
	bscale := (int32(m) << (b.es & 0x1f)) + int32(bexp)

	//Out

	// We are using a trick to not have to double-negate.
	// The simple way to do this would be to negate af and fb
	// on aneg and bneg respectively, and negate the result if the result should be negated.
	// Instead we take advantage that a>b, so we never have to negate a and
	// we can negate b if and only if bneg^aneg.

	combinedFrac := (uint64(afrac_2) << 31)

	if ascale-bscale < 31 {
		bf := uint64(bfrac_2) << ((31 - ascale + bscale) & 0x1f)
		if xneg {
			combinedFrac -= bf
		} else {
			combinedFrac += bf
		}
	}

	// combinedFrac looks like:
	// 1 bit for overflow 1 for hidden bit 62 for number

	if xneg {
		//This is faster than LeadingZeros64
		for (combinedFrac>>62)&1 == 0 {
			ascale--
			combinedFrac <<= 1
		}
	}
	overflow := int32(combinedFrac >> 63)
	ascale += overflow
	combinedFrac >>= overflow & 0b1 // leave the hidden bit in

	endm := ascale >> a.es
	endexp := uint32(ascale - (endm << a.es))

	var outPosit uint32
	var outn uint8
	if endm < 0 {
		outn = uint8(1 - endm)
		outPosit = (0b1 << 31) >> (outn & 0x0f)
	} else {
		outn = uint8(2 + endm)
		outPosit = 0x7fffffff - 0xffffffff>>(outn&0x0f)
	}
	//Recalculate the final fraction bits so that it matches the new exponent and m
	outFracBits := 31 - a.es - outn

	combinedFrac >>= (a.es + outn - 1)

	// toadd should be uint32 since it trades one instruction
	// outside the if-statement for inside
	var toadd uint32
	var outfrac uint32
	if outn-1 <= 32-a.es {
		var y uint32
		if ((combinedFrac) & 0x7fffffff) != 0 {
			y = 1
		}
		outfrac = uint32(combinedFrac >> 31)
		x := uint32(outfrac) & 1
		outfrac >>= 1
		toadd = x & (y | uint32(outfrac&1))
		outfrac &= 0x7fffffff >> (a.es + outn)
	} else {
		outfrac = uint32(combinedFrac >> outFracBits)
		outfrac >>= 1 + a.es + outn
	}

	outPosit |= uint32(endexp) << outFracBits
	outPosit |= outfrac
	outPosit += uint32(toadd)
	if aneg {
		outPosit = uint32(-int32(outPosit))
	}
	return posit{
		num: outPosit,
		es:  a.es,
	}
}

func SubPositSameES(a, b posit) posit {
	b.num = uint32(-int32(b.num))
	return AddPositSameES(a, b)
}

func MulPositSameES(a, b posit) posit {
	if a.es != b.es {
		panic("es must be the same for AddPositFast")
	}
	if a.num == 0 || b.num == 0 {
		return Newposit32FromBits(0, a.es)
	}
	if a.num == 1<<31 || b.num == 1<<31 {
		return Newposit32FromBits(1<<31, a.es)
	}
	//A
	aneg := a.num>>31 != 0
	if aneg {
		a.num = uint32(-int32(a.num))
	}

	an := uint8(1)
	var m int8
	if a.num&(1<<30) == 0 {
		an = uint8(bits.LeadingZeros32(a.num << 1))
		m = int8(-an)
	} else {
		an = uint8(bits.LeadingZeros32(-(a.num << 1)))
		m = int8(an - 1)
	}
	an++
	aexp := (a.num << ((an + 1) & 0x1f)) >> ((32 - a.es) & 0x1f)
	afracP1 := (a.num << ((a.es + an) & 0x1f)) | (0b1 << 31)
	ascale := (int32(m) << (b.es & 0x1f)) + int32(aexp)

	//B
	xneg := aneg
	bneg := b.num>>31 != 0
	if bneg {
		xneg = !xneg
		b.num = uint32(-int32(b.num))
	}

	var bn uint8
	if b.num&(1<<30) == 0 {
		bn = uint8(bits.LeadingZeros32(b.num << 1))
		m = int8(-bn)
	} else {
		bn = uint8(bits.LeadingZeros32(-(b.num << 1)))
		m = int8(bn - 1)
	}
	bn++
	bexp := (b.num << ((bn + 1) & 0x1f)) >> ((32 - a.es) & 0x1f)
	bfracP1 := (b.num << ((a.es + bn) & 0x1f)) | (0b1 << 31)
	ascale += (int32(m) << (b.es & 0x1f)) + int32(bexp)

	//Out

	af := int64(afracP1)
	bf := int64(bfracP1)
	combinedFrac := uint64(af * bf)

	// combinedFrac looks like:
	// 1 bit for overflow 1 for hidden bit 62 for number

	overflow := uint8((combinedFrac >> 63))
	ascale += int32(overflow)
	combinedFrac >>= overflow // leave the hidden bit in

	endm := ascale >> a.es
	endexp := uint32(ascale - (endm << a.es))

	var outPosit uint32
	var outn uint8
	if endm < 0 {
		outn = uint8(1 - endm)
		outPosit = (0b1 << 31) >> (outn & 0x0f)
	} else {
		outn = uint8(2 + endm)
		outPosit = 0x7fffffff - 0xffffffff>>(outn&0x0f)
	}
	//Recalculate the final fraction bits so that it matches the new exponent and m
	outFracBits := 31 - a.es - outn

	// combinedFrac >>= (30 - outFracBits)
	combinedFrac >>= (a.es + outn - 1)

	var toadd uint32
	var outfrac uint32
	if outn-1 <= 32-a.es {
		var y uint32
		if ((combinedFrac) & 0x7fffffff) != 0 {
			y = 1
		}
		outfrac = uint32(combinedFrac >> 31)
		x := uint32(outfrac) & 1
		outfrac >>= 1
		toadd = x & (y | uint32(outfrac&1))
		outfrac &= 0x7fffffff >> (a.es + outn)
	} else {
		outfrac = uint32(combinedFrac >> outFracBits)
		outfrac >>= 1 + a.es + outn
	}

	outPosit |= uint32(endexp) << outFracBits
	outPosit |= outfrac
	outPosit += toadd
	if xneg {
		outPosit = uint32(-int32(outPosit))
	}
	return posit{
		num: outPosit,
		es:  a.es,
	}
}

func DivPositSameES(a, b posit) posit {
	if a.es != b.es {
		panic("es must be the same for AddPositFast")
	}
	if a.num == 0 || b.num == 1<<31 {
		return Newposit32FromBits(0, a.es)
	}
	if b.num == 0 {
		return Newposit32FromBits(1<<31, a.es)
	}

	//A
	aneg := a.num>>31 != 0
	if aneg {
		a.num = uint32(-int32(a.num))
	}

	var an uint8
	var m int8
	if a.num&(1<<30) == 0 {
		an = uint8(bits.LeadingZeros32(a.num << 1))
		m = int8(-an)
	} else {
		an = uint8(bits.LeadingZeros32(-(a.num << 1)))
		m = int8(an - 1)
	}
	an++
	aexp := (a.num << ((an + 1) & 0x1f)) >> ((32 - a.es) & 0x1f)
	afracP1 := (a.num << ((a.es + an) & 0x1f)) | (0b1 << 31)
	ascale := (int32(m) << (b.es & 0x1f)) + int32(aexp)

	//B
	bneg := b.num>>31 != 0
	if bneg {
		b.num = uint32(-int32(b.num))
	}

	bn := uint8(1)
	if b.num&(1<<30) == 0 {
		for b.num&(1<<(30-bn)) == 0 {
			bn++
		}
		m = int8(-bn)
	} else {
		for b.num&(1<<(30-bn)) != 0 {
			bn++
		}
		m = int8(bn - 1)
	}

	bn++
	bexp := (b.num << ((bn + 1) & 0x1f)) >> ((32 - a.es) & 0x1f)
	bfracP1 := (b.num << ((a.es + bn) & 0x1f)) | (0b1 << 31)
	ascale -= (int32(m) << (b.es & 0x1f)) + int32(bexp)

	//Out

	af := int64(afracP1)
	bf := int64(bfracP1)
	combinedFrac := uint64(af<<30) / (uint64(bf))
	rem := uint64(af<<30) % (uint64(bf))

	// combinedFrac := uint64(af * bf)
	neg := (aneg || bneg) && (!aneg || !bneg)
	if combinedFrac == 0 {
		a.num = 0
		return a
	}
	// combinedFrac looks like:
	// 1 bit for overflow 1 for hidden bit 62 for number

	overflow := (combinedFrac >> (30)) == 0
	if overflow {
		ascale -= 1
		combinedFrac <<= 1 // remove the hidden bit
	}

	endm := ascale >> a.es
	endexp := uint32(ascale - (endm << a.es))

	var outPosit uint32
	var outn uint8
	if endm < 0 {
		outn = uint8(1 - endm)
		outPosit = (0b1 << 31) >> (outn & 0x0f)
	} else {
		outn = uint8(2 + endm)
		outPosit = 0x7fffffff - 0xffffffff>>(outn&0x0f)
	}
	//Recalculate the final fraction bits so that it matches the new exponent and m
	outFracBits := 31 - a.es - outn
	combinedFrac >>= (a.es + outn - 2)

	var toadd uint32
	var outfrac uint32
	if outn-1 <= 32-a.es {
		x := uint32(combinedFrac & 1)
		y_1 := ((combinedFrac) & 0x7fffffff) != 0
		y := uint32(0)
		if y_1 {
			y = 1
		}
		z := uint32((combinedFrac >> (32) & 1))
		c_1 := rem != 0
		c := uint32(0)
		if c_1 {
			c = 1
		}
		toadd = uint32(x & (y | z | c))
	}

	g := uint8((1 + a.es + outn) & 0x1f)
	outfrac = uint32(combinedFrac>>1) << g
	outfrac >>= g

	outPosit &^= 0xffffffff >> (1 + outn)
	outPosit |= uint32(endexp) << outFracBits
	outPosit |= outfrac
	outPosit += toadd
	if neg {
		outPosit = uint32(-int32(outPosit))
	}
	return posit{
		num: outPosit,
		es:  a.es,
	}
}
