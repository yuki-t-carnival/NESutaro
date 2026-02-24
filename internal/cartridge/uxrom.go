package cartridge

type UxROM struct {
	prgROM               []byte
	chrROM               [0x2000]byte
	hasPrgRAM            bool
	hasChrRAM            bool
	isVerticallyMirrored bool
	bank                 int
	header               *INESHeader
}

func NewUxROM(h *INESHeader, rom []byte) *UxROM {
	uxrom := &UxROM{}
	uxrom.header = h

	uxrom.prgROM = make([]byte, 0x4000*h.TotalPRGROMUnits)
	copy(uxrom.prgROM[:], rom[:0x4000*h.TotalPRGROMUnits])
	copy(uxrom.chrROM[:], rom[0x4000*h.TotalPRGROMUnits:])

	return uxrom
}

func (u *UxROM) ReadPRGROM(addr uint16) byte {
	switch {
	case 0x8000 <= addr && addr <= 0xBFFF:
		return u.prgROM[u.bank*0x4000+int(addr-0x8000)]
	case 0xC000 <= addr:
		return u.prgROM[(u.header.TotalPRGROMUnits-1)*0x4000+int(addr-0xC000)]
	default:
		return 0xFF
	}
}

func (u *UxROM) WritePRGROM(addr uint16, val byte) {
	u.bank = int(val & 0x0F)
	//fmt.Printf("UxROM select bank = %d\n", u.bank)
}

func (u *UxROM) ReadPRGRAM(addr uint16) byte {
	return 0xFF
}

func (u *UxROM) WritePRGRAM(addr uint16, val byte) {
}

func (u *UxROM) ReadCHRROM(addr uint16) byte {
	return u.chrROM[addr]

}

func (u *UxROM) WriteCHRROM(addr uint16, val byte) {
}

func (u *UxROM) GetHeaderInfo() []string {
	var strs []string
	return strs
}

func (u *UxROM) GetSaveData() []byte {
	return []byte{}
}

func (u *UxROM) IsVerticallyMirrored() bool {
	return u.isVerticallyMirrored
}
