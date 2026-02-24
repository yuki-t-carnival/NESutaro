package cpu

import (
	"nesutaro/internal/cpu/bus"
)

const (
	// Status Register(p)
	CarryFlagMask            byte = 1 << 0
	ZeroFlagMask             byte = 1 << 1
	InterruptDisableFlagMask byte = 1 << 2
	DecimalFlagMask          byte = 1 << 3
	TheBFlagMask             byte = 1 << 4
	OverflowFlagMask         byte = 1 << 6
	NegativeFlagMask         byte = 1 << 7
)

type CPU struct {
	Tracer *Tracer
	Bus    *bus.Bus

	// Registers
	a, x, y, s, p byte
	pc            uint16

	// Others
	cycles  int
	IsPanic bool

	isIFlagToggleDelayed bool

	testcnt int
}

func NewCPU(b *bus.Bus) *CPU {
	c := &CPU{
		Bus: b,
		s:   0xFD,
		p:   0x24,
	}
	lo := uint16(c.read(0xFFFC))
	hi := uint16(c.read(0xFFFD))
	c.pc = hi<<8 | lo
	return c
}

func (c *CPU) Step() int {
	c.cycles = 0

	//prevPC := c.pc
	if c.Bus.PPU.HasNMI() {
		c.nmi()
	}
	if c.Bus.HasIRQ && c.p&InterruptDisableFlagMask == 0 {
		c.irq()
	}
	if c.isIFlagToggleDelayed {
		c.p |= InterruptDisableFlagMask
		c.isIFlagToggleDelayed = false
	}
	op := c.fetch()

	/* if c.testcnt < 100 {
		opStr := fmt.Sprintf("%02X", op)
		for i := 0; i < opTable[op].Bytes-1; i++ {
			opStr += fmt.Sprintf(" %02X", c.read(c.pc+uint16(i)))
		}
		var name string
		if opTable[op].Name[0] == '*' {
			name = opTable[op].Name
		} else {
			name = " " + opTable[op].Name
		}
		fmt.Printf("%5d: %04X  %-8s %-20s", c.testcnt+1, prevPC, opStr, name)
		fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.a, c.x, c.y, c.p, c.s)
		c.testcnt++
	} */
	opTable[op].fn(c)

	return c.cycles
}

func (c *CPU) read(addr uint16) byte {
	return c.Bus.Read(addr)
}

func (c *CPU) write(addr uint16, val byte) {
	c.Bus.Write(addr, val)
}

func (c *CPU) fetch() byte {
	v := c.read(c.pc)
	c.pc++
	return v
}

func (c *CPU) nmi() {
	hi := byte(c.pc >> 8)
	c.write(0x0100+uint16(c.s), hi)
	c.s -= 1

	lo := byte(c.pc & 0x00FF)
	c.write(0x0100+uint16(c.s), lo)
	c.s -= 1

	p := c.p &^ TheBFlagMask
	c.write(0x100+uint16(c.s), p)
	c.s -= 1

	c.p |= InterruptDisableFlagMask
	nextLo := uint16(c.read(0xFFFA))
	nextHi := uint16(c.read(0xFFFB))
	c.pc = nextHi<<8 | nextLo

	c.Bus.PPU.DisableNMI()
	c.Bus.HasIRQ = false
}

func (c *CPU) irq() {
	hi := byte(c.pc >> 8)
	c.write(0x0100+uint16(c.s), hi)
	c.s -= 1

	lo := byte(c.pc & 0x00FF)
	c.write(0x0100+uint16(c.s), lo)
	c.s -= 1

	p := c.p &^ TheBFlagMask
	c.write(0x100+uint16(c.s), p)
	c.s -= 1

	c.p |= InterruptDisableFlagMask
	nextLo := uint16(c.read(0xFFFE))
	nextHi := uint16(c.read(0xFFFF))
	c.pc = nextHi<<8 | nextLo

	c.Bus.HasIRQ = false
}
