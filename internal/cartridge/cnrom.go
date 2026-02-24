package cartridge

type CNROM struct {
	prgROM               [0x8000]byte
	chrROM               [0x8000]byte
	hasPrgRAM            bool
	hasChrRAM            bool
	isVerticallyMirrored bool
	bank                 int
}

func NewCNROM(h *INESHeader, rom []byte) *CNROM {
	cnrom := &CNROM{}

	chrStart := 0x4000 * h.TotalPRGROMUnits
	chrEnd := chrStart + 0x2000*h.TotalCHRROMUnits
	copy(cnrom.prgROM[:], rom[:chrStart])
	copy(cnrom.chrROM[:], rom[chrStart:chrEnd])

	return cnrom
}

func (c *CNROM) ReadPRGROM(addr uint16) byte {
	return c.prgROM[addr-0x8000]
}

func (c *CNROM) WritePRGROM(addr uint16, val byte) {
	c.bank = int(val & 0x03)
	//fmt.Printf("CNROM select bank = %d\n", c.bank)
}

func (c *CNROM) ReadPRGRAM(addr uint16) byte {
	return 0xFF
}

func (c *CNROM) WritePRGRAM(addr uint16, val byte) {
}

func (c *CNROM) ReadCHRROM(addr uint16) byte {
	return c.chrROM[c.bank*0x2000+int(addr)]

}

func (c *CNROM) WriteCHRROM(addr uint16, val byte) {
}

func (c *CNROM) GetHeaderInfo() []string {
	var strs []string
	return strs
}

func (c *CNROM) GetSaveData() []byte {
	return []byte{}
}

func (c *CNROM) IsVerticallyMirrored() bool {
	return c.isVerticallyMirrored
}
