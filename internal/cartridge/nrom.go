package cartridge

type NROM struct {
	prgROM               [0x8000]byte
	chrROM               [0x2000]byte
	hasPrgRAM            bool
	hasChrRAM            bool
	isVerticallyMirrored bool
}

func NewNROM(h *INESHeader, rom []byte) *NROM {
	nrom := &NROM{}

	copy(nrom.prgROM[:], rom[:h.TotalPRGROMUnits*0x4000])

	if h.TotalCHRROMUnits >= 1 {
		copy(nrom.chrROM[:], rom[h.TotalPRGROMUnits*0x4000:])
	}
	return nrom
}

func (n *NROM) ReadPRGROM(addr uint16) byte {
	return n.prgROM[addr-0x8000]
}

func (n *NROM) WritePRGROM(addr uint16, val byte) {
}

func (n *NROM) ReadPRGRAM(addr uint16) byte {
	return 0xFF
}

func (n *NROM) WritePRGRAM(addr uint16, val byte) {
}

func (n *NROM) ReadCHRROM(addr uint16) byte {
	return n.chrROM[addr]

}

func (n *NROM) WriteCHRROM(addr uint16, val byte) {
}

func (n *NROM) GetHeaderInfo() []string {
	var strs []string
	return strs
}

func (n *NROM) GetSaveData() []byte {
	return []byte{}
}

func (n *NROM) IsVerticallyMirrored() bool {
	return n.isVerticallyMirrored
}
