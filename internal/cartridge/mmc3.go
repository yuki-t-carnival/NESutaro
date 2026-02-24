package cartridge

type MMC3 struct {
	prgROM               []byte
	chrROM               []byte
	hasPrgRAM            bool
	hasChrRAM            bool
	isVerticallyMirrored bool
	header               *INESHeader

	nextWrite    int
	r            [8]int
	prgBankMode  int
	chrInversion int
}

func NewMMC3(h *INESHeader, rom []byte) *MMC3 {
	mmc3 := &MMC3{}
	mmc3.header = h

	prgSize := 0x4000 * h.TotalPRGROMUnits
	chrSize := 0x2000 * h.TotalCHRROMUnits

	mmc3.prgROM = make([]byte, 0x4000*h.TotalPRGROMUnits)
	mmc3.chrROM = make([]byte, 0x2000*h.TotalCHRROMUnits)
	copy(mmc3.prgROM[:], rom[:prgSize])
	copy(mmc3.chrROM[:], rom[prgSize:prgSize+chrSize])

	return mmc3
}

func (m *MMC3) ReadPRGROM(addr uint16) byte {
	if m.prgBankMode == 0 {
		switch {
		case 0x8000 <= addr && addr <= 0x9FFF:
			return m.prgROM[0x2000*m.r[6]+int(addr-0x8000)]
		case 0xA000 <= addr && addr <= 0xBFFF:
			return m.prgROM[0x2000*m.r[7]+int(addr-0xA000)]
		case 0xC000 <= addr && addr <= 0xDFFF:
			return m.prgROM[0x2000*(m.header.TotalPRGROMUnits*2-2)+int(addr-0xC000)]
		case 0xE000 <= addr:
			return m.prgROM[0x2000*(m.header.TotalPRGROMUnits*2-1)+int(addr-0xE000)]
		default:
			return 0xFF
		}
	} else {
		switch {
		case 0x8000 <= addr && addr <= 0x9FFF:
			return m.prgROM[0x2000*(m.header.TotalPRGROMUnits*2-2)+int(addr-0x8000)]
		case 0xA000 <= addr && addr <= 0xBFFF:
			return m.prgROM[0x2000*m.r[7]+int(addr-0xA000)]
		case 0xC000 <= addr && addr <= 0xDFFF:
			return m.prgROM[0x2000*m.r[6]+int(addr-0xC000)]
		case 0xE000 <= addr:
			return m.prgROM[0x2000*(m.header.TotalPRGROMUnits*2-1)+int(addr-0xE000)]
		default:
			return 0xFF
		}
	}
}

func (m *MMC3) WritePRGROM(addr uint16, val byte) {
	switch {
	case 0x8000 <= addr && addr <= 0x9FFF:
		if addr&1 == 0 {
			m.nextWrite = int(val & 0x07)
			m.prgBankMode = int(val >> 6 & 1)
			m.chrInversion = int(val >> 7 & 1)
		} else {
			var bank int
			switch m.nextWrite {
			case 0, 1:
				bank = int(val & 0xFE)
			case 6, 7:
				bank = int(val & 0x3F)
			default:
				bank = int(val)
			}
			m.r[m.nextWrite] = bank
			/* for i := 0; i < 8; i++ {
				fmt.Printf("r%d:%03d ", i, m.r[i])
			}
			fmt.Println() */
		}
	}
	//fmt.Printf("MMC3 select bank = %d\n", m.bank)
}

func (m *MMC3) ReadPRGRAM(addr uint16) byte {
	return 0xFF
}

func (m *MMC3) WritePRGRAM(addr uint16, val byte) {
}

func (m *MMC3) ReadCHRROM(addr uint16) byte {
	if m.chrInversion == 0 {
		switch {
		case addr <= 0x07FF:
			return m.chrROM[m.r[0]*0x400+int(addr)]
		case 0x0800 <= addr && addr <= 0x0FFF:
			return m.chrROM[m.r[1]*0x400+int(addr-0x0800)]
		case 0x1000 <= addr && addr <= 0x13FF:
			return m.chrROM[m.r[2]*0x400+int(addr-0x1000)]
		case 0x1400 <= addr && addr <= 0x17FF:
			return m.chrROM[m.r[3]*0x400+int(addr-0x1400)]
		case 0x1800 <= addr && addr <= 0x1BFF:
			return m.chrROM[m.r[4]*0x400+int(addr-0x1800)]
		case 0x1C00 <= addr && addr <= 0x1FFF:
			return m.chrROM[m.r[5]*0x400+int(addr-0x1C00)]
		}
	} else {
		switch {
		case addr <= 0x03FF:
			return m.chrROM[m.r[2]*0x400+int(addr)]
		case 0x0400 <= addr && addr <= 0x07FF:
			return m.chrROM[m.r[3]*0x400+int(addr-0x0400)]
		case 0x0800 <= addr && addr <= 0x0BFF:
			return m.chrROM[m.r[4]*0x400+int(addr-0x0800)]
		case 0x0C00 <= addr && addr <= 0x0FFF:
			return m.chrROM[m.r[5]*0x400+int(addr-0x0C00)]
		case 0x1000 <= addr && addr <= 0x17FF:
			return m.chrROM[m.r[0]*0x400+int(addr-0x1000)]
		case 0x1800 <= addr && addr <= 0x1FFF:
			return m.chrROM[m.r[1]*0x400+int(addr-0x1800)]
		}
	}
	return 0xFF
}

func (m *MMC3) WriteCHRROM(addr uint16, val byte) {
}

func (m *MMC3) GetHeaderInfo() []string {
	var strs []string
	return strs
}

func (m *MMC3) GetSaveData() []byte {
	return []byte{}
}

func (m *MMC3) IsVerticallyMirrored() bool {
	return m.isVerticallyMirrored
}
