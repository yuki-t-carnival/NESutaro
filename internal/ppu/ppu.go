package ppu

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"nesutaro/internal/ppu/bus"
	"os"
)

const (
	VblankFlag      byte = 1 << 7
	VblankNMIEnable byte = 1 << 7
)

var xFlipLUT [256]byte

type PPU struct {
	Bus      *bus.Bus
	viewport [256 * 240]int
	screen   [2]*image.RGBA
	front    int

	oam [256]byte

	cycles int

	v         uint16
	t         uint16
	x         uint8
	w         bool
	ppuctrl   byte
	ppumask   byte
	ppustatus byte
	oamaddr   byte

	hasNMI bool

	ly int

	nesPal [64]color.RGBA
}

func NewPPU(b *bus.Bus) *PPU {
	palFilePath := "nes.pal"
	palFile, err := os.ReadFile(palFilePath)
	if err != nil {
		log.Fatal(err)
	}
	p := &PPU{
		Bus: b,
	}
	p.screen[0] = image.NewRGBA(image.Rect(0, 0, 256, 240))
	p.screen[1] = image.NewRGBA(image.Rect(0, 0, 256, 240))

	for i := 0; i < 64; i++ {
		p.nesPal[i] = color.RGBA{palFile[i*3+0], palFile[i*3+1], palFile[i*3+2], 255}
	}
	return p
}

func (p *PPU) Step(cpuCycles int) {
	p.cycles += cpuCycles * 3
	for p.cycles >= 341 {
		switch {
		case p.ly < 240: // HBlank: 0 ~ 239
			isBGEnabled := p.ppumask>>3&1 == 1
			if isBGEnabled {
				p.drawBGLine(p.ly)
			}
			if p.ly > 0 {
				isSpriteEnabled := p.ppumask>>4&1 == 1
				if isSpriteEnabled {
					p.drawSpritesLine(p.ly - 1)
				}
			}
		case p.ly == 240: // VBlank: 240 ~ 259
			p.front ^= 1
			backColor := p.nesPal[p.Bus.Read(0x3F00)]
			for y := 0; y < 240; y++ {
				for x := 0; x < 256; x++ {
					p.screen[p.front^1].SetRGBA(x, y, backColor)
				}
			}
			p.ppustatus |= VblankFlag
			if p.ppuctrl&VblankNMIEnable != 0 {
				p.hasNMI = true
			}
		case p.ly == 260: // Post(pre) render scanline: 260, 261
			p.ppustatus &^= VblankFlag
		}

		p.cycles -= 341
		p.ly += 1
		if p.ly >= 262 {
			p.ly = 0
		}

	}
}

func (p *PPU) drawBGLine(ly int) {
	currentCoarseY := ((p.getYScroll() + ly) % 480) / 8
	currentFineY := ((p.getYScroll() + ly) % 480) % 8
	prevCoarseX := -1

	var tileLine [2]byte
	for i := 0; i < 256; i++ {
		currentCoarseX := ((p.getXScroll() + i) % 512) / 8
		currentFineX := ((p.getXScroll() + i) % 512) % 8
		if currentCoarseX != prevCoarseX {
			tileIndex := p.getBGTileIndex(currentCoarseX, currentCoarseY)
			tileLine = p.getBGTileLine(tileIndex, currentFineY)
		}
		raw := p.getTileLinePixel(&tileLine, currentFineX)
		if raw != 0 {
			rgba := p.resolveBGPalette(currentCoarseX, currentCoarseY, raw)
			p.setPixel(i, p.ly, rgba)
		}
		prevCoarseX = currentCoarseX
	}
}

func (p *PPU) drawSpritesLine(ly int) {
	size := 8
	is16Size := p.ppuctrl>>5&1 == 1
	if is16Size {
		size = 16
	}
	var tileLine [2]byte
	for i := 63; i >= 0; i-- {
		posY := int(p.oam[i*4+0])
		tileIdx := p.oam[i*4+1]
		attr := int(p.oam[i*4+2])
		posX := int(p.oam[i*4+3])

		isXFlip := attr>>6&1 == 1
		isYFlip := attr>>7&1 == 1

		if posY <= ly && ly < posY+size {
			fineY := ly - posY
			tileY := fineY
			if isYFlip {
				tileY = size - 1 - tileY
			}
			tileLine = p.getSpriteTileLine(tileIdx, tileY)
			for j := 0; j < 8; j++ {
				if posX+j < 256 {
					tileX := j
					if isXFlip {
						tileX = 7 - tileX
					}
					raw := p.getTileLinePixel(&tileLine, tileX)
					palVal := attr & 0x03
					if raw != 0 {
						rgba := p.resolveSpritesPalette(palVal, raw)
						p.setPixel(posX+j, ly+1, rgba)
					}
					p.ppustatus |= 1 << 6
				}
			}
		}
	}
}

func (p *PPU) resolveSpritesPalette(palVal int, raw int) color.RGBA {
	colorVal := p.Bus.Read(0x3F10 + uint16(palVal)*4 + uint16(raw))
	return p.nesPal[colorVal&0x3F]
}

func (p *PPU) getBGTileIndex(coarseX, coarseY int) byte {
	mapCoarseX := coarseX % 32
	mapCoarseY := coarseY % 30
	nameStart := 0x2000
	nameStart += 0x0400 * (coarseX / 32)
	nameStart += 0x0800 * (coarseY / 30)
	nameAddr := nameStart + mapCoarseY*32 + mapCoarseX
	index := p.Bus.Read(uint16(nameAddr))
	return index
}

func (p *PPU) resolveBGPalette(coarseX, coarseY, raw int) color.RGBA {
	mapBlockX := coarseX % 32 / 4
	mapBlockY := coarseY % 30 / 4
	modX := coarseX % 32 % 4
	modY := coarseY % 30 % 4
	base := 0x23C0
	base += 0x0400 * (coarseX / 32)
	base += 0x0800 * (coarseY / 30)

	attrAddr := base + mapBlockY*8 + mapBlockX

	shift := 2 * (modX / 2)
	shift += 4 * (modY / 2)

	palVal := uint16(p.Bus.Read(uint16(attrAddr)) >> shift & 0x03)

	colorVal := p.Bus.Read(0x3F00 + palVal*4 + uint16(raw))

	return p.nesPal[colorVal&0x3F]
}

func (p *PPU) getBGTileLine(tileIndex byte, fineY int) [2]byte {
	var data [2]byte
	base := uint16(0x0000)
	if p.ppuctrl>>4&1 == 1 {
		base = 0x1000
	}
	data[0] = p.Bus.Read(base + uint16(tileIndex)*16 + uint16(fineY))
	data[1] = p.Bus.Read(base + uint16(tileIndex)*16 + 8 + uint16(fineY))
	return data
}

func (p *PPU) getTileLinePixel(tileLine *[2]byte, fineX int) int {
	shift := 7 - fineX
	lo := tileLine[0] & (1 << shift) >> shift
	hi := tileLine[1] & (1 << shift) >> shift
	px := int(hi<<1 | lo)
	return px & 0x03
}

func (p *PPU) getSpriteTileLine(tileIndex byte, fineY int) [2]byte {
	idx := uint16(tileIndex)
	var start uint16

	is16Size := p.ppuctrl>>5&1 == 1
	if is16Size {
		start = idx&1*0x1000 + idx&0xFE*0x10
	} else {
		spriteAddrNum := uint16(p.ppuctrl >> 3 & 1)
		start = spriteAddrNum*0x1000 + idx*0x10
	}

	var data [2]byte
	data[0] = p.Bus.Read(start + uint16(fineY))
	data[1] = p.Bus.Read(start + 8 + uint16(fineY))
	return data
}

func (p *PPU) setPixel(x, y int, rgba color.RGBA) {
	p.screen[p.front^1].SetRGBA(x, y, rgba)
}

func (p *PPU) getXScroll() int {
	xBit8 := int(p.t) >> 10 & 1
	coarseX := int(p.t & 0x001F)
	fineX := int(p.x & 0x07)
	return xBit8<<8 | coarseX<<3 | fineX
}

func (p *PPU) getYScroll() int {
	yBit8 := int(p.t) >> 11 & 1
	coarseY := int(p.t & 0x03E0 >> 5)
	fineY := int(p.t & 0x7000 >> 12)
	return yBit8<<8 | coarseY<<3 | fineY
}

// Get Viewport pixels converted from colorIndex to RGBA
func (p *PPU) GetGameScreen() *image.RGBA {
	return p.screen[p.front]
}

func (p *PPU) ReadOAM(addr uint16) byte {
	return p.oam[addr]
}
func (p *PPU) WriteOAM(addr uint16, val byte) {
	p.oam[addr] = val
}

func (p *PPU) WritePPUCTRL(val byte) {
	p.ppuctrl = val & 0xFC
	p.t = p.t&0x73FF | uint16(val)&0x03<<10
}

func (p *PPU) WritePPUMASK(val byte) {
	p.ppumask = val
}

func (p *PPU) ReadPPUSTATUS() byte {
	p.w = false
	val := p.ppustatus
	p.ppustatus &^= VblankFlag
	return val
}

func (p *PPU) WriteOAMADDR(val byte) {
	p.oamaddr = val
}

func (p *PPU) ReadOAMDATA() byte {
	return p.oam[p.oamaddr]
}

func (p *PPU) WriteOAMDATA(val byte) {
	p.oam[p.oamaddr] = val
	p.oamaddr += 1
}

func (p *PPU) WritePPUSCROLL(val byte) {
	isFirstWriting := !p.w
	if isFirstWriting {
		p.x = val & 0x07
		coarseX := uint16(val >> 3)
		p.t = p.t&0x7FE0 | coarseX
	} else {
		fineY := uint16(val & 0x07)
		coarseY := uint16(val >> 3)
		p.t = p.t&0x0C1F | fineY<<12 | coarseY<<5
	}
	p.w = !p.w
}

func (p *PPU) WritePPUADDR(val byte) {
	isFirstWriting := !p.w
	if isFirstWriting {
		// Bit 14 is forced 0 when writting the PPUADDR high byte.
		p.t = uint16(val&0x3F)<<8 | p.t&0x00FF
	} else {
		p.t = p.t&0x7F00 | uint16(val)
		p.v = p.t
	}
	p.w = !p.w
}

func (p *PPU) ReadPPUDATA() byte {
	val := p.Bus.Read(p.v)
	if p.ppuctrl>>2&1 == 0 {
		p.v += 1
	} else {
		p.v += 32
	}
	return val
}

func (p *PPU) WritePPUDATA(val byte) {
	p.Bus.Write(p.v, val)
	if p.ppuctrl>>2&1 == 0 {
		p.v += 1
	} else {
		p.v += 32
	}
}

func (p *PPU) WriteOAMDMA(data *[0x100]byte) {
	p.oam = *data
}

func (p *PPU) HasNMI() bool {
	return p.hasNMI
}

func (p *PPU) DisableNMI() {
	p.hasNMI = false
}

func (p *PPU) OAMLog() {
	for i, v := range p.oam {
		fmt.Printf("p.oam[%d]=%02X\n", i, v)
	}
}
