package bus

import (
	"nesutaro/internal/cartridge"
	"nesutaro/internal/joypad"
	"nesutaro/internal/ppu"
)

type Bus struct {
	Cart      *cartridge.Cartridge
	PPU       *ppu.PPU
	Joypad    *joypad.Joypad
	wram      [0x800]byte
	reg0x4010 byte
	reg0x4017 byte
	HasIRQ    bool
}

const ()

func NewBus(cart *cartridge.Cartridge, p *ppu.PPU, j *joypad.Joypad) *Bus {
	bus := &Bus{
		Cart:   cart,
		PPU:    p,
		Joypad: j,
	}

	return bus
}

func (b *Bus) Read(addr uint16) byte {
	switch {
	case addr <= 0x1FFF:
		return b.wram[addr&0x07FF]

	case addr == 0x2002:
		return b.PPU.ReadPPUSTATUS()
	case addr == 0x2004:
		return b.PPU.ReadOAMDATA()
	case addr == 0x2007:
		return b.PPU.ReadPPUDATA()
	case 0x2008 <= addr && addr <= 0x3FFF:
		return b.Read(0x2000 + addr&7)

	case addr == 0x4015:
		if b.reg0x4017>>6&0x03 == 0 {
			b.HasIRQ = true
		}
		return b.reg0x4017
	case addr == 0x4016:
		return b.Joypad.Read4016()

	case 0x6000 <= addr && addr <= 0x7FFF:
		return b.Cart.ReadPRGRAM(addr)
	case 0x8000 <= addr:
		if b.Cart.Header.TotalPRGROMUnits == 1 {
			return b.Cart.ReadPRGROM(0x8000 + addr&0x3FFF)
		} else {
			return b.Cart.ReadPRGROM(addr)
		}
	default:
		return 0xFF
	}
}

func (b *Bus) Write(addr uint16, val byte) {
	switch {
	case addr <= 0x1FFF:
		b.wram[addr&0x07FF] = val

	case addr == 0x2000:
		b.PPU.WritePPUCTRL(val)
	case addr == 0x2001:
		b.PPU.WritePPUMASK(val)
	case addr == 0x2003:
		b.PPU.WriteOAMADDR(val)
	case addr == 0x2004:
		b.PPU.WriteOAMDATA(val)
	case addr == 0x2005:
		b.PPU.WritePPUSCROLL(val)
	case addr == 0x2006:
		b.PPU.WritePPUADDR(val)
	case addr == 0x2007:
		b.PPU.WritePPUDATA(val)
	case 0x2008 <= addr && addr <= 0x3FFF:
		b.Write(0x2000+addr&7, val)

	case addr == 0x4010:
		if b.reg0x4010>>7&1 == 0 && val>>7&1 == 1 {
			b.HasIRQ = true
		}
		b.reg0x4010 = val

	case addr == 0x4014:
		var dmaData [256]byte
		for i := 0; i < 256; i++ {
			dmaData[i] = b.Read(uint16(val)<<8 + uint16(i))
		}
		b.PPU.WriteOAMDMA(&dmaData)

	case addr == 0x4015:
		if b.reg0x4010>>7&1 == 1 {
			b.HasIRQ = true
		}

	case addr == 0x4016:
		b.Joypad.Write4016(val)

	case addr == 0x4017:
		b.reg0x4017 = val

	case 0x6000 <= addr && addr <= 0x7FFF:
		b.Cart.WritePRGRAM(addr, val)
	case 0x8000 <= addr:
		if b.Cart.Header.TotalPRGROMUnits == 1 {
			b.Cart.WritePRGROM(0x8000+addr&0x3FFF, val)
		} else {
			b.Cart.WritePRGROM(addr, val)
		}
	}
}
