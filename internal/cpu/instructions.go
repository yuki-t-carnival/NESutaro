package cpu

func (c *CPU) adc(addr uint16) {
	a := c.a
	mem := c.read(addr)
	result16 := uint16(a) + uint16(mem) + uint16(c.p&CarryFlagMask)
	result := byte(result16)

	if result16 > 0xFF {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}

	if (result^a)&(result^mem)&0x80 != 0 {
		c.p |= OverflowFlagMask
	} else {
		c.p &^= OverflowFlagMask
	}

	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.a = result
}

func (c *CPU) ahx(addr uint16) {
	hi := byte(addr>>8) + 1
	c.write(addr, c.a&c.x&hi)
}

func (c *CPU) alr(addr uint16) {
	c.and(addr)
	c.lsrAccum()
}

func (c *CPU) arr(addr uint16) {
	c.and(addr)
	c.rorAccum()

	// ?
	bit6 := c.a >> 6 & 1
	bit5 := c.a >> 5 & 1
	if bit6^bit5 == 1 {
		c.p |= OverflowFlagMask
	} else {
		c.p &^= OverflowFlagMask
	}
	if bit6 == 1 {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
}

func (c *CPU) anc(addr uint16) {
	c.and(addr)

	if c.a&0x80 != 0 {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
}

func (c *CPU) and(addr uint16) {
	result := c.a & c.read(addr)

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}

	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.a = result
}

func (c *CPU) asl(addr uint16) {
	val := c.read(addr)

	if val&0x80 != 0 {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}

	result := val << 1

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.write(addr, result)
}

func (c *CPU) aslAccum() {
	val := c.a

	if val&0x80 != 0 {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}

	result := val << 1

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.a = result
}

// ?
func (c *CPU) axs(addr uint16) {
	imm := c.read(addr)
	t := c.a & c.x
	result := t - imm

	if t >= imm {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.x = result
}

func (c *CPU) bcc(addr uint16) bool {
	isBranched := c.p&CarryFlagMask == 0
	if isBranched {
		c.pc = addr
	}
	return isBranched
}

func (c *CPU) bcs(addr uint16) bool {
	isBranched := c.p&CarryFlagMask != 0
	if isBranched {
		c.pc = addr
	}
	return isBranched
}

func (c *CPU) beq(addr uint16) bool {
	isBranched := c.p&ZeroFlagMask != 0
	if isBranched {
		c.pc = addr
	}
	return isBranched
}

func (c *CPU) bit(addr uint16) {
	mem := c.read(addr)
	result := c.a & mem

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if mem&0x40 != 0 {
		c.p |= OverflowFlagMask
	} else {
		c.p &^= OverflowFlagMask
	}
	if mem&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}
}

func (c *CPU) bmi(addr uint16) bool {
	isBranched := c.p&NegativeFlagMask != 0
	if isBranched {
		c.pc = addr
	}
	return isBranched
}

func (c *CPU) bne(addr uint16) bool {
	isBranched := c.p&ZeroFlagMask == 0
	if isBranched {
		c.pc = addr
	}
	return isBranched
}

func (c *CPU) bpl(addr uint16) bool {
	isBranched := c.p&NegativeFlagMask == 0
	if isBranched {
		c.pc = addr
	}
	return isBranched
}

func (c *CPU) brk() {
	c.fetch()
	pushLo := byte(c.pc & 0x00FF)
	pushHi := byte(c.pc >> 8)
	c.write(0x0100+uint16(c.s), pushHi)
	c.s -= 1
	c.write(0x0100+uint16(c.s), pushLo)
	c.s -= 1
	c.write(0x0100+uint16(c.s), c.p|0x30)
	c.s -= 1
	lo := uint16(c.read(0xFFFE))
	hi := uint16(c.read(0xFFFF))
	c.pc = hi<<8 + lo

	c.p |= InterruptDisableFlagMask
}

func (c *CPU) bvc(addr uint16) bool {
	isBranched := c.p&OverflowFlagMask == 0
	if isBranched {
		c.pc = addr
	}
	return isBranched
}

func (c *CPU) bvs(addr uint16) bool {
	isBranched := c.p&OverflowFlagMask != 0
	if isBranched {
		c.pc = addr
	}
	return isBranched
}

func (c *CPU) clc() {
	c.p &^= CarryFlagMask
}

func (c *CPU) cld() {
	c.p &^= DecimalFlagMask
}

func (c *CPU) cli() {
	c.p &^= InterruptDisableFlagMask
}

func (c *CPU) clv() {
	c.p &^= OverflowFlagMask
}

func (c *CPU) cmp(addr uint16) {
	a := c.a
	mem := c.read(addr)
	result := a - mem

	if a >= mem {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
	if a == mem {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}
}

func (c *CPU) cpx(addr uint16) {
	x := c.x
	mem := c.read(addr)
	result := x - mem

	if x >= mem {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
	if x == mem {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}
}

func (c *CPU) cpy(addr uint16) {
	y := c.y
	mem := c.read(addr)
	result := y - mem

	if y >= mem {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
	if y == mem {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}
}

func (c *CPU) dcp(addr uint16) {
	c.dec(addr)
	c.cmp(addr)
}

func (c *CPU) dec(addr uint16) {
	result := c.read(addr) - 1

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}
	c.write(addr, result)
}

func (c *CPU) dex() {
	result := c.x - 1

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.x = result
}

func (c *CPU) dey() {
	result := c.y - 1

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.y = result
}

func (c *CPU) eor(addr uint16) {
	result := c.a ^ c.read(addr)

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.a = result
}

func (c *CPU) inc(addr uint16) {
	result := c.read(addr) + 1

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.write(addr, result)
}

func (c *CPU) inx() {
	result := c.x + 1

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.x = result
}

func (c *CPU) iny() {
	result := c.y + 1

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.y = result
}

func (c *CPU) isc(addr uint16) {
	c.inc(addr)
	c.sbc(addr)
}

func (c *CPU) jmp(addr uint16) {
	c.pc = addr
}

func (c *CPU) jsr(addr uint16) {
	pushAddr := c.pc - 1
	hi := byte(pushAddr >> 8)
	lo := byte(pushAddr & 0x00FF)

	c.write(0x0100+uint16(c.s), hi)
	c.s -= 1

	c.write(0x0100+uint16(c.s), lo)
	c.s -= 1

	c.pc = addr
}

func (c *CPU) las(addr uint16) {
	result := c.read(addr) & c.s
	c.a = result
	c.x = result
	c.s = result
}

func (c *CPU) laxAddr(addr uint16) {
	c.lda(addr)
	c.ldx(addr)
}

func (c *CPU) laxImm(addr uint16) {
	c.lda(addr)
	c.tax()
}

func (c *CPU) lda(addr uint16) {
	result := c.read(addr)

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.a = result
}

func (c *CPU) ldx(addr uint16) {
	result := c.read(addr)

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.x = result
}

func (c *CPU) ldy(addr uint16) {
	result := c.read(addr)

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.y = result
}

func (c *CPU) lsr(addr uint16) {
	val := c.read(addr)
	result := val >> 1

	if val&0x01 != 0 {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	c.p &^= NegativeFlagMask

	c.write(addr, result)
}

func (c *CPU) lsrAccum() {
	val := c.a
	result := val >> 1

	if val&0x01 != 0 {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	c.p &^= NegativeFlagMask

	c.a = result
}

func (c *CPU) nop() {
}

func (c *CPU) ora(addr uint16) {
	result := c.a | c.read(addr)

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.a = result
}

func (c *CPU) pha() {
	c.write(0x0100+uint16(c.s), c.a)
	c.s -= 1
}

func (c *CPU) php() {
	c.write(0x0100+uint16(c.s), c.p|(0b00110000))
	c.s -= 1
}

func (c *CPU) pla() {
	c.s += 1
	result := c.read(0x0100 + uint16(c.s))

	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.a = result
}

func (c *CPU) plp() {
	c.s += 1
	result := c.read(0x0100+uint16(c.s)) & 0xCF
	if result&0x04 != c.p&0x04 {
		c.isIFlagToggleDelayed = true
	}
	c.p = c.p&0x34 | result&0xCB
}

func (c *CPU) rla(addr uint16) {
	c.rol(addr)
	c.and(addr)
}

func (c *CPU) rol(addr uint16) {
	val := c.read(addr)
	result := val<<1 + c.p&CarryFlagMask

	if val&0x80 != 0 {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.write(addr, result)
}

func (c *CPU) rolAccum() {
	val := c.a
	result := val<<1 + c.p&CarryFlagMask

	if val&0x80 != 0 {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.a = result
}

func (c *CPU) ror(addr uint16) {
	val := c.read(addr)
	carry := c.p & 0x01 << 7
	result := carry + val>>1

	if val&0x01 != 0 {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.write(addr, result)
}

func (c *CPU) rorAccum() {
	val := c.a
	carry := c.p & 0x01 << 7
	result := carry + val>>1

	if val&0x01 != 0 {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.a = result
}

func (c *CPU) rra(addr uint16) {
	c.ror(addr)
	c.adc(addr)
}

func (c *CPU) rti() {
	c.s += 1
	status := c.read(0x0100+uint16(c.s)) & 0xCF
	c.p = c.p&0x30 | status

	c.s += 1
	lo := uint16(c.read(0x0100 + uint16(c.s)))

	c.s += 1
	hi := uint16(c.read(0x0100 + uint16(c.s)))

	c.pc = hi<<8 + lo
}

func (c *CPU) rts() {
	c.s += 1
	lo := uint16(c.read(0x0100 + uint16(c.s)))

	c.s += 1
	hi := uint16(c.read(0x0100 + uint16(c.s)))

	c.pc = hi<<8 + lo + 1
}

func (c *CPU) sax(addr uint16) {
	c.write(addr, c.a&c.x)
}

func (c *CPU) sbc(addr uint16) {
	mem := c.read(addr)
	var carryNot byte
	if c.p&CarryFlagMask == 0 {
		carryNot = 1
	} else {
		carryNot = 0
	}
	resultInt := int16(c.a) - int16(mem) - int16(carryNot)
	result := c.a - mem - carryNot

	if !(resultInt < 0) {
		c.p |= CarryFlagMask
	} else {
		c.p &^= CarryFlagMask
	}
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if (result^c.a)&(result^^mem)&0x80 != 0 {
		c.p |= OverflowFlagMask
	} else {
		c.p &^= OverflowFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}

	c.a = result
}

func (c *CPU) sec() {
	c.p |= CarryFlagMask
}

func (c *CPU) sed() {
	c.p |= DecimalFlagMask
}

func (c *CPU) sei() {
	if c.p&InterruptDisableFlagMask == 0 {
		c.isIFlagToggleDelayed = true
	}
}

func (c *CPU) shx(addr uint16) {
	hi := byte(addr>>8) + 1
	c.write(addr, c.x&hi)
}

func (c *CPU) shy(addr uint16) {
	hi := byte(addr>>8) + 1
	c.write(addr, c.y&hi)
}

func (c *CPU) slo(addr uint16) {
	c.asl(addr)
	c.ora(addr)
}

func (c *CPU) sre(addr uint16) {
	c.lsr(addr)
	c.eor(addr)
}

func (c *CPU) sta(addr uint16) {
	c.write(addr, c.a)
}

func (c *CPU) stx(addr uint16) {
	c.write(addr, c.x)
}

func (c *CPU) sty(addr uint16) {
	c.write(addr, c.y)
}

func (c *CPU) tas(addr uint16) {
	c.s = c.a & c.x
	hi := byte(addr>>8) + 1
	c.write(addr, c.s&hi)
}

func (c *CPU) tax() {
	result := c.a
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}
	c.x = result
}

func (c *CPU) tay() {
	result := c.a
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}
	c.y = result
}

func (c *CPU) tsx() {
	result := c.s
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}
	c.x = result
}

func (c *CPU) txa() {
	result := c.x
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}
	c.a = result
}

func (c *CPU) txs() {
	result := c.x
	c.s = result
}

func (c *CPU) tya() {
	result := c.y
	if result == 0 {
		c.p |= ZeroFlagMask
	} else {
		c.p &^= ZeroFlagMask
	}
	if result&0x80 != 0 {
		c.p |= NegativeFlagMask
	} else {
		c.p &^= NegativeFlagMask
	}
	c.a = result
}

func (c *CPU) xaa(addr uint16) {
	c.txa()
	c.and(addr)
}
