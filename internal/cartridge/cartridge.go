package cartridge

import (
	"fmt"
)

type Cartridge struct {
	Mapper Mapper
	Header *INESHeader
}

func NewCartridge(rom /* , sav */ []byte) *Cartridge {

	cart := &Cartridge{
		Header: NewINESHeader(rom[:0x10]),
	}
	switch cart.Header.MapperNum {
	case 0:
		cart.Mapper = NewNROM(cart.Header, rom[0x10:])
	case 2:
		cart.Mapper = NewUxROM(cart.Header, rom[0x10:])
	case 3:
		cart.Mapper = NewCNROM(cart.Header, rom[0x10:])
	case 4:
		cart.Mapper = NewMMC3(cart.Header, rom[0x10:])
	}

	return cart
}

type Mapper interface {
	ReadPRGROM(addr uint16) byte
	WritePRGROM(addr uint16, val byte)
	ReadPRGRAM(addr uint16) byte
	WritePRGRAM(addr uint16, val byte)
	ReadCHRROM(addr uint16) byte
	WriteCHRROM(addr uint16, val byte)
	GetSaveData() []byte
}

type INESHeader struct {
	IsVerticallyMirrored bool
	TotalPRGROMUnits     int
	TotalCHRROMUnits     int
	MapperNum            int
	PRGRAMBytes          int
	PRGNVRAMBytes        int
}

func NewINESHeader(rom []byte) *INESHeader {
	h := &INESHeader{}
	h.IsVerticallyMirrored = rom[6]&0x01 != 0
	h.TotalPRGROMUnits = int(rom[4])
	h.TotalCHRROMUnits = int(rom[5])
	h.MapperNum = int(rom[6] >> 4)
	if int(rom[10]&0x0F) > 0 {
		h.PRGRAMBytes = 64 << int(rom[10]&0x0F)
	}
	if int(rom[10]>>4) > 0 {
		h.PRGNVRAMBytes = 64 << int(rom[10]>>4)
	}

	fmt.Printf("Mapper number = %d\n", h.MapperNum)
	fmt.Printf("Total PRG ROM units = %d\n", h.TotalPRGROMUnits)
	fmt.Printf("Total CHR ROM units = %d\n", h.TotalCHRROMUnits)
	fmt.Printf("PRG RAM Bytes = %d\n", h.PRGRAMBytes)
	fmt.Printf("PRG NVRAM Bytes = %d\n", h.PRGNVRAMBytes)
	return h
}

func (c *Cartridge) ReadPRGROM(addr uint16) byte {
	return c.Mapper.ReadPRGROM(addr)
}

func (c *Cartridge) WritePRGROM(addr uint16, val byte) {
	c.Mapper.WritePRGROM(addr, val)
}

func (c *Cartridge) ReadPRGRAM(addr uint16) byte {
	return c.Mapper.ReadPRGRAM(addr)
}

func (c *Cartridge) WritePRGRAM(addr uint16, val byte) {
	c.Mapper.WritePRGRAM(addr, val)
}

func (c *Cartridge) ReadCHRROM(addr uint16) byte {
	return c.Mapper.ReadCHRROM(addr)

}

func (c *Cartridge) WriteCHRROM(addr uint16, val byte) {
	c.Mapper.WriteCHRROM(addr, val)
}

func (c *Cartridge) GetHeaderInfo() []string {
	var strs []string
	return strs
}

func (c *Cartridge) GetSaveData() []byte {
	return []byte{}
}

func (c *Cartridge) IsVerticallyMirrored() bool {
	return c.Header.IsVerticallyMirrored
}
