package posit

import (
	"fmt"
	"math"
	"math/bits"
)

type posit struct {
	num uint32
	es  uint8
}

func Newposit32FromBits(bits uint32, es uint8) posit {
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
	fmt.Println(regime, eseed, m)
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

	aexp := (a.num << (an + 1)) >> (32 - a.es)
	afrac_2 := (a.num << (a.es + an)) | (0b1 << 31)
	ascale := (int32(m) << a.es) + int32(aexp)

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

	bexp := (b.num << (bn + 1)) >> (32 - b.es)
	bfrac_2 := (b.num << (b.es + bn)) | (0b1 << 31)
	bscale := (int32(m) << b.es) + int32(bexp)

	//Out

	// We are using a trick to not have to double-negate.
	// The simple way to do this would be to negate af and fb
	// on aneg and bneg respectively, and negate the result if the result should be negated.
	// Instead we take advantage that a>b, so we never have to negate a and
	// we can negate b if and only if bneg^aneg.

	combinedFrac := (uint64(afrac_2) << 31)

	if ascale-bscale < 31 {
		bf := uint64(bfrac_2) << (31 - ascale + bscale)
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
	overflow := uint8(combinedFrac >> 63)
	ascale += int32(overflow)
	combinedFrac >>= overflow // leave the hidden bit in

	endm := ascale >> int32(a.es)
	endexp := uint32(ascale - (endm << a.es))

	var outPosit uint32
	var outn uint8
	if endm < 0 {
		outn = uint8(1 - endm)
		outPosit |= (0b1 << 31) >> outn
	} else {
		outn = uint8(2 + endm)
		outPosit = 0x7fffffff - 0xffffffff>>(outn)
	}
	//Recalculate the final fraction bits so that it matches the new exponent and m
	outFracBits := 31 - a.es - outn

	combinedFrac >>= (a.es + outn - 1)

	var toadd uint8
	var outfrac uint32
	if outn-1 <= 32-a.es {
		y_1 := ((combinedFrac) & 0x7FFF_FFFF) != 0
		var y uint8
		if y_1 {
			y = 1
		}
		outfrac = uint32(combinedFrac >> 31)

		x := uint8(outfrac & 1)
		outfrac >>= 1
		toadd = x & (y | uint8(outfrac&1))
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
		for a.num&(1<<(30-an)) == 0 {
			an++
		}
		m = int8(-an)
	} else {
		for a.num&(1<<(30-an)) != 0 {
			an++
		}
		m = int8(an - 1)
	}
	afracBits := (31 - a.es - 1 - an) //we need to bitshift by the remaining bits. this is 31 - 1 - n - es
	aexp := (uint32(a.num) & (0b00111111111111111111111111111111 >> an)) >> uint32(afracBits)

	afrac_2 := (a.num & ((0b00111111111111111111111111111111) >> (an + a.es)))
	afracP1 := afrac_2 | (0b1 << afracBits)
	ascale := (1<<int16(a.es))*int32(m) + int32(aexp)

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

	bfracBits := (31 - b.es - 1 - bn) //we need to bitshift by the remaining bits. this is 31 - 1 - n - es
	bexp := (uint32(b.num) & (0b00111111111111111111111111111111 >> bn)) >> uint32(bfracBits)

	bfrac_2 := (b.num & (0b00111111111111111111111111111111 >> (bn + b.es)))
	bfracP1 := bfrac_2 | (0b1 << bfracBits)
	bscale := (1<<int16(b.es))*int32(m) + int32(bexp)

	//Out
	var endScale = ascale + bscale
	afracP1 <<= 31 - afracBits
	bfracP1 <<= 31 - bfracBits

	af := int64(afracP1)
	bf := int64(bfracP1)
	combinedFrac := uint64(af * bf)
	neg := (aneg || bneg) && (!aneg || !bneg)
	if combinedFrac == 0 {
		a.num = 0
		return a
	}
	// combinedFrac looks like:
	// 1 bit for overflow 1 for hidden bit 62 for number

	overflow := uint8((combinedFrac >> (63)) & 1)
	endScale += int32(overflow)
	combinedFrac >>= overflow // leave the hidden bit in

	endexp, endm := splitExponent(endScale, a.es)

	var outPosit uint32
	var outn uint8
	if endm < 0 {
		outPosit |= ((0b1 << 31) >> (1 - endm))
		outn = uint8(1 - endm)
	} else {
		outPosit = 0x7fffffff
		outPosit ^= ((0b1 << 31) >> (2 + endm))
		outn = uint8(2 + endm)
	}
	//Recalculate the final fraction bits so that it matches the new exponent and m
	outFracBits := 31 - a.es - outn

	// combinedFrac >>= (30 - outFracBits)
	combinedFrac >>= (a.es + outn - 1)

	var toadd uint32
	if outn-1 <= 32-a.es {
		y_1 := ((combinedFrac) & 0x7FFF_FFFF) != 0
		y := uint64(0)
		if y_1 {
			y = 1
		}
		combinedFrac >>= 31
		x := combinedFrac & 1
		combinedFrac >>= 1
		z := combinedFrac & 1
		toadd = uint32(x & (y | z))
	} else {
		combinedFrac >>= 32
	}

	combinedFrac &= 0xffffffff >> (32 - outFracBits)

	outPosit &^= 0xffffffff >> (1 + outn)
	outPosit |= uint32(endexp) << outFracBits
	outPosit |= uint32(combinedFrac)
	outPosit += toadd
	if neg {
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
	if a.num == 0 {
		return Newposit32FromBits(0, a.es)
	}
	if b.num == 0 {
		return Newposit32FromBits(1<<31, a.es)
	}
	if b.num == 1<<31 {
		return Newposit32FromBits(0, a.es)
	}

	//A
	aneg := a.num>>31 != 0
	if aneg {
		a.num = uint32(-int32(a.num))
	}

	an := uint8(1)
	var m int8
	if a.num&(1<<30) == 0 {
		for a.num&(1<<(30-an)) == 0 {
			an++
		}
		m = int8(-an)
	} else {
		for a.num&(1<<(30-an)) != 0 {
			an++
		}
		m = int8(an - 1)
	}
	afracBits := (31 - a.es - 1 - an) //we need to bitshift by the remaining bits. this is 31 - 1 - n - es
	aexp := (uint32(a.num) & (0b00111111111111111111111111111111 >> an)) >> uint32(afracBits)

	afrac_2 := (a.num & ((0b00111111111111111111111111111111) >> (an + a.es)))
	afracP1 := afrac_2 | (0b1 << afracBits)
	ascale := (1<<int16(a.es))*int32(m) + int32(aexp)

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

	bfracBits := (31 - b.es - 1 - bn) //we need to bitshift by the remaining bits. this is 31 - 1 - n - es
	bexp := (uint32(b.num) & (0b00111111111111111111111111111111 >> bn)) >> uint32(bfracBits)

	bfrac_2 := (b.num & (0b00111111111111111111111111111111 >> (bn + b.es)))
	bfracP1 := bfrac_2 | (0b1 << bfracBits)
	bscale := (1<<int16(b.es))*int32(m) + int32(bexp)

	//Out

	var endScale = ascale - bscale
	afracP1 <<= 31 - afracBits
	bfracP1 <<= 31 - bfracBits

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
		endScale -= 1
		combinedFrac <<= 1 // remove the hidden bit
	}

	endexp, endm := splitExponent(endScale, a.es)

	var outPosit uint32
	var outn uint8
	if endm < 0 {
		outPosit |= ((0b1 << 31) >> (1 - endm))
		outn = uint8(1 - endm)
	} else {
		outPosit = 0x7fffffff
		outPosit ^= ((0b1 << 31) >> (2 + endm))
		outn = uint8(2 + endm)
	}
	//Recalculate the final fraction bits so that it matches the new exponent and m
	outFracBits := 31 - a.es - outn
	// 	combinedFrac >>= (29 - outFracBits)
	combinedFrac >>= (a.es + outn - 2)

	var toadd uint32
	if outn-1 <= 32-a.es {
		x := combinedFrac & 1
		y_1 := ((combinedFrac) & 0x7FFF_FFFF) != 0
		y := uint64(0)
		if y_1 {
			y = 1
		}
		z := (combinedFrac >> (32) & 1)
		a_1 := rem != 0
		a := uint64(0)
		if a_1 {
			a = 1
		}
		toadd = uint32(x & (y | z | a))
	}
	combinedFrac >>= 1
	combinedFrac &= 0xffffffff >> (32 - outFracBits)

	outPosit &^= 0xffffffff >> (1 + outn)
	outPosit |= uint32(endexp) << outFracBits
	outPosit |= uint32(combinedFrac)
	outPosit += toadd
	if neg {
		outPosit = uint32(-int32(outPosit))
	}
	return posit{
		num: outPosit,
		es:  a.es,
	}
}

func splitExponent(scale int32, es uint8) (endexp uint32, endm int32) {
	endexp1 := (scale % (1 << int16(es)))
	endm = scale / (1 << int16(es))

	if endexp1 < 0 {
		endexp = uint32(endexp1 + (1 << int16(es)))
		endm--
	} else {
		endexp = uint32(endexp1)
	}
	return
}
