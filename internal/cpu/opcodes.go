package cpu

type OpEntry struct {
	fn    func(c *CPU)
	Name  string
	Bytes int
}

var opTable = [256]OpEntry{
	// ADC - Add with Carry
	0x69: {fn: func(c *CPU) { addr := c.immediate(); c.adc(addr); c.cycles += 2 }, Name: "ADC #Immediate", Bytes: 2},
	0x65: {fn: func(c *CPU) { addr := c.zeroPage(); c.adc(addr); c.cycles += 3 }, Name: "ADC Zero Page", Bytes: 2},
	0x75: {fn: func(c *CPU) { addr := c.zeroPageX(); c.adc(addr); c.cycles += 4 }, Name: "ADC Zero Page, X", Bytes: 2},
	0x6D: {fn: func(c *CPU) { addr := c.absolute(); c.adc(addr); c.cycles += 4 }, Name: "ADC Absolute", Bytes: 3},
	0x7D: {fn: func(c *CPU) { addr, oops := c.absoluteX(); c.adc(addr); c.cycles += 4 + oops }, Name: "ADC Absolute, X", Bytes: 3},
	0x79: {fn: func(c *CPU) { addr, oops := c.absoluteY(); c.adc(addr); c.cycles += 4 + oops }, Name: "ADC Absolute, Y", Bytes: 3},
	0x61: {fn: func(c *CPU) { addr := c.indirectX(); c.adc(addr); c.cycles += 6 }, Name: "ADC (Indirect, X)", Bytes: 2},
	0x71: {fn: func(c *CPU) { addr, oops := c.indirectY(); c.adc(addr); c.cycles += 5 + oops }, Name: "ADC (Indirect), Y", Bytes: 2},
	// AND - Bitwise AND
	0x29: {fn: func(c *CPU) { addr := c.immediate(); c.and(addr); c.cycles += 2 }, Name: "AND #Immediate", Bytes: 2},
	0x25: {fn: func(c *CPU) { addr := c.zeroPage(); c.and(addr); c.cycles += 3 }, Name: "AND Zero Page", Bytes: 2},
	0x35: {fn: func(c *CPU) { addr := c.zeroPageX(); c.and(addr); c.cycles += 4 }, Name: "AND Zero Page, X", Bytes: 2},
	0x2D: {fn: func(c *CPU) { addr := c.absolute(); c.and(addr); c.cycles += 4 }, Name: "AND Absolute", Bytes: 3},
	0x3D: {fn: func(c *CPU) { addr, oops := c.absoluteX(); c.and(addr); c.cycles += 4 + oops }, Name: "AND Absolute, X", Bytes: 3},
	0x39: {fn: func(c *CPU) { addr, oops := c.absoluteY(); c.and(addr); c.cycles += 4 + oops }, Name: "AND Absolute, Y", Bytes: 3},
	0x21: {fn: func(c *CPU) { addr := c.indirectX(); c.and(addr); c.cycles += 6 }, Name: "AND (Indirect, X)", Bytes: 2},
	0x31: {fn: func(c *CPU) { addr, oops := c.indirectY(); c.and(addr); c.cycles += 5 + oops }, Name: "AND (Indirect), Y", Bytes: 2},
	// ASL - Arithmetic Shift Left
	0x0A: {fn: func(c *CPU) { c.aslAccum(); c.cycles += 2 }, Name: "ASL Accumlator", Bytes: 1},
	0x06: {fn: func(c *CPU) { addr := c.zeroPage(); c.asl(addr); c.cycles += 5 }, Name: "ASL Zero Page", Bytes: 2},
	0x16: {fn: func(c *CPU) { addr := c.zeroPageX(); c.asl(addr); c.cycles += 6 }, Name: "ASL Zero Page, X", Bytes: 2},
	0x0E: {fn: func(c *CPU) { addr := c.absolute(); c.asl(addr); c.cycles += 6 }, Name: "ASL Absolute", Bytes: 3},
	0x1E: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.asl(addr); c.cycles += 7 }, Name: "ASL Absolute, X", Bytes: 3},
	// BCC - Branch if Carry Clear
	0x90: {fn: func(c *CPU) {
		addr, oops := c.relative()
		isBranched := c.bcc(addr)
		c.cycles += c.calcBranchInstrCycles(isBranched, oops)
	}, Name: "BCC Relative", Bytes: 2},
	// BCS - Branch if Carry Set
	0xB0: {fn: func(c *CPU) {
		addr, oops := c.relative()
		isBranched := c.bcs(addr)
		c.cycles += c.calcBranchInstrCycles(isBranched, oops)
	}, Name: "BCS Relative", Bytes: 2},
	// BEQ - Branch if Equal
	0xF0: {fn: func(c *CPU) {
		addr, oops := c.relative()
		isBranched := c.beq(addr)
		c.cycles += c.calcBranchInstrCycles(isBranched, oops)
	}, Name: "BEQ Relative", Bytes: 2},
	// BIT - Bit Test
	0x24: {fn: func(c *CPU) { addr := c.zeroPage(); c.bit(addr); c.cycles += 3 }, Name: "BIT Zero Page", Bytes: 2},
	0x2C: {fn: func(c *CPU) { addr := c.absolute(); c.bit(addr); c.cycles += 4 }, Name: "Bit Absolute", Bytes: 3},
	// BMI - Branch if Minus
	0x30: {fn: func(c *CPU) {
		addr, oops := c.relative()
		isBranched := c.bmi(addr)
		c.cycles += c.calcBranchInstrCycles(isBranched, oops)
	}, Name: "BMI Relative", Bytes: 2},
	// BNE - Branch if Not Equal
	0xD0: {fn: func(c *CPU) {
		addr, oops := c.relative()
		isBranched := c.bne(addr)
		c.cycles += c.calcBranchInstrCycles(isBranched, oops)
	}, Name: "BNE Relative", Bytes: 2},
	// BPL - Branch if Plus
	0x10: {fn: func(c *CPU) {
		addr, oops := c.relative()
		isBranched := c.bpl(addr)
		c.cycles += c.calcBranchInstrCycles(isBranched, oops)
	}, Name: "BPL Relative", Bytes: 2},
	// BRK - Break (software IRQ)
	0x00: {fn: func(c *CPU) { c.brk(); c.cycles += 7 }, Name: "BRK #Immediate", Bytes: 2},
	// BVC - Branch if Overflow Clear
	0x50: {fn: func(c *CPU) {
		addr, oops := c.relative()
		isBranched := c.bvc(addr)
		c.cycles += c.calcBranchInstrCycles(isBranched, oops)
	}, Name: "BVC Relative", Bytes: 2},
	// BVS - Branch if Overflow Set
	0x70: {fn: func(c *CPU) {
		addr, oops := c.relative()
		isBranched := c.bvs(addr)
		c.cycles += c.calcBranchInstrCycles(isBranched, oops)
	}, Name: "BVS Relative", Bytes: 2},
	// CLC - Clear Carry
	0x18: {fn: func(c *CPU) { c.clc(); c.cycles += 2 }, Name: "CLC Implied", Bytes: 1},
	// CLD - Clear Decimal
	0xD8: {fn: func(c *CPU) { c.cld(); c.cycles += 2 }, Name: "CLD Implied", Bytes: 1},
	// CLI - Clear Interrupt Disable
	0x58: {fn: func(c *CPU) { c.cli(); c.cycles += 2 }, Name: "CLI Implied", Bytes: 1},
	// CLV - Clear Overflow
	0xB8: {fn: func(c *CPU) { c.clv(); c.cycles += 2 }, Name: "CLV Implied", Bytes: 1},
	// CMP - Compare A
	0xC9: {fn: func(c *CPU) { addr := c.immediate(); c.cmp(addr); c.cycles += 2 }, Name: "CMP #Immediate", Bytes: 2},
	0xC5: {fn: func(c *CPU) { addr := c.zeroPage(); c.cmp(addr); c.cycles += 3 }, Name: "CMP Zero Page", Bytes: 2},
	0xD5: {fn: func(c *CPU) { addr := c.zeroPageX(); c.cmp(addr); c.cycles += 4 }, Name: "CMP Zero Page, X", Bytes: 2},
	0xCD: {fn: func(c *CPU) { addr := c.absolute(); c.cmp(addr); c.cycles += 4 }, Name: "CMP Absolute", Bytes: 3},
	0xDD: {fn: func(c *CPU) { addr, oops := c.absoluteX(); c.cmp(addr); c.cycles += 4 + oops }, Name: "CMP Absolute, X", Bytes: 3},
	0xD9: {fn: func(c *CPU) { addr, oops := c.absoluteY(); c.cmp(addr); c.cycles += 4 + oops }, Name: "CMP Absolute, Y", Bytes: 3},
	0xC1: {fn: func(c *CPU) { addr := c.indirectX(); c.cmp(addr); c.cycles += 6 }, Name: "CMP (Indirect, X)", Bytes: 2},
	0xD1: {fn: func(c *CPU) { addr, oops := c.indirectY(); c.cmp(addr); c.cycles += 5 + oops }, Name: "CMP (Indirect), Y", Bytes: 2},
	// CPX - Compare X
	0xE0: {fn: func(c *CPU) { addr := c.immediate(); c.cpx(addr); c.cycles += 2 }, Name: "CPX #Immediate", Bytes: 2},
	0xE4: {fn: func(c *CPU) { addr := c.zeroPage(); c.cpx(addr); c.cycles += 3 }, Name: "CPX Zero Page", Bytes: 2},
	0xEC: {fn: func(c *CPU) { addr := c.absolute(); c.cpx(addr); c.cycles += 4 }, Name: "CPX Absolute", Bytes: 3},
	// CPY - Compare Y
	0xC0: {fn: func(c *CPU) { addr := c.immediate(); c.cpy(addr); c.cycles += 2 }, Name: "CPY #Immediate", Bytes: 2},
	0xC4: {fn: func(c *CPU) { addr := c.zeroPage(); c.cpy(addr); c.cycles += 3 }, Name: "CPY Zero Page", Bytes: 2},
	0xCC: {fn: func(c *CPU) { addr := c.absolute(); c.cpy(addr); c.cycles += 4 }, Name: "CPY Absolute", Bytes: 3},
	// DEC - Decrement Memory
	0xC6: {fn: func(c *CPU) { addr := c.zeroPage(); c.dec(addr); c.cycles += 5 }, Name: "DEC Zero Page", Bytes: 2},
	0xD6: {fn: func(c *CPU) { addr := c.zeroPageX(); c.dec(addr); c.cycles += 6 }, Name: "DEC Zero Page, X", Bytes: 2},
	0xCE: {fn: func(c *CPU) { addr := c.absolute(); c.dec(addr); c.cycles += 6 }, Name: "DEC Absolute", Bytes: 3},
	0xDE: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.dec(addr); c.cycles += 7 }, Name: "DEC Absolute, X", Bytes: 3},
	// DEX - Decrement X
	0xCA: {fn: func(c *CPU) { c.dex(); c.cycles += 2 }, Name: "DEX Implied", Bytes: 1},
	// DEY - Decrement Y
	0x88: {fn: func(c *CPU) { c.dey(); c.cycles += 2 }, Name: "DEY Implied", Bytes: 1},
	// EOR - Bitwise Exclusive OR
	0x49: {fn: func(c *CPU) { addr := c.immediate(); c.eor(addr); c.cycles += 2 }, Name: "EOR #Immediate", Bytes: 2},
	0x45: {fn: func(c *CPU) { addr := c.zeroPage(); c.eor(addr); c.cycles += 3 }, Name: "EOR Zero Page", Bytes: 2},
	0x55: {fn: func(c *CPU) { addr := c.zeroPageX(); c.eor(addr); c.cycles += 4 }, Name: "EOR Zero Page, X", Bytes: 2},
	0x4D: {fn: func(c *CPU) { addr := c.absolute(); c.eor(addr); c.cycles += 4 }, Name: "EOR Absolute", Bytes: 3},
	0x5D: {fn: func(c *CPU) { addr, oops := c.absoluteX(); c.eor(addr); c.cycles += 4 + oops }, Name: "EOR Absolute, X", Bytes: 3},
	0x59: {fn: func(c *CPU) { addr, oops := c.absoluteY(); c.eor(addr); c.cycles += 4 + oops }, Name: "EOR Absolute, Y", Bytes: 3},
	0x41: {fn: func(c *CPU) { addr := c.indirectX(); c.eor(addr); c.cycles += 6 }, Name: "EOR (Indirect, X)", Bytes: 2},
	0x51: {fn: func(c *CPU) { addr, oops := c.indirectY(); c.eor(addr); c.cycles += 5 + oops }, Name: "EOR (Indirect), Y", Bytes: 2},
	// INC - Increment Memory
	0xE6: {fn: func(c *CPU) { addr := c.zeroPage(); c.inc(addr); c.cycles += 5 }, Name: "INC Zero Page", Bytes: 2},
	0xF6: {fn: func(c *CPU) { addr := c.zeroPageX(); c.inc(addr); c.cycles += 6 }, Name: "INC Zero Page, X", Bytes: 2},
	0xEE: {fn: func(c *CPU) { addr := c.absolute(); c.inc(addr); c.cycles += 6 }, Name: "INC Absolute", Bytes: 3},
	0xFE: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.inc(addr); c.cycles += 7 }, Name: "INC Absolute, X", Bytes: 3},
	// INX - Increment X
	0xE8: {fn: func(c *CPU) { c.inx(); c.cycles += 2 }, Name: "INX Implied", Bytes: 1},
	// INY - Increment Y
	0xC8: {fn: func(c *CPU) { c.iny(); c.cycles += 2 }, Name: "INY Implied", Bytes: 1},
	// JMP - Jump
	0x4C: {fn: func(c *CPU) { addr := c.absolute(); c.jmp(addr); c.cycles += 3 }, Name: "JMP Absolute", Bytes: 3},
	0x6C: {fn: func(c *CPU) { addr := c.indirect(); c.jmp(addr); c.cycles += 5 }, Name: "JMP (Indirect)", Bytes: 3},
	// JSR - Jump to Subroutine
	0x20: {fn: func(c *CPU) { addr := c.absolute(); c.jsr(addr); c.cycles += 6 }, Name: "JSR Absolute", Bytes: 3},
	// LDA - Load A
	0xA9: {fn: func(c *CPU) { addr := c.immediate(); c.lda(addr); c.cycles += 2 }, Name: "LDA #Immediate", Bytes: 2},
	0xA5: {fn: func(c *CPU) { addr := c.zeroPage(); c.lda(addr); c.cycles += 3 }, Name: "LDA Zero Page", Bytes: 2},
	0xB5: {fn: func(c *CPU) { addr := c.zeroPageX(); c.lda(addr); c.cycles += 4 }, Name: "LDA Zero Page, X", Bytes: 2},
	0xAD: {fn: func(c *CPU) { addr := c.absolute(); c.lda(addr); c.cycles += 4 }, Name: "LDA Absolute", Bytes: 3},
	0xBD: {fn: func(c *CPU) { addr, oops := c.absoluteX(); c.lda(addr); c.cycles += 4 + oops }, Name: "LDA Absolute, X", Bytes: 3},
	0xB9: {fn: func(c *CPU) { addr, oops := c.absoluteY(); c.lda(addr); c.cycles += 4 + oops }, Name: "LDA Absolute, Y", Bytes: 3},
	0xA1: {fn: func(c *CPU) { addr := c.indirectX(); c.lda(addr); c.cycles += 6 }, Name: "LDA (Indirect, X)", Bytes: 2},
	0xB1: {fn: func(c *CPU) { addr, oops := c.indirectY(); c.lda(addr); c.cycles += 5 + oops }, Name: "LDA (Indirect), Y", Bytes: 2},
	// LDX - Load X
	0xA2: {fn: func(c *CPU) { addr := c.immediate(); c.ldx(addr); c.cycles += 2 }, Name: "LDX #Immediate", Bytes: 2},
	0xA6: {fn: func(c *CPU) { addr := c.zeroPage(); c.ldx(addr); c.cycles += 3 }, Name: "LDX Zero Page", Bytes: 2},
	0xB6: {fn: func(c *CPU) { addr := c.zeroPageY(); c.ldx(addr); c.cycles += 4 }, Name: "LDX Zero Page, Y", Bytes: 2},
	0xAE: {fn: func(c *CPU) { addr := c.absolute(); c.ldx(addr); c.cycles += 4 }, Name: "LDX Absolute", Bytes: 3},
	0xBE: {fn: func(c *CPU) { addr, oops := c.absoluteY(); c.ldx(addr); c.cycles += 4 + oops }, Name: "LDX Absolute, Y", Bytes: 3},
	// LDY - Load Y
	0xA0: {fn: func(c *CPU) { addr := c.immediate(); c.ldy(addr); c.cycles += 2 }, Name: "LDY #Immediate", Bytes: 2},
	0xA4: {fn: func(c *CPU) { addr := c.zeroPage(); c.ldy(addr); c.cycles += 3 }, Name: "LDY Zero Page", Bytes: 2},
	0xB4: {fn: func(c *CPU) { addr := c.zeroPageX(); c.ldy(addr); c.cycles += 4 }, Name: "LDY Zero Page, X", Bytes: 2},
	0xAC: {fn: func(c *CPU) { addr := c.absolute(); c.ldy(addr); c.cycles += 4 }, Name: "LDY Absolute", Bytes: 3},
	0xBC: {fn: func(c *CPU) { addr, oops := c.absoluteX(); c.ldy(addr); c.cycles += 4 + oops }, Name: "LDY Absolute, X", Bytes: 3},
	// LSR - Logical Shift Right
	0x4A: {fn: func(c *CPU) { c.lsrAccum(); c.cycles += 2 }, Name: "LSR Accumlator", Bytes: 1},
	0x46: {fn: func(c *CPU) { addr := c.zeroPage(); c.lsr(addr); c.cycles += 5 }, Name: "LSR Zero Page", Bytes: 2},
	0x56: {fn: func(c *CPU) { addr := c.zeroPageX(); c.lsr(addr); c.cycles += 6 }, Name: "LSR Zero Page, X", Bytes: 2},
	0x4E: {fn: func(c *CPU) { addr := c.absolute(); c.lsr(addr); c.cycles += 6 }, Name: "LSR Absolute", Bytes: 3},
	0x5E: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.lsr(addr); c.cycles += 7 }, Name: "LSR Absolute, X", Bytes: 3},
	// NOP - No Opertion
	0xEA: {fn: func(c *CPU) { c.nop(); c.cycles += 2 }, Name: "NOP Implied", Bytes: 1},
	// ORA - Bitwise OR
	0x09: {fn: func(c *CPU) { addr := c.immediate(); c.ora(addr); c.cycles += 2 }, Name: "ORA #Immediate", Bytes: 2},
	0x05: {fn: func(c *CPU) { addr := c.zeroPage(); c.ora(addr); c.cycles += 3 }, Name: "ORA Zero Page", Bytes: 2},
	0x15: {fn: func(c *CPU) { addr := c.zeroPageX(); c.ora(addr); c.cycles += 4 }, Name: "ORA Zero Page, X", Bytes: 2},
	0x0D: {fn: func(c *CPU) { addr := c.absolute(); c.ora(addr); c.cycles += 4 }, Name: "ORA Absolute", Bytes: 3},
	0x1D: {fn: func(c *CPU) { addr, oops := c.absoluteX(); c.ora(addr); c.cycles += 4 + oops }, Name: "ORA Absolute, X", Bytes: 3},
	0x19: {fn: func(c *CPU) { addr, oops := c.absoluteY(); c.ora(addr); c.cycles += 4 + oops }, Name: "ORA Absolute, Y", Bytes: 3},
	0x01: {fn: func(c *CPU) { addr := c.indirectX(); c.ora(addr); c.cycles += 6 }, Name: "ORA (Indirect, X)", Bytes: 2},
	0x11: {fn: func(c *CPU) { addr, oops := c.indirectY(); c.ora(addr); c.cycles += 5 + oops }, Name: "ORA (Indirect), Y", Bytes: 2},
	// PHA - Push A
	0x48: {fn: func(c *CPU) { c.pha(); c.cycles += 3 }, Name: "PHA Implied", Bytes: 1},
	// PHP - Push Processor Status
	0x08: {fn: func(c *CPU) { c.php(); c.cycles += 3 }, Name: "PHP Implied", Bytes: 1},
	// PLA - Pull A
	0x68: {fn: func(c *CPU) { c.pla(); c.cycles += 4 }, Name: "PLA Implied", Bytes: 1},
	// PLP - Pull Processor Status
	0x28: {fn: func(c *CPU) { c.plp(); c.cycles += 4 }, Name: "PLP Implied", Bytes: 1},
	// ROL - Rotate Left
	0x2A: {fn: func(c *CPU) { c.rolAccum(); c.cycles += 2 }, Name: "ROL Accumlator", Bytes: 1},
	0x26: {fn: func(c *CPU) { addr := c.zeroPage(); c.rol(addr); c.cycles += 5 }, Name: "ROL Zero Page", Bytes: 2},
	0x36: {fn: func(c *CPU) { addr := c.zeroPageX(); c.rol(addr); c.cycles += 6 }, Name: "ROL Zero Page, X", Bytes: 2},
	0x2E: {fn: func(c *CPU) { addr := c.absolute(); c.rol(addr); c.cycles += 6 }, Name: "ROL Absolute", Bytes: 3},
	0x3E: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.rol(addr); c.cycles += 7 }, Name: "ROL Absolute, X", Bytes: 3},
	// ROR - Rotate Right
	0x6A: {fn: func(c *CPU) { c.rorAccum(); c.cycles += 2 }, Name: "ROR Accumlator", Bytes: 1},
	0x66: {fn: func(c *CPU) { addr := c.zeroPage(); c.ror(addr); c.cycles += 5 }, Name: "ROR Zero Page", Bytes: 2},
	0x76: {fn: func(c *CPU) { addr := c.zeroPageX(); c.ror(addr); c.cycles += 6 }, Name: "ROR Zero Page, X", Bytes: 2},
	0x6E: {fn: func(c *CPU) { addr := c.absolute(); c.ror(addr); c.cycles += 6 }, Name: "ROR Absolute", Bytes: 3},
	0x7E: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.ror(addr); c.cycles += 7 }, Name: "ROR Absolute, X", Bytes: 3},
	// RTI - Return from Interrupt
	0x40: {fn: func(c *CPU) { c.rti(); c.cycles += 6 }, Name: "RTI Implied", Bytes: 1},
	// RTS - Return from Subroutine
	0x60: {fn: func(c *CPU) { c.rts(); c.cycles += 6 }, Name: "RTS Implied", Bytes: 1},
	// SBC - Subtract with Carry
	0xE9: {fn: func(c *CPU) { addr := c.immediate(); c.sbc(addr); c.cycles += 2 }, Name: "SBC #Immediate", Bytes: 2},
	0xE5: {fn: func(c *CPU) { addr := c.zeroPage(); c.sbc(addr); c.cycles += 3 }, Name: "SBC Zero Page", Bytes: 2},
	0xF5: {fn: func(c *CPU) { addr := c.zeroPageX(); c.sbc(addr); c.cycles += 4 }, Name: "SBC Zero Page, X", Bytes: 2},
	0xED: {fn: func(c *CPU) { addr := c.absolute(); c.sbc(addr); c.cycles += 4 }, Name: "SBC Absolute", Bytes: 3},
	0xFD: {fn: func(c *CPU) { addr, oops := c.absoluteX(); c.sbc(addr); c.cycles += 4 + oops }, Name: "SBC Absolute, X", Bytes: 3},
	0xF9: {fn: func(c *CPU) { addr, oops := c.absoluteY(); c.sbc(addr); c.cycles += 4 + oops }, Name: "SBC Absolute, Y", Bytes: 3},
	0xE1: {fn: func(c *CPU) { addr := c.indirectX(); c.sbc(addr); c.cycles += 6 }, Name: "SBC (Indirect, X)", Bytes: 2},
	0xF1: {fn: func(c *CPU) { addr, oops := c.indirectY(); c.sbc(addr); c.cycles += 5 + oops }, Name: "SBC (Indirect), Y", Bytes: 2},
	// SEC - Set Carry
	0x38: {fn: func(c *CPU) { c.sec(); c.cycles += 2 }, Name: "SEC Implied", Bytes: 1},
	// SED - Set Decimal
	0xF8: {fn: func(c *CPU) { c.sed(); c.cycles += 2 }, Name: "SED Implied", Bytes: 1},
	// SEI - Set Interrupt Disable
	0x78: {fn: func(c *CPU) { c.sei(); c.cycles += 2 }, Name: "SEI Implied", Bytes: 1},
	// STA - Store A
	0x85: {fn: func(c *CPU) { addr := c.zeroPage(); c.sta(addr); c.cycles += 3 }, Name: "STA Zero Page", Bytes: 2},
	0x95: {fn: func(c *CPU) { addr := c.zeroPageX(); c.sta(addr); c.cycles += 4 }, Name: "STA Zero Page, X", Bytes: 2},
	0x8D: {fn: func(c *CPU) { addr := c.absolute(); c.sta(addr); c.cycles += 4 }, Name: "STA Absolute", Bytes: 3},
	0x9D: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.sta(addr); c.cycles += 5 }, Name: "STA Absolute, X", Bytes: 3},
	0x99: {fn: func(c *CPU) { addr, _ := c.absoluteY(); c.sta(addr); c.cycles += 5 }, Name: "STA Absolute, Y", Bytes: 3},
	0x81: {fn: func(c *CPU) { addr := c.indirectX(); c.sta(addr); c.cycles += 6 }, Name: "STA (Indirect, X)", Bytes: 2},
	0x91: {fn: func(c *CPU) { addr, _ := c.indirectY(); c.sta(addr); c.cycles += 6 }, Name: "STA (Indirect), Y", Bytes: 2},
	// STX - Store X
	0x86: {fn: func(c *CPU) { addr := c.zeroPage(); c.stx(addr); c.cycles += 3 }, Name: "STX Zero Page", Bytes: 2},
	0x96: {fn: func(c *CPU) { addr := c.zeroPageY(); c.stx(addr); c.cycles += 4 }, Name: "STX Zero Page, Y", Bytes: 2},
	0x8E: {fn: func(c *CPU) { addr := c.absolute(); c.stx(addr); c.cycles += 4 }, Name: "STX Absolute", Bytes: 3},
	// STY - Store Y
	0x84: {fn: func(c *CPU) { addr := c.zeroPage(); c.sty(addr); c.cycles += 3 }, Name: "STY Zero Page", Bytes: 2},
	0x94: {fn: func(c *CPU) { addr := c.zeroPageX(); c.sty(addr); c.cycles += 4 }, Name: "STY Zero Page, X", Bytes: 2},
	0x8C: {fn: func(c *CPU) { addr := c.absolute(); c.sty(addr); c.cycles += 4 }, Name: "STY Absolute", Bytes: 3},
	// TAX - Transfer A to X
	0xAA: {fn: func(c *CPU) { c.tax(); c.cycles += 2 }, Name: "TAX Implied", Bytes: 1},
	// TAY - Transfer A to Y
	0xA8: {fn: func(c *CPU) { c.tay(); c.cycles += 2 }, Name: "TAY Implied", Bytes: 1},
	// TSX - Transfer Stack Pointer to X
	0xBA: {fn: func(c *CPU) { c.tsx(); c.cycles += 2 }, Name: "TSX Implied", Bytes: 1},
	// TXA - Transfer X to A
	0x8A: {fn: func(c *CPU) { c.txa(); c.cycles += 2 }, Name: "TXA Implied", Bytes: 1},
	// TXS - Transfer X to Stack Pointer
	0x9A: {fn: func(c *CPU) { c.txs(); c.cycles += 2 }, Name: "TXS Implied", Bytes: 1},
	// TYA - Transfer Y to A
	0x98: {fn: func(c *CPU) { c.tya(); c.cycles += 2 }, Name: "TYA Implied", Bytes: 1},

	// Unofficial
	// *AHX
	0x93: {fn: func(c *CPU) { addr, _ := c.indirectY(); c.ahx(addr); c.cycles += 6 }, Name: "*AHX (Indirect), Y", Bytes: 2},
	0x9F: {fn: func(c *CPU) { addr, _ := c.absoluteY(); c.ahx(addr); c.cycles += 5 }, Name: "*AHX Absolute, Y", Bytes: 3},
	// *ALR
	0x4B: {fn: func(c *CPU) { addr := c.immediate(); c.alr(addr); c.cycles += 2 }, Name: "*ALR #Immediate", Bytes: 2},
	// *ANC
	0x0B: {fn: func(c *CPU) { addr := c.immediate(); c.anc(addr); c.cycles += 2 }, Name: "*ANC #Immediate", Bytes: 2},
	0x2B: {fn: func(c *CPU) { addr := c.immediate(); c.anc(addr); c.cycles += 2 }, Name: "*ANC #Immediate", Bytes: 2},
	// *ARR
	0x6B: {fn: func(c *CPU) { addr := c.immediate(); c.arr(addr); c.cycles += 2 }, Name: "*ARR #Immediate", Bytes: 2},
	// *AXS
	0xCB: {fn: func(c *CPU) { addr := c.immediate(); c.axs(addr); c.cycles += 2 }, Name: "*AXS #Immediate", Bytes: 2},
	// *DCP
	0xC3: {fn: func(c *CPU) { addr := c.indirectX(); c.dcp(addr); c.cycles += 8 }, Name: "*DCP (Indirect, X)", Bytes: 2},
	0xC7: {fn: func(c *CPU) { addr := c.zeroPage(); c.dcp(addr); c.cycles += 5 }, Name: "*DCP Zero Page", Bytes: 2},
	0xCF: {fn: func(c *CPU) { addr := c.absolute(); c.dcp(addr); c.cycles += 6 }, Name: "*DCP Absolute", Bytes: 3},
	0xD3: {fn: func(c *CPU) { addr, _ := c.indirectY(); c.dcp(addr); c.cycles += 8 }, Name: "*DCP (Indirect), Y", Bytes: 2},
	0xD7: {fn: func(c *CPU) { addr := c.zeroPageX(); c.dcp(addr); c.cycles += 6 }, Name: "*DCP Zero Page, X", Bytes: 2},
	0xDB: {fn: func(c *CPU) { addr, _ := c.absoluteY(); c.dcp(addr); c.cycles += 7 }, Name: "*DCP Absolute, Y", Bytes: 3},
	0xDF: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.dcp(addr); c.cycles += 7 }, Name: "*DCP Absolute, X", Bytes: 3},
	// *ISC
	0xE3: {fn: func(c *CPU) { addr := c.indirectX(); c.isc(addr); c.cycles += 8 }, Name: "*ISC (Indirect, X)", Bytes: 2},
	0xE7: {fn: func(c *CPU) { addr := c.zeroPage(); c.isc(addr); c.cycles += 5 }, Name: "*ISC Zero Page", Bytes: 2},
	0xEF: {fn: func(c *CPU) { addr := c.absolute(); c.isc(addr); c.cycles += 6 }, Name: "*ISC Absolute", Bytes: 3},
	0xF3: {fn: func(c *CPU) { addr, _ := c.indirectY(); c.isc(addr); c.cycles += 8 }, Name: "*ISC (Indirect), Y", Bytes: 2},
	0xF7: {fn: func(c *CPU) { addr := c.zeroPageX(); c.isc(addr); c.cycles += 6 }, Name: "*ISC Zero Page, X", Bytes: 2},
	0xFB: {fn: func(c *CPU) { addr, _ := c.absoluteY(); c.isc(addr); c.cycles += 7 }, Name: "*ISC Absolute, Y", Bytes: 3},
	0xFF: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.isc(addr); c.cycles += 7 }, Name: "*ISC Absolute, X", Bytes: 3},
	// *KIL
	0x02: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	0x12: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	0x22: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	0x32: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	0x42: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	0x52: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	0x62: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	0x72: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	0x92: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	0xB2: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	0xD2: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	0xF2: {fn: func(c *CPU) { panic("*KIL\n") }, Name: "*KIL", Bytes: 1},
	// *LAS
	0xBB: {fn: func(c *CPU) { addr, oops := c.absoluteY(); c.las(addr); c.cycles += 4 + oops }, Name: "*LAS Absolute, Y", Bytes: 3},
	// *LAX
	0xA3: {fn: func(c *CPU) { addr := c.indirectX(); c.laxAddr(addr); c.cycles += 6 }, Name: "*LAX (Indirect, X)", Bytes: 2},
	0xA7: {fn: func(c *CPU) { addr := c.zeroPage(); c.laxAddr(addr); c.cycles += 3 }, Name: "*LAX Zero Page", Bytes: 2},
	0xAB: {fn: func(c *CPU) { addr := c.immediate(); c.laxImm(addr); c.cycles += 2 }, Name: "*LAX #Immediate", Bytes: 2},
	0xAF: {fn: func(c *CPU) { addr := c.absolute(); c.laxAddr(addr); c.cycles += 4 }, Name: "*LAX Absolute", Bytes: 3},
	0xB3: {fn: func(c *CPU) { addr, oops := c.indirectY(); c.laxAddr(addr); c.cycles += 5 + oops }, Name: "*LAX (Indirect), Y", Bytes: 2},
	0xB7: {fn: func(c *CPU) { addr := c.zeroPageY(); c.laxAddr(addr); c.cycles += 4 }, Name: "*LAX Zero Page, Y", Bytes: 2},
	0xBF: {fn: func(c *CPU) { addr, oops := c.absoluteY(); c.laxAddr(addr); c.cycles += 4 + oops }, Name: "*LAX Absolute, Y", Bytes: 3},
	// *NOP
	0x04: {fn: func(c *CPU) { _ = c.zeroPage(); c.nop(); c.cycles += 3 }, Name: "*NOP Zero Page", Bytes: 2},
	0x0C: {fn: func(c *CPU) { _ = c.absolute(); c.nop(); c.cycles += 4 }, Name: "*NOP Absolute", Bytes: 3},
	0x14: {fn: func(c *CPU) { _ = c.zeroPageX(); c.nop(); c.cycles += 4 }, Name: "*NOP Zero Page, X", Bytes: 2},
	0x1A: {fn: func(c *CPU) { c.nop(); c.cycles += 2 }, Name: "*NOP Implied", Bytes: 1},
	0x1C: {fn: func(c *CPU) { _, oops := c.absoluteX(); c.nop(); c.cycles += 4 + oops }, Name: "*NOP Absolute, X", Bytes: 3},
	0x34: {fn: func(c *CPU) { _ = c.zeroPageX(); c.nop(); c.cycles += 4 }, Name: "*NOP Zero Page, X", Bytes: 2},
	0x3A: {fn: func(c *CPU) { c.nop(); c.cycles += 2 }, Name: "*NOP Implied", Bytes: 1},
	0x3C: {fn: func(c *CPU) { _, oops := c.absoluteX(); c.nop(); c.cycles += 4 + oops }, Name: "*NOP Absolute, X", Bytes: 3},
	0x44: {fn: func(c *CPU) { _ = c.zeroPage(); c.nop(); c.cycles += 3 }, Name: "*NOP Zero Page", Bytes: 2},
	0x54: {fn: func(c *CPU) { _ = c.zeroPageX(); c.nop(); c.cycles += 4 }, Name: "*NOP Zero Page, X", Bytes: 2},
	0x5A: {fn: func(c *CPU) { c.nop(); c.cycles += 2 }, Name: "*NOP Implied", Bytes: 1},
	0x5C: {fn: func(c *CPU) { _, oops := c.absoluteX(); c.nop(); c.cycles += 4 + oops }, Name: "*NOP Absolute, X", Bytes: 3},
	0x64: {fn: func(c *CPU) { _ = c.zeroPage(); c.nop(); c.cycles += 3 }, Name: "*NOP Zero Page", Bytes: 2},
	0x74: {fn: func(c *CPU) { _ = c.zeroPageX(); c.nop(); c.cycles += 4 }, Name: "*NOP Zero Page, X", Bytes: 2},
	0x7A: {fn: func(c *CPU) { c.nop(); c.cycles += 2 }, Name: "*NOP Implied", Bytes: 1},
	0x7C: {fn: func(c *CPU) { _, oops := c.absoluteX(); c.nop(); c.cycles += 4 + oops }, Name: "*NOP Absolute, X", Bytes: 3},
	0x80: {fn: func(c *CPU) { _ = c.immediate(); c.nop(); c.cycles += 2 }, Name: "*NOP #Immediate", Bytes: 2},
	0x82: {fn: func(c *CPU) { _ = c.immediate(); c.nop(); c.cycles += 2 }, Name: "*NOP #Immediate", Bytes: 2},
	0x89: {fn: func(c *CPU) { _ = c.immediate(); c.nop(); c.cycles += 2 }, Name: "*NOP #Immediate", Bytes: 2},
	0xC2: {fn: func(c *CPU) { _ = c.immediate(); c.nop(); c.cycles += 2 }, Name: "*NOP #Immediate", Bytes: 2},
	0xD4: {fn: func(c *CPU) { _ = c.zeroPageX(); c.nop(); c.cycles += 4 }, Name: "*NOP Zero Page, X", Bytes: 2},
	0xDA: {fn: func(c *CPU) { c.nop(); c.cycles += 2 }, Name: "*NOP Implied", Bytes: 1},
	0xDC: {fn: func(c *CPU) { _, oops := c.absoluteX(); c.nop(); c.cycles += 4 + oops }, Name: "*NOP Absolute, X", Bytes: 3},
	0xE2: {fn: func(c *CPU) { _ = c.immediate(); c.nop(); c.cycles += 2 }, Name: "*NOP #Immediate", Bytes: 2},
	0xF4: {fn: func(c *CPU) { _ = c.zeroPageX(); c.nop(); c.cycles += 4 }, Name: "*NOP Zero Page, X", Bytes: 2},
	0xFA: {fn: func(c *CPU) { c.nop(); c.cycles += 2 }, Name: "*NOP Implied", Bytes: 1},
	0xFC: {fn: func(c *CPU) { _, oops := c.absoluteX(); c.nop(); c.cycles += 4 + oops }, Name: "*NOP Absolute, X", Bytes: 3},
	// *RLA
	0x23: {fn: func(c *CPU) { addr := c.indirectX(); c.rla(addr); c.cycles += 8 }, Name: "*RLA (Indirect, X)", Bytes: 2},
	0x27: {fn: func(c *CPU) { addr := c.zeroPage(); c.rla(addr); c.cycles += 5 }, Name: "*RLA Zero Page", Bytes: 2},
	0x2F: {fn: func(c *CPU) { addr := c.absolute(); c.rla(addr); c.cycles += 6 }, Name: "*RLA Absolute", Bytes: 3},
	0x33: {fn: func(c *CPU) { addr, _ := c.indirectY(); c.rla(addr); c.cycles += 8 }, Name: "*RLA (Indirect), Y", Bytes: 2},
	0x37: {fn: func(c *CPU) { addr := c.zeroPageX(); c.rla(addr); c.cycles += 6 }, Name: "*RLA Zero Page, X", Bytes: 2},
	0x3B: {fn: func(c *CPU) { addr, _ := c.absoluteY(); c.rla(addr); c.cycles += 7 }, Name: "*RLA Absolute, Y", Bytes: 3},
	0x3F: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.rla(addr); c.cycles += 7 }, Name: "*RLA Absolute, X", Bytes: 3},
	// *RRA
	0x63: {fn: func(c *CPU) { addr := c.indirectX(); c.rra(addr); c.cycles += 8 }, Name: "*RRA (Indirect, X)", Bytes: 2},
	0x67: {fn: func(c *CPU) { addr := c.zeroPage(); c.rra(addr); c.cycles += 5 }, Name: "*RRA Zero Page", Bytes: 2},
	0x6F: {fn: func(c *CPU) { addr := c.absolute(); c.rra(addr); c.cycles += 6 }, Name: "*RRA Absolute", Bytes: 3},
	0x73: {fn: func(c *CPU) { addr, _ := c.indirectY(); c.rra(addr); c.cycles += 8 }, Name: "*RRA (Indirect), Y", Bytes: 2},
	0x77: {fn: func(c *CPU) { addr := c.zeroPageX(); c.rra(addr); c.cycles += 6 }, Name: "*RRA Zero Page, X", Bytes: 2},
	0x7B: {fn: func(c *CPU) { addr, _ := c.absoluteY(); c.rra(addr); c.cycles += 7 }, Name: "*RRA Absolute, Y", Bytes: 3},
	0x7F: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.rra(addr); c.cycles += 7 }, Name: "*RRA Absolute, X", Bytes: 3},
	// *SAX
	0x83: {fn: func(c *CPU) { addr := c.indirectX(); c.sax(addr); c.cycles += 6 }, Name: "*SAX (Indirect, X)", Bytes: 2},
	0x87: {fn: func(c *CPU) { addr := c.zeroPage(); c.sax(addr); c.cycles += 3 }, Name: "*SAX Zero Page", Bytes: 2},
	0x8F: {fn: func(c *CPU) { addr := c.absolute(); c.sax(addr); c.cycles += 4 }, Name: "*SAX Absolute", Bytes: 3},
	0x97: {fn: func(c *CPU) { addr := c.zeroPageY(); c.sax(addr); c.cycles += 4 }, Name: "*SAX Zero Page, Y", Bytes: 2},
	// *SBC
	0xEB: {fn: func(c *CPU) { addr := c.immediate(); c.sbc(addr); c.cycles += 2 }, Name: "*SBC #Immediate", Bytes: 2},
	// *SHX
	0x9E: {fn: func(c *CPU) { addr, _ := c.absoluteY(); c.shx(addr); c.cycles += 5 }, Name: "*SHX Absolute, Y", Bytes: 3},
	// *SHY
	0x9C: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.shy(addr); c.cycles += 5 }, Name: "*SHY Absolute, X", Bytes: 3},
	// *SLO
	0x03: {fn: func(c *CPU) { addr := c.indirectX(); c.slo(addr); c.cycles += 8 }, Name: "*SLO (Indirect, X)", Bytes: 2},
	0x07: {fn: func(c *CPU) { addr := c.zeroPage(); c.slo(addr); c.cycles += 5 }, Name: "*SLO Zero Page", Bytes: 2},
	0x0F: {fn: func(c *CPU) { addr := c.absolute(); c.slo(addr); c.cycles += 6 }, Name: "*SLO Absolute", Bytes: 3},
	0x13: {fn: func(c *CPU) { addr, _ := c.indirectY(); c.slo(addr); c.cycles += 8 }, Name: "*SLO (Indirect), Y", Bytes: 2},
	0x17: {fn: func(c *CPU) { addr := c.zeroPageX(); c.slo(addr); c.cycles += 6 }, Name: "*SLO Zero Page, X", Bytes: 2},
	0x1B: {fn: func(c *CPU) { addr, _ := c.absoluteY(); c.slo(addr); c.cycles += 7 }, Name: "*SLO Absolute, Y", Bytes: 3},
	0x1F: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.slo(addr); c.cycles += 7 }, Name: "*SLO Absolute, X", Bytes: 3},
	// *SRE
	0x43: {fn: func(c *CPU) { addr := c.indirectX(); c.sre(addr); c.cycles += 8 }, Name: "*SRE (Indirect, X)", Bytes: 2},
	0x47: {fn: func(c *CPU) { addr := c.zeroPage(); c.sre(addr); c.cycles += 5 }, Name: "*SRE Zero Page", Bytes: 2},
	0x4F: {fn: func(c *CPU) { addr := c.absolute(); c.sre(addr); c.cycles += 6 }, Name: "*SRE Absolute", Bytes: 3},
	0x53: {fn: func(c *CPU) { addr, _ := c.indirectY(); c.sre(addr); c.cycles += 8 }, Name: "*SRE (Indirect), Y", Bytes: 2},
	0x57: {fn: func(c *CPU) { addr := c.zeroPageX(); c.sre(addr); c.cycles += 6 }, Name: "*SRE Zero Page, X", Bytes: 2},
	0x5B: {fn: func(c *CPU) { addr, _ := c.absoluteY(); c.sre(addr); c.cycles += 7 }, Name: "*SRE Absolute, Y", Bytes: 3},
	0x5F: {fn: func(c *CPU) { addr, _ := c.absoluteX(); c.sre(addr); c.cycles += 7 }, Name: "*SRE Absolute, X", Bytes: 3},
	// *TAS
	0x9B: {fn: func(c *CPU) { addr, _ := c.absoluteY(); c.tas(addr); c.cycles += 5 }, Name: "*TAS Absolute, Y", Bytes: 3},
	// *XAA
	0x8B: {fn: func(c *CPU) { addr := c.immediate(); c.xaa(addr); c.cycles += 2 }, Name: "*XAA #Immediate", Bytes: 2},
}

func (c *CPU) calcBranchInstrCycles(isBranched bool, oops int) int {
	cycles := 2
	if isBranched {
		cycles += 1
		cycles += oops
	}
	return cycles
}
