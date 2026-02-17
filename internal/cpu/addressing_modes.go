package cpu

func (c *CPU) immediate() uint16 {
	addr := c.pc
	c.pc += 1
	return addr
}

func (c *CPU) zeroPage() uint16 {
	zp := uint16(c.fetch())
	return zp
}

func (c *CPU) zeroPageX() uint16 {
	zp := uint16((c.fetch() + c.x) & 0xFF)
	return zp
}

func (c *CPU) zeroPageY() uint16 {
	zp := uint16((c.fetch() + c.y) & 0xFF)
	return zp
}

func (c *CPU) absolute() uint16 {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	return hi<<8 + lo
}

func (c *CPU) absoluteX() (uint16, int) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	base := hi<<8 + lo
	addr := base + uint16(c.x)

	oops := 0
	if base&0xFF00 != addr&0xFF00 {
		oops = 1
	}
	return addr, oops
}

func (c *CPU) absoluteY() (uint16, int) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	base := hi<<8 + lo
	addr := base + uint16(c.y)

	oops := 0
	if base&0xFF00 != addr&0xFF00 {
		oops = 1
	}
	return addr, oops
}

func (c *CPU) indirect() uint16 {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	ptr := hi<<8 + lo
	addrLo := uint16(c.read(ptr))
	addrHi := uint16(c.read((ptr & 0xFF00) + ((ptr + 1) & 0x00FF))) // 6502 bug
	return addrHi<<8 + addrLo
}
func (c *CPU) indirectX() uint16 {
	zp := uint16((c.fetch() + c.x) & 0xFF)
	lo := uint16(c.read(zp))
	hi := uint16(c.read((zp + 1) & 0xFF))
	return hi<<8 + lo
}

func (c *CPU) indirectY() (uint16, int) {
	zp := uint16(c.fetch())
	lo := uint16(c.read(zp))
	hi := uint16(c.read((zp + 1) & 0xFF))
	base := hi<<8 + lo
	addr := base + uint16(c.y)

	oops := 0
	if base&0xFF00 != addr&0xFF00 {
		oops = 1
	}
	return addr, oops
}

func (c *CPU) relative() (uint16, int) {
	offset := int8(c.fetch())
	addr := uint16(int32(c.pc) + int32(offset))

	oops := 0
	if c.pc&0xFF00 != addr&0xFF00 {
		oops = 1
	}

	return addr, oops
}
