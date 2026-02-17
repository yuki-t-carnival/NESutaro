package bus

import (
	"fmt"
	"nesutaro/internal/cartridge"
)

type Bus struct {
	Cart *cartridge.Cartridge
	vram [0x800]byte
	pram [0x20]byte
}

const ()

func NewBus(cart *cartridge.Cartridge) *Bus {
	bus := &Bus{
		Cart: cart,
	}
	return bus
}

func (b *Bus) Read(addr uint16) byte {
	switch {
	case addr <= 0x1FFF:
		return b.Cart.ReadCHRROM(addr)

	case 0x2000 <= addr && addr <= 0x23FF:
		return b.vram[0x000+addr&0x3FF]
	case 0x2400 <= addr && addr <= 0x27FF:
		if b.Cart.IsVerticallyMirrored() {
			return b.vram[0x400+addr&0x3FF]
		} else {
			return b.vram[0x000+addr&0x3FF]
		}
	case 0x2800 <= addr && addr <= 0x2BFF:
		if b.Cart.IsVerticallyMirrored() {
			return b.vram[0x000+addr&0x3FF]
		} else {
			return b.vram[0x400+addr&0x3FF]
		}
	case 0x2C00 <= addr && addr <= 0x2FFF:
		return b.vram[0x400+addr&0x3FF]
	case 0x3000 <= addr && addr <= 0x3EFF:
		return b.Read(addr - 0x1000)
	case 0x3F00 <= addr && addr <= 0x3FFF:
		pAddr := addr & 0x1F
		if pAddr == 0x10 || pAddr == 0x14 || pAddr == 0x18 || pAddr == 0x1C {
			pAddr -= 0x10
		}
		return b.pram[pAddr]

	default:
		return 0xFF
	}
}

func (b *Bus) Write(addr uint16, val byte) {
	switch {
	case addr <= 0x1FFF:
		b.Cart.WriteCHRROM(addr, val)

	case 0x2000 <= addr && addr <= 0x23FF:
		b.vram[0x000+addr&0x3FF] = val
	case 0x2400 <= addr && addr <= 0x27FF:
		if b.Cart.IsVerticallyMirrored() {
			b.vram[0x400+addr&0x3FF] = val
		} else {
			b.vram[0x000+addr&0x3FF] = val
		}
	case 0x2800 <= addr && addr <= 0x2BFF:
		if b.Cart.IsVerticallyMirrored() {
			b.vram[0x000+addr&0x3FF] = val
		} else {
			b.vram[0x400+addr&0x3FF] = val
		}
	case 0x2C00 <= addr && addr <= 0x2FFF:
		b.vram[0x400+addr&0x3FF] = val
	case 0x3000 <= addr && addr <= 0x3EFF:
		b.Write(addr-0x1000, val)
	case 0x3F00 <= addr && addr <= 0x3FFF:
		pAddr := addr & 0x1F
		if pAddr == 0x10 || pAddr == 0x14 || pAddr == 0x18 || pAddr == 0x1C {
			pAddr -= 0x10
		}
		b.pram[pAddr] = val
	}
}

func (b *Bus) VRAMLog() {
	for i, v := range b.vram {
		fmt.Printf("vram[%04X]=%02X\n", i, v)
	}
}
